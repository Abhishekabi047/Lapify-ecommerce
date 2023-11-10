package routes

import (
	"project/delivery/handlers"
	m "project/delivery/middleware"

	"github.com/gin-gonic/gin"
)

func UserRouter(r *gin.Engine, userHandler *handlers.UserHandler) *gin.Engine {

	
	r.POST("/signupwithotp", userHandler.SignupWithOtp)
	r.POST("/signupotpvalidation", userHandler.SignupOtpValidation)
	
	r.POST("/loginwithpassword", userHandler.LoginWithPassword)

	r.GET("/products", m.UserRetreiveCookie, userHandler.Products)
	r.GET("/productdetails/:productid", m.UserRetreiveCookie, userHandler.ProductDetails)

	return r
}
