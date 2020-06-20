// Copyright (c) 2020, NLP Odyssey Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package bpemodel

import (
	"reflect"
	"testing"
)

func TestSymbolMergeWith(t *testing.T) {
	a := WordSymbol{Symbol: Symbol{ID: 1, Length: 2}, Prev: 3, Next: 4}
	b := &WordSymbol{Symbol: Symbol{ID: 10, Length: 20}, Prev: 30, Next: 40}

	a.MergeWith(b, 100)

	exp := WordSymbol{Symbol: Symbol{ID: 100, Length: 22}, Prev: 3, Next: 40}
	if !reflect.DeepEqual(a, exp) {
		t.Errorf("expected %#v, actual %#v", exp, a)
	}
}

func TestNewWord(t *testing.T) {
	w := NewWord()
	if w == nil || len(*w) != 0 {
		t.Errorf("expected empty *Word, actual %#v", w)
	}
}

func TestWordAdd(t *testing.T) {
	w := NewWord()

	w.Add(11)
	expected := Word{
		&WordSymbol{Symbol: Symbol{ID: 11, Length: 1}, Prev: -1, Next: -1},
	}
	if !reflect.DeepEqual(*w, expected) {
		t.Errorf("expected %#v, actual %#v", expected, *w)
	}

	w.Add(22)
	expected = Word{
		&WordSymbol{Symbol: Symbol{ID: 11, Length: 1}, Prev: -1, Next: 1},
		&WordSymbol{Symbol: Symbol{ID: 22, Length: 1}, Prev: 0, Next: -1},
	}
	if !reflect.DeepEqual(*w, expected) {
		t.Errorf("expected %#v, actual %#v", expected, *w)
	}

	w.Add(33)
	expected = Word{
		&WordSymbol{Symbol: Symbol{ID: 11, Length: 1}, Prev: -1, Next: 1},
		&WordSymbol{Symbol: Symbol{ID: 22, Length: 1}, Prev: 0, Next: 2},
		&WordSymbol{Symbol: Symbol{ID: 33, Length: 1}, Prev: 1, Next: -1},
	}
	if !reflect.DeepEqual(*w, expected) {
		t.Errorf("expected %#v, actual %#v", expected, *w)
	}
}
