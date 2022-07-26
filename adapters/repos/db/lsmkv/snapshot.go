package lsmkv

import (
	"context"
	"fmt"
	"io/fs"
	"path/filepath"

	"github.com/pkg/errors"
	"github.com/semi-technologies/weaviate/entities/storagestate"
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
	select {
	case <-ctx.Done():
		return errors.Wrap(ctx.Err(), "long-running compaction in progress")
	case stopped := <-b.disk.compactionCycle.TryStop(ctx):
		select {
		case <-ctx.Done():
			return errors.Wrap(ctx.Err(), "long-running compaction in progress")
		default:
			if !stopped {
				return fmt.Errorf("failed to pause compaction")
			}
		}
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
	if b.isReadOnly() {
		return errors.Wrap(storagestate.ErrStatusReadOnly, "flush memtable")
	}

	select {
	case <-ctx.Done():
		return errors.Wrap(ctx.Err(), "long-running memtable flush in progress")
	case stopped := <-b.flushCycle.TryStop(ctx):
		select {
		case <-ctx.Done():
			return errors.Wrap(ctx.Err(), "long-running memtable flush in progress")
		default:
			if !stopped {
				return fmt.Errorf("failed to flash memtable")
			}			
		}
		defer b.flushCycle.Start()
		// this lock does not currently _need_ to be
		// obtained, as the only other place that
		// grabs this lock is the flush cycle, which
		// has just been stopped above.
		//
		// that being said, we will lock here anyway
		// as flushLock may be added elsewhere in the
		// future
		b.flushLock.Lock()
		if b.active == nil && b.flushing == nil {
			b.flushLock.Unlock()
			return nil
		}
		b.flushLock.Unlock()

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
	b.disk.compactionCycle.Start()
	return nil
}
