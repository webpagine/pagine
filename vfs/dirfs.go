package vfs

import (
	"io/fs"
	"os"
	"path/filepath"
)

type PathError = fs.PathError

func OsDirFS(path string) DirFS {
	return DirFS{Path: path}
}

type DirFS struct {
	Path string
}

func (dir DirFS) DirFS(path string) DirFS {
	return DirFS{Path: filepath.Join(dir.Path, path)}
}

func (dir DirFS) Open(name string) (fs.File, error) {
	return os.Open(filepath.Join(dir.Path, name))
}

func (dir DirFS) ReadFile(name string) ([]byte, error) {
	return os.ReadFile(filepath.Join(dir.Path, name))
}

func (dir DirFS) ReadDir(name string) ([]fs.DirEntry, error) {
	return os.ReadDir(filepath.Join(dir.Path, name))
}

func (dir DirFS) Stat(name string) (fs.FileInfo, error) {
	return os.Stat(filepath.Join(dir.Path, name))
}

func (dir DirFS) CreateFile(name string) (*os.File, error) {
	path := filepath.Join(dir.Path, name)

	parentDir := filepath.Dir(path)
	_, err := os.Stat(parentDir)
	switch {
	case err == nil:
	case os.IsNotExist(err):
		err = os.MkdirAll(parentDir, 0755)
		if err != nil {
			return nil, err
		}
	default:
		return nil, err
	}

	return os.Create(path)
}
