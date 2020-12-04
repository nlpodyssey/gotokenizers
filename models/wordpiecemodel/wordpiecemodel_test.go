// Copyright (c) 2020, NLP Odyssey Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package wordpiecemodel

import (
	"github.com/nlpodyssey/gotokenizers/models"
	"github.com/nlpodyssey/gotokenizers/pretokenizers"
	"github.com/nlpodyssey/gotokenizers/vocabulary"
	"reflect"
	"testing"
)

func TestWordPieceModelTokenize(t *testing.T) {
	terms := []string{
		"[UNK]",            // 0
		"foo",              // 1
		"##foo",            // 2
		"bar",              // 3
		"##bar",            // 4
		"baz",              // 5
		"##baz",            // 6
		"alpha",            // 7
		"##alpha",          // 8
		"gamma",            // 9
		"##gamma",          // 10
		"veryverylongterm", // 11
	}
	vocab := vocabulary.NewVocabulary()
	for _, term := range terms {
		vocab.AddTerm(term)
	}

	wordPiece := NewWordPieceModel(
		vocab,
		"[UNK]",
		"##",
		15,
	)

	sentence := []pretokenizers.PreToken{
		{String: "foo", Start: 0, End: 3},
		{String: "barbaz", Start: 3, End: 9},
		{String: "alphabetagamma", Start: 9, End: 23},
		{String: "foobarbaz", Start: 23, End: 32},
		{String: "qux", Start: 32, End: 35},
		{String: "veryverylongterm", Start: 35, End: 51},
	}

	tokens, err := wordPiece.Tokenize(sentence)
	if err != nil {
		t.Error(err)
	}
	expectedTokens := []models.Token{
		{ID: 1, Value: "foo", Offsets: models.TokenOffsets{Start: 0, End: 3}},
		{ID: 3, Value: "bar", Offsets: models.TokenOffsets{Start: 3, End: 6}},
		{ID: 6, Value: "##baz", Offsets: models.TokenOffsets{Start: 6, End: 9}},
		{ID: 0, Value: "[UNK]", Offsets: models.TokenOffsets{Start: 9, End: 23}},
		{ID: 1, Value: "foo", Offsets: models.TokenOffsets{Start: 23, End: 26}},
		{ID: 4, Value: "##bar", Offsets: models.TokenOffsets{Start: 26, End: 29}},
		{ID: 6, Value: "##baz", Offsets: models.TokenOffsets{Start: 29, End: 32}},
		{ID: 0, Value: "[UNK]", Offsets: models.TokenOffsets{Start: 32, End: 35}},
		{ID: 0, Value: "[UNK]", Offsets: models.TokenOffsets{Start: 35, End: 51}},
	}
	if !reflect.DeepEqual(tokens, expectedTokens) {
		t.Errorf("expected %+v, actual %+v", expectedTokens, tokens)
	}
}
