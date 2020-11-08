// Copyright (c) 2020, NLP Odyssey Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package vocabulary

import (
	"encoding/json"
	"io/ioutil"
)

// Vocabulary stores ID-term bidirectional associations.
type Vocabulary struct {
	termToID map[string]int
	idToTerm map[int]string
}

// NewVocabulary returns a new empty vocabulary.
func NewVocabulary() *Vocabulary {
	return &Vocabulary{
		termToID: make(map[string]int),
		idToTerm: make(map[int]string),
	}
}

// FromJSONFile reads a vocabulary from JSON file.
func FromJSONFile(filename string) (*Vocabulary, error) {
	rawData, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var termToID map[string]int
	err = json.Unmarshal(rawData, &termToID)
	if err != nil {
		return nil, err
	}

	idToTerm := make(map[int]string, len(termToID))
	for term, id := range termToID {
		idToTerm[id] = term
	}

	return &Vocabulary{termToID: termToID, idToTerm: idToTerm}, nil
}

// AddTerm adds a new term to the vocabulary.
// If the term does not yet exist in the vocabulary, a new ID is associated to
// it, corresponding to the current vocabulary size.
// Otherwise, if the string already exists, no insertion is performed.
func (v *Vocabulary) AddTerm(term string) {
	if _, ok := v.termToID[term]; ok {
		return
	}
	id := v.Size()
	v.termToID[term] = id
	v.idToTerm[id] = term
}

// Size returns the size of the Vocabulary.
func (v *Vocabulary) Size() int {
	return len(v.termToID)
}

// GetID returns the ID associated to the given string, and whether
// it was found.
func (v *Vocabulary) GetID(str string) (int, bool) {
	id, ok := v.termToID[str]
	return id, ok
}

// GetString returns the string associated to the given ID, and whether
// it was found.
func (v *Vocabulary) GetString(id int) (string, bool) {
	s, ok := v.idToTerm[id]
	return s, ok
}
