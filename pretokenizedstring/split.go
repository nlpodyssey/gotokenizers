// Copyright (c) 2020, NLP Odyssey Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package pretokenizedstring

import (
	"github.com/nlpodyssey/gotokenizers/models"
	"github.com/nlpodyssey/gotokenizers/normalizedstring"
)

// Split is a wrapper for a subpart of a NormalizedString.
//
// This Split contains the underlying NormalizedString as well as its offsets
// in the original string. These offsets are in the "original" referential.
// It also contains any Token associated to the current split.
type Split struct {
	// The underlying normalizedstring.NormalizedString.
	// Each SubString is represented by a normalizedstring.NormalizedString
	// and in the end we might be carrying a lot of SubString representing
	// various parts of the original input string.
	NormalizedString *normalizedstring.NormalizedString
	// Optional Tokens associated to this Split.
	Tokens *[]models.Token
}

type OriginalByteSplit struct {
	// A slice of the normalized string
	String string
	// The associated bytes offsets, in the original referential
	Offsets normalizedstring.Offsets
	// The potential tokens
	Tokens *[]models.Token
}

type NormalizedByteSplit struct {
	// A slice of the normalized string
	String string
	// The associated bytes offsets, in the normalized referential
	Offsets normalizedstring.Offsets
	// The potential tokens
	Tokens *[]models.Token
}

// SplitsFromNormalizedStrings transforms a slice of NormalizedStrings
// into a corresponding slice of Splits, with nil tokens.
func SplitsFromNormalizedStrings(nss []*normalizedstring.NormalizedString) []Split {
	splits := make([]Split, len(nss))
	for i, ns := range nss {
		splits[i] = Split{NormalizedString: ns, Tokens: nil}
	}
	return splits
}
