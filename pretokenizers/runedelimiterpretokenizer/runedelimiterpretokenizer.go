// Copyright (c) 2020, NLP Odyssey Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package runedelimiterpretokenizer

import (
	"github.com/nlpodyssey/gotokenizers/normalizedstring"
	"github.com/nlpodyssey/gotokenizers/pretokenizedstring"
	"github.com/nlpodyssey/gotokenizers/pretokenizers"
	"github.com/nlpodyssey/gotokenizers/splitpattern"
)

// RuneDelimiterPreTokenizer allows the generation of pre-tokens splitting
// by a specific rune.
type RuneDelimiterPreTokenizer struct {
	delimiter rune
}

var _ pretokenizers.PreTokenizer = &RuneDelimiterPreTokenizer{}

// New returns a new RuneDelimiterPreTokenizer,
// setting the given rune as delimiter.
func New(delimiter rune) *RuneDelimiterPreTokenizer {
	return &RuneDelimiterPreTokenizer{delimiter: delimiter}
}

// PreTokenize splits the NormalizedString by rune delimiter.
func (r *RuneDelimiterPreTokenizer) PreTokenize(pts *pretokenizedstring.PreTokenizedString) error {
	splittingPattern := splitpattern.FromRune(r.delimiter)
	return pts.Split(
		func(_ int, ns *normalizedstring.NormalizedString) ([]pretokenizedstring.Split, error) {
			nss, err := ns.Split(splittingPattern, normalizedstring.SplitDelimiterRemoved)
			if err != nil {
				return nil, err
			}
			return pretokenizedstring.SplitsFromNormalizedStrings(nss), nil
		},
	)
}
