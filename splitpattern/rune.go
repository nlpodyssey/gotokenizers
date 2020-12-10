// Copyright (c) 2020, NLP Odyssey Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package splitpattern

type RuneSplitPattern struct {
	r rune
	f *FuncSplitPattern
}

var _ SplitPattern = &RuneSplitPattern{}

func FromRune(r rune) *RuneSplitPattern {
	return &RuneSplitPattern{
		r: r,
		f: FromFunc(func(other rune) bool {
			return other == r
		}),
	}
}

func (sp *RuneSplitPattern) FindMatches(s string) ([]Capture, error) {
	return sp.f.FindMatches(s)
}
