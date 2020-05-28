// Copyright (c) 2020, NLP Odyssey Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package sequencenormalizer

import (
	. "github.com/nlpodyssey/gotokenizers/normalizers"
	. "github.com/nlpodyssey/gotokenizers/normalizers/normalizedstring"
)

// SequenceNormalizer allows concatenating multiple other Normalizers as a
// Sequence.
type SequenceNormalizer struct{ normalizers []Normalizer }

var _ Normalizer = &SequenceNormalizer{}

// NewSequenceNormalizer returns a new SequenceNormalizer, initializing it
// with the ordered sequence of Normalizers.
func NewSequenceNormalizer(normalizers []Normalizer) *SequenceNormalizer {
	return &SequenceNormalizer{normalizers: normalizers}
}

// Normalize transform the NormalizedString running the ordered sequence of
// normalizers (against the same NormalizedString).
//
// If one Normalizer returns an error, the same error is returned and
// the subsequent Normalizers (if any) are ignored.
func (sn *SequenceNormalizer) Normalize(ns *NormalizedString) error {
	for _, normalizer := range sn.normalizers {
		err := normalizer.Normalize(ns)
		if err != nil {
			return err
		}
	}
	return nil
}
