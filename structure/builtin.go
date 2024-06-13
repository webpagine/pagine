// Copyright 2024 Jelly Terra
// This Source Code Form is subject to the terms of the Mozilla Public License, v. 2.0
// that can be found in the LICENSE file and https://mozilla.org/MPL/2.0/.

package structure

func add(aInt, bInt any) int { return aInt.(int) + bInt.(int) }
func sub(aInt, bInt any) int { return aInt.(int) - bInt.(int) }
func mul(aInt, bInt any) int { return aInt.(int) * bInt.(int) }
func div(aInt, bInt any) int { return aInt.(int) / bInt.(int) }
func mod(aInt, bInt any) int { return aInt.(int) % bInt.(int) }

func divideSliceByN(s []any, nInt any) [][]any {
	n := nInt.(int)

	var divided [][]any
	for begin := 0; begin < len(s); begin += n {
		end := begin + n

		if end > len(s) {
			end = len(s)
		}

		divided = append(divided, s[begin:end])
	}
	return divided
}

func mapAsSlice(m map[string]any, keyName, valueName string) []any {
	slice := make([]any, len(m))
	i := 0
	for k, v := range m {
		slice[i] = map[string]any{
			keyName:   k,
			valueName: v,
		}
		i++
	}
	return slice
}
