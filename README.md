# Go Product Dashboard API 🚀

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
