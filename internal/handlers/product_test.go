package handlers

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"userBalanceAvito/internal/models"
)

func TestAddProduct(t *testing.T) {
	db := setupTestDB()
	router := gin.Default()
	router.POST("/product/add", AddProduct(db))

	product := models.Product{
		Name: "Iphone 13",
	}

	jsonBody, _ := json.Marshal(product)
	req, _ := http.NewRequest("POST", "/product/add", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	productResponse, ok := response["product"].(map[string]interface{})
	assert.True(t, ok)
	assert.NotEmpty(t, productResponse["name"])
}

func TestAddReservation(t *testing.T) {
	db := setupTestDB()
	router := gin.Default()
	router.POST("/product/reservation", ReservProduct(db))

	account := models.Account{}
	db.Create(&account)
	product := models.Product{
		Name: "Iphone 13",
	}
	db.Create(&product)

	headers := map[string]string{
		"AccountID": strconv.FormatUint(uint64(account.ID), 10),
		"ProductID": strconv.FormatUint(uint64(product.ID), 10),
		"Ammount":   strconv.FormatUint(uint64(25), 10),
	}

	req := createTestWithHeaderRequest("POST", "/product/reservation", headers)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	reservationResponse, ok := response["reservation"].(map[string]interface{})
	assert.True(t, ok)
	assert.NotEmpty(t, reservationResponse["id"])

	reservationID := uint(reservationResponse["id"].(float64))
	var reservation models.Reservation
	db.First(&reservation, reservationID)

	assert.NotEmpty(t, reservation)
	assert.Equal(t, account.ID, reservation.AccountID)
	assert.Equal(t, product.ID, reservation.ProductID)
	assert.Equal(t, 25, reservation.Amount)
}
