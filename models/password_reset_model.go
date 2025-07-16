package models

import "time"

// PasswordResetCode represents the 'password_reset_codes' table in the database
type PasswordResetCode struct {
	CodeID    int       `json:"code_id" db:"code_id"`    // Corresponds to code_id in DB
	UserID    int       `json:"user_id" db:"user_id"`    // Foreign key to users
	ResetCode string    `json:"reset_code" db:"reset_code"` // 6-digit reset code
	ExpiresAt time.Time `json:"expires_at" db:"expires_at"` // Expiration timestamp
	CreatedAt time.Time `json:"created_at" db:"created_at"` // Creation timestamp
}