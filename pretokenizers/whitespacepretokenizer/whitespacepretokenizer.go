// Copyright (c) 2020, NLP Odyssey Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package whitespacepretokenizer

import (
	"github.com/nlpodyssey/gotokenizers/normalizers/normalizedstring"
	"github.com/nlpodyssey/gotokenizers/pretokenizers"
	"regexp"
)

// WhiteSpacePreTokenizer allows the generation of pre-tokens made by
// distinct groups of unicode letters (words) and non-letter characters
// (such as punctuation signs or other symbols). Whitespace-lake characters
// are always identified as explicit tokens separators.
type WhiteSpacePreTokenizer struct{}

var _ pretokenizers.PreTokenizer = &WhiteSpacePreTokenizer{}
var tokensPattern *regexp.Regexp = regexp.MustCompile(`\pL+|[^\pL\s]+`)

// NewWhiteSpacePreTokenizer returns a new WhiteSpacePreTokenizer.
func NewWhiteSpacePreTokenizer() *WhiteSpacePreTokenizer {
	return &WhiteSpacePreTokenizer{}
}

// PreTokenize splits the NormalizedString into word and non-word groups
// separated by whitespace-like characters.
func (wt *WhiteSpacePreTokenizer) PreTokenize(
	ns *normalizedstring.NormalizedString,
) ([]pretokenizers.PreToken, error) {
	str := ns.Get()
	matches := tokensPattern.FindAllStringIndex(str, -1)
	tokens := make([]pretokenizers.PreToken, len(matches))

	for index, match := range matches {
		tokens[index] = pretokenizers.PreToken{
			String:    str[match[0]:match[1]],
			ByteStart: match[0],
			ByteEnd:   match[1],
		}
	}

	return tokens, nil
}
