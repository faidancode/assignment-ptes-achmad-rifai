package main

import (
	"database/sql"
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
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

	fmt.Println("🚀 Memulai proses seeding lengkap...")
	start := time.Now()

	// --- 1. SEED CATEGORIES ---
	categoryNames := []string{"Elektronik", "Pakaian", "Kesehatan", "Hobi", "Otomotif"}
	var categoryIDs []string
	for _, name := range categoryNames {
		newID := uuid.New().String()
		db.Exec("INSERT INTO categories (id, name, description) VALUES (?, ?, ?)", newID, name, "Deskripsi "+name)
		categoryIDs = append(categoryIDs, newID)
	}

	// --- 2. SEED PRODUCTS (10.000 data) ---
	var productIDs []string
	tx, _ := db.Begin()
	pStmt, _ := tx.Prepare(`INSERT INTO products (id, name, description, price, category_id, stock_quantity, is_active, created_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?)`)

	for i := 1; i <= 10000; i++ {
		pID := uuid.New().String()
		productIDs = append(productIDs, pID)
		price := decimal.NewFromFloat(rand.Float64() * 500000)
		pStmt.Exec(pID, fmt.Sprintf("Produk %d", i), "Desc", price, categoryIDs[rand.Intn(len(categoryIDs))], rand.Int31n(500), true, time.Now())
	}
	tx.Commit()
	fmt.Println("✅ 10.000 Produk Berhasil")

	// --- 3. SEED CUSTOMERS (500 data) ---
	var customerIDs []string
	cTx, _ := db.Begin()
	cStmt, _ := cTx.Prepare(`INSERT INTO customers (id, name, email) VALUES (?, ?, ?)`)

	for i := 1; i <= 500; i++ {
		cID := uuid.New().String()
		customerIDs = append(customerIDs, cID)
		cStmt.Exec(cID, fmt.Sprintf("Customer %d", i), fmt.Sprintf("user%d@example.com", i))
	}
	cTx.Commit()
	fmt.Println("✅ 500 Customers Berhasil")

	// --- 4. SEED ORDERS & ORDER ITEMS (1.000 Orders) ---
	oTx, _ := db.Begin()
	oStmt, _ := oTx.Prepare(`INSERT INTO orders (id, customer_id, total_quantity, total_price, created_at) VALUES (?, ?, ?, ?, ?)`)
	oiStmt, _ := oTx.Prepare(`INSERT INTO order_items (id, order_id, product_id, quantity, unit_price) VALUES (?, ?, ?, ?, ?)`)

	for i := 1; i <= 1000; i++ {
		orderID := uuid.New().String()
		custID := customerIDs[rand.Intn(len(customerIDs))]

		var totalQty int
		var totalPrice decimal.Decimal

		// INSERT ORDER DULU
		_, err := oStmt.Exec(orderID, custID, 0, 0, time.Now())
		if err != nil {
			log.Fatal("insert order gagal:", err)
		}

		itemsInOrder := rand.Intn(3) + 1
		for j := 0; j < itemsInOrder; j++ {
			pID := productIDs[rand.Intn(len(productIDs))]
			qty := rand.Intn(5) + 1
			uPrice := decimal.NewFromFloat(rand.Float64() * 200000)

			_, err := oiStmt.Exec(
				uuid.New().String(),
				orderID,
				pID,
				qty,
				uPrice,
			)
			if err != nil {
				log.Fatal("insert order item gagal:", err)
			}

			totalQty += qty
			totalPrice = totalPrice.Add(uPrice.Mul(decimal.NewFromInt(int64(qty))))
		}

		// UPDATE TOTAL SETELAH ITEMS
		_, err = oTx.Exec(
			`UPDATE orders SET total_quantity=?, total_price=? WHERE id=?`,
			totalQty,
			totalPrice,
			orderID,
		)
		if err != nil {
			log.Fatal("update order gagal:", err)
		}
	}

	oTx.Commit()
	fmt.Println("✅ 1.000 Orders & Items Berhasil")

	fmt.Printf("✨ SELESAI! Total waktu: %v\n", time.Since(start))
}
