// Copyright (c) 2020, NLP Odyssey Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package bpemodel

import (
	"github.com/nlpodyssey/gotokenizers/models"
	"github.com/nlpodyssey/gotokenizers/vocabulary"
	"reflect"
	"testing"
)

func TestTokenizeWithAndWithoutDropout(t *testing.T) {
	// Test tokenization. With dropout set to 0 tokenization is deterministic,
	// so we know exactly what the result should be.
	//
	// To test this, we'll build a simple model to tokenize the word 'unrelated'.
	terms := []string{
		"u",         // 0
		"n",         // 1
		"r",         // 2
		"e",         // 3
		"l",         // 4
		"a",         // 5
		"t",         // 6
		"d",         // 7
		"re",        // 8
		"at",        // 9
		"ed",        // 10
		"un",        // 11
		"ated",      // 12
		"rel",       // 13
		"related",   // 14
		"unrelated", // 15
	}
	vocab := vocabulary.NewVocabulary()
	for _, term := range terms {
		vocab.AddTerm(term)
	}

	mergeItems := []struct {
		k1, k2 string
		rank   int
		term   string
	}{
		{k1: "r", k2: "e", rank: 1, term: "re"},
		{k1: "a", k2: "t", rank: 2, term: "at"},
		{k1: "e", k2: "d", rank: 3, term: "ed"},
		{k1: "u", k2: "n", rank: 4, term: "un"},
		{k1: "at", k2: "ed", rank: 5, term: "ated"},
		{k1: "re", k2: "l", rank: 6, term: "rel"},
		{k1: "rel", k2: "ated", rank: 7, term: "related"},
		{k1: "un", k2: "related", rank: 8, term: "unrelated"},
	}
	merges := NewMergeMap()
	for _, m := range mergeItems {
		id1, _ := vocab.GetID(m.k1)
		id2, _ := vocab.GetID(m.k2)
		id3, _ := vocab.GetID(m.term)
		merges.Set(id1, id2, MergeValue{Rank: m.rank, ID: id3})
	}

	bpe := NewBPEModel(
		vocab,
		merges,
		DefaultCacheCapacity,
		0,
		"",
		"",
		"",
		false,
	)

	// With no dropout:
	tokens, err := bpe.Tokenize("unrelated")
	if err != nil {
		t.Error(err)
	}
	expectedTokens := []models.Token{
		{
			ID:      15,
			Value:   "unrelated",
			Offsets: models.TokenOffsets{Start: 0, End: 9},
		},
	}
	if !reflect.DeepEqual(tokens, expectedTokens) {
		t.Errorf("expected %+v, actual %+v", expectedTokens, tokens)
	}

	// Now set dropout to 1.0. Result should be no merges performed.
	bpe = NewBPEModel(
		vocab,
		merges,
		DefaultCacheCapacity,
		1,
		"",
		"",
		"",
		false,
	)
	tokens, err = bpe.Tokenize("unrelated")
	if err != nil {
		t.Error(err)
	}

	expectedTokens = []models.Token{
		{ID: 0, Value: "u", Offsets: models.TokenOffsets{Start: 0, End: 1}},
		{ID: 1, Value: "n", Offsets: models.TokenOffsets{Start: 1, End: 2}},
		{ID: 2, Value: "r", Offsets: models.TokenOffsets{Start: 2, End: 3}},
		{ID: 3, Value: "e", Offsets: models.TokenOffsets{Start: 3, End: 4}},
		{ID: 4, Value: "l", Offsets: models.TokenOffsets{Start: 4, End: 5}},
		{ID: 5, Value: "a", Offsets: models.TokenOffsets{Start: 5, End: 6}},
		{ID: 6, Value: "t", Offsets: models.TokenOffsets{Start: 6, End: 7}},
		{ID: 3, Value: "e", Offsets: models.TokenOffsets{Start: 7, End: 8}},
		{ID: 7, Value: "d", Offsets: models.TokenOffsets{Start: 8, End: 9}},
	}
	if !reflect.DeepEqual(tokens, expectedTokens) {
		t.Errorf("expected %+v, actual %+v", expectedTokens, tokens)
	}

	// Now try with dropout between 0 and 1.
	bpe = NewBPEModel(
		vocab,
		merges,
		DefaultCacheCapacity,
		0.3,
		"",
		"",
		"",
		false,
	)
	tokens, err = bpe.Tokenize("unrelated")
	if err != nil {
		t.Error(err)
	}
	t.Logf("BPE tokens with dropout 0.5: %v", tokens)
	if len(tokens) == 0 || len(tokens) > 9 {
		t.Errorf("expected 0 < len(tokens) < 0, got %v => %+v", len(tokens), tokens)
	}
}
