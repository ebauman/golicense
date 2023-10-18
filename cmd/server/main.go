package main

import (
	"fmt"
	"github.com/ebauman/golicense/pkg/server"
	"github.com/spf13/cobra"
	"os"
)

var (
	dsn       string
	httpPort  int
	httpsPort int
)

func init() {
	rootCmd.Flags().StringVar(&dsn, "dsn", "", "mysql dsn")
	rootCmd.Flags().IntVar(&httpPort, "http-port", 8080, "http port")
	rootCmd.Flags().IntVar(&httpsPort, "https-port", 8443, "https port")
}

var rootCmd = &cobra.Command{
	Use:   "golicense-server",
	Short: "manage golicense storage",
	RunE:  run,
}

func run(cmd *cobra.Command, args []string) error {
	svr, err := server.New(httpPort, httpsPort, dsn)
	if err != nil {
		return err
	}

	if err := svr.Run(cmd.Context()); err != nil {
		return err
	}

	<-cmd.Context().Done()

	return cmd.Context().Err()
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
