// Copyright (c) 2020, NLP Odyssey Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package splitpattern

import (
	"reflect"
	"testing"
)

func runTest(t *testing.T, sp SplitPattern, s string, expected []Capture) {
	t.Helper()

	captures, err := sp.FindMatches(s)
	if err != nil {
		t.Errorf("expected nil error, actual %#v", err)
		return
	}

	if !reflect.DeepEqual(captures, expected) {
		t.Errorf("expected\n  %#v\nactual\n  %#v", expected, captures)
	}
}
