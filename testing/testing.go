// Copyright (c) 2020, The GoTokenizers Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package testing

import "testing"

func AssertStringEqual(t *testing.T, what, actual, expected string) {
	if actual != expected {
		t.Errorf("Expected %v to be %#v, but got %#v", what, expected, actual)
	}
}

func AssertIntEqual(t *testing.T, what string, actual, expected int) {
	if actual != expected {
		t.Errorf("Expected %v to be %v, but got %v", what, expected, actual)
	}
}
