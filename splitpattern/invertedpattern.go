// Copyright (c) 2020, NLP Odyssey Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package splitpattern

type InvertedPattern struct {
	sp SplitPattern
}

var _ SplitPattern = &InvertedPattern{}

func Invert(sp SplitPattern) *InvertedPattern {
	return &InvertedPattern{sp: sp}
}

func (ip *InvertedPattern) FindMatches(s string) ([]Capture, error) {
	captures, err := ip.sp.FindMatches(s)
	if err != nil {
		return nil, err
	}
	for i, capture := range captures {
		captures[i].IsMatch = !capture.IsMatch
	}
	return captures, nil
}
