package validate

import (
	"fmt"
	"github.com/ebauman/golicense/pkg/certificate"
	"github.com/spf13/cobra"
	"os"
)

var (
	private = ""
	public  = ""
)

func init() {
	ValidateCertficateCmd.Flags().StringVar(&private, "private", "", "private key file")
	ValidateCertficateCmd.Flags().StringVar(&public, "public", "", "public key file")

	for _, f := range []string{"private", "public"} {
		if err := ValidateCertficateCmd.MarkFlagRequired(f); err != nil {
			panic(err)
		}
	}
}

var ValidateCertficateCmd = &cobra.Command{
	Use:     "certificate",
	Aliases: []string{"cert", "c"},
	Short:   "validate private/public keypair",
	RunE: func(cmd *cobra.Command, args []string) error {
		keyData, err := os.ReadFile(private)
		if err != nil {
			return err
		}

		certData, err := os.ReadFile(public)
		if err != nil {
			return err
		}

		privKey, err := certificate.PEMToPrivateKey(keyData)
		if err != nil {
			return err
		}

		pubKey, err := certificate.PEMToPublicKey(certData)
		if err != nil {
			return err
		}

		if !privKey.PublicKey.Equal(pubKey) {
			return fmt.Errorf("private key and public key do not match")
		}

		return nil
	},
}
