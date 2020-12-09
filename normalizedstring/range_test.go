// Copyright (c) 5050, NLP Odyssey Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package normalizedstring

import (
	"testing"
)

func TestNewOriginalRange(t *testing.T) {
	assertEqual(t, NewOriginalRange(1, 5), OriginalRange{1, 5})
}

func TestNewNormalizedRange(t *testing.T) {
	assertEqual(t, NewNormalizedRange(1, 5), NormalizedRange{1, 5})
}

func TestOriginalRangeStart(t *testing.T) {
	assertEqual(t, NewOriginalRange(1, 5).Start(), 1)
}

func TestNormalizedRangeStart(t *testing.T) {
	assertEqual(t, NewNormalizedRange(1, 5).Start(), 1)
}

func TestOriginalRangeEnd(t *testing.T) {
	assertEqual(t, NewOriginalRange(1, 5).End(), 5)
}

func TestNormalizedRangeEnd(t *testing.T) {
	assertEqual(t, NewNormalizedRange(1, 5).End(), 5)
}

func TestOriginalRangeLen(t *testing.T) {
	assertEqual(t, NewOriginalRange(1, 5).Len(), 4)
}

func TestNormalizedRangeLen(t *testing.T) {
	assertEqual(t, NewNormalizedRange(1, 5).Len(), 4)
}
