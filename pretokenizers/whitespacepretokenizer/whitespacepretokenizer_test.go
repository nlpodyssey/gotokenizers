// Copyright (c) 2020, NLP Odyssey Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package whitespacepretokenizer

import (
	"fmt"
	. "github.com/nlpodyssey/gotokenizers/normalizers/normalizedstring"
	. "github.com/nlpodyssey/gotokenizers/pretokenizers"
	"reflect"
	"testing"
)

func TestWhiteSpacePreTokenizer(t *testing.T) {
	wt := NewWhiteSpacePreTokenizer()

	tests := []struct {
		str      string
		expected []PreToken
	}{
		{str: "", expected: []PreToken{}},
		{str: " ", expected: []PreToken{}},
		{str: " \n\t", expected: []PreToken{}},
		{str: "x", expected: []PreToken{
			{String: "x", ByteStart: 0, ByteEnd: 1},
		}},
		{str: "foo", expected: []PreToken{
			{String: "foo", ByteStart: 0, ByteEnd: 3},
		}},
		{str: "foo bar", expected: []PreToken{
			{String: "foo", ByteStart: 0, ByteEnd: 3},
			{String: "bar", ByteStart: 4, ByteEnd: 7},
		}},
		{str: "foo \nbar", expected: []PreToken{
			{String: "foo", ByteStart: 0, ByteEnd: 3},
			{String: "bar", ByteStart: 5, ByteEnd: 8},
		}},
		{str: " foo bar ", expected: []PreToken{
			{String: "foo", ByteStart: 1, ByteEnd: 4},
			{String: "bar", ByteStart: 5, ByteEnd: 8},
		}},
		{str: "!", expected: []PreToken{
			{String: "!", ByteStart: 0, ByteEnd: 1},
		}},
		{str: "!?.", expected: []PreToken{
			{String: "!?.", ByteStart: 0, ByteEnd: 3},
		}},
		{str: "Süß!?", expected: []PreToken{
			{String: "Süß", ByteStart: 0, ByteEnd: 5},
			{String: "!?", ByteStart: 5, ByteEnd: 7},
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
