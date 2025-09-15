package repositories

import (
	"database/sql"
	"encoding/json"

	"github.com/salahfarzin/roja-shop/types"
)

type Product interface {
	FetchAll() ([]types.Product, error)
	Create(product types.Product) (string, error)
	CreateWithFile(product types.Product, file *types.File) (string, error)
}

type product struct {
	db *sql.DB
}

func NewProduct(db *sql.DB) Product {
	return &product{db: db}
}

func (c *product) FetchAll() ([]types.Product, error) {
	query := `SELECT p.id, p.brand, p.title, p.inventory, p.price, p.old_price, p.discount, p.description, p.details, p.style_notes,
		f.id, f.file_name, f.file_path, f.file_type, f.created_at
	FROM products p
	LEFT JOIN files f ON p.id = f.product_id`
	rows, err := c.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []types.Product
	for rows.Next() {
		var p types.Product
		var detailsJSON, styleNotesJSON string
		var file types.File
		var fileID sql.NullString
		var fileName, filePath, fileType, fileCreatedAt sql.NullString
		err := rows.Scan(&p.ID, &p.Brand, &p.Title, &p.Inventory, &p.Price, &p.OldPrice, &p.Discount, &p.Description, &detailsJSON, &styleNotesJSON,
			&fileID, &fileName, &filePath, &fileType, &fileCreatedAt)
		if err != nil {
			return nil, err
		}
		json.Unmarshal([]byte(detailsJSON), &p.Details)
		json.Unmarshal([]byte(styleNotesJSON), &p.StyleNotes)
		if fileID.Valid {
			file.ID = fileID.String
			file.Name = fileName.String
			file.Path = filePath.String
			file.Type = fileType.String
			file.CreatedAt = fileCreatedAt.String
			p.File = &file
		}
		products = append(products, p)
	}
	return products, nil
}

func (c *product) Create(product types.Product) (string, error) {
	return "", nil
}

func (c *product) CreateWithFile(product types.Product, file *types.File) (string, error) {
	tx, err := c.db.Begin()
	if err != nil {
		return "", err
	}
	defer tx.Rollback()

	detailsJSON, _ := json.Marshal(product.Details)
	styleNotesJSON, _ := json.Marshal(product.StyleNotes)

	id := product.ID
	if id == "" {
		// Generate a UUID in Go, or use SQLite's randomblob if you want
		row := tx.QueryRow("SELECT lower(hex(randomblob(16)))")
		row.Scan(&id)
	}

	_, err = tx.Exec(`INSERT INTO products (id, brand, title, inventory, price, old_price, discount, description, details, style_notes) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		id, product.Brand, product.Title, product.Inventory, product.Price, product.OldPrice, product.Discount, product.Description, string(detailsJSON), string(styleNotesJSON))
	if err != nil {
		return "", err
	}

	if file != nil {
		_, err = tx.Exec(`INSERT INTO files (id, product_id, file_name, file_path, file_type, created_at) VALUES (?, ?, ?, ?, ?, ?)`,
			file.ID, id, file.Name, file.Path, file.Type, file.CreatedAt)
		if err != nil {
			return "", err
		}
	}

	if err := tx.Commit(); err != nil {
		return "", err
	}
	return id, nil
}
