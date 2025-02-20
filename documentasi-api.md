# Dokumentasi API - Transaction

## Endpoint
**GET /transaction**

## Deskripsi
Mengambil data transaksi berdasarkan filter tertentu dengan paginasi.

## Parameter Query
| Nama         | Tipe   | Wajib | Deskripsi                               | Contoh  |
|--------------|--------|-------|------------------------------------------|---------|
| page_number  | int    | Ya    | Nomor halaman untuk paginasi            | 1       |
| page_size    | int    | Ya    | Jumlah data per halaman                  | 5       |
| status       | string | Tidak | Filter berdasarkan status transaksi      | pending |
| user_id      | int    | Tidak | Filter berdasarkan ID pengguna           | 1       |

## Response (Positive Case)
| Field                     | Tipe    | Deskripsi                                                        |
|----------------------------|---------|-------------------------------------------------------------------|
| status                     | string  | Status response (success atau error)                             |
| message                    | string  | Pesan deskriptif                                                 |
| data.data[]                | array   | Daftar transaksi                                                 |
| data.data[].id             | int     | ID transaksi                                                     |
| data.data[].user_id        | int     | ID pengguna                                                      |
| data.data[].amount         | int     | Jumlah nominal transaksi                                         |
| data.data[].status         | string  | Status transaksi (pending, success, atau lainnya)                |
| data.data[].created_at     | string  | Waktu pembuatan transaksi (format ISO 8601)                     |
| data.data[].updated_at     | string  | Waktu pembaruan transaksi (format ISO 8601)                     |
| data.page_number           | int     | Nomor halaman                                                    |
| data.page_size             | int     | Jumlah data per halaman                                           |
| data.total_record_count    | int     | Total jumlah data yang tersedia                                   |

### Contoh Response (Berhasil):
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

## Response (Negative Case)

| Skenario Kasus Negatif                                        | HTTP Status    | Response Status | Response Message                                 |
|---------------------------------------------------------------|----------------|-----------------|--------------------------------------------------|
| page_number atau page_size tidak diisi                         | 400 Bad Request| error           | page_number and page_size required               |
| page_number atau page_size bukan angka                         | 400 Bad Request| error           | page_number and page_size must be integer        |
| user_id bukan angka                                            | 400 Bad Request| error           | user_id must be integer                          |
| status mengandung karakter spesial yang tidak valid            | 400 Bad Request| error           | invalid status value                             |
| Data transaksi tidak ditemukan (contoh: user_id tidak ada)     | 200 OK         | success         | data not found dengan data array kosong          |

### Contoh Response (Gagal):
```json
{
  "status": "error",
  "message": "page_number and page_size required"
}
```
atau
```json
{
  "status": "success",
  "message": "data not found",
  "data": {
    "data": [],
    "page_number": 1,
    "page_size": 5,
    "total_record_count": 0
  }
}
```

## Endpoint
**GET /transaction/{id}**

## Deskripsi
Mengambil data transaksi berdasarkan ID.

## Parameter Path
| Nama | Tipe | Wajib | Deskripsi           | Contoh |
|------|------|-------|----------------------|--------|
| id   | int  | Ya    | ID transaksi yang dicari | 2      |

## Response (Positive Case)
| Field          | Tipe   | Deskripsi                        |
|----------------|--------|-----------------------------------|
| status         | string | Status response (success atau error) |
| message        | string | Pesan deskriptif                  |
| data.id        | int    | ID transaksi                      |
| data.user_id   | int    | ID pengguna                       |
| data.amount    | int    | Jumlah nominal transaksi          |
| data.status    | string | Status transaksi                  |
| data.created_at| string | Waktu pembuatan transaksi         |
| data.updated_at| string | Waktu pembaruan transaksi         |

### Contoh Response (Berhasil):
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

## Response (Negative Case)

| Skenario Kasus Negatif       | HTTP Status    | Response Status | Response Message |
|-------------------------------|----------------|-----------------|------------------|
| ID tidak ditemukan            | 404 Not Found  | error           | transaction not found |
| ID bukan angka                | 400 Bad Request| error           | invalid id format |

### Contoh Response (Gagal):
```json
{
  "status": "error",
  "message": "transaction not found",
  "data": null
}
```

## Endpoint
**DELETE /transaction/{id}**

## Deskripsi
Menghapus transaksi berdasarkan ID.

## Response (Positive Case)
```json
{
  "status": "success",
  "message": "Success delete",
  "data": null
}
```

## Response (Negative Case)
```json
{
  "status": "error",
  "message": "Transaction not found",
  "data": null
}
```

## Endpoint
**GET /dashboard/summary**

## Deskripsi
Mengambil ringkasan dashboard transaksi.

## Response (Positive Case)
```json
{
  "status": "success",
  "message": "Success get summary",
  "data": {
    "total_transactions_today": 0,
    "average_transaction_per_user": 2.5,
    "total_transactions": 5,
    "unique_users": 2,
    "total_pending_transactions": 0,
    "total_success_transactions": 2,
    "total_failed_transactions": 0
  }
}
```

## Endpoint
**GET /dashboard/report**

## Deskripsi
Mengambil laporan dashboard transaksi yang mencakup:
- Total transaksi sukses hari ini
- Rata-rata jumlah transaksi per user
- Daftar 10 transaksi terbaru

## Response (Positive Case)
```json
{
  "status": "success",
  "message": "success create",
  "data": {
    "average_transaction_per_user": 2.5,
    "latest_transactions": [...]

