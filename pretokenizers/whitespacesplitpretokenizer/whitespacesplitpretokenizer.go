// Copyright (c) 2020, NLP Odyssey Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package whitespacesplitpretokenizer

import (
	"github.com/nlpodyssey/gotokenizers/normalizers/normalizedstring"
	"github.com/nlpodyssey/gotokenizers/pretokenizers"
	"unicode"
)

// WhiteSpaceSplitPreTokenizer allows the generation of pre-tokens splitting
// by whitespace-lake characters.
type WhiteSpaceSplitPreTokenizer struct{}

var _ pretokenizers.PreTokenizer = &WhiteSpaceSplitPreTokenizer{}

// NewWhiteSpaceSplitPreTokenizer returns a new WhiteSpaceSplitPreTokenizer.
func NewWhiteSpaceSplitPreTokenizer() *WhiteSpaceSplitPreTokenizer {
	return &WhiteSpaceSplitPreTokenizer{}
}

// PreTokenize splits the NormalizedString by whitespace-like characters
func (w *WhiteSpaceSplitPreTokenizer) PreTokenize(
	ns *normalizedstring.NormalizedString,
) ([]pretokenizers.PreToken, error) {
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
