package cmd

import "github.com/spf13/cobra"

type commandArgs struct {
	inputFile string
}

func newArgs(cmd *cobra.Command) (*commandArgs, error) {
	var args commandArgs

	cmd.Flags().StringVarP(&args.inputFile, "input-file", "i", "", "The path to the file containing ONLY your go struct")

	if err := cmd.MarkFlagRequired("input-file"); err != nil {
		return nil, err
	}

	return &args, nil
}
