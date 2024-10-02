package handlers

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"net/http"
	"net/http/httptest"
	"testing"
	"userBalanceAvito/internal/models"
	"userBalanceAvito/internal/utils"
)

func setupTestDB() *gorm.DB {
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		panic("failed to connect to database")
	}
	db.AutoMigrate(&models.User{}, &models.Account{}, &models.Product{}, &models.Reservation{}, &models.Operation{})
	return db
}

func createTestRequest(method, url string, body interface{}) (*http.Request, *httptest.ResponseRecorder) {
	jsonData, _ := json.Marshal(body)
	req, _ := http.NewRequest(method, url, bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	return req, w
}

func TestRegistrationNewUser(t *testing.T) {
	db := setupTestDB()
	router := gin.Default()
	router.POST("/auth/registration", RegistrationNewUser(db))

	user := models.User{
		Username: "testusername",
		Password: "testpassword",
	}

	req, w := createTestRequest("POST", "/auth/registration", user)
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	assert.NotEmpty(t, response["user"])
	assert.NotEmpty(t, response["account"])
}

func TestAuthorizationUser(t *testing.T) {
	db := setupTestDB()
	router := gin.Default()
	router.POST("/auth/login", AuthorizationUser(db))

	hashedPassword, _ := utils.HashPassword("testpassword")
	testUser := models.User{
		Username: "testusername",
		Password: hashedPassword,
	}
	db.Create(&testUser)

	user := models.User{
		Username: "testusername",
		Password: "testpassword",
	}

	req, w := createTestRequest("POST", "/auth/login", user)
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	assert.Equal(t, true, response["status"])
	assert.NotEmpty(t, response["user"])
}

func TestRegistrationNewUserEmptyData(t *testing.T) {
	db := setupTestDB()
	router := gin.Default()
	router.POST("/auth/registration", RegistrationNewUser(db))

	user := models.User{
		Username: "",
		Password: "",
	}

	req, w := createTestRequest("POST", "/auth/registration", user)
	router.ServeHTTP(w, req)

	assert.Equal(t, 400, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	assert.NotEmpty(t, response["error"])
}

func TestAuthorizationUserInvalidData(t *testing.T) {
	db := setupTestDB()
	router := gin.Default()
	router.POST("/auth/login", AuthorizationUser(db))

	user := models.User{
		Username: "invalidusername",
		Password: "invalidpassword",
	}

	req, w := createTestRequest("POST", "/auth/login", user)
	router.ServeHTTP(w, req)

	assert.Equal(t, 400, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	assert.NotEmpty(t, response["error"])
}
