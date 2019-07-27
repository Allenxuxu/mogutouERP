package api

import (
	"net/http"
	"strconv"

	"github.com/Allenxuxu/mogutouERP/models"
	"github.com/Allenxuxu/mogutouERP/utils/response"
	"github.com/gin-gonic/gin"
)

// Finance 财务状况
func Finance(c *gin.Context) {
	year, err := strconv.Atoi(c.Param("year"))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, response.Error{Error: err.Error()})
		return
	}
	profit, turnover, orderQuantity, err := models.QueryYearFinance(year)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, response.Error{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"profit":        profit,
		"turnover":      turnover,
		"orderQuantity": orderQuantity,
	})
}
