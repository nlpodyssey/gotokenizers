// Copyright (c) 2020, NLP Odyssey Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package metaspacepretokenizer

import (
	"github.com/nlpodyssey/gotokenizers/normalizers/normalizedstring"
	"github.com/nlpodyssey/gotokenizers/pretokenizers"
	"unicode"
)

// MetaSpacePreTokenizer allows the generation of pre-tokens by virtually
// replacing all the whitespace-like characters with the provided
// meta-character (rune) and splitting the string by this character.
//
// A whitespace prefix (' ') can be optionally prepended to the input string,
// unless the first rune of the string is already a unicode whitespace.
type MetaSpacePreTokenizer struct {
	metaCharacter      rune
	prefixSpaceEnabled bool
}

var _ pretokenizers.PreTokenizer = &MetaSpacePreTokenizer{}

// NewMetaSpacePreTokenizer returns a new MetaSpacePreTokenizer, setting
// the meta-character, and the flag which indicates whether to add
// a whitespace prefix.
func NewMetaSpacePreTokenizer(
	metaCharacter rune, prefixSpaceEnabled bool,
) *MetaSpacePreTokenizer {
	return &MetaSpacePreTokenizer{
		metaCharacter:      metaCharacter,
		prefixSpaceEnabled: prefixSpaceEnabled,
	}
}

// DefaultMetaCharacter is the default meta-character (rune) used to
// initialize a DefaultMetaSpacePreTokenizer.
//
// This value is a lower one eighth block U+2581.
const DefaultMetaCharacter = '▁'

// DefaultMetaSpacePreTokenizer returns a new MetaSpacePreTokenizer with
// meta-character set to DefaultMetaCharacter ('▁', i.e. lower one eighth
// block U+2581), and prefix space enabled.
func DefaultMetaSpacePreTokenizer() *MetaSpacePreTokenizer {
	return NewMetaSpacePreTokenizer(DefaultMetaCharacter, true)
}

// PreTokenize virtually replaces all the whitespace-like characters with the
// meta-character and splits the NormalizedString by this character.
//
// If whitespace prefix is enabled, a whitespace (' ') is prepended to
// the NormalizedString, actually modifying its "normalized" value, only if
// the first rune of the string is not already a unicode whitespace.
func (m *MetaSpacePreTokenizer) PreTokenize(
	ns *normalizedstring.NormalizedString,
) ([]pretokenizers.PreToken, error) {
	if m.prefixSpaceEnabled && !startsWithWhitespace(ns.Get()) {
		ns.Prepend(" ")
	}

	tokens := make([]pretokenizers.PreToken, 0)
	word := make([]rune, 0)

	index := 0
	for _, r := range ns.Get() {
		if unicode.In(r, unicode.White_Space) {
			if len(word) > 0 {
				tokens = append(tokens, pretokenizers.PreToken{
					String: string(word),
					Start:  index - len(word),
					End:    index,
				})
				word = word[:0]
			}
			word = append(word, m.metaCharacter)
		} else {
			word = append(word, r)
		}
		index++
	}

	if len(word) > 0 {
		end := ns.Len()
		tokens = append(tokens, pretokenizers.PreToken{
			String: string(word),
			Start:  end - len(word),
			End:    end,
		})
	}

	return tokens, nil
}

func startsWithWhitespace(s string) bool {
	return len(s) != 0 && unicode.In([]rune(s)[0], unicode.White_Space)
}
