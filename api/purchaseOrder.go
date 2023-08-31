package api

import (
	"net/http"

	"github.com/Allenxuxu/mogutouERP/middleware"
	"github.com/Allenxuxu/mogutouERP/models"
	"github.com/Allenxuxu/mogutouERP/utils/response"
	"github.com/gin-gonic/gin"
)

// CreatePurchaseOrder 创建新的客户订单
func CreatePurchaseOrder(c *gin.Context) {
	userName := c.GetString(middleware.RequestUserNameKey)
	var data struct {
		Remarks string  `json:"remarks"`
		Amount  float32 `json:"amount" binding:"required"`
		Goods   []struct {
			ID     string `json:"id" binding:"required"`
			Number uint   `json:"number" binding:"required"`
		}
	}
	err := c.ShouldBindJSON(&data)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, response.Error{Error: err.Error()})
		return
	}
	if len(data.Goods) == 0 {
		c.AbortWithStatusJSON(http.StatusBadRequest, response.Error{Error: "没有商品"})
		return
	}

	goods := make([]models.PurchaseGoods, len(data.Goods))
	for i, v := range data.Goods {
		goods[i].GoodsID = v.ID
		goods[i].Number = v.Number
	}
	order := models.PurchaseOrder{
		Operator: userName,
		Remarks:  data.Remarks,
		Amount:   data.Amount,
	}
	orderInfo, err := models.CreatePurchaseOrder(&order, goods)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, response.Error{Error: err.Error()})
		return
	}

	c.JSON(http.StatusCreated, orderInfo)
}

// DeletePurchaseOrder 删除订单
func DeletePurchaseOrder(c *gin.Context) {
	orderID := c.Param("id")

	err := models.DeletePurchaseOrder(orderID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, response.Error{Error: err.Error()})
		return
	}

	c.JSON(http.StatusNoContent, gin.H{})
}

// ListPurchaseOrders 列举所有订单
func ListPurchaseOrders(c *gin.Context) {
	orders, err := models.ListPurchaseOrders()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, response.Error{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, orders)
}

// ConfirmPurchaseOrder 确认订单
func ConfirmPurchaseOrder(c *gin.Context) {
	orderID := c.Param("id")

	var data struct {
		Freight *float32 `json:"freight" binding:"required"`
	}
	err := c.ShouldBindJSON(&data)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, response.Error{Error: err.Error()})
		return
	}

	orderInfo, err := models.ConfirmPurchaseOrder(orderID, *data.Freight)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, response.Error{Error: err.Error()})
		return
	}

	c.JSON(http.StatusCreated, orderInfo)
}
