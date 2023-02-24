package sqlite

import (
	"fmt"
	"github.com/ebauman/golicense/pkg/types"
	"github.com/google/uuid"
)

func (p *Provider) listLicensesForCriteria(key string, criteria string) ([]*types.License, error) {
	rows, err := p.db.Query(fmt.Sprintf("SELECT id, licensee, notBefore, notAfter, key, certificate "+
		"FROM %s WHERE %s = ?", licenseTable, criteria), key)
	if err != nil {
		return nil, err
	}

	var out = make([]*types.License, 0)
	for rows.Next() {
		var license = &types.License{}

		if err = rows.Scan(&license.Id, &license.Licensee, &license.NotBefore, &license.NotAfter, &license.Key, &license.Certificate); err != nil {
			return nil, err
		}

		if err = p.fillGrants(license); err != nil {
			return nil, err
		}

		if err = p.fillMetadata(license); err != nil {
			return nil, err
		}

		out = append(out, license)
	}

	return out, nil
}

func (p *Provider) ListLicensesForLicensee(licensee string) ([]*types.License, error) {
	return p.listLicensesForCriteria(licensee, licenseeTable)
}

func (p *Provider) ListLicensesForCertificate(certificate string) ([]*types.License, error) {
	return p.listLicensesForCriteria(certificate, certificateTable)
}

func (p *Provider) GetLicense(id string) (*types.License, error) {
	row := p.db.QueryRow(fmt.Sprintf("SELECT id, licensee, notBefore, notAfter, key, certificate "+
		"FROM %s WHERE id = ?", licenseTable), id)
	if row.Err() != nil {
		return nil, row.Err()
	}

	var license = &types.License{}

	if err := row.Scan(&license.Id, &license.Licensee, &license.NotBefore, &license.NotAfter, &license.Key,
		&license.Certificate); err != nil {
		return nil, err
	}

	if err := p.fillGrants(license); err != nil {
		return nil, err
	}

	if err := p.fillMetadata(license); err != nil {
		return nil, err
	}

	return license, nil
}

func (p *Provider) CreateLicense(license *types.License) (*types.License, error) {
	tx, err := p.db.Begin()
	if err != nil {
		return nil, err
	}

	if license.Id == "" {
		license.Id = uuid.NewString()
	}

	_, err = tx.Exec(fmt.Sprintf("INSERT INTO %s (id, licensee, notBefore, notAfter, key, certificate) VALUES "+
		"(?, ?, ?, ?, ?, ?)", licenseTable), license.Id, license.Licensee, license.NotBefore, license.NotAfter,
		license.Key, license.Certificate)
	if err != nil {
		return nil, handleRollback(tx, err)
	}

	if err = p.storeGrants(tx, license); err != nil {
		return nil, err
	}

	if err = p.storeMetadata(tx, license); err != nil {
		return nil, err
	}

	if err = tx.Commit(); err != nil {
		return nil, err
	}

	return license, nil
}

func (p *Provider) UpdateLicense(license *types.License) (*types.License, error) {
	tx, err := p.db.Begin()
	if err != nil {
		return nil, err
	}

	_, err = p.db.Exec(fmt.Sprintf("UPDATE %s SET licensee = ?, notBefore = ?, notAfter = ?, key = ?, certificate = ? "+
		"WHERE id = ?", licenseTable), license.Licensee, license.NotBefore, license.NotAfter, license.Key,
		license.Certificate, license.Id)

	if err != nil {
		return nil, err
	}

	if err = p.storeGrants(tx, license); err != nil {
		return nil, err
	}

	if err = p.storeMetadata(tx, license); err != nil {
		return nil, err
	}

	if err = tx.Commit(); err != nil {
		return nil, err
	}

	return license, nil
}

func (p *Provider) DeleteLicense(id string) error {
	tx, err := p.db.Begin()
	if err != nil {
		return err
	}

	if err = p.storeMetadata(tx, &types.License{Id: id}); err != nil {
		return handleRollback(tx, err)
	}

	if err = p.storeGrants(tx, &types.License{Id: id}); err != nil {
		return handleRollback(tx, err)
	}

	if _, err = tx.Exec(fmt.Sprintf("DELETE FROM %s WHERE id = ?", licenseTable), id); err != nil {
		return handleRollback(tx, err)
	}

	if err = tx.Commit(); err != nil {
		return err
	}

	return nil
}
