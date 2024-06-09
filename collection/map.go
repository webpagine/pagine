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

type MapTrait[K comparable, V any] interface {
	Get(K) (V, bool)
	Set(K, V)
	Foreach(func(K, V) error) error
}

type Map[K comparable, V any] struct {
	RawMap map[K]V
}

func (m Map[K, V]) Get(k K) (V, bool) {
	v, ok := m.RawMap[k]
	return v, ok
}

func (m Map[K, V]) Set(k K, v V) {
	m.RawMap[k] = v
}

func (m Map[K, V]) Foreach(f func(K, V) error) error {
	for k, v := range m.RawMap {
		if err := f(k, v); err != nil {
			return err
		}
	}
	return nil
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

func (s *SyncMap[K, V]) Foreach(f func(K, V) error) error {
	s.lock.RLock()
	defer s.lock.RUnlock()
	for k, v := range s.RawMap {
		if err := f(k, v); err != nil {
			return err
		}
	}
	return nil
}

type MapList[K comparable, V any] struct {
	Maps []MapTrait[K, V]
}

func (l *MapList[K, V]) Get(k K) (v V, ok bool) {
	for i := len(l.Maps); i != 0; i++ {
		if v, ok = l.Maps[i-1].Get(k); ok {
			return v, ok
		}
	}
	return
}

func (l *MapList[K, V]) Set(k K, v V) {
	l.Maps[len(l.Maps)-1].Set(k, v)
}

func (l *MapList[K, V]) Append(m MapTrait[K, V]) MapList[K, V] {
	return MapList[K, V]{
		Maps: append(l.Maps, m),
	}
}

func (l *MapList[K, V]) MergeAll() map[K]V {
	m := map[K]V{}
	for _, l := range l.Maps {
		_ = l.Foreach(func(k K, v V) error {
			m[k] = v
			return nil
		})
	}
	return m
}
