// Copyright (c) 2020, NLP Odyssey Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package vocabulary

import (
	"testing"
)

func TestNewVocabulary(t *testing.T) {
	t.Parallel()

	v := NewVocabulary()
	if v == nil || len(v.idToTerm) != 0 || len(v.termToID) != 0 {
		t.Errorf("expected empty *Vocabulary, actual %#v", v)
	}
}

func TestFromJSONFile(t *testing.T) {
	t.Parallel()

	v, err := FromJSONFile("testdata/vocab.json")
	if err != nil {
		t.Fatal(err)
	}
	if v.Size() != 3 {
		t.Errorf("expected Size() == 3, actual %d", v.Size())
	}

	values := map[string]int{
		"foo": 0,
		"bar": 1,
		"baz": 2,
	}
	for term, id := range values {
		if i, b := v.GetID(term); !b || i != id {
			t.Errorf(" expected GetID(%#v) == (%d, true), actual (%d, %t)", term, id, i, b)
		}
		if s, b := v.GetString(id); !b || s != term {
			t.Errorf(" expected GetString(%d) == (%#v, true), actual (%#v, %t)", id, term, s, b)
		}
	}
}
