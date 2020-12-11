// Copyright (c) 2020, NLP Odyssey Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package wordpiecemodel

import (
	"fmt"
	"github.com/nlpodyssey/gotokenizers/models"
	"github.com/nlpodyssey/gotokenizers/strutils"
	"github.com/nlpodyssey/gotokenizers/vocabulary"
	"reflect"
	"testing"
)

func TestWordPieceModelTokenize(t *testing.T) {
	t.Parallel()

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

	wordPiece := New(
		vocab,
		"[UNK]",
		"##",
		15,
	)

	testCases := []struct {
		input    string
		expected []models.Token
	}{
		{
			"foo",
			[]models.Token{
				{ID: 1, Value: "foo", Offsets: strutils.ByteOffsets{Start: 0, End: 3}},
			},
		},
		{
			"barbaz",
			[]models.Token{
				{ID: 3, Value: "bar", Offsets: strutils.ByteOffsets{Start: 0, End: 3}},
				{ID: 6, Value: "##baz", Offsets: strutils.ByteOffsets{Start: 3, End: 6}},
			},
		},
		{
			"alphabetagamma",
			[]models.Token{
				{ID: 0, Value: "[UNK]", Offsets: strutils.ByteOffsets{Start: 0, End: 14}},
			},
		},
		{
			"foobarbaz",
			[]models.Token{
				{ID: 1, Value: "foo", Offsets: strutils.ByteOffsets{Start: 0, End: 3}},
				{ID: 4, Value: "##bar", Offsets: strutils.ByteOffsets{Start: 3, End: 6}},
				{ID: 6, Value: "##baz", Offsets: strutils.ByteOffsets{Start: 6, End: 9}},
			},
		},
		{
			"qux",
			[]models.Token{
				{ID: 0, Value: "[UNK]", Offsets: strutils.ByteOffsets{Start: 0, End: 3}},
			},
		},
		{
			"veryverylongterm",
			[]models.Token{
				{ID: 0, Value: "[UNK]", Offsets: strutils.ByteOffsets{Start: 0, End: 16}},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("%#v", tc.input), func(t *testing.T) {
			tokens, err := wordPiece.Tokenize(tc.input)
			if err != nil {
				t.Fatal(err)
			}
			assertEqual(t, tokens, tc.expected)
		})
	}
}

func assertEqual(t *testing.T, actual, expected interface{}) {
	t.Helper()
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("expected\n  %#v\nactual\n  %#v", expected, actual)
	}
}
