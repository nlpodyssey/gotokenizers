// Copyright (c) 2020, NLP Odyssey Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package bpemodel

import (
	"container/heap"
	"math/rand"
)

// Symbol is an abstract reference to a sequence of characters.
type Symbol struct {
	// Unique identifier, which implicitly refers to a sequence of characters.
	// For example, it might be the ID of a word in a vocabulary.
	ID int
	// The length of the implicit sequence of characters.
	Length int
}

// WordSymbol expands a Symbol with contextual information related to the
// Word that contains it.
type WordSymbol struct {
	Symbol
	// Prev is the index of the previous symbol in the Word.
	// -1 means no previous symbol.
	Prev int
	// Prev is the index of the next symbol in the Word.
	// -1 means no next symbol.
	Next int
}

// MergeWith merges the current WordSymbol with the other one.
// In order to update prev/next, we consider the receiver to be the WordSymbol
// on the left, and other to be the next one on the right.
func (s *WordSymbol) MergeWith(other *WordSymbol, newSymbolID int) {
	s.ID = newSymbolID
	s.Length += other.Length
	s.Next = other.Next
}

func (s *WordSymbol) HasPrev() bool {
	return s.Prev != -1
}

func (s *WordSymbol) HasNext() bool {
	return s.Next != -1
}

// Word is a slice of WordSymbol.
type Word []*WordSymbol

// NewWord returns a new empty Word.
func NewWord() *Word {
	w := make(Word, 0)
	return &w
}

func (w *Word) Len() int {
	return len(*w)
}

func (w *Word) getSymbolID(index int) int {
	return (*w)[index].ID
}

// Add appends a new symbol to the Word.
func (w *Word) Add(symbolID int) {
	sym := &WordSymbol{
		Symbol: Symbol{
			ID:     symbolID,
			Length: 1,
		},
		Prev: w.Len() - 1,
		Next: -1,
	}
	if sym.Prev != -1 {
		(*w)[sym.Prev].Next = w.Len()
	}
	*w = append(*w, sym)
}

func (w *Word) MergeAll(merges *MergeMap, dropout float64) {
	symbolsLen := w.Len()
	queue := make(WordMergeHeap, 0, symbolsLen)
	skip := make([]WordMerge, 0, symbolsLen)

	lastSymbolIndex := symbolsLen - 1
	for index := 0; index < lastSymbolIndex; index++ {
		if m, ok := merges.Get(w.getSymbolID(index), w.getSymbolID(index+1)); ok {
			heap.Push(&queue, WordMerge{MergeValue: m, Pos: index})
		}
	}

	hasDropout := dropout > 0
	for queue.Len() > 0 {
		top := heap.Pop(&queue).(WordMerge)

		if hasDropout && rand.Float64() < dropout {
			skip = append(skip, top)
			continue
		}

		// Re-insert the skipped elements
		for _, s := range skip {
			heap.Push(&queue, s)
		}
		skip = skip[:0] // empty `skip` without reallocating memory

		if (*w)[top.Pos].Length == 0 || !(*w)[top.Pos].HasNext() {
			// Do nothing if the symbol is empty, or if it's the last symbol
			continue
		}

		nextPos := (*w)[top.Pos].Next
		right := (*w)[nextPos]

		// Make sure we are not processing an expired queue entry
		if m, ok := merges.Get((*w)[top.Pos].ID, right.ID); !ok || m.ID != top.ID {
			continue
		}

		// Otherwise, let's merge
		(*w)[top.Pos].MergeWith(right, top.ID)
		// Tag the right part as removed
		(*w)[nextPos].Length = 0

		// Update `prev` on the new `next` to the current pos
		if right.HasNext() && right.Next < w.Len() {
			(*w)[right.Next].Prev = top.Pos
		}

		// Insert the new pair formed with the previous symbol
		current := (*w)[top.Pos]
		if current.HasPrev() {
			prevSymbol := (*w)[current.Prev]
			if m, ok := merges.Get(prevSymbol.ID, current.ID); ok {
				heap.Push(&queue, WordMerge{MergeValue: m, Pos: current.Prev})
			}
		}

		// Insert the new pair formed with the next symbol
		if current.HasNext() && current.Next < w.Len() {
			nextSymbol := (*w)[current.Next]
			if m, ok := merges.Get(current.ID, nextSymbol.ID); ok {
				heap.Push(&queue, WordMerge{MergeValue: m, Pos: top.Pos})
			}
		}
	}

	// Filter out the removed symbols
	for i := 0; i < w.Len(); {
		if (*w)[i].Length == 0 {
			w.removeSymbol(i)
			continue
		}
		i++
	}
}

func (w *Word) removeSymbol(index int) {
	*w = append((*w)[:index], (*w)[index+1:]...)
}
