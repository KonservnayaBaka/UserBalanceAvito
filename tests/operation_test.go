package tests

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
	"userBalanceAvito/internal/handlers"
	"userBalanceAvito/internal/models"
)

func TestCreateOperations(t *testing.T) {
	db := setupTestDB()
	router := gin.Default()
	router.POST("/operations/create", handlers.CreateOperations(db))

	account := models.Account{
		Balance: 0,
	}
	db.Create(&account)

	operation := models.Operation{
		AccountID:     account.ID,
		Amount:        1000,
		OperationType: "deposit",
		CreatedAt:     time.Now(),
	}

	req, w := createTestRequest("POST", "/operations/create", operation)
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	assert.NotEmpty(t, response["operation"])
	assert.Equal(t, float64(account.ID), response["operation"].(map[string]interface{})["account_id"])
	assert.Equal(t, float64(1000), response["operation"].(map[string]interface{})["amount"])
	assert.Equal(t, "deposit", response["operation"].(map[string]interface{})["operation_type"])
}

func TestCreateOperationsMissingFields(t *testing.T) {
	db := setupTestDB()
	router := gin.Default()
	router.POST("/operations/create", handlers.CreateOperations(db))

	operation := models.Operation{
		AccountID:     0,
		Amount:        0,
		OperationType: "",
	}

	req, w := createTestRequest("POST", "/operations/create", operation)
	router.ServeHTTP(w, req)

	assert.Equal(t, 400, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	assert.NotEmpty(t, response["error"])
}
