package routes

import (
	"pojok_baca_api/handlers"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func SetupRoutes(app *fiber.App) {
	app.Use(logger.New())

	// ===================================================================
	// PENTING: Middleware CORS harus diatur sebelum rute apapun,
	// termasuk sebelum app.Static()
	// ===================================================================
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*", // Ini mengizinkan semua origin. Untuk produksi, ubah ke origin spesifik.
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
		AllowMethods: "GET, POST, HEAD, PUT, DELETE, PATCH",
	}))

	// ===================================================================
	// PENTING: Tambahkan baris ini untuk melayani file statis (gambar)
	// Pastikan ini ada dan setelah CORS middleware.
	// ===================================================================
	app.Static("/public", "./public") // <--- PASTI BARIS INI ADA!

	api := app.Group("/api/v1")

	// --- Authentication Routes ---
	api.Post("/login", handlers.Login)
	api.Post("/register", handlers.Register)

	// Route untuk fungsionalitas Lupa Password
	api.Post("/password-reset/request", handlers.RequestPasswordReset)
	api.Post("/password-reset/verify", handlers.VerifyResetCode)
	api.Post("/password-reset/set-new-password", handlers.SetNewPassword)

	// --- Book Routes (CRUD) ---
	api.Get("/books", handlers.GetAllBooks)
	api.Get("/books/:id", handlers.GetBookByID)
	api.Post("/books", handlers.CreateBook)
	api.Put("/books/:id", handlers.UpdateBook)
	api.Delete("/books/:id", handlers.DeleteBook)

	// --- Category Routes (CRUD) ---
	api.Get("/categories", handlers.GetAllCategories)
	api.Get("/categories/:id", handlers.GetCategoryByID)
	api.Post("/categories", handlers.CreateCategory)
	api.Put("/categories/:id", handlers.UpdateCategory)
	api.Delete("/categories/:id", handlers.DeleteCategory)
}