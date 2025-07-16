package models

// Category represents the 'categories' table in the database
type Category struct {
	CategoryID  int    `json:"category_id" db:"category_id"` // Corresponds to category_id in DB
	NamaKategori string `json:"nama_kategori" db:"nama_kategori"`
	ImageURL    string `json:"image_url" db:"image_url"` // New field for image path/URL
}