// Copyright (c) 2020, NLP Odyssey Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package sequencenormalizer

import (
	"fmt"
	"github.com/nlpodyssey/gotokenizers/normalizedstring"
	"github.com/nlpodyssey/gotokenizers/normalizers"
	"github.com/nlpodyssey/gotokenizers/normalizers/lowercasenormalizer"
	"github.com/nlpodyssey/gotokenizers/normalizers/stripnormalizer"
	"testing"
)

func TestSequenceNormalizerWithTwoNormalizers(t *testing.T) {
	t.Parallel()

	sn := NewSequenceNormalizer([]normalizers.Normalizer{
		lowercasenormalizer.NewLowerCaseNormalizer(),
		stripnormalizer.NewStripNormalizer(true, true),
	})
	ns := normalizedstring.FromString("  Foo Bar SÜẞ  ")
	err := sn.Normalize(ns)
	if err != nil {
		t.Error(err)
	}
	expected := "foo bar süß"
	if actual := ns.Get(); actual != expected {
		t.Errorf("expected %#v, actual %#v", expected, actual)
	}
}

func TestSequenceNormalizerWithOneNormalizer(t *testing.T) {
	t.Parallel()

	sn := NewSequenceNormalizer([]normalizers.Normalizer{
		stripnormalizer.NewStripNormalizer(true, true),
	})
	ns := normalizedstring.FromString("  foo  ")
	err := sn.Normalize(ns)
	if err != nil {
		t.Error(err)
	}
	expected := "foo"
	if actual := ns.Get(); actual != expected {
		t.Errorf("expected %#v, actual %#v", expected, actual)
	}
}

func TestSequenceNormalizerWithEmptySequence(t *testing.T) {
	t.Parallel()

	sn := NewSequenceNormalizer([]normalizers.Normalizer{})
	ns := normalizedstring.FromString("  foo  ")
	err := sn.Normalize(ns)
	if err != nil {
		t.Error(err)
	}
	expected := "  foo  "
	if actual := ns.Get(); actual != expected {
		t.Errorf("expected %#v, actual %#v", expected, actual)
	}
}

type ErrorNormalizer struct{}

var _ normalizers.Normalizer = &ErrorNormalizer{}

func (sn *ErrorNormalizer) Normalize(_ *normalizedstring.NormalizedString) error {
	return fmt.Errorf("sample error")
}

func TestSequenceNormalizerReturnsTheFirstErrorEncountered(t *testing.T) {
	t.Parallel()

	sn := NewSequenceNormalizer([]normalizers.Normalizer{
		lowercasenormalizer.NewLowerCaseNormalizer(),
		&ErrorNormalizer{},
		stripnormalizer.NewStripNormalizer(true, true),
	})
	ns := normalizedstring.FromString("Foo")
	err := sn.Normalize(ns)
	if err == nil {
		t.Errorf("expected error, actual nil")
	}
}
