package routes

import (
	"project/delivery/handlers"

	m "project/delivery/middleware"

	"github.com/gin-gonic/gin"
)

func AdminRouter(r *gin.Engine, adminHandler *handlers.AdminHandler) *gin.Engine {
	r.POST("/admin/login", adminHandler.AdminLoginWithPassword)
	r.GET("/admin/home", m.AdminRetreiveToken, adminHandler.Home)

	r.GET("admin/users", m.AdminRetreiveToken, adminHandler.UsersList)
	r.PUT("/admin/users/toggle-permission/:id", m.AdminRetreiveToken, adminHandler.TogglePermission)
	r.GET("admin/search/users", m.AdminRetreiveToken, adminHandler.SearchUsers)

	r.POST("/admin/categories", m.AdminRetreiveToken, adminHandler.CreateCategory)
	r.PUT("/admin/categories/:id", m.AdminRetreiveToken, adminHandler.EditCategory)
	r.DELETE("/admin/categories/:id", m.AdminRetreiveToken, adminHandler.DeleteCategory)

	r.GET("/admin/products", m.AdminRetreiveToken, adminHandler.AdminProductlist)
	r.POST("/admin/products", m.AdminRetreiveToken, adminHandler.CreateProduct)

	r.PATCH("/admin/products/:id", m.AdminRetreiveToken, adminHandler.EditProduct)
	r.DELETE("admin/products/:id", m.AdminRetreiveToken, adminHandler.DeleteProduct)

	r.POST("/admin/coupons", m.AdminRetreiveToken, adminHandler.AddCoupon)
	r.GET("/admin/coupons", m.AdminRetreiveToken, adminHandler.AllCoupons)
	r.DELETE("/admin/coupons", m.AdminRetreiveToken, adminHandler.DeleteCoupon)

	// r.POST("/admin/offer", m.AdminRetreiveToken, adminHandler.AddOffer)
	// r.GET("/admin/offer", m.AdminRetreiveToken, adminHandler.AllOffer)

	r.GET("/admin/stockless/products", m.AdminRetreiveToken, adminHandler.StocklessProducts)

	r.POST("/admin/product/offer", m.AdminRetreiveToken, adminHandler.AddProductOffer)
	r.POST("/admin/category/offer", m.AdminRetreiveToken, adminHandler.AddCategoryOffer)
	return r
}
