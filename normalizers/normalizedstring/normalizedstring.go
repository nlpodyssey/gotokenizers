// Copyright (c) 2020, NLP Odyssey Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package normalizedstring

import (
	"strings"
	"unicode"
)

// NormalizedString takes care of processing an "original" string to modify it
// and obtain a "normalized" string.
//
// It keeps both versions of the string, alignments information between them,
// and provides an interface to retrieve ranges of each string, using offsets
// from any of them.
//
// It is possible to retrieve a part of the original string, by indexing it
// with offsets from the normalized one, and the other way around too.
// It is also possible to convert offsets from one referential string to the
// other one easily.
type NormalizedString struct {
	// The original version of the string, before any modification
	original string
	// The normalized version of the string, after all modifications
	normalized string
	// Mapping from normalized string to original one: (start, end) for each
	// rune of the normalized string
	alignments []AlignmentRange
}

// AlignmentRange represents a (start, end) range used for the alignment
// information of a NormalizedString.
type AlignmentRange struct {
	// Start rune position, inclusive
	start int
	// End rune position, exclusive
	end int
}

// NewNormalizedString builds and returns a new NormalizedString.
//
// It sets both the original and the normalized strings with the value of
// the given string, and initializes the alignments with a 1:1 mapping.
func NewNormalizedString(s string) *NormalizedString {
	return &NormalizedString{
		original:   s,
		normalized: s,
		alignments: newAlignments(s),
	}
}

// newAlignments returns a 1:1 alignment for the given string, treating it
// as the same value for both "normalized" and "original" versions.
func newAlignments(s string) []AlignmentRange {
	alignments := make([]AlignmentRange, len([]rune(s)))
	for index := range alignments {
		alignments[index] = AlignmentRange{
			start: index,
			end:   index + 1,
		}
	}
	return alignments
}

// Get returns the "normalized" version of the NormalizedString.
func (ns *NormalizedString) Get() string {
	return ns.normalized
}

// GetOriginal returns the "original" version of the NormalizedString.
func (ns *NormalizedString) GetOriginal() string {
	return ns.original
}

// Equal reports whether the "normalized" string of the receiver is equal
// to the "normalized" string of the given parameter.
//
// The "original" string and the alignments are ignored for this comparison.
func (ns *NormalizedString) Equal(other *NormalizedString) bool {
	return ns.normalized == other.normalized
}

// Len returns the length of the "normalized" string
// (counting runes, not bytes).
func (ns *NormalizedString) Len() int {
	return len([]rune(ns.normalized))
}

// LenOriginal returns the length of the "original" string
// (counting runes, not bytes).
func (ns *NormalizedString) LenOriginal() int {
	return len([]rune(ns.original))
}

// IsEmpty reports whether the "normalized" string is empty.
func (ns *NormalizedString) IsEmpty() bool {
	return len(ns.normalized) == 0
}

// ConvertOffset converts the given offsets range from one referential string
// to the other one: "original" to "normalized", or "normalized" to "original".
//
// If the range is out of bounds, then it returns (-1, -1, false).
func (ns *NormalizedString) ConvertOffset(nsRange NSRange) (start, end int, ok bool) {
	return nsRange.convertOffset(ns)
}

// GetRange returns a section of the "normalized" string according to the
// given range (remapped, if necessary).
//
// Indexing is on runes, not bytes.
// If the range is out of bounds, then it returns ("", false).
func (ns *NormalizedString) GetRange(nsRange NSRange) (value string, ok bool) {
	start, end, ok := nsRange.normalizedRange(ns)
	if !ok {
		return "", false
	}
	return getRangeOf(ns.normalized, start, end)
}

// GetRangeOriginal returns a section of the "original" string according to the
// given range (remapped, if necessary).
//
// Indexing is on runes, not bytes.
// If the range is out of bounds, then it returns ("", false).
func (ns *NormalizedString) GetRangeOriginal(nsRange NSRange) (value string, ok bool) {
	start, end, ok := nsRange.originalRange(ns)
	if !ok {
		return "", false
	}
	return getRangeOf(ns.original, start, end)
}

// RuneChanges is a single rune-based transformation that can be applied
// by "NormalizedString.Transform".
type RuneChanges struct {
	// Rune is a single rune of the new transformed "normalized" string.
	Rune rune

	// Changes is a number representing how the Rune must be handled when
	// it is inserted in the new "normalized" string.
	//
	// - "1" = this is a new rune
	// - "-N" = the rune is right before N removed runes
	// - "0" = this rune represents the old one (even if changed)
	//
	// Values greater than "1" are not allowed. To add multiple runes, each of
	// them must be represented by a separate RuneChanges value, with Changes
	// set to "1".
	Changes int
}

// Transform applies transformations to the current "normalized" version,
// updating the current alignments with the new ones.
//
// It expects a slice of RuneChanges, which are then applied in sequence to
// transform the "normalized" value.
// Since it is possible that the new result does not include some
// of the characters at the beginning of the original value, we need an
// "initialOffset" which represents the number of runes removed from the very
// beginning of the string.
func (ns *NormalizedString) Transform(dest []RuneChanges, initialOffset int) {
	var strBuilder strings.Builder

	// Pre-fill the new alignments with (0, 0)
	alignments := make([]AlignmentRange, len(dest))
	offset := -initialOffset

	for index, item := range dest {
		changes := item.Changes

		// A positive offset means we added characters. So we need to remove
		// this offset from the current index to find out the previous id
		oldIndex := index - offset

		if changes == 0 {
			// No changes required here
			alignments[index] = ns.alignments[oldIndex]
		} else if changes == 1 {
			// This is a newly inserted character
			offset++
			if oldIndex > 0 {
				prevEnd := ns.alignments[oldIndex-1].end
				alignments[index] = AlignmentRange{start: prevEnd, end: prevEnd}
			} // otherwise, it is already (0, 0)

		} else if changes < 0 { // changes < 0
			// Some characters where removed, nothing to change in alignments
			offset += changes
			alignments[index] = ns.alignments[oldIndex]
		} else { // changes > 1
			panic("invalid Changes > 1")
		}

		// Then we keep only the char for string reconstruction
		strBuilder.WriteRune(item.Rune)
	}

	ns.normalized = strBuilder.String()
	ns.alignments = alignments
}

// Filter applies filtering over the "normalized" string runes.
//
// The given filter callback is called for each rune of the "normalized"
// string. Only the runes for which the callback returns true are kept.
func (ns *NormalizedString) Filter(filter func(rune) bool) {
	runes := []rune(ns.normalized)
	runesLen := len(runes)

	if runesLen == 0 {
		return
	}

	removed := 0
	filtered := make([]RuneChanges, 0, runesLen)

	lastIndex := runesLen - 1
	_ = runes[lastIndex]

	for runeIndex := lastIndex; runeIndex >= 0; runeIndex-- {
		r := runes[runeIndex]

		if filter(r) {
			if removed > 0 {
				filtered = append(filtered, RuneChanges{
					Rune:    r,
					Changes: -removed,
				})
				removed = 0
			} else {
				filtered = append(filtered, RuneChanges{
					Rune:    r,
					Changes: 0,
				})
			}
		} else {
			removed++
		}
	}

	// Reverse `filtered`
	lastIndex = len(filtered) - 1
	if lastIndex > 0 {
		_ = filtered[lastIndex]
		for i, j := 0, lastIndex; i < j; i, j = i+1, j-1 {
			filtered[i], filtered[j] = filtered[j], filtered[i]
		}
	}

	ns.Transform(filtered, removed)
}

// Prepend prepends the given string to the "normalized" value.
func (ns *NormalizedString) Prepend(s string) {
	ns.normalized = s + ns.normalized
	alignments := make([]AlignmentRange, len([]rune(s))) // all (0, 0)
	ns.alignments = append(alignments, ns.alignments...)
}

// Append prepends the given string to the "normalized" value.
func (ns *NormalizedString) Append(s string) {
	ns.normalized += s

	lastOffset := AlignmentRange{} // (0, 0)
	alignmentsLen := len(ns.alignments)
	if alignmentsLen > 0 {
		lastAlignment := ns.alignments[alignmentsLen-1]
		lastOffset = AlignmentRange{
			start: lastAlignment.end,
			end:   lastAlignment.end,
		}
	}

	// TODO: compare this with appending a slice of alignments
	//       of size len([]rune(s))
	for range s { // note that this loops over the string's runes, not bytes
		ns.alignments = append(ns.alignments, lastOffset)
	}
}

// Map maps the runes of the "normalized" string.
//
// Each rune can be replaced with another rune (or the same one).
func (ns *NormalizedString) Map(f func(rune) rune) {
	ns.normalized = strings.Map(f, ns.normalized)
}

// ForEach calls the given function for each rune of the "normalized" string.
// FIXME: this might be inefficient: prefer direct iteration
//        on `range ns.Get()` and eventually remove this method
func (ns *NormalizedString) ForEach(f func(rune)) {
	for _, r := range ns.normalized {
		f(r)
	}
}

// ToLower remaps all Unicode letters of the "normalized" string to their
// lower case.
// FIXME: see Unicode special casing notes on `NormalizedString.ToUpper()`
func (ns *NormalizedString) ToLower() {
	newChars := make([]RuneChanges, 0, ns.Len())
	for _, r := range ns.normalized {
		newChars = append(
			newChars,
			RuneChanges{Rune: unicode.ToLower(r)}, // Changes: 0
		)
	}
	ns.Transform(newChars, 0)
}

// ToUpper remaps all Unicode letters of the "normalized" string to their
// upper case.
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
func (ns *NormalizedString) ToUpper() {
	newChars := make([]RuneChanges, 0, ns.Len())
	for _, r := range ns.normalized {
		newChars = append(
			newChars,
			RuneChanges{Rune: unicode.ToUpper(r)}, // Changes: 0
		)
	}
	ns.Transform(newChars, 0)
}

// SplitOff splits the receiver NormalizedString in two parts. The "normalized"
// version of the receiver is reduced to the interval [0, at), and a new
// NormalizedString is returned with the "normalized" string set to the
// original range [at, len).
//
// The provided `at` is an index on runes, not bytes.
func (ns *NormalizedString) SplitOff(at int) *NormalizedString {
	runes := []rune(ns.normalized)

	newNs := &NormalizedString{
		// Preserve the full original string to have meaningful alignments
		original:   ns.original,
		normalized: string(runes[at:]),
		alignments: ns.alignments[at:],
	}

	ns.normalized = string(runes[:at])
	ns.alignments = ns.alignments[:at]

	return newNs
}

// MergeWith appends the "original" and "normalized strings of the other
// NormalizedString to the values of the receiver.
func (ns *NormalizedString) MergeWith(other *NormalizedString) {
	alignmentsOffset := ns.LenOriginal()

	ns.original += other.original
	ns.normalized += other.normalized

	for _, alignment := range other.alignments {
		ns.alignments = append(ns.alignments, AlignmentRange{
			start: alignment.start + alignmentsOffset,
			end:   alignment.end + alignmentsOffset,
		})
	}
}

// Trim removes leading and trailing spaces from the "normalized" string.
func (ns *NormalizedString) Trim() {
	ns.TrimLeftRight(true, true)
}

// TrimLeft removes leading spaces from the "normalized" string.
func (ns *NormalizedString) TrimLeft() {
	ns.TrimLeftRight(true, false)
}

// TrimRight removes trailing spaces from the "normalized" string.
func (ns *NormalizedString) TrimRight() {
	ns.TrimLeftRight(false, true)
}

// TrimLeftRight removes leading (left) and/or trailing (right) spaces from
// the "normalized" string.
func (ns *NormalizedString) TrimLeftRight(left, right bool) {
	runes := []rune(ns.normalized)
	lenRunes := len(runes)

	leadingSpaces := 0
	if left {
		leadingSpaces = countLeadingSpaces(runes)
	}

	trailingSpaces := 0
	if right && leadingSpaces < lenRunes {
		trailingSpaces = countTrailingSpaces(runes)
	}

	if leadingSpaces == 0 && trailingSpaces == 0 {
		return
	}

	lastIndex := lenRunes - trailingSpaces
	ns.normalized = string(runes[leadingSpaces:lastIndex])
	ns.alignments = ns.alignments[leadingSpaces:lastIndex]
}

// Equal reports whether the bounds of the receiver AlignmentRange are equal
// to the bounds of the other AlignmentRange.
func (nsa *AlignmentRange) Equal(other AlignmentRange) bool {
	return nsa.start == other.start && nsa.end == other.end
}

// getRangeOf returns a section of the given string according to the given
// bounds.
//
// Indexing is on runes, not bytes.
// If the range is out of bounds, then it returns ("", false).
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

// countLeadingSpaces returns the number of leading spaces in the given slice
// of runes, if any.
func countLeadingSpaces(runes []rune) int {
	for i, r := range runes {
		if !unicode.In(r, unicode.White_Space) {
			return i
		}
	}
	return len(runes)
}

// countTrailingSpaces returns the number of trailing spaces in the given slice
// of runes, if any.
func countTrailingSpaces(runes []rune) int {
	runesLen := len(runes)
	if runesLen == 0 {
		return 0
	}

	lastIndex := runesLen - 1
	_ = runes[lastIndex]

	for i := lastIndex; i >= 0; i-- {
		if !unicode.In(runes[i], unicode.White_Space) {
			return lastIndex - i
		}
	}

	return len(runes)
}
