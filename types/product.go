package types

type Product struct {
	ID          string             `json:"id"`
	Brand       string             `json:"brand"`
	Title       string             `json:"title"`
	Inventory   int                `json:"inventory"`
	Price       float64            `json:"price"`
	OldPrice    *float64           `json:"old_price,omitempty"`
	Discount    *float64           `json:"discount,omitempty"`
	Description string             `json:"description"`
	Details     *map[string]string `json:"details"`
	StyleNotes  *map[string]string `json:"style_notes"`
	SoldCount   int                `json:"sold_count"`
	File        *File              `json:"file,omitempty"`
}
