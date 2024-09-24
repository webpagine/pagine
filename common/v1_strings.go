// Copyright 2024 Jelly Terra
// This Source Code Form is subject to the terms of the Mozilla Public License, v. 2.0
// that can be found in the LICENSE file and https://mozilla.org/MPL/2.0/.

package common

import "strings"

type V1Strings struct{}

func (V1Strings) Contains(s, substr string) bool     { return strings.Contains(s, substr) }
func (V1Strings) ContainsAny(s, chars string) bool   { return strings.Contains(s, chars) }
func (V1Strings) ContainsRune(s string, r rune) bool { return strings.ContainsRune(s, r) }

func (V1Strings) Count(s, substr string) int { return strings.Count(s, substr) }

func (V1Strings) HasPrefix(str, prefix string) bool { return strings.HasPrefix(str, prefix) }
func (V1Strings) HasSuffix(str, suffix string) bool { return strings.HasSuffix(str, suffix) }

func (V1Strings) Index(s, substr string) int     { return strings.Index(s, substr) }
func (V1Strings) IndexAny(s, chars string) int   { return strings.IndexAny(s, chars) }
func (V1Strings) IndexByte(s string, c byte) int { return strings.IndexByte(s, c) }
func (V1Strings) IndexRune(s string, r rune) int { return strings.IndexRune(s, r) }

func (V1Strings) Join(sep string, elems ...string) string { return strings.Join(elems, sep) }

func (V1Strings) LastIndex(s, substr string) int     { return strings.LastIndex(s, substr) }
func (V1Strings) LastIndexAny(s, chars string) int   { return strings.LastIndexAny(s, chars) }
func (V1Strings) LastIndexByte(s string, c byte) int { return strings.LastIndexByte(s, c) }

func (V1Strings) Repeat(s string, count int) string { return strings.Repeat(s, count) }

func (V1Strings) Replace(s, old, new string, n int) string { return strings.Replace(s, old, new, n) }
func (V1Strings) ReplaceAll(s, old, new string) string     { return strings.ReplaceAll(s, old, new) }

func (V1Strings) Split(s, substr string) []string      { return strings.Split(s, substr) }
func (V1Strings) SplitAfter(s, substr string) []string { return strings.SplitAfter(s, substr) }
func (V1Strings) SplitAfterN(s, substr string, n int) []string {
	return strings.SplitAfterN(s, substr, n)
}
func (V1Strings) SplitN(s, substr string, n int) []string { return strings.SplitN(s, substr, n) }

func (V1Strings) ToLower(s string) string { return strings.ToLower(s) }
func (V1Strings) ToTitle(s string) string { return strings.ToTitle(s) }
func (V1Strings) ToUpper(s string) string { return strings.ToUpper(s) }
func (V1Strings) ToValidUTF8(s, replacement string) string {
	return strings.ToValidUTF8(s, replacement)
}

func (V1Strings) Trim(s, cutset string) string         { return strings.Trim(s, cutset) }
func (V1Strings) TrimLeft(s, cutset string) string     { return strings.TrimLeft(s, cutset) }
func (V1Strings) TrimPrefix(str, prefix string) string { return strings.TrimPrefix(str, prefix) }
func (V1Strings) TrimRight(s, cutset string) string    { return strings.TrimRight(s, cutset) }
func (V1Strings) TrimSpace(s string) string            { return strings.TrimSpace(s) }
func (V1Strings) TrimSuffix(str, suffix string) string { return strings.TrimSuffix(str, suffix) }
