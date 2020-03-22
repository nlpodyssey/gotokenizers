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

func AssertBoolEqual(t *testing.T, what string, actual, expected bool) {
	if actual != expected {
		t.Errorf("Expected %v to be %v, but got %v", what, expected, actual)
	}
}

func AssertRuneSliceEqual(t *testing.T, what string, actual, expected []rune) {
	if len(actual) != len(expected) {
		t.Errorf(
			"Expected %v to be %v (%#v), but got %v (%#v) (lengths differ)",
			what, expected, string(expected), actual, string(actual))
		return
	}
	for index, actualItem := range actual {
		expectedItem := expected[index]
		if actualItem != expectedItem {
			t.Errorf(
				"Expected %v to be %v (%#v), but got %v (%#v)"+
					" (mismatch at index %v: expected %v %#v, actual %v %#v)",
				what, expected, string(expected), actual, string(actual),
				index, expectedItem, string(expectedItem),
				actualItem, string(actualItem))
		}
	}
}
