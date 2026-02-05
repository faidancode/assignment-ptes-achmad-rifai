# Go Product Dashboard API 

## Link
```json
https://github.com/faidancode/assignment-ptes-achmad-rifai
```

Project ini adalah implementasi sistem manajemen produk dan dashboard statistik menggunakan **Golang** dengan fokus pada performa tinggi melalui **Redis Caching** dan **SQL Aggregation**.

## Tech Stack
- **Backend:** Go 1.25
- **Database:** MySQL 8.0 (Primary) & Redis 7 (Caching)
- **SQL Tooling:** [SQLC](https://sqlc.dev/) (Type-safe SQL generator)
- **Migration:** [golang-migrate](https://github.com/golang-migrate/migrate)
- **Testing:** Uber GoMock & Testify
- **Containerization:** Docker & Docker Compose

## Performance Optimization
Project ini mengimplementasikan teknik **Cache-Aside** pada endpoint Dashboard. 
- **Database Only:** ~100ms - 250ms (untuk 10.000+ data).
- **With Redis Cache:** ~5ms - 70ms (reduksi waktu hingga >70%).


## Quick Start

### Setup Environment
1. Setup Environment
Salin file `.env.example` menjadi `.env` dan sesuaikan kredensialnya.
```bash
cp .env.example .env
```

2. Jalankan dengan Docker (Paling Cepat)
Jika Anda memiliki Make:

```bash
make docker-up
```

Jika TIDAK memiliki Make:

```bash
docker-compose up -d --build
```

3. Inisialisasi Database
Setelah container berjalan, jalankan migrasi dan seeder untuk mengisi 10.000 data:

Dengan Make:

```bash
make migrate-up
make seed
```

Tanpa Make:

### Migrasi

```bash
migrate -path internal/shared/database/migrations/ -database "mysql://user:password@tcp(localhost:3306)/assignment_ptes" up
```

### Seeder
```bash
go run internal/shared/database/seed/main.go
```

### Testing
Test menggunakan Uber GoMock untuk memastikan logika bisnis (Service Layer) teruji dengan baik tanpa bergantung pada database asli.

Jalankan Test:

```bash
go test ./... -v
```

### Tips Tambahan:

Pastikan variabel `DB_MIGRATION_URL` di `.env` Anda menggunakan format yang didukung `golang-migrate`, contoh:
`mysql://user:password@tcp(localhost:3306)/assignment_ptes`

## Swagger Documentation

Pastikan aplikasi, mysql dan redis sudah berjalan:

```json

http://localhost:3000/swagger/index.html
```

## Design Decisions & Technical Optimizations

### Resilient Server Orchestration

***Server Hardening***: Mengamankan server dengan konfigurasi Read, Write, dan Idle Timeouts yang ketat untuk mencegah resource exhaustion akibat koneksi buruk atau serangan Slowloris.

***Graceful Termination***: Mengimplementasikan mekanisme shutdown yang rapi. Server akan menyelesaikan semua request yang sedang berjalan sebelum mematikan koneksi, memastikan tidak ada data yang korup atau transaksi yang terputus di tengah jalan.

## High-Performance Database Strategy

Untuk menangani dataset besar (ribuan baris), sistem ini menghindari penggunaan ORM berat dan berfokus pada optimasi level SQL:

***UUID v7 (Time-Ordered)***: Berbeda dengan UUID v4 yang acak, UUID v7 bersifat kronologis. Hal ini sangat B-Tree Friendly untuk index MySQL, mencegah page fragmentation dan memastikan operasi INSERT tetap kencang meskipun data sudah mencapai jutaan baris.

***Eliminasi Masalah N+1***: Menggunakan teknik JSON Aggregation (JSON_ARRAYAGG & JSON_OBJECT). Data Order dan Items ditarik dalam satu kali round-trip ke database. Hal ini mengurangi latensi jaringan secara drastis dibanding melakukan looping query di aplikasi.

***Composite & Covering Index***: Menambahkan index strategis seperti idx_orders_customer_id_total_price. Database dapat melakukan Index Seek untuk kalkulasi SUM tanpa perlu membaca data baris secara utuh (Full Table Scan).

## Advanced Concurrency & Caching

Dashboard Report dirancang untuk menangani beban trafik tinggi dengan latensi minimal:

***Parallel Execution (Errgroup)***: Menggunakan errgroup untuk menjalankan query independen secara paralel. Hal ini memangkas waktu respon dari akumulatif ($T1 + T2$) menjadi waktu maksimal dari query terlama ($max(T1, T2)$).

***Cache Stampede Protection (Singleflight)***: Menggunakan golang.org/x/sync/singleflight untuk memastikan jika cache kadaluarsa di tengah trafik tinggi, hanya satu request yang menembak ke database, sementara request lainnya menunggu hasilnya. Ini mencegah "ledakan" beban pada database.

***Manual Cache Invalidation***: Menggunakan strategi active invalidation (menghapus cache saat terjadi mutasi data seperti Create/Update), sehingga dashboard tetap akurat (real-time) tanpa harus menunggu TTL cache habis.