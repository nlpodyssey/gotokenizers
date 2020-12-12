// Copyright (c) 2020, NLP Odyssey Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package encodings

import "github.com/nlpodyssey/gotokenizers/strutils"

// Encoding represents the output of a Tokenizer.
type Encoding struct {
	// IDs produced by the `Tokenizer`
	IDs []int
	// Type of the IDs
	TypeIDs []int
	// Tokens associated to each ID
	Tokens []string
	// Indice of the word associated to each token/ID
	// -1 means missing value
	Words []int
	// Offsets of the token/ID from the NormalizedString
	Offsets []strutils.ByteOffsets
	// Mask identifying special tokens
	SpecialTokensMask []int
	// Mask identifying padding tokens for the attention mechanism
	AttentionMask []int
	// A list of overflowing Encoding generated when we got truncated
	Overflowing []*Encoding
}

// EncodableToken represents a single token, expected to be part of a sequence
// which can be transformed into an Encoding.
type EncodableToken struct {
	ID      int
	Token   string
	Offsets strutils.ByteOffsets
	// -1 for "none"
	WordIndex int
	TypeID    int
}

// NewEncoding builds a new Encoding.
func NewEncoding(
	ids []int,
	typeIDs []int,
	tokens []string,
	words []int,
	offsets []strutils.ByteOffsets,
	specialTokensMask []int,
	attentionMask []int,
	overflowing []*Encoding,
) *Encoding {
	return &Encoding{
		IDs:               ids,
		TypeIDs:           typeIDs,
		Tokens:            tokens,
		Words:             words,
		Offsets:           offsets,
		SpecialTokensMask: specialTokensMask,
		AttentionMask:     attentionMask,
		Overflowing:       overflowing,
	}
}

// NewDefaultEncoding builds a new Encoding with empty data.
func NewDefaultEncoding() *Encoding {
	return NewEncodingWithCapacity(0)
}

// NewEncodingWithCapacity builds a new Encoding with empty data.
func NewEncodingWithCapacity(c int) *Encoding {
	return &Encoding{
		IDs:               make([]int, 0, c),
		TypeIDs:           make([]int, 0, c),
		Tokens:            make([]string, 0, c),
		Words:             make([]int, 0, c),
		Offsets:           make([]strutils.ByteOffsets, 0, c),
		SpecialTokensMask: make([]int, 0, c),
		AttentionMask:     make([]int, 0, c),
		Overflowing:       make([]*Encoding, 0),
	}
}

func EncodingFromEncodableTokens(tokens []EncodableToken) *Encoding {
	encoding := NewEncodingWithCapacity(len(tokens))
	for _, token := range tokens {
		encoding.IDs = append(encoding.IDs, token.ID)
		encoding.Tokens = append(encoding.Tokens, token.Token)
		encoding.Offsets = append(encoding.Offsets, token.Offsets)
		encoding.TypeIDs = append(encoding.TypeIDs, token.TypeID)
		encoding.Words = append(encoding.Words, token.WordIndex)
		encoding.SpecialTokensMask = append(encoding.SpecialTokensMask, 0)
		encoding.AttentionMask = append(encoding.AttentionMask, 1)
	}
	return encoding
}

// IsEmpty reports whether this Encoding is empty.
func (e *Encoding) IsEmpty() bool {
	return len(e.IDs) == 0
}

// Len returns the total length of this Encoding.
func (e *Encoding) Len() int {
	return len(e.IDs)
}
