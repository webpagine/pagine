// Copyright 2024 Jelly Terra
// This Source Code Form is subject to the terms of the Mozilla Public License, v. 2.0
// that can be found in the LICENSE file and https://mozilla.org/MPL/2.0/.

package common

import "errors"

type Slice struct {
	Raw []any
}

func (s *Slice) Len() int                    { return len(s.Raw) }
func (s *Slice) Index(i int) any             { return s.Raw[i] }
func (s *Slice) Slice(start, end int) *Slice { return &Slice{s.Raw[start:end]} }
func (s *Slice) Push(e ...any) (_ Void)      { s.Raw = append(s.Raw, e...); return }
func (s *Slice) StringSlice() []string {
	slice := make([]string, len(s.Raw))
	for i, e := range s.Raw {
		slice[i], _ = e.(string)
	}
	return slice
}

type Map struct {
	Raw map[string]any
}

type Void string

func (m *Map) Set(key string, value any) (_ Void) { m.Raw[key] = value; return }
func (m *Map) Get(key string) any                 { return m.Raw[key] }
func (m *Map) Has(key string) bool                { return m.Raw[key] != nil }
func (m *Map) Delete(key string) (_ Void)         { delete(m.Raw, key); return }

var V1Builtin = map[string]any{
	"panic": func(v string) (_ Void) { panic(errors.New(v)); return },

	"Any": func(v any) any { return v },
	"Int": func(v any) int { return v.(int) },
	"Str": func(v any) string { return v.(string) },

	"makeSlice": func() *Slice { return new(Slice) },
}
