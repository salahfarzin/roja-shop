package repositories

import (
	"database/sql"
	"time"

	"github.com/salahfarzin/roja-shop/types"
)

type Sale interface {
	Create(sale *types.Sale) error
	FetchAll() ([]*types.Sale, error)
}

type sale struct {
	db *sql.DB
}

func NewSaleRepo(db *sql.DB) Sale {
	return &sale{db: db}
}

func (r *sale) FetchAll() ([]*types.Sale, error) {
	rows, err := r.db.Query(`SELECT id, product_id, quantity, created_at FROM sales`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var sales []*types.Sale
	for rows.Next() {
		var s types.Sale
		var dt string
		if err := rows.Scan(&s.ID, &s.ProductID, &s.Quantity, &dt); err != nil {
			return nil, err
		}
		t, _ := time.Parse(time.RFC3339, dt)
		s.CreatedAt = t
		sales = append(sales, &s)
	}
	return sales, nil
}

func (r *sale) Create(sale *types.Sale) error {
	_, err := r.db.Exec(`INSERT INTO sales (id, product_id, quantity, created_at) VALUES (?, ?, ?, ?)`,
		sale.ID, sale.ProductID, sale.Quantity, sale.CreatedAt)
	return err
}
