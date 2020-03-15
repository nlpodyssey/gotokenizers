// Copyright (c) 2020, The GoTokenizers Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package normalizers

import "github.com/saientist/gotokenizers/normalizers/normalized_string"

// Takes care of pre-processing strings.
type Normalizer interface {
	Normalize(normalized *normalized_string.NormalizedString) error
}