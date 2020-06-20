// Copyright (c) 2020, NLP Odyssey Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package vocabulary

import "testing"

func TestNewVocabulary(t *testing.T) {
	t.Parallel()

	v := NewVocabulary()
	if v == nil || len(v.idToTerm) != 0 || len(v.termToID) != 0 {
		t.Errorf("expected empty *Vocabulary, actual %#v", v)
	}
}
