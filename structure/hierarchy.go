package structure

import (
	"github.com/webpagine/pagine/v2/collection"
	"github.com/webpagine/pagine/v2/util"
	"github.com/webpagine/pagine/v2/vfs"
	"maps"
	"os"
	"sync"
)

type Level struct {
	Root vfs.DirFS

	Data MetadataSet

	Units []Unit

	Levels  []Level
	Reports []LevelReport
}

type LevelReport struct {
	Level *Level
	Err   error
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

func ExecuteLevels(env *Env, root, dest vfs.DirFS, inherit MetadataSet) (Level, error) {

	var (
		wg sync.WaitGroup

		unitManifest UnitManifest
		metadata     MetadataSet

		units   collection.SyncArray[Unit]
		reports collection.SyncArray[LevelReport]
		levels  collection.SyncArray[Level]
	)

	err := func() error {

		err := util.UnmarshalTOMLFile(root, "/metadata.toml", &metadata)
		switch {
		case err == nil:
			metadata = metadata.Inherit(inherit)
		case os.IsNotExist(err):
			metadata = inherit
		default:
			if err != nil {
				return err
			}
		}

		err = util.UnmarshalTOMLFile(root, "/unit.toml", &unitManifest)
		switch {
		case err == nil:
			// No error will cause interrupt below.

			for _, unitItem := range unitManifest.Units {
				wg.Add(1)
				go func() {
					defer wg.Done()

					unit := Unit{
						Output:   unitItem.Output,
						Template: unitItem.Template,
					}

					unit.Report.TemplateErrors, unit.Report.Error = unit.Generate(env, root, dest, metadata, unitItem.Define)
					units.Append(unit)
				}()
			}
		case os.IsNotExist(err):
		default:
			return err
		}

		entries, err := root.ReadDir("/")
		if err != nil {
			return err
		}

		// No error will cause interrupt below.

		for _, entry := range entries {
			if !entry.IsDir() {
				continue
			}

			sub := root.DirFS(entry.Name())

			wg.Add(1)
			go func() {
				defer wg.Done()

				level, err := ExecuteLevels(env, sub, dest, metadata)
				if err != nil {
					reports.Append(LevelReport{
						Level: &level,
						Err:   err,
					})
					return
				}
				levels.Append(level)
			}()
		}

		return nil
	}()
	if err != nil {
		return Level{}, err
	}

	wg.Wait()

	return Level{
		Root:    root,
		Data:    metadata,
		Units:   units.RawArray,
		Levels:  levels.RawArray,
		Reports: reports.RawArray,
	}, nil
}
