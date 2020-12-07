// Copyright (c) 2020, NLP Odyssey Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package splitpattern

import (
	"github.com/dlclark/regexp2"
	"testing"
)

func TestRegexp2SplitPatternFindMatches(t *testing.T) {
	t.Parallel()

	r := regexp2.MustCompile(`\s+`, regexp2.None)
	sp := FromRegexp2(r)

	runTest(t, sp, "a   b", []Capture{
		{Offsets{0, 1}, false},
		{Offsets{1, 4}, true},
		{Offsets{4, 5}, false},
	})

	runTest(t, sp, "   a   b   ", []Capture{
		{Offsets{0, 3}, true},
		{Offsets{3, 4}, false},
		{Offsets{4, 7}, true},
		{Offsets{7, 8}, false},
		{Offsets{8, 11}, true},
	})

	runTest(t, sp, "", []Capture{
		{Offsets{0, 0}, false},
	})

	runTest(t, sp, "𝔾𝕠𝕠𝕕 𝕞𝕠𝕣𝕟𝕚𝕟𝕘", []Capture{
		{Offsets{0, 16}, false},
		{Offsets{16, 17}, true},
		{Offsets{17, 45}, false},
	})

	runTest(t, sp, "aaa", []Capture{
		{Offsets{0, 3}, false},
	})
}
