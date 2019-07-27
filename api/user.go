package api

import (
	"log"
	"net/http"
	"time"

	"github.com/Allenxuxu/mogutouERP/middleware"
	"github.com/Allenxuxu/mogutouERP/models"
	"github.com/Allenxuxu/mogutouERP/pkg/token"
	"github.com/Allenxuxu/mogutouERP/utils/response"
	"github.com/gin-gonic/gin"
)

// Login 登陆
func Login(c *gin.Context) {
	var data struct {
		Tel      string `json:"tel" binding:"required"`
		Password string `json:"password" binding:"required"`
	}
	err := c.ShouldBindJSON(&data)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, response.Error{Error: err.Error()})
		return
	}

	userInfo, err := models.VerifyUser(data.Tel, data.Password)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, response.Error{Error: err.Error()})
		return
	}

	role, err := models.GetUserRole(userInfo.UserID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, response.Error{Error: err.Error()})
		return
	}
	clientIP := c.ClientIP()

	tokenStr, err := token.Encode(userInfo.Name, userInfo.UserID, clientIP, role, "mogutou", time.Now().Add(time.Hour*24*3).Unix())
	if err != nil {
		log.Fatal(err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, response.Error{Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"token": tokenStr,
	})
}

// Logout 登出
func Logout(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{})
}

// CreateUser 创建一个普通员工账号
func CreateUser(c *gin.Context) {
	var data struct {
		Name     string `json:"name" binding:"required"`
		Tel      string `json:"tel" binding:"required"`
		Position string `json:"position" binding:"required"`
	}
	err := c.ShouldBindJSON(&data)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, response.Error{Error: err.Error()})
		return
	}
	if models.HaveTel(data.Tel) {
		c.AbortWithStatusJSON(http.StatusBadRequest, response.Error{Error: "账号已存在"})
		return
	}

	user := models.User{
		Name:     data.Name,
		Tel:      data.Tel,
		Position: data.Position,
	}

	err = models.CreateUser(&user)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, response.Error{Error: err.Error()})
		return
	}

	c.JSON(http.StatusCreated, models.GetUser(user.UserID))
}

// DeleteUser 删除员工信息
func DeleteUser(c *gin.Context) {
	userID := c.Param("id")
	if !models.HaveUser(userID) {
		c.AbortWithStatusJSON(http.StatusNotFound, response.Error{Error: "没有此用户信息"})
		return
	}

	roles, err := models.GetUserRole(userID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, response.Error{Error: err.Error()})
		return
	}
	for _, v := range roles {
		if v == "admin" {
			c.AbortWithStatusJSON(http.StatusBadRequest, response.Error{Error: "无法删除"})
			return
		}
	}

	err = models.DeleteUser(userID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, response.Error{Error: err.Error()})
		return
	}

	c.JSON(http.StatusNoContent, gin.H{})
}

// UpdateUser 更新用户信息
func UpdateUser(c *gin.Context) {
	userID := c.Param("id")
	if !models.HaveUser(userID) {
		c.AbortWithStatusJSON(http.StatusNotFound, response.Error{Error: "没有此用户信息"})
		return
	}

	var data struct {
		Name     string `json:"name"`
		Tel      string `json:"tel"`
		Password string `json:"password"`
		Position string `json:"position"`
	}
	err := c.ShouldBindJSON(&data)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, response.Error{Error: err.Error()})
		return
	}

	err = models.UpdateUser(userID, data.Name, data.Tel, data.Password, data.Position)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, response.Error{Error: err.Error()})
		return
	}

	c.JSON(http.StatusCreated, models.GetUser(userID))
}

// UpdatePassword 修改密码
func UpdatePassword(c *gin.Context) {
	userID := c.GetString(middleware.RequestUserIDKey)

	if !models.HaveUser(userID) {
		c.AbortWithStatusJSON(http.StatusNotFound, response.Error{Error: "没有此用户信息"})
		return
	}

	var data struct {
		Password string `json:"password" binding:"required"`
	}
	err := c.ShouldBindJSON(&data)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, response.Error{Error: err.Error()})
		return
	}

	err = models.UpdatePassword(userID, data.Password)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, response.Error{Error: err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{})
}

// GetUser 获取用户信息
func GetUser(c *gin.Context) {
	userID := c.GetString(middleware.RequestUserIDKey)

	if !models.HaveUser(userID) {
		c.AbortWithStatusJSON(http.StatusBadRequest, response.Error{Error: "没有此用户信息"})
		return
	}

	user := models.GetUser(userID)
	roles, err := models.GetUserRole(userID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, response.Error{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"name":     user.Name,
		"tel":      user.Tel,
		"position": user.Position,
		"roles":    roles,
	})

}

// ListUsers 列举所有用户信息
func ListUsers(c *gin.Context) {
	users, err := models.ListUsers()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, response.Error{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, users)
}
