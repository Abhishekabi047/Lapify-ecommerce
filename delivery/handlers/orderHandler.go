package handlers

import (
	"net/http"
	usecase "project/usecase/order"
	"strconv"

	"github.com/gin-gonic/gin"
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
	addressId, err := strconv.Atoi(straddress)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "strng conversion failed"})
		return
	}
	invoice, err := oh.OrderUseCase.ExecuteOrderCod(userid, addressId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	} else {
		c.JSON(http.StatusOK, gin.H{"invoice": invoice})
	}

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
	status := c.Param("status")
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
