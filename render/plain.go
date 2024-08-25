// Copyright 2024 Jelly Terra
// This Source Code Form is subject to the terms of the Mozilla Public License, v. 2.0
// that can be found in the LICENSE file and https://mozilla.org/MPL/2.0/.

package render

func Plain(content []byte) (string, error) {
	return string(content), nil
}

func init() {
	Renderers["text/plain"] = Plain
}
