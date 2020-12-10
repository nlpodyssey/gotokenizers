// Copyright (c) 2020, NLP Odyssey Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package bertpretokenizer

import (
	"github.com/nlpodyssey/gotokenizers/normalizedstring"
	"github.com/nlpodyssey/gotokenizers/pretokenizedstring"
	"github.com/nlpodyssey/gotokenizers/pretokenizers"
	"github.com/nlpodyssey/gotokenizers/splitpattern"
	"unicode"
)

// BertPreTokenizer allows the generation of pre-tokens suitable in the
// context of BERT (Bidirectional Encoder Representations from Transformers)
// models.
//
// Input strings are split by whitespace-like characters, which are excluded
// from the tokens, and by unicode punctuation marks, keeping each punctuation
// character on a separate token.
type BertPreTokenizer struct{}

var _ pretokenizers.PreTokenizer = &BertPreTokenizer{}

// New returns a new BertPreTokenizer.
func New() *BertPreTokenizer {
	return &BertPreTokenizer{}
}

// PreTokenize splits the NormalizedString into pre-tokens suitable for BERT
// models.
func (b *BertPreTokenizer) PreTokenize(pts *pretokenizedstring.PreTokenizedString) error {
	isWhitespacePattern := splitpattern.FromFunc(func(r rune) bool {
		return unicode.In(r, unicode.White_Space)
	})
	isBertPunctuationPattern := splitpattern.FromFunc(func(r rune) bool {
		//TODO: (from Rust) char::is_ascii_punctuation(&x) || ...
		return unicode.In(r, unicode.Punct)
	})
	err := pts.Split(
		func(_ int, ns *normalizedstring.NormalizedString) ([]pretokenizedstring.Split, error) {
			nss, err := ns.Split(isWhitespacePattern, normalizedstring.SplitDelimiterRemoved)
			if err != nil {
				return nil, err
			}
			return pretokenizedstring.SplitsFromNormalizedStrings(nss), nil
		},
	)
	if err != nil {
		return err
	}
	return pts.Split(
		func(_ int, ns *normalizedstring.NormalizedString) ([]pretokenizedstring.Split, error) {
			nss, err := ns.Split(isBertPunctuationPattern, normalizedstring.SplitDelimiterIsolated)
			if err != nil {
				return nil, err
			}
			return pretokenizedstring.SplitsFromNormalizedStrings(nss), nil
		},
	)
}
