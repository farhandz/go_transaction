# Dokumentasi API - Transaction

## Table of Contents
- [Setup Project](#setup-project)
- [Database](#database)
- [Endpoint API](#endpoint-api)
  - [POST /transactions](#post-transactions)
  - [GET /transactions](#get-transactions)
  - [GET /transactions/{id}](#get-transactionsid)
  - [PUT /transactions/{id}](#put-transactionsid)
  - [DELETE /transactions/{id}](#delete-transactionsid)
  - [GET /dashboard/summary](#get-dashboardsummary)
- [Testing](#testing)
- [Bonus](#bonus)
- [Dokumentasi API](#dokumentasi-api)
- [Kriteria Penilaian](#kriteria-penilaian)

## Setup Project

### Persiapan Awal
- clone project  dengan `git clone https://github.com/farhandz/go_transaction`.
- Salin `.env.example` menjadi `.env` dan isi konfigurasi yang dibutuhkan.
- Jalankan aplikasi dengan docker `docker-compose -f docker-compose-dev.yml up --build`.
- Akses aplikasi di [http://0.0.0.0:8000/health](http://0.0.0.0:8000/health) untuk pengecekan status.

### Middlewares
- Use Gin CORSMiddleware
go
router := gin.New()
router.Use(gin.Logger())
router.Use(gin.Recovery())
router.Use(middleware.CORSMiddleware())


### Code Structure
<pre>├── <font color="#3465A4"><b>config</b></font>
│   ├── config.go
│   ├── db.go
│   └── server.go
├── <font color="#3465A4"><b>controllers</b></font>
│   └── example_controller.go
├── docker-compose-dev.yml
├── docker-compose-prod.yml
├── Dockerfile
├── Dockerfile-dev
├── go.mod
├── go.sum
├── <font color="#3465A4"><b>helpers</b></font>
│   ├── response.go
│   └── search.go
├── <font color="#3465A4"><b>infra</b></font>
│   ├── <font color="#3465A4"><b>database</b></font>
│   │   └── database.go
│   └── <font color="#3465A4"><b>logger</b></font>
│       └── logger.go
├── main.go
├── Makefile
├── <font color="#3465A4"><b>migrations</b></font>
│   └── migration.go
├── <font color="#3465A4"><b>models</b></font>
│   └── example_model.go
├── README.md
├── <font color="#3465A4"><b>repository</b></font>
│   └── sql_repo.go
├── <font color="#3465A4"><b>routers</b></font>
│   ├── index.go
│   ├── <font color="#3465A4"><b>middleware</b></font>
│   │   └── cors.go
│   └── router.go
</pre>

## Database

### Skema Database
| Field      | Tipe     | Keterangan                              |
|------------|----------|------------------------------------------|
| id         | int      | Primary Key, auto increment              |
| user_id    | int      | ID pengguna yang melakukan transaksi    |
| amount     | int      | Jumlah transaksi                        |
| status     | string   | Status transaksi (success, pending, failed) |
| created_at | datetime | Waktu transaksi dibuat                  |
| updated_at | datetime | Waktu transaksi diperbarui               |

Gunakan GORM atau `database/sql` untuk migrasi database.

## Endpoint API

### **POST /transactions**
#### Deskripsi
Membuat transaksi baru.

#### Request Body
| Nama    | Tipe  | Wajib | Deskripsi             |
|---------|-------|-------|------------------------|
| user_id | int   | Ya    | ID pengguna           |
| amount  | int   | Ya    | Jumlah transaksi      |
| status  | string| Ya    | Status transaksi      |

#### Response (Positive Case)
```json
{
  "status": "success",
  "message": "transaction created successfully",
  "data": {
    "id": 1,
    "user_id": 1,
    "amount": 1000,
    "status": "success",
    "created_at": "2025-02-19T10:00:00Z",
    "updated_at": "2025-02-19T10:00:00Z"
  }
}
```

### **GET /transactions**
#### Deskripsi
Mengambil daftar transaksi berdasarkan filter tertentu dengan pagination.

#### Parameter Query
| Nama         | Tipe   | Wajib | Deskripsi                               | Contoh  |
|--------------|--------|-------|------------------------------------------|---------|
| page_number  | int    | Tidak | Nomor halaman untuk paginasi            | 1       |
| page_size    | int    | Tidak | Jumlah data per halaman                  | 5       |
| status       | string | Tidak | Filter berdasarkan status transaksi      | pending |
| user_id      | int    | Tidak | Filter berdasarkan ID pengguna           | 1       |

#### Response (Positive Case)
```json
{
  "status": "success",
  "message": "success get data",
  "data": {
    "data": [
      {
        "id": 3,
        "user_id": 0,
        "amount": 0,
        "status": "",
        "created_at": "2025-02-19T09:38:18.691092+06:00",
        "updated_at": "2025-02-19T09:38:18.691092+06:00"
      }
    ],
    "page_number": 1,
    "page_size": 5,
    "total_record_count": 5
  }
}
```

### **GET /transactions/{id}**
#### Deskripsi
Mengambil data transaksi berdasarkan ID.

#### Parameter Path
| Nama | Tipe | Wajib | Deskripsi           | Contoh |
|------|------|-------|----------------------|--------|
| id   | int  | Ya    | ID transaksi yang dicari | 2      |

#### Response (Positive Case)
```json
{
  "status": "success",
  "message": "success get by id",
  "data": {
    "id": 2,
    "user_id": 0,
    "amount": 0,
    "status": "success",
    "created_at": "2025-02-19T09:36:58.386317+06:00",
    "updated_at": "2025-02-19T13:58:18.079946+06:00"
  }
}
```

### **PUT /transactions/{id}**
#### Deskripsi
Mengupdate status transaksi berdasarkan ID.

#### Request Body
| Nama    | Tipe  | Wajib | Deskripsi             |
|---------|-------|-------|------------------------|
| status  | string| Ya    | Status transaksi baru |

### **DELETE /transactions/{id}**
#### Deskripsi
Menghapus transaksi berdasarkan ID.

### **GET /dashboard/summary**
#### Deskripsi
Mengambil ringkasan data transaksi untuk dashboard.

#### Response (Positive Case)
```json
{
  "status": "success",
  "message": "dashboard summary",
  "data": {
    "total_success_today": 100,
    "average_transaction_per_user": 5000,
    "latest_transactions": [
      { "id": 1, "user_id": 1, "amount": 1000, "status": "success", "created_at": "2025-02-19T10:00:00Z" }
    ]
  }
}
```

## Testing (86 coveragge)
- Gunakan library `testing` bawaan Go atau `testify` untuk unit test.
- jalankan perinta `docker exec -it dev_go_server sh`
- kemudian `go test ./controllers -cover`
- `ok      gin-boilerplate/controllers     (cached)        coverage: 86.7% of statements` result 

## Bonus
- Logging:  `logrus`
- Error Handling: Validasi input, data tidak ditemukan, dll.
- Pagination: Sudah didukung di endpoint `GET /transactions`.

## Dokumentasi API
- Postman collection (ada di sourcecode).

## Kriteria Penilaian
- Kualitas Kode: Struktur rapi, best practices Golang.
- Fungsionalitas: Semua endpoint berfungsi sesuai spesifikasi.
- Testing: Unit test mencakup fungsi utama dengan 86% coverage.
- Dokumentasi: README.md jelas dan lengkap, API dokumentasi tersedia.

