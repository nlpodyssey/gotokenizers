// Copyright (c) 2020, NLP Odyssey Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package whitespacesplitpretokenizer

import (
	"fmt"
	. "github.com/nlpodyssey/gotokenizers/normalizers/normalizedstring"
	. "github.com/nlpodyssey/gotokenizers/pretokenizers"
	"reflect"
	"testing"
)

func TestWhiteSpaceSplitPreTokenizer(t *testing.T) {
	wt := NewWhiteSpaceSplitPreTokenizer()

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
		{str: "foo bar", expected: []PreToken{
			{String: "foo", Start: 0, End: 3},
			{String: "bar", Start: 4, End: 7},
		}},
		{str: "foo \nbar", expected: []PreToken{
			{String: "foo", Start: 0, End: 3},
			{String: "bar", Start: 5, End: 8},
		}},
		{str: " foo bar ", expected: []PreToken{
			{String: "foo", Start: 1, End: 4},
			{String: "bar", Start: 5, End: 8},
		}},
		{str: "!", expected: []PreToken{
			{String: "!", Start: 0, End: 1},
		}},
		{str: "!?.", expected: []PreToken{
			{String: "!?.", Start: 0, End: 3},
		}},
		{str: "Süß!?", expected: []PreToken{
			{String: "Süß!?", Start: 0, End: 5},
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
