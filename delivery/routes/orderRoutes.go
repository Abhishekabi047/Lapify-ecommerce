package routes

import (
	"project/delivery/handlers"
	m "project/delivery/middleware"

	"github.com/gin-gonic/gin"
)

func OrderRouter(r *gin.Engine, orderHandler *handlers.OrderHandler) *gin.Engine {
	r.POST("/user/placeorder/:addressid", m.UserRetreiveCookie, orderHandler.PlaceOrder)
	r.GET("/user/order", m.UserRetreiveCookie, orderHandler.OrderHistory)
	r.PATCH("/user/cancelorder/:orderid", m.UserRetreiveCookie, orderHandler.CancelOrder)

	r.PATCH("/admin/order/:orderid/:status", m.AdminRetreiveToken, orderHandler.AdminOrderUpdate)
	r.GET("/admin/order", m.AdminRetreiveToken, orderHandler.AdminOrderDetails)
	r.PATCH("/admin/cancelorder/:orderid", m.AdminRetreiveToken, orderHandler.AdminCancelOrder)

	return r
}
