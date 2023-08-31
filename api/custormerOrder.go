package api

import (
	"net/http"

	"github.com/Allenxuxu/mogutouERP/middleware"
	"github.com/Allenxuxu/mogutouERP/models"
	"github.com/Allenxuxu/mogutouERP/utils/response"
	"github.com/gin-gonic/gin"
)

// CreateCustormerOrder 创建新的客户订单
func CreateCustormerOrder(c *gin.Context) {
	userName := c.GetString(middleware.RequestUserNameKey)
	var data struct {
		Remarks         string  `json:"remarks"`
		Amount          float32 `json:"amount" binding:"required"`
		Name            string  `json:"name" binding:"required"`
		Tel             string  `json:"tel" binding:"required"`
		DeliveryAddress string  `json:"deliveryAddress" binding:"required"`
		DeliveryTime    string  `json:"deliveryTime" binding:"required"`
		Deposit         float32 `json:"deposit" binding:"required"`
		Goods           []struct {
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

	goods := make([]models.CustormerGoods, len(data.Goods))
	for i, v := range data.Goods {
		goods[i].GoodsID = v.ID
		goods[i].Number = v.Number
	}
	order := models.CustormerOrder{
		Operator:        userName,
		Remarks:         data.Remarks,
		Amount:          data.Amount,
		Name:            data.Name,
		Tel:             data.Tel,
		DeliveryAddress: data.DeliveryAddress,
		DeliveryTime:    data.DeliveryTime,
		Deposit:         data.Deposit,
	}

	orderInfo, err := models.CreateCustormerOrder(&order, goods)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, response.Error{Error: err.Error()})
		return
	}

	c.JSON(http.StatusCreated, orderInfo)
}

// DeleteCustormerOrder 删除订单
func DeleteCustormerOrder(c *gin.Context) {
	orderID := c.Param("id")

	err := models.DeleteCustormerOrder(orderID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, response.Error{Error: err.Error()})
		return
	}

	c.JSON(http.StatusNoContent, gin.H{})
}

// ListCustormerOrders 列举所有订单
func ListCustormerOrders(c *gin.Context) {
	orders, err := models.ListCustormerOrders()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, response.Error{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, orders)
}

// ConfirmCustormerOrder 确认订单
func ConfirmCustormerOrder(c *gin.Context) {
	orderID := c.Param("id")

	var data struct {
		Freight *float32 `json:"freight" binding:"required"`
	}
	err := c.ShouldBindJSON(&data)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, response.Error{Error: err.Error()})
		return
	}

	orderInfo, err := models.ConfirmCustormerOrder(orderID, *data.Freight)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, response.Error{Error: err.Error()})
		return
	}

	c.JSON(http.StatusCreated, orderInfo)
}
