package types

type File struct {
	ID        string `json:"id"`
	ProductID string `json:"product_id"`
	Name      string `json:"name"`
	Path      string `json:"path"`
	Type      string `json:"type"`
	CreatedAt string `json:"created_at"`
}
