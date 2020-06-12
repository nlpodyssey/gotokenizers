// Copyright (c) 2020, NLP Odyssey Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package bytelevelpretokenizer

import (
	"fmt"
	"github.com/dlclark/regexp2"
	. "github.com/nlpodyssey/gotokenizers/normalizers/normalizedstring"
	. "github.com/nlpodyssey/gotokenizers/pretokenizers"
	"reflect"
	"testing"
)

func TestByteLevelPreTokenizerWithPrefixSpaceDisabled(t *testing.T) {
	t.Parallel()

	b := NewByteLevelPreTokenizer(DefaultSplittingRegexp, false)

	tests := []struct {
		str            string
		expectedStr    string
		expectedTokens []PreToken
	}{
		{
			str:            "",
			expectedStr:    "",
			expectedTokens: []PreToken{},
		},
		{
			str:         " ",
			expectedStr: "Ġ",
			expectedTokens: []PreToken{
				{String: "Ġ", Start: 0, End: 1},
			},
		},
		{
			str:         " \n\r",
			expectedStr: "ĠĊč",
			expectedTokens: []PreToken{
				{String: "ĠĊč", Start: 0, End: 3},
			},
		},
		{
			str:         "x",
			expectedStr: "x",
			expectedTokens: []PreToken{
				{String: "x", Start: 0, End: 1},
			},
		},
		{
			str:         "Foo",
			expectedStr: "Foo",
			expectedTokens: []PreToken{
				{String: "Foo", Start: 0, End: 3},
			},
		},
		{
			str:         "Foo bar baz",
			expectedStr: "FooĠbarĠbaz",
			expectedTokens: []PreToken{
				{String: "Foo", Start: 0, End: 3},
				{String: "Ġbar", Start: 3, End: 7},
				{String: "Ġbaz", Start: 7, End: 11},
			},
		},
		{
			str:         "  Foo  bar  ",
			expectedStr: "ĠĠFooĠĠbarĠĠ",
			expectedTokens: []PreToken{
				{String: "Ġ", Start: 0, End: 1},
				{String: "ĠFoo", Start: 1, End: 5},
				{String: "Ġ", Start: 5, End: 6},
				{String: "Ġbar", Start: 6, End: 10},
				{String: "ĠĠ", Start: 10, End: 12},
			},
		},
		{
			str:         "Foo\nbar baz",
			expectedStr: "FooĊbarĠbaz",
			expectedTokens: []PreToken{
				{String: "Foo", Start: 0, End: 3},
				{String: "Ċ", Start: 3, End: 4},
				{String: "bar", Start: 4, End: 7},
				{String: "Ġbaz", Start: 7, End: 11},
			},
		},
		{
			str:         "\nSüß Café!?\r",
			expectedStr: "ĊSÃ¼ÃŁĠCafÃ©!?č",
			expectedTokens: []PreToken{
				{String: "Ċ", Start: 0, End: 1},
				{String: "SÃ¼ÃŁ", Start: 1, End: 6},
				{String: "ĠCafÃ©", Start: 6, End: 12},
				{String: "!?", Start: 12, End: 14},
				{String: "č", Start: 14, End: 15},
			},
		},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("%#v", test.str), func(t *testing.T) {
			ns := NewNormalizedString(test.str)
			tokens, err := b.PreTokenize(ns)
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

func TestByteLevelPreTokenizerWithPrefixSpaceEnabled(t *testing.T) {
	t.Parallel()

	b := NewByteLevelPreTokenizer(DefaultSplittingRegexp, true)

	tests := []struct {
		str            string
		expectedStr    string
		expectedTokens []PreToken
	}{
		{
			str:         "",
			expectedStr: "Ġ",
			expectedTokens: []PreToken{
				{String: "Ġ", Start: 0, End: 1},
			},
		},
		{
			str:         " ",
			expectedStr: "Ġ",
			expectedTokens: []PreToken{
				{String: "Ġ", Start: 0, End: 1},
			},
		},
		{
			str:         " \n\r",
			expectedStr: "ĠĊč",
			expectedTokens: []PreToken{
				{String: "ĠĊč", Start: 0, End: 3},
			},
		},
		{
			str:         "x",
			expectedStr: "Ġx",
			expectedTokens: []PreToken{
				{String: "Ġx", Start: 0, End: 2},
			},
		},
		{
			str:         "Foo",
			expectedStr: "ĠFoo",
			expectedTokens: []PreToken{
				{String: "ĠFoo", Start: 0, End: 4},
			},
		},
		{
			str:         "Foo bar baz",
			expectedStr: "ĠFooĠbarĠbaz",
			expectedTokens: []PreToken{
				{String: "ĠFoo", Start: 0, End: 4},
				{String: "Ġbar", Start: 4, End: 8},
				{String: "Ġbaz", Start: 8, End: 12},
			},
		},
		{
			str:         "  Foo  bar  ",
			expectedStr: "ĠĠFooĠĠbarĠĠ",
			expectedTokens: []PreToken{
				{String: "Ġ", Start: 0, End: 1},
				{String: "ĠFoo", Start: 1, End: 5},
				{String: "Ġ", Start: 5, End: 6},
				{String: "Ġbar", Start: 6, End: 10},
				{String: "ĠĠ", Start: 10, End: 12},
			},
		},
		{
			str:         "\nSüß Café!?\r",
			expectedStr: "ĊSÃ¼ÃŁĠCafÃ©!?č",
			expectedTokens: []PreToken{
				{String: "Ċ", Start: 0, End: 1},
				{String: "SÃ¼ÃŁ", Start: 1, End: 6},
				{String: "ĠCafÃ©", Start: 6, End: 12},
				{String: "!?", Start: 12, End: 14},
				{String: "č", Start: 14, End: 15},
			},
		},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("%#v", test.str), func(t *testing.T) {
			ns := NewNormalizedString(test.str)
			tokens, err := b.PreTokenize(ns)
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

func TestDefaultByteLevelPreTokenizer(t *testing.T) {
	t.Parallel()

	b := DefaultByteLevelPreTokenizer()

	tests := []struct {
		str            string
		expectedStr    string
		expectedTokens []PreToken
	}{
		{
			str:         "Foo bar baz",
			expectedStr: "ĠFooĠbarĠbaz",
			expectedTokens: []PreToken{
				{String: "ĠFoo", Start: 0, End: 4},
				{String: "Ġbar", Start: 4, End: 8},
				{String: "Ġbaz", Start: 8, End: 12},
			},
		},
		{
			str:         "\nSüß Café!?\r",
			expectedStr: "ĊSÃ¼ÃŁĠCafÃ©!?č",
			expectedTokens: []PreToken{
				{String: "Ċ", Start: 0, End: 1},
				{String: "SÃ¼ÃŁ", Start: 1, End: 6},
				{String: "ĠCafÃ©", Start: 6, End: 12},
				{String: "!?", Start: 12, End: 14},
				{String: "č", Start: 14, End: 15},
			},
		},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("%#v", test.str), func(t *testing.T) {
			ns := NewNormalizedString(test.str)
			tokens, err := b.PreTokenize(ns)
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

func TestByteLevelPreTokenizerWithCustomRegexp(t *testing.T) {
	t.Parallel()

	splittingRegepx := regexp2.MustCompile(`foo|süß|bar`, regexp2.None)
	b := NewByteLevelPreTokenizer(splittingRegepx, false)

	tests := []struct {
		str            string
		expectedStr    string
		expectedTokens []PreToken
	}{
		{
			str:         " x foo baz süß café bar qux ",
			expectedStr: "ĠxĠfooĠbazĠsÃ¼ÃŁĠcafÃ©ĠbarĠquxĠ",
			expectedTokens: []PreToken{
				{String: "ĠxĠ", Start: 0, End: 3},
				{String: "foo", Start: 3, End: 6},
				{String: "ĠbazĠ", Start: 6, End: 11},
				{String: "sÃ¼ÃŁ", Start: 11, End: 16},
				{String: "ĠcafÃ©Ġ", Start: 16, End: 23},
				{String: "bar", Start: 23, End: 26},
				{String: "ĠquxĠ", Start: 26, End: 31},
			},
		},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("%#v", test.str), func(t *testing.T) {
			ns := NewNormalizedString(test.str)
			tokens, err := b.PreTokenize(ns)
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
