package main

import "github.com/CanobbioE/go-struct-optimiser/cmd"

func main() {
	if err := cmd.Command().Execute(); err != nil {
		panic(err)
	}
}
