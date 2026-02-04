// main.go
package main

import (
	"assignment-ptes-achmad-rifai/internal/bootstrap"
	"assignment-ptes-achmad-rifai/internal/category"
	"assignment-ptes-achmad-rifai/internal/customer"
	"assignment-ptes-achmad-rifai/internal/dashboard"
	"assignment-ptes-achmad-rifai/internal/order"
	"assignment-ptes-achmad-rifai/internal/product"
	"assignment-ptes-achmad-rifai/internal/shared/database/dbgen"
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql" // Driver diganti ke MySQL
	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"

	_ "assignment-ptes-achmad-rifai/docs"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// ControllerRegistry untuk mengelompokkan handler
type ControllerRegistry struct {
	Category  *category.Handler
	Product   *product.Handler
	Customer  *customer.Handler
	Order     *order.Handler
	Dashboard *dashboard.Handler
}

func connectDBWithRetry(dsn string, maxRetries int) (*sql.DB, error) {
	var db *sql.DB
	var err error

	for i := 1; i <= maxRetries; i++ {
		db, err = sql.Open("mysql", dsn)
		if err == nil {
			err = db.Ping()
			if err == nil {
				log.Println("✅ Successfully connected to MySQL database")
				return db, nil
			}
		}

		log.Printf("⚠️  MySQL connection attempt %d/%d failed: %v", i, maxRetries, err)

		if i < maxRetries {
			time.Sleep(time.Second * 5)
		}
	}

	return nil, err
}

func connectRedisWithRetry(addr string, maxRetries int) (*redis.Client, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: "",
		DB:       0,
	})

	for i := 1; i <= maxRetries; i++ {
		ctx := context.Background()
		_, err := rdb.Ping(ctx).Result()
		if err == nil {
			log.Println("✅ Successfully connected to Redis")
			return rdb, nil
		}

		log.Printf("⚠️  Redis connection attempt %d/%d failed: %v", i, maxRetries, err)

		if i < maxRetries {
			time.Sleep(time.Second * 5)
		}
	}

	return nil, fmt.Errorf("failed to connect to Redis after %d attempts", maxRetries)
}

// @title           Assignment PTES API
// @version         1.0
// @description     API Server for Order and Dashboard Management.
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:3000
// @BasePath  /api/v1
func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found")
	}

	// Connect to MySQL with retry (max 10x, timeout 50s)
	db, err := connectDBWithRetry(os.Getenv("DB_URL"), 10)
	if err != nil {
		log.Fatal("❌ Cannot connect to database after retries:", err)
	}
	defer db.Close()

	queries := dbgen.New(db)

	// Connect to Redis with retry (max 10x, timeout 50s)
	rdb, err := connectRedisWithRetry(os.Getenv("REDIS_URL"), 10)
	if err != nil {
		log.Fatal("❌ Cannot connect to Redis after retries:", err)
	}
	defer rdb.Close()

	log.Println("🚀 All services connected successfully, starting server...")

	// Dependency Injection (DI)

	categoryRepo := category.NewRepository(queries)
	categoryService := category.NewService(categoryRepo)
	categoryHandler := category.NewHandler(categoryService)

	productRepo := product.NewRepository(queries)
	productService := product.NewService(productRepo, rdb)
	productHandler := product.NewHandler(productService)

	customerRepo := customer.NewRepository(queries)
	customerService := customer.NewService(customerRepo)
	customerHandler := customer.NewHandler(customerService)

	orderRepo := order.NewRepository(queries)
	orderService := order.NewService(db, orderRepo)
	orderHandler := order.NewHandler(orderService)

	dashboardRepo := dashboard.NewRepository(queries)
	dashboardService := dashboard.NewService(dashboardRepo, rdb)
	dashboardHandler := dashboard.NewHandler(dashboardService)

	registry := ControllerRegistry{
		Category:  categoryHandler,
		Product:   productHandler,
		Customer:  customerHandler,
		Order:     orderHandler,
		Dashboard: dashboardHandler,
	}

	// Router Setup
	r := gin.Default()

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	// API Grouping
	api := r.Group("/api/v1")
	{
		category.RegisterRoutes(api, registry.Category)
		product.RegisterRoutes(api, registry.Product)
		customer.RegisterRoutes(api, registry.Customer)
		order.RegisterRoutes(api, registry.Order)
		dashboard.RegisterRoutes(api, registry.Dashboard)
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
