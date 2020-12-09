// Copyright (c) 2020, NLP Odyssey Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package pretokenizers

import "github.com/nlpodyssey/gotokenizers/pretokenizedstring"

// PreTokenizer is implemented by any value that has a PreTokenize method,
// which takes care of performing a pre-segmentation step.
//
// Pre-tokenization splits the given string into multiple substrings, keeping
// track of the offsets between the original string and the substrings.
// In some occasions, the NormalizedString might be modified.
type PreTokenizer interface {
	PreTokenize(pts *pretokenizedstring.PreTokenizedString) error
}

// PreToken represents a pre-tokenized substring, along with its offsets
// position on the original string.
type PreToken struct {
	// The pre-tokenized substring
	String string
	// Start rune position on the original string, inclusive
	Start int
	// End rune position on the original string, exclusive
	End int
}
