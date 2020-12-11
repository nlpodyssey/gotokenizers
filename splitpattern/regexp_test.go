// Copyright (c) 2020, NLP Odyssey Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package splitpattern

import (
	"github.com/nlpodyssey/gotokenizers/strutils"
	"regexp"
	"testing"
)

func TestRegexpSplitPatternFindMatches(t *testing.T) {
	t.Parallel()

	r := regexp.MustCompile(`\s+`)
	sp := FromRegexp(r)

	runTest(t, sp, "a   b", []Capture{
		{strutils.ByteOffsets{Start: 0, End: 1}, false},
		{strutils.ByteOffsets{Start: 1, End: 4}, true},
		{strutils.ByteOffsets{Start: 4, End: 5}, false},
	})

	runTest(t, sp, "   a   b   ", []Capture{
		{strutils.ByteOffsets{Start: 0, End: 3}, true},
		{strutils.ByteOffsets{Start: 3, End: 4}, false},
		{strutils.ByteOffsets{Start: 4, End: 7}, true},
		{strutils.ByteOffsets{Start: 7, End: 8}, false},
		{strutils.ByteOffsets{Start: 8, End: 11}, true},
	})

	runTest(t, sp, "", []Capture{
		{strutils.ByteOffsets{Start: 0, End: 0}, false},
	})

	runTest(t, sp, "𝔾𝕠𝕠𝕕 𝕞𝕠𝕣𝕟𝕚𝕟𝕘", []Capture{
		{strutils.ByteOffsets{Start: 0, End: 16}, false},
		{strutils.ByteOffsets{Start: 16, End: 17}, true},
		{strutils.ByteOffsets{Start: 17, End: 45}, false},
	})

	runTest(t, sp, "aaa", []Capture{
		{strutils.ByteOffsets{Start: 0, End: 3}, false},
	})
}
