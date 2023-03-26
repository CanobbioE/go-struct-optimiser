package cmd

import (
	"github.com/spf13/cobra"
	"io/ioutil"
	"os"
)

func Command() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "go-struct-optimiser",
		Short: "Optimise your Golang structure's memory allocation",
		Run: func(cmd *cobra.Command, args []string) {
			if err := cmd.Usage(); err != nil {
				panic(err)
			}
		},
	}

	args, err := newArgs(cmd)
	if err != nil {
		panic(err)
	}

	cmd.Run = func(*cobra.Command, []string) {
		file, err := os.Open(args.inputFile)
		if err != nil {
			panic(err)
		}

		content, err := ioutil.ReadAll(file)
		if err != nil {
			panic(err)
		}

		details := parseStruct(string(content))
		printOptimisedStruct(details)
	}
	return cmd
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
	tmp := make([]anyType, 0)
	for typeName, typeOccurrences := range types {
		// skip types that are not 4bytes
		if !typeName.isSize(4) {
			continue
		}

		// Since a full block is 8 bytes, and we are working with 4 bytes types,
		// create a defined number (N) of 2-sized arrays,
		// where N is defined as half the occurrences of the current 4 bytes type.
		for i := 0; i < typeOccurrences/2; i++ {
			cpuPasses = append(cpuPasses, nSizedArray(typeName, 2))
		}

		if typeOccurrences%2 == 0 {
			continue
		}

		// if count / 2 has a reminder, combine with the next 4bytes type
		switch len(tmp) {
		case 0, 1:
			tmp = append(tmp, typeName)
		case 2:
			cpuPasses = append(cpuPasses, tmp)
			tmp = make([]anyType, 0)
		}
	}

	// Calculate the total number of 4bytes blocks,
	// if the result is an odd number: notify the caller that we have a remainder (i.e. extra free space)
	fourBytesBlocks := types["int32"] + types["float32"] + types["rune"]
	if fourBytesBlocks%2 == 0 {
		return cpuPasses, nil
	}
	return cpuPasses, tmp
}

func handleOneByteBlocks(cpuPasses [][]anyType, types map[anyType]int, remainder []anyType) [][]anyType {
	// calculate the total number of 1 byte blocks already parsed
	oneByteBlocks := types["bool"] + types["byte"]

	// collect all the unused 1 byte blocks in a map
	unusedBlocks := make(map[anyType]int)
	for typeName, typeOccurrences := range types {
		if typeName.isSize(1) {
			unusedBlocks[typeName] = typeOccurrences
		}
	}

	var usedBlocks int
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
