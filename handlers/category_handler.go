package handlers

import (
	"database/sql"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"pojok_baca_api/database"
	"pojok_baca_api/models"
	"pojok_baca_api/utils"
	"strconv"
)

// GetAllCategories gets all categories from the database
// GET /api/v1/categories
func GetAllCategories(c *fiber.Ctx) error {
	rows, err := database.DB.Query("SELECT category_id, nama_kategori, image_url FROM categories")
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, fmt.Sprintf("Failed to retrieve categories: %v", err))
	}
	defer rows.Close()

	var categories []models.Category
	for rows.Next() {
		var category models.Category
		// Scan into category struct, including the image_url field
		if err := rows.Scan(&category.CategoryID, &category.NamaKategori, &category.ImageURL); err != nil {
			return utils.ErrorResponse(c, fiber.StatusInternalServerError, fmt.Sprintf("Failed to scan category data: %v", err))
		}
		categories = append(categories, category)
	}

	if err = rows.Err(); err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, fmt.Sprintf("Error during rows iteration: %v", err))
	}

	if len(categories) == 0 {
		return utils.JSONResponse(c, fiber.StatusOK, "No categories found", []models.Category{})
	}

	return utils.JSONResponse(c, fiber.StatusOK, "Categories retrieved successfully", categories)
}

// GetCategoryByID gets a single category by its ID
// GET /api/v1/categories/:id
func GetCategoryByID(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid category ID")
	}

	category := new(models.Category)
	err = database.DB.QueryRow("SELECT category_id, nama_kategori, image_url FROM categories WHERE category_id = ?", id).
		Scan(&category.CategoryID, &category.NamaKategori, &category.ImageURL)
	if err != nil {
		if err == sql.ErrNoRows {
			return utils.ErrorResponse(c, fiber.StatusNotFound, "Category not found")
		}
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, fmt.Sprintf("Failed to retrieve category: %v", err))
	}

	return utils.JSONResponse(c, fiber.StatusOK, "Category retrieved successfully", category)
}

// CreateCategory adds a new category to the database
// POST /api/v1/categories
func CreateCategory(c *fiber.Ctx) error {
	category := new(models.Category)
	if err := c.BodyParser(category); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body")
	}

	if category.NamaKategori == "" {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Category name is required")
	}

	result, err := database.DB.Exec(
		"INSERT INTO categories (nama_kategori, image_url) VALUES (?, ?)",
		category.NamaKategori, category.ImageURL,
	)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, fmt.Sprintf("Failed to create category: %v", err))
	}

	id, _ := result.LastInsertId()
	category.CategoryID = int(id)

	return utils.JSONResponse(c, fiber.StatusCreated, "Category created successfully", category)
}

// UpdateCategory updates an existing category in the database
// PUT /api/v1/categories/:id
func UpdateCategory(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid category ID")
	}

	category := new(models.Category)
	if err := c.BodyParser(category); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body")
	}

	if category.NamaKategori == "" {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Category name is required")
	}

	res, err := database.DB.Exec(
		"UPDATE categories SET nama_kategori = ?, image_url = ? WHERE category_id = ?",
		category.NamaKategori, category.ImageURL, id,
	)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, fmt.Sprintf("Failed to update category: %v", err))
	}

	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 0 {
		return utils.ErrorResponse(c, fiber.StatusNotFound, "Category not found or no changes made")
	}

	category.CategoryID = id // Set the ID back for response
	return utils.JSONResponse(c, fiber.StatusOK, "Category updated successfully", category)
}

// DeleteCategory deletes a category from the database
// DELETE /api/v1/categories/:id
func DeleteCategory(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid category ID")
	}

	res, err := database.DB.Exec("DELETE FROM categories WHERE category_id = ?", id)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, fmt.Sprintf("Failed to delete category: %v", err))
	}

	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 0 {
		return utils.ErrorResponse(c, fiber.StatusNotFound, "Category not found")
	}

	return utils.JSONResponse(c, fiber.StatusOK, "Category deleted successfully", nil)
}