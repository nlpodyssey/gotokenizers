// Copyright (c) 2020, NLP Odyssey Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package bpemodel

import (
	"fmt"
	"github.com/nlpodyssey/gotokenizers/models"
	"github.com/nlpodyssey/gotokenizers/pretokenizers"
	"github.com/nlpodyssey/gotokenizers/vocabulary"
)

var ErrUnknownTokenOutOfVocabulary = fmt.Errorf("the provided unk token is out of vocabulary")

// BpeModel is a Byte Pair Encoding (BPE) model.
//
// See: https://www.aclweb.org/anthology/P16-1162/
type BpeModel struct {
	// Vocabulary of tokens
	vocab *vocabulary.Vocabulary
	// Mapping between symbol pairs and their (rank, new_id).
	merges *MergeMap
	// Cache for optimizing the encoding step.
	cache *WordCache
	// Dropout probability for merges.
	// A value of 0 means no dropout (default).
	// With value 1.0, tokenization will perform no merges, so the result will
	// just be characters.
	// See: https://arxiv.org/abs/1910.13267
	dropout float64
	// The unknown token to be used in the vocabulary when we an unknown
	// token is encountered.
	// Set to empty string to disable.
	unknownToken string
	// An optional prefix to use on any subword that exists only behind
	// another one.
	// Set to empty string to disable.
	continuingSubwordPrefix string
	// An optional suffix to characterize and end-of-word subword.
	// Set to empty string to disable.
	endOfWordSuffix string
}

// NewBpeModel returns a new BpeModel initialized with the given options.
func NewBpeModel(
	vocab *vocabulary.Vocabulary,
	merges *MergeMap,
	cacheCapacity int,
	dropout float64,
	unknownToken string,
	continuingSubwordPrefix string,
	endOfWordSuffix string,
) *BpeModel {
	return &BpeModel{
		vocab:                   vocab,
		merges:                  merges,
		cache:                   NewCache(cacheCapacity),
		dropout:                 dropout,
		unknownToken:            unknownToken,
		continuingSubwordPrefix: continuingSubwordPrefix,
		endOfWordSuffix:         endOfWordSuffix,
	}
}

func NewDefaultBpeModel() *BpeModel {
	return &BpeModel{
		vocab:                   vocabulary.NewVocabulary(),
		merges:                  NewMergeMap(),
		cache:                   NewDefaultCache(),
		dropout:                 0,
		unknownToken:            "",
		continuingSubwordPrefix: "",
		endOfWordSuffix:         "",
	}
}

func (m *BpeModel) Tokenize(sentence []pretokenizers.PreToken) ([]models.Token, error) {
	if len(sentence) == 0 {
		return []models.Token{}, nil
	}

	if m.hasDropout() {
		// If using dropout we don't want to use the cache.
		return m.tokenizeWithoutCache(sentence)
	}

	encoded := make([]models.Token, 0, len(sentence))
	shouldUpdateCache := false
	cacheKeys := m.extractSentenceStrings(sentence)
	cachedWords := m.cache.GetValues(cacheKeys)

	for i, preToken := range sentence {
		word := cachedWords[i]

		if word == nil {
			// No cache hit: re-compute merges and add to cache.
			var err error
			word, err = m.mergeWord(preToken.String)
			if err != nil {
				return nil, err
			}
			cachedWords[i] = word
			shouldUpdateCache = true
		} else {
			// Avoid a possible needless cache update later.
			cachedWords[i] = nil
		}

		tokens, err := m.wordToTokens(i, word, preToken.Start)
		if err != nil {
			return nil, err
		}
		encoded = append(encoded, tokens...)
	}

	// Try updating the cache if we need to.
	if shouldUpdateCache {
		m.cache.SetValues(cacheKeys, cachedWords)
	}

	return encoded, nil
}

func (m *BpeModel) tokenizeWithoutCache(
	sentence []pretokenizers.PreToken,
) ([]models.Token, error) {
	encoded := make([]models.Token, 0, len(sentence))

	for i, preToken := range sentence {
		word, err := m.mergeWord(preToken.String)
		if err != nil {
			return nil, err
		}
		tokens, err := m.wordToTokens(i, word, preToken.Start)
		if err != nil {
			return nil, err
		}
		encoded = append(encoded, tokens...)
	}

	return encoded, nil
}

func (m *BpeModel) hasDropout() bool {
	return m.dropout > 0
}

func (m *BpeModel) hasUnknownToken() bool {
	return len(m.unknownToken) != 0
}

func (m *BpeModel) hasContinuingSubwordPrefix() bool {
	return len(m.continuingSubwordPrefix) != 0
}

func (m *BpeModel) hasEndOfWordSuffix() bool {
	return len(m.endOfWordSuffix) != 0
}

func (m *BpeModel) extractSentenceStrings(sentence []pretokenizers.PreToken) []string {
	s := make([]string, len(sentence))
	for i, preToken := range sentence {
		s[i] = preToken.String
	}
	return s
}

func (m *BpeModel) mergeWord(w string) (*Word, error) {
	word := NewWord()

	hasUnknownToken := m.hasUnknownToken()
	hasContinuingSubwordPrefix := m.hasContinuingSubwordPrefix()
	hasEndOfWordSuffix := m.hasEndOfWordSuffix()

	runes := []rune(w)
	lastRuneIndex := len(runes) - 1

	for i, r := range runes {
		s := string(r)

		if hasContinuingSubwordPrefix && i > 0 {
			s = m.continuingSubwordPrefix + s
		}
		if hasEndOfWordSuffix && i == lastRuneIndex {
			s = s + m.endOfWordSuffix
		}

		id, ok := m.vocab.GetID(s)
		if !ok {
			if !hasUnknownToken {
				continue
			}
			id, ok = m.vocab.GetID(m.unknownToken)
			if !ok {
				return nil, ErrUnknownTokenOutOfVocabulary
			}
		}
		word.Add(id)
	}

	word.MergeAll(m.merges, m.dropout)

	return word, nil
}

func (m *BpeModel) wordToTokens(
	wordIndex int,
	word *Word,
	initialOffset int,
) ([]models.Token, error) {
	tokens := make([]models.Token, word.Len())
	offset := initialOffset
	for i, wordSymbol := range *word {
		value, ok := m.vocab.GetString(wordSymbol.ID)
		if !ok {
			return nil, fmt.Errorf("ID %d not found in vocabulary", wordSymbol.ID)
		}
		tokens[i] = models.Token{
			ID:    wordSymbol.ID,
			Value: value,
			Offsets: models.TokenOffsets{
				Start: offset,
				End:   offset + wordSymbol.Length,
			},
			WordIndex: wordIndex,
		}
		offset += wordSymbol.Length
	}
	return tokens, nil
}
