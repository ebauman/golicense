package client

import (
	"crypto/rand"
	"crypto/rsa"
	"github.com/ebauman/golicense/pkg/license"
	"github.com/ebauman/golicense/pkg/types"
	"github.com/google/uuid"
	"testing"
	"time"
)

func validKeyLicense() (*rsa.PrivateKey, string, error) {
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, "", err
	}

	// license
	lic := types.License{
		Id:       uuid.NewString(),
		Licensee: "testing",
		Metadata: map[string]string{
			"golicense.io/authority": "testing-authority",
		},
		Grants: map[string]int{
			"golicense.io/nodes": 50,
		},
		NotBefore: time.Now(),
		NotAfter:  time.Now().Add(time.Second * 31536000), // five hundred twenty-five thousand six hundred minutes
	}

	licenseKey, err := license.GenerateLicenseKey(key, lic)
	if err != nil {
		return nil, "", err
	}

	return key, licenseKey, nil
}

func TestBlocking_ValidLicenseValidKeys(t *testing.T) {
	privateKey, licenseKey, err := validKeyLicense()
	if err != nil {
		t.Error(err)
	}

	response := Blocking(licenseKey, map[string]int{"golicense.io/nodes": 10}, []*rsa.PublicKey{&privateKey.PublicKey})
	if !response.Success || response.Error != nil {
		t.Errorf("unsuccessful licensing, should have succeeded. error: %s", err.Error())
	}
}

func TestAsync_ValidLicenseValidKeys(t *testing.T) {
	privateKey, licenseKey, err := validKeyLicense()
	if err != nil {
		t.Error(err)
	}

	resChan := make(chan types.LicensingResponse)

	Async(licenseKey, map[string]int{"golicense.io/nodes": 10}, []*rsa.PublicKey{&privateKey.PublicKey}, resChan)

	response := <-resChan

	if !response.Success || response.Error != nil {
		t.Errorf("unsuccessful licensing, should have succeeded. error: %s", err.Error())
	}
}

func TestBlocking_InvalidLicenseValidKeys(t *testing.T) {
	privateKey, _, err := validKeyLicense()
	if err != nil {
		t.Error(err)
	}

	response := Blocking("lkjsdflkjsdf", map[string]int{"golicense.io/nodes": 10}, []*rsa.PublicKey{&privateKey.PublicKey})

	if response.Success || response.Error == nil {
		t.Errorf("successful licensing, should return invalid license")
	}

	if !types.IsError(response.Error, types.InvalidLicenseError{}) {
		t.Errorf("error returned is not an InvalidLicenseError")
	}
}

func TestAsync_InvalidLicenseValidKeys(t *testing.T) {
	privateKey, _, err := validKeyLicense()
	if err != nil {
		t.Error(err)
	}

	resChan := make(chan types.LicensingResponse)
	Async("lkjsdflkjsdflkjsdf", map[string]int{"golicense.io/nodes": 10}, []*rsa.PublicKey{&privateKey.PublicKey}, resChan)

	response := <-resChan
	if response.Success || response.Error == nil {
		t.Errorf("successful licensing, should return invalid license")
	}

	if !types.IsErrorReflect(response.Error, types.InvalidLicenseError{}) {
		t.Errorf("error returned is not an InvalidLicenseError")
	}
}

func TestBlocking_ValidLicenseInvalidKeys(t *testing.T) {
	_, licenseKey, err := validKeyLicense()
	if err != nil {
		t.Error(err)
	}

	// make new cert that's invalid for what we just generated
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	response := Blocking(licenseKey, map[string]int{"golicense.io/nodes": 10}, []*rsa.PublicKey{&privateKey.PublicKey})

	if response.Success || response.Error == nil {
		t.Errorf("successful licensing, should return invalid")
	}
}

func TestAsync_ValidLicenseInvalidKeys(t *testing.T) {
	_, licenseKey, err := validKeyLicense()
	if err != nil {
		t.Error(err)
	}

	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)

	resChan := make(chan types.LicensingResponse)
	Async(licenseKey, map[string]int{"golicense.io/nodes": 10}, []*rsa.PublicKey{&privateKey.PublicKey}, resChan)

	response := <-resChan

	if response.Success || response.Error == nil {
		t.Errorf("successful licensing, should return invalid")
	}
}
