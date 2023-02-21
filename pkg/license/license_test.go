package license

import (
	"crypto/rand"
	"crypto/rsa"
	"github.com/ebauman/golicense/pkg/types"
	"github.com/google/uuid"
	"testing"
	"time"
)

func Test_LicenseGenAndValidate(t *testing.T) {
	// use a random key
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Error(err)
	}

	license := types.License{
		Id:       uuid.NewString(),
		Licensee: "eamon",
		Metadata: nil,
		Grants: map[string]int{
			"app.local/nodes": 5,
		},
		NotBefore: time.Now(),
		NotAfter:  time.Now().Add(time.Second * 3600),
	}

	lic, err := GenerateLicenseKey(key, license)
	if err != nil {
		t.Error(err)
	}

	// validate using same data
	_, err = ParseLicenseKey([]byte(lic), []*rsa.PublicKey{&key.PublicKey})
	if err != nil {
		t.Error(err)
	}
}
