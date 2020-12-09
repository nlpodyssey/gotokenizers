// Copyright (c) 2020, NLP Odyssey Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package whitespacepretokenizer

import (
	"fmt"
	"github.com/nlpodyssey/gotokenizers/normalizedstring"
	"github.com/nlpodyssey/gotokenizers/pretokenizedstring"
	"reflect"
	"testing"
)

func TestWhiteSpacePreTokenizer_PreTokenize(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		input          string
		expectedSplits []pretokenizedstring.OriginalByteSplit
	}{
		{
			"Hey man!",
			[]pretokenizedstring.OriginalByteSplit{
				{String: "Hey", Offsets: normalizedstring.Offsets{Start: 0, End: 3}},
				{String: "man", Offsets: normalizedstring.Offsets{Start: 4, End: 7}},
				{String: "!", Offsets: normalizedstring.Offsets{Start: 7, End: 8}},
			},
		},
		{
			"How are you doing?",
			[]pretokenizedstring.OriginalByteSplit{
				{String: "How", Offsets: normalizedstring.Offsets{Start: 0, End: 3}},
				{String: "are", Offsets: normalizedstring.Offsets{Start: 4, End: 7}},
				{String: "you", Offsets: normalizedstring.Offsets{Start: 8, End: 11}},
				{String: "doing", Offsets: normalizedstring.Offsets{Start: 12, End: 17}},
				{String: "?", Offsets: normalizedstring.Offsets{Start: 17, End: 18}},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("%#v", tc.input), func(t *testing.T) {
			t.Parallel()

			pt := NewDefault()
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
