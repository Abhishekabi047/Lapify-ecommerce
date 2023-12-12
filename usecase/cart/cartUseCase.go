package usecase

import (
	"errors"
	"log"
	"project/domain/entity"
	repository "project/repository/cart"
	productrepository "project/repository/product"
)

type CartUseCase struct {
	cartRepo    *repository.CartRepository
	productRepo *productrepository.ProductRepository
}

func NewCart(cartRepo *repository.CartRepository, productRepo *productrepository.ProductRepository) *CartUseCase {
	return &CartUseCase{cartRepo: cartRepo, productRepo: productRepo}
}

func (cu *CartUseCase) ExecuteAddToCart( id int, quantity int, userid int) error {
	var usercart *entity.Cart
	var cartid int
	usercart, err := cu.cartRepo.GetByUserid(userid)
	if err != nil {
		cart, err1 := cu.cartRepo.Create(userid)
		if err1 != nil {
			return errors.New("failed to create userid")
		}
		usercart = cart
		cartid = int(cart.ID)
	} else {
		cartid = int(usercart.ID)
	}

	prod, err := cu.productRepo.GetProductById(id)
	if err != nil {
		return errors.New("product not found")
	}
	cartitem := &entity.CartItem{
		CartId:      cartid,
		ProductId:   int(prod.ID),
		Category:    prod.Category,
		Quantity:    quantity,
		ProductName: prod.Name,
		Price:       int(prod.OfferPrize),
	}
	existingProduct, _ := cu.cartRepo.GetByName(prod.Name, cartid)

	if cartitem.Price == 0 {
		cartitem.Price = int(prod.Price)
	}
	if existingProduct == nil {
		err := cu.cartRepo.CreateCartItem(cartitem)
		if err != nil {
			return errors.New("error creating new  cart item")
		}
	} else {
		existingProduct.Quantity += quantity
		err := cu.cartRepo.UpdateCartItem(existingProduct)
		if err != nil {
			return errors.New("error updating new cart item")
		}
	}
	usercart.TotalPrize += cartitem.Price * int(quantity)
	usercart.ProductQuantity += quantity
	err1 := cu.cartRepo.UpdateCart(usercart)
	if err1 != nil {
		return errors.New("cart updation failed")
	}
	return nil
}

func (cu *CartUseCase) ExecuteCartItems(userId int) ([]entity.CartItem, error) {
	usercart, err := cu.cartRepo.GetByUserid(userId)
	if err != nil {
		return nil, errors.New("Failed to find usercart")
	}
	cartItems, err := cu.cartRepo.GetAllCartItems(int(usercart.ID))
	if err != nil {
		return nil, errors.New("errror finding cartitems")
	}
	return cartItems, nil
}

func (cu *CartUseCase) ExecuteRemoveCartItem(userid, id int) error {
	usercart, err := cu.cartRepo.GetByUserid(userid)
	if err != nil {
		return errors.New("error finding user cart")
	}
	prod, err := cu.productRepo.GetProductById(id)
	if err != nil {
		return errors.New("product not found")
	}
	// existingProd, err := cu.cartRepo.GetByName(product, int(usercart.ID))
	// if err != nil {
	// 	return errors.New("removing product failed")
	// }
	existingProd, err := cu.cartRepo.GetByName(prod.Name, int(usercart.ID))
	if err != nil {
		return errors.New("removing product failed")
	}


	if existingProd.Quantity == 1 {
		err := cu.cartRepo.RemoveCartItem(existingProd)
		if err != nil {
			log.Printf("Error removing product from cart: %v", err)
			return errors.New("reomving products failed")
		}
	} else {
		existingProd.Quantity -= 1
		err := cu.cartRepo.UpdateCartItem(existingProd)
		if err != nil {
			return errors.New("error upadting user")
		}
	}
	if prod.OfferPrize != 0{
		usercart.TotalPrize -= int(prod.OfferPrize)
	}else{
	usercart.TotalPrize -= int(prod.Price)
	}
	usercart.ProductQuantity -= 1
	if usercart.OfferPrize > 0 {
		usercart.OfferPrize = 0
	}
	err1 := cu.cartRepo.UpdateCart(usercart)
	if err1 != nil {
		return errors.New("Remove from cart failed")
	}
	return nil
}

func (cu *CartUseCase) ExecuteAddWishlist(productid int, userid int) error {
	product, err := cu.productRepo.GetProductById(productid)
	if err != nil {
		return errors.New("product not found")
	}
	exisiting, err := cu.cartRepo.GetProductsFromWishlist( product.ID, userid)
	if err != nil {
		return errors.New("error finding exisiting product")
	}
	if exisiting == true {
		return errors.New("product already exist")
	} else {
		wishprod := &entity.WishList{
			UserId:      userid,
			Category:    product.Category,
			ProductId:   product.ID,
			ProductName: product.Name,
			Prize:       int(product.Price),
		}
		err := cu.cartRepo.AddProductToWishlist(wishprod)
		if err != nil {
			return errors.New("error adding to wishlist")
		}
	}
	return nil
}

func (cu *CartUseCase) ExecuteRemoveFromWishList( productid, userid int) error {
	exisiting, err := cu.cartRepo.GetProductsFromWishlist( productid, userid)
	if err != nil {
		return errors.New("error getting products")
	}
	if exisiting == true {
		err := cu.cartRepo.RemoveFromWishlist(productid, userid)
		if err != nil {
			return errors.New("error removing products from wishlist")
		}
	} else {
		return errors.New("product not found")
	}
	return nil
}

func (cu *CartUseCase) ExecuteViewWishlist(userid int) ([]entity.WishList, error) {
	wishlist, err := cu.cartRepo.GetWishlist(userid)
	if err != nil {
		return nil, errors.New("Error getting wishlist")
	}
	return *wishlist, nil
}

func (c *CartUseCase) ExecuteApplyCoupon(userId int, code string) (int, error) {
	var totaloffer, totalprize int
	usercart, err := c.cartRepo.GetByUserid(userId)
	if err != nil {
		return 0, errors.New("failed to find user cart")
	}
	coupon, err := c.productRepo.GetCouponByCode(code)
	if err != nil {
		return 0, errors.New("coupon not found")
	}
	if coupon.UsedCount >= coupon.UsageLimit{
		return 0,errors.New("coupon usage exeeded")
	}else{
	// cartitems, err := c.cartRepo.GetAllCartItems(int(usercart.ID))
	// if err != nil {
	// 	return 0, errors.New("user cart item not found")
	// }
	// for _, cartitem := range cartitems {
	// 	if cartitem.Category == coupon.Category {
	// 		totalprize += int(cartitem.Price) * cartitem.Quantity
	// 	}
	// }
	totalprize = usercart.TotalPrize

	if totalprize > 0 {
		if coupon.Type == "percentage" {
			totaloffer = totalprize / coupon.Amount
		} else {
			totaloffer = coupon.Amount
		}
	} else {
		return 0, errors.New("Add more products")
	}
	if usercart.OfferPrize != 0 {
		return 0, errors.New("usercart offer already applied")
	} else {
		usercart.OfferPrize = totaloffer
		err := c.cartRepo.UpdateCart(usercart)
		if err != nil {
			return 0, errors.New("user cart update failed")
		}
		var UsedCoupon = entity.UsedCoupon{
			UserId:     userId,
			CouponCode: code,
		}
		err1 := c.productRepo.UpdateCouponUsage(&UsedCoupon)
		if err1 != nil {
			return 0, errors.New("user coupon usage updation failed")
		}
		coupon.UsedCount=coupon.UsedCount + 1
		
		err = c.productRepo.UpdateCouponCount(coupon)
		if err != nil {
			return 0, errors.New("user update coupon count failed")
		}
	}

	}
	return totaloffer, nil

}

func (c *CartUseCase) ExecuteOfferCheck(userid int) (*[]entity.Offer, error) {
	usercart, err := c.cartRepo.GetByUserid(userid)
	if err != nil {
		return nil, errors.New("failed to find user cart")
	}
	offer, err := c.productRepo.GetOfferByPrize(int(usercart.TotalPrize))
	if err != nil {
		return nil, errors.New("add few more products")
	}
	return offer, nil
}

func (cu *CartUseCase) ExecuteCart(userid int) (*entity.Cart, error) {
	userCart, err := cu.cartRepo.GetByUserid(userid)
	if err != nil {
		return nil, errors.New("failed to find user")
	} else {
		return userCart, nil
	}
}

func (cu *CartUseCase) ExecuteCartitem(userid int) (*[]entity.CartItem, error) {
	userCart, err := cu.cartRepo.GetByUserid(userid)
	if err != nil {
		return nil, errors.New("Failed to find user")
	}
	cartitems, err := cu.cartRepo.GetAllCartItems(int(userCart.ID))
	if err != nil {
		return nil, err
	}
	return &cartitems, nil
}
