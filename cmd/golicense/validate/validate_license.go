package validate

import (
	"crypto/rsa"
	"encoding/json"
	"fmt"
	"github.com/ebauman/golicense/pkg/certificate"
	"github.com/ebauman/golicense/pkg/license"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
	"os"
)

var (
	licenseKeyFile = ""
	privateKeyFile = ""
	publicKeyFile  = ""

	pubKey     *rsa.PublicKey
	licenseKey string

	outFormat = ""
)

func init() {
	ValidateLicenseCmd.Flags().StringVar(&licenseKeyFile, "key-file", "", "location of file containing license key")
	ValidateLicenseCmd.Flags().StringVar(&privateKeyFile, "private", "", "path to private key file")
	ValidateLicenseCmd.Flags().StringVar(&publicKeyFile, "public", "", "path to public key file")
	ValidateLicenseCmd.Flags().StringVar(&outFormat, "out", "", "format in which to print license details. possible values are yaml, json")
}

var ValidateLicenseCmd = &cobra.Command{
	Use:     "license [license key]",
	Aliases: []string{"l", "lic"},
	Args:    cobra.MaximumNArgs(1),
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if privateKeyFile == "" && publicKeyFile == "" {
			return fmt.Errorf("one of --private, --public must be specified")
		}

		if len(args) == 0 && licenseKeyFile == "" {
			return fmt.Errorf("license key must be provided as argument, or --key-file flag must be set")
		}

		var err error
		if privateKeyFile != "" {
			if pubKey, err = loadPrivateKeyFile(privateKeyFile); err != nil {
				return err
			}
		} else {
			if pubKey, err = loadPublicKeyFile(publicKeyFile); err != nil {
				return err
			}
		}

		if len(args) > 0 {
			licenseKey = args[0]
		} else {
			keyData, err := os.ReadFile(licenseKeyFile)
			if err != nil {
				return err
			}

			licenseKey = string(keyData)
		}

		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		license, err := license.ParseLicenseKey([]byte(licenseKey), []*rsa.PublicKey{pubKey})
		if err != nil {
			return err
		}

		switch outFormat {
		case "yaml":
			out, err := yaml.Marshal(license)
			if err != nil {
				return err
			}
			fmt.Println(string(out))
		case "json":
			out, err := json.Marshal(license)
			if err != nil {
				return err
			}
			fmt.Println(string(out))
		default:
			fmt.Println("valid")
		}

		return nil
	},
}

func loadPrivateKeyFile(location string) (*rsa.PublicKey, error) {
	privateData, err := os.ReadFile(location)
	if err != nil {
		return nil, err
	}

	privKey, err := certificate.PEMToPrivateKey(privateData)
	if err != nil {
		return nil, err
	}

	return &privKey.PublicKey, nil
}

func loadPublicKeyFile(location string) (*rsa.PublicKey, error) {
	publicData, err := os.ReadFile(location)
	if err != nil {
		return nil, err
	}

	pubKey, err := certificate.PEMToPublicKey(publicData)
	if err != nil {
		return nil, err
	}

	return pubKey, nil
}
