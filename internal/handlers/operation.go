package handlers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"userBalanceAvito/internal/models"
	"userBalanceAvito/internal/utils"
)

func GetOperations(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		operationID, err := utils.GetHeaderString(c, "OperationID")
		if err != nil {
			utils.HandleError(c, 400, err)
			return
		}

		var operation models.Operation
		if err := db.Where("id = ?", operationID).First(&operation).Error; err != nil {
			utils.HandleError(c, 400, err)
			return
		}

		c.JSON(200, gin.H{"operation": operation})
	}
}

func CreateOperations(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var operation models.Operation

		if err := c.ShouldBindJSON(&operation); err != nil {
			utils.HandleError(c, 400, err)
			return
		}

		if operation.AccountID == 0 || operation.Amount == 0 || operation.OperationType == "" {
			utils.HandleError(c, 400, fmt.Errorf("AccountID, Amount, and OperationType are required"))
			return
		}

		if err := db.Create(&operation).Error; err != nil {
			utils.HandleError(c, 400, err)
			return
		}

		c.JSON(200, gin.H{"operation": operation})
	}
}
