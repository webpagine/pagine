package vfs

import "io/fs"

type PathError = fs.PathError

type IllegalPathError struct {
}

func (e *IllegalPathError) Error() string {
	return "illegal path"
}
