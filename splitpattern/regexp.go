// Copyright (c) 2020, NLP Odyssey Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package splitpattern

import "regexp"

type RegexpSplitPattern struct {
	r *regexp.Regexp
}

var _ SplitPattern = &RegexpSplitPattern{}

func FromRegexp(r *regexp.Regexp) *RegexpSplitPattern {
	return &RegexpSplitPattern{r: r}
}

func (sp *RegexpSplitPattern) FindMatches(s string) ([]Capture, error) {
	if len(s) == 0 {
		return []Capture{{Offsets: Offsets{Start: 0, End: 0}, IsMatch: false}}, nil
	}

	prev := 0
	splits := make([]Capture, 0, len(s))

	matches := sp.r.FindAllStringIndex(s, -1)
	for _, match := range matches {
		startByte := match[0]
		endByte := match[1]

		if prev != startByte {
			splits = append(splits, Capture{
				Offsets: Offsets{Start: prev, End: startByte},
				IsMatch: false,
			})
		}
		splits = append(splits, Capture{
			Offsets: Offsets{Start: startByte, End: endByte},
			IsMatch: true,
		})
		prev = endByte
	}

	if prev != len(s) {
		splits = append(splits, Capture{
			Offsets: Offsets{Start: prev, End: len(s)},
			IsMatch: false,
		})
	}

	return splits, nil
}
