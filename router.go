package main

import (
	"embed"
	"net/http"

	"github.com/Allenxuxu/mogutouERP/api"
	"github.com/Allenxuxu/mogutouERP/middleware"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

//go:embed index.html  log.png  static
var static embed.FS

func initRouter() *gin.Engine {
	r := gin.New()

	r.StaticFS("/ui", http.FS(static))
	
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowAllOrigins = true
	corsConfig.AddAllowHeaders("Authorization")

	r.Use(cors.New(corsConfig))
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	rr := r.Group("/api/v1")
	rr.POST("/login", api.Login)

	authRouter := rr.Group("/")
	authRouter.Use(middleware.Auth)
	{
		authRouter.GET("/logout", api.Logout)

		authRouter.PATCH("/userPassword", api.UpdatePassword)
		authRouter.GET("/users", api.ListUsers)
		authRouter.GET("/user", api.GetUser)

		authRouter.GET("/commodities", api.ListCommodities)

		authRouter.POST("/order/custormer", api.CreateCustormerOrder)
		authRouter.GET("/orders/custormer", api.ListCustormerOrders)
		authRouter.DELETE("/order/custormer/:id", api.DeleteCustormerOrder)
		authRouter.PATCH("/order/custormer/:id/confirm", api.ConfirmCustormerOrder)

		//需要admin权限的API
		adminRouter := authRouter.Group("/")
		adminRouter.Use(middleware.Admin)

		adminRouter.POST("/order/purchase", api.CreatePurchaseOrder)
		adminRouter.GET("/orders/purchase", api.ListPurchaseOrders)
		adminRouter.DELETE("/order/purchase/:id", api.DeletePurchaseOrder)
		adminRouter.PATCH("/order/purchase/:id/confirm", api.ConfirmPurchaseOrder)

		adminRouter.POST("/user", api.CreateUser)
		adminRouter.DELETE("/user/:id", api.DeleteUser)
		adminRouter.PATCH("/user/:id", api.UpdateUser)

		adminRouter.POST("/commodity", api.CreateCommodity)
		adminRouter.PATCH("/commodity/:id", api.UpdateCommodity)
		adminRouter.DELETE("/commodity/:id", api.DeleteCommodity)
		adminRouter.GET("/admin/commodities", api.ListCommoditiesAsAdmin)
		adminRouter.GET("/admin/finance/:year", api.Finance)
	}

	return r
}
