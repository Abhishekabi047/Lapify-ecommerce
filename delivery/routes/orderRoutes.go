package routes

import (
	"project/delivery/handlers"
	m "project/delivery/middleware"

	"github.com/gin-gonic/gin"
)

func OrderRouter(r *gin.Engine, orderHandler *handlers.OrderHandler) *gin.Engine {
	r.POST("/user/order/place", m.UserRetreiveCookie, orderHandler.PlaceOrder)
	r.POST("/user/payment/verify", m.UserRetreiveCookie, orderHandler.PaymentVerification)
	r.GET("/user/order/history", m.UserRetreiveCookie, orderHandler.OrderHistory)
	r.PATCH("/user/order/cancel/:orderid", m.UserRetreiveCookie, orderHandler.CancelOrder)

	r.PATCH("/admin/order/update/:orderid", m.AdminRetreiveToken, orderHandler.AdminOrderUpdate)
	r.GET("/admin/order/details", m.AdminRetreiveToken, orderHandler.AdminOrderDetails)
	r.PATCH("/admin/order/cancel/:orderid", m.AdminRetreiveToken, orderHandler.AdminCancelOrder)

	r.GET("/admin/salesreport/period/:period", m.AdminRetreiveToken, orderHandler.SalesReportByPeriod)
	r.GET("/admin/salesreport/date/:start/:end", m.AdminRetreiveToken, orderHandler.SalesReportByDate)
	r.GET("/admin/salesreport/payment/:start/:end/:payment", m.AdminRetreiveToken, orderHandler.SalesReportByPayment)

	r.GET("/user/order/invoice", m.UserRetreiveCookie, orderHandler.PrintInvoice)

	r.GET("/user/stripe", m.UserRetreiveCookie, orderHandler.ExecutePaymentStripe)
	r.POST("/webhook", orderHandler.HandleWebhook)
	r.GET("/user/orderstatus", m.UserRetreiveCookie, orderHandler.OrderStatus)
	return r
}
