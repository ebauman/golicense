package sqlite

import (
	"fmt"
	"github.com/ebauman/golicense/pkg/types"
	"github.com/google/uuid"
)

func (p *Provider) ListLicensees(authority string) ([]*types.Licensee, error) {
	rows, err := p.db.Query(fmt.Sprintf("SELECT id, name, authority FROM %s WHERE authority = ?",
		licenseeTable), authority)

	if err != nil {
		return nil, err
	}

	var out = make([]*types.Licensee, 0)
	for rows.Next() {
		var licensee = &types.Licensee{}
		if err := rows.Scan(&licensee.Id, &licensee.Name, &licensee.Authority); err != nil {
			return nil, err
		}

		if err = p.fillMetadata(licensee); err != nil {
			return nil, err
		}

		out = append(out, licensee)
	}

	return out, nil
}

func (p *Provider) GetLicensee(id string) (*types.Licensee, error) {
	row := p.db.QueryRow(fmt.Sprintf("SELECT id, name, authority FROM %s WHERE id = ?",
		licenseeTable), id)

	if row.Err() != nil {
		return nil, row.Err()
	}

	var licensee = &types.Licensee{}

	if err := row.Scan(&licensee.Id, &licensee.Name, &licensee.Authority); err != nil {
		return nil, err
	}

	if err := p.fillMetadata(licensee); err != nil {
		return nil, err
	}

	return licensee, nil
}

func (p *Provider) CreateLicensee(inputLicensee *types.Licensee) (*types.Licensee, error) {
	licensee := inputLicensee.DeepCopyObject().(*types.Licensee)
	tx, err := p.db.Begin()
	if err != nil {
		return nil, err
	}

	if licensee.Id == "" {
		licensee.Id = uuid.NewString()
	}

	_, err = tx.Exec(fmt.Sprintf("INSERT INTO %s (id, name, authority) VALUES (?, ?, ?)",
		licenseeTable), licensee.Id, licensee.Name, licensee.Authority)

	if err != nil {
		return nil, handleRollback(tx, err)
	}

	if err = p.storeMetadata(tx, licensee); err != nil {
		return nil, handleRollback(tx, err)
	}

	if err = tx.Commit(); err != nil {
		return nil, err
	}

	return licensee, nil
}

func (p *Provider) UpdateLicensee(inputLicensee *types.Licensee) (*types.Licensee, error) {
	licensee := inputLicensee.DeepCopyObject().(*types.Licensee)
	tx, err := p.db.Begin()
	if err != nil {
		return nil, err
	}

	if _, err = tx.Exec(fmt.Sprintf("UPDATE %s SET name = ?, authority = ? WHERE id = ?",
		licenseeTable),
		licensee.Name, licensee.Authority, licensee.Id); err != nil {
		return nil, handleRollback(tx, err)
	}

	if err = p.storeMetadata(tx, licensee); err != nil {
		return nil, handleRollback(tx, err)
	}

	if err = tx.Commit(); err != nil {
		return nil, err
	}

	return licensee, nil
}

func (p *Provider) DeleteLicensee(id string) error {
	tx, err := p.db.Begin()
	if err != nil {
		return err
	}

	_, err = tx.Exec(fmt.Sprintf("DELETE FROM %s WHERE licensee = %s", licenseTable), id)
	if err != nil {
		return handleRollback(tx, err)
	}

	_, err = tx.Exec(fmt.Sprintf("DELETE FROM %s WHERE id = ?", licenseeTable), id)
	if err != nil {
		return handleRollback(tx, err)
	}

	if err = tx.Commit(); err != nil {
		return err
	}

	return nil
}
