package api

import (
	"net/http"

	"github.com/Allenxuxu/mogutouERP/models"
	"github.com/Allenxuxu/mogutouERP/utils/response"
	"github.com/gin-gonic/gin"
)

// CreateCommodity 创建新商品类型
func CreateCommodity(c *gin.Context) {
	var data struct {
		ID            string  `json:"id" binding:"required"`
		Name          string  `json:"name" binding:"required"`
		Colour        string  `json:"colour" binding:"required"`
		Size          string  `json:"size" binding:"required"`
		Brand         string  `json:"brand" binding:"required"`
		Price         float32 `json:"price"`
		PurchasePrice float32 `json:"purchase_price"`
	}
	err := c.ShouldBindJSON(&data)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, response.Error{Error: err.Error()})
		return
	}

	commodity := models.Commodity{
		ID:            data.ID,
		Name:          data.Name,
		Colour:        data.Colour,
		Size:          data.Size,
		Brand:         data.Brand,
		Price:         data.Price,
		PurchasePrice: data.PurchasePrice,
	}
	err = models.CreateCommodity(&commodity)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, response.Error{Error: err.Error()})
		return
	}

	c.JSON(http.StatusCreated, models.GetCommodity(data.ID))
}

// UpdateCommodity 更新商品信息
func UpdateCommodity(c *gin.Context) {
	commodityID := c.Param("id")

	if !models.HaveCommodity(commodityID) {
		c.AbortWithStatusJSON(http.StatusNotFound, response.Error{Error: "没有此商品信息"})
		return
	}

	var data struct {
		Name          string  `json:"name"`
		Colour        string  `json:"colour"`
		Size          string  `json:"size"`
		Brand         string  `json:"brand"`
		Price         float32 `json:"price"`
		PurchasePrice float32 `json:"purchase_price"`
	}
	err := c.ShouldBindJSON(&data)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, response.Error{Error: err.Error()})
		return
	}

	commodity := models.Commodity{
		Name:          data.Name,
		Colour:        data.Colour,
		Size:          data.Size,
		Brand:         data.Brand,
		Price:         data.Price,
		PurchasePrice: data.PurchasePrice,
	}
	err = models.UpdateCommodity(commodityID, &commodity)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, response.Error{Error: err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{})
}

// DeleteCommodity 删除商品
func DeleteCommodity(c *gin.Context) {
	commodityID := c.Param("id")

	if !models.HaveCommodity(commodityID) {
		c.AbortWithStatusJSON(http.StatusNotFound, response.Error{Error: "没有此商品信息"})
		return
	}

	err := models.DeleteCommodity(commodityID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, response.Error{Error: err.Error()})
		return
	}

	c.JSON(http.StatusNoContent, gin.H{})
}

// ListCommodities 列举所有商品
func ListCommodities(c *gin.Context) {
	commodities, err := models.GetCommodities()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, response.Error{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, commodities)
}

// ListCommoditiesAsAdmin list所有商品信息
func ListCommoditiesAsAdmin(c *gin.Context) {
	commodities, err := models.GetCommoditiesWithPurchasePrice()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, response.Error{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, commodities)
}
