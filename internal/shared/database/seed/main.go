package main

import (
	"database/sql"
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/google/uuid" // Menggunakan Google UUID
	"github.com/joho/godotenv"
	"github.com/shopspring/decimal"
)

func main() {
	_ = godotenv.Load()
	db, err := sql.Open("mysql", os.Getenv("DB_URL"))
	if err != nil {
		log.Fatal("Gagal koneksi DB:", err)
	}
	defer db.Close()

	fmt.Println("🚀 Memulai proses seeding dengan UUID...")
	start := time.Now()

	// 1. Definisikan Nama Kategori
	categoryNames := []string{"Elektronik", "Pakaian", "Kesehatan", "Hobi", "Otomotif"}

	// Slice untuk menampung UUID kategori yang baru dibuat
	var categoryIDs []string

	fmt.Println("📦 Mengisi data kategori dengan UUID...")
	for _, name := range categoryNames {
		// Generate UUID menggunakan google/uuid
		newID := uuid.New().String()

		_, err := db.Exec("INSERT INTO categories (id, name, description) VALUES (?, ?, ?)",
			newID, name, "Deskripsi kategori "+name)

		if err != nil {
			log.Printf("Gagal insert kategori %s: %v", name, err)
			continue
		}
		categoryIDs = append(categoryIDs, newID)
	}

	// 2. Gunakan Transaction untuk 10.000 Produk
	tx, err := db.Begin()
	if err != nil {
		log.Fatal(err)
	}

	stmt, err := tx.Prepare(`INSERT INTO products 
		(id, name, description, price, category_id, stock_quantity, is_active, created_at) 
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)`)
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	fmt.Println("🛒 Mengisi 10.000 data produk...")
	for i := 1; i <= 10000; i++ {
		productID := uuid.New().String()
		name := fmt.Sprintf("Produk Performa Test %d", i)
		desc := "Deskripsi produk skala besar"
		price := decimal.NewFromFloat(rand.Float64() * 1000000)
		stock := rand.Int31n(500)

		// Pilih ID kategori secara acak dari slice UUID yang tadi dibuat
		catID := categoryIDs[rand.Intn(len(categoryIDs))]

		_, err = stmt.Exec(productID, name, desc, price, catID, stock, true, time.Now())
		if err != nil {
			log.Printf("Gagal insert produk ke-%d: %v", i, err)
			continue
		}

		if i%2500 == 0 {
			fmt.Printf("✅ %d data produk berhasil dimasukkan...\n", i)
		}
	}

	if err := tx.Commit(); err != nil {
		log.Fatal("Gagal commit transaksi:", err)
	}

	fmt.Printf("✨ SELESAI! Data dengan UUID berhasil di-generate.\n")
	fmt.Printf("⏱️ Total waktu: %v\n", time.Since(start))
}
