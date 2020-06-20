// Copyright (c) 2020, NLP Odyssey Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package bpemodel

// MergeMap maps pairs of Symbol IDs to (Rank, ID) values.
type MergeMap map[symbolIDPair]MergeValue

// MergeValue is a (Rank, ID) pair.
type MergeValue struct {
	// Rank determines the order in which a merge is applied during
	// tokenization.
	Rank int
	// ID is the vocabulary ID of the symbol resulting from merging a pair of
	// symbols.
	ID int
}

type symbolIDPair [2]int

// NewMergeMap a new empty MergeMap.
func NewMergeMap() *MergeMap {
	m := make(MergeMap)
	return &m
}

// Get returns a value associated to the given pair of ID, and whether
// the value exists in the map.
func (m *MergeMap) Get(firstID, secondID int) (MergeValue, bool) {
	v, ok := (*m)[symbolIDPair{firstID, secondID}]
	return v, ok
}

func (m *MergeMap) Set(firstID, secondID int, v MergeValue) {
	(*m)[symbolIDPair{firstID, secondID}] = v
}
