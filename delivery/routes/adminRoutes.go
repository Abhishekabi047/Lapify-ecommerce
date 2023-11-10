package routes

import (
	"project/delivery/handlers"

	m "project/delivery/middleware"

	"github.com/gin-gonic/gin"
)

func AdminRouter(r *gin.Engine, adminHandler *handlers.AdminHandler) *gin.Engine {
	r.POST("/adminloginpassword", adminHandler.AdminLoginWithPassword)
	

	r.GET("/usermanagement", m.AdminRetreiveToken, adminHandler.UsersList)
	r.POST("/userpermission/:id", m.AdminRetreiveToken, adminHandler.TogglePermission)

	r.POST("/addcategory", m.AdminRetreiveToken, adminHandler.CreateCategory)
	r.PUT("/editcategory:id", m.AdminRetreiveToken, adminHandler.EditCategory)
	r.DELETE("/deletecategory/:id", m.AdminRetreiveToken, adminHandler.DeleteCategory)

	r.POST("/addproduct", m.AdminRetreiveToken, adminHandler.CreateProduct)
	r.PUT("/editproduct/:id", m.AdminRetreiveToken, adminHandler.EditProduct)
	r.DELETE("/deleteproduct/:id", m.AdminRetreiveToken, adminHandler.DeleteProduct)
	return r
}
