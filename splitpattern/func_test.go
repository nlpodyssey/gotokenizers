// Copyright (c) 2020, NLP Odyssey Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package splitpattern

import (
	"github.com/nlpodyssey/gotokenizers/strutils"
	"testing"
)

func TestFuncSplitPatternFindMatches(t *testing.T) {
	t.Parallel()

	sp := FromFunc(func(r rune) bool {
		return r == 'b'
	})

	runTest(t, sp, "aba", []Capture{
		{strutils.ByteOffsets{Start: 0, End: 1}, false},
		{strutils.ByteOffsets{Start: 1, End: 2}, true},
		{strutils.ByteOffsets{Start: 2, End: 3}, false},
	})
	runTest(t, sp, "aaaab", []Capture{
		{strutils.ByteOffsets{Start: 0, End: 4}, false},
		{strutils.ByteOffsets{Start: 4, End: 5}, true},
	})
	runTest(t, sp, "bbaaa", []Capture{
		{strutils.ByteOffsets{Start: 0, End: 1}, true},
		{strutils.ByteOffsets{Start: 1, End: 2}, true},
		{strutils.ByteOffsets{Start: 2, End: 5}, false},
	})
	runTest(t, sp, "", []Capture{
		{strutils.ByteOffsets{Start: 0, End: 0}, false},
	})
	runTest(t, sp, "aaa", []Capture{
		{strutils.ByteOffsets{Start: 0, End: 3}, false},
	})
}
