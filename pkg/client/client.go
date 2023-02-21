package client

import (
	"crypto/rsa"
	"fmt"
	license2 "github.com/ebauman/golicense/pkg/license"
	"github.com/ebauman/golicense/pkg/types"
)

func Blocking(licenseKey string, usages map[string]int, publicKeys []*rsa.PublicKey) types.LicensingResponse {
	if licenseKey == "" {
		return types.NewLicensingResponse(false, types.NewInvalidLicenseError())
	}

	if len(publicKeys) == 0 {
		return types.NewLicensingResponse(false, types.NewInvalidPublicKeysError())
	}

	if len(usages) == 0 {
		return types.NewLicensingResponse(false, fmt.Errorf("usages map must not be empty"))
	}

	license, err := license2.ParseLicenseKey([]byte(licenseKey), publicKeys)
	if err != nil {
		return types.NewLicensingResponse(false, err)
	}

	// valid license, check usages
	for k, v := range usages {
		// check to see if there is a grant for what we're trying to use
		if _, ok := license.Grants[k]; !ok {
			return types.NewLicensingResponse(false, types.NewGrantNotFoundError(k, v))
		}

		// for each usage, look it up in the grants map
		// if the usage is higher than the grant, deny
		if v > license.Grants[k] {
			return types.NewLicensingResponse(false, types.NewGrantExceededError(k, v, license.Grants[k]))
		}
	}

	// at this point we have enough capacity to proceed, return nil
	return types.NewLicensingResponse(true, nil)
}

func Async(licenseKey string, usages map[string]int, publicKeys []*rsa.PublicKey, response chan<- types.LicensingResponse) {
	go func() {
		response <- Blocking(licenseKey, usages, publicKeys)
	}()
}
