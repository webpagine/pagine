package structure

func add(aInt, bInt any) int { return aInt.(int) + bInt.(int) }

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

var builtinFuncMap = map[string]any{
	"add":            add,
	"divideSliceByN": divideSliceByN,
	"mapAsSlice":     mapAsSlice,
}
