
package controllers
import (
	"net/http"
	"net/http/httptest"
	"testing"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gin-boilerplate/models"
    "errors"
    "bytes"
    "encoding/json"
    "gin-boilerplate/repository"
)

// MockTransactionRepository mocks the TransactionRepository interface

type MockTransactionRepository struct {
	mock.Mock
}

func (m *MockTransactionRepository) GetTransactionByID(id int) (*models.Transaction, error) {
	args := m.Called(id)
	if args.Get(0) != nil {
		return args.Get(0).(*models.Transaction), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockTransactionRepository) CountSuccessToday() (int, error) {
	args := m.Called()
	return args.Int(0), args.Error(1)
}

func (m *MockTransactionRepository) CountTotalTransactions() (int, error) {
	args := m.Called()
	return args.Int(0), args.Error(1)
}

func (m *MockTransactionRepository) CountUniqueUsers() (int, error) {
	args := m.Called()
	return args.Int(0), args.Error(1)
}

func (m *MockTransactionRepository) GetLatestTransactions(limit int) ([]models.Transaction, error) {
	args := m.Called(limit)
	return args.Get(0).([]models.Transaction), args.Error(1)
}

func (m *MockTransactionRepository) GetTransactionSummary() (repository.TransactionSummary, error) {
	args := m.Called()
	return args.Get(0).(repository.TransactionSummary), args.Error(1)
}

func (m *MockTransactionRepository) UpdateTransactionStatus(id int, status string) (*models.Transaction, error) {
	args := m.Called(id, status)
	if args.Get(0) != nil {
		return args.Get(0).(*models.Transaction), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockTransactionRepository) DeleteTransactionByID(id int) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockTransactionRepository) Save(transaction *models.Transaction) error {
	args := m.Called(transaction)
	return args.Error(0)
}

func (m *MockTransactionRepository) GetTransactionsWithFilters(transactions *[]models.Transaction, pageNumber, pageSize int, status string, userID int) (int64, error) {
	args := m.Called(transactions, pageNumber, pageSize, status, userID)
	return args.Get(0).(int64), args.Error(1)
}


func TestGetTransactionByID_Success(t *testing.T) {
    gin.SetMode(gin.TestMode)

    mockRepo := new(MockTransactionRepository)
    controller := &TransactionController{Repo: mockRepo}

    transaction := &models.Transaction{
        ID:     1,
        UserID: 123,
        Amount: 1000,
        Status: "completed",
    }
    mockRepo.On("GetTransactionByID", 1).Return(transaction, nil)

    w := httptest.NewRecorder()
    ctx, _ := gin.CreateTestContext(w)
    ctx.Params = []gin.Param{{Key: "id", Value: "1"}}

    controller.GetTransactionByID(ctx)

    assert.Equal(t, http.StatusOK, w.Code)
    assert.Contains(t, w.Body.String(), `"id":1`)
    assert.Contains(t, w.Body.String(), `"user_id":123`)
    assert.Contains(t, w.Body.String(), `"amount":1000`)
    assert.Contains(t, w.Body.String(), `"status":"completed"`)

    mockRepo.AssertExpectations(t)
}

func TestGetTransactionByID_NotFound(t *testing.T) {
    gin.SetMode(gin.TestMode)

    mockRepo := new(MockTransactionRepository)
    controller := &TransactionController{Repo: mockRepo}

    mockRepo.On("GetTransactionByID", 1).Return(nil, errors.New("not found"))

    w := httptest.NewRecorder()
    ctx, _ := gin.CreateTestContext(w)
    ctx.Params = []gin.Param{{Key: "id", Value: "1"}}

    controller.GetTransactionByID(ctx)

    assert.Contains(t, w.Body.String(), "Transaction not found")
}

func TestGetDashboardReport_Success(t *testing.T) {
    gin.SetMode(gin.TestMode)
    w := httptest.NewRecorder()
    ctx, _ := gin.CreateTestContext(w)

    mockRepo := new(MockTransactionRepository)
    mockRepo.On("CountSuccessToday").Return(5, nil)
    mockRepo.On("CountTotalTransactions").Return(100, nil)
    mockRepo.On("CountUniqueUsers").Return(10, nil)
    mockRepo.On("GetLatestTransactions", 10).Return([]models.Transaction{
        {ID: 1, Status: "success"},
        {ID: 2, Status: "failed"},
    }, nil)

    controller := &TransactionController{Repo: mockRepo}
    controller.GetDashboardReport(ctx)

    assert.Equal(t, http.StatusOK, w.Code)
    assert.Contains(t, w.Body.String(), `"total_success_today":5`)
    assert.Contains(t, w.Body.String(), `"average_transaction_per_user":10`)
    assert.Contains(t, w.Body.String(), `"latest_transactions"`)
}

func TestGetDashboardSummary_Success(t *testing.T) {
    gin.SetMode(gin.TestMode)
    w := httptest.NewRecorder()
    ctx, _ := gin.CreateTestContext(w)

    mockRepo := new(MockTransactionRepository)
    mockRepo.On("GetTransactionSummary").Return(repository.TransactionSummary{
        TotalTransactions: 200,
        UniqueUsers:        50,
    }, nil)

    controller := &TransactionController{Repo: mockRepo}
    controller.GetDashboardSummary(ctx)

    assert.Equal(t, http.StatusOK, w.Code)
    assert.Contains(t, w.Body.String(), `"total_transactions":200`)
    assert.Contains(t, w.Body.String(), `"unique_users":50`)
}


func TestGetDashboardSummary_Failure(t *testing.T) {
	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)

	mockRepo := new(MockTransactionRepository)

    mockRepo.On("GetTransactionSummary").Return(repository.TransactionSummary{
        TotalTransactions: 200,
        UniqueUsers:        50,
    }, errors.New("database error"))
	controller := &TransactionController{Repo: mockRepo}
	controller.GetDashboardSummary(ctx)

	assert.Equal(t, http.StatusBadRequest, w.Code) // Karena pakai helpers.Error jadi 400
	assert.Contains(t, w.Body.String(), `"status":"error"`)
	assert.Contains(t, w.Body.String(), `"message":"Failed to fetch summary data"`)
}

func TestUpdateTransactionStatus_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	ctx.Params = gin.Params{{Key: "id", Value: "1"}}

	body := `{"status":"success"}`
	ctx.Request = httptest.NewRequest(http.MethodPut, "/transactions/1", bytes.NewBufferString(body))
	ctx.Request.Header.Set("Content-Type", "application/json")

	mockRepo := new(MockTransactionRepository)
	mockTransaction := &models.Transaction{ID: 1, Status: "success"}

	mockRepo.On("UpdateTransactionStatus", 1, "success").Return(mockTransaction, nil)

	controller := &TransactionController{Repo: mockRepo}
	controller.UpdateTransactionStatus(ctx)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), `"status":"success"`)
	assert.Contains(t, w.Body.String(), `"message":"Success Update Transaction Status"`)
	assert.Contains(t, w.Body.String(), `"id":1`)
	assert.Contains(t, w.Body.String(), `"status":"success"`)
}

func TestUpdateTransactionStatus_InvalidID(t *testing.T) {
	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	ctx.Params = gin.Params{{Key: "id", Value: "abc"}} // Invalid ID

	controller := &TransactionController{}
	controller.UpdateTransactionStatus(ctx)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), `"status":"error"`)
	assert.Contains(t, w.Body.String(), `"message":"Invalid transaction ID"`)
}

func TestUpdateTransactionStatus_InvalidRequestBody(t *testing.T) {
	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	ctx.Params = gin.Params{{Key: "id", Value: "1"}}

	body := `{"invalid_key":"value"}`
	ctx.Request = httptest.NewRequest(http.MethodPut, "/transactions/1", bytes.NewBufferString(body))
	ctx.Request.Header.Set("Content-Type", "application/json")

	controller := &TransactionController{}
	controller.UpdateTransactionStatus(ctx)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), `"status":"error"`)
	assert.Contains(t, w.Body.String(), `"message":"Invalid request body"`)
}

func TestUpdateTransactionStatus_InvalidStatusValue(t *testing.T) {
	mockRepo := new(MockTransactionRepository)
	controller := &TransactionController{Repo: mockRepo}

	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)

	// Simulasi request dengan status yang tidak valid
	reqBody := `{"status": "unknown"}`
	ctx.Request = httptest.NewRequest("PUT", "/transactions/1", bytes.NewBufferString(reqBody))
	ctx.Params = []gin.Param{{Key: "id", Value: "1"}}

	controller.UpdateTransactionStatus(ctx)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	body := w.Body.String()
	assert.Contains(t, body, "Invalid status value")
}

func TestUpdateTransactionStatus_TransactionNotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	ctx.Params = gin.Params{{Key: "id", Value: "1"}}

	body := `{"status":"success"}`
	ctx.Request = httptest.NewRequest(http.MethodPut, "/transactions/1", bytes.NewBufferString(body))
	ctx.Request.Header.Set("Content-Type", "application/json")

	mockRepo := new(MockTransactionRepository)
	mockRepo.On("UpdateTransactionStatus", 1, "success").Return(nil, errors.New("not found"))

	controller := &TransactionController{Repo: mockRepo}
	controller.UpdateTransactionStatus(ctx)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), `"status":"error"`)
	assert.Contains(t, w.Body.String(), `"message":"Transaction not found"`)
}


func TestDeleteTransaction_Success(t *testing.T) {
	mockRepo := new(MockTransactionRepository)
	mockRepo.On("DeleteTransactionByID", 1).Return(nil)

	controller := &TransactionController{Repo: mockRepo}

	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	ctx.Params = []gin.Param{{Key: "id", Value: "1"}}

	controller.DeleteTransaction(ctx)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Success delete")
}

func TestDeleteTransaction_InvalidID(t *testing.T) {
	controller := &TransactionController{}

	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	ctx.Params = []gin.Param{{Key: "id", Value: "abc"}} // ID invalid

	controller.DeleteTransaction(ctx)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "Invalid transaction ID")
}

func TestDeleteTransaction_NotFound(t *testing.T) {
	mockRepo := new(MockTransactionRepository)
	mockRepo.On("DeleteTransactionByID", 1).Return(errors.New("not found"))

	controller := &TransactionController{Repo: mockRepo}

	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	ctx.Params = []gin.Param{{Key: "id", Value: "1"}}

	controller.DeleteTransaction(ctx)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "Transaction not found")
}

func TestCreateTransaction_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockRepo := new(MockTransactionRepository)
	controller := &TransactionController{Repo: mockRepo}

	transaction := &models.Transaction{
		ID:     1,
		Amount: 100000,
		Status: "pending",
	}

	mockRepo.On("Save", transaction).Return(nil)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	body, _ := json.Marshal(transaction)
	c.Request, _ = http.NewRequest(http.MethodPost, "/transactions", bytes.NewBuffer(body))
	c.Request.Header.Set("Content-Type", "application/json")

	controller.CreateTransaction(c)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.Equal(t, "success", response["status"])
	assert.Equal(t, "Success Create", response["message"])
	assert.NotNil(t, response["data"])

	mockRepo.AssertExpectations(t)
}

func TestCreateTransaction_InvalidRequestBody(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockRepo := new(MockTransactionRepository)
	controller := &TransactionController{Repo: mockRepo}

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// Invalid JSON
	c.Request, _ = http.NewRequest(http.MethodPost, "/transactions", bytes.NewBuffer([]byte(`{invalid json}`)))
	c.Request.Header.Set("Content-Type", "application/json")

	controller.CreateTransaction(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.Equal(t, "error", response["status"])
	assert.Contains(t, response["message"], "invalid character")

	mockRepo.AssertExpectations(t)
}

func TestCreateTransaction_FailedToSave(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockRepo := new(MockTransactionRepository)
	controller := &TransactionController{Repo: mockRepo}

	transaction := &models.Transaction{
		ID:     1,
		Amount: 100000,
		Status: "pending",
	}

	mockRepo.On("Save", transaction).Return(errors.New("failed to save transaction"))

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	body, _ := json.Marshal(transaction)
	c.Request, _ = http.NewRequest(http.MethodPost, "/transactions", bytes.NewBuffer(body))
	c.Request.Header.Set("Content-Type", "application/json")

	controller.CreateTransaction(c)

	assert.Equal(t, 400, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.Equal(t, "error", response["status"])
	assert.Equal(t, "failed to save transaction", response["message"])

	mockRepo.AssertExpectations(t)
}

func TestGetTransactions_Success(t *testing.T) {
    mockRepo := new(MockTransactionRepository)
    controller := TransactionController{Repo: mockRepo}

    w := httptest.NewRecorder()
    c, _ := gin.CreateTestContext(w)
    req, _ := http.NewRequest(http.MethodGet, "/transactions?page_number=1&page_size=10", nil)
    c.Request = req

    mockRepo.On("GetTransactionsWithFilters", mock.Anything, 1, 10, "", 0).
        Return(int64(2), nil).
        Run(func(args mock.Arguments) {
            ptr := args.Get(0).(*[]models.Transaction)
            *ptr = []models.Transaction{
                {ID: 1, Status: "success"},
                {ID: 2, Status: "failed"},
            }
        })

    controller.GetTransactions(c)

    assert.Equal(t, http.StatusOK, w.Code)

    var response struct {
        Status           string                 `json:"status"`
        Message          string                 `json:"message"`
        Data             map[string]interface{} `json:"data"`
    }

    err := json.Unmarshal(w.Body.Bytes(), &response)
    assert.NoError(t, err)
    assert.Equal(t, "success", response.Status)
    assert.Equal(t, "success get data", response.Message)

    // Biasanya number jadi float64 karena JSON parsing
    assert.EqualValues(t, 2, response.Data["total_record_count"])
    assert.Len(t, response.Data["data"], 2)
}





// func TestGetTransactions_ErrorFetching(t *testing.T) {
// 	gin.SetMode(gin.TestMode)

// 	mockRepo := new(MockTransactionRepository)
// 	controller := &TransactionController{Repo: mockRepo}

//     var transactions []models.Transaction
//     mockRepo.On("GetTransactionsWithFilters", &transactions, 1, 10, "", 0).
//     Return(int64(0), errors.New("database error"))

// 	w := httptest.NewRecorder()
// 	c, _ := gin.CreateTestContext(w)
// 	c.Request, _ = http.NewRequest(http.MethodGet, "/transactions", nil)

// 	controller.GetTransactions(c)

// 	assert.Equal(t, http.StatusBadRequest, w.Code)

// 	mockRepo.AssertExpectations(t)
// }