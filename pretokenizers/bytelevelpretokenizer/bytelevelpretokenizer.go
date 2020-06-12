// Copyright (c) 2020, NLP Odyssey Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package bytelevelpretokenizer

import (
	"github.com/dlclark/regexp2"
	"github.com/nlpodyssey/gotokenizers/normalizers/normalizedstring"
	"github.com/nlpodyssey/gotokenizers/pretokenizers"
	"strings"
	"unicode"
)

// ByteLevelPreTokenizer allows the generation of pre-tokens suitable in the
// context of BPE (Byte Pair Encoding) models.
//
// This pre-tokenizer splits the input string into tokens according to the
// given regular expression. It also transforms the string expanding each
// UTF-8 character (rune) into bytes, which are then re-mapped to new custom
// runes.
//
// A whitespace prefix (' ') can be optionally prepended to the input string,
// unless the first rune of the string is already a unicode whitespace.
//
// Offsets trimming can be enabled to exclude whitespaces in the post-processing
// step.
type ByteLevelPreTokenizer struct {
	splittingRegexp    *regexp2.Regexp
	prefixSpaceEnabled bool
}

var _ pretokenizers.PreTokenizer = &ByteLevelPreTokenizer{}

// DefaultSplittingRegexp is a simple default regular expression that
// can be used for ByteLevelPreTokenizer.
//
// This MUST be treated as a read-only constant value .
var DefaultSplittingRegexp = regexp2.MustCompile(
	`'s|'t|'re|'ve|'m|'ll|'d| ?\p{L}+| ?\p{N}+| ?[^\s\p{L}\p{N}]+|\s+(?!\S)|\s+`,
	regexp2.IgnoreCase|regexp2.Multiline)

// NewByteLevelPreTokenizer returns a new ByteLevelPreTokenizer, setting
// the splitting regular expression and the flag to enable or disable prefix
// space insertion.
func NewByteLevelPreTokenizer(
	splittingRegexp *regexp2.Regexp,
	prefixSpaceEnabled bool,
) *ByteLevelPreTokenizer {
	return &ByteLevelPreTokenizer{
		splittingRegexp:    splittingRegexp,
		prefixSpaceEnabled: prefixSpaceEnabled,
	}
}

// DefaultByteLevelPreTokenizer returns a new ByteLevelPreTokenizer, using
// DefaultSplittingRegexp, and enabling prefix space insertion.
func DefaultByteLevelPreTokenizer() *ByteLevelPreTokenizer {
	return NewByteLevelPreTokenizer(DefaultSplittingRegexp, true)
}

// PreTokenize transforms the NormalizedString, remapping its bytes,
// and splits it according to the configured regular expression.
//
// If whitespace prefix is enabled, a whitespace (' ') is prepended to
// the NormalizedString, actually modifying its "normalized" value, only if
// the first rune of the string is not already a unicode whitespace.
func (b *ByteLevelPreTokenizer) PreTokenize(
	ns *normalizedstring.NormalizedString,
) ([]pretokenizers.PreToken, error) {
	if b.prefixSpaceEnabled && !startsWithWhitespace(ns.Get()) {
		ns.Prepend(" ")
	}

	str := ns.Get()
	runes := []rune(str)

	tokens, err := b.makeInitialTokens(str)
	if err != nil {
		return nil, err
	}

	changes := make([]normalizedstring.RuneChanges, 0, len(str))

	lastIndex := 0
	for index, token := range tokens {
		originalString := string(runes[token.Start:token.End])
		tokenChanges, newString := makeTokenChangesAndNewString(originalString)
		changes = append(changes, tokenChanges...)
		// The length of the token's changes corresponds to the
		// length of the new string in runes.
		end := lastIndex + len(tokenChanges)
		tokens[index] = pretokenizers.PreToken{
			String: newString,
			Start:  lastIndex,
			End:    end,
		}
		lastIndex = end
	}

	ns.Transform(changes, 0)

	return tokens, nil
}

// makeInitialTokens splits the input strings according to the configured
// regular expression. Start and End rune indices are set for each token,
// while String is not set.
func (b *ByteLevelPreTokenizer) makeInitialTokens(s string) ([]pretokenizers.PreToken, error) {
	match, err := b.splittingRegexp.FindStringMatch(s)
	if err != nil {
		return nil, err
	}
	tokens := make([]pretokenizers.PreToken, 0)

	lastIndex := 0
	for match != nil {
		if match.Index != lastIndex {
			tokens = append(tokens, pretokenizers.PreToken{
				Start: lastIndex,
				End:   match.Index,
			})
		}

		lastIndex = match.Index + match.Length
		tokens = append(tokens, pretokenizers.PreToken{
			Start: match.Index,
			End:   lastIndex,
		})

		match, err = b.splittingRegexp.FindNextMatch(match)
		if err != nil {
			return nil, err
		}
	}

	endIndex := len([]rune(s))
	if lastIndex != endIndex {
		tokens = append(tokens, pretokenizers.PreToken{
			Start: lastIndex,
			End:   endIndex,
		})
	}

	return tokens, nil
}

// makeTokenChangesAndNewString creates the sequence of changes for the input
// token string, expanding each rune into re-mapped bytes, and also returns
// the new transformed string.
func makeTokenChangesAndNewString(s string) ([]normalizedstring.RuneChanges, string) {
	var sb strings.Builder
	sb.Grow(len(s))

	runeChanges := make([]normalizedstring.RuneChanges, 0, len(s))
	for _, r := range s {
		for i, b := range []byte(string(r)) {
			newRune := byteToRune[b]
			changes := 0
			if i != 0 {
				changes = 1
			}
			runeChanges = append(runeChanges, normalizedstring.RuneChanges{
				Rune:    newRune,
				Changes: changes,
			})
			sb.WriteRune(newRune)
		}
	}
	return runeChanges, sb.String()
}

func startsWithWhitespace(s string) bool {
	return len(s) != 0 && unicode.In([]rune(s)[0], unicode.White_Space)
}

var byteToRune [0x100]rune

func init() {
	n := 0
	for i := range byteToRune {
		if (i >= '!' && i <= '~') || (i >= 0xA1 && i <= 0xAC) || (i >= 0xAE && i <= 0xFF) {
			byteToRune[i] = rune(i)
		} else {
			byteToRune[i] = rune(0x100 + n)
			n++
		}
	}
}
