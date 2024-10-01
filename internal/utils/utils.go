package utils

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"strconv"
)

func HandleError(c *gin.Context, statusCode int, err error) {
	c.JSON(statusCode, gin.H{"error": err.Error()})
}

func GetHeaderInt64(c *gin.Context, key string) (int64, error) {
	header := c.GetHeader(key)
	if header == "" {
		return 0, fmt.Errorf("%s header is empty", key)
	}
	return strconv.ParseInt(header, 10, 64)
}

func GetHeaderString(c *gin.Context, key string) (string, error) {
	header := c.GetHeader(key)
	if header == "" {
		return "", fmt.Errorf("%s header is empty", key)
	}
	return header, nil
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
