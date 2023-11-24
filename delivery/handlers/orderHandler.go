package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
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
}

func NewOrderHandler(OrderUseCase *usecase.OrderUseCase) *OrderHandler {
	return &OrderHandler{OrderUseCase}
}

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
		razorId, orderId, err1 := oh.OrderUseCase.ExecuteRazorPay(userid, addressId, c)
		if err1 != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err1.Error()})
		}
		c.HTML(http.StatusOK, "razor.html", gin.H{
			"message": "make payment",
			"razorId": razorId,
			"orderid": orderId})
		c.JSON(http.StatusOK, gin.H{"message": "complete your razor pay through ", "razorId": razorId, "orderid": orderId, "userid": userid})
	}
}

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
	c.JSON(http.StatusOK, gin.H{"message ": "order cancelled"})
}

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

func (or *OrderHandler) SalesReportByPeriod(c *gin.Context) {
	period := c.Param("period")

	report, err := or.OrderUseCase.ExecuteSalesReportByPeriod(period)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
	c.JSON(http.StatusOK, gin.H{"report": report})
}

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

	fmt.Println("hello")
	// var order entity.Order

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
				// Access resul.ID and use it as needed
				err := cr.OrderUseCase.UpdateInvoiceStatus(int(Resul.OrderId), "succesfull")
				if err != nil {
					fmt.Fprintf(os.Stderr, "Error updating invoice status: %v\n", err)
					// Handle the error as needed
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

			// Log confirmation success
			fmt.Println("PaymentIntent confirmed successfully")
		} else {
			// Log that PaymentIntent is already succeeded
			fmt.Println("PaymentIntent is already succeeded")
		}

		// Then define and call a func to handle the successful payment intent.
		// handlePaymentIntentSucceeded(paymentIntent)
	case "payment_method.attached":
		var paymentMethod stripe.PaymentMethod
		err := json.Unmarshal(event.Data.Raw, &paymentMethod)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error parsing webhook JSON: %v\n", err)
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}
	case "payment_intent.failed":
		// Handle failed payment intents here
		// You can create an invoice or perform other actions
		// based on the failed payment intent.
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
		// Example: Create an invoice for the failed payment intent
		invoice, err := cr.OrderUseCase.CreateInvoiceForFailedPayment(userid, addressid)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		Resul = invoice
		fmt.Println("Invoice Created successfully fail")
		// Then define and call a func to handle the successful attachment of a PaymentMethod.
		// handlePaymentMethodAttached(paymentMethod)
	// ... handle other event types
	default:
		fmt.Fprintf(os.Stderr, "Unhandled event type: %s\n", event.Type)

	}

	c.JSON(http.StatusOK, gin.H{"status": "suscces", "invoice": Resul})

}
func (or *OrderHandler) OrderStatus(c *gin.Context){
	strorderid:=c.PostForm("orderid")
	orderid,err:=strconv.Atoi(strorderid)
	if err != nil{
		c.JSON(http.StatusBadRequest,gin.H{"errror":"str failed"})
		return
	}
	order,err1:=or.OrderUseCase.ExecuteOrderid(orderid)
	if err != nil{
		c.JSON(http.StatusBadRequest,gin.H{"error":err1.Error()})
		return
	}
	c.JSON(http.StatusOK,gin.H{"order":order})

}
