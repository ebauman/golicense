package license

import (
	"bytes"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/ebauman/golicense/pkg/types"
	"strings"
)

func ParseLicenseKey(licenseBytes []byte, publicKeys []*rsa.PublicKey) (*types.License, error) {
	if len(licenseBytes) == 0 {
		return nil, types.InvalidLicenseError{}
	}

	licenseSlice := strings.Split(string(licenseBytes), ".")

	if len(licenseSlice) < 2 {
		return nil, types.NewInvalidLicenseError()
	}

	hash := sha256.New()
	_, err := hash.Write([]byte(licenseSlice[0]))
	if err != nil {
		return nil, types.InvalidLicenseError{}
	}

	hashSum := hash.Sum(nil)

	signature, err := base64.StdEncoding.DecodeString(licenseSlice[1])
	if err != nil {
		return nil, types.InvalidLicenseError{}
	}

	var valid = false
	for _, key := range publicKeys {
		err = rsa.VerifyPSS(key, crypto.SHA256, hashSum, signature, nil)
		if err == nil {
			valid = true
			break
		}
	}

	if !valid {
		return nil, types.InvalidLicenseError{}
	}

	licenseJson, err := base64.StdEncoding.DecodeString(licenseSlice[0])
	if err != nil {
		return nil, types.InvalidLicenseError{}
	}

	var license = types.License{}
	if err = json.Unmarshal(licenseJson, &license); err != nil {
		return nil, types.InvalidLicenseError{}
	}

	return &license, nil
}

func GenerateLicenseKey(key *rsa.PrivateKey, license types.License) (string, error) {
	licenseJson, err := json.Marshal(license)
	if err != nil {
		return "", nil
	}

	base64License := encode(licenseJson)

	hash := sha256.New()
	_, err = hash.Write(base64License)
	if err != nil {
		return "", err
	}

	hashSum := hash.Sum(nil)

	signature, err := rsa.SignPSS(rand.Reader, key, crypto.SHA256, hashSum, nil)
	if err != nil {
		return "", err
	}

	base64Signature := encode(signature)

	return fmt.Sprintf("%s.%s", base64License, base64Signature), nil
}

func encode(in []byte) []byte {
	var buf = &bytes.Buffer{}

	encoder := base64.NewEncoder(base64.StdEncoding, buf)

	encoder.Write(in)

	encoder.Close()

	return buf.Bytes()
}
