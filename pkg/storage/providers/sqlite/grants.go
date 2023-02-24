package sqlite

import (
	"database/sql"
	"fmt"
	"github.com/ebauman/golicense/pkg/types"
)

func (p *Provider) fillGrants(license *types.License) error {
	rows, err := p.db.Query(fmt.Sprintf("SELECT key, value FROM %s WHERE license = ?", grantsTable),
		license.Id)

	if err != nil {
		return err
	}

	grants := make(map[string]int)
	for rows.Next() {
		var (
			key   string
			value int
		)

		err = rows.Scan(&key, &value)
		if err != nil {
			return err
		}

		grants[key] = value
	}

	license.Grants = grants

	return nil
}

func (p *Provider) storeGrants(tx *sql.Tx, license *types.License) error {
	var ourCommit = false // we don't want to commit a tx that's only half done
	var err error
	if tx == nil {
		ourCommit = true
		tx, err = p.db.Begin()
		if err != nil {
			return err
		}
	}

	_, err = tx.Exec(fmt.Sprintf("DELETE FROM %s WHERE license = ?", grantsTable), license.Id)
	if err != nil {
		if ourCommit {
			return handleRollback(tx, err)
		}
		return err
	}

	statement, err := tx.Prepare(fmt.Sprintf("INSERT INTO %s (key, value, license) VALUES (?, ?, ?)",
		grantsTable))
	defer statement.Close()

	if err != nil {
		if ourCommit {
			return handleRollback(tx, err)
		}
		return err
	}

	for k, v := range license.Grants {
		_, err := statement.Exec(k, v, license.Id)
		if err != nil {
			if ourCommit {
				return handleRollback(tx, err)
			}
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
