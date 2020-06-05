// Copyright (c) 2020, NLP Odyssey Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package bertnormalizer

import (
	. "github.com/nlpodyssey/gotokenizers/normalizers/normalizedstring"
	"testing"
)

func TestDefaultBertNormalizer(t *testing.T) {
	sn := DefaultBertNormalizer()
	ns := NewNormalizedString("(\n\t) (\a\b) (咖啡) (o\u0302) (BAR)")
	err := sn.Normalize(ns)
	if err != nil {
		t.Error(err)
	}
	expected := "(  ) () ( 咖  啡 ) (o) (bar)"
	if actual := ns.Get(); actual != expected {
		t.Errorf("expected %#v, actual %#v", expected, actual)
	}
}
func TestBertNormalizerWithAllFlagsEnabled(t *testing.T) {
	sn := NewBertNormalizer(true, true, true, true)
	ns := NewNormalizedString("(\n\t) (\a\b) (咖啡) (o\u0302) (BAR)")
	err := sn.Normalize(ns)
	if err != nil {
		t.Error(err)
	}
	expected := "(  ) () ( 咖  啡 ) (o) (bar)"
	if actual := ns.Get(); actual != expected {
		t.Errorf("expected %#v, actual %#v", expected, actual)
	}
}

func TestBertNormalizerWithAllFlagsDisabled(t *testing.T) {
	sn := NewBertNormalizer(false, false, false, false)
	ns := NewNormalizedString("(\n\t) (\a\b) (咖啡) (o\u0302) (BAR)")
	err := sn.Normalize(ns)
	if err != nil {
		t.Error(err)
	}
	expected := "(\n\t) (\a\b) (咖啡) (o\u0302) (BAR)"
	if actual := ns.Get(); actual != expected {
		t.Errorf("expected %#v, actual %#v", expected, actual)
	}
}

func TestBertNormalizerWithTextCleaningOnly(t *testing.T) {
	sn := NewBertNormalizer(true, false, false, false)
	ns := NewNormalizedString("(\n\t) (\a\b) (咖啡) (o\u0302) (BAR)")
	err := sn.Normalize(ns)
	if err != nil {
		t.Error(err)
	}
	expected := "(  ) () (咖啡) (o\u0302) (BAR)"
	if actual := ns.Get(); actual != expected {
		t.Errorf("expected %#v, actual %#v", expected, actual)
	}
}

func TestBertNormalizerWithChineseCharsHandlingOnly(t *testing.T) {
	sn := NewBertNormalizer(false, true, false, false)
	ns := NewNormalizedString("(\n\t) (\a\b) (咖啡) (o\u0302) (BAR)")
	err := sn.Normalize(ns)
	if err != nil {
		t.Error(err)
	}
	expected := "(\n\t) (\a\b) ( 咖  啡 ) (o\u0302) (BAR)"
	if actual := ns.Get(); actual != expected {
		t.Errorf("expected %#v, actual %#v", expected, actual)
	}
}

func TestBertNormalizerWithAccentsStrippingOnly(t *testing.T) {
	sn := NewBertNormalizer(false, false, true, false)
	ns := NewNormalizedString("(\n\t) (\a\b) (咖啡) (o\u0302) (BAR)")
	err := sn.Normalize(ns)
	if err != nil {
		t.Error(err)
	}
	expected := "(\n\t) (\a\b) (咖啡) (o) (BAR)"
	if actual := ns.Get(); actual != expected {
		t.Errorf("expected %#v, actual %#v", expected, actual)
	}
}

func TestBertNormalizerWithLowerCaseOnly(t *testing.T) {
	sn := NewBertNormalizer(false, false, false, true)
	ns := NewNormalizedString("(\n\t) (\a\b) (咖啡) (o\u0302) (BAR)")
	err := sn.Normalize(ns)
	if err != nil {
		t.Error(err)
	}
	expected := "(\n\t) (\a\b) (咖啡) (o\u0302) (bar)"
	if actual := ns.Get(); actual != expected {
		t.Errorf("expected %#v, actual %#v", expected, actual)
	}
}
