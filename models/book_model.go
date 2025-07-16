package models

// Book represents the 'books' table in the database
type Book struct {
	BookID      int    `json:"book_id" db:"book_id"` // Corresponds to book_id in DB
	Judul       string `json:"judul" db:"judul"`
	Penulis     string `json:"penulis" db:"penulis"`
	Penerbit    string `json:"penerbit" db:"penerbit"`
	TahunTerbit int    `json:"tahun_terbit" db:"tahun_terbit"`
	Sinopsis    string `json:"sinopsis" db:"sinopsis"`
	ImageURL    string `json:"image_url" db:"image_url"` // New field for image path/URL
	CategoryID  int    `json:"category_id" db:"category_id"` // Foreign key to categories
}