// Copyright (c) 2020, NLP Odyssey Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package bpemodel

import (
	"fmt"
	"github.com/nlpodyssey/gotokenizers/models"
	"github.com/nlpodyssey/gotokenizers/vocabulary"
)

var ErrUnknownTokenOutOfVocabulary = fmt.Errorf("the provided unk token is out of vocabulary")

// BPEModel is a Byte Pair Encoding (BPE) model.
//
// See: https://www.aclweb.org/anthology/P16-1162/
type BPEModel struct {
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
	// Whether to fuse multiple unknown tokens
	unknownFusionEnabled bool
}

var _ models.Model = &BPEModel{}

// NewBPEModel returns a new BPEModel initialized with the given options.
func NewBPEModel(
	vocab *vocabulary.Vocabulary,
	merges *MergeMap,
	cacheCapacity int,
	dropout float64,
	unknownToken string,
	continuingSubwordPrefix string,
	endOfWordSuffix string,
	unknownFusionEnabled bool,
) *BPEModel {
	return &BPEModel{
		vocab:                   vocab,
		merges:                  merges,
		cache:                   NewCache(cacheCapacity),
		dropout:                 dropout,
		unknownToken:            unknownToken,
		continuingSubwordPrefix: continuingSubwordPrefix,
		endOfWordSuffix:         endOfWordSuffix,
		unknownFusionEnabled:    unknownFusionEnabled,
	}
}

func NewDefaultBPEModel() *BPEModel {
	return &BPEModel{
		vocab:                   vocabulary.NewVocabulary(),
		merges:                  NewMergeMap(),
		cache:                   NewDefaultCache(),
		dropout:                 0,
		unknownToken:            "",
		continuingSubwordPrefix: "",
		endOfWordSuffix:         "",
		unknownFusionEnabled:    false,
	}
}

func (m *BPEModel) Tokenize(sequence string) ([]models.Token, error) {
	if len(sequence) == 0 {
		return nil, nil
	}

	if !m.hasDropout() {
		return m.tokenizeWithCache(sequence)
	}

	word, err := m.mergeWord(sequence)
	if err != nil {
		return nil, err
	}
	return m.wordToTokens(word)
}

func (m *BPEModel) tokenizeWithCache(sequence string) ([]models.Token, error) {
	hit := m.cache.Get(sequence)
	if hit != nil {
		return m.wordToTokens(hit)
	}

	word, err := m.mergeWord(sequence)
	if err != nil {
		return nil, err
	}

	tokens, err := m.wordToTokens(word)
	if err != nil {
		return nil, err
	}
	m.cache.Set(sequence, word)
	return tokens, nil
}

func (m *BPEModel) mergeWord(w string) (*Word, error) {
	word := NewWordWithCapacity(len(w))

	var unkTokenID int

	if m.hasUnknownToken() {
		var unkTokenExists bool
		unkTokenID, unkTokenExists = m.vocab.GetID(m.unknownToken)
		if !unkTokenExists {
			return nil, ErrUnknownTokenOutOfVocabulary
		}
	}

	var unk *Symbol

	runes := []rune(w)
	lastRuneIndex := len(runes) - 1

	for i, r := range runes {
		s := string(r)
		byteLen := len(s)

		if i == 0 && m.hasContinuingSubwordPrefix() {
			s = m.continuingSubwordPrefix + s
		}
		if i == lastRuneIndex && m.hasEndOfWordSuffix() {
			s = s + m.endOfWordSuffix
		}

		id, foundInVocab := m.vocab.GetID(s)
		if foundInVocab {
			if unk != nil {
				word.Add(unk.ID, unk.Length)
				unk = nil
			}
			word.Add(id, byteLen)
		} else if m.hasUnknownToken() {
			if unk != nil {
				if m.unknownFusionEnabled {
					unk.Length += byteLen
				} else {
					// Do not fuse unk, add the previous one
					word.Add(unk.ID, unk.Length)

					unk.ID = unkTokenID
					unk.Length = byteLen
				}
			} else {
				unk = &Symbol{
					ID:     unkTokenID,
					Length: byteLen,
				}
			}
		}
	}

	if unk != nil {
		word.Add(unk.ID, unk.Length)
	}

	word.MergeAll(m.merges, m.dropout)

	return word, nil
}

func (m *BPEModel) wordToTokens(word *Word) ([]models.Token, error) {
	tokens := make([]models.Token, word.Len())
	offsetStart := 0
	for i, wordSymbol := range *word {
		value, ok := m.vocab.GetString(wordSymbol.ID)
		if !ok {
			return nil, fmt.Errorf("ID %d not found in vocabulary", wordSymbol.ID)
		}
		offsetEnd := offsetStart + wordSymbol.Length
		tokens[i] = models.Token{
			ID:    wordSymbol.ID,
			Value: value,
			Offsets: models.TokenOffsets{
				Start: offsetStart,
				End:   offsetEnd,
			},
		}
		offsetStart = offsetEnd
	}
	return tokens, nil
}

func (m *BPEModel) hasDropout() bool {
	return m.dropout > 0
}

func (m *BPEModel) hasContinuingSubwordPrefix() bool {
	return len(m.continuingSubwordPrefix) != 0
}

func (m *BPEModel) hasEndOfWordSuffix() bool {
	return len(m.endOfWordSuffix) != 0
}

func (m *BPEModel) hasUnknownToken() bool {
	return len(m.unknownToken) != 0
}
