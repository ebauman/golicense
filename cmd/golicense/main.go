package main

import (
	"fmt"
	"github.com/ebauman/golicense/cmd/golicense/create"
	"github.com/ebauman/golicense/cmd/golicense/validate"
	"github.com/spf13/cobra"
	"os"
)

var rootCmd = &cobra.Command{
	Use:   "golicense",
	Short: "manage licenses for go programs",
}

func main() {
	rootCmd.AddCommand(create.CreateCmd)
	rootCmd.AddCommand(validate.ValidateCmd)

	if err := rootCmd.Execute(); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
