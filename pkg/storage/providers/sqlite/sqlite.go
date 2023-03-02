package sqlite

import (
	"database/sql"
	"errors"
	"github.com/ebauman/golicense/pkg/storage"
	_ "modernc.org/sqlite"
)

const (
	authorityTable        = "authority"
	productTable          = "product"
	authorityProductTable = "authority_product"
	licenseeTable         = "licensee"
	licenseTable          = "license"
	certificateTable      = "certificate"
	metadataTable         = "metadata"
	grantsTable           = "grants"
)

type Provider struct {
	db *sql.DB
}

func NewSqliteProvider(dsn string) (storage.GolicenseStore, error) {
	db, err := sql.Open("sqlite", dsn)
	if err != nil {
		return nil, err
	}

	return &Provider{db: db}, nil
}

func handleRollback(transaction *sql.Tx, err error) error {
	rollbackError := transaction.Rollback()
	if rollbackError != nil {
		return errors.Join(rollbackError, err)
	}

	return err
}
