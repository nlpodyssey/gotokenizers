// Copyright (c) 2020, NLP Odyssey Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package metaspacepretokenizer

import (
	"github.com/nlpodyssey/gotokenizers/normalizedstring"
	"github.com/nlpodyssey/gotokenizers/pretokenizedstring"
	"github.com/nlpodyssey/gotokenizers/pretokenizers"
	"github.com/nlpodyssey/gotokenizers/splitpattern"
	"strings"
)

// MetaSpacePreTokenizer allows the generation of pre-tokens by virtually
// replacing all the whitespace-like characters with the provided
// meta-character (rune) and splitting the string by this character.
//
// A whitespace prefix (' ') can be optionally prepended to the input string,
// unless the first rune of the string is already a unicode whitespace.
type MetaSpacePreTokenizer struct {
	replacement        rune
	strReplacement     string
	prefixSpaceEnabled bool
}

var _ pretokenizers.PreTokenizer = &MetaSpacePreTokenizer{}

// DefaultReplacementCharacter is the default meta-character (rune) used to
// initialize a NewDefault.
//
// This value is a lower one eighth block U+2581.
const DefaultReplacementCharacter = '▁'

// New returns a new MetaSpacePreTokenizer.
func New(replacement rune, prefixSpaceEnabled bool) *MetaSpacePreTokenizer {
	return &MetaSpacePreTokenizer{
		replacement:        replacement,
		strReplacement:     string(replacement),
		prefixSpaceEnabled: prefixSpaceEnabled,
	}
}

// NewDefault returns a new MetaSpacePreTokenizer with
// meta-character set to DefaultReplacementCharacter ('▁', i.e. lower one eighth
// block U+2581), and prefix space enabled.
func NewDefault() *MetaSpacePreTokenizer {
	return New(DefaultReplacementCharacter, true)
}

// PreTokenize virtually replaces all the whitespace-like characters with the
// meta-character and splits the NormalizedString by this character.
//
// If whitespace prefix is enabled, a whitespace (' ') is prepended to
// the NormalizedString, actually modifying its "normalized" value, only if
// the first rune of the string is not already a unicode whitespace.
func (m *MetaSpacePreTokenizer) PreTokenize(pts *pretokenizedstring.PreTokenizedString) error {
	splittingPattern := splitpattern.FromRune(m.replacement)
	return pts.Split(
		func(_ int, ns *normalizedstring.NormalizedString) ([]pretokenizedstring.Split, error) {
			if m.prefixSpaceEnabled && !strings.HasPrefix(ns.Get(), m.strReplacement) {
				ns.Prepend(m.strReplacement)
			}
			err := ns.Replace(splitpattern.FromRune(' '), m.strReplacement)
			if err != nil {
				return nil, err
			}
			nss, err := ns.Split(splittingPattern, normalizedstring.SplitDelimiterMergedWithNext)
			if err != nil {
				return nil, err
			}
			return pretokenizedstring.SplitsFromNormalizedStrings(nss), nil
		},
	)
}
