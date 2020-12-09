// Copyright (c) 2020, NLP Odyssey Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package whitespacepretokenizer

import (
	"github.com/dlclark/regexp2"
	"github.com/nlpodyssey/gotokenizers/normalizedstring"
	"github.com/nlpodyssey/gotokenizers/pretokenizedstring"
	"github.com/nlpodyssey/gotokenizers/pretokenizers"
	"github.com/nlpodyssey/gotokenizers/splitpattern"
)

// WhiteSpacePreTokenizer allows the generation of pre-tokens made by
// distinct groups of unicode letters (words) and non-letter characters
// (such as punctuation signs or other symbols). Whitespace-like characters
// are always identified as explicit tokens separators.
type WhiteSpacePreTokenizer struct {
	r *regexp2.Regexp
}

// (readonly)
var DefaultWordRegexp = regexp2.MustCompile(`\w+|[^\w\s]+`, regexp2.IgnoreCase|regexp2.Multiline)

var _ pretokenizers.PreTokenizer = &WhiteSpacePreTokenizer{}

// New returns a new WhiteSpacePreTokenizer.
func New(r *regexp2.Regexp) *WhiteSpacePreTokenizer {
	return &WhiteSpacePreTokenizer{r: r}
}

func NewDefault() *WhiteSpacePreTokenizer {
	return New(DefaultWordRegexp)
}

// PreTokenize splits the NormalizedString into word and non-word groups
// separated by whitespace-like characters.
func (w *WhiteSpacePreTokenizer) PreTokenize(pts *pretokenizedstring.PreTokenizedString) error {
	splittingPattern := splitpattern.Invert(splitpattern.FromRegexp2(w.r))

	err := pts.Split(
		func(_ int, ns *normalizedstring.NormalizedString) ([]pretokenizedstring.Split, error) {
			nss, err := ns.Split(splittingPattern, normalizedstring.SplitDelimiterRemoved)
			if err != nil {
				return nil, err
			}
			return pretokenizedstring.SplitsFromNormalizedStrings(nss), nil
		},
	)
	if err != nil {
		return err
	}
	return nil
}
