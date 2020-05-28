// Copyright (c) 2020, NLP Odyssey Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package stripnormalizer

import (
	. "github.com/nlpodyssey/gotokenizers/normalizers/normalizedstring"
	"testing"
)

func TestStripNormalizerLeftOnly(t *testing.T) {
	sn := NewStripNormalizer(true, false)
	ns := NewNormalizedString(" \n\tfoo\t\n ")
	err := sn.Normalize(ns)
	if err != nil {
		t.Error(err)
	}
	expected := "foo\t\n "
	if actual := ns.Get(); actual != expected {
		t.Errorf("expected %#v, actual %#v", expected, actual)
	}
}

func TestStripNormalizerRightOnly(t *testing.T) {
	sn := NewStripNormalizer(false, true)
	ns := NewNormalizedString(" \n\tfoo\t\n ")
	err := sn.Normalize(ns)
	if err != nil {
		t.Error(err)
	}
	expected := " \n\tfoo"
	if actual := ns.Get(); actual != expected {
		t.Errorf("expected %#v, actual %#v", expected, actual)
	}
}

func TestStripNormalizerLeftAndRight(t *testing.T) {
	sn := NewStripNormalizer(true, true)
	ns := NewNormalizedString(" \n\tfoo\t\n ")
	err := sn.Normalize(ns)
	if err != nil {
		t.Error(err)
	}
	expected := "foo"
	if actual := ns.Get(); actual != expected {
		t.Errorf("expected %#v, actual %#v", expected, actual)
	}
}

func TestStripNormalizerNoLeftAndNoRight(t *testing.T) {
	sn := NewStripNormalizer(false, false)
	ns := NewNormalizedString(" \n\tfoo\t\n ")
	err := sn.Normalize(ns)
	if err != nil {
		t.Error(err)
	}
	expected := " \n\tfoo\t\n "
	if actual := ns.Get(); actual != expected {
		t.Errorf("expected %#v, actual %#v", expected, actual)
	}
}
