package sqlite

import (
	"database/sql"
	"fmt"
	"github.com/ebauman/golicense/pkg/types"
)

func (p *Provider) fillMetadata(obj types.MetaObject) error {
	rows, err := p.db.Query(fmt.Sprintf("SELECT key, value FROM %s WHERE kind = ? AND id = ?",
		metadataTable), obj.GetKind(), obj.GetId())
	if err != nil {
		return err
	}

	out := make(map[string]string)

	for rows.Next() {
		var (
			key   string
			value string
		)

		if err := rows.Scan(&key, &value); err != nil {
			return err
		}

		out[key] = value
	}

	obj.SetMetadata(out)

	return nil
}

func (p *Provider) storeMetadata(tx *sql.Tx, obj types.MetaObject) error {
	var err error
	var ourCommit = false
	if tx == nil {
		ourCommit = true
		tx, err = p.db.Begin()
		if err != nil {
			return err
		}
	}

	// first, delete existing metadata
	_, err = tx.Exec(fmt.Sprintf("DELETE FROM %s WHERE kind = ? AND id = ?", metadataTable),
		obj.GetKind(), obj.GetId())
	if err != nil {
		if ourCommit {
			return handleRollback(tx, err)
		}
		return err
	}

	statement, err := tx.Prepare(fmt.Sprintf("INSERT INTO %s (key, value, kind, id) VALUES "+
		"(?, ? ,? ,?)", metadataTable))
	if err != nil {
		if ourCommit {
			return handleRollback(tx, err)
		}
		return err
	}

	defer statement.Close()

	for k, v := range obj.GetMetadata() {
		if _, err := statement.Exec(k, v, obj.GetKind(), obj.GetId()); err != nil {
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
