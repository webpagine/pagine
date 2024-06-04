// Copyright 2024 Jelly Terra
// This Source Code Form is subject to the terms of the Mozilla Public License, v. 2.0
// that can be found in the LICENSE file and https://mozilla.org/MPL/2.0/.

package collection

import "sync"

func ValuesOfRawMap[K comparable, V any](m map[K]V) []V {
	values := make([]V, len(m))
	i := 0
	for _, v := range m {
		values[i] = v
		i++
	}
	return values
}

func MergeRawMap[K comparable, V any](set map[K]V, subset map[K]V) {
	for k, v := range subset {
		set[k] = v
	}
}

type Map[K comparable, V any] struct {
	RawMap map[K]V
}

func NewMap[K comparable, V any]() Map[K, V] {
	return Map[K, V]{
		RawMap: map[K]V{},
	}
}

type SyncMap[K comparable, V any] struct {
	Map[K, V]
	lock sync.RWMutex
}

func (s *SyncMap[K, V]) Set(k K, v V) {
	s.lock.Lock()
	defer s.lock.Unlock()

	if s.RawMap == nil {
		s.RawMap = map[K]V{}
	}

	s.RawMap[k] = v
}

func (s *SyncMap[K, V]) Get(k K) (V, bool) {
	s.lock.RLock()
	defer s.lock.RUnlock()
	v, ok := s.RawMap[k]
	return v, ok
}
