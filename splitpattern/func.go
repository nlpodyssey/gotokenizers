// Copyright (c) 2020, NLP Odyssey Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package splitpattern

import "github.com/nlpodyssey/gotokenizers/strutils"

type FuncSplitPattern struct {
	f func(rune) bool
}

var _ SplitPattern = &FuncSplitPattern{}

func FromFunc(f func(rune) bool) *FuncSplitPattern {
	return &FuncSplitPattern{f: f}
}

func (sp *FuncSplitPattern) FindMatches(s string) ([]Capture, error) {
	if len(s) == 0 {
		return []Capture{{
			Offsets: strutils.ByteOffsets{Start: 0, End: 0},
			IsMatch: false,
		}}, nil
	}

	lastOffset := 0
	lastSeen := 0

	splits := make([]Capture, 0, len(s))

	for i, r := range s {
		runeBytesLen := len(string(r))
		lastSeen = i + runeBytesLen

		if !sp.f(r) {
			continue
		}

		if lastOffset < i {
			// We need to emit what was before this match
			splits = append(splits, Capture{
				Offsets: strutils.ByteOffsets{Start: lastOffset, End: i},
				IsMatch: false,
			})
		}
		splits = append(splits, Capture{
			Offsets: strutils.ByteOffsets{Start: i, End: lastSeen},
			IsMatch: true,
		})
		lastOffset = lastSeen
	}

	// Do not forget the last potential split
	if lastSeen > lastOffset {
		splits = append(splits, Capture{
			Offsets: strutils.ByteOffsets{Start: lastOffset, End: lastSeen},
			IsMatch: false,
		})
	}

	return splits, nil
}
