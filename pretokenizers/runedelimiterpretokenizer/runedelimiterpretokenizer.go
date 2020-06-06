// Copyright (c) 2020, NLP Odyssey Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package runedelimiterpretokenizer

import (
	"github.com/nlpodyssey/gotokenizers/normalizers/normalizedstring"
	"github.com/nlpodyssey/gotokenizers/pretokenizers"
)

// RuneDelimiterPreTokenizer allows the generation of pre-tokens splitting
// by a specific rune.
type RuneDelimiterPreTokenizer struct {
	delimiter rune
}

var _ pretokenizers.PreTokenizer = &RuneDelimiterPreTokenizer{}

// NewRuneDelimiterPreTokenizer returns a new RuneDelimiterPreTokenizer,
// setting the given rune as delimiter.
func NewRuneDelimiterPreTokenizer(delimiter rune) *RuneDelimiterPreTokenizer {
	return &RuneDelimiterPreTokenizer{delimiter: delimiter}
}

// PreTokenize splits the NormalizedString by rune delimiter.
func (rd *RuneDelimiterPreTokenizer) PreTokenize(
	ns *normalizedstring.NormalizedString,
) ([]pretokenizers.PreToken, error) {
	tokens := make([]pretokenizers.PreToken, 0)
	word := make([]rune, 0)

	str := string(ns.Get())
	index := 0
	for _, r := range str {
		if r == rd.delimiter {
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
