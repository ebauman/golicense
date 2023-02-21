package create

import "github.com/spf13/cobra"

var CreateCmd = &cobra.Command{
	Use:   "create",
	Short: "create resources",
}

func init() {
	CreateCmd.AddCommand(CreateCertCmd)
	CreateCmd.AddCommand(CreateLicenseCmd)
}
