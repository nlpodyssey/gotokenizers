// Copyright (c) 2020, NLP Odyssey Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package runedelimiterpretokenizer

import (
	"fmt"
	. "github.com/nlpodyssey/gotokenizers/normalizers/normalizedstring"
	. "github.com/nlpodyssey/gotokenizers/pretokenizers"
	"reflect"
	"testing"
)

func TestRuneDelimiterPreTokenizer(t *testing.T) {
	t.Parallel()

	wt := NewRuneDelimiterPreTokenizer('Ʊ')

	tests := []struct {
		str      string
		expected []PreToken
	}{
		{str: "", expected: []PreToken{}},
		{str: "Ʊ", expected: []PreToken{}},
		{str: "ƱƱ", expected: []PreToken{}},
		{str: " ", expected: []PreToken{
			{String: " ", Start: 0, End: 1},
		}},
		{str: " \n\t", expected: []PreToken{
			{String: " \n\t", Start: 0, End: 3},
		}},
		{str: "x", expected: []PreToken{
			{String: "x", Start: 0, End: 1},
		}},
		{str: "foo", expected: []PreToken{
			{String: "foo", Start: 0, End: 3},
		}},
		{str: "foo bar", expected: []PreToken{
			{String: "foo bar", Start: 0, End: 7},
		}},
		{str: "fooƱbarƱbaz", expected: []PreToken{
			{String: "foo", Start: 0, End: 3},
			{String: "bar", Start: 4, End: 7},
			{String: "baz", Start: 8, End: 11},
		}},
		{str: "Ʊfoo", expected: []PreToken{
			{String: "foo", Start: 1, End: 4},
		}},
		{str: "fooƱ", expected: []PreToken{
			{String: "foo", Start: 0, End: 3},
		}},
		{str: "ƱSüßƱCafé!?Ʊ", expected: []PreToken{
			{String: "Süß", Start: 1, End: 4},
			{String: "Café!?", Start: 5, End: 11},
		}},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("%#v", test.str), func(t *testing.T) {
			ns := NewNormalizedString(test.str)
			tokens, err := wt.PreTokenize(ns)
			if err != nil {
				t.Error(err)
			}
			if !reflect.DeepEqual(tokens, test.expected) {
				t.Errorf("expected %v, actual %v", test.expected, tokens)
			}
		})
	}
}
