// Copyright (c) 2020, NLP Odyssey Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package strutils

// IsRuneBoundary checks that `index`-th byte is the first byte in a UTF-8 code
// point sequence or the end of the string.
//
// The start and end of the string (when `index == len(s)`) are
// considered to be boundaries.
//
// Returns `false` if `index` is greater than `len(s)`.
//
// # Examples
//
// ```
// ```
func IsRuneBoundary(s string, index int) bool {
	// 0 and len are always ok.
	// Test for 0 explicitly so that it can optimize out the check
	// easily and skip reading string data for that case.
	if index == 0 || index == len(s) {
		return true
	}
	if index < 0 || index > len(s) {
		return false
	}
	// This is bit magic equivalent to: s[index] < 128 || s[index] >= 192
	return int8(s[index]) >= -0x40
}

// ByteOffsets represents a (start, end) range of byte offset positions.
type ByteOffsets struct {
	// Start byte position, inclusive.
	Start int
	// End byte position, exclusive.
	End int
}

// RuneOffsets represents a (start, end) range of byte rune positions.
type RuneOffsets struct {
	// Start rune position, inclusive.
	Start int
	// End rune position, exclusive.
	End int
}

type BytesToRuneOffsetConverter struct {
	mapping []int
}

func NewBytesToRuneOffsetConverter(sequence string) *BytesToRuneOffsetConverter {
	mapping := make([]int, 0, len(sequence))
	for runeIndex, r := range sequence {
		for range []byte(string(r)) {
			mapping = append(mapping, runeIndex)
		}
	}
	return &BytesToRuneOffsetConverter{mapping: mapping}
}

func (b *BytesToRuneOffsetConverter) Convert(offsets ByteOffsets) RuneOffsets {
	if len(b.mapping) == 0 {
		return RuneOffsets{Start: 0, End: 0}
	}

	start := b.mapping[offsets.Start]

	var end int
	if offsets.End == len(b.mapping) {
		end = b.mapping[offsets.End-1] + 1
	} else {
		end = b.mapping[offsets.End]
	}

	return RuneOffsets{Start: start, End: end}
}
