// Copyright (c) 2020, The GoTokenizers Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package normalized_string

import (
	"strings"
	"unicode"
)

// A normalized string takes care of keeping both versions of a `string`, and
// provides necessary alignments to retrieve ranges of both strings.
type NormalizedString struct {
	original   string
	normalized string

	// Mapping from the normalized string to the original one
	alignments []NormalizedStringAlignment
}

// A single alignment information for `NormalizedString`
type NormalizedStringAlignment struct {
	// The position in the modified string
	pos int

	// The number of insertions or deletions
	changes int
}

func NewNormalizedString(s string) NormalizedString {
	return NormalizedString{
		original:   s,
		normalized: s,
		alignments: newAlignments(s),
	}
}

func (ns *NormalizedString) Equal(other NormalizedString) bool {
	return ns.normalized == other.normalized
}

// Returns the length of the normalized string (counting runes ,not bytes)
func (ns *NormalizedString) Len() int {
	return len([]rune(ns.normalized))
}

// Returns the length of the original string (counting runes, not bytes)
func (ns *NormalizedString) LenOriginal() int {
	return len([]rune(ns.original))
}

func (ns *NormalizedString) IsEmpty() bool {
	return len(ns.normalized) == 0
}

// Returns the normalized string.
func (ns *NormalizedString) Get() string {
	return ns.normalized
}

// Returns the original string.
func (ns *NormalizedString) GetOriginal() string {
	return ns.original
}

// Returns a range of the normalized string (indexing on runes, not bytes)
func (ns *NormalizedString) GetRange(start, end int) (string, bool) {
	return getRangeOf(ns.normalized, start, end)
}

// Return a range of the original string, using a range from the
// normalized string
func (ns *NormalizedString) GetRangeOriginal(start, end int) (string, bool) {
	originalStart, originalEnd := ns.GetOriginalOffsets(start, end)
	if originalStart == -1 {
		return "", false
	}
	return getRangeOf(ns.original, originalStart, originalEnd)
}

// Returns the `(start, end)` range of the original string corresponding to the
// received range on the normalized string.
// Returns `(-1, -1)` if out of bounds.
func (ns *NormalizedString) GetOriginalOffsets(start, end int) (int, int) {
	if end <= start || start < 0 || end > len(ns.alignments) {
		return -1, -1
	}
	return ns.alignments[start].pos, ns.alignments[end-1].changes
}

// See `NormalizedString.Transform`.
type RuneChanges struct {
	// A rune of the normalized string
	Rune rune

	// `1` = this is a new rune
	// `-N` = the rune is right before N removed runes
	// `0` = this rune represents the old one (even if changed)
	// Values greater than `1` are not allowed: if multiple chars are added,
	// each of them must have a `change` of `1`.
	Changes int
}

// Applies transformations to the current normalized version, updating the
// current alignments with the new ones.
//
// Since it is possible that the normalized string doesn't include some of the
// runes at the beginning of the original one, we need an `initialOffset` which
// represents the number of removed runes at the very beginning.
func (ns *NormalizedString) Transform(dest []RuneChanges, initialOffset int) {
	offset := 0
	remainingOffset := initialOffset

	var strBuilder strings.Builder
	alignments := make([]NormalizedStringAlignment, 0, len(ns.alignments))

	for index, runeSizePair := range dest {
		c := runeSizePair.Rune
		changes := runeSizePair.Changes

		if remainingOffset != 0 {
			changes -= remainingOffset
			remainingOffset = 0
		}

		var uof int
		if offset < 0 {
			uof = -offset
		} else {
			uof = offset
		}

		// A positive offset means we added characters. So we need to remove
		// this offset from the current index to find out the previous id
		var idx int
		if offset < 0 {
			idx = index + uof
		} else {
			idx = index - uof
		}

		alignmentPair := NormalizedStringAlignment{-1, -1}

		if changes > 0 {
			// This is a newly inserted character, so we use the alignment
			// from the previous one
			offset += 1
			if idx < 1 {
				alignmentPair.pos = 0
				alignmentPair.changes = 0
			} else {
				alignmentPair = ns.alignments[idx-1]
			}

		} else if changes == 0 {
			alignmentPair = ns.alignments[idx]

		} else { // changes < 0
			// Some characters where removed, so we merge our range with the
			// one from the removed characters as the new alignment
			uch := -changes
			offset += changes

			lastIndex := idx + uch

			alignmentPair = ns.alignments[idx]

			for alIndex := idx + 1; alIndex <= lastIndex; alIndex++ {
				al := ns.alignments[alIndex]
				if al.pos < alignmentPair.pos {
					alignmentPair.pos = al.pos
				}
				if al.changes < alignmentPair.pos {
					alignmentPair.pos = al.changes
				}
				if al.pos > alignmentPair.changes {
					alignmentPair.changes = al.pos
				}
				if al.changes > alignmentPair.changes {
					alignmentPair.changes = al.changes
				}
			}
		}

		// Then we keep only the char for string reconstruction
		strBuilder.WriteRune(c)

		if alignmentPair.pos == -1 {
			// TODO: is this really possible??
			panic("Bad alignement in NormalizedString.Transform")
		}
		alignments = append(alignments, alignmentPair)
	}

	ns.normalized = strBuilder.String()
	ns.alignments = alignments
}

func (ns *NormalizedString) Filter(filter func(rune) bool) {
	runes := []rune(ns.normalized)
	runesLen := len(runes)

	if runesLen == 0 {
		return
	}

	removed := 0
	filtered := make([]RuneChanges, 0, runesLen)

	lastIndex := len(runes) - 1
	_ = runes[lastIndex]

	for runeIndex := lastIndex; runeIndex >= 0; runeIndex-- {
		c := runes[runeIndex]
		keep := filter(c)
		if keep {
			if removed > 0 {
				filtered = append(filtered, RuneChanges{
					Rune:    c,
					Changes: -removed,
				})
				removed = 0
			} else {
				filtered = append(filtered, RuneChanges{
					Rune:    c,
					Changes: 0,
				})
			}
		} else {
			removed += 1
		}
	}

	// Reverse `filtered`
	lastIndex = len(filtered) - 1
	if lastIndex >= 0 {
		_ = filtered[lastIndex]
		for i, j := 0, lastIndex; i < j; i, j = i+1, j-1 {
			filtered[i], filtered[j] = filtered[j], filtered[i]
		}
	}

	ns.Transform(filtered, removed)
}

func (ns *NormalizedString) Prepend(s string) {
	ns.normalized = s + ns.normalized
	runes := []rune(s)
	runesLen := len(runes)
	alignments := make([]NormalizedStringAlignment, runesLen+len(ns.alignments))
	// By default, all the new alignments have already {pos: 0, changes: 0}
	ns.alignments = append(alignments, ns.alignments...)
}

func (ns *NormalizedString) Append(s string) {
	ns.normalized += s

	lastOffset := NormalizedStringAlignment{} // {pos: 0, changes: 0}
	alignmentsLen := len(ns.alignments)
	if alignmentsLen > 0 {
		lastAlignment := ns.alignments[alignmentsLen-1]
		lastOffset = NormalizedStringAlignment{
			pos:     lastAlignment.changes,
			changes: lastAlignment.changes,
		}
	}

	for range s { // note that this loops over the string's runes, not bytes
		ns.alignments = append(ns.alignments, lastOffset)
	}
}

// FIXME: this might be inefficient: prefer direct use of `strings.Map`
//        on `ns.Get()` and eventually remove this method
func (ns *NormalizedString) Map(f func(rune) rune) {
	ns.normalized = strings.Map(f, ns.normalized)
}

// FIXME: this might be inefficient: prefer direct iteration
//        on `range ns.Get()` and eventually remove this method
func (ns *NormalizedString) ForEach(f func(rune)) {
	for _, r := range ns.normalized {
		f(r)
	}
}

// FIXME: see Unicode special casing notes on `NormalizedString.Uppercase()`
func (ns *NormalizedString) Lowercase() {
	newChars := make([]RuneChanges, 0, ns.Len())
	for _, r := range ns.normalized {
		newChars = append(
			newChars,
			RuneChanges{Rune: unicode.ToLower(r)}, // Changes: 0
		)
	}
	ns.Transform(newChars, 0)
}

// FIXME: Go `unicode` package does not consider Unicode special casing
//        (see https://www.unicode.org/Public/UCD/latest/ucd/SpecialCasing.txt)
//        As a result, every single rune is always transformed to a single
//        new rune. This is not the case in the original rust implementation.
//        Should we consider using another "better" package?
//
//        Explanatory example:
//        ```rust
//        fn main() {
//            let a = '\u{00DF}';
//            println!("{} {:?}", a, a.to_string().chars().count());
//            // => ß 1
//            let b = a.to_uppercase();
//            println!("{} {}", b, b.to_string().chars().count());
//            // => SS 2
//        }
//        ```
//
//        ```go
//        package main
//
//        import (
// 	          "fmt"
// 	          "unicode"
//        )
//
//        func main() {
// 	          a := '\u00DF'
// 	          fmt.Println(string(a), len([]rune(string(a))))
//            // => ß 1
// 	          b := unicode.ToUpper(a)
// 	          fmt.Println(string(b), len([]rune(string(b))))
//            // => ß 1
//        }
//        ```
func (ns *NormalizedString) Uppercase() {
	newChars := make([]RuneChanges, 0, ns.Len())
	for _, r := range ns.normalized {
		newChars = append(
			newChars,
			RuneChanges{Rune: unicode.ToUpper(r)}, // Changes: 0
		)
	}
	ns.Transform(newChars, 0)
}

// Split off the normalized string, returning a new NormalizedString that
// contains the range [at, len). The original NormalizedString itself will
// then contain the range [0, at).
// The provided `at` indexes on `rune`, not bytes.
func (ns *NormalizedString) SplitOff(at int) NormalizedString {
	runes := []rune(ns.normalized)
	runesLen := len(runes)

	if at > runesLen {
		return NewNormalizedString("")
	}

	// Split normalized

	byteIndex := len(string(runes[:at]))

	normalized := ns.normalized[byteIndex:]
	ns.normalized = ns.normalized[:byteIndex]

	alignments := ns.alignments[at:]
	ns.alignments = ns.alignments[:at]

	// Split original

	originalAt := 0
	alignmentsLen := len(ns.alignments)
	if alignmentsLen > 0 {
		originalAt = ns.alignments[alignmentsLen-1].changes
	}

	originalRunes := []rune(ns.original)
	originalByteIndex := len(string(originalRunes[:originalAt]))

	original := ns.original[originalByteIndex:]
	ns.original = ns.original[:originalByteIndex]

	return NormalizedString{
		original:   original,
		normalized: normalized,
		alignments: alignments,
	}
}

func (ns *NormalizedString) MergeWith(other NormalizedString) {
	ns.original += other.original
	ns.normalized += other.normalized

	nsLen := ns.Len()
	for _, alignment := range other.alignments {
		ns.alignments = append(ns.alignments, NormalizedStringAlignment{
			pos:     alignment.pos + nsLen,
			changes: alignment.changes + nsLen,
		})
	}
}

func (nsa *NormalizedStringAlignment) Equal(
	other NormalizedStringAlignment,
) bool {
	return nsa.pos == other.pos && nsa.changes == other.changes
}

func newAlignments(s string) []NormalizedStringAlignment {
	alignments := make([]NormalizedStringAlignment, len([]rune(s)))

	for index := range alignments {
		alignments[index] = NormalizedStringAlignment{
			pos:     index,
			changes: index + 1,
		}
	}

	return alignments
}

func getRangeOf(s string, start, end int) (string, bool) {
	if end <= start || start < 0 {
		return "", false
	}

	runes := []rune(s)
	if end > len(runes) {
		return "", false
	}

	return string(runes[start:end]), true
}
