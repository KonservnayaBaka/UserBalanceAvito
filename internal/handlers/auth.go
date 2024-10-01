package handlers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"userBalanceAvito/internal/models"
	"userBalanceAvito/internal/utils"
)

func RegistrationNewUser(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var user models.User
		if err := c.ShouldBindJSON(&user); err != nil {
			utils.HandleError(c, 400, err)
			return
		}

		if user.Username == "" || user.Password == "" {
			utils.HandleError(c, 400, fmt.Errorf("username or password is empty"))
			return
		}

		hashedPassword, err := utils.HashPassword(user.Password)
		if err != nil {
			utils.HandleError(c, 500, fmt.Errorf("error while hashing password"))
			return
		}
		user.Password = hashedPassword

		if err := db.Create(&user).Error; err != nil {
			utils.HandleError(c, 400, err)
			return
		}

		account := models.Account{}
		if err := db.Create(&account).Error; err != nil {
			utils.HandleError(c, 400, err)
			return
		}

		c.JSON(200, gin.H{"user": user, "account": account})
	}
}

func AuthorizationUser(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var user models.User
		if err := c.ShouldBindJSON(&user); err != nil {
			utils.HandleError(c, 400, err)
			return
		}

		var dbUser models.User
		if err := db.Where("username = ?", user.Username).First(&dbUser).Error; err != nil {
			utils.HandleError(c, 400, fmt.Errorf("invalid username or password"))
			return
		}

		if !utils.CheckPasswordHash(user.Password, dbUser.Password) {
			utils.HandleError(c, 400, fmt.Errorf("invalid username or password"))
			return
		}

		c.JSON(200, gin.H{"user": dbUser, "status": true})
	}
}
