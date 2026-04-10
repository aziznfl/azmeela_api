package main

import (
	"log"
	"time"
	_ "time/tzdata"

	"github.com/gin-gonic/gin"

	"github.com/azmeela/sispeg-api/internal/config"
	httpDelivery "github.com/azmeela/sispeg-api/internal/delivery/http"
	"github.com/azmeela/sispeg-api/internal/delivery/http/middleware"
	"github.com/azmeela/sispeg-api/internal/domain"
	"github.com/azmeela/sispeg-api/pkg/cache"
	"github.com/azmeela/sispeg-api/pkg/database"
	"github.com/azmeela/sispeg-api/pkg/logger"
	"github.com/azmeela/sispeg-api/pkg/token"
)

// @title Sispeg API
// @version 1.0
// @description This is the core API for Sispeg Project in Azmeela
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /api/v1
func main() {
	// Initialize Configuration
	cfg, err := config.LoadConfig(".")
	if err != nil {
		log.Fatalf("Cannot load config: %v", err)
	}

	// Initialize Logger
	logger.InitLogger()
	logger.Log.Info("Starting Azmeela Internal API")

	// Set Timezone to Asia/Jakarta
	loc, err := time.LoadLocation("Asia/Jakarta")
	if err != nil {
		logger.Log.Error("Failed to load location Asia/Jakarta, using system default")
	} else {
		time.Local = loc
	}

	// Initialize Postgres Database Connection
	db, err := database.NewPostgresConn(cfg)
	if err != nil {
		logger.Log.Fatal("Failed to connect database")
	}

	// Migrations
	// Pre-migration: Create enum type if not exists (Postgres specific)
	db.Exec("DO $$ BEGIN IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'membership_status_type') THEN CREATE TYPE membership_status_type AS ENUM ('0', '1', '2'); END IF; END $$;")

	db.AutoMigrate(
		&domain.AdminType{},
		&domain.Employee{},
		&domain.CustomerType{}, 
		&domain.Customer{}, 
		&domain.TransactionStatus{}, 
		&domain.Transaction{},
		&domain.ProductType{},
		&domain.ProductSize{},
		&domain.ProductCode{},
		&domain.Product{},
		&domain.ProductPrice{},
		&domain.ProductStockLog{},
		&domain.Attendance{},
		&domain.Holiday{},
		&domain.Leave{},
		&domain.Overtime{},
		&domain.CashAdvance{},
		&domain.CashAdvanceHistory{},
		&domain.SalaryVariable{},
	)

	// Initialize Redis Connection
	rdb := cache.NewRedisClient(cfg)

	// Setup TokenMaker
	tokenMaker, err := token.NewJWTMaker(cfg)
	if err != nil {
		logger.Log.Fatal("Failed to setup JWT maker")
	}

	// Setup Gin Engine
	if cfg.GinMode == "release" {
		gin.SetMode(gin.ReleaseMode)
	}
	r := gin.Default()
	r.RedirectTrailingSlash = false
	r.RedirectFixedPath = false

	// CORS Middleware (allow Vue frontend on different port)
	r.Use(middleware.CORSMiddleware())
	r.Use(middleware.RecoveryMiddleware())

	// Rate Limiting Middleware
	r.Use(middleware.RateLimiter())

	// Define API routes
	api := r.Group("/api/v1")
	{
		httpDelivery.SetupRouter(api, db, rdb, tokenMaker)
	}

	// Start server
	port := cfg.ServerPort
	if port == "" {
		port = "8080"
	}
	logger.Log.Info("Listening and serving HTTP on :" + port)
	if err := r.Run(":" + port); err != nil {
		logger.Log.Fatal("Failed to start server: " + err.Error())
	}
}
