// Copyright (c) 2020, NLP Odyssey Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package normalizers

import . "github.com/nlpodyssey/gotokenizers/normalizers/normalizedstring"

// Normalizer is implemented by any value that has a Normalize method,
// which takes care of pre-processing strings.
type Normalizer interface {
	Normalize(ns *NormalizedString) error
}
