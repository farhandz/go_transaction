package controllers

import (
	"gin-boilerplate/models"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"math"
	"gin-boilerplate/helpers"
	"gin-boilerplate/repository"
)





type TransactionController struct {
	Repo repository.TransactionRepository
}


func (tc *TransactionController) GetTransactions(ctx *gin.Context) {
	// Ambil query params

	pageStr := ctx.Query("page_number")
	pageSizeStr := ctx.Query("page_size")
	status := ctx.Query("status")
	userIDStr := ctx.Query("user_id")

	// Default pagination
	pageNumber := 1
	pageSize := 10
	var userID int

	if pageStr != "" {
		pageNumber, _ = strconv.Atoi(pageStr)
	}
	if pageSizeStr != "" {
		pageSize, _ = strconv.Atoi(pageSizeStr)
	}
	if userIDStr != "" {
		userID, _ = strconv.Atoi(userIDStr)
	}

	// Ambil data dari repository
	var transactions []models.Transaction
	totalRecordCount, err := tc.Repo.GetTransactionsWithFilters(&transactions, pageNumber, pageSize, status, userID)
	if err != nil {
		helpers.Error(ctx, err.Error(), nil)
		return
	}


	data :=  gin.H{
		"page_number":        pageNumber,
		"page_size":          pageSize,
		"total_record_count": totalRecordCount,
		"data":               transactions,
	}

	helpers.Success(ctx, "success get data", data)
	
}

func (tc *TransactionController) GetDashboardReport(ctx *gin.Context) {
	// Total transaksi sukses hari ini
	totalSuccessToday, err := tc.Repo.CountSuccessToday()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Total transaksi dan unique user
	totalTransactions, err := tc.Repo.CountTotalTransactions()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	uniqueUsers, err := tc.Repo.CountUniqueUsers()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Rata-rata jumlah transaksi per user
	var avgTransactionPerUser float64
	if uniqueUsers > 0 {
		avgTransactionPerUser = float64(totalTransactions) / float64(uniqueUsers)
	} else {
		avgTransactionPerUser = 0
	}

	// Daftar 10 transaksi terbaru
	latestTransactions, err := tc.Repo.GetLatestTransactions(10)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	data := gin.H{
		"total_success_today":         totalSuccessToday,
		"average_transaction_per_user": math.Round(avgTransactionPerUser*100) / 100, // dibulatkan 2 desimal
		"latest_transactions":         latestTransactions,
	}

	helpers.Success(ctx, "success create", data)
	
}

func (tc *TransactionController) GetTransactionByID(ctx *gin.Context) {

	// Ambil param ID dari URL
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)

	// Ambil data transaksi dari repository
	transaction, err := tc.Repo.GetTransactionByID(id)
	if err != nil {
		helpers.Error(ctx, "Transaction not found", nil)
		return
	}
	helpers.Success(ctx, "success get by id", transaction)
}

type UpdateStatusRequest struct {
	Status string `json:"status" binding:"required"`
}


func (tc *TransactionController)  UpdateTransactionStatus(ctx *gin.Context) {
	// Ambil param ID dari URL
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		helpers.Error(ctx, "Invalid transaction ID", nil)
		return
	}

	// Binding request body
	var req UpdateStatusRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		helpers.Error(ctx, "Invalid request body", nil)
		return
	}

	// Validasi status hanya boleh: "success", "pending", "failed"
	validStatuses := map[string]bool{
		"success": true,
		"pending": true,
		"failed":  true,
	}

	if !validStatuses[req.Status] {
		helpers.Error(ctx, "Invalid status value. Allowed values: success, pending, failed", nil)
		return
	}

	// Update transaksi lewat repository
	transaction, err := tc.Repo.UpdateTransactionStatus(id, req.Status)
	if err != nil {
		helpers.Error(ctx, "Transaction not found", nil)
		return
	}

	// Return response sukses
	helpers.Success(ctx, "Success Update Transaction Status", transaction)
}


func (tc *TransactionController) GetDashboardSummary(ctx *gin.Context) {
	summary, err := tc.Repo.GetTransactionSummary()
	if err != nil {
		helpers.Error(ctx, "Failed to fetch summary data", nil)
		return
	}

	helpers.Success(ctx, "Success get summary", summary)
}


func (tc *TransactionController) DeleteTransaction(ctx *gin.Context) {
	// Ambil param ID dari URL
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		helpers.Error(ctx, "Invalid transaction ID", nil)
		return
	}

	// Hapus transaksi pakai repository
	err = tc.Repo.DeleteTransactionByID(id)
	if err != nil {
		helpers.Error(ctx, "Transaction not found", nil)
		return
	}

	helpers.Success(ctx, "Success delete", nil)
}


func (tc *TransactionController) CreateTransaction(ctx *gin.Context) {
	transaction := new(models.Transaction)

	if err := ctx.ShouldBindJSON(transaction); err != nil {
		helpers.Error(ctx, err.Error(), nil)
		return
	}

	if err := tc.Repo.Save(transaction); err != nil {
		helpers.Error(ctx, err.Error(), nil)
		return
	}

	helpers.Success(ctx, "Success Create", transaction)
}

