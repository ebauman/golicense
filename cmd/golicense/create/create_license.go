package create

import (
	"fmt"
	"github.com/ebauman/golicense/pkg/certificate"
	"github.com/ebauman/golicense/pkg/flag"
	golicense "github.com/ebauman/golicense/pkg/license"
	"github.com/ebauman/golicense/pkg/types"
	"github.com/google/uuid"
	"github.com/spf13/cobra"
	"os"
	"time"
)

var (
	metadataSlice = []string{}
	grantSlice    = []string{}
	licensee      = ""
	privateKey    = ""
	notBefore     = ""
	notAfter      = ""
)

func init() {
	CreateLicenseCmd.Flags().StringVar(&licensee, "licensee", "", "licensee")
	CreateLicenseCmd.Flags().StringSliceVar(&grantSlice, "grant", []string{}, "grant")
	CreateLicenseCmd.Flags().StringSliceVar(&metadataSlice, "metadata", []string{}, "metadata")
	CreateLicenseCmd.Flags().StringVar(&privateKey, "key", "", "path to private key")
	CreateLicenseCmd.Flags().StringVar(&notBefore, "not-before", time.Now().Format("2006-01-02"), "license not before")
	CreateLicenseCmd.Flags().StringVar(&notAfter, "not-after", "", "license not after")

	if err := CreateLicenseCmd.MarkFlagFilename("key"); err != nil {
		panic(err)
	}

	for _, v := range []string{"licensee", "grant", "key", "not-after"} {
		if err := CreateLicenseCmd.MarkFlagRequired(v); err != nil {
			panic(err)
		}
	}
}

var CreateLicenseCmd = &cobra.Command{
	Use:     "license",
	Short:   "create license",
	Aliases: []string{"l", "lic"},
	RunE: func(cmd *cobra.Command, args []string) error {
		grants, err := flag.ParseGrantFlags(grantSlice)
		if err != nil {
			return err
		}

		metadata, err := flag.ParseMetadataFlags(metadataSlice)
		if err != nil {
			return err
		}

		var notAfterTime time.Time
		if notAfterTime, err = flag.ParseTime(notAfter); err != nil {
			return err
		}

		var notBeforeTime time.Time
		if notBefore != "" {
			if notBeforeTime, err = flag.ParseTime(notBefore); err != nil {
				return err
			}
		}

		keyData, err := os.ReadFile(privateKey)
		if err != nil {
			return err
		}

		key, err := certificate.PEMToPrivateKey(keyData)
		if err != nil {
			return err
		}

		license := types.License{
			Id:        uuid.NewString(),
			Grants:    grants,
			Metadata:  metadata,
			NotAfter:  notAfterTime,
			NotBefore: notBeforeTime,
			Licensee:  licensee,
		}

		keystring, err := golicense.GenerateLicenseKey(key, license)
		if err != nil {
			return err
		}

		fmt.Println(keystring)

		return nil
	},
}
