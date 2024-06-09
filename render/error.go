package render

import "fmt"

type NoRendererFoundError struct {
	Path string
}

func (e *NoRendererFoundError) Error() string {
	return fmt.Sprint("no renderer found for the content type: ", e.Path)
}
