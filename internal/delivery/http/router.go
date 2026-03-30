package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"

	"github.com/azmeela/sispeg-api/internal/delivery/http/handler"
	"github.com/azmeela/sispeg-api/internal/delivery/http/middleware"
	"github.com/azmeela/sispeg-api/internal/repository"
	"github.com/azmeela/sispeg-api/internal/usecase"
	"github.com/azmeela/sispeg-api/pkg/token"
)

// SetupRouter will handle all routes initialization
func SetupRouter(r *gin.RouterGroup, db *gorm.DB, rdb *redis.Client, tokenMaker token.TokenMaker) {

	// Health Check
	r.GET("/health", HealthCheck)

	// Initialize Repositories
	employeeRepo := repository.NewEmployeeRepository(db)
	redisRepo := repository.NewRedisRepository(rdb)
	attendanceRepo := repository.NewAttendanceRepository(db)
	leaveRepo := repository.NewLeaveRepository(db)
	overtimeRepo := repository.NewOvertimeRepository(db)
	cashAdvanceRepo := repository.NewCashAdvanceRepository(db)
	holidayRepo := repository.NewHolidayRepository(db)
	salaryVarRepo := repository.NewSalaryVariableRepository(db)
	adminTypeRepo := repository.NewAdminTypeRepository(db)
	reportRepo := repository.NewReportRepository(db)
	customerRepo := repository.NewCustomerRepository(db)
	transactionRepo := repository.NewTransactionRepository(db)
	productRepo := repository.NewProductRepository(db)

	// Initialize Usecases
	employeeUcase := usecase.NewEmployeeUsecase(employeeRepo)
	authUcase := usecase.NewAuthUsecase(employeeRepo, redisRepo, tokenMaker)
	attendanceUcase := usecase.NewAttendanceUsecase(attendanceRepo)
	leaveUcase := usecase.NewLeaveUsecase(leaveRepo)
	overtimeUcase := usecase.NewOvertimeUsecase(overtimeRepo)
	cashAdvanceUcase := usecase.NewCashAdvanceUsecase(cashAdvanceRepo)
	holidayUcase := usecase.NewHolidayUsecase(holidayRepo)
	reportUcase := usecase.NewReportUsecase(reportRepo)
	salaryVarUcase := usecase.NewSalaryVariableUsecase(salaryVarRepo)
	customerUcase := usecase.NewCustomerUsecase(customerRepo)
	transactionUcase := usecase.NewTransactionUsecase(transactionRepo, productRepo)
	payrollUcase := usecase.NewPayrollUsecase(employeeRepo, overtimeRepo, cashAdvanceRepo, salaryVarRepo)
	productUcase := usecase.NewProductUsecase(productRepo)
	adminTypeUcase := usecase.NewAdminTypeUsecase(adminTypeRepo)

	// App level Auth Middleware
	authMiddleware := middleware.AuthMiddleware(tokenMaker)

	// ---------------------------------------------------------
	// API Routes grouped by entities
	// based on api_design.md
	// ---------------------------------------------------------

	// 0. Auth
	authHandler := handler.NewAuthHandler(authUcase)
	authRoutes := r.Group("/auth")
	{
		authRoutes.POST("/login", authHandler.Login)
		authRoutes.POST("/refresh", authHandler.Refresh)
		authRoutes.POST("/logout", authHandler.Logout)
	}

	// 1. Employees (Implementation complete, protected by auth)
	employeeHandler := handler.NewEmployeeHandler(employeeUcase)
	eRoutes := r.Group("/employees")
	eRoutes.Use(authMiddleware)
	{
		eRoutes.GET("", employeeHandler.Fetch)
		eRoutes.GET("/:id", employeeHandler.GetByID)
		eRoutes.POST("", employeeHandler.Store)
		eRoutes.PUT("/:id", employeeHandler.Update)
		eRoutes.DELETE("/:id", employeeHandler.Delete)
		eRoutes.GET("/:id/leaves", notImplemented) // Specific employee leaves
	}

	// 1b. Admin Types (Employee Roles)
	adminTypeHandler := handler.NewAdminTypeHandler(adminTypeUcase)
	adminTypes := r.Group("/admin-types").Use(authMiddleware)
	{
		adminTypes.GET("", adminTypeHandler.Fetch)
	}

	// 2. Attendance
	attendanceHandler := handler.NewAttendanceHandler(attendanceUcase, employeeUcase)
	attendances := r.Group("/attendances").Use(authMiddleware)
	{
		attendances.GET("", attendanceHandler.Fetch)
		attendances.POST("/clock-in", attendanceHandler.ClockIn)
		attendances.POST("/clock-out", attendanceHandler.ClockOut)
		attendances.GET("/today", attendanceHandler.GetToday)
	}

	// 3. Leaves (Cuti / Sakit)
	leaveHandler := handler.NewLeaveHandler(leaveUcase)
	leaves := r.Group("/leaves").Use(authMiddleware)
	{
		leaves.GET("", leaveHandler.Fetch)
		leaves.POST("", leaveHandler.Store)
		leaves.PUT("/:id/status", leaveHandler.UpdateStatus)
	}

	// 4. Overtime (Lembur)
	overtimeHandler := handler.NewOvertimeHandler(overtimeUcase, employeeUcase)
	overtimes := r.Group("/overtimes").Use(authMiddleware)
	{
		overtimes.GET("", overtimeHandler.Fetch)
		overtimes.POST("", overtimeHandler.Store)
		overtimes.PUT("/:id/status", overtimeHandler.UpdateStatus)
	}

	// 5. Holidays (Libur)
	holidayHandler := handler.NewHolidayHandler(holidayUcase)
	holidays := r.Group("/holidays").Use(authMiddleware)
	{
		holidays.GET("", holidayHandler.Fetch)
		holidays.POST("", holidayHandler.Store)
		holidays.PUT("/:id", holidayHandler.Update)
		holidays.DELETE("/:id", holidayHandler.Delete)
	}

	// 6. Cash Advances (Kasbon)
	cashAdvanceHandler := handler.NewCashAdvanceHandler(cashAdvanceUcase, employeeUcase)
	cashAdvances := r.Group("/cash-advances").Use(authMiddleware)
	{
		cashAdvances.GET("", cashAdvanceHandler.Fetch)
		cashAdvances.POST("", cashAdvanceHandler.Store)
		cashAdvances.PUT("/:id/status", cashAdvanceHandler.UpdateStatus)
		cashAdvances.POST("/payment", cashAdvanceHandler.AddPayment)
	}

	// 7. Reports & Analytics
	reportHandler := handler.NewReportHandler(reportUcase, employeeUcase)
	reports := r.Group("/reports").Use(authMiddleware)
	{
		reports.GET("/monthly-summary", reportHandler.GetMonthlySummary)
		reports.GET("/dashboard-stats", reportHandler.GetDashboardStats)
		reports.GET("/pending-approvals", reportHandler.GetPendingApprovals)
		reports.GET("/recent-activities", reportHandler.GetRecentActivities)
		reports.GET("/commerce-stats", reportHandler.GetCommerceStats)
	}

	// 8. Salary Variables (Variabel Gaji)
	salaryVarHandler := handler.NewSalaryVariableHandler(salaryVarUcase)
	salaryVars := r.Group("/salary-variables").Use(authMiddleware)
	{
		salaryVars.GET("", salaryVarHandler.Fetch)
		salaryVars.GET("/:id", salaryVarHandler.GetByID)
		salaryVars.POST("", salaryVarHandler.Store)
		salaryVars.PUT("/:id", salaryVarHandler.Update)
		salaryVars.DELETE("/:id", salaryVarHandler.Delete)
	}

	// 9. Payroll (Penggajian) - Superadmin only
	payrollHandler := handler.NewPayrollHandler(payrollUcase, employeeUcase)
	payroll := r.Group("/payroll").Use(authMiddleware)
	{
		payroll.GET("", payrollHandler.Generate)
		payroll.GET("/:id", payrollHandler.GenerateByEmployee)
	}

	// 10. Customers (Commerce)
	customerHandler := handler.NewCustomerHandler(customerUcase)
	customers := r.Group("/customers").Use(authMiddleware)
	{
		customers.GET("", customerHandler.Fetch)
		customers.GET("/:id", customerHandler.GetByID)
		customers.POST("", customerHandler.Store)
		customers.PUT("/:id", customerHandler.Update)
		customers.DELETE("/:id", customerHandler.Delete)
		customers.GET("/types", customerHandler.GetTypes)
		customers.POST("/types", customerHandler.CreateType)
		customers.PUT("/types/:id", customerHandler.UpdateType)
		customers.DELETE("/types/:id", customerHandler.DeleteType)
		customers.GET("/:id/addresses", customerHandler.GetAddresses)
		customers.POST("/addresses", customerHandler.CreateAddress)
		customers.PUT("/addresses/:address_id", customerHandler.UpdateAddress)
		customers.DELETE("/addresses/:address_id", customerHandler.DeleteAddress)
	}

	// 11. Transactions (Commerce)
	transactionHandler := handler.NewTransactionHandler(transactionUcase)
	transactions := r.Group("/transactions").Use(authMiddleware)
	{
		transactions.GET("", transactionHandler.Fetch)
		transactions.GET("/:id", transactionHandler.GetByID)
		transactions.POST("", transactionHandler.Store)
		transactions.PUT("/:id", transactionHandler.Update)
		transactions.DELETE("/:id", transactionHandler.Delete)
		transactions.GET("/statuses", transactionHandler.GetStatuses)
		transactions.GET("/generate-code", transactionHandler.GenerateCode)
		transactions.GET("/:id/logs", transactionHandler.GetLogs)
	}

	// 12. Products & Inventory
	productHandler := handler.NewProductHandler(productUcase)
	products := r.Group("/products").Use(authMiddleware)
	{
		products.GET("/inventory", productHandler.GetInventory)
		products.GET("/inventory/logs/:id", productHandler.GetStockLogs)
		products.GET("/sizes", productHandler.GetSizes)
		products.GET("/sizes_type", productHandler.GetSizesType)
		products.GET("/colors", productHandler.GetColors)
		products.PUT("/inventory/:id/stock", productHandler.UpdateStock)

		products.GET("/codes", productHandler.GetCodes)
		products.GET("/code_with_types", productHandler.GetCodesWithTypes)
		products.POST("/codes", productHandler.CreateCode)
		products.PUT("/codes/:id", productHandler.UpdateCode)
		products.DELETE("/codes/:id", productHandler.DeleteCode)

		products.GET("/types", productHandler.GetTypes)
		products.POST("/types", productHandler.CreateType)
		products.PUT("/types/:id", productHandler.UpdateType)
		products.DELETE("/types/:id", productHandler.DeleteType)

		products.POST("/", productHandler.CreateProduct)
		products.PUT("/:id", productHandler.UpdateProduct)
		products.DELETE("/:id", productHandler.DeleteProduct)
	}
}

// notImplemented is a placeholder handler for routes that are defined but not yet implemented
func notImplemented(c *gin.Context) {
	handler.ErrorResponse(c, http.StatusNotImplemented, "Endpoint is not implemented yet")
}
