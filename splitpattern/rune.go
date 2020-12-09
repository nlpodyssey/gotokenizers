// Copyright (c) 2020, NLP Odyssey Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package splitpattern

type RuneSplitPattern struct {
	r rune
}

var _ SplitPattern = &RuneSplitPattern{}

func FromRune(r rune) *RuneSplitPattern {
	return &RuneSplitPattern{r: r}
}

func (sp *RuneSplitPattern) FindMatches(s string) ([]Capture, error) {
	return FromFunc(func(r rune) bool {
		return r == sp.r
	}).FindMatches(s)
}
