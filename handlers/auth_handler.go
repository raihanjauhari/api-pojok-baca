package handlers

import (
	"database/sql"
	"fmt"
	"log"
	"math/rand"
	"net/smtp"
	"strings"
	"time"

	"pojok_baca_api/database"
	"pojok_baca_api/models"
	"pojok_baca_api/utils"

	"github.com/gofiber/fiber/v2"
)

// Login handles user authentication
func Login(c *fiber.Ctx) error {
	userLogin := new(models.UserLogin)

	if err := c.BodyParser(userLogin); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body")
	}

	userLogin.Email = strings.TrimSpace(userLogin.Email)
	userLogin.Password = strings.TrimSpace(userLogin.Password)

	if userLogin.Email == "" || userLogin.Password == "" {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Email dan Password harus diisi")
	}

	user := new(models.User)
	err := database.DB.QueryRow("SELECT user_id, nama_lengkap, nim, email, password FROM users WHERE email = ?", userLogin.Email).Scan(&user.UserID, &user.NamaLengkap, &user.NIM, &user.Email, &user.Password)
	if err != nil {
		if err == sql.ErrNoRows {
			return utils.ErrorResponse(c, fiber.StatusUnauthorized, "Email atau password salah")
		}
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, fmt.Sprintf("Database error during login attempt: %v", err))
	}

	if userLogin.Password != user.Password {
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "Email atau password salah")
	}

	return utils.JSONResponse(c, fiber.StatusOK, "Login successful", fiber.Map{
		"user_id":      user.UserID,
		"nim":          user.NIM,
		"nama_lengkap": user.NamaLengkap,
		"email":        user.Email,
	})
}

// Register handles new user registration
func Register(c *fiber.Ctx) error {
	user := new(models.User)

	if err := c.BodyParser(user); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body")
	}

	user.NamaLengkap = strings.TrimSpace(user.NamaLengkap)
	user.NIM = strings.TrimSpace(user.NIM)
	user.Email = strings.TrimSpace(user.Email)
	user.Password = strings.TrimSpace(user.Password)

	if user.NIM == "" || user.Email == "" || user.Password == "" || user.NamaLengkap == "" {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Semua kolom (Nama Lengkap, NIM, Email, Password) harus diisi")
	}

	var count int
	database.DB.QueryRow("SELECT COUNT(*) FROM users WHERE email = ? OR nim = ?", user.Email, user.NIM).Scan(&count)
	if count > 0 {
		return utils.ErrorResponse(c, fiber.StatusConflict, "NIM atau Email sudah terdaftar")
	}

	result, err := database.DB.Exec(
		"INSERT INTO users (nama_lengkap, nim, email, password) VALUES (?, ?, ?, ?)",
		user.NamaLengkap, user.NIM, user.Email, user.Password,
	)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, fmt.Sprintf("Gagal mendaftarkan pengguna: %v", err))
	}

	id, _ := result.LastInsertId()
	user.UserID = int(id)
	user.Password = ""

	return utils.JSONResponse(c, fiber.StatusCreated, "Pengguna berhasil didaftarkan", user)
}

// RequestPasswordReset handles request to initiate password reset process
// POST /api/v1/password-reset/request
func RequestPasswordReset(c *fiber.Ctx) error {
	type RequestBody struct {
		Email string `json:"email"`
	}
	req := new(RequestBody)
	if err := c.BodyParser(req); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body")
	}

	req.Email = strings.TrimSpace(req.Email)
	if req.Email == "" {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Email diperlukan")
	}

	user := new(models.User)
	err := database.DB.QueryRow("SELECT user_id FROM users WHERE email = ?", req.Email).Scan(&user.UserID)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Printf("Password reset requested for non-existent email: %s", req.Email)
			return utils.JSONResponse(c, fiber.StatusOK, "Jika email terdaftar, kode reset akan dikirim.", nil)
		}
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, fmt.Sprintf("Database error: %v", err))
	}

	_, err = database.DB.Exec("DELETE FROM password_reset_codes WHERE user_id = ?", user.UserID)
	if err != nil {
		log.Printf("Failed to delete old reset codes for user %d: %v", user.UserID, err)
	}

	rand.Seed(time.Now().UnixNano())
	resetCode := fmt.Sprintf("%06d", rand.Intn(1000000))

	_, err = database.DB.Exec(
		"INSERT INTO password_reset_codes (user_id, reset_code) VALUES (?, ?)",
		user.UserID, resetCode,
	)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, fmt.Sprintf("Failed to save reset code: %v", err))
	}

	const (
		smtpHost     = "smtp.gmail.com"
		smtpPort     = "587"
		senderEmail  = "asrama.putera.murakatahst@gmail.com"
		senderPass   = "yblp ditv bddp dclj" // GANTI INI DENGAN PASSWORD APLIKASI ASLI AKUN GOOGLE ANDA!
	)

	to := []string{req.Email}
	msg := []byte(
		"To: " + req.Email + "\r\n" +
			"From: " + senderEmail + "\r\n" +
			"Subject: Kode Reset Password Pojok Baca\r\n" +
			"Content-Type: text/plain; charset=UTF-8\r\n" +
			"\r\n" +
			"Halo Pengguna Pojok Baca,\n\n" +
			"Anda telah meminta reset password. Berikut adalah kode reset Anda:\n\n" +
			"Kode Reset: " + resetCode + "\n\n" +
			"Kode ini tidak memiliki masa kedaluwarsa, namun hanya dapat digunakan satu kali.\n\n" +
			"Jika Anda tidak meminta reset password ini, harap abaikan email ini.\n\n" +
			"Terima kasih,\nTim Pojok Baca\n",
	)

	auth := smtp.PlainAuth("", senderEmail, senderPass, smtpHost)

	err = smtp.SendMail(smtpHost+":"+smtpPort, auth, senderEmail, to, msg)
	if err != nil {
		log.Printf("Gagal mengirim email ke %s: %v", req.Email, err)
	}

	return utils.JSONResponse(c, fiber.StatusOK, "Jika email terdaftar, kode reset akan dikirim.", nil)
}

// VerifyResetCode handles verification of the reset code
// POST /api/v1/password-reset/verify
func VerifyResetCode(c *fiber.Ctx) error {
	type RequestBody struct {
		Email     string `json:"email"`
		ResetCode string `json:"reset_code"`
	}
	req := new(RequestBody)
	if err := c.BodyParser(req); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body")
	}

	req.Email = strings.TrimSpace(req.Email)
	req.ResetCode = strings.TrimSpace(req.ResetCode)

	if req.Email == "" || req.ResetCode == "" {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Email dan kode reset diperlukan")
	}

	var storedUserID int

	err := database.DB.QueryRow(
		`SELECT prc.user_id FROM password_reset_codes prc
		 JOIN users u ON prc.user_id = u.user_id
		 WHERE u.email = ? AND prc.reset_code = ?`,
		req.Email, req.ResetCode,
	).Scan(&storedUserID)

	if err != nil {
		if err == sql.ErrNoRows {
			return utils.ErrorResponse(c, fiber.StatusBadRequest, "Kode reset tidak valid atau tidak cocok dengan email")
		}
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, fmt.Sprintf("Database error: %v", err))
	}

	_, err = database.DB.Exec("DELETE FROM password_reset_codes WHERE user_id = ?", storedUserID)
	if err != nil {
		log.Printf("Gagal menghapus kode reset yang sudah digunakan untuk user %d: %v", storedUserID, err)
	}

	return utils.JSONResponse(c, fiber.StatusOK, "Kode reset valid", nil)
}

// SetNewPassword handles setting a new password after successful code verification
// POST /api/v1/password-reset/set-new-password
func SetNewPassword(c *fiber.Ctx) error {
	type RequestBody struct {
		Email       string `json:"email"`
		NewPassword string `json:"new_password"`
	}
	req := new(RequestBody)
	if err := c.BodyParser(req); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body")
	}

	req.Email = strings.TrimSpace(req.Email)
	req.NewPassword = strings.TrimSpace(req.NewPassword)

	if req.Email == "" || req.NewPassword == "" {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Email dan password baru diperlukan")
	}

	res, err := database.DB.Exec(
		"UPDATE users SET password = ? WHERE email = ?",
		req.NewPassword,
		req.Email,
	)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, fmt.Sprintf("Failed to update password: %v", err))
	}

	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 0 {
		return utils.ErrorResponse(c, fiber.StatusNotFound, "Pengguna tidak ditemukan")
	}

	return utils.JSONResponse(c, fiber.StatusOK, "Password berhasil diubah", nil)
}