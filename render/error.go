// Copyright 2024 Jelly Terra
// This Source Code Form is subject to the terms of the Mozilla Public License, v. 2.0
// that can be found in the LICENSE file and https://mozilla.org/MPL/2.0/.

package render

type UnknownMimeTypeError struct {
	MimeType string
}

func (e *UnknownMimeTypeError) Error() string {
	return "unknown mime type: " + e.MimeType
}

type UnknownExtError struct {
	Mime, Ext string
}

func (e *UnknownExtError) Error() string {
	return "unknown ext: " + e.Ext
}

type NoRenderFoundError struct {
	MimeType string
}

func (e *NoRenderFoundError) Error() string {
	return "no render found for mime type: " + e.MimeType
}
