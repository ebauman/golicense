package sqlite

import (
	"fmt"
	"github.com/ebauman/golicense/pkg/types"
)

func (p *Provider) ListCertificates(authority string) ([]*types.Certificate, error) {
	rows, err := p.db.Query(fmt.Sprintf("SELECT id, authority, privateKey FROM %s WHERE authority = ?",
		certificateTable), authority)
	if err != nil {
		return nil, err
	}

	var out = make([]*types.Certificate, 0)
	for rows.Next() {
		var certificate = &types.Certificate{}
		err = rows.Scan(&certificate.Id, &certificate.Authority, &certificate.PrivateKey)
		if err != nil {
			return nil, err
		}

		if err = p.fillMetadata(certificate); err != nil {
			return nil, err
		}

		out = append(out, certificate)
	}

	return out, nil
}

func (p *Provider) GetCertificate(id string) (*types.Certificate, error) {
	row := p.db.QueryRow(fmt.Sprintf("SELECT id, authority, privateKey FROM %s WHERE id = ?",
		certificateTable), id)

	if row.Err() != nil {
		return nil, row.Err()
	}

	var certificate = &types.Certificate{}

	err := row.Scan(&certificate.Id, &certificate.Authority, &certificate.PrivateKey)
	if err != nil {
		return nil, err
	}

	if err = p.fillMetadata(certificate); err != nil {
		return nil, err
	}

	return certificate, nil
}

// @TODO - Finish up this provider. Need to figure out how to store private key - probably as
// a base64 string? That way there's no weirdness with storing newlines, etc.
// Gotta de-base64 that on the way out and parse it into a PrivateKey. Use the existing
// methods in github.com/ebauman/golicense/pkg/certificate
//
// Then test it!
