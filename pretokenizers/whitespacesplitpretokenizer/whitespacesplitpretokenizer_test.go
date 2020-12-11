// Copyright (c) 2020, NLP Odyssey Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package whitespacesplitpretokenizer

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
			"Hey man!",
			[]pretokenizedstring.OriginalByteSplit{
				{String: "Hey", Offsets: strutils.ByteOffsets{Start: 0, End: 3}},
				{String: "man!", Offsets: strutils.ByteOffsets{Start: 4, End: 8}},
			},
		},
		{
			"Hey, man, Good?",
			[]pretokenizedstring.OriginalByteSplit{
				{String: "Hey,", Offsets: strutils.ByteOffsets{Start: 0, End: 4}},
				{String: "man,", Offsets: strutils.ByteOffsets{Start: 5, End: 9}},
				{String: "Good?", Offsets: strutils.ByteOffsets{Start: 10, End: 15}},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("%#v", tc.input), func(t *testing.T) {
			pt := New()
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
