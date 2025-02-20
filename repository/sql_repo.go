package repository

import (
	"gin-boilerplate/infra/database"
	"gin-boilerplate/models"
	"gorm.io/gorm"
	"time"
)

type TransactionSummary struct {
	TotalTransactionsToday    int     `json:"total_transactions_today"`
	AverageTransactionPerUser float64 `json:"average_transaction_per_user"`
	TotalTransactions         int     `json:"total_transactions"`
	UniqueUsers               int     `json:"unique_users"`
	TotalPendingTransactions  int     `json:"total_pending_transactions"`
	TotalSuccessTransactions  int     `json:"total_success_transactions"`
	TotalFailedTransactions   int     `json:"total_failed_transactions"`
}

type TransactionRepository interface {
	GetTransactionByID(id int) (*models.Transaction, error)
	CountSuccessToday() (int, error)
	CountTotalTransactions() (int, error)
	CountUniqueUsers() (int, error)
	GetLatestTransactions(limit int) ([]models.Transaction, error)
	GetTransactionSummary() (TransactionSummary, error)
	UpdateTransactionStatus(id int, status string) (*models.Transaction, error)
	DeleteTransactionByID(id int) error
	Save(transaction *models.Transaction) error
	GetTransactionsWithFilters(transactions *[]models.Transaction, pageNumber, pageSize int, status string, userID int) (int64, error)
}

type TransactionRepositoryImpl struct{}

func (r *TransactionRepositoryImpl) GetTransactionByID(id int) (*models.Transaction, error) {
	var transaction models.Transaction
	err := database.DB.First(&transaction, id).Error
	if err != nil {
		return nil, err
	}
	return &transaction, nil
}

func (r *TransactionRepositoryImpl) CountSuccessToday() (int, error) {
	var count int64
	today := time.Now().Format("2006-01-02")
	err := database.DB.Model(&models.Transaction{}).
		Where("status = ? AND DATE(created_at) = ?", "success", today).
		Count(&count).Error
	return int(count), err
}

func (r *TransactionRepositoryImpl) CountTotalTransactions() (int, error) {
	var count int64
	err := database.DB.Model(&models.Transaction{}).Count(&count).Error
	return int(count), err
}

func (r *TransactionRepositoryImpl) CountUniqueUsers() (int, error) {
	var count int64
	err := database.DB.Model(&models.Transaction{}).Distinct("user_id").Count(&count).Error
	return int(count), err
}

func (r *TransactionRepositoryImpl) GetLatestTransactions(limit int) ([]models.Transaction, error) {
	var transactions []models.Transaction
	err := database.DB.Order("created_at DESC").Limit(limit).Find(&transactions).Error
	return transactions, err
}

func (r *TransactionRepositoryImpl) GetTransactionSummary() (TransactionSummary, error) {
	var summary TransactionSummary
	var totalToday, totalTransactions, uniqueUsers, totalPending, totalSuccess, totalFailed int64

	today := time.Now().Format("2006-01-02")
	database.DB.Model(&models.Transaction{}).Where("DATE(created_at) = ?", today).Count(&totalToday)
	database.DB.Model(&models.Transaction{}).Count(&totalTransactions)
	database.DB.Model(&models.Transaction{}).Distinct("user_id").Count(&uniqueUsers)
	database.DB.Model(&models.Transaction{}).Where("status = ?", "pending").Count(&totalPending)
	database.DB.Model(&models.Transaction{}).Where("status = ?", "success").Count(&totalSuccess)
	database.DB.Model(&models.Transaction{}).Where("status = ?", "failed").Count(&totalFailed)

	averageTransactionPerUser := 0.0
	if uniqueUsers > 0 {
		averageTransactionPerUser = float64(totalTransactions) / float64(uniqueUsers)
	}

	summary = TransactionSummary{
		TotalTransactionsToday:    int(totalToday),
		AverageTransactionPerUser: averageTransactionPerUser,
		TotalTransactions:         int(totalTransactions),
		UniqueUsers:                int(uniqueUsers),
		TotalPendingTransactions:  int(totalPending),
		TotalSuccessTransactions:  int(totalSuccess),
		TotalFailedTransactions:   int(totalFailed),
	}

	return summary, nil
}

func (r *TransactionRepositoryImpl) UpdateTransactionStatus(id int, status string) (*models.Transaction, error) {
	var transaction models.Transaction
	err := database.DB.First(&transaction, id).Error
	if err != nil {
		return nil, err
	}

	transaction.Status = status
	err = database.DB.Save(&transaction).Error
	if err != nil {
		return nil, err
	}

	return &transaction, nil
}

func (r *TransactionRepositoryImpl) DeleteTransactionByID(id int) error {
	result := database.DB.Delete(&models.Transaction{}, id)
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return result.Error
}

func (r *TransactionRepositoryImpl) Save(transaction *models.Transaction) error {
	return database.DB.Create(transaction).Error
}

func (r *TransactionRepositoryImpl) GetTransactionsWithFilters(transactions *[]models.Transaction, pageNumber, pageSize int, status string, userID int) (int64, error) {
	var totalRecordCount int64
	query := database.DB.Model(&models.Transaction{})

	if status != "" {
		query = query.Where("status = ?", status)
	}

	if userID != 0 {
		query = query.Where("user_id = ?", userID)
	}

	err := query.Count(&totalRecordCount).Error
	if err != nil {
		return 0, err
	}
	offset := (pageNumber - 1) * pageSize
	err = query.Limit(pageSize).Offset(offset).Find(transactions).Error
	return totalRecordCount, err
}
