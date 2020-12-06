// Copyright (c) 2020, NLP Odyssey Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package splitpattern

// SplitPattern is implemented by any value which represents a pattern
// for splitting a string.
type SplitPattern interface {
	// FindMatches slices the given string in a list of Capture values.
	//
	// This method MUST cover the whole string in its outputs, with
	// contiguous ordered slices.
	FindMatches(string) ([]Capture, error)
}

// Capture is a single pattern capture of text within the input string,
// which provides offset positions and a flag reporting whether this is a
// match or not.
type Capture struct {
	Offsets Offsets
	IsMatch bool
}

// Offsets represents a (start, end) range of offset positions.
type Offsets struct {
	// Start byte position, inclusive.
	Start int
	// End byte position, exclusive.
	End int
}
