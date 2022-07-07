//                           _       _
// __      _____  __ ___   ___  __ _| |_ ___
// \ \ /\ / / _ \/ _` \ \ / / |/ _` | __/ _ \
//  \ V  V /  __/ (_| |\ V /| | (_| | ||  __/
//   \_/\_/ \___|\__,_| \_/ |_|\__,_|\__\___|
//
//  Copyright © 2016 - 2022 SeMI Technologies B.V. All rights reserved.
//
//  CONTACT: hello@semi.technology
//

package visited

// ListSet is a reusable list with very efficient resets. Inspired by the C++
// implementation in hnswlib it can be reset with zero memrory writes in the
// array by moving the match target instead of altering the list. Only after a
// version overflow do we need to actually reset
type ListSet struct {
	set []uint8
}

//  Len returns the length of 
func (l ListSet) Len() uint64 { return uint64(len(l.set)) - 1 }

// Free allocated slice. This list not resuable after this call
func (l *ListSet) Free() { l.set = nil }

func NewList(size int) ListSet {
	set := make([]uint8, size+1)
	// start at 1 since the initial value of the list is already 0, so we need
	// something to differentiate from that
	set[0] = 1 // version
	return ListSet{set: set}
}

func (l *ListSet) Visit(node uint64) {
	if node >= l.Len() {
		l.resize(node + 1024)
	}

	l.set[node+1] = l.set[0]
}

func (l *ListSet) Visited(node uint64) bool {
	return node < l.Len() && l.set[node+1] == l.set[0]
}

func (l *ListSet) resize(target uint64) {
	newStore := make([]uint8, target)
	copy(newStore, l.set)
	l.set = newStore
}

// Reset list only in case of an overflow.
func (l *ListSet) Reset() {
	l.set[0]++
	if l.set[0] == 0 { // if overflowed
		for i := range l.set {
			l.set[i] = 0
		}
		l.set[0] = 1
	}
}
