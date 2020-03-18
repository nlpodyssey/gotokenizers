// Copyright (c) 2020, The GoTokenizers Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package normalized_string

import (
	"testing"

	. "github.com/saientist/gotokenizers/testing"
)

func TestNewNormalizedString(t *testing.T) {
	t.Run("with an empty string", func(t *testing.T) {
		assertNormalizedStringEqual(t, NewNormalizedString(""),
			NormalizedString{
				original:   "",
				normalized: "",
				alignments: []NormalizedStringAlignment{},
			},
		)
	})

	t.Run("with a simple string", func(t *testing.T) {
		assertNormalizedStringEqual(t, NewNormalizedString("Abc"),
			NormalizedString{
				original:   "Abc",
				normalized: "Abc",
				alignments: []NormalizedStringAlignment{
					{0, 1}, {1, 2}, {2, 3},
				},
			},
		)
	})

	t.Run("with a string containing non-ASCII characters", func(t *testing.T) {
		assertNormalizedStringEqual(t, NewNormalizedString("Süß"),
			NormalizedString{
				original:   "Süß",
				normalized: "Süß",
				alignments: []NormalizedStringAlignment{
					{0, 1}, {1, 2}, {2, 3},
				},
			},
		)
	})
}

func TestNormalizedStringEqual(t *testing.T) {
	t.Run("true if `normalized` is the same", func(t *testing.T) {
		a := NormalizedString{
			original:   "a",
			normalized: "ab",
			alignments: []NormalizedStringAlignment{{0, 1}, {1, 1}},
		}
		b := NormalizedString{
			original:   "b",
			normalized: "ab",
			alignments: []NormalizedStringAlignment{{0, 0}, {0, 1}},
		}
		if !a.Equal(b) {
			t.Fail()
		}
	})

	t.Run("false if `normalized` differ", func(t *testing.T) {
		a := NormalizedString{
			original:   "a",
			normalized: "ab",
			alignments: []NormalizedStringAlignment{{0, 1}, {1, 1}},
		}
		b := NormalizedString{
			original:   "a",
			normalized: "az",
			alignments: []NormalizedStringAlignment{{0, 1}, {1, 1}},
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
			alignments: []NormalizedStringAlignment{},
		}
		AssertIntEqual(t, "Len()", ns.Len(), 0)
	})

	t.Run("with a simple normalized string", func(t *testing.T) {
		ns := NormalizedString{
			original:   "a",
			normalized: "ab",
			alignments: []NormalizedStringAlignment{{0, 1}, {1, 1}},
		}
		AssertIntEqual(t, "Len()", ns.Len(), 2)
	})

	t.Run("with a normalized string containing non-ASCII characters",
		func(t *testing.T) {
			ns := NormalizedString{
				original:   "S",
				normalized: "Süß",
				alignments: []NormalizedStringAlignment{{0, 1}, {1, 1}, {1, 1}},
			}
			AssertIntEqual(t, "Len()", ns.Len(), 3)
		})
}

func TestNormalizedStringLenOriginal(t *testing.T) {
	t.Run("with an empty original string", func(t *testing.T) {
		ns := NormalizedString{
			original:   "",
			normalized: "a",
			alignments: []NormalizedStringAlignment{{0, 0}},
		}
		AssertIntEqual(t, "LenOriginal()", ns.LenOriginal(), 0)
	})

	t.Run("with a simple original string", func(t *testing.T) {
		ns := NormalizedString{
			original:   "abc",
			normalized: "",
			alignments: []NormalizedStringAlignment{},
		}
		AssertIntEqual(t, "LenOriginal()", ns.LenOriginal(), 3)
	})

	t.Run("with an original string containing non-ASCII characters",
		func(t *testing.T) {
			ns := NormalizedString{
				original:   "Süß",
				normalized: "",
				alignments: []NormalizedStringAlignment{},
			}
			AssertIntEqual(t, "LenOriginal()", ns.LenOriginal(), 3)
		})
}

func TestNormalizedStringIsEmpty(t *testing.T) {
	t.Run("true if `normalized` is empty", func(t *testing.T) {
		ns := NormalizedString{
			original:   "abc",
			normalized: "",
			alignments: []NormalizedStringAlignment{},
		}
		if !ns.IsEmpty() {
			t.Fail()
		}
	})

	t.Run("false if `normalized` is not empty", func(t *testing.T) {
		ns := NormalizedString{
			original:   "",
			normalized: "a",
			alignments: []NormalizedStringAlignment{{0, 0}},
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
		alignments: []NormalizedStringAlignment{{0, 1}, {1, 1}},
	}
	AssertStringEqual(t, "Get()", ns.Get(), "ab")
}

func TestNormalizedStringGetOriginal(t *testing.T) {
	ns := NormalizedString{
		original:   "a",
		normalized: "ab",
		alignments: []NormalizedStringAlignment{{0, 1}, {1, 1}},
	}
	AssertStringEqual(t, "GetOriginal()", ns.GetOriginal(), "a")
}

func TestNormalizedStringGetRange(t *testing.T) {
	run := func(
		name string,
		ns NormalizedString,
		start, end int,
		expStr string,
		expFlag bool,
	) {
		t.Run(name, func(t *testing.T) {
			if s, f := ns.GetRange(start, end); s != expStr || f != expFlag {
				t.Errorf("Expected (%#v, %v), but got (%#v, %v)",
					expStr, expFlag, s, f)
			}
		})
	}

	ns := NormalizedString{
		original:   "",
		normalized: "Foo süß bar",
		alignments: []NormalizedStringAlignment{
			{0, 0}, {0, 0}, {0, 0}, {0, 0}, {0, 0}, {0, 0},
			{0, 0}, {0, 0}, {0, 0}, {0, 0}, {0, 0},
		},
	}

	run("blank result with start < end", ns, 3, 2, "", false)
	run("blank result with start = end", ns, 3, 3, "", false)
	run("blank result with start < 0", ns, -1, 3, "", false)
	run("blank result with start > runes length", ns, 1, 12, "", false)
	run("valid result with a left-most range", ns, 0, 3, "Foo", true)
	run("valid result with a right-most range", ns, 8, 11, "bar", true)
	run("valid result with a range in the middle", ns, 4, 7, "süß", true)
	run("valid result with full string range", ns, 0, 11, "Foo süß bar", true)

	ns = NormalizedString{
		original:   "foo",
		normalized: "",
		alignments: []NormalizedStringAlignment{},
	}

	run("blank result with an empty string", ns, 0, 0, "", false)
}

func TestNormalizedStringGetRangeOriginal(t *testing.T) {
	run := func(
		name string,
		ns NormalizedString,
		start, end int,
		expStr string,
		expFlag bool,
	) {
		t.Run(name, func(t *testing.T) {
			s, f := ns.GetRangeOriginal(start, end)
			if s != expStr || f != expFlag {
				t.Errorf("Expected (%#v, %v), but got (%#v, %v)",
					expStr, expFlag, s, f)
			}
		})
	}

	ns := NormalizedString{
		original:   "süß",
		normalized: "Foo süß bar",
		alignments: []NormalizedStringAlignment{
			{0, 0}, {0, 0}, {0, 0}, {0, 0}, {0, 1}, {1, 2},
			{2, 3}, {3, 3}, {3, 3}, {3, 3}, {3, 3},
		},
	}

	run("blank result with start < end", ns, 3, 2, "", false)
	run("blank result with start < end", ns, 3, 2, "", false)
	run("blank result with start = end", ns, 3, 3, "", false)
	run("blank result with start < 0", ns, -1, 3, "", false)
	run("blank result with start > runes length", ns, 1, 12, "", false)
	run("blank result starting from a prepended range", ns, 0, 4, "", false)
	run("blank result starting from an appended range", ns, 7, 11, "", false)
	run("valid result with partially mapped left-most range",
		ns, 0, 5, "s", true)
	run("valid result with partially mapped right-most range",
		ns, 6, 11, "ß", true)
	run("valid result with completely mapped range", ns, 4, 7, "süß", true)
	run("valid result with full string range", ns, 0, 11, "süß", true)

	ns = NormalizedString{
		original:   "süß",
		normalized: "süß",
		alignments: []NormalizedStringAlignment{{0, 1}, {1, 2}, {2, 3}},
	}

	run("unmodified string - valid result with a left-most range",
		ns, 0, 1, "s", true)
	run("unmodified string - valid result with a right-most range",
		ns, 2, 3, "ß", true)
	run("unmodified string - valid result with a range in the middle",
		ns, 1, 2, "ü", true)
	run("unmodified string - valid result with full string range",
		ns, 0, 3, "süß", true)

	ns = NormalizedString{
		original:   "süß",
		normalized: "sß",
		alignments: []NormalizedStringAlignment{{0, 2}, {2, 3}},
	}

	run("string with deletion - valid result with a left-most range",
		ns, 0, 1, "sü", true)
	run("string with deletion - valid result with a right-most range",
		ns, 1, 2, "ß", true)
	run("string with deletion - valid result with full string range",
		ns, 0, 2, "süß", true)

	ns = NormalizedString{
		original:   "foo",
		normalized: "",
		alignments: []NormalizedStringAlignment{},
	}

	run("blank result with an empty normalized string", ns, 0, 0, "", false)

	ns = NormalizedString{
		original:   "",
		normalized: "foo",
		alignments: []NormalizedStringAlignment{{0, 0}, {0, 0}, {0, 0}},
	}

	run("blank result with an empty original string", ns, 0, 3, "", false)
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
			alignments: []NormalizedStringAlignment{},
		})

	run("non-empty string, empty changes", NewNormalizedString("Bar"),
		[]RuneChanges{}, 0,
		NormalizedString{
			original:   "Bar",
			normalized: "",
			alignments: []NormalizedStringAlignment{},
		})

	run("non-empty string, empty changes, offset", NewNormalizedString("Bar"),
		[]RuneChanges{}, 3,
		NormalizedString{
			original:   "Bar",
			normalized: "",
			alignments: []NormalizedStringAlignment{},
		})

	run("1:1 mapping (all Changes = 0)", NewNormalizedString("Bär"),
		[]RuneChanges{{'S', 0}, {'ü', 0}, {'ß', 0}}, 0,
		NormalizedString{
			original:   "Bär",
			normalized: "Süß",
			alignments: []NormalizedStringAlignment{{0, 1}, {1, 2}, {2, 3}},
		})

	run("1:1 mapping (all Changes == 0), with offset",
		NewNormalizedString("Bär"),
		[]RuneChanges{{'ü', 0}, {'ß', 0}}, 1,
		NormalizedString{
			original:   "Bär",
			normalized: "üß",
			alignments: []NormalizedStringAlignment{{0, 2}, {2, 3}},
		})

	run("adding a rune and deleting the rest (only one Change = 1)",
		NewNormalizedString("Bar"),
		[]RuneChanges{{'x', 1}}, 0,
		NormalizedString{
			original:   "Bar",
			normalized: "x",
			alignments: []NormalizedStringAlignment{{0, 0}},
		})

	run("adding a rune at the beginning", NewNormalizedString("x"),
		[]RuneChanges{{'a', 1}, {'x', 0}}, 0,
		NormalizedString{
			original:   "x",
			normalized: "ax",
			alignments: []NormalizedStringAlignment{{0, 0}, {0, 1}},
		})

	run("adding more runes at the beginning", NewNormalizedString("x"),
		[]RuneChanges{{'a', 1}, {'b', 1}, {'x', 0}}, 0,
		NormalizedString{
			original:   "x",
			normalized: "abx",
			alignments: []NormalizedStringAlignment{{0, 0}, {0, 0}, {0, 1}},
		})

	run("adding a rune at the end", NewNormalizedString("x"),
		[]RuneChanges{{'x', 0}, {'a', 1}}, 0,
		NormalizedString{
			original:   "x",
			normalized: "xa",
			alignments: []NormalizedStringAlignment{
				// FIXME: shouldn't be {0, 1}, {1, 1}?
				{0, 1}, {0, 1},
			},
		})

	run("adding more runes at the end", NewNormalizedString("x"),
		[]RuneChanges{{'x', 0}, {'a', 1}, {'b', 1}}, 0,
		NormalizedString{
			original:   "x",
			normalized: "xab",
			alignments: []NormalizedStringAlignment{
				// FIXME: shouldn't be {0, 1}, {1, 1}, {1, 1}?
				{0, 1}, {0, 1}, {0, 1},
			},
		})

	run("adding runes at begin and end", NewNormalizedString("x"),
		[]RuneChanges{{'a', 1}, {'x', 0}, {'b', 1}}, 0,
		NormalizedString{
			original:   "x",
			normalized: "axb",
			alignments: []NormalizedStringAlignment{
				// FIXME: shouldn't be {0, 0}, {0, 1}, {1, 1}?
				{0, 0}, {0, 1}, {0, 1},
			},
		})

	run("adding a rune in the middle", NewNormalizedString("ab"),
		[]RuneChanges{{'a', 0}, {'x', 1}, {'b', 0}}, 0,
		NormalizedString{
			original:   "ab",
			normalized: "axb",
			alignments: []NormalizedStringAlignment{
				// FIXME: shouldn't be {0, 1}, {1, 1}, {1, 2}?
				{0, 1}, {0, 1}, {1, 2},
			},
		})

	run("adding multiple runes in the middle", NewNormalizedString("ab"),
		[]RuneChanges{{'a', 0}, {'x', 1}, {'y', 1}, {'b', 0}}, 0,
		NormalizedString{
			original:   "ab",
			normalized: "axyb",
			alignments: []NormalizedStringAlignment{
				// FIXME: shouldn't be {0, 1}, {1, 1}, {1, 1}, {1, 2}?
				{0, 1}, {0, 1}, {0, 1}, {1, 2},
			},
		})

	run("Change -1 at the beginning", NewNormalizedString("Bar"),
		[]RuneChanges{{'Q', -1}, {'r', 0}}, 0,
		NormalizedString{
			original:   "Bar",
			normalized: "Qr",
			alignments: []NormalizedStringAlignment{{0, 2}, {2, 3}},
		})

	run("Change -1 at the end", NewNormalizedString("Bar"),
		[]RuneChanges{{'B', 0}, {'x', -1}}, 0,
		NormalizedString{
			original:   "Bar",
			normalized: "Bx",
			alignments: []NormalizedStringAlignment{{0, 1}, {1, 3}},
		})

	run("Change -1 in the middle", NewNormalizedString("abcd"),
		[]RuneChanges{{'a', 0}, {'x', -1}, {'d', 0}}, 0,
		NormalizedString{
			original:   "abcd",
			normalized: "axd",
			alignments: []NormalizedStringAlignment{{0, 1}, {1, 3}, {3, 4}},
		})

	run("Change -2 in the middle", NewNormalizedString("abcde"),
		[]RuneChanges{{'a', 0}, {'x', -2}, {'e', 0}}, 0,
		NormalizedString{
			original:   "abcde",
			normalized: "axe",
			alignments: []NormalizedStringAlignment{{0, 1}, {1, 4}, {4, 5}},
		})
}

func TestNormalizedStringAlignmentEqual(t *testing.T) {
	t.Run("true if `pos` and `changes` are the same", func(t *testing.T) {
		a := NormalizedStringAlignment{pos: 1, changes: 2}
		b := NormalizedStringAlignment{pos: 1, changes: 2}
		if !a.Equal(b) {
			t.Fail()
		}
	})

	t.Run("false if `pos` differ", func(t *testing.T) {
		a := NormalizedStringAlignment{pos: 1, changes: 2}
		b := NormalizedStringAlignment{pos: 3, changes: 2}
		if a.Equal(b) {
			t.Fail()
		}
	})

	t.Run("false if `changes` differ", func(t *testing.T) {
		a := NormalizedStringAlignment{pos: 1, changes: 2}
		b := NormalizedStringAlignment{pos: 1, changes: 3}
		if a.Equal(b) {
			t.Fail()
		}
	})

	t.Run("false if `pos` and `changes` are different", func(t *testing.T) {
		a := NormalizedStringAlignment{pos: 1, changes: 2}
		b := NormalizedStringAlignment{pos: 3, changes: 4}
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
	actual, expected []NormalizedStringAlignment,
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
