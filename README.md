# 📚 KlinikKU API

API sederhana untuk manajemen operasional Klinik/Puskesmas, dengan autentikasi JWT dan dokumentasi Swagger.

## 🚀 Jalankan Proyek

1. Copy `.env.example` ke `.env` dan isi konfigurasi database
2. Jalankan migration:
   ```bash
   go run main.go

    Akses Swagger:

    {http/https}://{Host}:{Port}/swagger/index.html

🔐 Otentikasi JWT

Login ke API:

POST /api/auth/login
Content-Type: application/json

{
  "username": "admin",
  "password": "password123"
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

    URL: http://localhost:8080/api/auth/login

    Body (raw JSON):

{
  "username": "admin",
  "password": "password123"
}

Get All Kunjungan

    Method: GET

    URL: http://localhost:8080/api/kunjungan

    Header:

Authorization: Bearer <your_token>