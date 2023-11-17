package routes

import (
	"project/delivery/handlers"

	m "project/delivery/middleware"

	"github.com/gin-gonic/gin"
)

func AdminRouter(r *gin.Engine, adminHandler *handlers.AdminHandler) *gin.Engine {
	r.POST("/admin/login", adminHandler.AdminLoginWithPassword)

	r.GET("admin/usermanagement", m.AdminRetreiveToken, adminHandler.UsersList)
	r.POST("admin/usermanagement/:id", m.AdminRetreiveToken, adminHandler.TogglePermission)

	r.POST("/admin/category", m.AdminRetreiveToken, adminHandler.CreateCategory)
	r.PATCH("/admin/category/:id", m.AdminRetreiveToken, adminHandler.EditCategory)
	r.DELETE("/admin/category/:id", m.AdminRetreiveToken, adminHandler.DeleteCategory)

	r.GET("/admin/products", m.AdminRetreiveToken, adminHandler.AdminProductlist)
	r.POST("/admin/product", m.AdminRetreiveToken, adminHandler.CreateProduct)
	r.PATCH("/admin/product/:id", m.AdminRetreiveToken, adminHandler.EditProduct)
	r.DELETE("admin/product/:id", m.AdminRetreiveToken, adminHandler.DeleteProduct)
	return r
}
