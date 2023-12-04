package routes

import (
	"project/delivery/handlers"

	m "project/delivery/middleware"

	"github.com/gin-gonic/gin"
)

func AdminRouter(r *gin.Engine, adminHandler *handlers.AdminHandler) *gin.Engine {
	r.POST("/admin/login", adminHandler.AdminLoginWithPassword)
	r.GET("/admin/home", m.AdminRetreiveToken, adminHandler.Home)

	r.GET("admin/usermanagement", m.AdminRetreiveToken, adminHandler.UsersList)
	r.POST("admin/usermanagement/:id", m.AdminRetreiveToken, adminHandler.TogglePermission)

	r.POST("/admin/category", m.AdminRetreiveToken, adminHandler.CreateCategory)
	r.PATCH("/admin/category/:id", m.AdminRetreiveToken, adminHandler.EditCategory)
	r.DELETE("/admin/category/:id", m.AdminRetreiveToken, adminHandler.DeleteCategory)

	r.GET("/admin/products", m.AdminRetreiveToken, adminHandler.AdminProductlist)
	r.POST("/admin/product", m.AdminRetreiveToken, adminHandler.CreateProduct)

	r.PATCH("/admin/product/:id", m.AdminRetreiveToken, adminHandler.EditProduct)
	r.DELETE("admin/product/:id", m.AdminRetreiveToken, adminHandler.DeleteProduct)

	r.POST("/admin/coupon", m.AdminRetreiveToken, adminHandler.AddCoupon)
	r.GET("/admin/coupon", m.AdminRetreiveToken, adminHandler.AllCoupons)
	r.DELETE("/admin/coupon", m.AdminRetreiveToken, adminHandler.DeleteCoupon)

	r.POST("/admin/offer", m.AdminRetreiveToken, adminHandler.AddOffer)
	r.GET("/admin/offer", m.AdminRetreiveToken, adminHandler.AllOffer)

	r.POST("/admin/product/offer", m.AdminRetreiveToken, adminHandler.AddProductOffer)
	r.POST("/admin/category/offer", m.AdminRetreiveToken, adminHandler.AddCategoryOffer)
	return r
}
