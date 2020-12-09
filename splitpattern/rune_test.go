// Copyright (c) 2020, NLP Odyssey Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package splitpattern

import (
	"testing"
)

func TestRuneSplitPatternFindMatches(t *testing.T) {
	t.Parallel()

	sp := FromRune('a')

	runTest(t, sp, "aba", []Capture{
		{Offsets{0, 1}, true},
		{Offsets{1, 2}, false},
		{Offsets{2, 3}, true},
	})
	runTest(t, sp, "bbbba", []Capture{
		{Offsets{0, 4}, false},
		{Offsets{4, 5}, true},
	})
	runTest(t, sp, "aabbb", []Capture{
		{Offsets{0, 1}, true},
		{Offsets{1, 2}, true},
		{Offsets{2, 5}, false},
	})
	runTest(t, sp, "", []Capture{
		{Offsets{0, 0}, false},
	})
	runTest(t, sp, "bbb", []Capture{
		{Offsets{0, 3}, false},
	})
}
