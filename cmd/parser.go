package cmd

import (
	"fmt"
	"strings"
)

type structDetails struct {
	countByType map[anyType]int
	namesByType map[anyType][]string
	structName  string
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

		if strings.Contains(line, "//") {
			panic("comments are not supported")
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
