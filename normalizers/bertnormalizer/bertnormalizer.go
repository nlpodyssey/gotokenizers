// Copyright (c) 2020, NLP Odyssey Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package bertnormalizer

import (
	"github.com/nlpodyssey/gotokenizers/normalizers"
	"github.com/nlpodyssey/gotokenizers/normalizers/normalizedstring"
	"unicode"
)

// BertNormalizer allows string normalizations especially suitable for
// BERT-based models.
type BertNormalizer struct {
	// Whether to do the BERT basic cleaning, replacing whitespace-like
	// characters with simple whitespaces, and removing control characters.
	textCleaning bool
	// Whether to put spaces around Chinese characters, so that they can be
	// split.
	chineseCharsHandling bool
	// Whether to strip accents.
	accentsStripping bool
	// Whether to lowercase the input.
	lowerCaseEnabled bool
}

var _ normalizers.Normalizer = &BertNormalizer{}

// NewBertNormalizer returns a new BertNormalizer.
func NewBertNormalizer(
	textCleaning, chineseCharsHandling, accentsStripping, lowerCaseEnabled bool,
) *BertNormalizer {
	return &BertNormalizer{
		textCleaning:         textCleaning,
		chineseCharsHandling: chineseCharsHandling,
		accentsStripping:     accentsStripping,
		lowerCaseEnabled:     lowerCaseEnabled,
	}
}

// DefaultBertNormalizer returns a new BertNormalizer with all
// normalizations enabled.
func DefaultBertNormalizer() *BertNormalizer {
	return NewBertNormalizer(true, true, true, true)
}

// Normalize transform the NormalizedString in place.
func (sn *BertNormalizer) Normalize(ns *normalizedstring.NormalizedString) error {
	if sn.textCleaning {
		sn.cleanText(ns)
	}
	if sn.chineseCharsHandling {
		sn.handleChineseChars(ns)
	}
	if sn.accentsStripping {
		sn.stripAccents(ns)
	}
	if sn.lowerCaseEnabled {
		ns.ToLower()
	}
	return nil
}

// cleanText replaces whitespace-like characters with simple whitespaces, and
// removes control characters.
func (sn *BertNormalizer) cleanText(ns *normalizedstring.NormalizedString) {
	// Since '\t', '\n', and '\r' are also control characters, it's important
	// to first map them to simple whitespaces, and only after that apply the
	// filtering.
	ns.Map(mapWhiteSpace)
	ns.Filter(isNotControlCharacter)
}

// isNotControlCharacter reports whether the given rune is not a control
// character (or a similar control-like character).
//
// Control characters should be excluded from the normalized string, when text
// cleaning is enabled.
func isNotControlCharacter(r rune) bool {
	return r != unicode.ReplacementChar && !unicode.In(r, unicode.Other)
}

// mapWhiteSpace returns a simple whitespace rune (' ') if the given rune is
// a whitespace-like character; any other rune is returned as-is.
func mapWhiteSpace(r rune) rune {
	if unicode.In(r, unicode.White_Space) {
		return ' '
	}
	return r
}

// chineseCharacters defines sets of chinese characters.
//
// A "chinese character" is defined as anything in the CJK Unicode block:
// https://en.wikipedia.org/wiki/CJK_Unified_Ideographs_(Unicode_block)
//
// Note that the CJK Unicode block does NOT include all Japanese and Korean
// characters, despite its name.
// The modern Korean Hangul alphabet is a different block, as well as Japanese
// Hiragana and Katakana. Those alphabets are used to write space-separated
// words, so they are not treated specially, and handled like all the other
// languages.
var chineseCharacters = &unicode.RangeTable{
	R16: []unicode.Range16{
		{0x3400, 0x4dbf, 1},
		{0x4e00, 0x9fff, 1},
		{0xf900, 0xfaff, 1},
	},
	R32: []unicode.Range32{
		{0x20000, 0x2a6df, 1},
		{0x2a700, 0x2b73f, 1},
		{0x2b740, 0x2b81f, 1},
		{0x2b920, 0x2ceaf, 1},
		{0x2f800, 0x2fa1f, 1},
	},
}

// handleChineseChars puts spaces around Chinese characters, so that they can
// be split. All non-Chinese characters are left unchanged.
func (sn *BertNormalizer) handleChineseChars(ns *normalizedstring.NormalizedString) {
	runeChanges := make([]normalizedstring.RuneChanges, 0, ns.Len())
	for _, r := range ns.Get() {
		if unicode.In(r, chineseCharacters) {
			runeChanges = append(
				runeChanges,
				normalizedstring.RuneChanges{Rune: ' ', Changes: 1},
				normalizedstring.RuneChanges{Rune: r, Changes: 0},
				normalizedstring.RuneChanges{Rune: ' ', Changes: 1},
			)
		} else {
			runeChanges = append(runeChanges, normalizedstring.RuneChanges{Rune: r, Changes: 0})
		}
	}
	ns.Transform(runeChanges, 0)
}

// stripAccents removes accent characters (Mn: Mark, non-spacing) from the
// normalized string.
func (sn *BertNormalizer) stripAccents(ns *normalizedstring.NormalizedString) {
	// TODO: ns.Nfd()
	ns.Filter(isNotMarkNonSpacing)
}

// isNotMarkNonSpacing reports whether the given rune is not a Unicode
// character in category Mn (Mark, non-spacing)
func isNotMarkNonSpacing(r rune) bool {
	return !unicode.In(r, unicode.Mn)
}
