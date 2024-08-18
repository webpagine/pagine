// Copyright 2024 Jelly Terra
// This Source Code Form is subject to the terms of the Mozilla Public License, v. 2.0
// that can be found in the LICENSE file and https://mozilla.org/MPL/2.0/.

package structure

import (
	"github.com/jellyterra/collection-go"
	"github.com/webpagine/pagine/v2/vfs"
	"maps"
	"sync"
)

type Level struct {
	Root *vfs.DirFS

	Data MetadataSet

	Units []*Unit

	Levels []*Level
}

type LevelReport struct {
	Level *Level

	Error       error
	UnitReports []*UnitReport
}

func (l *Level) Generate(env *Env, dest *vfs.DirFS) (*LevelReport, error) {
	var (
		unitReports collection.SyncVector[*UnitReport]

		wg sync.WaitGroup
	)

	for _, unit := range l.Units {
		wg.Add(1)
		go func() {
			defer wg.Done()

			errors, err := unit.Generate(env, l.Root, dest, l.Data)

			unitReports.Push(&UnitReport{
				Unit:           unit,
				Error:          err,
				TemplateErrors: errors,
			})
		}()
	}
	wg.Wait()

	return &LevelReport{
		Level:       l,
		UnitReports: unitReports.It.Raw,
	}, nil
}

type MetadataSet map[string]map[string]any // map[namespace]map[dataKey]dataValue

func (m MetadataSet) Clone() MetadataSet {
	cloned := maps.Clone(m)

	for namespace, dataMap := range m {
		cloned[namespace] = maps.Clone(dataMap)
	}

	return cloned
}

func (m MetadataSet) Merge(new MetadataSet) {
	for namespace, newDataMap := range new {
		originalDataMap, ok := m[namespace]
		if !ok {
			originalDataMap = map[string]any{}
		}
		for k, v := range newDataMap {
			originalDataMap[k] = v
		}

		m[namespace] = originalDataMap
	}
}

func (m MetadataSet) Inherit(old MetadataSet) MetadataSet {
	cloned := old.Clone()
	cloned.Merge(m)
	return cloned
}
