// Copyright (c) 2020, NLP Odyssey Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package normalizedstring

import (
	"github.com/nlpodyssey/gotokenizers/splitpattern"
	"github.com/nlpodyssey/gotokenizers/strutils"
	"reflect"
	"regexp"
	"testing"
	"unicode"
)

func TestNew(t *testing.T) {
	t.Parallel()

	ns := New(
		"Foo",
		"Bar",
		[]AlignmentRange{{0, 1}, {1, 2}, {2, 3}},
		42,
	)

	if ns == nil {
		t.Fatal("expected *NormalizedString, actual nil")
	}

	expected := NormalizedString{
		original:      "Foo",
		normalized:    "Bar",
		alignments:    []AlignmentRange{{0, 1}, {1, 2}, {2, 3}},
		originalShift: 42,
	}

	assertEqual(t, *ns, expected)
}

func TestFromString(t *testing.T) {
	t.Parallel()

	assertEqual(t, FromString(""), New(
		"",
		"",
		[]AlignmentRange{},
		0,
	))

	assertEqual(t, FromString("a"), New(
		"a",
		"a",
		[]AlignmentRange{{0, 1}},
		0,
	))

	assertEqual(t, FromString("Foo"), New(
		"Foo",
		"Foo",
		[]AlignmentRange{{0, 1}, {1, 2}, {2, 3}},
		0,
	))

	assertEqual(t, FromString("ß"), New(
		"ß",
		"ß",
		[]AlignmentRange{{0, 2}, {0, 2}},
		0,
	))

	assertEqual(t, FromString("ß"), New(
		"ß",
		"ß",
		[]AlignmentRange{{0, 2}, {0, 2}},
		0,
	))

	assertEqual(t, FromString("ℝ"), New(
		"ℝ",
		"ℝ",
		[]AlignmentRange{{0, 3}, {0, 3}, {0, 3}},
		0,
	))

	assertEqual(t, FromString("💣"), New(
		"💣",
		"💣",
		[]AlignmentRange{{0, 4}, {0, 4}, {0, 4}, {0, 4}},
		0,
	))

	assertEqual(t, FromString("aßz"), New(
		"aßz",
		"aßz",
		[]AlignmentRange{{0, 1}, {1, 3}, {1, 3}, {3, 4}},
		0,
	))

	assertEqual(t, FromString("ßxß"), New(
		"ßxß",
		"ßxß",
		[]AlignmentRange{{0, 2}, {0, 2}, {2, 3}, {3, 5}, {3, 5}},
		0,
	))
}

func TestNormalizedStringGet(t *testing.T) {
	t.Parallel()

	ns := New(
		"Foo",
		"Bar",
		[]AlignmentRange{{0, 1}, {1, 2}, {2, 3}},
		42,
	)
	assertEqual(t, ns.Get(), "Bar")
}

func TestNormalizedStringGetOriginal(t *testing.T) {
	t.Parallel()

	ns := New(
		"Foo",
		"Bar",
		[]AlignmentRange{{0, 1}, {1, 2}, {2, 3}},
		42,
	)
	assertEqual(t, ns.GetOriginal(), "Foo")
}

func TestNormalizedStringLen(t *testing.T) {
	t.Parallel()

	ns := New(
		"Foo",
		"Ba",
		[]AlignmentRange{{0, 1}, {1, 2}},
		42,
	)
	assertEqual(t, ns.Len(), 2)

	assertEqual(t, FromString("ℝ").Len(), 3)
}

func TestNormalizedStringOriginalLen(t *testing.T) {
	t.Parallel()

	ns := New(
		"Foo",
		"Ba",
		[]AlignmentRange{{0, 1}, {1, 2}},
		42,
	)
	assertEqual(t, ns.OriginalLen(), 3)

	assertEqual(t, FromString("ℝ").OriginalLen(), 3)
}

func TestNormalizedStringIsEmpty(t *testing.T) {
	t.Parallel()

	ns := New("x", "", []AlignmentRange{}, 42)
	assertEqual(t, ns.IsEmpty(), true)

	ns = New("", "x", []AlignmentRange{}, 42)
	assertEqual(t, ns.IsEmpty(), false)
}

func TestNormalizedStringOriginalOffsets(t *testing.T) {
	t.Parallel()

	ns := New("", "x", []AlignmentRange{}, 0)
	assertEqual(t, ns.OriginalOffsets(), strutils.ByteOffsets{Start: 0, End: 0})

	ns = New("", "x", []AlignmentRange{}, 42)
	assertEqual(t, ns.OriginalOffsets(), strutils.ByteOffsets{Start: 42, End: 42})

	ns = New("Foo", "", []AlignmentRange{}, 0)
	assertEqual(t, ns.OriginalOffsets(), strutils.ByteOffsets{Start: 0, End: 3})

	ns = New("Foo", "", []AlignmentRange{}, 42)
	assertEqual(t, ns.OriginalOffsets(), strutils.ByteOffsets{Start: 42, End: 45})

	ns = New("ℝ", "", []AlignmentRange{}, 0)
	assertEqual(t, ns.OriginalOffsets(), strutils.ByteOffsets{Start: 0, End: 3})

	ns = New("ℝ", "", []AlignmentRange{}, 42)
	assertEqual(t, ns.OriginalOffsets(), strutils.ByteOffsets{Start: 42, End: 45})
}

func TestNormalizedStringPrepend(t *testing.T) {
	t.Parallel()

	ns := FromString("there")
	ns.Prepend("Hey ")
	assertEqual(t, ns, New(
		"there",
		"Hey there",
		[]AlignmentRange{
			{0, 1},
			{0, 1},
			{0, 1},
			{0, 1},
			{0, 1},
			{1, 2},
			{2, 3},
			{3, 4},
			{4, 5},
		},
		0,
	))

	or, ok := ns.CoerceRangeToOriginal(NewNormalizedRange(0, 4))
	assertEqual(t, ok, true)
	assertEqual(t, or, NewOriginalRange(0, 1))
}

func TestNormalizedStringAppend(t *testing.T) {
	t.Parallel()

	ns := FromString("Hey")
	ns.Append(" there")
	assertEqual(t, ns, New(
		"Hey",
		"Hey there",
		[]AlignmentRange{
			{0, 1},
			{1, 2},
			{2, 3},
			{2, 3},
			{2, 3},
			{2, 3},
			{2, 3},
			{2, 3},
			{2, 3},
		},
		0,
	))

	or, ok := ns.CoerceRangeToOriginal(NewNormalizedRange(3, len(" there")))
	assertEqual(t, ok, true)
	assertEqual(t, or, NewOriginalRange(2, 3))
}

func TestNormalizedStringFilter(t *testing.T) {
	t.Parallel()

	t.Run("Remove runes", func(t *testing.T) {
		t.Parallel()

		ns := FromString("élégant")
		ns.Filter(func(r rune) bool {
			return r != 'n'
		})
		assertEqual(t, ns, New(
			"élégant",
			"élégat",
			[]AlignmentRange{
				{0, 2},
				{0, 2},
				{2, 3},
				{3, 5},
				{3, 5},
				{5, 6},
				{6, 7},
				// Skipped range
				{8, 9},
			},
			0,
		))
		assertEqual(t, ns.OriginalAlignments(), []AlignmentRange{
			{0, 2},
			{0, 2},
			{2, 3},
			{3, 5},
			{3, 5},
			{5, 6},
			{6, 7},
			{7, 7}, // Eaten n
			{7, 8},
		})
	})

	t.Run("Remove at beginning", func(t *testing.T) {
		t.Parallel()

		ns := FromString("     Hello")
		ns.Filter(func(r rune) bool {
			return !unicode.In(r, unicode.White_Space)
		})
		assertEqual(t, ns.Get(), "Hello")

		r, ok := ns.GetOriginalRange(NewNormalizedRange(1, len("Hello")))
		assertEqual(t, ok, true)
		assertEqual(t, r, "ello")

		r, ok = ns.GetOriginalRange(NewNormalizedRange(0, ns.Len()))
		assertEqual(t, ok, true)
		assertEqual(t, r, "Hello")
	})

	t.Run("Remove at end", func(t *testing.T) {
		t.Parallel()

		ns := FromString("Hello    ")
		ns.Filter(func(r rune) bool {
			return !unicode.In(r, unicode.White_Space)
		})
		assertEqual(t, ns.Get(), "Hello")

		r, ok := ns.GetOriginalRange(NewNormalizedRange(0, 4))
		assertEqual(t, ok, true)
		assertEqual(t, r, "Hell")

		r, ok = ns.GetOriginalRange(NewNormalizedRange(0, ns.Len()))
		assertEqual(t, ok, true)
		assertEqual(t, r, "Hello")
	})

	t.Run("Remove around both edges", func(t *testing.T) {
		t.Parallel()

		ns := FromString("  Hello  ")
		ns.Filter(func(r rune) bool {
			return !unicode.In(r, unicode.White_Space)
		})
		assertEqual(t, ns.Get(), "Hello")

		r, ok := ns.GetOriginalRange(NewNormalizedRange(0, len("Hello")))
		assertEqual(t, ok, true)
		assertEqual(t, r, "Hello")

		r, ok = ns.GetOriginalRange(NewNormalizedRange(0, len("Hell")))
		assertEqual(t, ok, true)
		assertEqual(t, r, "Hell")
	})
}

func Test_RangeConversion(t *testing.T) {
	t.Parallel()

	ns := FromString("    __Hello__   ")
	ns.Filter(func(r rune) bool {
		return !unicode.In(r, unicode.White_Space)
	})
	ns.ToLower()

	normRange, nrOk := ns.CoerceRangeToNormalized(NewOriginalRange(6, 11))
	assertEqual(t, nrOk, true)
	assertEqual(t, normRange, NewNormalizedRange(2, 7))
	{
		s, ok := ns.GetRange(normRange)
		assertEqual(t, ok, true)
		assertEqual(t, s, "hello")
	}
	{
		s, ok := ns.GetOriginalRange(normRange)
		assertEqual(t, ok, true)
		assertEqual(t, s, "Hello")
	}
	{
		s, ok := ns.GetRange(NewOriginalRange(6, 11))
		assertEqual(t, ok, true)
		assertEqual(t, s, "hello")
	}
	{
		s, ok := ns.GetOriginalRange(NewOriginalRange(6, 11))
		assertEqual(t, ok, true)
		assertEqual(t, s, "Hello")
	}

	// Make sure we get empty results, or false, only in specific cases
	{
		rng, ok := ns.CoerceRangeToNormalized(NewOriginalRange(0, 0))
		assertEqual(t, ok, true)
		assertEqual(t, rng, NewNormalizedRange(0, 0))
	}
	{
		rng, ok := ns.CoerceRangeToNormalized(NewOriginalRange(3, 3))
		assertEqual(t, ok, true)
		assertEqual(t, rng, NewNormalizedRange(3, 3))
	}
	{
		rng, ok := ns.CoerceRangeToNormalized(NewOriginalRange(15, ns.OriginalLen()))
		assertEqual(t, ok, true)
		assertEqual(t, rng, NewNormalizedRange(9, 9))
	}
	{
		rng, ok := ns.CoerceRangeToNormalized(NewOriginalRange(16, ns.OriginalLen()))
		assertEqual(t, ok, true)
		assertEqual(t, rng, NewNormalizedRange(16, 16))
	}
	{
		rng, ok := ns.CoerceRangeToNormalized(NewOriginalRange(17, ns.OriginalLen()))
		assertEqual(t, ok, false)
		assertEqual(t, rng, NewNormalizedRange(0, 0))
	}
	// ---
	{
		rng, ok := ns.CoerceRangeToOriginal(NewNormalizedRange(0, 0))
		assertEqual(t, ok, true)
		assertEqual(t, rng, NewOriginalRange(0, 0))
	}
	{
		rng, ok := ns.CoerceRangeToOriginal(NewNormalizedRange(3, 3))
		assertEqual(t, ok, true)
		assertEqual(t, rng, NewOriginalRange(3, 3))
	}
	{
		rng, ok := ns.CoerceRangeToOriginal(NewNormalizedRange(9, ns.Len()))
		assertEqual(t, ok, true)
		assertEqual(t, rng, NewOriginalRange(9, 9))
	}
	{
		rng, ok := ns.CoerceRangeToOriginal(NewNormalizedRange(10, ns.Len()))
		assertEqual(t, ok, false)
		assertEqual(t, rng, NewOriginalRange(0, 0))
	}
}

func Test_OriginalRange(t *testing.T) {
	t.Parallel()

	ns := FromString("Hello_______ World!")
	ns.Filter(func(r rune) bool {
		return r != '_'
	})
	ns.ToLower()
	{
		s, ok := ns.GetRange(NewNormalizedRange(6, 11))
		assertEqual(t, ok, true)
		assertEqual(t, s, "world")
	}
	{
		s, ok := ns.GetOriginalRange(NewNormalizedRange(6, 11))
		assertEqual(t, ok, true)
		assertEqual(t, s, "World")
	}

	originalRange, originalRangeOk := ns.CoerceRangeToOriginal(NewNormalizedRange(6, 11))
	assertEqual(t, originalRangeOk, true)
	assertEqual(t, originalRange, NewOriginalRange(13, 18))

	{
		s, ok := ns.GetRange(originalRange)
		assertEqual(t, ok, true)
		assertEqual(t, s, "world")
	}
	{
		s, ok := ns.GetOriginalRange(originalRange)
		assertEqual(t, ok, true)
		assertEqual(t, s, "World")
	}
}

func TestNormalizedStringTransform(t *testing.T) {
	t.Parallel()

	t.Run("Added characters alignments", func(t *testing.T) {
		t.Parallel()

		ns := FromString("野口 No")

		transformations := make([]RuneChange, 0)
		for _, r := range ns.Get() {
			if r > 0x4E00 {
				transformations = append(transformations, []RuneChange{
					{Rune: ' ', Change: 0},
					{Rune: r, Change: 1},
					{Rune: ' ', Change: 1},
				}...)
			} else {
				transformations = append(transformations, RuneChange{
					Rune:   r,
					Change: 0,
				})
			}
		}
		ns.Transform(transformations, 0)

		assertEqual(t, ns, New(
			"野口 No",
			" 野  口  No",
			[]AlignmentRange{
				{0, 3},
				{0, 3},
				{0, 3},
				{0, 3},
				{0, 3},
				{3, 6},
				{3, 6},
				{3, 6},
				{3, 6},
				{3, 6},
				{6, 7},
				{7, 8},
				{8, 9},
			},
			0,
		))
		assertEqual(t, ns.OriginalAlignments(), []AlignmentRange{
			{0, 5},
			{0, 5},
			{0, 5},
			{5, 10},
			{5, 10},
			{5, 10},
			{10, 11},
			{11, 12},
			{12, 13},
		})
	})

	t.Run("Added around edges", func(t *testing.T) {
		t.Parallel()

		ns := FromString("Hello")
		ns.Transform([]RuneChange{
			{' ', 1},
			{'H', 0},
			{'e', 0},
			{'l', 0},
			{'l', 0},
			{'o', 0},
			{' ', 1},
		}, 0)
		assertEqual(t, ns.Get(), " Hello ")

		r, ok := ns.GetOriginalRange(NewNormalizedRange(1, ns.Len()-1))
		assertEqual(t, ok, true)
		assertEqual(t, r, "Hello")
	})
}

func TestNormalizedStringSplit(t *testing.T) {
	t.Parallel()

	ns := FromString("The-final--countdown")

	test := func(behaviour SplitDelimiterBehavior, expected []string) {
		t.Helper()
		nss, err := ns.Split(splitpattern.FromRune('-'), behaviour)
		if err != nil {
			t.Error(err)
			return
		}
		if nss == nil {
			t.Error("the result is nil")
			return
		}
		actual := make([]string, len(nss))
		for i, ns := range nss {
			actual[i] = ns.Get()
		}
		assertEqual(t, actual, expected)
	}

	test(SplitDelimiterRemoved, []string{"The", "final", "countdown"})
	test(SplitDelimiterIsolated, []string{"The", "-", "final", "-", "-", "countdown"})
	test(SplitDelimiterMergedWithPrevious, []string{"The-", "final-", "-", "countdown"})
	test(SplitDelimiterMergedWithNext, []string{"The", "-final", "-", "-countdown"})
	test(SplitDelimiterContiguous, []string{"The", "-", "final", "--", "countdown"})
}

func TestNormalizedStringReplace(t *testing.T) {
	t.Parallel()

	t.Run("Simple 1", func(t *testing.T) {
		t.Parallel()

		ns := FromString(" Hello   friend ")
		err := ns.Replace(splitpattern.FromRune(' '), "_")
		if err != nil {
			t.Fatal(err)
		}
		assertEqual(t, ns.Get(), "_Hello___friend_")
	})

	t.Run("Simple 2", func(t *testing.T) {
		t.Parallel()

		ns := FromString("aaaab")
		err := ns.Replace(splitpattern.FromRune('a'), "b")
		if err != nil {
			t.Fatal(err)
		}
		assertEqual(t, ns.Get(), "bbbbb")
	})

	t.Run("Overlapping", func(t *testing.T) {
		t.Parallel()

		ns := FromString("aaaab")
		err := ns.Replace(splitpattern.FromString("aaa"), "b")
		if err != nil {
			t.Fatal(err)
		}
		assertEqual(t, ns.Get(), "bab")
	})

	t.Run("Regex", func(t *testing.T) {
		t.Parallel()

		ns := FromString(" Hello   friend ")
		r := regexp.MustCompile(`\s+`)

		err := ns.Replace(splitpattern.FromRegexp(r), "_")
		if err != nil {
			t.Fatal(err)
		}
		assertEqual(t, ns.Get(), "_Hello_friend_")
	})
}

func TestNormalizedStringTransformRange(t *testing.T) {
	t.Parallel()

	t.Run("Single bytes", func(t *testing.T) {
		t.Parallel()

		t.Run("Removing at the beginning", func(t *testing.T) {
			t.Parallel()

			ns := FromString("Hello friend")
			ns.TransformRange(
				NewOriginalRange(0, 4),
				[]RuneChange{{'Y', 0}},
				3,
			)
			assertEqual(t, ns, New(
				"Hello friend",
				"Yo friend",
				[]AlignmentRange{
					{3, 4},
					{4, 5},
					{5, 6},
					{6, 7},
					{7, 8},
					{8, 9},
					{9, 10},
					{10, 11},
					{11, 12},
				},
				0,
			))
			assertEqual(t, ns.OriginalAlignments(), []AlignmentRange{
				{0, 0},
				{0, 0},
				{0, 0},
				{0, 1},
				{1, 2},
				{2, 3},
				{3, 4},
				{4, 5},
				{5, 6},
				{6, 7},
				{7, 8},
				{8, 9},
			})
		})

		t.Run("Removing in the middle", func(t *testing.T) {
			t.Parallel()

			ns := FromString("Hello friend")
			ns.TransformRange(
				NewOriginalRange(3, 10),
				[]RuneChange{{'_', 0}, {'F', 0}, {'R', -2}},
				2,
			)
			assertEqual(t, ns, New(
				"Hello friend",
				"Hel_FRnd",
				[]AlignmentRange{
					{0, 1},
					{1, 2},
					{2, 3},
					{5, 6},
					{6, 7},
					{7, 8},
					{10, 11},
					{11, 12},
				},
				0,
			))
			assertEqual(t, ns.OriginalAlignments(), []AlignmentRange{
				{0, 1},
				{1, 2},
				{2, 3},
				{3, 3},
				{3, 3},
				{3, 4},
				{4, 5},
				{5, 6},
				{6, 6},
				{6, 6},
				{6, 7},
				{7, 8},
			})
		})

		t.Run("Removing at the end", func(t *testing.T) {
			t.Parallel()

			ns := FromString("Hello friend")
			ns.TransformRange(
				NewOriginalRange(5, ns.OriginalLen()),
				[]RuneChange{{'_', 0}, {'F', -5}},
				0,
			)
			assertEqual(t, ns, New(
				"Hello friend",
				"Hello_F",
				[]AlignmentRange{
					{0, 1},
					{1, 2},
					{2, 3},
					{3, 4},
					{4, 5},
					{5, 6},
					{6, 7},
				},
				0,
			))
			assertEqual(t, ns.OriginalAlignments(), []AlignmentRange{
				{0, 1},
				{1, 2},
				{2, 3},
				{3, 4},
				{4, 5},
				{5, 6},
				{6, 7},
				{7, 7},
				{7, 7},
				{7, 7},
				{7, 7},
				{7, 7},
			})
		})

		t.Run("Adding at the beginning", func(t *testing.T) {
			t.Parallel()

			ns := FromString("Hello friend")
			ns.TransformRange(
				NewOriginalRange(0, 1),
				[]RuneChange{{'H', 1}, {'H', 0}},
				0,
			)
			assertEqual(t, ns, New(
				"Hello friend",
				"HHello friend",
				[]AlignmentRange{
					{0, 0},
					{0, 1},
					{1, 2},
					{2, 3},
					{3, 4},
					{4, 5},
					{5, 6},
					{6, 7},
					{7, 8},
					{8, 9},
					{9, 10},
					{10, 11},
					{11, 12},
				},
				0,
			))
			assertEqual(t, ns.OriginalAlignments(), []AlignmentRange{
				{1, 2},
				{2, 3},
				{3, 4},
				{4, 5},
				{5, 6},
				{6, 7},
				{7, 8},
				{8, 9},
				{9, 10},
				{10, 11},
				{11, 12},
				{12, 13},
			})
		})

		t.Run("Adding at the beginning - alternative", func(t *testing.T) {
			t.Parallel()

			ns := FromString("Hello friend")
			ns.TransformRange(
				NewOriginalRange(0, 0),
				[]RuneChange{{'H', 1}},
				0,
			)
			assertEqual(t, ns, New(
				"Hello friend",
				"HHello friend",
				[]AlignmentRange{
					{0, 0},
					{0, 1},
					{1, 2},
					{2, 3},
					{3, 4},
					{4, 5},
					{5, 6},
					{6, 7},
					{7, 8},
					{8, 9},
					{9, 10},
					{10, 11},
					{11, 12},
				},
				0,
			))
			assertEqual(t, ns.OriginalAlignments(), []AlignmentRange{
				{1, 2},
				{2, 3},
				{3, 4},
				{4, 5},
				{5, 6},
				{6, 7},
				{7, 8},
				{8, 9},
				{9, 10},
				{10, 11},
				{11, 12},
				{12, 13},
			})
		})

		t.Run("Adding as part of the first character", func(t *testing.T) {
			t.Parallel()

			ns := FromString("Hello friend")
			ns.TransformRange(
				NewOriginalRange(0, 1),
				[]RuneChange{{'H', 0}, {'H', 1}},
				0,
			)
			assertEqual(t, ns, New(
				"Hello friend",
				"HHello friend",
				[]AlignmentRange{
					{0, 1},
					{0, 1},
					{1, 2},
					{2, 3},
					{3, 4},
					{4, 5},
					{5, 6},
					{6, 7},
					{7, 8},
					{8, 9},
					{9, 10},
					{10, 11},
					{11, 12},
				},
				0,
			))
			assertEqual(t, ns.OriginalAlignments(), []AlignmentRange{
				{0, 2},
				{2, 3},
				{3, 4},
				{4, 5},
				{5, 6},
				{6, 7},
				{7, 8},
				{8, 9},
				{9, 10},
				{10, 11},
				{11, 12},
				{12, 13},
			})
		})

		t.Run("Adding in the middle", func(t *testing.T) {
			t.Parallel()

			ns := FromString("Hello friend")
			ns.TransformRange(
				NewOriginalRange(5, 6),
				[]RuneChange{{'_', 0}, {'m', 1}, {'y', 1}, {'_', 1}},
				0,
			)
			assertEqual(t, ns, New(
				"Hello friend",
				"Hello_my_friend",
				[]AlignmentRange{
					{0, 1},
					{1, 2},
					{2, 3},
					{3, 4},
					{4, 5},
					{5, 6},
					{5, 6},
					{5, 6},
					{5, 6},
					{6, 7},
					{7, 8},
					{8, 9},
					{9, 10},
					{10, 11},
					{11, 12},
				},
				0,
			))
			assertEqual(t, ns.OriginalAlignments(), []AlignmentRange{
				{0, 1},
				{1, 2},
				{2, 3},
				{3, 4},
				{4, 5},
				{5, 9},
				{9, 10},
				{10, 11},
				{11, 12},
				{12, 13},
				{13, 14},
				{14, 15},
			})
		})

		t.Run("Adding at the end", func(t *testing.T) {
			t.Parallel()

			ns := FromString("Hello friend")
			ns.TransformRange(
				NewOriginalRange(11, ns.OriginalLen()),
				[]RuneChange{{'d', 0}, {'_', 1}, {'!', 1}},
				0,
			)
			assertEqual(t, ns, New(
				"Hello friend",
				"Hello friend_!",
				[]AlignmentRange{
					{0, 1},
					{1, 2},
					{2, 3},
					{3, 4},
					{4, 5},
					{5, 6},
					{6, 7},
					{7, 8},
					{8, 9},
					{9, 10},
					{10, 11},
					{11, 12},
					{11, 12},
					{11, 12},
				},
				0,
			))
			assertEqual(t, ns.OriginalAlignments(), []AlignmentRange{
				{0, 1},
				{1, 2},
				{2, 3},
				{3, 4},
				{4, 5},
				{5, 6},
				{6, 7},
				{7, 8},
				{8, 9},
				{9, 10},
				{10, 11},
				{11, 14},
			})
		})
	})

	t.Run("Multiple bytes", func(t *testing.T) {
		t.Parallel()

		t.Run("Removing at the beginning", func(t *testing.T) {
			t.Parallel()

			ns := FromString("𝔾𝕠𝕠𝕕")
			ns.TransformRange(
				NewOriginalRange(0, 9),
				[]RuneChange{{'G', -1}},
				0,
			)
			assertEqual(t, ns, New(
				"𝔾𝕠𝕠𝕕",
				"G𝕠𝕕",
				[]AlignmentRange{
					{0, 4},
					{8, 12},
					{8, 12},
					{8, 12},
					{8, 12},
					{12, 16},
					{12, 16},
					{12, 16},
					{12, 16},
				},
				0,
			))
			assertEqual(t, ns.OriginalAlignments(), []AlignmentRange{
				{0, 1},
				{0, 1},
				{0, 1},
				{0, 1},
				{1, 1},
				{1, 1},
				{1, 1},
				{1, 1},
				{1, 5},
				{1, 5},
				{1, 5},
				{1, 5},
				{5, 9},
				{5, 9},
				{5, 9},
				{5, 9},
			})

			r, ok := ns.GetRange(NewOriginalRange(0, 8))
			assertEqual(t, ok, true)
			assertEqual(t, r, "G")

			r, ok = ns.GetRange(NewOriginalRange(0, 4))
			assertEqual(t, ok, true)
			assertEqual(t, r, "G")

			r, ok = ns.GetOriginalRange(NewOriginalRange(0, 4))
			assertEqual(t, ok, true)
			assertEqual(t, r, "𝔾")

			r, ok = ns.GetOriginalRange(NewOriginalRange(0, 8))
			assertEqual(t, ok, true)
			assertEqual(t, r, "𝔾𝕠")
		})

		t.Run("Removing in the middle", func(t *testing.T) {
			t.Parallel()

			ns := FromString("𝔾𝕠𝕠𝕕")
			ns.TransformRange(
				NewOriginalRange(4, 12),
				[]RuneChange{{'o', -1}},
				0,
			)
			assertEqual(t, ns, New(
				"𝔾𝕠𝕠𝕕",
				"𝔾o𝕕",
				[]AlignmentRange{
					{0, 4},
					{0, 4},
					{0, 4},
					{0, 4},
					{4, 8},
					{12, 16},
					{12, 16},
					{12, 16},
					{12, 16},
				},
				0,
			))
			assertEqual(t, ns.OriginalAlignments(), []AlignmentRange{
				{0, 4},
				{0, 4},
				{0, 4},
				{0, 4},
				{4, 5},
				{4, 5},
				{4, 5},
				{4, 5},
				{5, 5},
				{5, 5},
				{5, 5},
				{5, 5},
				{5, 9},
				{5, 9},
				{5, 9},
				{5, 9},
			})
		})

		t.Run("Removing at the end", func(t *testing.T) {
			t.Parallel()

			ns := FromString("𝔾𝕠𝕠𝕕")
			ns.TransformRange(
				NewOriginalRange(12, ns.OriginalLen()),
				[]RuneChange{{'d', 0}, {'!', 1}},
				0,
			)
			assertEqual(t, ns, New(
				"𝔾𝕠𝕠𝕕",
				"𝔾𝕠𝕠d!",
				[]AlignmentRange{
					{0, 4},
					{0, 4},
					{0, 4},
					{0, 4},
					{4, 8},
					{4, 8},
					{4, 8},
					{4, 8},
					{8, 12},
					{8, 12},
					{8, 12},
					{8, 12},
					{12, 16},
					{12, 16},
				},
				0,
			))
		})

		t.Run("Adding at the beginning", func(t *testing.T) {
			t.Parallel()

			ns := FromString("𝔾𝕠𝕠𝕕")
			ns.TransformRange(
				NewOriginalRange(0, 4),
				[]RuneChange{{'_', 1}, {'𝔾', 0}},
				0,
			)
			assertEqual(t, ns, New(
				"𝔾𝕠𝕠𝕕",
				"_𝔾𝕠𝕠𝕕",
				[]AlignmentRange{
					{0, 0},
					{0, 4},
					{0, 4},
					{0, 4},
					{0, 4},
					{4, 8},
					{4, 8},
					{4, 8},
					{4, 8},
					{8, 12},
					{8, 12},
					{8, 12},
					{8, 12},
					{12, 16},
					{12, 16},
					{12, 16},
					{12, 16},
				},
				0,
			))
			assertEqual(t, ns.OriginalAlignments(), []AlignmentRange{
				{1, 5},
				{1, 5},
				{1, 5},
				{1, 5},
				{5, 9},
				{5, 9},
				{5, 9},
				{5, 9},
				{9, 13},
				{9, 13},
				{9, 13},
				{9, 13},
				{13, 17},
				{13, 17},
				{13, 17},
				{13, 17},
			})

			r, ok := ns.GetRange(NewOriginalRange(0, 8))
			assertEqual(t, ok, true)
			assertEqual(t, r, "𝔾𝕠")

			r, ok = ns.GetRange(NewOriginalRange(0, 4))
			assertEqual(t, ok, true)
			assertEqual(t, r, "𝔾")

			r, ok = ns.GetOriginalRange(NewOriginalRange(0, 4))
			assertEqual(t, ok, true)
			assertEqual(t, r, "𝔾")

			r, ok = ns.GetOriginalRange(NewOriginalRange(0, 8))
			assertEqual(t, ok, true)
			assertEqual(t, r, "𝔾𝕠")
		})

		t.Run("Adding at the beginning - alternative", func(t *testing.T) {
			t.Parallel()

			ns := FromString("𝔾𝕠𝕠𝕕")
			ns.TransformRange(
				NewOriginalRange(0, 0),
				[]RuneChange{{'_', 1}},
				0,
			)
			assertEqual(t, ns, New(
				"𝔾𝕠𝕠𝕕",
				"_𝔾𝕠𝕠𝕕",
				[]AlignmentRange{
					{0, 0},
					{0, 4},
					{0, 4},
					{0, 4},
					{0, 4},
					{4, 8},
					{4, 8},
					{4, 8},
					{4, 8},
					{8, 12},
					{8, 12},
					{8, 12},
					{8, 12},
					{12, 16},
					{12, 16},
					{12, 16},
					{12, 16},
				},
				0,
			))
			assertEqual(t, ns.OriginalAlignments(), []AlignmentRange{
				{1, 5},
				{1, 5},
				{1, 5},
				{1, 5},
				{5, 9},
				{5, 9},
				{5, 9},
				{5, 9},
				{9, 13},
				{9, 13},
				{9, 13},
				{9, 13},
				{13, 17},
				{13, 17},
				{13, 17},
				{13, 17},
			})

			r, ok := ns.GetRange(NewOriginalRange(0, 8))
			assertEqual(t, ok, true)
			assertEqual(t, r, "𝔾𝕠")

			r, ok = ns.GetRange(NewOriginalRange(0, 4))
			assertEqual(t, ok, true)
			assertEqual(t, r, "𝔾")

			r, ok = ns.GetOriginalRange(NewOriginalRange(0, 4))
			assertEqual(t, ok, true)
			assertEqual(t, r, "𝔾")

			r, ok = ns.GetOriginalRange(NewOriginalRange(0, 8))
			assertEqual(t, ok, true)
			assertEqual(t, r, "𝔾𝕠")
		})

		t.Run("Adding as part of the first character", func(t *testing.T) {
			t.Parallel()

			ns := FromString("𝔾𝕠𝕠𝕕")
			ns.TransformRange(
				NewOriginalRange(0, 4),
				[]RuneChange{{'𝔾', 0}, {'o', 1}},
				0,
			)
			assertEqual(t, ns, New(
				"𝔾𝕠𝕠𝕕",
				"𝔾o𝕠𝕠𝕕",
				[]AlignmentRange{
					{0, 4},
					{0, 4},
					{0, 4},
					{0, 4},
					{0, 4},
					{4, 8},
					{4, 8},
					{4, 8},
					{4, 8},
					{8, 12},
					{8, 12},
					{8, 12},
					{8, 12},
					{12, 16},
					{12, 16},
					{12, 16},
					{12, 16},
				},
				0,
			))

			assertEqual(t, ns.OriginalAlignments(), []AlignmentRange{
				{0, 5},
				{0, 5},
				{0, 5},
				{0, 5},
				{5, 9},
				{5, 9},
				{5, 9},
				{5, 9},
				{9, 13},
				{9, 13},
				{9, 13},
				{9, 13},
				{13, 17},
				{13, 17},
				{13, 17},
				{13, 17},
			})

			r, ok := ns.GetRange(NewOriginalRange(0, 8))
			assertEqual(t, ok, true)
			assertEqual(t, r, "𝔾o𝕠")

			r, ok = ns.GetRange(NewOriginalRange(0, 4))
			assertEqual(t, ok, true)
			assertEqual(t, r, "𝔾o")

			r, ok = ns.GetOriginalRange(NewOriginalRange(0, 4))
			assertEqual(t, ok, true)
			assertEqual(t, r, "𝔾")

			r, ok = ns.GetOriginalRange(NewOriginalRange(0, 8))
			assertEqual(t, ok, true)
			assertEqual(t, r, "𝔾𝕠")
		})

		t.Run("Adding in the middle", func(t *testing.T) {
			t.Parallel()

			ns := FromString("𝔾𝕠𝕠𝕕")
			ns.TransformRange(
				NewOriginalRange(4, 8),
				[]RuneChange{{'𝕠', 0}, {'o', 1}, {'o', 1}, {'o', 1}},
				0,
			)
			assertEqual(t, ns, New(
				"𝔾𝕠𝕠𝕕",
				"𝔾𝕠ooo𝕠𝕕",
				[]AlignmentRange{
					{0, 4},
					{0, 4},
					{0, 4},
					{0, 4},
					{4, 8},
					{4, 8},
					{4, 8},
					{4, 8},
					{4, 8},
					{4, 8},
					{4, 8},
					{8, 12},
					{8, 12},
					{8, 12},
					{8, 12},
					{12, 16},
					{12, 16},
					{12, 16},
					{12, 16},
				},
				0,
			))
			assertEqual(t, ns.OriginalAlignments(), []AlignmentRange{
				{0, 4},
				{0, 4},
				{0, 4},
				{0, 4},
				{4, 11},
				{4, 11},
				{4, 11},
				{4, 11},
				{11, 15},
				{11, 15},
				{11, 15},
				{11, 15},
				{15, 19},
				{15, 19},
				{15, 19},
				{15, 19},
			})
		})

		t.Run("Adding at the end", func(t *testing.T) {
			t.Parallel()

			ns := FromString("𝔾𝕠𝕠𝕕")
			ns.TransformRange(
				NewOriginalRange(16, ns.OriginalLen()),
				[]RuneChange{{'!', 1}},
				0,
			)
			assertEqual(t, ns, New(
				"𝔾𝕠𝕠𝕕",
				"𝔾𝕠𝕠𝕕!",
				[]AlignmentRange{
					{0, 4},
					{0, 4},
					{0, 4},
					{0, 4},
					{4, 8},
					{4, 8},
					{4, 8},
					{4, 8},
					{8, 12},
					{8, 12},
					{8, 12},
					{8, 12},
					{12, 16},
					{12, 16},
					{12, 16},
					{12, 16},
					{12, 16},
				},
				0,
			))
			assertEqual(t, ns.OriginalAlignments(), []AlignmentRange{
				{0, 4},
				{0, 4},
				{0, 4},
				{0, 4},
				{4, 8},
				{4, 8},
				{4, 8},
				{4, 8},
				{8, 12},
				{8, 12},
				{8, 12},
				{8, 12},
				{12, 17},
				{12, 17},
				{12, 17},
				{12, 17},
			})
		})
	})
}

func assertEqual(t *testing.T, actual, expected interface{}) {
	t.Helper()
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("expected\n  %#v\nactual\n  %#v", expected, actual)
	}
}
