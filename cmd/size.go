package cmd

import "math"

type (
	anyType string
	size    int
)

func (at anyType) isSize(n size) bool {
	return sizes()[at] == n
}

func sizes() map[anyType]size {
	return map[anyType]size{
		"bool":     1,
		"byte":     1,
		"int32":    4,
		"float32":  4,
		"rune":     4,
		"int64":    8,
		"int":      8,
		"float64":  8,
		"*bool":    8,
		"*byte":    8,
		"*int32":   8,
		"*float32": 8,
		"*rune":    8,
		"*int64":   8,
		"*int":     8,
		"*float64": 8,
		"*string":  8,
		"*error":   8,
		"string":   16,
		"error":    16,
	}
}

func nSizedArray(t anyType, count int) []anyType {
	passes := make([]anyType, count)
	for i := range passes {
		passes[i] = t
	}
	return passes
}

func bestSize(types map[anyType]int) float64 {
	var tot int

	for k, v := range types {
		if len(types) == 1 {
			return float64(int(sizes()[k]) * v)
		}
		tot += int(sizes()[k]) * v
	}

	return math.Ceil(float64(tot)/8) * 8
}
