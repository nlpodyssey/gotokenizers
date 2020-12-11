// Copyright (c) 2020, NLP Odyssey Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package models

import "github.com/nlpodyssey/gotokenizers/strutils"

// Model represents a model used during Tokenization (like BPE or Word or Unigram).
type Model interface {
	// Tokenize tokenizes the given sequence into multiple underlying Tokens.
	// The Token.Offsets are expected to be relative to the given sequence.
	Tokenize(sequence string) ([]Token, error)
}

type Token struct {
	ID      int
	Value   string
	Offsets strutils.ByteOffsets
}
