// Copyright (c) 2020, NLP Odyssey Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package metaspacepretokenizer

import (
	"fmt"
	"github.com/nlpodyssey/gotokenizers/pretokenizedstring"
	"github.com/nlpodyssey/gotokenizers/strutils"
	"reflect"
	"testing"
)

func TestMetaSpacePreTokenizer_PreTokenize(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		input                    string
		expectedOriginalSplits   []pretokenizedstring.OriginalByteSplit
		expectedNormalizedSplits []pretokenizedstring.NormalizedByteSplit
	}{
		{
			"Hey friend!",
			[]pretokenizedstring.OriginalByteSplit{
				{String: "▁Hey", Offsets: strutils.ByteOffsets{Start: 0, End: 3}},
				{String: "▁friend!", Offsets: strutils.ByteOffsets{Start: 3, End: 11}},
			},
			[]pretokenizedstring.NormalizedByteSplit{
				{String: "▁Hey", Offsets: strutils.ByteOffsets{Start: 0, End: 6}},
				{String: "▁friend!", Offsets: strutils.ByteOffsets{Start: 6, End: 16}},
			},
		},
		{
			"Hey   friend!",
			[]pretokenizedstring.OriginalByteSplit{
				{String: "▁Hey", Offsets: strutils.ByteOffsets{Start: 0, End: 3}},
				{String: "▁", Offsets: strutils.ByteOffsets{Start: 3, End: 4}},
				{String: "▁", Offsets: strutils.ByteOffsets{Start: 4, End: 5}},
				{String: "▁friend!", Offsets: strutils.ByteOffsets{Start: 5, End: 13}},
			},
			[]pretokenizedstring.NormalizedByteSplit{
				{String: "▁Hey", Offsets: strutils.ByteOffsets{Start: 0, End: 6}},
				{String: "▁", Offsets: strutils.ByteOffsets{Start: 6, End: 9}},
				{String: "▁", Offsets: strutils.ByteOffsets{Start: 9, End: 12}},
				{String: "▁friend!", Offsets: strutils.ByteOffsets{Start: 12, End: 22}},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("%#v", tc.input), func(t *testing.T) {
			pt := NewDefault()
			pts := pretokenizedstring.FromString(tc.input)
			err := pt.PreTokenize(pts)
			if err != nil {
				t.Fatal(err)
			}

			assertEqual(t, pts.GetOriginalByteSplits(), tc.expectedOriginalSplits)
			assertEqual(t, pts.GetNormalizedByteSplits(), tc.expectedNormalizedSplits)
		})
	}
}

func assertEqual(t *testing.T, actual, expected interface{}) {
	t.Helper()
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("expected\n  %#v\nactual\n  %#v", expected, actual)
	}
}
