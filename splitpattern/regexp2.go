// Copyright (c) 2020, NLP Odyssey Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package splitpattern

import "github.com/dlclark/regexp2"

type Regexp2SplitPattern struct {
	r *regexp2.Regexp
}

var _ SplitPattern = &Regexp2SplitPattern{}

func FromRegexp2(r *regexp2.Regexp) *Regexp2SplitPattern {
	return &Regexp2SplitPattern{r: r}
}

func (sp *Regexp2SplitPattern) FindMatches(s string) ([]Capture, error) {
	if len(s) == 0 {
		return []Capture{{Offsets: Offsets{Start: 0, End: 0}, IsMatch: false}}, nil
	}

	runes := []rune(s)
	prev := 0
	splits := make([]Capture, 0, len(s))

	match, err := sp.r.FindStringMatch(s)
	if err != nil {
		return nil, err
	}

	for match != nil {
		startByte := len(string(runes[:match.Index]))
		endByte := startByte + len(match.String())

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

		match, err = sp.r.FindNextMatch(match)
		if err != nil {
			return nil, err
		}
	}

	if prev != len(s) {
		splits = append(splits, Capture{
			Offsets: Offsets{Start: prev, End: len(s)},
			IsMatch: false,
		})
	}

	return splits, nil
}
