// Copyright (c) 2020, NLP Odyssey Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package bytelevelpretokenizer

import (
	"github.com/dlclark/regexp2"
	"github.com/nlpodyssey/gotokenizers/normalizedstring"
	"github.com/nlpodyssey/gotokenizers/pretokenizedstring"
	"github.com/nlpodyssey/gotokenizers/pretokenizers"
	"github.com/nlpodyssey/gotokenizers/splitpattern"
	"unicode"
)

// ByteLevelPreTokenizer allows the generation of pre-tokens suitable in the
// context of BPE (Byte Pair Encoding) models.
//
// This pre-tokenizer splits the input string into tokens according to the
// given regular expression. It also transforms the string expanding each
// UTF-8 character (rune) into bytes, which are then re-mapped to new custom
// runes.
//
// A whitespace prefix (' ') can be optionally prepended to the input string,
// unless the first rune of the string is already a unicode whitespace.
//
// Offsets trimming can be enabled to exclude whitespaces in the post-processing
// step.
type ByteLevelPreTokenizer struct {
	splittingRegexp        *regexp2.Regexp
	prefixSpaceEnabled     bool
	offsetsTrimmingEnabled bool
}

var _ pretokenizers.PreTokenizer = &ByteLevelPreTokenizer{}

// DefaultSplittingRegexp is a simple default regular expression that
// can be used for ByteLevelPreTokenizer.
//
// This MUST be treated as a read-only constant value .
var DefaultSplittingRegexp = regexp2.MustCompile(
	`'s|'t|'re|'ve|'m|'ll|'d| ?\p{L}+| ?\p{N}+| ?[^\s\p{L}\p{N}]+|\s+(?!\S)|\s+`,
	regexp2.IgnoreCase|regexp2.Multiline)

// New returns a new ByteLevelPreTokenizer.
func New(
	splittingRegexp *regexp2.Regexp,
	prefixSpaceEnabled bool,
	offsetsTrimmingEnabled bool,
) *ByteLevelPreTokenizer {
	return &ByteLevelPreTokenizer{
		splittingRegexp:        splittingRegexp,
		prefixSpaceEnabled:     prefixSpaceEnabled,
		offsetsTrimmingEnabled: offsetsTrimmingEnabled,
	}
}

// NewDefault returns a new ByteLevelPreTokenizer, using
// DefaultSplittingRegexp, and enabling prefix space insertion.
func NewDefault() *ByteLevelPreTokenizer {
	return New(DefaultSplittingRegexp, true, true)
}

// PreTokenize is in charge of transforming all the unicode characters into
// their byte-level counterpart. It also splits the input according to the
// configured regex.
func (b *ByteLevelPreTokenizer) PreTokenize(pts *pretokenizedstring.PreTokenizedString) error {
	//let re_ref: &Regex = &RE;
	splittingPattern := splitpattern.FromRegexp2(b.splittingRegexp)

	err := pts.Split(
		func(_ int, ns *normalizedstring.NormalizedString) ([]pretokenizedstring.Split, error) {
			if b.prefixSpaceEnabled && !startsWithWhitespace(ns.Get()) {
				ns.Prepend(" ")
			}
			nss, err := ns.Split(splittingPattern, normalizedstring.SplitDelimiterIsolated)
			if err != nil {
				return nil, err
			}
			return pretokenizedstring.SplitsFromNormalizedStrings(nss), nil
		},
	)
	if err != nil {
		return err
	}

	return pts.Normalize(func(ns *normalizedstring.NormalizedString) error {
		s := ns.Get()
		transformations := make([]normalizedstring.RuneChange, 0, len(s))
		i := 0
		for _, r := range s {
			size := len(string(r))
			bytes := []byte(s[i : i+size])
			i += size

			for byteIndex, byteVal := range bytes {
				change := 0
				if byteIndex > 0 {
					change = 1
				}
				transformations = append(transformations, normalizedstring.RuneChange{
					Rune:   byteToRune[byteVal],
					Change: change,
				})
			}
		}
		ns.Transform(transformations, 0)
		return nil
	})
}

func startsWithWhitespace(s string) bool {
	return len(s) > 0 && unicode.In([]rune(s)[0], unicode.White_Space)
}

var byteToRune [0x100]rune

func init() {
	n := 0
	for i := range byteToRune {
		if (i >= '!' && i <= '~') || (i >= 0xA1 && i <= 0xAC) || (i >= 0xAE && i <= 0xFF) {
			byteToRune[i] = rune(i)
		} else {
			byteToRune[i] = rune(0x100 + n)
			n++
		}
	}
}
