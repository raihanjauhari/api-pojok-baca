package models

import "time"

type User struct {
    UserID     int       `json:"user_id" db:"user_id"`
    NamaLengkap string    `json:"nama_lengkap" db:"nama_lengkap"`
    NIM        string    `json:"nim" db:"nim"`
    Email      string    `json:"email" db:"email"`
    Password   string    `json:"password" db:"password"`
    CreatedAt  time.Time `json:"created_at" db:"created_at"`
    UpdatedAt  time.Time `json:"updated_at" db:"updated_at"`
}

// UserLogin struct for handling login requests
type UserLogin struct {
    Email    string `json:"email" validate:"required"` // UBAH DARI NIM KE EMAIL
    Password string `json:"password" validate:"required"`
}