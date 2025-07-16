package handlers

import (
	"database/sql"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"pojok_baca_api/database"
	"pojok_baca_api/models"
	"pojok_baca_api/utils"
	"strconv" // For converting string to int
)

// GetAllBooks gets all books from the database
// GET /api/v1/books
func GetAllBooks(c *fiber.Ctx) error {
	rows, err := database.DB.Query("SELECT book_id, judul, penulis, penerbit, tahun_terbit, sinopsis, image_url, category_id FROM books")
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, fmt.Sprintf("Failed to retrieve books: %v", err))
	}
	defer rows.Close()

	var books []models.Book
	for rows.Next() {
		var book models.Book
		// Scan into book struct, including the image_url field
		if err := rows.Scan(&book.BookID, &book.Judul, &book.Penulis, &book.Penerbit, &book.TahunTerbit, &book.Sinopsis, &book.ImageURL, &book.CategoryID); err != nil {
			return utils.ErrorResponse(c, fiber.StatusInternalServerError, fmt.Sprintf("Failed to scan book data: %v", err))
		}
		books = append(books, book)
	}

	// Check for any errors during row iteration
	if err = rows.Err(); err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, fmt.Sprintf("Error during rows iteration: %v", err))
	}

	// If no books are found, return an empty array with a success status
	if len(books) == 0 {
		return utils.JSONResponse(c, fiber.StatusOK, "No books found", []models.Book{})
	}

	return utils.JSONResponse(c, fiber.StatusOK, "Books retrieved successfully", books)
}

// GetBookByID gets a single book by its ID
// GET /api/v1/books/:id
func GetBookByID(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id")) // Get ID from URL parameter and convert to int
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid book ID")
	}

	book := new(models.Book)
	// Query a single book by ID
	err = database.DB.QueryRow("SELECT book_id, judul, penulis, penerbit, tahun_terbit, sinopsis, image_url, category_id FROM books WHERE book_id = ?", id).
		Scan(&book.BookID, &book.Judul, &book.Penulis, &book.Penerbit, &book.TahunTerbit, &book.Sinopsis, &book.ImageURL, &book.CategoryID)
	if err != nil {
		// Handle case where book is not found
		if err == sql.ErrNoRows {
			return utils.ErrorResponse(c, fiber.StatusNotFound, "Book not found")
		}
		// Handle other database errors
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, fmt.Sprintf("Failed to retrieve book: %v", err))
	}

	return utils.JSONResponse(c, fiber.StatusOK, "Book retrieved successfully", book)
}

// CreateBook adds a new book to the database
// POST /api/v1/books
func CreateBook(c *fiber.Ctx) error {
	book := new(models.Book)

	// Parse request body into Book struct
	if err := c.BodyParser(book); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body")
	}

	// Basic validation for required fields
	if book.Judul == "" || book.Penulis == "" || book.Penerbit == "" || book.TahunTerbit == 0 || book.CategoryID == 0 {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Judul, Penulis, Penerbit, Tahun Terbit, and Category ID are required")
	}

	// Insert the new book into the database
	result, err := database.DB.Exec(
		"INSERT INTO books (judul, penulis, penerbit, tahun_terbit, sinopsis, image_url, category_id) VALUES (?, ?, ?, ?, ?, ?, ?)",
		book.Judul, book.Penulis, book.Penerbit, book.TahunTerbit, book.Sinopsis, book.ImageURL, book.CategoryID,
	)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, fmt.Sprintf("Failed to create book: %v", err))
	}

	// Get the ID of the newly inserted book
	id, _ := result.LastInsertId()
	book.BookID = int(id) // Set the newly generated ID back to the struct for response

	return utils.JSONResponse(c, fiber.StatusCreated, "Book created successfully", book)
}

// UpdateBook updates an existing book in the database
// PUT /api/v1/books/:id
func UpdateBook(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id")) // Get ID from URL parameter
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid book ID")
	}

	book := new(models.Book)
	// Parse request body for updated book data
	if err := c.BodyParser(book); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body")
	}

	// Basic validation for required fields
	if book.Judul == "" || book.Penulis == "" || book.Penerbit == "" || book.TahunTerbit == 0 || book.CategoryID == 0 {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Judul, Penulis, Penerbit, Tahun Terbit, and Category ID are required")
	}

	// Update the book in the database
	res, err := database.DB.Exec(
		"UPDATE books SET judul = ?, penulis = ?, penerbit = ?, tahun_terbit = ?, sinopsis = ?, image_url = ?, category_id = ? WHERE book_id = ?",
		book.Judul, book.Penulis, book.Penerbit, book.TahunTerbit, book.Sinopsis, book.ImageURL, book.CategoryID, id,
	)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, fmt.Sprintf("Failed to update book: %v", err))
	}

	// Check if any rows were affected (meaning book was found and updated)
	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 0 {
		return utils.ErrorResponse(c, fiber.StatusNotFound, "Book not found or no changes made")
	}

	book.BookID = id // Set the ID back to the struct for response
	return utils.JSONResponse(c, fiber.StatusOK, "Book updated successfully", book)
}

// DeleteBook deletes a book from the database
// DELETE /api/v1/books/:id
func DeleteBook(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id")) // Get ID from URL parameter
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid book ID")
	}

	// Delete the book from the database
	res, err := database.DB.Exec("DELETE FROM books WHERE book_id = ?", id)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, fmt.Sprintf("Failed to delete book: %v", err))
	}

	// Check if any rows were affected (meaning book was found and deleted)
	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 0 {
		return utils.ErrorResponse(c, fiber.StatusNotFound, "Book not found")
	}

	return utils.JSONResponse(c, fiber.StatusOK, "Book deleted successfully", nil) // Return nil data for successful deletion
}