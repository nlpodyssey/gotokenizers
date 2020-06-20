// Copyright (c) 2020, NLP Odyssey Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package bpemodel

import "container/heap"

type WordMerge struct {
	MergeValue
	Pos int
}

type WordMergeHeap []WordMerge

var _ heap.Interface = &WordMergeHeap{}

func (h *WordMergeHeap) Len() int {
	return len(*h)
}

func (h *WordMergeHeap) Less(i, j int) bool {
	// By manually implementing this, we make the heap a min-heap,
	// ordered first on the rank, and on the pos otherwise.
	if (*h)[i].Rank == (*h)[j].Rank {
		return (*h)[i].Pos > (*h)[j].Pos
	}
	return (*h)[i].Rank > (*h)[j].Rank
}

func (h *WordMergeHeap) Swap(i, j int) {
	(*h)[i], (*h)[j] = (*h)[j], (*h)[i]
}

func (h *WordMergeHeap) Push(x interface{}) {
	*h = append(*h, x.(WordMerge))
}

func (h *WordMergeHeap) Pop() interface{} {
	lastIndex := len(*h) - 1
	x := (*h)[lastIndex]
	*h = (*h)[0:lastIndex]
	return x
}
