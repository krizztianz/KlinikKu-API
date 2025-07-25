# 📚 Book Library API

API sederhana untuk manajemen buku dan kategori, dengan autentikasi JWT dan dokumentasi Swagger.

## 🚀 Jalankan Proyek

1. Copy `.env.example` ke `.env` dan isi konfigurasi database
2. Jalankan migration:
   ```bash
   go run main.go

    Akses Swagger:

    {http/https}://{Host}:{Port}/swagger/index.html

🔐 Otentikasi JWT

Login ke API:

POST /api/users/login
Content-Type: application/json

{
  "username": "admin",
  "password": "Rahasia123"
}

Response:

{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6..."
}

Gunakan token tersebut di header:

Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6...

🧪 Contoh di Swagger

    Klik tombol 🔒 Authorize

    Masukkan:

    Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6...

    Klik Authorize

📬 Contoh Request via Postman
Login

    Method: POST

    URL: http://localhost:8080/api/users/login

    Body (raw JSON):

{
  "username": "admin",
  "password": "Rahasia123"
}

Get All Books

    Method: GET

    URL: http://localhost:8080/api/books

    Header:

Authorization: Bearer <your_token>