package main

import (
	"fmt"
	"math"
	"strings"
)

const input = `
type SpecialStruct struct {
	a,b,c int32
	ptr *error
	f float64
	flag,flag2 bool
}
`

func main() {
	details := parseStruct(input)
	printOptimisedStruct(details)
}

type (
	anyType       string
	size          int
	structDetails struct {
		countByType map[anyType]int
		namesByType map[anyType][]string
		structName  string
	}
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

func printOptimisedStruct(details *structDetails) {
	fmt.Println("the best size for your struct is", bestSize(details.countByType))

	cpuPasses := optimize(details.countByType)

	indexByType := make(map[anyType]int)
	for t := range details.countByType {
		indexByType[t] = 0
	}

	fmt.Printf("type %s struct {\n", details.structName)
	for i := range cpuPasses {
		for j := range cpuPasses[i] {
			t := cpuPasses[i][j]

			varName := details.namesByType[t][indexByType[t]]
			indexByType[t]++

			fmt.Printf("\t%s\t%v\n", varName, t)
		}
	}
	fmt.Println("}")
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

func optimize(types map[anyType]int) [][]anyType {

	cpuPasses := handleFullBlocks(types)

	cpuPasses, remainder := handleFourBytesBlocks(cpuPasses, types)

	cpuPasses = handleOneByteBlocks(cpuPasses, types, remainder)

	return cpuPasses
}

func handleFullBlocks(types map[anyType]int) [][]anyType {
	passes := make([][]anyType, 0)
	for k, v := range types {
		if k.isSize(8) || k.isSize(16) {
			passes = append(passes, nSizedArray(k, v))
		}
	}
	return passes
}

func handleFourBytesBlocks(cpuPasses [][]anyType, types map[anyType]int) ([][]anyType, []anyType) {
	fourBytesBlocks := types["int32"] + types["float32"] + types["rune"]

	tmp := make([]anyType, 0)
	for k, v := range types {
		// skip types that are not 4bytes
		if !k.isSize(4) {
			continue
		}

		// create N 2-sized arrays for each 4bytes type, where N = count/2
		for i := 0; i < v/2; i++ {
			cpuPasses = append(cpuPasses, nSizedArray(k, 2))
		}

		if v%2 == 0 {
			continue
		}

		// if count / 2 has a reminder, combine with the next 4bytes type
		switch len(tmp) {
		case 0, 1:
			tmp = append(tmp, k)
		case 2:
			cpuPasses = append(cpuPasses, tmp)
			tmp = make([]anyType, 0)
		}
	}

	if fourBytesBlocks%2 == 0 {
		return cpuPasses, nil
	}

	// if we have an odd number of 4bytes blocks, notify the caller that we have a remainder
	return cpuPasses, tmp
}

func handleOneByteBlocks(cpuPasses [][]anyType, types map[anyType]int, remainder []anyType) [][]anyType {
	oneByteBlocks := types["bool"] + types["byte"]

	var usedBlocks int
	unusedBlocks := make(map[anyType]int)
	for k, v := range types {
		if k.isSize(1) {
			unusedBlocks[k] = v
		}
	}

	if remainder != nil {
		for k, v := range types {
			if usedBlocks == 4 || usedBlocks >= oneByteBlocks {
				break
			}

			if !k.isSize(1) {
				continue
			}

			remainder = append(remainder, nSizedArray(k, v)...)
			usedBlocks++
			unusedBlocks[k] = 0
		}
		cpuPasses = append(cpuPasses, remainder)
	}

	if usedBlocks >= oneByteBlocks {
		return cpuPasses
	}

	tmp := make([]anyType, 0)
	for k, v := range unusedBlocks {
		tmp = append(tmp, nSizedArray(k, v)...)
	}

	return append(cpuPasses, tmp)

}

func nSizedArray(t anyType, count int) []anyType {
	passes := make([]anyType, count)
	for i := range passes {
		passes[i] = t
	}
	return passes
}

func parseStruct(s string) *structDetails {
	details := &structDetails{
		countByType: make(map[anyType]int),
		namesByType: make(map[anyType][]string),
	}

	lines := strings.Split(s, "\n")

	for _, line := range lines {
		if line == "" {
			continue
		}
		if strings.Contains(line, "type") {
			details.structName = strings.Split(line, " ")[1]
			continue
		}

		if strings.Contains(line, "}") {
			return details
		}

		varNames := strings.Split(line, " ")[0]
		varNames = strings.TrimSpace(varNames)

		count := len(strings.Split(varNames, ","))

		t := anyType(strings.Split(line, " ")[1])

		if v, ok := details.countByType[t]; ok {
			count += v
		}
		details.countByType[t] = count
		details.namesByType[t] = strings.Split(varNames, ",")
	}

	return details
}
