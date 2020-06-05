// Copyright (c) 2020, NLP Odyssey Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package testing

import "testing"

// AssertStringEqual causes a failure if the actual string value is
// different from the expected one.
func AssertStringEqual(t *testing.T, what, actual, expected string) {
	if actual != expected {
		t.Errorf("Expected %v to be %#v, but got %#v", what, expected, actual)
	}
}

// AssertIntEqual causes a failure if the actual int value is
// different from the expected one.
func AssertIntEqual(t *testing.T, what string, actual, expected int) {
	if actual != expected {
		t.Errorf("Expected %v to be %v, but got %v", what, expected, actual)
	}
}

// AssertBoolEqual causes a failure if the actual bool value is
// different from the expected one.
func AssertBoolEqual(t *testing.T, what string, actual, expected bool) {
	if actual != expected {
		t.Errorf("Expected %v to be %v, but got %v", what, expected, actual)
	}
}

// AssertRuneSliceEqual causes a failure if the actual []rune value is
// different from the expected one.
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

// AssertPanic causes a failure if the callback invocation does not cause
// a panic.
func AssertPanic(t *testing.T, what string, callback func()) {
	defer func() {
		if recover() == nil {
			t.Errorf("%s was expected to panic, but recover() is nil", what)
		}
	}()
	callback()
}
