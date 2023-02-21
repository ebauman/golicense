package create

import (
	"crypto/rand"
	"crypto/rsa"
	"fmt"
	"github.com/ebauman/golicense/pkg/certificate"
	"github.com/spf13/cobra"
)

var CreateCertCmd = &cobra.Command{
	Use:     "certificate",
	Aliases: []string{"cert", "c", "crt"},
	RunE: func(cmd *cobra.Command, args []string) error {
		// generate a new private key
		privKey, err := rsa.GenerateKey(rand.Reader, 2048)
		if err != nil {
			return fmt.Errorf("error generating rsa key: %s", err.Error())
		}

		priv, err := certificate.KeyToPrivatePEM(privKey)
		if err != nil {
			return fmt.Errorf("error pem formatting private key: %s", err.Error())
		}

		pub, err := certificate.KeyToPublicPEM(privKey)
		if err != nil {
			return fmt.Errorf("error pem formatting public key: %s", err.Error())
		}

		fmt.Printf("%s%s", string(priv), string(pub))

		return nil
	},
}
