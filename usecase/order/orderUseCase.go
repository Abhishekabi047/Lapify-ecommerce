package order

import (
	"errors"
	"fmt"
	"project/domain/entity"
	"project/domain/utils"
	cartrepository "project/repository/cart"
	repository "project/repository/order"
	productrepository "project/repository/product"
	userrepository "project/repository/user"
	"time"

	"github.com/gin-gonic/gin"
	razorpay "github.com/razorpay/razorpay-go"
)

type OrderUseCase struct {
	orderRepo   *repository.OrderRepository
	cartRepo    *cartrepository.CartRepository
	userRepo    *userrepository.UserRepository
	productRepo *productrepository.ProductRepository
}

func NewOrder(orderRepo *repository.OrderRepository, cartRepo *cartrepository.CartRepository, userRepo *userrepository.UserRepository, productRepo *productrepository.ProductRepository) *OrderUseCase {
	return &OrderUseCase{orderRepo: orderRepo, cartRepo: cartRepo, userRepo: userRepo, productRepo: productRepo}
}

func (or *OrderUseCase) ExecuteOrderCod(userid int, address int) (*entity.Invoice, error) {
	var orderitems []entity.OrderItem
	cart, err := or.cartRepo.GetByUserid(userid)
	if err != nil {
		return nil, errors.New("cart not found")
	}

	cartitems, err := or.cartRepo.GetAllCartItems(int(cart.ID))
	if err != nil {
		return nil, errors.New("cart items not found")
	}

	for _, cartitem := range cartitems {
		fmt.Printf("ProductID: %d, Category: %d, Quantity: %d, Price: %d\n",
			cartitem.ProductId, cartitem.Category, cartitem.Quantity, cartitem.Price)
	}

	useraddress, err := or.userRepo.GetAddressByID(userid)
	if err != nil {
		return nil, errors.New("address not found")
	}
	Total := cart.TotalPrize - int(cart.OfferPrize)
	order := &entity.Order{
		UserId:        cart.UserId,
		Addressid:     useraddress.Id,
		Total:         Total,
		Status:        "pending",
		PaymentMethod: "cod",
		PaymentStatus: "pending",
	}
	orderID, err := or.orderRepo.Create(order)
	if err != nil {
		return nil, errors.New("order placing failed")
	}
	InvoiceData := &entity.Invoice{
		OrderId:     orderID,
		UserId:      userid,
		AddressType: useraddress.Type,
		Quantity:    cart.ProductQuantity,
		Price:       float64(order.Total),
		Payment:     order.PaymentMethod,
		Status:      order.PaymentStatus,
		PaymentId:   "nil",
		Remark:      "nil",
	}
	invoice, err := or.orderRepo.CreateInvoice(InvoiceData)
	if err != nil {
		return nil, errors.New("error creating invoice")
	}
	for _, cartitem := range cartitems {
		orderitem := entity.OrderItem{
			OrderId:   orderID,
			ProductId: cartitem.ProductId,
			Category:  cartitem.Category,
			Quantity:  cartitem.Quantity,
			Prize:     cartitem.Price,
		}
		orderitems = append(orderitems, orderitem)
		inventory := entity.Inventory{
			ProductId:       cartitem.ProductId,
			ProductCategory: cartitem.Category,
			Quantity:        cartitem.Quantity,
		}
		// fmt.Println("prod:", cartitem.ProductId)
		err := or.productRepo.DecreaseProductQuantity(&inventory)
		if err != nil {
			return nil, err
		}
	}
	err1 := or.orderRepo.CreateOrderItems(orderitems)
	if err1 != nil {
		return nil, errors.New("failed to create order items")
	}
	err = or.cartRepo.RemoveCartItems(int(cart.ID))
	if err != nil {
		return nil, errors.New("removing cart failed")
	}
	cart.ProductQuantity = 0
	cart.TotalPrize = 0
	cart.OfferPrize = 0
	cart.ProductQuantity = 0
	err = or.cartRepo.UpdateCart(cart)
	if err != nil {
		return nil, errors.New("error upadting cart")
	}
	return invoice, nil
}

func (co *OrderUseCase) ExecuteCancelOrder(orderid int) error {
	result, err := co.orderRepo.GetOrderById(orderid)
	if err != nil {
		return errors.New("erroro getting orderid")
	}
	// userid,err:=co.userRepo.GetById(result.UserId)
	// if err != nil{
	// 	return errors.New("error getting user")
	// }
	if result.Status != "pending" && result.Status != "confirmed" {
		return errors.New("order cancel time exceeded")
	}
	result.Status = "cancelled"
	err1 := co.orderRepo.Update(result)
	if err1 != nil {
		return errors.New("order cancellation failed")
	}
	return nil
}

func (co *OrderUseCase) ExecuteOrderHistory(userid, page, limit int) ([]entity.Order, error) {
	offset := (page - 1) * limit
	orderList, err := co.orderRepo.GetAllOrders(userid, offset, limit)
	if err != nil {
		return nil, errors.New("failed to get order list")
	}
	return orderList, nil
}

func (co *OrderUseCase) ExecuteOrderUpdate(OrderId int, status string) error {
	result, err := co.orderRepo.GetOrderById(OrderId)
	if err != nil {
		return errors.New("error finding order")
	}
	result.Status = status
	err1 := co.orderRepo.Update(result)
	if err1 != nil {
		return errors.New("error updating  order status")
	}
	return nil
}

func (co *OrderUseCase) UpdatedUser(orderid int) (*entity.Order, error) {

	result, err := co.orderRepo.GetOrderById(orderid)
	if err != nil {
		return nil, errors.New("error finding order")
	}
	return result, nil
}

func (co *OrderUseCase) ExecuteAdminOrder(page, limit int) ([]entity.Order, error) {
	offset := (page - 1) * limit
	result, err := co.orderRepo.GetAllOrderList(offset, limit)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (co *OrderUseCase) ExecuteAdminCancelOrder(orderid int) error {
	result, err := co.orderRepo.GetOrderById(orderid)
	if err != nil {
		return err
	}
	if result.Status != "pending" {
		return errors.New("admin already confirmed order")
	}
	result.Status = "cancelled"
	err1 := co.orderRepo.Update(result)
	if err1 != nil {
		return errors.New("updation failed")
	}
	return nil
}

func (rp *OrderUseCase) ExecuteRazorPay(userId, address int, c *gin.Context) (string, int, error) {
	var orderItems []entity.OrderItem
	cart, err := rp.cartRepo.GetByUserid(userId)
	if err != nil {
		return "", 0, errors.New("Cart not Found")
	}
	cartitems, err1 := rp.cartRepo.GetAllCartItems(int(cart.ID))
	if err1 != nil {
		return "", 0, errors.New("cartitems ")
	}
	usraddress, err2 := rp.userRepo.GetAddressByID(userId)
	if err2 != nil {
		return "", 0, errors.New("address not found")
	}
	client := razorpay.NewClient("rzp_test_leWrFNIomWqk5W", "R59k58EhgS48BaauF22urj5A")

	data := map[string]interface{}{
		"amount":   int(cart.TotalPrize),
		"currency": "INR",
		"receipt":  "101",
	}

	body, err := client.Order.Create(data, nil)
	if err != nil {
		return "", 0, errors.New("Errro creating order")
	}

	razorId, _ := body["id"].(string)
	Total := cart.TotalPrize
	order := &entity.Order{
		UserId:        cart.UserId,
		Addressid:     usraddress.Id,
		Total:         Total,
		Status:        "pending",
		PaymentMethod: "razorpay",
		PaymentStatus: "pending",
		PaymentId:     razorId,
	}
	orderId, err := rp.orderRepo.Create(order)
	if err != nil {
		return "", 0, errors.New("Order placing failed")
	}
	for _, cartItem := range cartitems {
		orderitem := entity.OrderItem{
			OrderId:   orderId,
			ProductId: cartItem.ProductId,
			Category:  cartItem.Category,
			Quantity:  cartItem.Quantity,
			Prize:     cartItem.Price,
		}
		orderItems = append(orderItems, orderitem)
		inventory := entity.Inventory{
			ProductId:       cartItem.ProductId,
			ProductCategory: cartItem.Category,
			Quantity:        cartItem.Quantity,
		}

		err := rp.productRepo.DecreaseProductQuantity(&inventory)
		if err != nil {
			return "", 0, err
		}
	}
	err3 := rp.orderRepo.CreateOrderItems(orderItems)
	if err3 != nil {
		return "", 0, errors.New("Errro creating orderitems")
	}
	return razorId, orderId, nil
}

func (rv *OrderUseCase) ExecuteRazorPayVerification(signature, razorid, PaymentId string) (*entity.Invoice, error) {
	result, err := rv.orderRepo.GetByRazorId(razorid)
	if err != nil {
		return nil, errors.New("order not found")
	}
	err1 := utils.RazorPaymentVerification(signature, razorid, PaymentId)
	if err1 != nil {
		result.PaymentStatus = "failed"
		result.PaymentId = PaymentId
		err2 := rv.orderRepo.Update(result)
		if err2 != nil {
			return nil, errors.New("payment updation failed")
		}
	}
	result.PaymentStatus = "succesfull"
	result.PaymentId = PaymentId
	err3 := rv.orderRepo.Update(result)
	if err3 != nil {
		return nil, errors.New("payment updation failed")
	}
	userCart, err := rv.cartRepo.GetByUserid(result.UserId)
	if err != nil {
		return nil, errors.New("usercart not found")
	}
	useraddress, err := rv.userRepo.GetAddressById(result.Addressid)
	if err != nil {
		return nil, errors.New("useraddress not found")
	}
	Total := userCart.TotalPrize
	InvoiceData := &entity.Invoice{
		OrderId:     result.ID,
		UserId:      result.UserId,
		AddressType: useraddress.Type,
		Quantity:    userCart.ProductQuantity,
		Price:       float64(Total),
		Payment:     "razorpay",
		Status:      result.PaymentStatus,
		PaymentId:   PaymentId,
		Remark:      "",
	}

	Invoice, err := rv.orderRepo.CreateInvoice(InvoiceData)
	if err != nil {
		return nil, errors.New("Invoice creating failed")
	}
	err4 := rv.cartRepo.RemoveCartItems(int(userCart.ID))
	if err4 != nil {
		return nil, errors.New("removing cart items failed")
	}
	userCart.ProductQuantity = 0
	userCart.TotalPrize = 0
	userCart.OfferPrize = 0
	userCart.ProductQuantity = 0
	err = rv.cartRepo.UpdateCart(userCart)
	if err != nil {
		return nil, errors.New("error upadting cart")
	}
	return Invoice, nil

}
func (sr *OrderUseCase) ExecuteSalesReportByPeriod(period string) (*entity.SalesReport, error) {
	startdate, enddate := utils.CalcualtePeriodDate(period)

	orders, err := sr.orderRepo.GetByDate(startdate, enddate)
	if err != nil {
		return nil, errors.New("report fetching failed")
	}
	return orders, nil
}

func (sr *OrderUseCase) ExecuteSalesReportByDate(startdate, enddate time.Time) (*entity.SalesReport, error) {
	orders, err := sr.orderRepo.GetByDate(startdate, enddate)
	if err != nil {
		return nil, errors.New("report fetching failed")
	}
	return orders, nil
}

func (sr *OrderUseCase) ExecuteSalesReportByPaymentMethod(startdate, enddate time.Time, paymentmethod string) (*entity.SalesReport, error) {
	orders, err := sr.orderRepo.GetByPaymentMethod(startdate, enddate, paymentmethod)
	if err != nil {
		return nil, errors.New("report fetching failed")
	}
	return orders, nil
}
func (cu *OrderUseCase) ExecuteCartit(userid int) (*entity.Cart, error) {
	userCart, err := cu.cartRepo.GetByUserid(userid)
	if err != nil {
		return nil, errors.New("failed to find user")
	} else {
		return userCart, nil
	}
}

func (or *OrderUseCase) ExecuteInvoiceStripe(userid int, address int) (*entity.Invoice, error) {
	var orderitems []entity.OrderItem
	cart, err := or.cartRepo.GetByUserid(userid)
	if err != nil {
		return nil, errors.New("cart not found")
	}

	cartitems, err := or.cartRepo.GetAllCartItems(int(cart.ID))
	if err != nil {
		return nil, errors.New("cart items not found")
	}

	for _, cartitem := range cartitems {
		fmt.Printf("ProductID: %d, Category: %d, Quantity: %d, Price: %d\n",
			cartitem.ProductId, cartitem.Category, cartitem.Quantity, cartitem.Price)
	}

	useraddress, err := or.userRepo.GetAddressByID(userid)
	if err != nil {
		return nil, errors.New("address not found")
	}
	Total := cart.TotalPrize - int(cart.OfferPrize)
	order := &entity.Order{
		UserId:        cart.UserId,
		Addressid:     useraddress.Id,
		Total:         Total,
		Status:        "pending",
		PaymentMethod: "Stripe",
		PaymentStatus: "pending",
	}
	orderID, err := or.orderRepo.Create(order)
	if err != nil {
		return nil, errors.New("order placing failed")
	}
	InvoiceData := &entity.Invoice{
		OrderId:     orderID,
		UserId:      userid,
		AddressType: useraddress.Type,
		Quantity:    cart.ProductQuantity,
		Price:       float64(order.Total),
		Payment:     order.PaymentMethod,
		Status:      order.PaymentStatus,
		PaymentId:   "nil",
		Remark:      "nil",
	}
	invoice, err := or.orderRepo.CreateInvoice(InvoiceData)
	if err != nil {
		return nil, errors.New("error creating invoice")
	}

	for _, cartitem := range cartitems {
		orderitem := entity.OrderItem{
			OrderId:   orderID,
			ProductId: cartitem.ProductId,
			Category:  cartitem.Category,
			Quantity:  cartitem.Quantity,
			Prize:     cartitem.Price,
		}
		orderitems = append(orderitems, orderitem)
		inventory := entity.Inventory{
			ProductId:       cartitem.ProductId,
			ProductCategory: cartitem.Category,
			Quantity:        cartitem.Quantity,
		}
		err := or.productRepo.DecreaseProductQuantity(&inventory)
		if err != nil {
			return nil, err
		}
	}
	err1 := or.orderRepo.CreateOrderItems(orderitems)
	if err1 != nil {
		return nil, errors.New("failed to create order items")
	}

	err = or.cartRepo.RemoveCartItems(int(cart.ID))
	if err != nil {
		return nil, errors.New("removing cart failed")
	}
	cart.ProductQuantity = 0
	cart.TotalPrize = 0
	cart.OfferPrize = 0
	cart.ProductQuantity = 0
	err = or.cartRepo.UpdateCart(cart)
	if err != nil {
		return nil, errors.New("error upadting cart")

	}
	return invoice, nil
}

func (or *OrderUseCase) CreateInvoiceForFailedPayment(userid int, address int) (*entity.Invoice, error) {
	var orderitems []entity.OrderItem
	cart, err := or.cartRepo.GetByUserid(userid)
	if err != nil {
		return nil, errors.New("cart not found")
	}

	cartitems, err := or.cartRepo.GetAllCartItems(int(cart.ID))
	if err != nil {
		return nil, errors.New("cart items not found")
	}

	for _, cartitem := range cartitems {
		fmt.Printf("ProductID: %d, Category: %d, Quantity: %d, Price: %d\n",
			cartitem.ProductId, cartitem.Category, cartitem.Quantity, cartitem.Price)
	}

	useraddress, err := or.userRepo.GetAddressById(userid)
	if err != nil {
		return nil, errors.New("address not found")
	}
	Total := cart.TotalPrize - int(cart.OfferPrize)
	order := &entity.Order{
		UserId:        cart.UserId,
		Addressid:     useraddress.Id,
		Total:         Total,
		Status:        "pending",
		PaymentMethod: "Stripe",
		PaymentStatus: "Failed",
	}
	orderID, err := or.orderRepo.Create(order)
	if err != nil {
		return nil, errors.New("order placing failed")
	}
	InvoiceData := &entity.Invoice{
		OrderId:     orderID,
		UserId:      userid,
		AddressType: useraddress.Type,
		Quantity:    cart.ProductQuantity,
		Price:       float64(order.Total),
		Payment:     order.PaymentMethod,
		Status:      order.PaymentStatus,
		PaymentId:   "nil",
		Remark:      "nil",
	}
	invoice, err := or.orderRepo.CreateInvoice(InvoiceData)
	if err != nil {
		return nil, errors.New("error creating invoice")
	}

	for _, cartitem := range cartitems {
		orderitem := entity.OrderItem{
			OrderId:   orderID,
			ProductId: cartitem.ProductId,
			Category:  cartitem.Category,
			Quantity:  cartitem.Quantity,
			Prize:     cartitem.Price,
		}
		orderitems = append(orderitems, orderitem)
		inventory := entity.Inventory{
			ProductId:       cartitem.ProductId,
			ProductCategory: cartitem.Category,
			Quantity:        cartitem.Quantity,
		}

		err := or.productRepo.DecreaseProductQuantity(&inventory)
		if err != nil {
			return nil, err
		}
	}
	err1 := or.orderRepo.CreateOrderItems(orderitems)
	if err1 != nil {
		return nil, errors.New("failed to create order items")
	}
	if order.PaymentStatus == "succesful" {
		err = or.cartRepo.RemoveCartItems(int(cart.ID))
		if err != nil {
			return nil, errors.New("removing cart failed")
		}
		cart.ProductQuantity = 0
		cart.TotalPrize = 0
		cart.OfferPrize = 0
		cart.ProductQuantity = 0
		err = or.cartRepo.UpdateCart(cart)
		if err != nil {
			return nil, errors.New("error upadting cart")
		}
	}
	return invoice, nil
}
func (uc *OrderUseCase) UpdateInvoiceStatus(orderID int, status string) error {

	invoice, err := uc.orderRepo.GetOrderById(orderID)
	if err != nil {
		return err
	}
	invoice.PaymentStatus = status
	err = uc.orderRepo.Update(invoice)
	if err != nil {
		return err
	}

	return nil
}

func (co *OrderUseCase) ExecuteOrderid(OrderId int) (*entity.Order, error) {
	result, err := co.orderRepo.GetOrderById(OrderId)
	if err != nil {
		return nil, err
	}
	return result, nil

}
