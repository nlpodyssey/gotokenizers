// Copyright (c) 2020, NLP Odyssey Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package bertpretokenizer

import (
	"fmt"
	. "github.com/nlpodyssey/gotokenizers/normalizers/normalizedstring"
	. "github.com/nlpodyssey/gotokenizers/pretokenizers"
	"reflect"
	"testing"
)

func TestBertPreTokenizer(t *testing.T) {
	t.Parallel()

	wt := NewBertPreTokenizer()

	tests := []struct {
		str      string
		expected []PreToken
	}{
		{str: "", expected: []PreToken{}},
		{str: " ", expected: []PreToken{}},
		{str: " \n\t", expected: []PreToken{}},
		{str: "x", expected: []PreToken{
			{String: "x", Start: 0, End: 1},
		}},
		{str: "foo", expected: []PreToken{
			{String: "foo", Start: 0, End: 3},
		}},
		{str: "foo bar baz", expected: []PreToken{
			{String: "foo", Start: 0, End: 3},
			{String: "bar", Start: 4, End: 7},
			{String: "baz", Start: 8, End: 11},
		}},
		{str: " foo ", expected: []PreToken{
			{String: "foo", Start: 1, End: 4},
		}},
		{str: "  foo  ", expected: []PreToken{
			{String: "foo", Start: 2, End: 5},
		}},
		{str: " \nfoo   bar \t  baz\n\r", expected: []PreToken{
			{String: "foo", Start: 2, End: 5},
			{String: "bar", Start: 8, End: 11},
			{String: "baz", Start: 15, End: 18},
		}},
		{str: "!", expected: []PreToken{
			{String: "!", Start: 0, End: 1},
		}},
		{str: "!?", expected: []PreToken{
			{String: "!", Start: 0, End: 1},
			{String: "?", Start: 1, End: 2},
		}},
		{str: " ! ? ", expected: []PreToken{
			{String: "!", Start: 1, End: 2},
			{String: "?", Start: 3, End: 4},
		}},
		{str: "Süß Café!?", expected: []PreToken{
			{String: "Süß", Start: 0, End: 3},
			{String: "Café", Start: 4, End: 8},
			{String: "!", Start: 8, End: 9},
			{String: "?", Start: 9, End: 10},
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
