// Copyright (c) 2020, NLP Odyssey Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package bytelevelpretokenizer

import (
	"fmt"
	"github.com/nlpodyssey/gotokenizers/pretokenizedstring"
	"github.com/nlpodyssey/gotokenizers/strutils"
	"reflect"
	"testing"
)

func TestByteLevelPreTokenizer_PreTokenize(t *testing.T) {
	t.Parallel()

	t.Run("Prefix space disabled", func(t *testing.T) {
		t.Parallel()

		pt := New(DefaultSplittingRegexp, false, true)
		pts := pretokenizedstring.FromString("Hello my friend, how is your day going?")

		err := pt.PreTokenize(pts)
		if err != nil {
			t.Fatal(err)
		}

		assertEqual(t, pts.GetOriginalByteSplits(), []pretokenizedstring.OriginalByteSplit{
			{String: "Hello", Offsets: strutils.ByteOffsets{Start: 0, End: 5}},
			{String: "Ġmy", Offsets: strutils.ByteOffsets{Start: 5, End: 8}},
			{String: "Ġfriend", Offsets: strutils.ByteOffsets{Start: 8, End: 15}},
			{String: ",", Offsets: strutils.ByteOffsets{Start: 15, End: 16}},
			{String: "Ġhow", Offsets: strutils.ByteOffsets{Start: 16, End: 20}},
			{String: "Ġis", Offsets: strutils.ByteOffsets{Start: 20, End: 23}},
			{String: "Ġyour", Offsets: strutils.ByteOffsets{Start: 23, End: 28}},
			{String: "Ġday", Offsets: strutils.ByteOffsets{Start: 28, End: 32}},
			{String: "Ġgoing", Offsets: strutils.ByteOffsets{Start: 32, End: 38}},
			{String: "?", Offsets: strutils.ByteOffsets{Start: 38, End: 39}},
		})
	})

	t.Run("Prefix space enabled", func(t *testing.T) {
		t.Parallel()

		pt := New(DefaultSplittingRegexp, true, true)

		strings := []string{
			" Hello my friend, how is your day going?",
			"Hello my friend, how is your day going?",
		}

		for _, str := range strings {
			t.Run(fmt.Sprintf("with string %#v", str), func(t *testing.T) {
				pts := pretokenizedstring.FromString(str)

				err := pt.PreTokenize(pts)
				if err != nil {
					t.Fatal(err)
				}

				assertEqual(t, pts.GetNormalizedByteSplits(), []pretokenizedstring.NormalizedByteSplit{
					{String: "ĠHello", Offsets: strutils.ByteOffsets{Start: 0, End: 7}},
					{String: "Ġmy", Offsets: strutils.ByteOffsets{Start: 7, End: 11}},
					{String: "Ġfriend", Offsets: strutils.ByteOffsets{Start: 11, End: 19}},
					{String: ",", Offsets: strutils.ByteOffsets{Start: 19, End: 20}},
					{String: "Ġhow", Offsets: strutils.ByteOffsets{Start: 20, End: 25}},
					{String: "Ġis", Offsets: strutils.ByteOffsets{Start: 25, End: 29}},
					{String: "Ġyour", Offsets: strutils.ByteOffsets{Start: 29, End: 35}},
					{String: "Ġday", Offsets: strutils.ByteOffsets{Start: 35, End: 40}},
					{String: "Ġgoing", Offsets: strutils.ByteOffsets{Start: 40, End: 47}},
					{String: "?", Offsets: strutils.ByteOffsets{Start: 47, End: 48}},
				})
			})
		}
	})

	t.Run("Handling of newlines", func(t *testing.T) {
		t.Parallel()

		pt := New(DefaultSplittingRegexp, false, true)
		pts := pretokenizedstring.FromString("Hello there\nHello there")

		err := pt.PreTokenize(pts)
		if err != nil {
			t.Fatal(err)
		}

		assertEqual(t, pts.GetOriginalByteSplits(), []pretokenizedstring.OriginalByteSplit{
			{String: "Hello", Offsets: strutils.ByteOffsets{Start: 0, End: 5}},
			{String: "Ġthere", Offsets: strutils.ByteOffsets{Start: 5, End: 11}},
			{String: "Ċ", Offsets: strutils.ByteOffsets{Start: 11, End: 12}},
			{String: "Hello", Offsets: strutils.ByteOffsets{Start: 12, End: 17}},
			{String: "Ġthere", Offsets: strutils.ByteOffsets{Start: 17, End: 23}},
		})
	})

	t.Run("Handling of multiple whitespaces", func(t *testing.T) {
		t.Parallel()

		pt := New(DefaultSplittingRegexp, false, true)
		pts := pretokenizedstring.FromString("Hello there       dear")

		err := pt.PreTokenize(pts)
		if err != nil {
			t.Fatal(err)
		}

		assertEqual(t, pts.GetOriginalByteSplits(), []pretokenizedstring.OriginalByteSplit{
			{String: "Hello", Offsets: strutils.ByteOffsets{Start: 0, End: 5}},
			{String: "Ġthere", Offsets: strutils.ByteOffsets{Start: 5, End: 11}},
			{String: "ĠĠĠĠĠĠ", Offsets: strutils.ByteOffsets{Start: 11, End: 17}},
			{String: "Ġdear", Offsets: strutils.ByteOffsets{Start: 17, End: 22}},
		})
	})

	t.Run("Offsets when character splits up", func(t *testing.T) {
		t.Parallel()

		pt := New(DefaultSplittingRegexp, false, true)

		input := "i⭢j"
		pts := pretokenizedstring.FromString(input)

		err := pt.PreTokenize(pts)
		if err != nil {
			t.Fatal(err)
		}

		assertEqual(t, pts.GetOriginalByteSplits(), []pretokenizedstring.OriginalByteSplit{
			{String: "i", Offsets: strutils.ByteOffsets{Start: 0, End: 1}},
			{String: "âŃ¢", Offsets: strutils.ByteOffsets{Start: 1, End: 4}},
			{String: "j", Offsets: strutils.ByteOffsets{Start: 4, End: 5}},
		})

		assertEqual(t, pts.GetNormalizedByteSplits(), []pretokenizedstring.NormalizedByteSplit{
			{String: "i", Offsets: strutils.ByteOffsets{Start: 0, End: 1}},
			{String: "âŃ¢", Offsets: strutils.ByteOffsets{Start: 1, End: 7}},
			{String: "j", Offsets: strutils.ByteOffsets{Start: 7, End: 8}},
		})

		strings := make([]string, 0)
		for _, split := range pts.GetOriginalByteSplits() {
			strings = append(strings, input[split.Offsets.Start:split.Offsets.End])
		}
		assertEqual(t, strings, []string{"i", "⭢", "j"})
	})
}

func assertEqual(t *testing.T, actual, expected interface{}) {
	t.Helper()
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("expected\n  %#v\nactual\n  %#v", expected, actual)
	}
}
