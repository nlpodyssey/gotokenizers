// Copyright (c) 2020, NLP Odyssey Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package bertpretokenizer

import (
	"github.com/nlpodyssey/gotokenizers/normalizedstring"
	"github.com/nlpodyssey/gotokenizers/pretokenizedstring"
	"github.com/nlpodyssey/gotokenizers/strutils"
	"reflect"
	"testing"
)

func TestBertPreTokenizer_PreTokenize(t *testing.T) {
	t.Parallel()

	t.Run("Basic", func(t *testing.T) {
		pt := New()
		pts := pretokenizedstring.FromString("Hey friend!     How are you?!?")
		err := pt.PreTokenize(pts)
		if err != nil {
			t.Fatal(err)
		}

		assertEqual(t, pts.GetOriginalByteSplits(), []pretokenizedstring.OriginalByteSplit{
			{String: "Hey", Offsets: strutils.ByteOffsets{Start: 0, End: 3}},
			{String: "friend", Offsets: strutils.ByteOffsets{Start: 4, End: 10}},
			{String: "!", Offsets: strutils.ByteOffsets{Start: 10, End: 11}},
			{String: "How", Offsets: strutils.ByteOffsets{Start: 16, End: 19}},
			{String: "are", Offsets: strutils.ByteOffsets{Start: 20, End: 23}},
			{String: "you", Offsets: strutils.ByteOffsets{Start: 24, End: 27}},
			{String: "?", Offsets: strutils.ByteOffsets{Start: 27, End: 28}},
			{String: "!", Offsets: strutils.ByteOffsets{Start: 28, End: 29}},
			{String: "?", Offsets: strutils.ByteOffsets{Start: 29, End: 30}},
		})
	})

	t.Run("Chinese characters", func(t *testing.T) {
		pt := New()

		ns := normalizedstring.FromString("野口里佳 Noguchi Rika")
		runeChanges := make([]normalizedstring.RuneChange, 0, ns.Len())
		for _, r := range ns.Get() {
			if r > 0x4E00 {
				runeChanges = append(runeChanges,
					normalizedstring.RuneChange{Rune: ' ', Change: 0},
					normalizedstring.RuneChange{Rune: r, Change: 1},
					normalizedstring.RuneChange{Rune: ' ', Change: 1},
				)
			} else {
				runeChanges = append(runeChanges,
					normalizedstring.RuneChange{Rune: r, Change: 0})
			}
		}
		ns.Transform(runeChanges, 0)

		pts := pretokenizedstring.FromNormalizedString(ns)

		err := pt.PreTokenize(pts)
		if err != nil {
			t.Fatal(err)
		}

		assertEqual(t, pts.GetOriginalByteSplits(), []pretokenizedstring.OriginalByteSplit{
			{String: "野", Offsets: strutils.ByteOffsets{Start: 0, End: 3}},
			{String: "口", Offsets: strutils.ByteOffsets{Start: 3, End: 6}},
			{String: "里", Offsets: strutils.ByteOffsets{Start: 6, End: 9}},
			{String: "佳", Offsets: strutils.ByteOffsets{Start: 9, End: 12}},
			{String: "Noguchi", Offsets: strutils.ByteOffsets{Start: 13, End: 20}},
			{String: "Rika", Offsets: strutils.ByteOffsets{Start: 21, End: 25}},
		})
	})
}

func assertEqual(t *testing.T, actual, expected interface{}) {
	t.Helper()
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("expected\n  %#v\nactual\n  %#v", expected, actual)
	}
}
