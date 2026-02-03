// main.go
package main

import (
	"assignment-ptes-achmad-rifai/internal/bootstrap"
	"assignment-ptes-achmad-rifai/internal/category"
	"assignment-ptes-achmad-rifai/internal/customer"
	"assignment-ptes-achmad-rifai/internal/order"
	"assignment-ptes-achmad-rifai/internal/product"
	"assignment-ptes-achmad-rifai/internal/shared/database/dbgen"
	"database/sql"
	"log"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql" // Driver diganti ke MySQL
	"github.com/joho/godotenv"
)

// ControllerRegistry untuk mengelompokkan handler
type ControllerRegistry struct {
	Category *category.Handler
	Product  *product.Handler
	Customer *customer.Handler
	Order    *order.Handler
}

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found")
	}

	db, err := sql.Open("mysql", os.Getenv("DB_URL"))
	if err != nil {
		log.Fatal("Cannot connect to database:", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatal("Database unreachable:", err)
	}

	queries := dbgen.New(db)

	// Dependency Injection (DI)

	categoryRepo := category.NewRepository(queries)
	categoryService := category.NewService(categoryRepo)
	categoryHandler := category.NewHandler(categoryService)

	productRepo := product.NewRepository(queries)
	productService := product.NewService(productRepo)
	productHandler := product.NewHandler(productService)

	customerRepo := customer.NewRepository(queries)
	customerService := customer.NewService(customerRepo)
	customerHandler := customer.NewHandler(customerService)

	orderRepo := order.NewRepository(queries)
	orderService := order.NewService(db, orderRepo)
	orderHandler := order.NewHandler(orderService)

	registry := ControllerRegistry{
		Category: categoryHandler,
		Product:  productHandler,
		Customer: customerHandler,
		Order:    orderHandler,
	}

	// Router Setup
	r := gin.Default()

	// API Grouping
	api := r.Group("/api/v1")
	{
		category.RegisterRoutes(api, registry.Category)
		product.RegisterRoutes(api, registry.Product)
		customer.RegisterRoutes(api, registry.Customer)
		order.RegisterRoutes(api, registry.Order)
	}

	// Audit logger & Server Config
	auditLogger := bootstrap.NewStdoutAuditLogger()

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	// Server Hardening and Graceful Management
	bootstrap.StartHTTPServer(
		r,
		bootstrap.ServerConfig{
			Port:         port,
			ReadTimeout:  5 * time.Second,
			WriteTimeout: 10 * time.Second,
			IdleTimeout:  60 * time.Second,
		},
		auditLogger,
	)
}
