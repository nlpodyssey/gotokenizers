// Copyright (c) 2020, NLP Odyssey Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package models

type Token struct {
	ID        int
	Value     string
	Offsets   TokenOffsets
	WordIndex int
}

type TokenOffsets struct {
	Start int
	End   int
}
