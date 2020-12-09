// Copyright (c) 2020, NLP Odyssey Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package splitpattern

import (
	"testing"
)

func TestFuncSplitPatternFindMatches(t *testing.T) {
	t.Parallel()

	sp := FromFunc(func(r rune) bool {
		return r == 'b'
	})

	runTest(t, sp, "aba", []Capture{
		{Offsets{0, 1}, false},
		{Offsets{1, 2}, true},
		{Offsets{2, 3}, false},
	})
	runTest(t, sp, "aaaab", []Capture{
		{Offsets{0, 4}, false},
		{Offsets{4, 5}, true},
	})
	runTest(t, sp, "bbaaa", []Capture{
		{Offsets{0, 1}, true},
		{Offsets{1, 2}, true},
		{Offsets{2, 5}, false},
	})
	runTest(t, sp, "", []Capture{
		{Offsets{0, 0}, false},
	})
	runTest(t, sp, "aaa", []Capture{
		{Offsets{0, 3}, false},
	})
}
