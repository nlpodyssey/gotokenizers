// Copyright (c) 2020, The GoTokenizers Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package normalized_string

// Represents a range usable by `NormalizedString` to index its content.
type NSRange interface {
	Start() int
	End() int
	Get() (int, int)
	SetStart(start int)
	SetEnd(end int)
	Set(start, end int)
	originalRange(ns *NormalizedString) (int, int, bool)
	normalizedRange(ns *NormalizedString) (int, int, bool)
	convertOffset(ns *NormalizedString) (int, int, bool)
}

type nsRange struct {
	start int
	end   int
}

func (r *nsRange) Start() int         { return r.start }
func (r *nsRange) End() int           { return r.end }
func (r *nsRange) Get() (int, int)    { return r.start, r.end }
func (r *nsRange) SetStart(start int) { r.start = start }
func (r *nsRange) SetEnd(end int)     { r.end = end }
func (r *nsRange) Set(start, end int) { r.start, r.end = start, end }

// A range using indices relative to the `original` string
type NSOriginalRange struct{ nsRange }

var _ NSRange = &NSOriginalRange{}

func NewNSOriginalRange(start, end int) NSOriginalRange {
	return NSOriginalRange{nsRange{start: start, end: end}}
}

func (r *NSOriginalRange) originalRange(ns *NormalizedString) (int, int, bool) {
	runes := []rune(ns.original)
	if r.start > r.end || r.start < 0 || r.end > len(runes) {
		return 0, 0, false
	}

	return r.start, r.end, true
}

func (r *NSOriginalRange) normalizedRange(ns *NormalizedString) (int, int, bool) {
	runes := []rune(ns.original)
	if r.start > r.end || r.start < 0 || r.end > len(runes) {
		return 0, 0, false
	}

	start, end := 0, 0

	for i, alignment := range ns.alignments {
		if r.end < alignment.end {
			break
		}
		end = i + 1
		if alignment.start <= r.start {
			start = i
		}
	}

	return start, end, true
}

func (r *NSOriginalRange) convertOffset(ns *NormalizedString) (int, int, bool) {
	return r.normalizedRange(ns)
}

// A range using indices relative to the `normalized` string
type NSNormalizedRange struct{ nsRange }

var _ NSRange = &NSNormalizedRange{}

func NewNSNormalizedRange(start, end int) NSNormalizedRange {
	return NSNormalizedRange{nsRange{start: start, end: end}}
}
func (r *NSNormalizedRange) originalRange(ns *NormalizedString) (int, int, bool) {
	if r.start > r.end || r.start < 0 || r.end > len(ns.alignments) {
		return 0, 0, false
	}
	if r.start == r.end {
		return 0, 0, true
	}
	return ns.alignments[r.start].start, ns.alignments[r.end-1].end, true
}

func (r *NSNormalizedRange) normalizedRange(ns *NormalizedString) (int, int, bool) {
	if r.start > r.end || r.start < 0 || r.end > len(ns.alignments) {
		return 0, 0, false
	}
	return r.start, r.end, true
}

func (r *NSNormalizedRange) convertOffset(ns *NormalizedString) (int, int, bool) {
	return r.originalRange(ns)
}
