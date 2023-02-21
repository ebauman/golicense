package validate

import "github.com/spf13/cobra"

var ValidateCmd = &cobra.Command{
	Use:   "validate",
	Short: "validate resources",
}

func init() {
	ValidateCmd.AddCommand(ValidateCertficateCmd)
	ValidateCmd.AddCommand(ValidateLicenseCmd)
}
