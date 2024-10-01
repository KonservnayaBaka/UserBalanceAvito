package handlers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"userBalanceAvito/internal/models"
	"userBalanceAvito/internal/utils"
)

func AddProduct(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var product models.Product

		if err := c.ShouldBind(&product); err != nil {
			utils.HandleError(c, 400, err)
			return
		}

		if product.Name == "" {
			utils.HandleError(c, 400, fmt.Errorf("product name is empty"))
			return
		}

		if err := db.Create(&product).Error; err != nil {
			utils.HandleError(c, 400, err)
			return
		}

		c.JSON(200, gin.H{"product": product})
	}
}

func ReservProduct(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		accountID, err := utils.GetHeaderInt64(c, "AccountID")
		if err != nil {
			utils.HandleError(c, 400, err)
			return
		}
		productID, err := utils.GetHeaderInt64(c, "ProductID")
		if err != nil {
			utils.HandleError(c, 400, err)
			return
		}
		ammountCount, err := utils.GetHeaderInt64(c, "Ammount")
		if err != nil {
			utils.HandleError(c, 400, err)
			return
		}

		var reservation models.Reservation
		reservation.AccountID = uint(accountID)
		reservation.ProductID = uint(productID)
		reservation.Amount = int(ammountCount)

		if err := db.Create(&reservation).Error; err != nil {
			utils.HandleError(c, 400, err)
			return
		}

		c.JSON(200, gin.H{"reservation": reservation})
	}
}
