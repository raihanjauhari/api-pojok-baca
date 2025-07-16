package database

import (
	"database/sql" // Package untuk interaksi SQL generik
	"fmt"          // Package untuk format string
	"log"          // Package untuk logging error dan informasi
	"os"           // Package untuk mengakses variabel lingkungan

	_ "github.com/go-sql-driver/mysql" // Driver MySQL untuk MariaDB. Gunakan underscore (_) karena kita hanya mengimpor efek samping (register driver)
)

// DB adalah variabel global yang akan menyimpan koneksi database.
// Ini bisa diakses dari package lain yang mengimpor package database.
var DB *sql.DB

// ConnectDB bertanggung jawab untuk menginisialisasi dan membuka koneksi ke database.
func ConnectDB() {
	var err error

	// Membuat DSN (Data Source Name) string dari variabel lingkungan.
	// Format DSN untuk MySQL/MariaDB: "user:password@tcp(host:port)/dbname?param=value"
	// os.Getenv() digunakan untuk mengambil nilai dari file .env yang sudah Anda buat.
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		os.Getenv("DB_USER"),     // Mengambil nama pengguna DB dari .env
		os.Getenv("DB_PASSWORD"), // Mengambil password DB dari .env (bisa kosong jika tidak ada password)
		os.Getenv("DB_HOST"),     // Mengambil host DB dari .env
		os.Getenv("DB_PORT"),     // Mengambil port DB dari .env
		os.Getenv("DB_NAME"),     // Mengambil nama database dari .env
	)

	// Membuka koneksi database menggunakan driver "mysql".
	// Jika ada error, akan dicatat dan aplikasi akan berhenti.
	DB, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("Gagal membuka koneksi database: %v", err)
	}

	// Memverifikasi koneksi ke database dengan mengirimkan ping.
	// Jika gagal ping, berarti koneksi tidak berhasil atau ada masalah jaringan/DB.
	if err = DB.Ping(); err != nil {
		log.Fatalf("Gagal melakukan ping ke database: %v", err)
	}

	log.Println("Berhasil terhubung ke MariaDB! 🥳")
}

// CloseDB bertanggung jawab untuk menutup koneksi database ketika aplikasi berhenti.
// Ini penting untuk membebaskan sumber daya database.
func CloseDB() {
	if DB != nil {
		err := DB.Close()
		if err != nil {
			log.Printf("Error saat menutup koneksi database: %v", err)
		} else {
			log.Println("Koneksi database ditutup. 👋")
		}
	}
}