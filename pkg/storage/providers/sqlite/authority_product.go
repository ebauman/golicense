package sqlite

import (
	"database/sql"
	"fmt"
	"github.com/ebauman/golicense/pkg/types"
)

func (p *Provider) storeProductsForAuthority(tx *sql.Tx, authority *types.Authority) error {
	var err error
	var ourCommit = false
	if tx == nil {
		ourCommit = true
		tx, err = p.db.Begin()
		if err != nil {
			return err
		}
	}

	// first delete associated products
	_, err = tx.Exec(fmt.Sprintf("DELETE FROM %s WHERE authority = ?", authorityProductTable), authority.Id)
	if err != nil {
		if ourCommit {
			return handleRollback(tx, err)
		}
		return err
	}

	// now associate products
	statement, err := tx.Prepare(fmt.Sprintf("INSERT INTO %s (authority, product) VALUES (?, ?)", authorityProductTable))
	if err != nil {
		if ourCommit {
			return handleRollback(tx, err)
		}
		return err
	}

	defer statement.Close()

	for _, prod := range authority.Products {
		if _, err := statement.Exec(authority.Id, prod); err != nil {
			return err
		}
	}

	if ourCommit {
		if err = tx.Commit(); err != nil {
			return err
		}
	}

	return nil
}

func (p *Provider) fillProductsForAuthority(authority *types.Authority) error {
	// get associated products
	rows, err := p.db.Query(fmt.Sprintf("SELECT product_id FROM %s WHERE authority = ?", authorityProductTable), authority.Id)
	if err != nil {
		return err
	}

	var out = make([]string, 0)
	for rows.Next() {
		var product string
		if err := rows.Scan(&product); err != nil {
			return nil
		}

		out = append(out, product)
	}

	authority.Products = out

	return nil
}
