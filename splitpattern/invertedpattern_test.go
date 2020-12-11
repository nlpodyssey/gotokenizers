// Copyright (c) 2020, NLP Odyssey Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package splitpattern

import (
	"github.com/nlpodyssey/gotokenizers/strutils"
	"testing"
)

func TestInvertedPatternFindMatches(t *testing.T) {
	t.Parallel()

	sp := Invert(FromRune('a'))

	runTest(t, sp, "aba", []Capture{
		{strutils.ByteOffsets{Start: 0, End: 1}, false},
		{strutils.ByteOffsets{Start: 1, End: 2}, true},
		{strutils.ByteOffsets{Start: 2, End: 3}, false},
	})
}
