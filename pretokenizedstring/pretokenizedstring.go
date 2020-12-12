// Copyright (c) 2020, NLP Odyssey Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package pretokenizedstring

import (
	"fmt"
	"github.com/nlpodyssey/gotokenizers/encodings"
	"github.com/nlpodyssey/gotokenizers/models"
	"github.com/nlpodyssey/gotokenizers/normalizedstring"
	"github.com/nlpodyssey/gotokenizers/strutils"
)

// PreTokenizedString is in charge of splitting an underlying string,
// making sure everything is fine while doing so, and providing ways to
// normalize and tokenize these splits.
//
// Once everything has been normalized and tokenized, the PreTokenizedString
// is able to build an Encoding with all the relevant offsets and word ids,
// relative to the original string.
type PreTokenizedString struct {
	original string
	splits   []Split
}

func FromString(s string) *PreTokenizedString {
	return FromNormalizedString(normalizedstring.FromString(s))
}

func FromNormalizedString(ns *normalizedstring.NormalizedString) *PreTokenizedString {
	return &PreTokenizedString{
		original: ns.GetOriginal(),
		splits: []Split{
			{
				NormalizedString: ns,
				Tokens:           nil,
			},
		},
	}
}

// SplitFunc (used by PreTokenizedString.Split) takes a
// normalizedstring.NormalizedString and is in charge of returning an iterator
// over the produced normalizedstring.NormalizedString.
//
//SplitFunc is free of modifying these NormalizedString as relevant, as long
// as it respects the constraint stated below.
//
// There is only one constraint that MUST be respected: the produced
// normalizedstring.NormalizedString, if combined back together, must have the
// same "original" string as the original one given to SplitFunc.
// This concretely means that, for the offset tracking to work as expected,
// SplitFunc must produce "splits" of the original string.
type SplitFunc func(
	index int,
	ns *normalizedstring.NormalizedString,
) ([]Split, error)

// Split splits the PreTokenizedString by providing a SplitFunc in charge of
// splitting each substring (normalizedstring.NormalizedString) into multiple
// parts.
func (p *PreTokenizedString) Split(splitFunc SplitFunc) error {
	// newSplits is at least as big as p.splits
	newSplits := make([]Split, 0, len(p.splits))

	for index, originalSplit := range p.splits {
		if originalSplit.Tokens != nil {
			newSplits = append(newSplits, originalSplit)
			continue
		}

		items, err := splitFunc(index, originalSplit.NormalizedString)
		if err != nil {
			return err
		}

		for _, item := range items {
			if item.NormalizedString.IsEmpty() {
				continue
			}
			newSplits = append(newSplits, item)
		}
	}

	p.splits = newSplits
	return nil
}

// Normalize normalizes all the splits that do not have attached Split.Tokens,
// using the provided normalization function.
func (p *PreTokenizedString) Normalize(
	normalize func(ns *normalizedstring.NormalizedString) error,
) error {
	for _, split := range p.splits {
		if split.Tokens != nil {
			continue
		}
		err := normalize(split.NormalizedString)
		if err != nil {
			return err
		}
	}
	return nil
}

// Tokenize tokenizes all the splits that do not have attached Split.Tokens,
// using the provided tokenization function.
func (p *PreTokenizedString) Tokenize(
	tokenize func(ns *normalizedstring.NormalizedString) ([]models.Token, error),
) error {
	for i, split := range p.splits {
		if split.Tokens != nil {
			continue
		}
		tokens, err := tokenize(split.NormalizedString)
		if err != nil {
			return err
		}
		p.splits[i].Tokens = &tokens
	}
	return nil
}

// GetOriginalByteSplits returns a list of OriginalByteSplit.
func (p *PreTokenizedString) GetOriginalByteSplits() []OriginalByteSplit {
	result := make([]OriginalByteSplit, len(p.splits))
	for i, split := range p.splits {
		result[i] = OriginalByteSplit{
			String:  split.NormalizedString.Get(),
			Offsets: split.NormalizedString.OriginalOffsets(),
			Tokens:  split.Tokens,
		}
	}
	return result
}

// GetNormalizedByteSplits returns a list of NormalizedByteSplit.
func (p *PreTokenizedString) GetNormalizedByteSplits() []NormalizedByteSplit {
	result := make([]NormalizedByteSplit, len(p.splits))
	offset := 0
	for i, split := range p.splits {
		start := offset
		offset += split.NormalizedString.Len()

		result[i] = NormalizedByteSplit{
			String: split.NormalizedString.Get(),
			Offsets: strutils.ByteOffsets{
				Start: start,
				End:   offset,
			},
			Tokens: split.Tokens,
		}
	}
	return result
}

func (p *PreTokenizedString) Splits() []Split {
	return p.splits
}

// IntoEncoding transforms the current PreTokenizedString into an
// encodings.Encoding.
//
// If a wordIndex is provided (i.e. >= 0), any word in the generated Encoding
// will be set to this value. This is generally used with pre-tokenized
// input, that does not need the PreTokenizedString to generate word ids.
//
// This method will fail if some splits do not have associated Token.
//
// Offset indices are based on bytes (not runes).
func (p *PreTokenizedString) IntoEncoding(wordIndex int, typeID int) (*encodings.Encoding, error) {
	if len(p.splits) == 0 {
		return encodings.NewDefaultEncoding(), nil
	}
	if !p.allSplitsHaveTokens() {
		return nil, fmt.Errorf("splits have not been tokenized, call `PreTokenizedString.Tokenize` first")
	}

	sequence := make([]encodings.EncodableToken, 0)

	for splitIndex, split := range p.splits {
		nsOffsets := split.NormalizedString.OriginalOffsets()

		actualWordIndex := wordIndex
		if actualWordIndex < 0 {
			actualWordIndex = splitIndex
		}

		for _, token := range *split.Tokens {
			var offsets strutils.ByteOffsets

			tokenOrigRange, ok := split.NormalizedString.CoerceRangeToOriginal(
				normalizedstring.NewNormalizedRange(token.Offsets.Start, token.Offsets.End))
			if ok {
				offsets = strutils.ByteOffsets{
					Start: nsOffsets.Start + tokenOrigRange.Start(),
					End:   nsOffsets.Start + tokenOrigRange.End(),
				}
			} else {
				offsets = token.Offsets
			}

			sequence = append(sequence, encodings.EncodableToken{
				ID:        token.ID,
				Token:     token.Value,
				Offsets:   offsets,
				WordIndex: actualWordIndex,
				TypeID:    typeID,
			})
		}
	}

	return encodings.EncodingFromEncodableTokens(sequence), nil
}

func (p *PreTokenizedString) allSplitsHaveTokens() bool {
	for _, split := range p.splits {
		if split.Tokens == nil {
			return false
		}
	}
	return true
}
