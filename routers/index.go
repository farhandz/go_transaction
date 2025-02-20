package routers

import (
	"gin-boilerplate/controllers"
	"gin-boilerplate/repository"

	"github.com/gin-gonic/gin"
	"net/http"
)

// RegisterRoutes add all routing list here automatically get main router
func RegisterRoutes(route *gin.Engine) {
	route.NoRoute(func(ctx *gin.Context) {
		ctx.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "Route Not Found"})
	})

	// Inisialisasi Repository dan Controller
	transactionRepo := &repository.TransactionRepositoryImpl{}
	transactionController := &controllers.TransactionController{
		Repo: transactionRepo,
	}

	route.GET("/health", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{"live": "sip ss sudahh runningg"})
	})

	// Add All route
	route.GET("/transaction", transactionController.GetTransactions)
	route.GET("/transaction/:id", transactionController.GetTransactionByID)
	route.GET("/dashboard/summary", transactionController.GetDashboardSummary)
	route.DELETE("/transaction/:id", transactionController.DeleteTransaction)
	route.PUT("/transaction/:id", transactionController.UpdateTransactionStatus)
	route.POST("/transaction", transactionController.CreateTransaction)
	route.GET("/dashboard/report", transactionController.GetDashboardReport)
}
