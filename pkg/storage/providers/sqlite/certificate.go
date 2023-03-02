package sqlite

import (
	"fmt"
	certificate2 "github.com/ebauman/golicense/pkg/certificate"
	"github.com/ebauman/golicense/pkg/types"
	"github.com/google/uuid"
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
		var base64Key string
		err = rows.Scan(&certificate.Id, &certificate.Authority, &base64Key)
		if err != nil {
			return nil, err
		}

		privKey, err := certificate2.PEMToPrivateKey([]byte(base64Key))
		if err != nil {
			return nil, err
		}

		certificate.PrivateKey = *privKey

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

	var b64Cert string
	err := row.Scan(&certificate.Id, &certificate.Authority, &b64Cert)
	if err != nil {
		return nil, err
	}

	// convert back into PrivateKey
	cert, err := certificate2.PEMToPrivateKey([]byte(b64Cert))
	if err != nil {
		return nil, err
	}

	certificate.PrivateKey = *cert

	if err = p.fillMetadata(certificate); err != nil {
		return nil, err
	}

	return certificate, nil
}

func (p *Provider) CreateCertificate(inputCertificate *types.Certificate) (*types.Certificate, error) {
	certificate := inputCertificate.DeepCopyObject().(*types.Certificate)
	tx, err := p.db.Begin()
	if err != nil {
		return nil, err
	}

	if certificate.Id == "" {
		certificate.Id = uuid.NewString()
	}

	b64Key, err := certificate2.KeyToPrivatePEM(&certificate.PrivateKey)
	if err != nil {
		return nil, handleRollback(tx, err)
	}

	_, err = tx.Exec(fmt.Sprintf("INSERT INTO %s (id, authority, privateKey) VALUES (?, ?, ?)",
		certificateTable), certificate.Id, certificate.Authority, b64Key)
	if err != nil {
		return nil, handleRollback(tx, err)
	}

	// persist metadata
	if err = p.storeMetadata(tx, certificate); err != nil {
		return nil, err
	}

	return certificate, nil
}

func (p *Provider) UpdateCertificate(inputCertificate *types.Certificate) (*types.Certificate, error) {
	certificate := inputCertificate.DeepCopyObject().(*types.Certificate)
	tx, err := p.db.Begin()
	if err != nil {
		return nil, err
	}

	b64Key, err := certificate2.KeyToPrivatePEM(&certificate.PrivateKey)
	if err != nil {
		return nil, handleRollback(tx, err)
	}

	_, err = tx.Exec(fmt.Sprintf("UPDATE %s SET authority = ?, privateKey = ? WHERE id = ?",
		certificateTable), certificate.Authority, b64Key, certificate.Id)
	if err != nil {
		return nil, handleRollback(tx, err)
	}

	if err = p.storeMetadata(tx, certificate); err != nil {
		return nil, handleRollback(tx, err)
	}

	return certificate, nil
}

func (p *Provider) DeleteCertificate(id string) error {
	tx, err := p.db.Begin()
	if err != nil {
		return err
	}

	if err = p.storeMetadata(tx, &types.Certificate{Id: id}); err != nil {
		return handleRollback(tx, err)
	}

	if _, err = tx.Exec(fmt.Sprintf("DELETE FROM %s WHERE id = ?", certificateTable), id); err != nil {
		return handleRollback(tx, err)
	}

	return nil
}
