package cmd

import (
	"math"
	"strings"
)

type (
	anyType string
	size    int
)

func (at anyType) isSize(n size) bool {
	return getSize(at) == n
}

func getSize(t anyType) size {
	if strings.Index(string(t), "*") == 0 {
		return 8
	}

	if strings.Contains(string(t), "[]") {
		return 24
	}

	return map[anyType]size{
		"bool":      1,
		"byte":      1,
		"int32":     4,
		"float32":   4,
		"rune":      4,
		"int64":     8,
		"int":       8,
		"float64":   8,
		"uint64":    8,
		"string":    16,
		"error":     16,
		"time.Time": 24,
	}[t]
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
			return float64(int(getSize(k)) * v)
		}
		tot += int(getSize(k)) * v
	}

	return math.Ceil(float64(tot)/8) * 8
}
