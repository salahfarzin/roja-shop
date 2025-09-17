package repositories

import (
	"database/sql"
	"encoding/json"

	"github.com/salahfarzin/roja-shop/types"
)

type Product interface {
	FetchAll(limit, offset int) ([]types.Product, error)
	FetchOne(id string) (*types.Product, error)
	Create(product types.Product) (string, error)
	CreateWithFile(product types.Product, file *types.File) (string, error)
	Update(id string, input types.Product) error
}

type product struct {
	db *sql.DB
}

func NewProduct(db *sql.DB) Product {
	return &product{db: db}
}

func (c *product) FetchOne(id string) (*types.Product, error) {
	query := `SELECT p.id, p.brand, p.title, p.inventory, p.sold_count, p.price, p.old_price, p.discount, p.description, p.details, p.style_notes,
		f.id, f.file_name, f.file_path, f.file_type, f.created_at
	FROM products p
	LEFT JOIN files f ON p.id = f.product_id
	WHERE p.id = ?`
	row := c.db.QueryRow(query, id)
	var p types.Product
	var detailsJSON, styleNotesJSON string
	var file types.File
	var fileID, fileName, filePath, fileType, fileCreatedAt sql.NullString
	err := row.Scan(&p.ID, &p.Brand, &p.Title, &p.Inventory, &p.SoldCount, &p.Price, &p.OldPrice, &p.Discount, &p.Description, &detailsJSON, &styleNotesJSON,
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
	return &p, nil
}

func (c *product) FetchAll(limit, offset int) ([]types.Product, error) {
	query := `SELECT p.id, p.brand, p.title, p.inventory, p.price, p.old_price, p.discount, p.description, p.details, p.style_notes,
		f.id, f.file_name, f.file_path, f.file_type, f.created_at
	FROM products p
	LEFT JOIN files f ON p.id = f.product_id
	LIMIT ? OFFSET ?`
	rows, err := c.db.Query(query, limit, offset)
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

func (c *product) Update(id string, product types.Product) error {
	// Build dynamic SQL for only provided fields
	fields := []string{}
	args := []interface{}{}

	if product.Brand != "" {
		fields = append(fields, "brand = ?")
		args = append(args, product.Brand)
	}
	if product.Title != "" {
		fields = append(fields, "title = ?")
		args = append(args, product.Title)
	}
	if product.Inventory >= 0 {
		fields = append(fields, "inventory = ?")
		args = append(args, product.Inventory)
	}
	if product.SoldCount != 0 {
		fields = append(fields, "sold_count = ?")
		args = append(args, product.SoldCount)
	}
	if product.Price >= 0 {
		fields = append(fields, "price = ?")
		args = append(args, product.Price)
	}
	if product.OldPrice != nil {
		fields = append(fields, "old_price = ?")
		args = append(args, product.OldPrice)
	}
	if product.Discount != nil {
		fields = append(fields, "discount = ?")
		args = append(args, product.Discount)
	}
	if product.Description != "" {
		fields = append(fields, "description = ?")
		args = append(args, product.Description)
	}
	if product.Details != nil {
		detailsJSON, _ := json.Marshal(product.Details)
		fields = append(fields, "details = ?")
		args = append(args, string(detailsJSON))
	}
	if product.StyleNotes != nil {
		styleNotesJSON, _ := json.Marshal(product.StyleNotes)
		fields = append(fields, "style_notes = ?")
		args = append(args, string(styleNotesJSON))
	}

	if len(fields) == 0 {
		return nil // nothing to update
	}

	query := "UPDATE products SET " + fields[0]
	for _, f := range fields[1:] {
		query += ", " + f
	}
	query += " WHERE id = ?"
	args = append(args, id)

	_, err := c.db.Exec(query, args...)
	return err
}
