// Copyright (c) 2020, NLP Odyssey Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package metaspacepretokenizer

import (
	"fmt"
	. "github.com/nlpodyssey/gotokenizers/normalizers/normalizedstring"
	. "github.com/nlpodyssey/gotokenizers/pretokenizers"
	"reflect"
	"testing"
)

func TestMetaSpacePreTokenizerWithoutWhiteSpacePrefix(t *testing.T) {
	t.Parallel()

	wt := NewMetaSpacePreTokenizer('Ʊ', false)

	tests := []struct {
		str      string
		expected []PreToken
	}{
		{
			str: "", expected: []PreToken{}},
		{
			str: " ", expected: []PreToken{
				{String: "Ʊ", Start: 0, End: 1},
			}},
		{
			str: " \n\t", expected: []PreToken{
				{String: "Ʊ", Start: 0, End: 1},
				{String: "Ʊ", Start: 1, End: 2},
				{String: "Ʊ", Start: 2, End: 3},
			}},
		{
			str: "x", expected: []PreToken{
				{String: "x", Start: 0, End: 1},
			}},
		{
			str: "foo", expected: []PreToken{
				{String: "foo", Start: 0, End: 3},
			}},
		{
			str: "foo bar\tbaz", expected: []PreToken{
				{String: "foo", Start: 0, End: 3},
				{String: "Ʊbar", Start: 3, End: 7},
				{String: "Ʊbaz", Start: 7, End: 11},
			}},
		{
			str: "foo \nbar", expected: []PreToken{
				{String: "foo", Start: 0, End: 3},
				{String: "Ʊ", Start: 3, End: 4},
				{String: "Ʊbar", Start: 4, End: 8},
			}},
		{
			str: " foo bar ", expected: []PreToken{
				{String: "Ʊfoo", Start: 0, End: 4},
				{String: "Ʊbar", Start: 4, End: 8},
				{String: "Ʊ", Start: 8, End: 9},
			}},
		{
			str: "\nSüß Café!?\r", expected: []PreToken{
				{String: "ƱSüß", Start: 0, End: 4},
				{String: "ƱCafé!?", Start: 4, End: 11},
				{String: "Ʊ", Start: 11, End: 12},
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
			if ns.Get() != test.str {
				t.Errorf(
					"the input string %#v was modified, but it shouldn't: %#v",
					test.str, ns.Get())
			}
		})
	}
}

func TestMetaSpacePreTokenizerWithWhiteSpacePrefix(t *testing.T) {
	t.Parallel()

	wt := NewMetaSpacePreTokenizer('Ʊ', true)

	tests := []struct {
		str            string
		expectedStr    string
		expectedTokens []PreToken
	}{
		{
			str:         "",
			expectedStr: " ",
			expectedTokens: []PreToken{
				{String: "Ʊ", Start: 0, End: 1},
			},
		},
		{
			str:         " ",
			expectedStr: " ",
			expectedTokens: []PreToken{
				{String: "Ʊ", Start: 0, End: 1},
			},
		},
		{
			str:         " \n\t",
			expectedStr: " \n\t",
			expectedTokens: []PreToken{
				{String: "Ʊ", Start: 0, End: 1},
				{String: "Ʊ", Start: 1, End: 2},
				{String: "Ʊ", Start: 2, End: 3},
			},
		},
		{
			str:         "x",
			expectedStr: " x",
			expectedTokens: []PreToken{
				{String: "Ʊx", Start: 0, End: 2},
			},
		},
		{
			str:         "foo",
			expectedStr: " foo",
			expectedTokens: []PreToken{
				{String: "Ʊfoo", Start: 0, End: 4},
			},
		},
		{
			str:         "foo bar\tbaz",
			expectedStr: " foo bar\tbaz",
			expectedTokens: []PreToken{
				{String: "Ʊfoo", Start: 0, End: 4},
				{String: "Ʊbar", Start: 4, End: 8},
				{String: "Ʊbaz", Start: 8, End: 12},
			},
		},
		{
			str:         "foo \nbar",
			expectedStr: " foo \nbar",
			expectedTokens: []PreToken{
				{String: "Ʊfoo", Start: 0, End: 4},
				{String: "Ʊ", Start: 4, End: 5},
				{String: "Ʊbar", Start: 5, End: 9},
			},
		},
		{
			str:         " foo bar ",
			expectedStr: " foo bar ",
			expectedTokens: []PreToken{
				{String: "Ʊfoo", Start: 0, End: 4},
				{String: "Ʊbar", Start: 4, End: 8},
				{String: "Ʊ", Start: 8, End: 9},
			},
		},
		{
			str:         "\nSüß Café!?\r",
			expectedStr: "\nSüß Café!?\r",
			expectedTokens: []PreToken{
				{String: "ƱSüß", Start: 0, End: 4},
				{String: "ƱCafé!?", Start: 4, End: 11},
				{String: "Ʊ", Start: 11, End: 12},
			},
		},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("%#v", test.str), func(t *testing.T) {
			ns := NewNormalizedString(test.str)
			tokens, err := wt.PreTokenize(ns)
			if err != nil {
				t.Error(err)
			}
			if !reflect.DeepEqual(tokens, test.expectedTokens) {
				t.Errorf("expected %v, actual %v", test.expectedTokens, tokens)
			}
			if ns.Get() != test.expectedStr {
				t.Errorf("expected %#v, actual %#v", test.expectedStr, ns.Get())
			}
		})
	}
}

func TestDefaultMetaSpacePreTokenizer(t *testing.T) {
	t.Parallel()

	wt := DefaultMetaSpacePreTokenizer()

	tests := []struct {
		str            string
		expectedStr    string
		expectedTokens []PreToken
	}{
		{
			str:         "foo",
			expectedStr: " foo",
			expectedTokens: []PreToken{
				{String: "▁foo", Start: 0, End: 4},
			},
		},
		{
			str:         "foo",
			expectedStr: " foo",
			expectedTokens: []PreToken{
				{String: "▁foo", Start: 0, End: 4},
			},
		},
		{
			str:         "\nSüß Café!?\r",
			expectedStr: "\nSüß Café!?\r",
			expectedTokens: []PreToken{
				{String: "▁Süß", Start: 0, End: 4},
				{String: "▁Café!?", Start: 4, End: 11},
				{String: "▁", Start: 11, End: 12},
			},
		},
		{
			// This test makes sure that the default meta-character is actually used
			str:         " ",
			expectedStr: " ",
			expectedTokens: []PreToken{
				{String: string(DefaultMetaCharacter), Start: 0, End: 1},
			},
		},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("%#v", test.str), func(t *testing.T) {
			ns := NewNormalizedString(test.str)
			tokens, err := wt.PreTokenize(ns)
			if err != nil {
				t.Error(err)
			}
			if !reflect.DeepEqual(tokens, test.expectedTokens) {
				t.Errorf("expected %v, actual %v", test.expectedTokens, tokens)
			}
			if ns.Get() != test.expectedStr {
				t.Errorf("expected %#v, actual %#v", test.expectedStr, ns.Get())
			}
		})
	}
}
