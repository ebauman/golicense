package certificate

import (
	"bytes"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
)

var (
	rsaPrivateKeyBlock = "RSA PRIVATE KEY"
	publicCertBlock    = "CERTIFICATE"
)

func KeyToPrivatePEM(key *rsa.PrivateKey) ([]byte, error) {
	keyPem := new(bytes.Buffer)
	err := pem.Encode(keyPem, &pem.Block{
		Type:  rsaPrivateKeyBlock,
		Bytes: x509.MarshalPKCS1PrivateKey(key),
	})
	if err != nil {
		return nil, err
	}

	return keyPem.Bytes(), nil
}

func KeyToPublicPEM(key *rsa.PrivateKey) ([]byte, error) {
	publicPem := new(bytes.Buffer)
	err := pem.Encode(publicPem, &pem.Block{
		Type:  publicCertBlock,
		Bytes: x509.MarshalPKCS1PublicKey(&key.PublicKey),
	})
	if err != nil {
		return nil, err
	}

	return publicPem.Bytes(), nil
}

func PEMToPublicKey(data []byte) (*rsa.PublicKey, error) {
	pBlock, err := pemDecode(data)
	if err != nil {
		return nil, err
	}

	if pBlock.Type != publicCertBlock {
		return nil, fmt.Errorf("invalid pem data, not public key")
	}

	return x509.ParsePKCS1PublicKey(pBlock.Bytes)
}

func PEMToPrivateKey(data []byte) (*rsa.PrivateKey, error) {
	pBlock, err := pemDecode(data)
	if err != nil {
		return nil, err
	}

	if pBlock.Type != rsaPrivateKeyBlock {
		return nil, fmt.Errorf("invalid pem data, not private key")
	}

	return x509.ParsePKCS1PrivateKey(pBlock.Bytes)
}

func pemDecode(data []byte) (*pem.Block, error) {
	pBlock, _ := pem.Decode(data)

	if pBlock == nil {
		return nil, fmt.Errorf("error decoding pem data")
	}

	return pBlock, nil
}
