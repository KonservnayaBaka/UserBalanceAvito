package tests

import (
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"userBalanceAvito/internal/handlers"
	"userBalanceAvito/internal/models"
)

func createTestWithHeaderRequest(method, url string, headers map[string]string) *http.Request {
	req, _ := http.NewRequest(method, url, nil)
	for key, value := range headers {
		req.Header.Set(key, value)
	}
	return req
}

func TestEnrollmentBalance(t *testing.T) {
	db := setupTestDB()
	router := gin.Default()
	router.PUT("/balance/enrollment", handlers.EnrollmentBalance(db))

	account := models.Account{
		Balance: 0,
	}
	db.Create(&account)

	headers := map[string]string{
		"UserID":       strconv.FormatUint(uint64(account.ID), 10),
		"ChangeNumber": "5000",
	}

	req := createTestWithHeaderRequest("PUT", "/balance/enrollment", headers)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)

	var updatedAccount models.Account
	db.First(&updatedAccount, account.ID)

	assert.Equal(t, 5000, updatedAccount.Balance)
}

func TestWriteOffBalance(t *testing.T) {
	db := setupTestDB()
	router := gin.Default()
	router.PUT("/balance/off", handlers.WriteOffBalance(db))

	account := models.Account{
		Balance: 1500,
	}
	db.Create(&account)

	headers := map[string]string{
		"UserID":       strconv.FormatUint(uint64(account.ID), 10),
		"ChangeNumber": "1000",
	}

	req := createTestWithHeaderRequest("PUT", "/balance/off", headers)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)

	var updatedAccount models.Account
	db.First(&updatedAccount, account.ID)

	assert.Equal(t, 500, updatedAccount.Balance)
}

func TestTransferBalance(t *testing.T) {
	db := setupTestDB()
	router := gin.Default()
	router.PUT("/balance/transfer", handlers.TransferBalance(db))

	account1 := models.Account{
		Balance: 1000,
	}
	account2 := models.Account{
		Balance: 0,
	}
	db.Create(&account1)
	db.Create(&account2)

	headers := map[string]string{
		"RecipientID":  strconv.FormatUint(uint64(account2.ID), 10),
		"SenderID":     strconv.FormatUint(uint64(account1.ID), 10),
		"ChangeNumber": strconv.FormatUint(uint64(500), 10),
	}

	req := createTestWithHeaderRequest("PUT", "/balance/transfer", headers)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)

	var updatedAccount1 models.Account
	db.First(&updatedAccount1, account1.ID)
	var updatedAccount2 models.Account
	db.First(&updatedAccount2, account2.ID)

	assert.Equal(t, 500, updatedAccount1.Balance)
	assert.Equal(t, 500, updatedAccount2.Balance)
}
