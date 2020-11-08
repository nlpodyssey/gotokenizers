// Copyright (c) 2020, NLP Odyssey Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package bpemodel

import (
	"github.com/nlpodyssey/gotokenizers/vocabulary"
	"reflect"
	"testing"
)

func TestMergeMapFromFile(t *testing.T) {
	t.Parallel()

	vocabTerms := []string{
		"ab",       // 0
		"cd",       // 1
		"efg",      // 2
		"hij",      // 3
		"klmn",     // 4
		"opqr",     // 5
		"abcd",     // 6
		"efghij",   // 7
		"klmnopqr", // 8
	}
	vocab := vocabulary.NewVocabulary()
	for _, term := range vocabTerms {
		vocab.AddTerm(term)
	}

	m, err := MergeMapFromFile("testdata/merges.txt", vocab, 0)
	if err != nil {
		t.Fatal(err)
	}
	expected := MergeMap{
		symbolIDPair{0, 1}: MergeValue{Rank: 0, ID: 6},
		symbolIDPair{2, 3}: MergeValue{Rank: 1, ID: 7},
		symbolIDPair{4, 5}: MergeValue{Rank: 2, ID: 8},
	}
	if !reflect.DeepEqual(*m, expected) {
		t.Errorf("expected:\n  %#v\nactual:\n  %#v\n", expected, *m)
	}
}

func TestMergeMapFromFileWithPrefixLength(t *testing.T) {
	t.Parallel()

	vocabTerms := []string{
		"ab",      // 0
		"cd",      // 1
		"efg",     // 2
		"hij",     // 3
		"klmn",    // 4
		"opqr",    // 5
		"abd",     // 6
		"efgij",   // 7
		"klmnpqr", // 8
	}
	vocab := vocabulary.NewVocabulary()
	for _, term := range vocabTerms {
		vocab.AddTerm(term)
	}

	m, err := MergeMapFromFile("testdata/merges.txt", vocab, 1)
	if err != nil {
		t.Fatal(err)
	}
	expected := MergeMap{
		symbolIDPair{0, 1}: MergeValue{Rank: 0, ID: 6},
		symbolIDPair{2, 3}: MergeValue{Rank: 1, ID: 7},
		symbolIDPair{4, 5}: MergeValue{Rank: 2, ID: 8},
	}
	if !reflect.DeepEqual(*m, expected) {
		t.Errorf("expected:\n  %#v\nactual:\n  %#v\n", expected, *m)
	}
}
