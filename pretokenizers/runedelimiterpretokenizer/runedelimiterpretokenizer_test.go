// Copyright (c) 2020, NLP Odyssey Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package runedelimiterpretokenizer

import (
	"fmt"
	"github.com/nlpodyssey/gotokenizers/pretokenizedstring"
	"github.com/nlpodyssey/gotokenizers/strutils"
	"reflect"
	"testing"
)

func TestWhiteSpaceSplitPreTokenizer_PreTokenize(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		input          string
		expectedSplits []pretokenizedstring.OriginalByteSplit
	}{
		{
			"",
			[]pretokenizedstring.OriginalByteSplit{}},
		{
			"Ʊ",
			[]pretokenizedstring.OriginalByteSplit{}},
		{
			"ƱƱ",
			[]pretokenizedstring.OriginalByteSplit{}},
		{
			" ",
			[]pretokenizedstring.OriginalByteSplit{
				{String: " ", Offsets: strutils.ByteOffsets{Start: 0, End: 1}},
			},
		},
		{
			" \n\t",
			[]pretokenizedstring.OriginalByteSplit{
				{String: " \n\t", Offsets: strutils.ByteOffsets{Start: 0, End: 3}},
			},
		},
		{
			"x",
			[]pretokenizedstring.OriginalByteSplit{
				{String: "x", Offsets: strutils.ByteOffsets{Start: 0, End: 1}},
			},
		},
		{
			"foo",
			[]pretokenizedstring.OriginalByteSplit{
				{String: "foo", Offsets: strutils.ByteOffsets{Start: 0, End: 3}},
			},
		},
		{
			"foo bar",
			[]pretokenizedstring.OriginalByteSplit{
				{String: "foo bar", Offsets: strutils.ByteOffsets{Start: 0, End: 7}},
			},
		},
		{
			"fooƱbarƱbaz",
			[]pretokenizedstring.OriginalByteSplit{
				{String: "foo", Offsets: strutils.ByteOffsets{Start: 0, End: 3}},
				{String: "bar", Offsets: strutils.ByteOffsets{Start: 5, End: 8}},
				{String: "baz", Offsets: strutils.ByteOffsets{Start: 10, End: 13}},
			},
		},
		{
			"Ʊfoo",
			[]pretokenizedstring.OriginalByteSplit{
				{String: "foo", Offsets: strutils.ByteOffsets{Start: 2, End: 5}},
			},
		},
		{
			"fooƱ",
			[]pretokenizedstring.OriginalByteSplit{
				{String: "foo", Offsets: strutils.ByteOffsets{Start: 0, End: 3}},
			},
		},
		{
			"ƱSüßƱCafé!?Ʊ",
			[]pretokenizedstring.OriginalByteSplit{
				{String: "Süß", Offsets: strutils.ByteOffsets{Start: 2, End: 7}},
				{String: "Café!?", Offsets: strutils.ByteOffsets{Start: 9, End: 16}},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("%#v", tc.input), func(t *testing.T) {
			pt := New('Ʊ')
			pts := pretokenizedstring.FromString(tc.input)
			err := pt.PreTokenize(pts)
			if err != nil {
				t.Fatal(err)
			}

			assertEqual(t, pts.GetOriginalByteSplits(), tc.expectedSplits)
		})
	}
}

func assertEqual(t *testing.T, actual, expected interface{}) {
	t.Helper()
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("expected\n  %#v\nactual\n  %#v", expected, actual)
	}
}
