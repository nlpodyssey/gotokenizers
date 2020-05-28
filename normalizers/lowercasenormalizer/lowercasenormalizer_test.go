// Copyright (c) 2020, NLP Odyssey Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package lowercasenormalizer

import (
	. "github.com/nlpodyssey/gotokenizers/normalizers/normalizedstring"
	"testing"
)

func TestLowerCaseNormalizer(t *testing.T) {
	sn := NewLowerCaseNormalizer()
	ns := NewNormalizedString("Foo Bar SÜẞ")
	err := sn.Normalize(ns)
	if err != nil {
		t.Error(err)
	}
	expected := "foo bar süß"
	if actual := ns.Get(); actual != expected {
		t.Errorf("expected %#v, actual %#v", expected, actual)
	}
}
