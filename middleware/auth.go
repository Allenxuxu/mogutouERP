package middleware

import (
	"log"
	"net/http"

	"github.com/Allenxuxu/mogutouERP/pkg/token"
	"github.com/Allenxuxu/mogutouERP/utils/response"
	"github.com/gin-gonic/gin"
)

// RequestRoleKey context key值
const RequestRoleKey string = "mogutou-Role"

// RequestUserIDKey context key值
const RequestUserIDKey string = "mogutou-UserID"

// RequestUserNameKey context key值
const RequestUserNameKey string = "mogutou-UserName"

// Auth jwt鉴权
func Auth(c *gin.Context) {
	jwtToken := c.GetHeader("Authorization")

	info, err := token.Decode(jwtToken)
	if err != nil || info.PerAddr != c.ClientIP() {
		c.AbortWithStatusJSON(http.StatusUnauthorized, response.Error{Error: "非法访问"})
		return
	}

	c.Set(RequestRoleKey, info.Roles)
	c.Set(RequestUserIDKey, info.UserID)
	c.Set(RequestUserNameKey, info.UserName)
}

// Admin jwt鉴权
func Admin(c *gin.Context) {
	roles := c.GetStringSlice(RequestRoleKey)
	log.Println(roles)
	for _, v := range roles {
		if v == "admin" {
			return
		}
	}

	c.AbortWithStatusJSON(http.StatusBadRequest, response.Error{Error: "非法访问"})
}
