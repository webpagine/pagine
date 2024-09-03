// Copyright 2024 Jelly Terra
// This Source Code Form is subject to the terms of the Mozilla Public License, v. 2.0
// that can be found in the LICENSE file and https://mozilla.org/MPL/2.0/.

package structure

import . "strings"

type v1Strings struct{}

func (v1Strings) Contains(s, substr string) bool     { return Contains(s, substr) }
func (v1Strings) ContainsAny(s, chars string) bool   { return Contains(s, chars) }
func (v1Strings) ContainsRune(s string, r rune) bool { return ContainsRune(s, r) }

func (v1Strings) Count(s, substr string) int { return Count(s, substr) }

func (v1Strings) HasPrefix(str, prefix string) bool { return HasPrefix(str, prefix) }
func (v1Strings) HasSuffix(str, suffix string) bool { return HasSuffix(str, suffix) }

func (v1Strings) Index(s, substr string) int     { return Index(s, substr) }
func (v1Strings) IndexAny(s, chars string) int   { return IndexAny(s, chars) }
func (v1Strings) IndexByte(s string, c byte) int { return IndexByte(s, c) }
func (v1Strings) IndexRune(s string, r rune) int { return IndexRune(s, r) }

func (v1Strings) Join(sep string, elems ...string) string { return Join(elems, sep) }

func (v1Strings) LastIndex(s, substr string) int     { return LastIndex(s, substr) }
func (v1Strings) LastIndexAny(s, chars string) int   { return LastIndexAny(s, chars) }
func (v1Strings) LastIndexByte(s string, c byte) int { return LastIndexByte(s, c) }

func (v1Strings) Repeat(s string, count int) string { return Repeat(s, count) }

func (v1Strings) Replace(s, old, new string, n int) string { return Replace(s, old, new, n) }
func (v1Strings) ReplaceAll(s, old, new string) string     { return ReplaceAll(s, old, new) }

func (v1Strings) Split(s, substr string) []string      { return Split(s, substr) }
func (v1Strings) SplitAfter(s, substr string) []string { return SplitAfter(s, substr) }
func (v1Strings) SplitAfterN(s, substr string, n int) []string {
	return SplitAfterN(s, substr, n)
}
func (v1Strings) SplitN(s, substr string, n int) []string { return SplitN(s, substr, n) }

func (v1Strings) ToLower(s string) string { return ToLower(s) }
func (v1Strings) ToTitle(s string) string { return ToTitle(s) }
func (v1Strings) ToUpper(s string) string { return ToUpper(s) }
func (v1Strings) ToValidUTF8(s, replacement string) string {
	return ToValidUTF8(s, replacement)
}

func (v1Strings) Trim(s, cutset string) string         { return Trim(s, cutset) }
func (v1Strings) TrimLeft(s, cutset string) string     { return TrimLeft(s, cutset) }
func (v1Strings) TrimPrefix(str, prefix string) string { return TrimPrefix(str, prefix) }
func (v1Strings) TrimRight(s, cutset string) string    { return TrimRight(s, cutset) }
func (v1Strings) TrimSpace(s string) string            { return TrimSpace(s) }
func (v1Strings) TrimSuffix(str, suffix string) string { return TrimSuffix(str, suffix) }
