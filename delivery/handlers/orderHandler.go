package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"project/config"
	"project/domain/entity"
	usecase "project/usecase/order"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/paymentintent"
	"github.com/stripe/stripe-go/webhook"
)

type OrderHandler struct {
	OrderUseCase *usecase.OrderUseCase
	Razorpay     config.Razopay
}

func NewOrderHandler(OrderUseCase *usecase.OrderUseCase, Razorpay config.Razopay) *OrderHandler {
	return &OrderHandler{OrderUseCase, Razorpay}
}

// PlaceOrder godoc
// @Summary Place an order
// @Description Places an order for the authenticated user based on the selected payment method.
// @ID place-order
// @Accept json
// @Tags User Orders
// @Produce json
// @Param addressid path int true "Address ID for the order"
// @Param payment path string true "Payment method ('cod', 'razorpay', 'wallet')"
// @Success 200 {string} string "Invoice details" "Successful response for COD payment"
// @Success 200 {string} string "Complete your Razorpay payment through. Razorpay ID: {razorId}, Order ID: {orderid}, User ID: {userid}" "Successful response for Razorpay payment"
// @Success 200 {string} string "Invoice details" "Successful response for Wallet payment"
// @Failure 400 {string} string "Bad request"
// @Router /user/order/place/{addressid}/{payment} [post]
func (oh *OrderHandler) PlaceOrder(c *gin.Context) {
	userID, _ := c.Get("userId")
	userid := userID.(int)
	straddress := c.Param("addressid")
	PaymentMethod := c.Param("payment")
	addressId, err := strconv.Atoi(straddress)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "strng conversion failed"})
		return
	}
	if PaymentMethod == "cod" {
		invoice, err := oh.OrderUseCase.ExecuteOrderCod(userid, addressId)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		} else {
			c.JSON(http.StatusOK, gin.H{"invoice": invoice})
		}
	} else if PaymentMethod == "razorpay" {
		razorId, orderId, err1 := oh.OrderUseCase.ExecuteRazorPay(userid, addressId)
		if err1 != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err1.Error()})
			return
		}
		c.HTML(http.StatusOK, "razor.html", gin.H{
			"message": "make payment",
			"razorId": razorId,
			"orderid": orderId})
		// c.JSON(http.StatusOK, gin.H{"message": "complete your razor pay through ", "razorId": razorId, "orderid": orderId, "userid": userid})
	} else if PaymentMethod == "wallet" {
		invoice, err2 := oh.OrderUseCase.ExecutePaymentWallet(userid, addressId)
		if err2 != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err2.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"Invoice": invoice})
	}
}

// PaymentVerification godoc
// @Summary Verify payment for Razorpay
// @Description Verifies the payment for Razorpay based on the provided signature, Razorpay ID, and payment ID.
// @ID verify-payment-razorpay
// @Accept multipart/form-data
// @Tags User Orders
// @Produce json
// @Param sign formData string true "Signature for payment verification"
// @Param razorid formData string true "Razorpay ID"
// @Param paymentid formData string true "Payment ID"
// @Success 200 {string} string "Payment successful. Invoice details: {invoice}"
// @Failure 400 {string} string "Bad request"
// @Router /order/payment/verify [post]
func (co *OrderHandler) PaymentVerification(c *gin.Context) {
	Signature := c.PostForm("sign")
	razorId := c.PostForm("razorid")
	paymentId := c.PostForm("paymentid")
	invoice, err1 := co.OrderUseCase.ExecuteRazorPayVerification(Signature, razorId, paymentId)
	if err1 != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err1.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "payment succesful", "invoice": invoice})
}

// CancelOrder godoc
// @Summary Cancel an order
// @Description Cancels an order based on the provided order ID.
// @ID cancel-order
// @Tags User Orders
// @Produce json
// @Param orderid path int true "Order ID to be canceled"
// @Success 200 {string} string "Order canceled successfully"
// @Failure 400 {string} string "Bad request"
// @Router /user/order/cancel/{orderid} [patch]
func (co *OrderHandler) CancelOrder(c *gin.Context) {
	strorderId := c.Param("orderid")
	orderid, err := strconv.Atoi(strorderId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "string convertion failed"})
		return
	}
	err1 := co.OrderUseCase.ExecuteCancelOrder(orderid)
	if err1 != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err1.Error()})
		return
	}
	updatedorder, err2 := co.OrderUseCase.UpdatedUser(orderid)
	if err2 != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err2.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"order cancceled ": updatedorder})
}

// OrderHistory godoc
// @Summary Retrieve order history for the authenticated user
// @Description Retrieves the order history for the authenticated user based on pagination parameters.
// @ID get-order-history
// @Tags User Orders
// @Produce json
// @Param page query int false "Page number for pagination (default is 1)"
// @Param limit query int false "Number of items per page (default is 5)"
// @Success 200 {string} string "Order history retrieved successfully"
// @Failure 400 {string} string "Bad request"
// @Router /user/order/history [get]
func (co *OrderHandler) OrderHistory(c *gin.Context) {
	userID, _ := c.Get("userId")
	userid := userID.(int)
	strpage := c.DefaultQuery("page", "1")
	page, err := strconv.Atoi(strpage)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "string convertion failed"})
		return
	}
	limitstr := c.DefaultQuery("limit", "5")
	limit, err := strconv.Atoi(limitstr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "string convertion failed"})
		return
	}
	orderlist, err := co.OrderUseCase.ExecuteOrderHistory(userid, page, limit)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"orderlist": orderlist})
}

// AdminOrderUpdate godoc
// @Summary Update order status (Admin)
// @Description Updates the status of an order based on the provided order ID and status (for admin use).
// @ID admin-update-order
// @Tags Admin Orders
// @Produce json
// @Param orderid path int true "Order ID to be updated"
// @Param status formData string true "New status for the order"
// @Success 200 {string} string "Order updated successfully. Updated order status: {updated order status}"
// @Failure 400 {string} string "Bad request"
// @Router /admin/order/update/{orderid} [patch]
func (op *OrderHandler) AdminOrderUpdate(c *gin.Context) {
	strorderid := c.Param("orderid")
	orderid, err := strconv.Atoi(strorderid)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "strign conversion failed"})
		return
	}
	status := c.PostForm("status")
	fmt.Println("status", status)
	err1 := op.OrderUseCase.ExecuteOrderUpdate(orderid, status)
	if err1 != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err1.Error()})
		return
	}
	orderstatus, err2 := op.OrderUseCase.UpdatedUser(orderid)
	if err2 != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err2.Error()})
	}
	c.JSON(http.StatusOK, gin.H{"success": "order updated", "updated order": orderstatus})
}

// AdminOrderDetails godoc
// @Summary Retrieve order details for admin
// @Description Retrieves the order details for admin based on pagination parameters.
// @ID get-admin-order-details
// @Tags Admin Orders
// @Produce json
// @Param page query int false "Page number for pagination (default is 1)"
// @Param limit query int false "Number of items per page (default is 5)"
// @Success 200 {string} string "Order details retrieved successfully"
// @Failure 400 {string} string "Bad request"
// @Router /admin/order/details [get]
func (op *OrderHandler) AdminOrderDetails(c *gin.Context) {
	strpage := c.DefaultQuery("page", "1")
	page, err := strconv.Atoi(strpage)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "string convertion failed"})
		return
	}
	strlimit := c.DefaultQuery("limit", "5")
	limit, err := strconv.Atoi(strlimit)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "string convertion failed"})
		return
	}
	orderlist, err1 := op.OrderUseCase.ExecuteAdminOrder(page, limit)
	if err1 != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err1.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"orderlist": orderlist})
}

// AdminCancelOrder godoc
// @Summary Cancel order (Admin)
// @Description Cancels an order based on the provided order ID (for admin use).
// @ID admin-cancel-order
// @Tags Admin Orders
// @Produce json
// @Param orderid path int true "Order ID to be canceled"
// @Success 200 {string} string "Order cancelled successfully"
// @Failure 400 {string} string "Bad request"
// @Router /admin/order/cancel/{orderid} [patch]
func (op *OrderHandler) AdminCancelOrder(c *gin.Context) {
	strOrderid := c.Param("orderid")
	orderid, err := strconv.Atoi(strOrderid)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "string covertion failed"})
		return
	}
	err1 := op.OrderUseCase.ExecuteAdminCancelOrder(orderid)
	if err1 != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err1.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "order cancelled"})
}

// SalesReportByDate godoc
// @Summary Generate sales report by date range
// @Description Generates a sales report based on the provided start and end dates.
// @ID sales-report-by-date
// @Tags Admin Report
// @Produce json
// @Param start path string true "Start date for the report (format: 2-1-2006)"
// @Param end path string true "End date for the report (format: 2-1-2006)"
// @Success 200 {string} string "Sales report generated successfully"
// @Failure 400 {string} string "Bad request"
// @Router /admin/salesreport/date/{start}/{end} [get]
func (or *OrderHandler) SalesReportByDate(c *gin.Context) {
	startDateStr := c.Param("start")
	endDateStr := c.Param("end")
	startDate, err := time.Parse("2-1-2006", startDateStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	endDate, err := time.Parse("2-1-2006", endDateStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	report, err := or.OrderUseCase.ExecuteSalesReportByDate(startDate, endDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"report": report})
}

// SalesReportByPeriod godoc
// @Summary Generate sales report by period
// @Description Generates a sales report based on the provided period.
// @ID sales-report-by-period
// @Tags Admin Report
// @Produce json
// @Param period path string true "Period for the report (e.g., 'monthly', 'quarterly', 'yearly')"
// @Success 200 {string} string "Sales report generated successfully"
// @Failure 400 {string} string "Bad request"
// @Router /admin/salesreport/period/{period} [get]
func (or *OrderHandler) SalesReportByPeriod(c *gin.Context) {
	period := c.Param("period")

	report, err := or.OrderUseCase.ExecuteSalesReportByPeriod(period)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
	c.JSON(http.StatusOK, gin.H{"report": report})
}

// SalesReportByPayment godoc
// @Summary Generate sales report by payment method and date range
// @Description Generates a sales report based on the provided start and end dates and payment method.
// @ID sales-report-by-payment
// @Tags Admin Report
// @Produce json
// @Param start path string true "Start date for the report (format: 2-1-2006)"
// @Param end path string true "End date for the report (format: 2-1-2006)"
// @Param paymentmethod path string true "Payment method for the report"
// @Success 200 {string} string "Sales report generated successfully"
// @Failure 400 {string} string "Bad request"
// @Router /admin/salesreport/payment/{start}/{end}/{paymentmethod} [get]
func (or *OrderHandler) SalesReportByPayment(c *gin.Context) {
	startDateStr := c.Param("start")
	endDateStr := c.Param("end")
	paymentmethod := c.Param("paymentmethod")
	startDate, err := time.Parse("2-1-2006", startDateStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
	endDate, err := time.Parse("2-1-2006", endDateStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
	report, err := or.OrderUseCase.ExecuteSalesReportByPaymentMethod(startDate, endDate, paymentmethod)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
	c.JSON(http.StatusOK, gin.H{"report": report})
}

func (or *OrderHandler) ExecutePaymentStripe(c *gin.Context) {
	UserID, _ := c.Get("userId")
	userId := UserID.(int)
	straddress := c.PostForm("address")
	addresid, err := strconv.Atoi(straddress)

	stripe.Key = "sk_test_51OFC9hSJxogb8Is5XfYeIuKpqDOMKzH7NPVdwTZDVu0I6wc0sOX4CCZ66scJRKM7iYemPXk2D5fvRLKGrHFe60OF00psv6EYzW"

	orders, err := or.OrderUseCase.ExecuteCartit(userId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	params := &stripe.PaymentIntentParams{
		Params: stripe.Params{
			Metadata: map[string]string{
				"user_id":    strconv.Itoa(userId),
				"address_id": strconv.Itoa(addresid),
			},
		},
		Amount:   stripe.Int64(int64(orders.TotalPrize)),
		Currency: stripe.String("INR"),
	}

	intent, err := paymentintent.New(params)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.HTML(http.StatusOK, "stripe.html", gin.H{
		"ClientSecret": intent.ClientSecret,
	})

}

var invoiceMap = make(map[string]*entity.Invoice)

func (cr *OrderHandler) HandleWebhook(c *gin.Context) {
	var Resul *entity.Invoice

	const MaxBodyBytes = int64(65536)
	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, MaxBodyBytes)
	payload, err := io.ReadAll(c.Request.Body)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading request body: %v\n", err)
		c.AbortWithStatus(http.StatusServiceUnavailable)
		return
	}

	endpointSecret := "whsec_907f81648d3e7562e20257062f4a8a175b9584e0eb9d2c97f772d99237824749"

	event, err := webhook.ConstructEvent(payload, c.Request.Header.Get("Stripe-Signature"), endpointSecret)

	if err != nil {

		fmt.Fprintf(os.Stderr, "Error verifying webhook signature: %v\n", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	switch event.Type {
	case "payment_intent.created":
		var paymentIntent stripe.PaymentIntent
		err := json.Unmarshal(event.Data.Raw, &paymentIntent)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error parsing webhook JSON: %v\n", err)
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}
		userID, _ := paymentIntent.Metadata["user_id"]
		addressID, _ := paymentIntent.Metadata["address_id"]

		userid, _ := strconv.Atoi(userID)
		addressid, _ := strconv.Atoi(addressID)

		result, err := cr.OrderUseCase.ExecuteInvoiceStripe(userid, addressid)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"errror": "invoice creation failed"})
			return
		}
		invoiceMap[paymentIntent.ID] = result
		Resul = result
		fmt.Println("Invoice Created successfully")
		fmt.Println("res", Resul.OrderId)

	case "payment_intent.succeeded":
		var paymentIntent stripe.PaymentIntent
		err := json.Unmarshal(event.Data.Raw, &paymentIntent)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error parsing webhook JSON: %v\n", err)
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}
		Resul = invoiceMap[paymentIntent.ID]

		fmt.Println("res", Resul.OrderId)
		if paymentIntent.Status == "succeeded" {
			if Resul != nil {

				err := cr.OrderUseCase.UpdateInvoiceStatus(int(Resul.OrderId), "succesfull")
				if err != nil {
					fmt.Fprintf(os.Stderr, "Error updating invoice status: %v\n", err)

				}
				fmt.Println("Invoice Created successfully succes")
			}
		}

		if paymentIntent.Status != "succeeded" {
			_, confirmErr := paymentintent.Confirm(
				paymentIntent.ID,
				&stripe.PaymentIntentConfirmParams{
					PaymentMethod: stripe.String(paymentIntent.PaymentMethod.ID),
				},
			)

			if confirmErr != nil {
				fmt.Fprintf(os.Stderr, "Error confirming PaymentIntent: %v\n", confirmErr)
				c.AbortWithStatus(http.StatusBadRequest)
				return
			}

			fmt.Println("PaymentIntent confirmed successfully")
		} else {

			fmt.Println("PaymentIntent is already succeeded")
		}

	case "payment_method.attached":
		var paymentMethod stripe.PaymentMethod
		err := json.Unmarshal(event.Data.Raw, &paymentMethod)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error parsing webhook JSON: %v\n", err)
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}
	case "payment_intent.failed":

		var paymentIntent stripe.PaymentIntent
		err := json.Unmarshal(event.Data.Raw, &paymentIntent)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error parsing webhook JSON: %v\n", err)
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}
		userID, _ := paymentIntent.Metadata["user_id"]
		addressID, _ := paymentIntent.Metadata["address_id"]

		userid, _ := strconv.Atoi(userID)
		addressid, _ := strconv.Atoi(addressID)

		invoice, err := cr.OrderUseCase.CreateInvoiceForFailedPayment(userid, addressid)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		Resul = invoice
		fmt.Println("Invoice Created successfully fail")

	default:
		fmt.Fprintf(os.Stderr, "Unhandled event type: %s\n", event.Type)

	}

	c.JSON(http.StatusOK, gin.H{"status": "suscces", "invoice": Resul})

}

// OrderStatus godoc
// @Summary Get the status of an order
// @Description Retrieves the status of an order based on the provided order ID.
// @ID get-order-status
// @Tags User Orders
// @Produce json
// @Param orderid path int true "Order ID for which the status should be retrieved"
// @Success 200 {string} string "Order status retrieved successfully"
// @Failure 400 {string} string "Bad request"
// @Router /user/orderstatus/{orderid} [get]
func (or *OrderHandler) OrderStatus(c *gin.Context) {
	strorderid := c.Param("orderid")
	orderid, err := strconv.Atoi(strorderid)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errror": "str failed"})
		return
	}
	order, err1 := or.OrderUseCase.ExecuteOrderid(orderid)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err1.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"order": order})

}

// PrintInvoice godoc
// @Summary Print invoice for an order
// @Description Generates and downloads the invoice for a specific order.
// @ID print-invoice
// @Tags User Orders
// @Produce json
// @Param orderid query int true "Order ID for which the invoice should be generated"
// @Success 200 {string} string "Invoice generated and downloaded successfully"
// @Failure 400 {string} string "Bad request"
// @Router /user/order/invoice [get]
func (or *OrderHandler) PrintInvoice(c *gin.Context) {
	userID, _ := c.Get("userId")
	userid := userID.(int)
	strorderId := c.Query("orderid")
	orderid, err := strconv.Atoi(strorderId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	pdf, err := or.OrderUseCase.ExecutPrintInvoice(orderid, userid)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.Header("Content-Disposition", "attachment;filename=Invoice.pdf")
	c.Header("Content_Type", "application/pdf")

	err = pdf.Output(c.Writer)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errror": err.Error()})
		return
	}

	pdfFilePath := "salesreport/invoice.pdf"

	err = pdf.OutputFileAndClose(pdfFilePath)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.File(pdfFilePath)

	c.JSON(http.StatusOK, gin.H{"pdf": pdfFilePath})

}
