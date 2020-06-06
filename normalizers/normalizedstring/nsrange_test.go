// Copyright (c) 2020, NLP Odyssey Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package normalizedstring

import (
	"testing"
)

func TestNewNSOriginalRange(t *testing.T) {
	t.Parallel()

	r := NewNSOriginalRange(1, 2)
	expected := NSOriginalRange{baseNsRange{start: 1, end: 2}}
	if *r != expected {
		t.Errorf("Expected %+v, got %+v", expected, r)
	}
}

func TestNSOriginalRangeStart(t *testing.T) {
	t.Parallel()

	r := NewNSOriginalRange(1, 2)
	if r.Start() != 1 {
		t.Errorf("expected Start() 1, actual %v", r.Start())
	}
}

func TestNSOriginalRangeEnd(t *testing.T) {
	t.Parallel()

	r := NewNSOriginalRange(1, 2)
	if r.End() != 2 {
		t.Errorf("expected End() 1, actual %v", r.End())
	}
}

func TestNSOriginalRangeGet(t *testing.T) {
	t.Parallel()

	r := NewNSOriginalRange(1, 2)
	start, end := r.Get()
	if start != 1 {
		t.Errorf("expected Get() start 1, actual %v", start)
	}
	if end != 2 {
		t.Errorf("expected Get() end 2, actual %v", end)
	}
}

func TestNSOriginalSetStart(t *testing.T) {
	t.Parallel()

	r := NewNSOriginalRange(2, 3)
	r.SetStart(1)
	expected := NSOriginalRange{baseNsRange{start: 1, end: 3}}
	if *r != expected {
		t.Errorf("Expected %+v, got %+v", expected, r)
	}
}

func TestNSOriginalSetEnd(t *testing.T) {
	t.Parallel()

	r := NewNSOriginalRange(2, 3)
	r.SetEnd(5)
	expected := NSOriginalRange{baseNsRange{start: 2, end: 5}}
	if *r != expected {
		t.Errorf("Expected %+v, got %+v", expected, r)
	}
}

func TestNSOriginalSet(t *testing.T) {
	t.Parallel()

	r := NewNSOriginalRange(2, 3)
	r.Set(1, 5)
	expected := NSOriginalRange{baseNsRange{start: 1, end: 5}}
	if *r != expected {
		t.Errorf("Expected %+v, got %+v", expected, r)
	}
}

func TestNewNSNormalizedRange(t *testing.T) {
	t.Parallel()

	r := NewNSNormalizedRange(1, 2)
	expected := NSNormalizedRange{baseNsRange{start: 1, end: 2}}
	if *r != expected {
		t.Errorf("Expected %+v, got %+v", expected, r)
	}
}

func TestNSNormalizedRangeStart(t *testing.T) {
	t.Parallel()

	r := NewNSNormalizedRange(1, 2)
	if r.Start() != 1 {
		t.Errorf("expected Start() 1, actual %v", r.Start())
	}
}

func TestNSNormalizedRangeEnd(t *testing.T) {
	t.Parallel()

	r := NewNSNormalizedRange(1, 2)
	if r.End() != 2 {
		t.Errorf("expected End() 1, actual %v", r.End())
	}
}

func TestNSNormalizedRangeGet(t *testing.T) {
	t.Parallel()

	r := NewNSNormalizedRange(1, 2)
	start, end := r.Get()
	if start != 1 {
		t.Errorf("expected Get() start 1, actual %v", start)
	}
	if end != 2 {
		t.Errorf("expected Get() end 2, actual %v", end)
	}
}

func TestNSNormalizedSetStart(t *testing.T) {
	t.Parallel()

	r := NewNSNormalizedRange(2, 3)
	r.SetStart(1)
	expected := NSNormalizedRange{baseNsRange{start: 1, end: 3}}
	if *r != expected {
		t.Errorf("Expected %+v, got %+v", expected, r)
	}
}

func TestNSNormalizedSetEnd(t *testing.T) {
	t.Parallel()

	r := NewNSNormalizedRange(2, 3)
	r.SetEnd(5)
	expected := NSNormalizedRange{baseNsRange{start: 2, end: 5}}
	if *r != expected {
		t.Errorf("Expected %+v, got %+v", expected, r)
	}
}

func TestNSNormalizedSet(t *testing.T) {
	t.Parallel()

	r := NewNSNormalizedRange(2, 3)
	r.Set(1, 5)
	expected := NSNormalizedRange{baseNsRange{start: 1, end: 5}}
	if *r != expected {
		t.Errorf("Expected %+v, got %+v", expected, r)
	}
}
