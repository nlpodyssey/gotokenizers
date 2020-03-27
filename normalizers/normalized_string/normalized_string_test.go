// Copyright (c) 2020, The GoTokenizers Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package normalized_string

import (
	"fmt"
	"testing"

	. "github.com/nlpodyssey/gotokenizers/testing"
)

func TestNewNormalizedString(t *testing.T) {
	t.Run("with an empty string", func(t *testing.T) {
		assertNormalizedStringEqual(t, NewNormalizedString(""),
			NormalizedString{
				original:   "",
				normalized: "",
				alignments: []AlignmentRange{},
			},
		)
	})

	t.Run("with a simple string", func(t *testing.T) {
		assertNormalizedStringEqual(t, NewNormalizedString("Abc"),
			NormalizedString{
				original:   "Abc",
				normalized: "Abc",
				alignments: []AlignmentRange{{0, 1}, {1, 2}, {2, 3}},
			},
		)
	})

	t.Run("with a string containing non-ASCII characters", func(t *testing.T) {
		assertNormalizedStringEqual(t, NewNormalizedString("Süß"),
			NormalizedString{
				original:   "Süß",
				normalized: "Süß",
				alignments: []AlignmentRange{{0, 1}, {1, 2}, {2, 3}},
			},
		)
	})
}

func TestNormalizedStringEqual(t *testing.T) {
	t.Run("true if `normalized` is the same", func(t *testing.T) {
		a := NormalizedString{
			original:   "a",
			normalized: "ab",
			alignments: []AlignmentRange{{0, 1}, {1, 1}},
		}
		b := NormalizedString{
			original:   "b",
			normalized: "ab",
			alignments: []AlignmentRange{{0, 0}, {0, 1}},
		}
		if !a.Equal(b) {
			t.Fail()
		}
	})

	t.Run("false if `normalized` differ", func(t *testing.T) {
		a := NormalizedString{
			original:   "a",
			normalized: "ab",
			alignments: []AlignmentRange{{0, 1}, {1, 1}},
		}
		b := NormalizedString{
			original:   "a",
			normalized: "az",
			alignments: []AlignmentRange{{0, 1}, {1, 1}},
		}
		if a.Equal(b) {
			t.Fail()
		}
	})
}

func TestNormalizedStringLen(t *testing.T) {
	t.Run("with an empty normalized string", func(t *testing.T) {
		ns := NormalizedString{
			original:   "ab",
			normalized: "",
			alignments: []AlignmentRange{},
		}
		AssertIntEqual(t, "Len()", ns.Len(), 0)
	})

	t.Run("with a simple normalized string", func(t *testing.T) {
		ns := NormalizedString{
			original:   "a",
			normalized: "ab",
			alignments: []AlignmentRange{{0, 1}, {1, 1}},
		}
		AssertIntEqual(t, "Len()", ns.Len(), 2)
	})

	t.Run("with a normalized string containing non-ASCII characters",
		func(t *testing.T) {
			ns := NormalizedString{
				original:   "S",
				normalized: "Süß",
				alignments: []AlignmentRange{{0, 1}, {1, 1}, {1, 1}},
			}
			AssertIntEqual(t, "Len()", ns.Len(), 3)
		})
}

func TestNormalizedStringLenOriginal(t *testing.T) {
	t.Run("with an empty original string", func(t *testing.T) {
		ns := NormalizedString{
			original:   "",
			normalized: "a",
			alignments: []AlignmentRange{{0, 0}},
		}
		AssertIntEqual(t, "LenOriginal()", ns.LenOriginal(), 0)
	})

	t.Run("with a simple original string", func(t *testing.T) {
		ns := NormalizedString{
			original:   "abc",
			normalized: "",
			alignments: []AlignmentRange{},
		}
		AssertIntEqual(t, "LenOriginal()", ns.LenOriginal(), 3)
	})

	t.Run("with an original string containing non-ASCII characters",
		func(t *testing.T) {
			ns := NormalizedString{
				original:   "Süß",
				normalized: "",
				alignments: []AlignmentRange{},
			}
			AssertIntEqual(t, "LenOriginal()", ns.LenOriginal(), 3)
		})
}

func TestNormalizedStringIsEmpty(t *testing.T) {
	t.Run("true if `normalized` is empty", func(t *testing.T) {
		ns := NormalizedString{
			original:   "abc",
			normalized: "",
			alignments: []AlignmentRange{},
		}
		if !ns.IsEmpty() {
			t.Fail()
		}
	})

	t.Run("false if `normalized` is not empty", func(t *testing.T) {
		ns := NormalizedString{
			original:   "",
			normalized: "a",
			alignments: []AlignmentRange{{0, 0}},
		}
		if ns.IsEmpty() {
			t.Fail()
		}
	})
}

func TestNormalizedStringGet(t *testing.T) {
	ns := NormalizedString{
		original:   "a",
		normalized: "ab",
		alignments: []AlignmentRange{{0, 1}, {1, 1}},
	}
	AssertStringEqual(t, "Get()", ns.Get(), "ab")
}

func TestNormalizedStringGetOriginal(t *testing.T) {
	ns := NormalizedString{
		original:   "a",
		normalized: "ab",
		alignments: []AlignmentRange{{0, 1}, {1, 1}},
	}
	AssertStringEqual(t, "GetOriginal()", ns.GetOriginal(), "a")
}

func TestNormalizedStringConvertOffsetCommonCases(t *testing.T) {
	run := func(
		name string,
		ns NormalizedString,
		rangeStart, rangeEnd int,
		expectedStart, expectedEnd int,
		expectedFlag bool,
	) {
		t.Run(fmt.Sprintf("NSOriginalRange | %s", name), func(t *testing.T) {
			r := NewNSOriginalRange(rangeStart, rangeEnd)
			start, end, flag := ns.ConvertOffset(&r)
			AssertIntEqual(t, "start", start, expectedStart)
			AssertIntEqual(t, "end", end, expectedEnd)
			AssertBoolEqual(t, "flag", flag, expectedFlag)
		})

		t.Run(fmt.Sprintf("NSNormalizedRange | %s", name), func(t *testing.T) {
			r := NewNSNormalizedRange(rangeStart, rangeEnd)
			start, end, flag := ns.ConvertOffset(&r)
			AssertIntEqual(t, "start", start, expectedStart)
			AssertIntEqual(t, "end", end, expectedEnd)
			AssertBoolEqual(t, "flag", flag, expectedFlag)
		})
	}

	run("empty string, start < 0", NewNormalizedString(""), -1, 0, 0, 0, false)
	run("empty string, end < start", NewNormalizedString(""), 1, 0, 0, 0, false)
	run("empty string, end > 0", NewNormalizedString(""), 0, 1, 0, 0, false)
	run("start < 0", NewNormalizedString("Bar"), -1, 0, 0, 0, false)
	run("end < start", NewNormalizedString("Bar"), 1, 0, 0, 0, false)
	run("end > len", NewNormalizedString("Bar"), 0, 4, 0, 0, false)

	run("empty string, empty range", NewNormalizedString(""), 0, 0, 0, 0, true)
	run("one rune", NewNormalizedString("Bar"), 1, 2, 1, 2, true)
	run("more runes", NewNormalizedString("Bar Baz"), 2, 5, 2, 5, true)
	run("runes at beginning", NewNormalizedString("Bar"), 0, 2, 0, 2, true)
	run("runes at end", NewNormalizedString("Bar"), 1, 3, 1, 3, true)
	run("full string range", NewNormalizedString("Bar"), 0, 3, 0, 3, true)
}

func TestNormalizedStringConvertOffsetFromOriginalRange(t *testing.T) {
	run := func(
		name string,
		ns NormalizedString,
		rangeStart, rangeEnd int,
		expectedStart, expectedEnd int,
		expectedFlag bool,
	) {
		t.Run(name, func(t *testing.T) {
			r := NewNSOriginalRange(rangeStart, rangeEnd)
			start, end, flag := ns.ConvertOffset(&r)
			AssertIntEqual(t, "start", start, expectedStart)
			AssertIntEqual(t, "end", end, expectedEnd)
			AssertBoolEqual(t, "flag", flag, expectedFlag)
		})
	}

	run("chars removed at the beginning",
		NormalizedString{
			original:   "Bar",
			normalized: "ar",
			alignments: []AlignmentRange{{1, 2}, {2, 3}},
		},
		0, 2, 0, 1, true)

	run("chars removed at the end",
		NormalizedString{
			original:   "Bar",
			normalized: "Ba",
			alignments: []AlignmentRange{{0, 1}, {1, 2}},
		},
		2, 3, 1, 2, true)

	run("chars removed in the middle",
		NormalizedString{
			original:   "Bar Qux",
			normalized: "Baux",
			alignments: []AlignmentRange{{0, 1}, {1, 2}, {5, 6}, {6, 7}},
		},
		1, 6, 1, 3, true)

	run("chars added at the beginning",
		NormalizedString{
			original:   "Bar",
			normalized: "xyBar",
			alignments: []AlignmentRange{
				{0, 0}, {0, 0}, {0, 1}, {1, 2}, {2, 3},
			},
		},
		0, 2, 2, 4, true)

	run("chars added at the end",
		NormalizedString{
			original:   "Bar",
			normalized: "BarXy",
			alignments: []AlignmentRange{
				{0, 1}, {1, 2}, {2, 3}, {3, 3}, {3, 3},
			},
		},
		1, 3, 1, 5, true) // FIXME: it should be be 1, 3

	run("chars added in the middle",
		NormalizedString{
			original:   "abcd",
			normalized: "abXYcd",
			// FIXME: should be {0, 1}, {1, 2}, {2, 2}, {2, 2}, {2, 3}, {3, 4}
			alignments: []AlignmentRange{
				{0, 1}, {1, 2}, {1, 2}, {1, 2}, {2, 3}, {3, 4},
			},
		},
		1, 3, 3, 5, true) // FIXME: should be 2, 5

	run("range of removed chars only",
		NormalizedString{
			original:   "Bar Qux",
			normalized: "Baux",
			alignments: []AlignmentRange{{0, 1}, {1, 2}, {5, 6}, {6, 7}},
		},
		2, 5, 1, 2, true) // FIXME: should be (0, 0) or (2, 2) ...or even false?
}

func TestNormalizedStringConvertOffsetFromNormalizedlRange(t *testing.T) {
	run := func(
		name string,
		ns NormalizedString,
		rangeStart, rangeEnd int,
		expectedStart, expectedEnd int,
		expectedFlag bool,
	) {
		t.Run(name, func(t *testing.T) {
			r := NewNSNormalizedRange(rangeStart, rangeEnd)
			start, end, flag := ns.ConvertOffset(&r)
			AssertIntEqual(t, "start", start, expectedStart)
			AssertIntEqual(t, "end", end, expectedEnd)
			AssertBoolEqual(t, "flag", flag, expectedFlag)
		})
	}

	run("chars removed at the beginning",
		NormalizedString{
			original:   "Bar",
			normalized: "ar",
			alignments: []AlignmentRange{{1, 2}, {2, 3}},
		},
		0, 2, 1, 3, true)

	run("chars removed at the end",
		NormalizedString{
			original:   "Bar",
			normalized: "Ba",
			alignments: []AlignmentRange{{0, 1}, {1, 2}},
		},
		0, 2, 0, 2, true)

	run("chars removed in the middle",
		NormalizedString{
			original:   "Bar Qux",
			normalized: "Baux",
			alignments: []AlignmentRange{{0, 1}, {1, 2}, {5, 6}, {6, 7}},
		},
		1, 3, 1, 6, true)

	run("chars added at the beginning",
		NormalizedString{
			original:   "Bar",
			normalized: "xyBar",
			alignments: []AlignmentRange{
				{0, 0}, {0, 0}, {0, 1}, {1, 2}, {2, 3},
			},
		},
		0, 3, 0, 1, true)

	run("chars added at the end",
		NormalizedString{
			original:   "Bar",
			normalized: "BarXy",
			alignments: []AlignmentRange{
				{0, 1}, {1, 2}, {2, 3}, {3, 3}, {3, 3},
			},
		},
		2, 4, 2, 3, true)

	run("chars added in the middle",
		NormalizedString{
			original:   "abcd",
			normalized: "abXYcd",
			// FIXME: should be {0, 1}, {1, 2}, {2, 2}, {2, 2}, {2, 3}, {3, 4}
			alignments: []AlignmentRange{
				{0, 1}, {1, 2}, {1, 2}, {1, 2}, {2, 3}, {3, 4},
			},
		},
		1, 5, 1, 3, true)

	run("range of new chars only",
		NormalizedString{
			original:   "abcd",
			normalized: "abXYcd",
			// FIXME: should be {0, 1}, {1, 2}, {2, 2}, {2, 2}, {2, 3}, {3, 4}
			alignments: []AlignmentRange{
				{0, 1}, {1, 2}, {1, 2}, {1, 2}, {2, 3}, {3, 4},
			},
		},
		2, 4, 1, 2, true) // FIXME: should be (0, 0) or (2, 2) ...or even false?
}

func TestNormalizedStringGetRange(t *testing.T) {
	run := func(
		name string,
		ns NormalizedString,
		nsRange NSRange,
		expStr string,
		expFlag bool,
	) {
		t.Run(name, func(t *testing.T) {
			if s, f := ns.GetRange(nsRange); s != expStr || f != expFlag {
				t.Errorf("Expected (%#v, %v), but got (%#v, %v)",
					expStr, expFlag, s, f)
			}
		})
	}

	runOriginal := func(
		name string,
		ns NormalizedString,
		start, end int,
		expStr string,
		expFlag bool,
	) {
		nsRange := NewNSOriginalRange(start, end)
		run(fmt.Sprintf("NSOriginalRange | %s", name),
			ns, &nsRange, expStr, expFlag)
	}

	runNormalized := func(
		name string,
		ns NormalizedString,
		start, end int,
		expStr string,
		expFlag bool,
	) {
		nsRange := NewNSNormalizedRange(start, end)
		run(fmt.Sprintf("NSNormalizedRange | %s", name),
			ns, &nsRange, expStr, expFlag)
	}

	runBoth := func(
		name string,
		ns NormalizedString,
		start, end int,
		expStr string,
		expFlag bool,
	) {
		runOriginal(name, ns, start, end, expStr, expFlag)
		runNormalized(name, ns, start, end, expStr, expFlag)
	}

	ns := NewNormalizedString("")
	runBoth("empty string, empty range", ns, 0, 0, "", false)
	runBoth("empty string, start > end", ns, 1, 0, "", false)
	runBoth("empty string, start < 0", ns, -1, 0, "", false)
	runBoth("empty string, end > 0", ns, 0, 1, "", false)

	ns = NewNormalizedString("Bar")
	runBoth("no transformations, empty range", ns, 0, 0, "", false)
	runBoth("no transformations, start > end", ns, 1, 0, "", false)
	runBoth("no transformations, start < 0", ns, -1, 0, "", false)
	runBoth("no transformations, end > len", ns, 0, 4, "", false)
	runBoth("no transformations, leftmost range", ns, 0, 2, "Ba", true)
	runBoth("no transformations, rightmost", ns, 1, 3, "ar", true)
	runBoth("no transformations, middle range", ns, 1, 2, "a", true)
	runBoth("no transformations, full string range", ns, 0, 3, "Bar", true)

	runNormalized("can get newly inserted characters",
		NormalizedString{
			original:   "",
			normalized: "Bar",
			alignments: []AlignmentRange{{0, 0}, {0, 0}, {0, 0}},
		},
		1, 2, "a", true,
	)

	runOriginal("cannot get deleted characters",
		NormalizedString{
			original:   "Bar",
			normalized: "",
			alignments: []AlignmentRange{},
		},
		1, 2, "", false,
	)

	runOriginal("range including some deleted characters",
		NormalizedString{
			original:   "Bar Qux",
			normalized: "Baux",
			alignments: []AlignmentRange{{0, 1}, {1, 2}, {5, 6}, {6, 7}},
		},
		1, 6, "au", true,
	)

	runOriginal("range including some added characters",
		NormalizedString{
			original:   "abcd",
			normalized: "abXYcd",
			// FIXME: should be {0, 1}, {1, 2}, {2, 2}, {2, 2}, {2, 3}, {3, 4}
			alignments: []AlignmentRange{
				{0, 1}, {1, 2}, {1, 2}, {1, 2}, {2, 3}, {3, 4},
			},
		},
		1, 3, "Yc", true, // FIXME: wrong! It must be "abXYx"
	)
}

func TestNormalizedStringGetRangeOriginal(t *testing.T) {
	run := func(
		name string,
		ns NormalizedString,
		nsRange NSRange,
		expStr string,
		expFlag bool,
	) {
		t.Run(name, func(t *testing.T) {
			if s, f := ns.GetRangeOriginal(nsRange); s != expStr || f != expFlag {
				t.Errorf("Expected (%#v, %v), but got (%#v, %v)",
					expStr, expFlag, s, f)
			}
		})
	}

	runOriginal := func(
		name string,
		ns NormalizedString,
		start, end int,
		expStr string,
		expFlag bool,
	) {
		nsRange := NewNSOriginalRange(start, end)
		run(fmt.Sprintf("NSOriginalRange | %s", name),
			ns, &nsRange, expStr, expFlag)
	}

	runNormalized := func(
		name string,
		ns NormalizedString,
		start, end int,
		expStr string,
		expFlag bool,
	) {
		nsRange := NewNSNormalizedRange(start, end)
		run(fmt.Sprintf("NSNormalizedRange | %s", name),
			ns, &nsRange, expStr, expFlag)
	}

	runBoth := func(
		name string,
		ns NormalizedString,
		start, end int,
		expStr string,
		expFlag bool,
	) {
		runOriginal(name, ns, start, end, expStr, expFlag)
		runNormalized(name, ns, start, end, expStr, expFlag)
	}

	ns := NewNormalizedString("")
	runBoth("empty string, empty range", ns, 0, 0, "", false)
	runBoth("empty string, start > end", ns, 1, 0, "", false)
	runBoth("empty string, start < 0", ns, -1, 0, "", false)
	runBoth("empty string, end > 0", ns, 0, 1, "", false)

	ns = NewNormalizedString("Bar")
	runBoth("no transformations, empty range", ns, 0, 0, "", false)
	runBoth("no transformations, start > end", ns, 1, 0, "", false)
	runBoth("no transformations, start < 0", ns, -1, 0, "", false)
	runBoth("no transformations, end > len", ns, 0, 4, "", false)
	runBoth("no transformations, leftmost range", ns, 0, 2, "Ba", true)
	runBoth("no transformations, rightmost", ns, 1, 3, "ar", true)
	runBoth("no transformations, middle range", ns, 1, 2, "a", true)
	runBoth("no transformations, full string range", ns, 0, 3, "Bar", true)

	runOriginal("can get deleted characters",
		NormalizedString{
			original:   "Bar",
			normalized: "",
			alignments: []AlignmentRange{},
		},
		1, 2, "a", true,
	)

	runNormalized("cannot get newly inserted characters",
		NormalizedString{
			original:   "",
			normalized: "Bar",
			alignments: []AlignmentRange{{0, 0}, {0, 0}, {0, 0}},
		},
		1, 2, "", false,
	)

	runNormalized("range including some deleted characters",
		NormalizedString{
			original:   "Bar Qux",
			normalized: "Baux",
			alignments: []AlignmentRange{{0, 1}, {1, 2}, {5, 6}, {6, 7}},
		},
		1, 3, "ar Qu", true,
	)

	runNormalized("range including some added characters",
		NormalizedString{
			original:   "abcd",
			normalized: "abXYcd",
			// FIXME: should be {0, 1}, {1, 2}, {2, 2}, {2, 2}, {2, 3}, {3, 4}
			alignments: []AlignmentRange{
				{0, 1}, {1, 2}, {1, 2}, {1, 2}, {2, 3}, {3, 4},
			},
		},
		1, 5, "bc", true,
	)
}

func TestNormalizedStringTransform(t *testing.T) {
	run := func(
		name string,
		ns NormalizedString,
		dest []RuneChanges,
		initialOffset int,
		expected NormalizedString,
	) {
		t.Run(name, func(t *testing.T) {
			ns.Transform(dest, initialOffset)
			assertNormalizedStringEqual(t, ns, expected)
		})
	}

	run("empty string, empty changes", NewNormalizedString(""),
		[]RuneChanges{}, 0,
		NormalizedString{
			original:   "",
			normalized: "",
			alignments: []AlignmentRange{},
		})

	run("non-empty string, empty changes", NewNormalizedString("Bar"),
		[]RuneChanges{}, 0,
		NormalizedString{
			original:   "Bar",
			normalized: "",
			alignments: []AlignmentRange{},
		})

	run("non-empty string, empty changes, offset", NewNormalizedString("Bar"),
		[]RuneChanges{}, 3,
		NormalizedString{
			original:   "Bar",
			normalized: "",
			alignments: []AlignmentRange{},
		})

	run("1:1 mapping (all changes = 0)", NewNormalizedString("Bär"),
		[]RuneChanges{{'S', 0}, {'ü', 0}, {'ß', 0}}, 0,
		NormalizedString{
			original:   "Bär",
			normalized: "Süß",
			alignments: []AlignmentRange{{0, 1}, {1, 2}, {2, 3}},
		})

	run("1:1 mapping (all changes == 0), with offset",
		NewNormalizedString("Bär"),
		[]RuneChanges{{'ü', 0}, {'ß', 0}}, 1,
		NormalizedString{
			original:   "Bär",
			normalized: "üß",
			alignments: []AlignmentRange{{1, 2}, {2, 3}},
		})

	run("adding a rune and deleting the rest (only one Change = 1)",
		NewNormalizedString("Bar"),
		[]RuneChanges{{'x', 1}}, 0,
		NormalizedString{
			original:   "Bar",
			normalized: "x",
			alignments: []AlignmentRange{{0, 0}},
		})

	run("adding a rune at the beginning", NewNormalizedString("x"),
		[]RuneChanges{{'a', 1}, {'x', 0}}, 0,
		NormalizedString{
			original:   "x",
			normalized: "ax",
			alignments: []AlignmentRange{{0, 0}, {0, 1}},
		})

	run("adding more runes at the beginning", NewNormalizedString("x"),
		[]RuneChanges{{'a', 1}, {'b', 1}, {'x', 0}}, 0,
		NormalizedString{
			original:   "x",
			normalized: "abx",
			alignments: []AlignmentRange{{0, 0}, {0, 0}, {0, 1}},
		})

	run("adding a rune at the end", NewNormalizedString("x"),
		[]RuneChanges{{'x', 0}, {'a', 1}}, 0,
		NormalizedString{
			original:   "x",
			normalized: "xa",
			alignments: []AlignmentRange{
				// FIXME: should be {0, 1}, {1, 1}?
				{0, 1}, {0, 1},
			},
		})

	run("adding more runes at the end", NewNormalizedString("x"),
		[]RuneChanges{{'x', 0}, {'a', 1}, {'b', 1}}, 0,
		NormalizedString{
			original:   "x",
			normalized: "xab",
			alignments: []AlignmentRange{
				// FIXME: shouldn't be {0, 1}, {1, 1}, {1, 1}?
				{0, 1}, {0, 1}, {0, 1},
			},
		})

	run("adding runes at beginning and end", NewNormalizedString("x"),
		[]RuneChanges{{'a', 1}, {'x', 0}, {'b', 1}}, 0,
		NormalizedString{
			original:   "x",
			normalized: "axb",
			alignments: []AlignmentRange{
				// FIXME: shouldn't be {0, 0}, {0, 1}, {1, 1}?
				{0, 0}, {0, 1}, {0, 1},
			},
		})

	run("adding a rune in the middle", NewNormalizedString("ab"),
		[]RuneChanges{{'a', 0}, {'x', 1}, {'b', 0}}, 0,
		NormalizedString{
			original:   "ab",
			normalized: "axb",
			alignments: []AlignmentRange{
				// FIXME: shouldn't be {0, 1}, {1, 1}, {1, 2}?
				{0, 1}, {0, 1}, {1, 2},
			},
		})

	run("adding multiple runes in the middle", NewNormalizedString("ab"),
		[]RuneChanges{{'a', 0}, {'x', 1}, {'y', 1}, {'b', 0}}, 0,
		NormalizedString{
			original:   "ab",
			normalized: "axyb",
			alignments: []AlignmentRange{
				// FIXME: shouldn't be {0, 1}, {1, 1}, {1, 1}, {1, 2}?
				{0, 1}, {0, 1}, {0, 1}, {1, 2},
			},
		})

	run("change -1 at the beginning", NewNormalizedString("Bar"),
		[]RuneChanges{{'Q', -1}, {'r', 0}}, 0,
		NormalizedString{
			original:   "Bar",
			normalized: "Qr",
			alignments: []AlignmentRange{{0, 1}, {2, 3}},
		})

	run("change -1 at the end", NewNormalizedString("Bar"),
		[]RuneChanges{{'B', 0}, {'x', -1}}, 0,
		NormalizedString{
			original:   "Bar",
			normalized: "Bx",
			alignments: []AlignmentRange{{0, 1}, {1, 2}},
		})

	run("change -1 in the middle", NewNormalizedString("abcd"),
		[]RuneChanges{{'a', 0}, {'x', -1}, {'d', 0}}, 0,
		NormalizedString{
			original:   "abcd",
			normalized: "axd",
			alignments: []AlignmentRange{{0, 1}, {1, 2}, {3, 4}},
		})

	run("change -2 in the middle", NewNormalizedString("abcde"),
		[]RuneChanges{{'a', 0}, {'x', -2}, {'e', 0}}, 0,
		NormalizedString{
			original:   "abcde",
			normalized: "axe",
			alignments: []AlignmentRange{{0, 1}, {1, 2}, {4, 5}},
		})

	t.Run("it panics with Changes > 1", func(t *testing.T) {
		ns := NewNormalizedString("Bar")
		AssertPanic(t, "using a Change = 2", func() {
			ns.Transform([]RuneChanges{{'a', 2}}, 0)
		})
	})
}

func TestNormalizedStringFilter(t *testing.T) {
	run := func(
		name string,
		ns NormalizedString,
		filter func(rune) bool,
		expected NormalizedString,
	) {
		t.Run(name, func(t *testing.T) {
			ns.Filter(filter)
			assertNormalizedStringEqual(t, ns, expected)
		})
	}

	run("filter empty string", NewNormalizedString(""),
		func(r rune) bool { return true },
		NormalizedString{
			original:   "",
			normalized: "",
			alignments: []AlignmentRange{},
		})

	run("filter all characters true", NewNormalizedString("Bar"),
		func(r rune) bool { return true },
		NormalizedString{
			original:   "Bar",
			normalized: "Bar",
			alignments: []AlignmentRange{{0, 1}, {1, 2}, {2, 3}},
		})

	run("filter all characters false", NewNormalizedString("Bar"),
		func(r rune) bool { return false },
		NormalizedString{
			original:   "Bar",
			normalized: "",
			alignments: []AlignmentRange{},
		})

	run("remove one character at the beginning", NewNormalizedString("abcd"),
		func(r rune) bool { return r > 'a' },
		NormalizedString{
			original:   "abcd",
			normalized: "bcd",
			alignments: []AlignmentRange{{1, 2}, {2, 3}, {3, 4}},
		})

	run("remove more characters at the beginning", NewNormalizedString("abcde"),
		func(r rune) bool { return r > 'b' },
		NormalizedString{
			original:   "abcde",
			normalized: "cde",
			alignments: []AlignmentRange{{2, 3}, {3, 4}, {4, 5}},
		})

	run("remove one character at the end", NewNormalizedString("abcd"),
		func(r rune) bool { return r < 'd' },
		NormalizedString{
			original:   "abcd",
			normalized: "abc",
			alignments: []AlignmentRange{{0, 1}, {1, 2}, {2, 3}},
		})

	run("remove more characters at the end", NewNormalizedString("abcde"),
		func(r rune) bool { return r < 'd' },
		NormalizedString{
			original:   "abcde",
			normalized: "abc",
			alignments: []AlignmentRange{{0, 1}, {1, 2}, {2, 3}},
		})

	run("remove one character in the middle", NewNormalizedString("axb"),
		func(r rune) bool { return r < 'x' },
		NormalizedString{
			original:   "axb",
			normalized: "ab",
			alignments: []AlignmentRange{{0, 1}, {2, 3}},
		})

	run("remove more characters in the middle", NewNormalizedString("axyb"),
		func(r rune) bool { return r < 'x' },
		NormalizedString{
			original:   "axyb",
			normalized: "ab",
			alignments: []AlignmentRange{{0, 1}, {3, 4}},
		})

	run("remove characters in various places", NewNormalizedString("awxbycz"),
		func(r rune) bool { return r < 'w' },
		NormalizedString{
			original:   "awxbycz",
			normalized: "abc",
			alignments: []AlignmentRange{{0, 1}, {3, 4}, {5, 6}},
		})

	run("filter non-ASCII runes", NewNormalizedString("süß!"),
		func(r rune) bool { return r < 'z' },
		NormalizedString{
			original:   "süß!",
			normalized: "s!",
			alignments: []AlignmentRange{{0, 1}, {3, 4}},
		})
}

func TestNormalizedStringPrepend(t *testing.T) {
	run := func(
		name string,
		ns NormalizedString,
		s string,
		expected NormalizedString,
	) {
		t.Run(name, func(t *testing.T) {
			ns.Prepend(s)
			assertNormalizedStringEqual(t, ns, expected)
		})
	}

	run("prepend empty string to empty string", NewNormalizedString(""), "",
		NormalizedString{
			original:   "",
			normalized: "",
			alignments: []AlignmentRange{},
		})

	run("prepend empty string ", NewNormalizedString("ab"), "",
		NormalizedString{
			original:   "ab",
			normalized: "ab",
			alignments: []AlignmentRange{{0, 1}, {1, 2}},
		})

	run("prepend to empty string ", NewNormalizedString(""), "ab",
		NormalizedString{
			original:   "",
			normalized: "ab",
			alignments: []AlignmentRange{{0, 0}, {0, 0}},
		})

	run("prepend one rune", NewNormalizedString("ab"), "x",
		NormalizedString{
			original:   "ab",
			normalized: "xab",
			alignments: []AlignmentRange{{0, 0}, {0, 1}, {1, 2}},
		})

	run("prepend more runes", NewNormalizedString("ab"), "xy",
		NormalizedString{
			original:   "ab",
			normalized: "xyab",
			alignments: []AlignmentRange{
				{0, 0}, {0, 0}, {0, 1}, {1, 2},
			},
		})

	run("non-ASCII runes", NewNormalizedString("äö"), "süß",
		NormalizedString{
			original:   "äö",
			normalized: "süßäö",
			alignments: []AlignmentRange{
				{0, 0}, {0, 0}, {0, 0}, {0, 1}, {1, 2},
			},
		})
}

func TestNormalizedStringAppend(t *testing.T) {
	run := func(
		name string,
		ns NormalizedString,
		s string,
		expected NormalizedString,
	) {
		t.Run(name, func(t *testing.T) {
			ns.Append(s)
			assertNormalizedStringEqual(t, ns, expected)
		})
	}

	run("append empty string to empty string", NewNormalizedString(""), "",
		NormalizedString{
			original:   "",
			normalized: "",
			alignments: []AlignmentRange{},
		})

	run("append empty string ", NewNormalizedString("ab"), "",
		NormalizedString{
			original:   "ab",
			normalized: "ab",
			alignments: []AlignmentRange{{0, 1}, {1, 2}},
		})

	run("append to empty string ", NewNormalizedString(""), "ab",
		NormalizedString{
			original:   "",
			normalized: "ab",
			alignments: []AlignmentRange{{0, 0}, {0, 0}},
		})

	run("append one rune", NewNormalizedString("ab"), "x",
		NormalizedString{
			original:   "ab",
			normalized: "abx",
			alignments: []AlignmentRange{{0, 1}, {1, 2}, {2, 2}},
		})

	run("append more runes", NewNormalizedString("ab"), "xy",
		NormalizedString{
			original:   "ab",
			normalized: "abxy",
			alignments: []AlignmentRange{
				{0, 1}, {1, 2}, {2, 2}, {2, 2},
			},
		})

	run("non-ASCII runes", NewNormalizedString("äö"), "süß",
		NormalizedString{
			original:   "äö",
			normalized: "äösüß",
			alignments: []AlignmentRange{
				{0, 1}, {1, 2}, {2, 2}, {2, 2}, {2, 2},
			},
		})
}

func TestNormalizedStringMap(t *testing.T) {
	ns := NewNormalizedString("abc")

	ns.Map(func(r rune) rune { return r + 3 })

	assertNormalizedStringEqual(t, ns, NormalizedString{
		original:   "abc",
		normalized: "def",
		alignments: []AlignmentRange{{0, 1}, {1, 2}, {2, 3}},
	})
}

func TestNormalizedStringForEach(t *testing.T) {
	ns := NewNormalizedString("abc")

	visitedRunes := make([]rune, 0, 3)
	ns.ForEach(func(r rune) { visitedRunes = append(visitedRunes, r) })

	AssertRuneSliceEqual(t, "visited runes", visitedRunes, []rune{'a', 'b', 'c'})
}

func TestNormalizedStringLowercase(t *testing.T) {
	run := func(s string, expected NormalizedString) {
		t.Run(fmt.Sprintf("%#v", s), func(t *testing.T) {
			ns := NewNormalizedString(s)
			ns.Lowercase()
			assertNormalizedStringEqual(t, ns, expected)
		})
	}

	run("", NormalizedString{
		original:   "",
		normalized: "",
		alignments: []AlignmentRange{},
	})

	run("AË", NormalizedString{
		original:   "AË",
		normalized: "aë",
		alignments: []AlignmentRange{{0, 1}, {1, 2}},
	})
}

func TestNormalizedStringUppercase(t *testing.T) {
	run := func(s string, expected NormalizedString) {
		t.Run(fmt.Sprintf("%#v", s), func(t *testing.T) {
			ns := NewNormalizedString(s)
			ns.Uppercase()
			assertNormalizedStringEqual(t, ns, expected)
		})
	}

	run("", NormalizedString{
		original:   "",
		normalized: "",
		alignments: []AlignmentRange{},
	})

	run("aë", NormalizedString{
		original:   "aë",
		normalized: "AË",
		alignments: []AlignmentRange{{0, 1}, {1, 2}},
	})
}

func TestNormalizedStringSplitOff(t *testing.T) {
	run := func(
		name string,
		ns NormalizedString,
		at int,
		expModified, expNewNs NormalizedString,
	) {
		t.Run(name, func(t *testing.T) {
			actualNewNs := ns.SplitOff(at)
			assertNormalizedStringEqual(t, ns, expModified)
			assertNormalizedStringEqual(t, actualNewNs, expNewNs)
		})
	}

	run("empty string, at 0", NewNormalizedString(""), 0,
		NewNormalizedString(""), NewNormalizedString(""))
	run("empty string, at 1", NewNormalizedString(""), 1,
		NewNormalizedString(""), NewNormalizedString(""))

	run("no transformations, split at 0", NewNormalizedString("Bar"), 0,
		NewNormalizedString(""), NewNormalizedString("Bar"))
	run("no transformations, split at len", NewNormalizedString("Bar"), 3,
		NewNormalizedString("Bar"), NewNormalizedString(""))

	// FIXME: fix SplitOff creating new invalid mappings and add more tests

	{
		ns := NewNormalizedString("Foo")
		AssertPanic(t, "using a negative index", func() {
			ns.SplitOff(-1)
		})
	}
}

func TestNormalizedStringMergeWith(t *testing.T) {
	run := func(name string, ns, other, expected NormalizedString) {
		t.Run(name, func(t *testing.T) {
			ns.MergeWith(other)
			assertNormalizedStringEqual(t, ns, expected)
		})
	}

	run("empty strings", NewNormalizedString(""), NewNormalizedString(""),
		NewNormalizedString(""))

	// FIXME: fix MergeWith creating new invalid mappings and add more tests
}

func TestNormalizedStringAlignmentEqual(t *testing.T) {
	t.Run("true if `pos` and `changes` are the same", func(t *testing.T) {
		a := AlignmentRange{start: 1, end: 2}
		b := AlignmentRange{start: 1, end: 2}
		if !a.Equal(b) {
			t.Fail()
		}
	})

	t.Run("false if `pos` differ", func(t *testing.T) {
		a := AlignmentRange{start: 1, end: 2}
		b := AlignmentRange{start: 3, end: 2}
		if a.Equal(b) {
			t.Fail()
		}
	})

	t.Run("false if `changes` differ", func(t *testing.T) {
		a := AlignmentRange{start: 1, end: 2}
		b := AlignmentRange{start: 1, end: 3}
		if a.Equal(b) {
			t.Fail()
		}
	})

	t.Run("false if `pos` and `changes` are different", func(t *testing.T) {
		a := AlignmentRange{start: 1, end: 2}
		b := AlignmentRange{start: 3, end: 4}
		if a.Equal(b) {
			t.Fail()
		}
	})
}

func assertNormalizedStringEqual(
	t *testing.T,
	actual, expected NormalizedString,
) {
	AssertStringEqual(t, "original", actual.original, expected.original)
	AssertStringEqual(t, "normalized", actual.normalized, expected.normalized)
	assertAlignmentsEqual(t, actual.alignments, expected.alignments)
}

func assertAlignmentsEqual(
	t *testing.T,
	actual, expected []AlignmentRange,
) {
	if len(expected) != len(actual) {
		t.Errorf(
			"Expected alignments to be %v, but got %v (lengths differ)",
			expected, actual)
		return
	}
	for index := range expected {
		if !expected[index].Equal(actual[index]) {
			t.Errorf(
				"Expected alignments to be %v, but got %v"+
					" (mismatch at index %v: expected %v, actual %v)",
				expected, actual, index, expected[index], actual[index])
		}
	}
}
