package sqlite

import (
	"fmt"
	"github.com/ebauman/golicense/pkg/types"
	"github.com/google/uuid"
)

func (p *Provider) ListProducts() ([]*types.Product, error) {
	rows, err := p.db.Query(fmt.Sprintf("SELECT id, name, unit FROM %s", productTable))
	if err != nil {
		return nil, err
	}

	products := make([]*types.Product, 0)
	for rows.Next() {
		var product = &types.Product{}

		if err := rows.Scan(&product.Id, &product.Name, &product.Unit); err != nil {
			return nil, err
		}

		if err = p.fillMetadata(product); err != nil {
			return nil, err
		}

		products = append(products, product)
	}

	return products, nil
}

func (p *Provider) GetProduct(id string) (*types.Product, error) {
	row := p.db.QueryRow(fmt.Sprintf("SELECT id, name, unit FROM %s WHERE id = ?", productTable),
		id)
	if row.Err() != nil {
		return nil, row.Err()
	}

	var product = &types.Product{}

	if err := row.Scan(&product.Id, &product.Name, &product.Unit); err != nil {
		return nil, err
	}

	if err := p.fillMetadata(product); err != nil {
		return nil, err
	}

	return product, nil
}

func (p *Provider) CreateProduct(inputProduct *types.Product) (*types.Product, error) {
	product := inputProduct.DeepCopyObject().(*types.Product)
	tx, err := p.db.Begin()
	if err != nil {
		return nil, err
	}

	if product.Id == "" {
		product.Id = uuid.NewString()
	}

	if _, err = tx.Exec(fmt.Sprintf("INSERT INTO %s (id, name, unit) VALUES (?, ?, ?)", productTable),
		product.Id, product.Name, product.Unit); err != nil {
		return nil, handleRollback(tx, err)
	}

	if err = p.storeMetadata(tx, product); err != nil {
		return nil, handleRollback(tx, err)
	}

	if err = tx.Commit(); err != nil {
		return nil, err
	}

	return product, nil
}

func (p *Provider) UpdateProduct(inputProduct *types.Product) (*types.Product, error) {
	product := inputProduct.DeepCopyObject().(*types.Product)
	tx, err := p.db.Begin()
	if err != nil {
		return nil, err
	}

	if _, err := p.db.Exec(fmt.Sprintf("UPDATE %s SET name = ?, unit = ? WHERE id = ?", productTable),
		product.Name, product.Unit); err != nil {
		return nil, handleRollback(tx, err)
	}

	if err = p.storeMetadata(tx, product); err != nil {
		return nil, handleRollback(tx, err)
	}

	if err = tx.Commit(); err != nil {
		return nil, err
	}

	return product, nil
}

func (p *Provider) DeleteProduct(id string) error {
	tx, err := p.db.Begin()
	if err != nil {
		return err
	}

	if _, err = tx.Exec(fmt.Sprintf("DELETE FROM %s WHERE product = ?", authorityProductTable), id); err != nil {
		return handleRollback(tx, err)
	}

	if _, err = tx.Exec(fmt.Sprintf("DELETE FROM %s WHERE id = ?", productTable), id); err != nil {
		return handleRollback(tx, err)
	}

	if err = tx.Commit(); err != nil {
		return err
	}

	return nil
}
