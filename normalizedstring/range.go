// Copyright (c) 2020, NLP Odyssey Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package normalizedstring

// Range represents a range which can be used by NormalizedString to index
// its content.
//
// A Range can use indices relative to either the "original" string
// (see OriginalRange), or the "normalized" string (see NormalizedRange).
type Range interface {
	// Start returns the start index, inclusive.
	Start() int
	// End returns the end index, exclusive.
	End() int
	// Len returns the length of the range (i.e. End - Start).
	Len() int
}

// OriginalRange represents a range usable by the NormalizedString to index
// its content, using indices relative to the "original" string.
type OriginalRange struct {
	start int
	end   int
}

var _ Range = OriginalRange{}

func NewOriginalRange(start, end int) OriginalRange {
	return OriginalRange{start: start, end: end}
}

func (r OriginalRange) Start() int {
	return r.start
}

func (r OriginalRange) End() int {
	return r.end
}

func (r OriginalRange) Len() int {
	return r.end - r.start
}

// NormalizedRange represents a range usable by the NormalizedString to index
// its content, using indices relative to the "normalized" string.
type NormalizedRange struct {
	start int
	end   int
}

var _ Range = NormalizedRange{}

func NewNormalizedRange(start, end int) NormalizedRange {
	return NormalizedRange{start: start, end: end}
}

func (r NormalizedRange) Start() int {
	return r.start
}

func (r NormalizedRange) End() int {
	return r.end
}

func (r NormalizedRange) Len() int {
	return r.end - r.start
}
