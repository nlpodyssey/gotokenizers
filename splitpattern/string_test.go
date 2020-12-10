// Copyright (c) 2020, NLP Odyssey Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package splitpattern

import (
	"testing"
)

func TestStringSplitPatternFindMatches(t *testing.T) {
	t.Parallel()

	{
		sp := FromString("a")

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
	}

	{
		sp := FromString("ab")

		runTest(t, sp, "aabbb", []Capture{
			{Offsets{0, 1}, false},
			{Offsets{1, 3}, true},
			{Offsets{3, 5}, false},
		})
		runTest(t, sp, "aabbab", []Capture{
			{Offsets{0, 1}, false},
			{Offsets{1, 3}, true},
			{Offsets{3, 4}, false},
			{Offsets{4, 6}, true},
		})
	}

	{
		sp := FromString("")

		runTest(t, sp, "", []Capture{
			{Offsets{0, 0}, false},
		})
		runTest(t, sp, "aaa", []Capture{
			{Offsets{0, 3}, false},
		})
	}

	runTest(t, FromString("b"), "aaa", []Capture{
		{Offsets{0, 3}, false},
	})
}
