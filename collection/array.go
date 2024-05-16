// Copyright 2024 Jelly Terra
// This Source Code Form is subject to the terms of the Mozilla Public License, v. 2.0
// that can be found in the LICENSE file and https://mozilla.org/MPL/2.0/.

package collection

import "sync"

type Array[E any] struct {
	RawArray []E
}

func (a *Array[E]) Len() uint { return uint(len(a.RawArray)) }

func (a *Array[E]) Cap() uint { return uint(len(a.RawArray)) }

func (a *Array[E]) Append(e E) { a.RawArray = append(a.RawArray, e) }

type ComparableArray[E comparable] struct {
	Array[E]
}

func (a *ComparableArray[E]) Contains(e E) bool {
	for _, each := range a.RawArray {
		if each == e {
			return true
		}
	}

	return false
}

type SyncArray[E any] struct {
	Array[E]

	lock sync.RWMutex
}

func (a *SyncArray[E]) Len() uint {
	a.lock.RLock()
	defer a.lock.RUnlock()
	return uint(len(a.RawArray))
}

func (a *SyncArray[E]) Cap() uint {
	a.lock.RLock()
	defer a.lock.RUnlock()
	return uint(cap(a.RawArray))
}

func (a *SyncArray[E]) Append(e E) {
	a.lock.Lock()
	defer a.lock.Unlock()
	a.RawArray = append(a.RawArray, e)
}

type ComparableSyncArray[E comparable] struct {
	SyncArray[E]
}

func (a *ComparableSyncArray[E]) Contains(e E) bool {
	a.lock.RLock()
	defer a.lock.RUnlock()

	for _, each := range a.RawArray {
		if each == e {
			return true
		}
	}

	return false
}
