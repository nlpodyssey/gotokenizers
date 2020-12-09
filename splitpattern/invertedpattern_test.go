// Copyright (c) 2020, NLP Odyssey Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package splitpattern

import (
	"testing"
)

func TestInvertedPatternFindMatches(t *testing.T) {
	t.Parallel()

	sp := Invert(FromRune('a'))

	runTest(t, sp, "aba", []Capture{
		{Offsets{0, 1}, false},
		{Offsets{1, 2}, true},
		{Offsets{2, 3}, false},
	})
}
