package lsmkv

import (
	"context"
	"io/fs"
	"path/filepath"

	"github.com/pkg/errors"
)

// PauseCompaction waits for all ongoing compactions to finish,
// then makes sure that no new compaction can be started.
//
// This is a preparatory stage for taking snapshots.
//
// A timeout should be specified for the input context as some
// compactions are long-running, in which case it may be better
// to fail the backup attempt and retry later, than to block
// indefinitely.
func (b *Bucket) PauseCompaction(ctx context.Context) error {
	compactionHalted := make(chan struct{})

	go func() {
		b.disk.stopCompactionCycle <- struct{}{}
		compactionHalted <- struct{}{}
	}()

	select {
	case <-ctx.Done():
		// resume the compaction cycle, as the
		// context deadline was exceeded
		defer b.disk.initCompactionCycle(DefaultCompactionInterval)
		return errors.Wrap(ctx.Err(), "long-running compaction in progress")
	case <-compactionHalted:
		return nil
	}
}

// FlushMemtable flushes any active memtable and returns only once the memtable
// has been fully flushed and a stable state on disk has been reached.
//
// This is a preparatory stage for taking snapshots.
//
// A timeout should be specified for the input context as some
// flushes are long-running, in which case it may be better
// to fail the backup attempt and retry later, than to block
// indefinitely.
func (b *Bucket) FlushMemtable(ctx context.Context) error {
	defer b.initFlushCycle()
	flushed := make(chan struct{})

	go func() {
		b.stopFlushCycle <- struct{}{}
		flushed <- struct{}{}
	}()

	select {
	case <-ctx.Done():
		return errors.Wrap(ctx.Err(), "long-running memtable flush in progress")
	case <-flushed:
		if b.active == nil && b.flushing == nil {
			return nil
		}

		return b.FlushAndSwitch()
	}
}

// ListFiles lists all files that currently exist in the Bucket. The files are only
// in a stable state if the memtable is empty, and if compactions are paused. If one
// of those conditions is not given, it errors
func (b *Bucket) ListFiles(ctx context.Context) ([]string, error) {
	var (
		bucketRoot = b.disk.dir
		files      []string
	)

	err := filepath.WalkDir(bucketRoot, func(path string, d fs.DirEntry, err error) error {
		if d.IsDir() {
			return nil
		}
		files = append(files, path)
		return nil
	})
	if err != nil {
		return nil, errors.Errorf("failed to list files for bucket: %s", err)
	}

	return files, nil
}

// ResumeCompaction starts the compaction cycle again.
// It errors if compactions were not paused
func (b *Bucket) ResumeCompaction(ctx context.Context) error {
	b.disk.initCompactionCycle(DefaultCompactionInterval)
	return nil
}
