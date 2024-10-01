package handlers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"strconv"
	"userBalanceAvito/internal/models"
	"userBalanceAvito/internal/utils"
)

func GetBalance(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, err := utils.GetHeaderString(c, "UserID")
		if err != nil {
			utils.HandleError(c, 400, err)
			return
		}

		account := models.Account{}
		if err := db.Where("id = ?", userID).First(&account).Error; err != nil {
			utils.HandleError(c, 400, err)
			return
		}

		c.JSON(200, gin.H{"balance": account.Balance})
	}
}

func EnrollmentBalance(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, err := utils.GetHeaderString(c, "UserID")
		if err != nil {
			utils.HandleError(c, 400, err)
			return
		}
		changeNumberStr, err := utils.GetHeaderString(c, "ChangeNumber")
		if err != nil {
			utils.HandleError(c, 400, err)
			return
		}

		parseСhangeNumber, err := strconv.ParseInt(changeNumberStr, 10, 64)
		if err != nil {
			utils.HandleError(c, 400, err)
			return
		}

		account := models.Account{}
		if err := db.Where("id = ?", userID).First(&account).Error; err != nil {
			utils.HandleError(c, 400, err)
			return
		}

		account.Balance += int(parseСhangeNumber)
		if err := db.Save(&account).Error; err != nil {
			utils.HandleError(c, 400, err)
		}

		c.JSON(200, gin.H{"balance": account.Balance})
	}
}

func WriteOffBalance(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, err := utils.GetHeaderString(c, "UserID")
		if err != nil {
			utils.HandleError(c, 400, err)
			return
		}
		changeNumberStr, err := utils.GetHeaderString(c, "ChangeNumber")
		if err != nil {
			utils.HandleError(c, 400, err)
			return
		}

		parseСhangeNumber, err := strconv.ParseInt(changeNumberStr, 10, 64)
		if err != nil {
			utils.HandleError(c, 400, err)
			return
		}

		account := models.Account{}
		if err := db.Where("id = ?", userID).First(&account).Error; err != nil {
			utils.HandleError(c, 400, err)
			return
		}

		if account.Balance > 0 {
			account.Balance -= int(parseСhangeNumber)
		} else {
			utils.HandleError(c, 400, fmt.Errorf("balance is zero"))
			return
		}
		if err := db.Save(&account).Error; err != nil {
			utils.HandleError(c, 400, err)
		}

		c.JSON(200, gin.H{"balance": account.Balance})
	}
}

func TransferBalance(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		tx := db.Begin()
		defer func() {
			if r := recover(); r != nil {
				tx.Rollback()
			}
		}()

		recipientID, err := utils.GetHeaderString(c, "RecipientID")
		if err != nil {
			tx.Rollback()
			utils.HandleError(c, 400, err)
			return
		}
		senderID, err := utils.GetHeaderString(c, "SenderID")
		if err != nil {
			tx.Rollback()
			utils.HandleError(c, 400, err)
			return
		}
		changeNumberStr, err := utils.GetHeaderString(c, "ChangeNumber")
		if err != nil {
			tx.Rollback()
			utils.HandleError(c, 400, err)
			return
		}

		parseChangeNumber, err := strconv.ParseInt(changeNumberStr, 10, 64)
		if err != nil {
			tx.Rollback()
			utils.HandleError(c, 400, err)
			return
		}

		var accountRecipient, accountSender models.Account
		if err := tx.Set("gorm:query_option", "FOR UPDATE").Where("id = ?", recipientID).First(&accountRecipient).Error; err != nil {
			tx.Rollback()
			utils.HandleError(c, 400, err)
			return
		}

		if err := tx.Set("gorm:query_option", "FOR UPDATE").Where("id = ?", senderID).First(&accountSender).Error; err != nil {
			tx.Rollback()
			utils.HandleError(c, 400, err)
			return
		}

		if accountSender.Balance < int(parseChangeNumber) {
			tx.Rollback()
			utils.HandleError(c, 400, fmt.Errorf("Insufficient balance"))
			return
		}

		accountSender.Balance -= int(parseChangeNumber)
		accountRecipient.Balance += int(parseChangeNumber)

		if err := tx.Save(&accountSender).Error; err != nil {
			tx.Rollback()
			utils.HandleError(c, 400, err)
			return
		}

		if err := tx.Save(&accountRecipient).Error; err != nil {
			tx.Rollback()
			utils.HandleError(c, 400, err)
			return
		}

		if err := tx.Commit().Error; err != nil {
			utils.HandleError(c, 500, err)
			return
		}

		c.JSON(200, gin.H{"SenderBalance": accountSender.Balance, "RecipientBalance": accountRecipient.Balance})
	}
}
