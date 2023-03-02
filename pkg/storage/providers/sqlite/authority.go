package sqlite

import (
	"fmt"
	"github.com/ebauman/golicense/pkg/types"
	"github.com/google/uuid"
)

func (p *Provider) ListAuthorities() ([]*types.Authority, error) {
	rows, err := p.db.Query(fmt.Sprintf("SELECT id, name FROM %s", authorityTable))
	if err != nil {
		return nil, err
	}

	var out = make([]*types.Authority, 0)
	for rows.Next() {
		var authority = &types.Authority{}
		err := rows.Scan(&authority.Id, &authority.Name)
		if err != nil {
			return nil, err
		}

		if err := p.fillMetadata(authority); err != nil {
			return nil, err
		}

		if err := p.fillProductsForAuthority(authority); err != nil {
			return nil, err
		}

		out = append(out, authority)
	}

	return out, nil
}

func (p *Provider) GetAuthority(id string) (*types.Authority, error) {
	row := p.db.QueryRow(fmt.Sprintf("SELECT id, name FROM %s WHERE id = ?", authorityTable), id)
	if row.Err() != nil {
		return nil, row.Err()
	}

	var authority = &types.Authority{}

	err := row.Scan(&authority.Id, &authority.Name)
	if err != nil {
		return nil, err
	}

	if err = p.fillMetadata(authority); err != nil {
		return nil, err
	}

	if err = p.fillProductsForAuthority(authority); err != nil {
		return nil, err
	}

	return authority, nil
}

func (p *Provider) CreateAuthority(inputAuthority *types.Authority) (*types.Authority, error) {
	authority := inputAuthority.DeepCopyObject().(*types.Authority)
	tx, err := p.db.Begin()
	if err != nil {
		return nil, err
	}

	if authority.GetId() == "" {
		authority.Id = uuid.NewString()
	}

	_, err = tx.Exec(fmt.Sprintf("INSERT INTO %s (id, name) VALUES (?, ?)", authorityTable),
		authority.Id, authority.Name)
	if err != nil {
		return nil, handleRollback(tx, err)
	}

	if err = p.storeMetadata(tx, authority); err != nil {
		return nil, handleRollback(tx, err)
	}

	if err = p.storeProductsForAuthority(tx, authority); err != nil {
		return nil, handleRollback(tx, err)
	}

	if err = tx.Commit(); err != nil {
		return nil, err
	}

	return authority, nil
}

func (p *Provider) UpdateAuthority(inputAuthority *types.Authority) (*types.Authority, error) {
	authority := inputAuthority.DeepCopyObject().(*types.Authority)
	tx, err := p.db.Begin()
	if err != nil {
		return nil, err
	}

	_, err = tx.Exec(fmt.Sprintf("UPDATE %s SET name = ? WHERE id = ?", authorityTable),
		authority.Name, authority.Id)
	if err != nil {
		return nil, handleRollback(tx, err)
	}

	if err = p.storeMetadata(tx, authority); err != nil {
		return nil, handleRollback(tx, err)
	}

	if err = p.storeProductsForAuthority(tx, authority); err != nil {
		return nil, handleRollback(tx, err)
	}

	if err = tx.Commit(); err != nil {
		return nil, err
	}

	return authority, nil
}

func (p *Provider) DeleteAuthority(id string) error {
	tx, err := p.db.Begin()
	if err != nil {
		return err
	}

	// go through each table that references authority and delete
	for _, kind := range []string{licenseeTable, certificateTable, authorityProductTable} {
		_, err = tx.Exec(fmt.Sprintf("DELETE FROM %s WHERE authority = ?", kind), id)
		if err != nil {
			return handleRollback(tx, err)
		}
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}
