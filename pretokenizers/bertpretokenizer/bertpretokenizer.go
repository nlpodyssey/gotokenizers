// Copyright (c) 2020, NLP Odyssey Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package bertpretokenizer

import (
	"github.com/nlpodyssey/gotokenizers/normalizers/normalizedstring"
	"github.com/nlpodyssey/gotokenizers/pretokenizers"
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

// NewBertPreTokenizer returns a new BertPreTokenizer.
func NewBertPreTokenizer() *BertPreTokenizer {
	return &BertPreTokenizer{}
}

// PreTokenize splits the NormalizedString into pre-tokens suitable for BERT
// models.
func (rd *BertPreTokenizer) PreTokenize(
	ns *normalizedstring.NormalizedString,
) ([]pretokenizers.PreToken, error) {
	tokens := make([]pretokenizers.PreToken, 0)
	word := make([]rune, 0)

	index := 0
	for _, r := range ns.Get() {
		switch {
		case unicode.In(r, unicode.White_Space):
			if len(word) > 0 {
				tokens = append(tokens, pretokenizers.PreToken{
					String: string(word),
					Start:  index - len(word),
					End:    index,
				})
				word = word[:0]
			}
		case unicode.In(r, unicode.Punct):
			if len(word) > 0 {
				tokens = append(tokens, pretokenizers.PreToken{
					String: string(word),
					Start:  index - len(word),
					End:    index,
				})
				word = word[:0]
			}
			tokens = append(tokens, pretokenizers.PreToken{
				String: string(r),
				Start:  index,
				End:    index + 1,
			})
		default:
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
