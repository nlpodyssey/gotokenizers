// Copyright (c) 2020, The GoTokenizers Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package normalized_string

import (
	"testing"

	. "github.com/saientist/gotokenizers/testing"
)

func TestNewNSOriginalRange(t *testing.T) {
	r := NewNSOriginalRange(1, 2)
	expected := NSOriginalRange{baseNsRange{start: 1, end: 2}}
	if r != expected {
		t.Errorf("Expected %+v, got %+v", expected, r)
	}
}

func TestNSOriginalRangeStart(t *testing.T) {
	r := NewNSOriginalRange(1, 2)
	AssertIntEqual(t, "Start()", r.Start(), 1)
}

func TestNSOriginalRangeEnd(t *testing.T) {
	r := NewNSOriginalRange(1, 2)
	AssertIntEqual(t, "End()", r.End(), 2)
}

func TestNSOriginalRangeGet(t *testing.T) {
	r := NewNSOriginalRange(1, 2)
	start, end := r.Get()
	AssertIntEqual(t, "Get() start", start, 1)
	AssertIntEqual(t, "Get() end", end, 2)
}

func TestNSOriginalSetStart(t *testing.T) {
	r := NewNSOriginalRange(2, 3)
	r.SetStart(1)
	expected := NSOriginalRange{baseNsRange{start: 1, end: 3}}
	if r != expected {
		t.Errorf("Expected %+v, got %+v", expected, r)
	}
}

func TestNSOriginalSetEnd(t *testing.T) {
	r := NewNSOriginalRange(2, 3)
	r.SetEnd(5)
	expected := NSOriginalRange{baseNsRange{start: 2, end: 5}}
	if r != expected {
		t.Errorf("Expected %+v, got %+v", expected, r)
	}
}

func TestNSOriginalSet(t *testing.T) {
	r := NewNSOriginalRange(2, 3)
	r.Set(1, 5)
	expected := NSOriginalRange{baseNsRange{start: 1, end: 5}}
	if r != expected {
		t.Errorf("Expected %+v, got %+v", expected, r)
	}
}

func TestNewNSNormalizedRange(t *testing.T) {
	r := NewNSNormalizedRange(1, 2)
	expected := NSNormalizedRange{baseNsRange{start: 1, end: 2}}
	if r != expected {
		t.Errorf("Expected %+v, got %+v", expected, r)
	}
}

func TestNSNormalizedRangeStart(t *testing.T) {
	r := NewNSNormalizedRange(1, 2)
	AssertIntEqual(t, "Start()", r.Start(), 1)
}

func TestNSNormalizedRangeEnd(t *testing.T) {
	r := NewNSNormalizedRange(1, 2)
	AssertIntEqual(t, "End()", r.End(), 2)
}

func TestNSNormalizedRangeGet(t *testing.T) {
	r := NewNSNormalizedRange(1, 2)
	start, end := r.Get()
	AssertIntEqual(t, "Get() start", start, 1)
	AssertIntEqual(t, "Get() end", end, 2)
}

func TestNSNormalizedSetStart(t *testing.T) {
	r := NewNSNormalizedRange(2, 3)
	r.SetStart(1)
	expected := NSNormalizedRange{baseNsRange{start: 1, end: 3}}
	if r != expected {
		t.Errorf("Expected %+v, got %+v", expected, r)
	}
}

func TestNSNormalizedSetEnd(t *testing.T) {
	r := NewNSNormalizedRange(2, 3)
	r.SetEnd(5)
	expected := NSNormalizedRange{baseNsRange{start: 2, end: 5}}
	if r != expected {
		t.Errorf("Expected %+v, got %+v", expected, r)
	}
}

func TestNSNormalizedSet(t *testing.T) {
	r := NewNSNormalizedRange(2, 3)
	r.Set(1, 5)
	expected := NSNormalizedRange{baseNsRange{start: 1, end: 5}}
	if r != expected {
		t.Errorf("Expected %+v, got %+v", expected, r)
	}
}
