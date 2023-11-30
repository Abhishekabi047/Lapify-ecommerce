package handlers

import (
	"fmt"
	"net/http"
	"project/delivery/middleware"
	"project/delivery/models"
	"project/domain/entity"
	Cartusecase "project/usecase/cart"
	Productusecase "project/usecase/product"
	usecase "project/usecase/user"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
)

type UserHandler struct {
	UserUseCase    *usecase.UserUseCase
	ProductUseCase *Productusecase.ProductUseCase
	CartUSeCase    *Cartusecase.CartUseCase
}

func NewUserhandler(UserUseCase *usecase.UserUseCase, ProductUseCase *Productusecase.ProductUseCase, CartUseCase *Cartusecase.CartUseCase) *UserHandler {
	return &UserHandler{UserUseCase, ProductUseCase, CartUseCase}
}

func (uh *UserHandler) SignupWithOtp(c *gin.Context) {
	var user models.Signup
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	key, err := uh.UserUseCase.ExecuteSignupWithOtp(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	} else {
		c.JSON(http.StatusOK, gin.H{"otp send to ": user.Phone, "key": key})
	}
}

func (uh *UserHandler) SignupOtpValidation(c *gin.Context) {
	key := c.PostForm("key")
	otp := c.PostForm("otp")
	err := uh.UserUseCase.ExecuteSignupOtpValidation(key, otp)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	} else {
		c.JSON(http.StatusOK, gin.H{"message": "user signup succesfull"})
	}
}

func (uh *UserHandler) LoginWithPassword(c *gin.Context) {

	phone := c.PostForm("phone")
	password := c.PostForm("password")

	userId, err := uh.UserUseCase.ExecuteLoginWithPassword(phone, password)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	} else {
		fmt.Println("userId:", userId)
		middleware.CreateToken(userId, phone, "user", c)
		c.JSON(http.StatusOK, gin.H{"message": "user logged in succesfully and cookie stored"})
	}

}

func (po *UserHandler) Products(c *gin.Context) {

	userID, exists := c.Get("userId")
	if !exists || userID == nil {
		// Handle the case where "userId" is not set in the context or is nil
		c.JSON(http.StatusBadRequest, gin.H{"error": "userId not found in the context"})
		return
	}

	pagestr := c.DefaultQuery("page", "1")
	limitstr := c.DefaultQuery("limit", "10")
	page, err := strconv.Atoi(pagestr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed page conv"})
		return
	}
	limit, err := strconv.Atoi(limitstr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed limit conv"})
		return
	}
	productlist, err := po.ProductUseCase.ExecuteProductList(page, limit)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "product not found"})
		return
	}
	prodlist := make([]entity.Product, len(productlist))
	for i, product := range productlist {
		prodlist[i] = entity.Product{
			ID:       product.ID,
			Name:     product.Name,
			Price:    product.Price,
			ImageURL: product.ImageURL,
			Size:     product.Size,
			Category: product.Category,
		}
	}
	c.JSON(http.StatusOK, gin.H{"userId": userID, "products": prodlist})
}

func (pd *UserHandler) ProductDetails(c *gin.Context) {
	ids := c.Param("productid")
	id, err := strconv.Atoi(ids)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id string conv failed"})
		return
	}
	product, productdetails, err1 := pd.ProductUseCase.ExecuteProductDetails(id)
	if err1 != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "product not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"products": product, "product details": productdetails})
}

func (ac *UserHandler) AddToCart(c *gin.Context) {
	userID, exists := c.Get("userId")
	if !exists || userID == nil {
		// Handle the case where "userid" is not set in the context or is nil
		c.JSON(http.StatusBadRequest, gin.H{"error": "userid not found in the context"})
		return
	}
	userid := userID.(int)
	fmt.Println("userId", userid)
	product := c.PostForm("category")
	strid := c.PostForm("productid")
	strquantity := c.PostForm("quantity")
	id, err := strconv.Atoi(strid)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "str convertion failed"})
		return
	}
	quantity, err := strconv.Atoi(strquantity)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "str convertion failed"})
		return
	}

	err = ac.CartUSeCase.ExecuteAddToCart(product, id, quantity, userid)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	addedProduct, err := ac.CartUSeCase.ExecuteCartItems(userid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "product addedto car succesfully", "addedproduct": addedProduct})
}

func (rc *UserHandler) RemoveFromCart(c *gin.Context) {
	userID, _ := c.Get("userId")
	userid := userID.(int)
	id := c.PostForm("id")
	product := c.PostForm("product")
	Id, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "str convertion failed"})
		return
	}
	err1 := rc.CartUSeCase.ExecuteRemoveCartItem(userid, Id, product)
	if err1 != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err1.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "product removed from cart"})
}

func (cu *UserHandler) Cart(c *gin.Context) {
	userID, _ := c.Get("userId")
	userid := userID.(int)
	var usercartresponse entity.Cart
	usercart, err1 := cu.CartUSeCase.ExecuteCart(userid)
	if err1 != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err1.Error()})
		return
	}
	copier.Copy(&usercartresponse, &usercart)
	c.JSON(http.StatusOK, gin.H{"usercart": usercartresponse})
}

func (cu *UserHandler) AddToWishList(c *gin.Context) {
	userID, _ := c.Get("userId")
	userid := userID.(int)
	strproduct := c.Param("category")
	product, err := strconv.Atoi(strproduct)
	strId := c.Param("productid")
	id, err := strconv.Atoi(strId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "string conc failed"})
		return
	}
	err1 := cu.CartUSeCase.ExecuteAddWishlist(product, id, userid)
	if err1 != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err1.Error()})
		return
	}
	wishlistItems, err := cu.CartUSeCase.ExecuteViewWishlist(userid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to retrieve wishlist items"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "product added to wishlist", "wishlist": wishlistItems})
}

func (cu *UserHandler) RemoveFromWishlist(c *gin.Context) {
	userID, _ := c.Get("userId")
	userid := userID.(int)
	strproduct := c.Param("product")
	product, err := strconv.Atoi(strproduct)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "strng conversion failed"})
		return
	}
	ID := c.Param("id")
	id, err := strconv.Atoi(ID)

	err1 := cu.CartUSeCase.ExecuteRemoveFromWishList(product, id, userid)
	if err1 != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err1.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "succesfully removed from wishlist"})
}
func (cu *UserHandler) ViewWishlist(c *gin.Context) {
	userID, _ := c.Get("userId")
	userid := userID.(int)
	wishlist, err := cu.CartUSeCase.ExecuteViewWishlist(userid)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
	c.JSON(http.StatusOK, gin.H{"message": wishlist})
}

func (cu *UserHandler) Logout(c *gin.Context) {
	err := middleware.DeleteToken(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errror": "cookie delete failed"})
		return
	} else {
		c.JSON(http.StatusOK, gin.H{"message": "user logged out succesfully"})
	}
}

func (cu *UserHandler) AddAddress(c *gin.Context) {
	userID, _ := c.Get("userId")
	userid := userID.(int)
	var address entity.UserAddress
	if err := c.ShouldBindJSON(&address); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	address.User_id = userid
	err := cu.UserUseCase.ExecuteAddAddress(&address)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	} else {
		c.JSON(http.StatusOK, gin.H{"message": "address added succesfully"})
	}
}

func (cy *UserHandler) ShowUserDetails(c *gin.Context) {
	userID, _ := c.Get("userId")
	userid := userID.(int)
	userdetails, address, err := cy.UserUseCase.ExecuteShowUserDetails(userid)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"userdetails": userdetails, "address": address})
}

func (eu *UserHandler) EditProfile(c *gin.Context) {
	userID, _ := c.Get("userId")
	userid := userID.(int)
	var userinput models.EditUser
	if err := c.ShouldBindJSON(&userinput); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var user entity.User
	copier.Copy(&user, &userinput)
	err := eu.UserUseCase.ExecuteEditProfile(user, userid)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	updatedUser, _, err := eu.UserUseCase.ExecuteShowUserDetails(userid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch updated user details"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "user edited succesfully", "updateduser": updatedUser})
}

func (cp *UserHandler) ChangePassword(c *gin.Context) {
	userID, _ := c.Get("userId")
	userid := userID.(int)
	key, err := cp.UserUseCase.ExecuteChangePassword(userid)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "otp send succesfully", "key": key})
}

func (cp *UserHandler) OtpValidationPassword(c *gin.Context) {
	userID, _ := c.Get("userId")
	userid := userID.(int)
	password := c.PostForm("password")
	otp := c.PostForm("otp")
	err := cp.UserUseCase.ExecuteOtpValidationPassword(password, otp, userid)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "password changed succesfully"})
}

func (cl *UserHandler) CartItems(c *gin.Context) {

	userID, _ := c.Get("userId")
	userid := userID.(int)
	cartitems, err := cl.CartUSeCase.ExecuteCartitem(userid)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"cartlist": cartitems})
}

func (ea *UserHandler) EditAddress(c *gin.Context) {
	var useraddress entity.UserAddress
	if err := c.ShouldBindJSON(&useraddress); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	userID, _ := c.Get("userId")
	userid := userID.(int)
	addresstype := c.Param("type")
	err := ea.UserUseCase.ExecuteEditAddress(useraddress, userid, addresstype)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": "address edited succesfully"})
}

func (da *UserHandler) DeleteAddress(c *gin.Context) {
	userID, _ := c.Get("userId")
	userid := userID.(int)
	addresstype := c.Param("type")
	err := da.UserUseCase.ExecuteDeleteAddress(userid, addresstype)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": "address Deleted succesfully"})
}

func (or *UserHandler) SearchProduct(c *gin.Context) {
	pagestr := c.DefaultQuery("page", "1")
	limistr := c.DefaultQuery("limit", "5")
	search := c.Query("search")
	page, err := strconv.Atoi(pagestr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "str convertion failed"})
		return
	}
	limit, err := strconv.Atoi(limistr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "str convertion failed"})
		return
	}
	productlist, err := or.ProductUseCase.ExecuteProductSearch(page, limit, search)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	responselist := make([]entity.Product, len(productlist))
	for i, product := range productlist {
		responselist[i] = entity.Product{
			ID:       product.ID,
			Name:     product.Name,
			Price:    product.Price,
			Category: product.Category,
			ImageURL: product.ImageURL,
			Size:     product.Size,
		}
	}
	c.JSON(http.StatusOK, gin.H{"products": responselist})
}

func (sc *UserHandler) SortByCategory(c *gin.Context) {
	pagestr := c.DefaultQuery("page", "1")
	limitstr := c.DefaultQuery("limit", "5")
	ID := c.Query("id")
	id, err := strconv.Atoi(ID)
	page, err := strconv.Atoi(pagestr)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": "str conv failed"})
		return
	}
	limit, err := strconv.Atoi(limitstr)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": "str conv failed"})
		return
	}
	productlist, err1 := sc.ProductUseCase.ExecuteProductByCategory(page, limit, id)
	if err1 != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err1.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"Products": productlist})
}

func (sc *UserHandler) SortByFilter(c *gin.Context) {
	strminPrize:=c.Query("minprize")
	strmaxPrize:=c.Query("maxprize")
	strcategory:=c.Query("category")
	size:=c.Query("size")
	minPrize,_:=strconv.Atoi(strminPrize)
	maxPrize,_:=strconv.Atoi(strmaxPrize)
	category,_:=strconv.Atoi(strcategory)
	productlist,err1:=sc.ProductUseCase.ExecuteProductFilter(size,minPrize,maxPrize,category)
	if err1 != nil{
		c.JSON(http.StatusBadRequest,gin.H{"error":err1.Error()})
		return
	}
	c.JSON(http.StatusOK,gin.H{"products":productlist})
}
func (sc *UserHandler) ApplyCoupon(c *gin.Context){
	userID,_:=c.Get("userId")
	userid:=userID.(int)
	code:=c.PostForm("code")

	totaloffer,err:=sc.CartUSeCase.ExecuteApplyCoupon(userid,code)
	if err != nil{
		c.JSON(http.StatusBadRequest,gin.H{"error":err.Error()})
		return
	}
	c.JSON(http.StatusOK,gin.H{"offer prize":totaloffer,"offer":"applied succesfully"})
}
