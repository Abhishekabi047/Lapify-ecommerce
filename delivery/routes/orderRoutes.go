package routes

import (
	"project/delivery/handlers"
	m "project/delivery/middleware"

	"github.com/gin-gonic/gin"
)

func OrderRouter(r *gin.Engine, orderHandler *handlers.OrderHandler) *gin.Engine {
	r.POST("/user/placeorder", m.UserRetreiveCookie, orderHandler.PlaceOrder)
	r.POST("/paymentverification", m.UserRetreiveCookie, orderHandler.PaymentVerification)
	r.GET("/user/order", m.UserRetreiveCookie, orderHandler.OrderHistory)
	r.PATCH("/user/cancelorder/:orderid", m.UserRetreiveCookie, orderHandler.CancelOrder)

	r.PATCH("/admin/order/:orderid", m.AdminRetreiveToken, orderHandler.AdminOrderUpdate)
	r.GET("/admin/order", m.AdminRetreiveToken, orderHandler.AdminOrderDetails)
	r.PATCH("/admin/cancelorder/:orderid", m.AdminRetreiveToken, orderHandler.AdminCancelOrder)

	r.GET("/admin/salesreportbyperiod", m.AdminRetreiveToken, orderHandler.SalesReportByPeriod)
	r.GET("/admin/salesreportbydate", m.AdminRetreiveToken, orderHandler.SalesReportByDate)
	r.GET("/admin/salesreportbypayment", m.AdminRetreiveToken, orderHandler.SalesReportByPayment)
	
	r.GET("/user/printinvioce", m.UserRetreiveCookie, orderHandler.PrintInvoice)

	r.GET("/user/stripe", m.UserRetreiveCookie, orderHandler.ExecutePaymentStripe)
	r.POST("/webhook", orderHandler.HandleWebhook)
	r.GET("/user/orderstatus", m.UserRetreiveCookie, orderHandler.OrderStatus)
	return r
}
