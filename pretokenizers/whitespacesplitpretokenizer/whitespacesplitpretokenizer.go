// Copyright (c) 2020, NLP Odyssey Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package whitespacesplitpretokenizer

import (
	"github.com/nlpodyssey/gotokenizers/normalizedstring"
	"github.com/nlpodyssey/gotokenizers/pretokenizedstring"
	"github.com/nlpodyssey/gotokenizers/pretokenizers"
	"github.com/nlpodyssey/gotokenizers/splitpattern"
	"unicode"
)

// WhiteSpaceSplitPreTokenizer allows the generation of pre-tokens splitting
// by whitespace-like characters.
type WhiteSpaceSplitPreTokenizer struct{}

var _ pretokenizers.PreTokenizer = &WhiteSpaceSplitPreTokenizer{}

// New returns a new WhiteSpaceSplitPreTokenizer.
func New() *WhiteSpaceSplitPreTokenizer {
	return &WhiteSpaceSplitPreTokenizer{}
}

// PreTokenize splits the NormalizedString by whitespace-like characters
func (w *WhiteSpaceSplitPreTokenizer) PreTokenize(pts *pretokenizedstring.PreTokenizedString) error {
	splittingPattern := splitpattern.FromFunc(func(r rune) bool {
		return unicode.In(r, unicode.White_Space)
	})
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
