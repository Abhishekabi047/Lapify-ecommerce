package order

import (
	"errors"
	"fmt"
	"project/domain/entity"
	cartrepository "project/repository/cart"
	repository "project/repository/order"
	productrepository "project/repository/product"
	userrepository "project/repository/user"
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
	fmt.Println("Cart Items:")
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
			return nil, errors.New("quantity decrease failed")
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
