// Copyright (c) 2020, NLP Odyssey Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package splitpattern

import "regexp"

type StringSplitPattern struct {
	s string
	r *RegexpSplitPattern
}

var _ SplitPattern = &StringSplitPattern{}

func FromString(s string) *StringSplitPattern {
	sp := &StringSplitPattern{s: s}
	if len(s) > 0 {
		sp.r = FromRegexp(regexp.MustCompile(regexp.QuoteMeta(s)))
	}
	return sp
}

func (sp *StringSplitPattern) FindMatches(s string) ([]Capture, error) {
	if sp.r == nil {
		// If we try to find the matches with an empty string, just don't match anything
		return []Capture{{
			// FIXME: is len of runes (and not bytes) correct?
			Offsets: Offsets{0, len([]rune(s))},
			IsMatch: false,
		}}, nil
	}
	return sp.r.FindMatches(s)
}
