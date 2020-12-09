// Copyright (c) 2020, NLP Odyssey Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package pretokenizedstring

import "github.com/nlpodyssey/gotokenizers/normalizedstring"

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

/// Normalize normalizes all the splits that do not have attached Split.Tokens,
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

// GetOriginalByteSplits returns a list of NormalizedByteSplit.
func (p *PreTokenizedString) GetNormalizedByteSplits() []NormalizedByteSplit {
	result := make([]NormalizedByteSplit, len(p.splits))
	offset := 0
	for i, split := range p.splits {
		start := offset
		offset += split.NormalizedString.Len()

		result[i] = NormalizedByteSplit{
			String: split.NormalizedString.Get(),
			Offsets: normalizedstring.Offsets{
				Start: start,
				End:   offset,
			},
			Tokens: split.Tokens,
		}
	}
	return result
}
