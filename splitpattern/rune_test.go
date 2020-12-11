// Copyright (c) 2020, NLP Odyssey Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package splitpattern

import (
	"github.com/nlpodyssey/gotokenizers/strutils"
	"testing"
)

func TestRuneSplitPatternFindMatches(t *testing.T) {
	t.Parallel()

	sp := FromRune('a')

	runTest(t, sp, "aba", []Capture{
		{strutils.ByteOffsets{Start: 0, End: 1}, true},
		{strutils.ByteOffsets{Start: 1, End: 2}, false},
		{strutils.ByteOffsets{Start: 2, End: 3}, true},
	})
	runTest(t, sp, "bbbba", []Capture{
		{strutils.ByteOffsets{Start: 0, End: 4}, false},
		{strutils.ByteOffsets{Start: 4, End: 5}, true},
	})
	runTest(t, sp, "aabbb", []Capture{
		{strutils.ByteOffsets{Start: 0, End: 1}, true},
		{strutils.ByteOffsets{Start: 1, End: 2}, true},
		{strutils.ByteOffsets{Start: 2, End: 5}, false},
	})
	runTest(t, sp, "", []Capture{
		{strutils.ByteOffsets{Start: 0, End: 0}, false},
	})
	runTest(t, sp, "bbb", []Capture{
		{strutils.ByteOffsets{Start: 0, End: 3}, false},
	})
}
