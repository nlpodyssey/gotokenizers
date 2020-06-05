// Copyright (c) 2020, NLP Odyssey Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package normalizedstring

// NSRange is implemented by values which represent a range usable by the
// NormalizedString to index its content.
//
// This common interface is implemented by both NSOriginalRange and
// NSNormalizedRange, and it contributes avoiding some duplicate code in the
// implementation.
type NSRange interface {
	Start() int
	End() int
	Get() (start, end int)
	SetStart(start int)
	SetEnd(end int)
	Set(start, end int)
	originalRange(ns *NormalizedString) (start, end int, ok bool)
	normalizedRange(ns *NormalizedString) (start, end int, ok bool)
	convertOffset(ns *NormalizedString) (start, end int, ok bool)
}

// baseNsRange is the internal common struct shared by both NSOriginalRange and
// NSNormalizedRange.
//
// It implements common fields and methods, avoiding duplicate code.
type baseNsRange struct{ start, end int }

// Start returns the lower bound (inclusive) of the range.
func (r *baseNsRange) Start() int { return r.start }

// End returns the upper bound (exclusive) of the range.
func (r *baseNsRange) End() int { return r.end }

// Get returns the both the upper bound (exclusive) and the upper bound
// (exclusive) of the range.
func (r *baseNsRange) Get() (start, end int) { return r.start, r.end }

// SetStart sets the lower bound (inclusive) of the range.
func (r *baseNsRange) SetStart(start int) { r.start = start }

// SetEnd sets the upper bound (exclusive) of the range.
func (r *baseNsRange) SetEnd(end int) { r.end = end }

// Set sets the both the upper bound (exclusive) and the upper bound
// (exclusive) of the range.
func (r *baseNsRange) Set(start, end int) { r.start, r.end = start, end }

// NSOriginalRange identifies a (start, end) range of indices on the
// "original" value of a NormalizedString.
type NSOriginalRange struct{ baseNsRange }

var _ NSRange = &NSOriginalRange{}

// NewNSOriginalRange returns a new NSOriginalRange initialized with
// the given (start, end) values.
func NewNSOriginalRange(start, end int) *NSOriginalRange {
	return &NSOriginalRange{baseNsRange{start: start, end: end}}
}

// originalRange returns the range bounds of the "original" string, that is,
// the very same bounds of the NSOriginalRange.
// If the range is out of bounds, then it returns (-1, -1, false).
func (r *NSOriginalRange) originalRange(ns *NormalizedString) (start, end int, ok bool) {
	runes := []rune(ns.original)
	if r.start > r.end || r.start < 0 || r.end > len(runes) {
		return -1, -1, false
	}
	return r.start, r.end, true
}

// normalizedRange returns the range bounds of the "normalized" string,
// remapped from the bounds of the "original" string.
// If the range is out of bounds, then it returns (-1, -1, false).
func (r *NSOriginalRange) normalizedRange(ns *NormalizedString) (start, end int, ok bool) {
	runes := []rune(ns.original)
	if r.start > r.end || r.start < 0 || r.end > len(runes) {
		return -1, -1, false
	}

	start, end = 0, 0

	for i, alignment := range ns.alignments {
		if alignment.start == alignment.end {
			continue
		}
		if r.end < alignment.end {
			break
		}
		end = i + 1
		if alignment.start == r.start {
			start = i
		} else if alignment.start < r.start {
			start = i + 1
		}
	}

	return start, end, true
}

// convertOffset is an alias for "NSOriginalRange.normalizedRange".
func (r *NSOriginalRange) convertOffset(ns *NormalizedString) (start, end int, ok bool) {
	return r.normalizedRange(ns)
}

// NSNormalizedRange identifies a (start, end) range of indices on the
// "normalized" value of a NormalizedString.
type NSNormalizedRange struct{ baseNsRange }

var _ NSRange = &NSNormalizedRange{}

// NewNSNormalizedRange returns a new NSOriginalRange initialized with
// the given (start, end) values.
func NewNSNormalizedRange(start, end int) *NSNormalizedRange {
	return &NSNormalizedRange{baseNsRange{start: start, end: end}}
}

// originalRange returns the range bounds of the "original" string,
// remapped from the bounds of the "normalized" string.
// If the range is out of bounds, then it returns (-1, -1, false).
func (r *NSNormalizedRange) originalRange(ns *NormalizedString) (start, end int, ok bool) {
	if r.start > r.end || r.start < 0 || r.end > len(ns.alignments) {
		return -1, -1, false
	}
	if r.start == r.end {
		return 0, 0, true
	}
	return ns.alignments[r.start].start, ns.alignments[r.end-1].end, true
}

// normalizedRange returns the range bounds of the "normalized" string,
// that is,  the very same bounds of the NSNormalizedRange.
// If the range is out of bounds, then it returns (-1, -1, false).
func (r *NSNormalizedRange) normalizedRange(ns *NormalizedString) (start, end int, ok bool) {
	if r.start > r.end || r.start < 0 || r.end > len(ns.alignments) {
		return -1, -1, false
	}
	return r.start, r.end, true
}

// convertOffset is an alias for "NSNormalizedRange.originalRange".
func (r *NSNormalizedRange) convertOffset(ns *NormalizedString) (start, end int, ok bool) {
	return r.originalRange(ns)
}
