// Copyright (c) 2020, NLP Odyssey Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package stripnormalizer

import (
	. "github.com/nlpodyssey/gotokenizers/normalizers"
	. "github.com/nlpodyssey/gotokenizers/normalizers/normalizedstring"
)

// StripNormalizer allows string normalization removing leading spaces,
// trailing spaces, or both.
type StripNormalizer struct{ left, right bool }

var _ Normalizer = &StripNormalizer{}

// NewStripNormalizer returns a new StripNormalizer, initialized for stripping
// leading spaces (left) and/or trailing spaces (right).
func NewStripNormalizer(left, right bool) *StripNormalizer {
	return &StripNormalizer{left: left, right: right}
}

// Normalize strips the NormalizedString in place.
func (sn *StripNormalizer) Normalize(ns *NormalizedString) error {
	ns.TrimLeftRight(sn.left, sn.right)
	return nil
}
