package routes

import (
	"project/delivery/handlers"
	m "project/delivery/middleware"

	"github.com/gin-gonic/gin"
)

func UserRouter(r *gin.Engine, userHandler *handlers.UserHandler) *gin.Engine {

	r.POST("/user/signup", userHandler.SignupWithOtp)
	r.POST("user/signup/otp", userHandler.SignupOtpValidation)

	r.POST("/user/login", userHandler.LoginWithPassword)
	r.POST("/user/address", m.UserRetreiveCookie, userHandler.AddAddress)
	r.PATCH("/user/address/:type", m.UserRetreiveCookie, userHandler.EditAddress)
	r.DELETE("/user/address/:type", m.UserRetreiveCookie, userHandler.DeleteAddress)
	r.GET("/user/details", m.UserRetreiveCookie, userHandler.ShowUserDetails)
	r.PATCH("/user/profile", m.UserRetreiveCookie, userHandler.EditProfile)
	r.POST("/user/changepassword", m.UserRetreiveCookie, userHandler.ChangePassword)
	r.POST("/user/changepassword/validation", m.UserRetreiveCookie, userHandler.OtpValidationPassword)

	r.GET("/user/products", m.UserRetreiveCookie, userHandler.Products)
	r.GET("/user/products/details/:productid", m.UserRetreiveCookie, userHandler.ProductDetails)
	r.GET("/user/products/search", m.UserRetreiveCookie, userHandler.SearchProduct)
	r.GET("/user/products/sort", m.UserRetreiveCookie, userHandler.SortByCategory)
	r.GET("/user/products/filter", m.UserRetreiveCookie, userHandler.SortByFilter)

	r.POST("/user/cart", m.UserRetreiveCookie, userHandler.AddToCart)
	r.DELETE("user/cart/:id", m.UserRetreiveCookie, userHandler.RemoveFromCart)
	r.GET("/user/cart", m.UserRetreiveCookie, userHandler.Cart)
	r.GET("/user/cartlist", m.UserRetreiveCookie, userHandler.CartItems)
	r.POST("/user/wishlist/:category/:productid", m.UserRetreiveCookie, userHandler.AddToWishList)
	r.DELETE("/user/wishlist/:product/:id", m.UserRetreiveCookie, userHandler.RemoveFromWishlist)
	r.GET("/user/wishlist", m.UserRetreiveCookie, userHandler.ViewWishlist)
	r.GET("/user/coupons", m.UserRetreiveCookie, userHandler.AvailableCoupons)
	r.POST("/user/cart/coupon", m.UserRetreiveCookie, userHandler.ApplyCoupon)
	r.POST("/logout", userHandler.Logout)
	return r
}
