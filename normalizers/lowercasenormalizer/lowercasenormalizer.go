// Copyright (c) 2020, NLP Odyssey Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package lowercasenormalizer

import (
	"github.com/nlpodyssey/gotokenizers/normalizers"
	"github.com/nlpodyssey/gotokenizers/normalizers/normalizedstring"
)

// LowerCaseNormalizer allows string normalization remapping all Unicode
// letters to their lower case.
type LowerCaseNormalizer struct{}

var _ normalizers.Normalizer = &LowerCaseNormalizer{}

// NewLowerCaseNormalizer returns a new LowerCaseNormalizer.
func NewLowerCaseNormalizer() *LowerCaseNormalizer {
	return &LowerCaseNormalizer{}
}

// Normalize transform the NormalizedString to lowercase in place.
func (sn *LowerCaseNormalizer) Normalize(ns *normalizedstring.NormalizedString) error {
	ns.ToLower()
	return nil
}
