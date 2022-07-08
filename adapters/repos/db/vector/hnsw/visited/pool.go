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

import (
	"sync"
)

type Pool struct {
	sync.Mutex
	listSetSize int
	listSets    []ListSet
}

func NewPool(size int, listSetSize int) *Pool {
	p := &Pool{
		listSetSize: listSetSize,
		listSets:    make([]ListSet, size, size+32), // make enough room
	}

	for i := 0; i < size; i++ {
		p.listSets[i] = NewList(listSetSize)
	}

	return p
}

func (p *Pool) Borrow() ListSet {
	p.Lock()
	defer p.Unlock()

	if n := len(p.listSets); n > 0 {
		l := p.listSets[n-1]
		p.listSets[n-1].Free() // prevent memory leak
		p.listSets = p.listSets[:n-1]
		return l
	}

	return NewList(p.listSetSize)
}

func (p *Pool) Return(l ListSet) {
	p.Lock()
	defer p.Unlock()

	if l.Len() > p.listSetSize*5/4 { // 5/4 could be tuned
		return // discard this list, it does not match our current criteria
	}

	l.Reset()
	p.listSets = append(p.listSets, l)
}

func (p *Pool) Destroy() {
	p.Lock()
	defer p.Unlock()
	for i := range p.listSets {
		p.listSets[i].Free()
	}

	p.listSets = nil
}
