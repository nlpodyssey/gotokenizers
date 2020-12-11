// Copyright (c) 2020, NLP Odyssey Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package strutils

import (
	"reflect"
	"testing"
)

func TestIsRuneBoundary(t *testing.T) {
	s := "Löwe 老虎 Léopard"

	assertEqual(t, IsRuneBoundary(s, 0), true)
	assertEqual(t, IsRuneBoundary(s, len(s)), true)

	assertEqual(t, IsRuneBoundary(s, -1), false)
	assertEqual(t, IsRuneBoundary(s, 1000), false)

	// start of `老`
	assertEqual(t, IsRuneBoundary(s, 6), true)

	// second byte of `ö`
	assertEqual(t, IsRuneBoundary(s, 2), false)

	// third byte of `老`
	assertEqual(t, IsRuneBoundary(s, 8), false)
}

func assertEqual(t *testing.T, actual, expected interface{}) {
	t.Helper()
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("expected\n  %#v\nactual\n  %#v", expected, actual)
	}
}
