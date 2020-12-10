// Copyright (c) 2020, NLP Odyssey Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package wordpiecemodel

import (
	"fmt"
	"github.com/nlpodyssey/gotokenizers/models"
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

var _ models.Model = &WordPieceModel{}

func New(
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

func NewDefault() *WordPieceModel {
	return &WordPieceModel{
		vocab:                   vocabulary.NewVocabulary(),
		unknownToken:            "[UNK]",
		continuingSubwordPrefix: "##",
		maxInputCharsPerWord:    100,
	}
}

func (m *WordPieceModel) Tokenize(sequence string) ([]models.Token, error) {
	if len([]rune(sequence)) > m.maxInputCharsPerWord {
		unkTokenID, unkTokenExists := m.vocab.GetID(m.unknownToken)
		if !unkTokenExists {
			return nil, ErrUnknownTokenOutOfVocabulary
		}
		return []models.Token{{
			ID:      unkTokenID,
			Value:   m.unknownToken,
			Offsets: models.TokenOffsets{Start: 0, End: len(sequence)},
		}}, nil
	}

	isBad := false
	start := 0
	subTokens := make([]models.Token, 0)

	for start < len(sequence) {
		end := len(sequence)
		hasCurToken := false
		var curToken models.Token

		for start < end {
			subStr := sequence[start:end]

			if start > 0 {
				subStr = m.continuingSubwordPrefix + subStr
			}

			if id, ok := m.vocab.GetID(subStr); ok {
				hasCurToken = true
				curToken = models.Token{
					ID:      id,
					Value:   subStr,
					Offsets: models.TokenOffsets{Start: start, End: end},
				}
				break
			}

			if len(subStr) > 0 {
				end -= len(string(subStr[len(subStr)-1]))
			} else {
				end -= 1
			}
		}

		if !hasCurToken {
			isBad = true
			break
		}

		subTokens = append(subTokens, curToken)
		start = end
	}

	if isBad {
		unkTokenID, unkTokenExists := m.vocab.GetID(m.unknownToken)
		if !unkTokenExists {
			return nil, ErrUnknownTokenOutOfVocabulary
		}
		return []models.Token{{
			ID:      unkTokenID,
			Value:   m.unknownToken,
			Offsets: models.TokenOffsets{Start: 0, End: len(sequence)},
		}}, nil
	}

	return subTokens, nil
}
