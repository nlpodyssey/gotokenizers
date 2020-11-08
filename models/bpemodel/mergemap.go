// Copyright (c) 2020, NLP Odyssey Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package bpemodel

import (
	"bufio"
	"fmt"
	"github.com/nlpodyssey/gotokenizers/vocabulary"
	"os"
	"strings"
)

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

// MergeMapFromFile reads merges from file.
func MergeMapFromFile(
	filename string,
	vocab *vocabulary.Vocabulary,
	prefixLength int,
) (m *MergeMap, err error) {
	m = NewMergeMap()

	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer func() {
		if e := file.Close(); e != nil && err == nil {
			err = e
		}
	}()

	scanner := bufio.NewScanner(file)
	for lineCount, rank := 1, 0; scanner.Scan(); lineCount++ {
		line := scanner.Text()

		if strings.HasPrefix(line, "#version") {
			continue
		}
		terms := strings.Split(line, " ")
		if len(terms) != 2 {
			return nil, fmt.Errorf("line %d: malformed merges", lineCount)
		}

		leftID, leftOK := vocab.GetID(terms[0])
		if !leftOK {
			return nil, fmt.Errorf("line %d: left merge token is out of vocabulary", lineCount)
		}
		rightID, rightOK := vocab.GetID(terms[1])
		if !rightOK {
			return nil, fmt.Errorf("line %d: right merge token is out of vocabulary", lineCount)
		}

		mergedTerm := fmt.Sprintf("%s%s", terms[0], terms[1][prefixLength:])
		mergedID, mergedOK := vocab.GetID(mergedTerm)
		if !mergedOK {
			return nil, fmt.Errorf("line %d: merged token is out of vocabulary", lineCount)
		}

		m.Set(leftID, rightID, MergeValue{Rank: rank, ID: mergedID})
		rank++
	}
	if err = scanner.Err(); err != nil {
		return nil, err
	}

	return m, nil
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
