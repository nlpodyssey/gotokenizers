// Copyright (c) 2020, NLP Odyssey Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package normalizedstring

import (
	"fmt"
	"github.com/nlpodyssey/gotokenizers/splitpattern"
	"github.com/nlpodyssey/gotokenizers/strutils"
	"strings"
	"unicode"
)

// NormalizedString takes care of processing an "original" string to modify
// it and obtain a "normalized" string.
//
// It keeps both versions of the string, alignments information between both
// of them, and provides an interface to retrieve ranges of each string,
// using offsets from any of them.
//
// It is possible to retrieve a part of the original string, by indexing it
// with offsets from the normalized one, and the other way around too.
// It is also possible to convert offsets from one referential to the other
// one easily.
type NormalizedString struct {
	// The original version of the string, before any modification.
	original string
	// The normalized version of the string, after all modifications.
	normalized string
	// Mapping from the normalized string to the original one: (start, end)
	// for each byte of the normalized string.
	alignments []AlignmentRange
	// If this NormalizedString is a slice of a bigger one, we keep track
	// of the missing part, so that we can still give offsets from this
	// original string.
	originalShift int
}

// AlignmentRange represents a (start, end) range for representing the
// alignments of a NormalizedString.
type AlignmentRange struct {
	// Start byte position, inclusive.
	start int
	// End byte position, exclusive.
	end int
}

// New returns a new NormalizedString.
func New(
	original string,
	normalized string,
	alignments []AlignmentRange,
	originalShift int,
) *NormalizedString {
	return &NormalizedString{
		original:      original,
		normalized:    normalized,
		alignments:    alignments,
		originalShift: originalShift,
	}
}

// FromString returns a new NormalizedString built from the given string.
func FromString(s string) *NormalizedString {
	return &NormalizedString{
		original:      s,
		normalized:    s,
		alignments:    alignmentsFromString(s),
		originalShift: 0,
	}
}

func alignmentsFromString(s string) []AlignmentRange {
	alignments := make([]AlignmentRange, 0, len(s))

	for runeIndex, r := range s {
		strRune := string(r)
		bytesLen := len(strRune)
		for byteIndex := 0; byteIndex < bytesLen; byteIndex++ {
			alignments = append(alignments, AlignmentRange{
				start: runeIndex,
				end:   runeIndex + bytesLen,
			})
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

// Len returns the length in bytes of the "normalized" string.
func (ns *NormalizedString) Len() int {
	return len(ns.normalized)
}

// OriginalLen returns the length in bytes of the "original" string.
func (ns *NormalizedString) OriginalLen() int {
	return len(ns.original)
}

// IsEmpty reports whether the "normalized" string is empty.
func (ns *NormalizedString) IsEmpty() bool {
	return len(ns.normalized) == 0
}

// OriginalOffsets returns the original offsets.
func (ns *NormalizedString) OriginalOffsets() strutils.ByteOffsets {
	return strutils.ByteOffsets{
		Start: ns.originalShift,
		End:   ns.originalShift + ns.OriginalLen(),
	}
}

// RuneChange is a single rune-based transformation that can be applied
// by NormalizedString.TransformRange
type RuneChange struct {
	// Rune is a rune of the new transformed "normalized" string.
	Rune rune

	// Change is a number representing how the Rune must be handled when
	// it is inserted in the new "normalized" string.

	// - "1" = this is a new rune
	// - "-N" = the rune is right before N removed runes
	// - "0" = this rune is replacing the existing one

	// Values greater than "1" are not allowed. To add multiple runes, each of
	// them must be represented by a separate RuneChange value, with Changes
	// set to "1".
	Change int
}

// TransformRange applies transformations to the current normalized version of
// the string, while updating the alignments.
//
// This method expect a slice of RuneChange, where each item corresponds to one
// rune of the new normalized string.
//
// Since it is possible that the normalized string doesn't include some of the
// characters at the beginning of the original one, we need an initialOffset
// which represents the number of removed chars at the very beginning.
func (ns *NormalizedString) TransformRange(
	rng Range,
	dest []RuneChange,
	initialOffset int,
) {
	normalizedRange, normalizedRangeOk := ns.CoerceRangeToNormalized(rng)
	if !normalizedRangeOk {
		return
	}

	// Retrieve the original characters that are being replaced. This lets us
	// compute the change in byte sizes along the way.
	replacedNormalizedString := ns.normalized[normalizedRange.Start():normalizedRange.End()]

	initialRemoved := len(replacedNormalizedString[:initialOffset])
	replacedNormalizedRunes := []rune(replacedNormalizedString[initialOffset:])

	offset := initialRemoved + normalizedRange.Start()
	alignments := make([]AlignmentRange, 0, normalizedRange.Len())

	var normalizedBuilder strings.Builder

	for _, runeChange := range dest {
		idx := offset
		var align AlignmentRange

		if runeChange.Change >= 1 {
			if idx < 1 {
				align = AlignmentRange{start: 0, end: 0}
			} else {
				// This is a newly inserted character, so it shares the same
				// alignment than the previous one.
				align = ns.alignments[idx-1]
			}
		} else {
			align = ns.alignments[idx]
		}

		// If we are replacing a character, find it and compute the change in size
		replacedRuneSize := 0
		if runeChange.Change <= 0 {
			replacedRune := replacedNormalizedRunes[0]
			replacedNormalizedRunes = replacedNormalizedRunes[1:]
			replacedRuneSize = len(string(replacedRune))
		}

		// If we are removing some characters, find them too
		totalBytesToRemove := 0
		if runeChange.Change <= -1 {
			totalBytesToRemove = len(string(replacedNormalizedRunes[:-runeChange.Change]))
			replacedNormalizedRunes = replacedNormalizedRunes[-runeChange.Change:]
		}

		// Keep track of the changes for next offsets
		offset += replacedRuneSize
		offset += totalBytesToRemove

		newBytesCount := len(string(runeChange.Rune))
		for i := 0; i < newBytesCount; i++ {
			alignments = append(alignments, align)
		}

		// Then we keep only the rune for string reconstruction
		normalizedBuilder.WriteRune(runeChange.Rune)
	}

	newAlignments := make([]AlignmentRange, 0, len(ns.alignments)-normalizedRange.Len()+len(alignments))
	newAlignments = append(newAlignments, ns.alignments[:normalizedRange.start]...)
	newAlignments = append(newAlignments, alignments...)
	newAlignments = append(newAlignments, ns.alignments[normalizedRange.end:]...)
	ns.alignments = newAlignments

	ns.normalized = fmt.Sprintf("%s%s%s",
		ns.normalized[:normalizedRange.start],
		normalizedBuilder.String(),
		ns.normalized[normalizedRange.end:],
	)
}

// Transform applies transformations to the current normalized version of
// the string, while updating the alignments.
//
// It is the same as calling TransformRange with the full normalized string range.
func (ns *NormalizedString) Transform(dest []RuneChange, initialOffset int) {
	ns.TransformRange(NewNormalizedRange(0, ns.Len()), dest, initialOffset)
}

// OriginalAlignments recalculates the original alignments.
func (ns *NormalizedString) OriginalAlignments() []AlignmentRange {
	// (start, end) are in alignments
	// (offset, length) are in originalAlignments
	originalAlignments := make([]AlignmentRange, 0, len(ns.original))

	// Eventual gap before first group

	if start := ns.alignments[0].start; start > 0 {
		for i := 0; i < start; i++ {
			originalAlignments = append(originalAlignments, AlignmentRange{
				start: 0,
				end:   0,
			})
		}
	}

	last := ns.alignments[0]
	offset := 0
	length := 0

	for _, alignment := range ns.alignments {
		if last == alignment {
			// This is the same group
			length++
		} else {
			// This is a new group
			if alignment.start < last.end {
				panic("Invalid overlapping ranges.")
			}

			// Add the old group
			lastLen := last.end - last.start
			for i := 0; i < lastLen; i++ {
				originalAlignments = append(originalAlignments, AlignmentRange{
					start: offset,
					end:   offset + length,
				})
			}
			offset += length
			length = 1

			// Eventual gap between the 2 groups
			gapLen := alignment.start - last.end
			for i := 0; i < gapLen; i++ {
				originalAlignments = append(originalAlignments, AlignmentRange{
					start: offset,
					end:   offset,
				})
			}
		}

		last = alignment
	}

	// Add the last group
	lastLen := last.end - last.start
	for i := 0; i < lastLen; i++ {
		originalAlignments = append(originalAlignments, AlignmentRange{
			start: offset,
			end:   offset + length,
		})
	}

	// Add eventual last gap
	offset += length

	lastGapLen := len(ns.original) - len(originalAlignments)
	for i := 0; i < lastGapLen; i++ {
		originalAlignments = append(originalAlignments, AlignmentRange{
			start: offset,
			end:   offset,
		})
	}

	return originalAlignments
}

// GetRange returns a range of the normalized string, and a flag which reports
// whether the operation is successful or not.
func (ns *NormalizedString) GetRange(rng Range) (string, bool) {
	nr, ok := ns.CoerceRangeToNormalized(rng)
	if !ok {
		return "", false
	}
	return ns.normalized[nr.Start():nr.End()], true
}

// GetOriginalRange returns a range of the original string, and a flag which
// reports whether the operation is successful or not.
func (ns *NormalizedString) GetOriginalRange(rng Range) (string, bool) {
	or, ok := ns.CoerceRangeToOriginal(rng)
	if !ok {
		return "", false
	}
	return ns.original[or.Start():or.End()], true
}

// CoerceRangeToNormalized coerces the given Range (either an OriginalRange or
// a NormalizedRange) to a NormalizedRange, performing a conversion if needed.
// It also returns a flag which reports whether the operation was successful.
//
// The operation is unsuccessful if the range is targeting something out of
// range.
func (ns *NormalizedString) CoerceRangeToNormalized(r Range) (NormalizedRange, bool) {
	// If the string range is already in the normalized referential, return it as it is
	if nr, isNormalized := r.(NormalizedRange); isNormalized {
		return nr, true
	}

	// If we target an empty range, let's return the same
	if r.Len() == 0 {
		return NewNormalizedRange(r.Start(), r.End()), true
	}

	// If the target goes reverse, return invalid status
	if r.Start() > r.End() {
		return NewNormalizedRange(0, 0), false
	}

	// If we target (0, 0) on an empty string, we want to expand to the entire equivalent
	if len(ns.original) == 0 && r.Start() == 0 && r.End() == 0 {
		return NewNormalizedRange(0, ns.Len()), true
	}

	start := -1
	end := -1

	for i, alignment := range ns.alignments {
		if alignment.end > r.End() {
			break
		}

		if start == -1 && alignment.start >= r.Start() {
			// For now, don't update if width == 0
			if alignment.start != alignment.end {
				start = i
			}
		}

		if alignment.end <= r.End() {
			end = i + 1
		}
	}

	switch {
	case start != -1 && end != -1:
		return NewNormalizedRange(start, end), true
	case start != -1:
		return NewNormalizedRange(start, start), true
	case end != -1:
		return NewNormalizedRange(end, end), true
	default:
		return NewNormalizedRange(0, 0), false
	}
}

// CoerceRangeToOriginal coerces the given Range (either an OriginalRange or
// a NormalizedRange) to a OriginalRange, performing a conversion if needed.
// It also returns a flag which reports whether the operation was successful.
//
// The operation is unsuccessful if the range is targeting something out of
// range.
func (ns *NormalizedString) CoerceRangeToOriginal(r Range) (OriginalRange, bool) {
	// If the string range is already in the original referential, return it as it is
	if or, isOriginal := r.(OriginalRange); isOriginal {
		return or, true
	}

	// If we target an empty range, let's return the same
	if r.Len() == 0 {
		return NewOriginalRange(r.Start(), r.End()), true
	}

	// If the target goes reverse, return invalid status
	if r.Start() > r.End() {
		return NewOriginalRange(0, 0), false
	}

	// If we target (0, 0) on an empty string, we want to expand to the entire equivalent
	if len(ns.normalized) == 0 && r.Start() == 0 && r.End() == 0 {
		return NewOriginalRange(0, ns.OriginalLen()), true
	}

	alignments := ns.alignments[r.Start():r.End()]

	// Expand alignments, returning the range covered by this slice
	return NewOriginalRange(alignments[0].start, alignments[len(alignments)-1].end), true
}

// Prepend prepends the given string to the NormalizedString.
//
// FIXME: Prepend does nothing if the normalized string is empty
func (ns *NormalizedString) Prepend(s string) {
	if len(ns.normalized) == 0 {
		return
	}

	transformations := make([]RuneChange, 0, len(s)+1)

	for i, r := range s {
		change := 0
		if i > 0 {
			change = 1
		}
		transformations = append(transformations, RuneChange{Rune: r, Change: change})
	}

	firstRune := []rune(ns.normalized)[0]
	firstRuneBytesLen := len(string(firstRune))
	transformations = append(transformations, RuneChange{Rune: firstRune, Change: 1})

	ns.TransformRange(NewNormalizedRange(0, firstRuneBytesLen), transformations, 0)
}

// Append appends the given string to the NormalizedString.
//
// FIXME: Append does nothing if the normalized string is empty
func (ns *NormalizedString) Append(s string) {
	if len(ns.normalized) == 0 {
		return
	}

	transformations := make([]RuneChange, 0, len(s)+1)

	runes := []rune(ns.normalized)
	lastRuneIndex := len(runes) - 1
	lastRune := runes[lastRuneIndex]

	transformations = append(transformations, RuneChange{Rune: lastRune, Change: 0})

	for _, r := range s {
		transformations = append(transformations, RuneChange{Rune: r, Change: 1})
	}

	// FIXME: it's probably wrong to pass the rune index, instead of byte index
	ns.TransformRange(NewNormalizedRange(lastRuneIndex, ns.Len()), transformations, 0)
}

// SplitDelimiterBehavior is used by NormalizedString.Split to define the
// expected behavior for the delimiter.
//
// For example, when splitting on '-' with input `the-final--countdown`:
// - SplitDelimiterRemoved            => [ "the", "final", "countdown" ]
// - SplitDelimiterIsolated           => [ "the", "-", "final", "-", "-", "countdown" ]
// - SplitDelimiterMergedWithPrevious => [ "the-", "final-", "-", "countdown" ]
// - SplitDelimiterMergedWithNext     => [ "the", "-final", "-", "-countdown" ]
// - SplitDelimiterContiguous         => [ "the", "-", "final", "--", "countdown" ]
type SplitDelimiterBehavior uint8

const (
	SplitDelimiterRemoved            SplitDelimiterBehavior = iota
	SplitDelimiterIsolated                                  = iota
	SplitDelimiterMergedWithPrevious                        = iota
	SplitDelimiterMergedWithNext                            = iota
	SplitDelimiterContiguous                                = iota
)

// Split splits the current string in many subparts.
func (ns *NormalizedString) Split(
	pattern splitpattern.SplitPattern,
	behaviour SplitDelimiterBehavior,
) ([]*NormalizedString, error) {
	captures, err := pattern.FindMatches(ns.normalized)
	if err != nil {
		return nil, err
	}

	type SplitMatch struct {
		Offsets      strutils.ByteOffsets
		ShouldRemove bool
	}

	splits := make([]SplitMatch, 0, len(captures))

	// Process the matches according to the selected behavior
	switch behaviour {
	case SplitDelimiterRemoved:
		for _, capture := range captures {
			splits = append(splits, SplitMatch{
				Offsets:      capture.Offsets,
				ShouldRemove: capture.IsMatch,
			})
		}
	case SplitDelimiterIsolated:
		for _, capture := range captures {
			splits = append(splits, SplitMatch{
				Offsets:      capture.Offsets,
				ShouldRemove: false,
			})
		}
	case SplitDelimiterMergedWithPrevious:
		previousMatch := false
		for _, capture := range captures {
			if capture.IsMatch && !previousMatch && len(splits) > 0 {
				splits[len(splits)-1].Offsets.End = capture.Offsets.End
			} else {
				splits = append(splits, SplitMatch{
					Offsets:      capture.Offsets,
					ShouldRemove: false,
				})
			}
			previousMatch = capture.IsMatch
		}
	case SplitDelimiterMergedWithNext:
		previousMatch := false
		for i := len(captures) - 1; i >= 0; i-- {
			capture := captures[i]

			if capture.IsMatch && !previousMatch && len(splits) > 0 {
				splits[len(splits)-1].Offsets.Start = capture.Offsets.Start
			} else {
				splits = append(splits, SplitMatch{
					Offsets:      capture.Offsets,
					ShouldRemove: false,
				})
			}
			previousMatch = capture.IsMatch
		}
		// Reverse splits
		for i := len(splits)/2 - 1; i >= 0; i-- {
			opp := len(splits) - 1 - i
			splits[i], splits[opp] = splits[opp], splits[i]
		}
	case SplitDelimiterContiguous:
		previousMatch := false
		for _, capture := range captures {
			if capture.IsMatch == previousMatch && len(splits) > 0 {
				splits[len(splits)-1].Offsets.End = capture.Offsets.End
			} else {
				splits = append(splits, SplitMatch{
					Offsets:      capture.Offsets,
					ShouldRemove: false,
				})
			}
			previousMatch = capture.IsMatch
		}
	default:
		panic(fmt.Sprintf("unexpected SplitDelimiterBehavior: %d", behaviour))
	}

	// Then we split according to the computed splits

	result := make([]*NormalizedString, 0)
	for _, split := range splits {
		if split.ShouldRemove {
			continue
		}
		sliced, ok := ns.Slice(NewNormalizedRange(split.Offsets.Start, split.Offsets.End))
		if !ok {
			return nil, fmt.Errorf("NormalizedString bad split")
		}
		result = append(result, sliced)
	}
	return result, nil
}

// Replace replaces anything that matches the pattern with the given content.
func (ns *NormalizedString) Replace(
	pattern splitpattern.SplitPattern,
	content string,
) error {
	offset := 0

	captures, err := pattern.FindMatches(ns.normalized)
	if err != nil {
		return err
	}
	for _, capture := range captures {
		if !capture.IsMatch {
			continue
		}
		rStart := capture.Offsets.Start + offset
		if rStart < 0 {
			rStart = 0
		}
		rEnd := capture.Offsets.End + offset
		if rEnd < 0 {
			rEnd = 0
		}
		rng := NewNormalizedRange(rStart, rEnd)

		changes := make([]RuneChange, 0, len(content))
		for _, r := range content {
			changes = append(changes, RuneChange{Rune: r, Change: 1})
		}

		removedRunes := len([]rune(ns.normalized[rng.Start():rng.End()]))

		ns.TransformRange(rng, changes, removedRunes)

		newLen := len(content)
		oldLen := capture.Offsets.End - capture.Offsets.Start
		offset += newLen - oldLen
	}

	return nil
}

// Slice returns a slice of the current NormalizedString. It also returns a flag
// which reports whether the operation is successful or not.
func (ns *NormalizedString) Slice(rng Range) (*NormalizedString, bool) {
	if !ns.rangeIsValid(rng) {
		return nil, false
	}

	normalizedRange, nrOk := ns.CoerceRangeToNormalized(rng)
	if !nrOk {
		return nil, false
	}
	originalRange, orOk := ns.CoerceRangeToOriginal(rng)
	if !orOk {
		return nil, false
	}

	nShift := originalRange.Start()

	original, _ := ns.GetOriginalRange(rng)
	normalized, _ := ns.GetRange(rng)

	alignments := make([]AlignmentRange, normalizedRange.Len())

	for i, alignment := range ns.alignments[normalizedRange.start:normalizedRange.end] {
		alignments[i] = AlignmentRange{
			start: alignment.start - nShift,
			end:   alignment.end - nShift,
		}
	}

	originalShift := ns.originalShift + originalRange.start

	return New(original, normalized, alignments, originalShift), true
}

// Filter applies filtering over the characters of the NormalizedString.
func (ns *NormalizedString) Filter(keep func(rune) bool) {
	removed := 0
	removedStart := 0
	transforms := make([]RuneChange, 0, ns.Len())

	hasLastRune := false
	var lastRune rune

	for _, r := range ns.normalized {
		if !keep(r) {
			removed++
			continue
		}

		if hasLastRune {
			transforms = append(transforms, RuneChange{
				Rune:   lastRune,
				Change: -removed,
			})
		} else {
			removedStart = removed
		}

		lastRune = r
		hasLastRune = true
		removed = 0
	}

	if hasLastRune {
		transforms = append(transforms, RuneChange{
			Rune:   lastRune,
			Change: -removed,
		})
	}
	ns.Transform(transforms, removedStart)
}

// Map maps the characters of the NormalizedString.
func (ns *NormalizedString) Map(mapFunc func(rune) rune) {
	transforms := make([]RuneChange, 0, ns.Len())
	for _, r := range ns.normalized {
		transforms = append(transforms, RuneChange{
			Rune:   mapFunc(r),
			Change: 0,
		})
	}
	ns.Transform(transforms, 0)
}

// ToUpper remaps all Unicode letters of the "normalized" string to their
// upper case.
// FIXME: Go `unicode` package does not consider Unicode special casing
//        (see https://www.unicode.org/Public/UCD/latest/ucd/SpecialCasing.txt)
//        As a result, every single rune is always transformed to a single
//        new rune. This is not the case in the original rust implementation.
//        Should we consider using another "better" package?
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
//        ```go
//        package main
//        import (
// 	          "fmt"
// 	          "unicode"
//        )
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
	ns.Map(unicode.ToUpper)
}

// ToLower remaps all Unicode letters of the "normalized" string to their
// lower case.
// FIXME: see Unicode special casing notes on `NormalizedString.ToUpper()`
func (ns *NormalizedString) ToLower() {
	ns.Map(unicode.ToLower)
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

	lastIndex := lenRunes - trailingSpaces - 1

	transforms := make([]RuneChange, 0, lastIndex-leadingSpaces+1)
	for i := leadingSpaces; i <= lastIndex; i++ {
		change := 0
		if i == lastIndex {
			change = -trailingSpaces
		}
		transforms = append(transforms, RuneChange{
			Rune:   runes[i],
			Change: change,
		})
	}
	ns.Transform(transforms, leadingSpaces)
}

// rangeIsValid validates the given range, to make sure it is on rune boundaries.
func (ns *NormalizedString) rangeIsValid(r Range) bool {
	var s string

	switch r.(type) {
	case OriginalRange:
		s = ns.original
	case NormalizedRange:
		s = ns.normalized
	default:
		panic(fmt.Sprintf("unexpected Range implementation: %#v", r))
	}

	return strutils.IsRuneBoundary(s, r.Start()) &&
		strutils.IsRuneBoundary(s, r.End())
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
