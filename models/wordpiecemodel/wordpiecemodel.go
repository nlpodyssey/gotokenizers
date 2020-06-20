// Copyright (c) 2020, NLP Odyssey Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package wordpiecemodel

import (
	"fmt"
	"github.com/nlpodyssey/gotokenizers/models"
	"github.com/nlpodyssey/gotokenizers/pretokenizers"
	"github.com/nlpodyssey/gotokenizers/vocabulary"
)

var ErrUnknownTokenOutOfVocabulary = fmt.Errorf("the provided unk token is out of vocabulary")

// WordPieceModel is a WordPiece model.
//
// See: https://static.googleusercontent.com/media/research.google.com/en//pubs/archive/37842.pdf
type WordPieceModel struct {
	// Vocabulary of (token -> ID) mappings.
	vocab *vocabulary.Vocabulary
	// The unknown token for the vocabulary.
	unknownToken string
	// Prefix for continuing subwords.
	continuingSubwordPrefix string
	// Maximum number of input characters per word.
	maxInputCharsPerWord int
}

func NewWordPieceModel(
	vocab *vocabulary.Vocabulary,
	unknownToken string,
	continuingSubwordPrefix string,
	maxInputCharsPerWord int,
) *WordPieceModel {
	return &WordPieceModel{
		vocab:                   vocab,
		unknownToken:            unknownToken,
		continuingSubwordPrefix: continuingSubwordPrefix,
		maxInputCharsPerWord:    maxInputCharsPerWord,
	}
}

func NewDefaultWordPieceModel() *WordPieceModel {
	return &WordPieceModel{
		vocab:                   vocabulary.NewVocabulary(),
		unknownToken:            "[UNK]",
		continuingSubwordPrefix: "##",
		maxInputCharsPerWord:    100,
	}
}

func (m *WordPieceModel) Tokenize(sentence []pretokenizers.PreToken) ([]models.Token, error) {
	outputTokens := make([]models.Token, 0, len(sentence))

	for index, preToken := range sentence {
		runes := []rune(preToken.String)
		runesLen := len(runes)

		if len(runes) > m.maxInputCharsPerWord {
			unknownTokenID, ok := m.vocab.GetID(m.unknownToken)
			if !ok {
				return nil, ErrUnknownTokenOutOfVocabulary
			}
			outputTokens = append(outputTokens, models.Token{
				ID:    unknownTokenID,
				Value: m.unknownToken,
				Offsets: models.TokenOffsets{
					Start: preToken.Start,
					End:   preToken.End,
				},
				WordIndex: index,
			})
			continue
		}

		isBad := false
		start := 0
		subTokens := make([]models.Token, 0)

		for start < runesLen {
			end := runesLen
			var curToken models.Token
			tokenFound := false

			for ; start < end; end-- {
				subStr := string(runes[start:end])
				if start > 0 {
					subStr = m.continuingSubwordPrefix + subStr
				}

				if id, ok := m.vocab.GetID(subStr); ok {
					curToken = models.Token{
						ID:    id,
						Value: subStr,
						Offsets: models.TokenOffsets{
							Start: preToken.Start + start,
							End:   preToken.Start + end,
						},
						WordIndex: index,
					}
					tokenFound = true
					break
				}
			}

			if !tokenFound {
				isBad = true
				break
			}

			subTokens = append(subTokens, curToken)
			start = end
		}

		if isBad {
			unknownTokenID, ok := m.vocab.GetID(m.unknownToken)
			if !ok {
				return nil, ErrUnknownTokenOutOfVocabulary
			}
			outputTokens = append(outputTokens, models.Token{
				ID:    unknownTokenID,
				Value: m.unknownToken,
				Offsets: models.TokenOffsets{
					Start: preToken.Start,
					End:   preToken.End,
				},
				WordIndex: index,
			})
		} else {
			outputTokens = append(outputTokens, subTokens...)
		}
	}

	return outputTokens, nil
}
