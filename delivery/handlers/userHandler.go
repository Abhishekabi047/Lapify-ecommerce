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

// SignupWithOtp godoc
// @Summary Sign up a user with OTP
// @Description Registers a new user using OTP verification.
// @ID signup-with-otp
// @Accept json
// @Tags User
// @Produce json
// @Param user body models.Signup true "User details for signup with OTP"
// @Success 200 {string} string "OTP sent successfully to the provided phone number"
// @Failure 400 {string} string "Bad request: Invalid request"
// @Failure 500 {string} string "Internal Server Error: Something went wrong"
// @Router /user/signup [post]
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

// SignupOtpValidation godoc
// @Summary Validate OTP for user signup
// @Description Validates the provided OTP for user signup.
// @ID signup-otp-validation
// @Accept multipart/form-data
// @Tags User
// @Produce json
// @Param key formData string true "Key associated with the OTP validation"
// @Param otp formData string true "OTP to be validated"
// @Success 200 {string} string "User signup successful"
// @Failure 401 {string} string "Unauthorized: Invalid key or OTP"
// @Router /user/signup/otpvalidation [post]
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

// LoginWithPassword godoc
// @Summary Log in a user using phone and password
// @Description Authenticate a user using phone and password, and generate an authentication token
// @ID loginWithPassword
// @Tags User
// @Accept multipart/form-data
// @Produce json
// @Param phone formData string true "Phone number of the user"
// @Param password formData string true "User password"
// @Success 200 {string} string "message: User logged in successfully and cookie stored"
// @Failure 400 {string} string "error: Invalid phone number or password"
// @Router /user/login [post]
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

// Products godoc
// @Summary Get a list of products
// @Description Retrieve a list of products with pagination
// @ID getProducts
// @Tags User Products
// @Produce json
// @Param page query string false "Page number for pagination (default: 1)"
// @Param limit query string false "Limit the number of products per page (default: 10)"
// @Success 200 {string} string "userId: <userID>, products: []entity.Product"
// @Failure 400 {string} string "error: userId not found in the context"
// @Failure 400 {string} string "error: Failed to convert string to integer (page or limit)"
// @Failure 400 {string} string "error: Product not found"
// @Router /user/products [get]
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

// ProductDetails godoc
// @Summary Get details of a specific product
// @Description Retrieve details of a product based on the provided product ID
// @ID getProductDetails
// @Tags User Products
// @Produce json
// @Param productid path string true "Product ID to get details for"
// @Success 200 {string} string "products: entity.Product, product details: entity.ProductDetails"
// @Failure 400 {string} string "error: Failed to convert string to integer (product ID)"
// @Failure 400 {string} string "error: Product not found"
// @Router /user/products/details/{productid} [get]
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

// AddToCart godoc
// @Summary Add a product to the user's cart
// @Description Add a product to the user's cart based on the provided product ID and quantity
// @ID addToCart
// @Tags User Products
// @Accept multipart/form-data
// @Produce json
// @Param productid formData string true "Product ID to add to the cart"
// @Param quantity formData string true "Quantity of the product to add to the cart"
// @Success 200 {string} string "message: Product added to cart successfully, addedproduct: []entity.CartItem"
// @Failure 400 {string} string "error: userId not found in the context"
// @Failure 400 {string} string "error: Failed to convert string to integer (product ID or quantity)"
// @Failure 400 {string} string "error: Failed to add product to cart"
// @Failure 500 {string} string "error: Failed to retrieve cart items"
// @Router /user/cart [post]
func (ac *UserHandler) AddToCart(c *gin.Context) {
	userID, exists := c.Get("userId")
	if !exists || userID == nil {
		// Handle the case where "userid" is not set in the context or is nil
		c.JSON(http.StatusBadRequest, gin.H{"error": "userid not found in the context"})
		return
	}
	userid := userID.(int)
	fmt.Println("userId", userid)

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

	err = ac.CartUSeCase.ExecuteAddToCart(id, quantity, userid)
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

// RemoveFromCart godoc
// @Summary Remove a product from the user's cart
// @Description Remove a product from the user's cart based on the provided cart item ID
// @ID removeFromCart
// @Tags User Products
// @Produce json
// @Param id path string true "Cart item ID to remove from the cart"
// @Success 200 {string} string "message: Product removed from cart"
// @Failure 400 {string} string "error: userId not found in the context"
// @Failure 400 {string} string "error: Failed to convert string to integer (cart item ID)"
// @Failure 400 {string} string "error: Failed to remove product from cart"
// @Router /user/cart/{id} [delete]
func (rc *UserHandler) RemoveFromCart(c *gin.Context) {
	userID, _ := c.Get("userId")
	userid := userID.(int)
	id := c.Param("id")

	Id, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "str convertion failed"})
		return
	}
	err1 := rc.CartUSeCase.ExecuteRemoveCartItem(userid, Id)
	if err1 != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err1.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "product removed from cart"})
}

// Cart godoc
// @Summary Get the user's cart
// @Description Retrieve the user's cart based on the provided user ID
// @ID getCart
// @Tags User Products
// @Produce json
// @Success 200 {string} string "usercart: entity.Cart"
// @Failure 400 {string} string "error: userId not found in the context"
// @Failure 400 {string} string "error: Failed to retrieve user's cart"
// @Router /user/cart [get]
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

// AddToWishList handles the endpoint to add a product to the user's wishlist.
// @Summary Add a product to the wishlist
// @Description Adds the specified product to the user's wishlist.
// @ID addToWishList
// @Accept multipart/form-data
// @Tags User Products
// @Produce json
// @Param productid formData int true "Product ID to add to wishlist"
// @Success 200 {string} string "product added to wishlist"
// @Failure 400 {string} string "Bad Request: string conc failed"
// @Failure 400 {string} string "Bad Request: error message"
// @Failure 500 {string} string "Internal Server Error: failed to retrieve wishlist items"
// @Router /user/wishlist [post]
func (cu *UserHandler) AddToWishList(c *gin.Context) {
	userID, _ := c.Get("userId")
	userid := userID.(int)
	strId := c.PostForm("productid")
	id, err := strconv.Atoi(strId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "string conc failed"})
		return
	}
	err1 := cu.CartUSeCase.ExecuteAddWishlist(id, userid)
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

// RemoveFromWishlist handles the endpoint to remove a product from the user's wishlist.
// @Summary Remove a product from the wishlist
// @Description Removes the specified product from the user's wishlist.
// @ID removeFromWishlist
// @Accept multipart/form-data
// @Tags User Products
// @Produce json
// @Param id path int true "Product ID to remove from wishlist"
// @Success 200 {object} string "message": "successfully removed from wishlist"
// @Failure 400 {object} string "error": "Bad Request: error message"
// @Failure 500 {object} string "error": "Internal Server Error: failed to remove from wishlist"
// @Router /user/wishlist/{id} [delete]
func (cu *UserHandler) RemoveFromWishlist(c *gin.Context) {
	userID, _ := c.Get("userId")
	userid := userID.(int)
	ID := c.Param("id")
	id, err := strconv.Atoi(ID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err1 := cu.CartUSeCase.ExecuteRemoveFromWishList(id, userid)
	if err1 != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err1.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "succesfully removed from wishlist"})
}

// ViewWishlist handles the endpoint to view the user's wishlist.
// @Summary View user's wishlist
// @Description Retrieves and returns the products in the user's wishlist.
// @ID viewWishlist
// @Tags User Products
// @Produce json
// @Success 200 {string} string "wishlist retrieved successfully"
// @Failure 400 {string} string "Bad Request: error message"
// @Failure 500 {string} string "Internal Server Error: failed to retrieve wishlist"
// @Router /user/wishlist [get]
func (cu *UserHandler) ViewWishlist(c *gin.Context) {
	userID, _ := c.Get("userId")
	userid := userID.(int)
	wishlist, err := cu.CartUSeCase.ExecuteViewWishlist(userid)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
	c.JSON(http.StatusOK, gin.H{"message": wishlist})
}

// Logout godoc
// @Summary Logs out the user
// @Description Deletes the authentication token cookie to log the user out
// @Tags User 
// @Produce json
// @Success 200 {string} string "user logged out successfully"
// @Failure 400 {string} string "cookie delete failed"
// @Router /logout [post]
func (cu *UserHandler) Logout(c *gin.Context) {
	err := middleware.DeleteToken(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errror": "cookie delete failed"})
		return
	} else {
		c.JSON(http.StatusOK, gin.H{"message": "user logged out succesfully"})
	}
}

// AddAddress godoc
// @Summary Adds a new address for the user
// @Description Adds a new address associated with the authenticated user
// @Produce json
// @Tags User Address
// @Param address body entity.UserAddress true "Address information to be added"
// @Success 200 {string} string "address added successfully"
// @Failure 400 {string} string "Bad Request"
// @Router /user/address [post]
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

// ShowUserDetails godoc
// @Summary Retrieve details for the authenticated user
// @Description Retrieve user details and associated address based on the provided user ID
// @Produce json
// @Tags User
// @Success 200 {string} string "userdetails: entity.UserDetails, address: entity.UserAddress"
// @Failure 400 {string} string "error: Bad Request"
// @Router /user/details [get]
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

// EditProfile godoc
// @Summary Edit the profile for the authenticated user
// @Description Edit user profile based on the provided user ID and input data
// @Produce json
// @Tags User
// @Param userinput body models.EditUser true "User information to be edited"
// @Success 200 {string} string "message: user edited successfully, updateduser: entity.UserDetails"
// @Failure 400 {string} string "error: Bad Request"
// @Failure 500 {string} string "error: Failed to fetch updated user details"
// @Router /user/profile [patch]
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

// ChangePassword godoc
// @Summary Request OTP for changing user password
// @Description Initiates the process of changing the user password by sending an OTP.
// @ID change-password
// @Accept json
// @Tags User 
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {string} string "OTP sent successfully for password change"
// @Failure 400 {string} string "Bad request: Unable to initiate password change"
// @Router /user/change-password [post]
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

// OtpValidationPassword godoc
// @Summary Validate OTP and change user password
// @Description Validates the provided OTP and changes the user password.
// @ID otp-validation-password
// @Accept multipart/form-data
// @Produce json
// @Tags User 
// @Security ApiKeyAuth
// @Param password formData string true "New password for the user"
// @Param otp formData string true "OTP for validation"
// @Success 200 {string} string "Password changed successfully"
// @Failure 400 {string} string "Bad request: Unable to change password"
// @Router /user/change-password/validation [post]
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

//	 CartItems godoc
//		@Summary Get the items in the user's cart
//		@Description Retrieve the items in the user's cart based on the provided user ID
//		@Produce json
// @Tags User Products
//		@Param userId path int true "User ID"
//		@Success 200 {string} string "cartlist: []entity.CartItem"
//		@Failure 400 {string} string "error: Bad Request"
//		@Router /user/cartlist [get]
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

// EditAddress godoc
//
//		@Summary Edit the user's address
//		@Description Edit the user's address of a specific type (e.g., home, work)
//		@Produce json
//	 @Tags User Address
//		@Param type path string true "Address type (e.g., home, work)"
//		@Param useraddress body entity.UserAddress true "Updated address information"
//		@Success 200 {string} string "success: address edited successfully"
//		@Failure 400 {string} string "error: Bad Request"
//		@Router /user/address/{type} [patch]
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

// @Summary Delete user address
// @Description Deletes a specific type of address for the authenticated user.
// @ID delete-user-address
// @Tags User Address
// @Produce json
// @Param type path string true "Type of address to be deleted (e.g., 'home', 'work')"
// @Success 200 {object} string "success": "address Deleted successfully" "Successful response"
// @Failure 400 {object} string "error": "Error message" "Error response"
// @Router /user/address/{type} [delete]
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

// SearchProduct godoc
// @Summary Search products
// @Description Searches for products based on the provided search query.
// @ID search-products
// @Tags User Sort
// @Produce json
// @Param page query int false "Page number for pagination (default is 1)"
// @Param limit query int false "Number of items per page (default is 5)"
// @Param search query string true "Search query string"
// @Success 200 {string} entity.product
// @Failure 400 {object} string "error": "Error message" "Error response"
// @Router /user/products/search [get]
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

// SortByCategory godoc
// @Summary Sort products by category
// @Description Retrieves a list of products sorted by category based on the provided category ID.
// @ID sort-products-by-category
// @Tags User Sort
// @Produce json
// @Param page query int false "Page number for pagination (default is 1)"
// @Param limit query int false "Number of items per page (default is 5)"
// @Param id query int true "Category ID for sorting products"
// @Success 200 {string} string "Product list retrieved successfully"
// @Failure 400 {string} string "Bad request"
// @Failure 502 {string} string "Bad gateway"
// @Router /user/products/sort [get]
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

// SortByFilter godoc
// @Summary Sort products by filter
// @Description Retrieves a list of products based on the provided filter criteria.
// @ID sort-products-by-filter
// @Tags User Sort
// @Produce json
// @Param minprize query int false "Minimum prize for product filtering"
// @Param maxprize query int false "Maximum prize for product filtering"
// @Param category query int false "Category ID for product filtering"
// @Param size query string false "Product size for filtering"
// @Success 200 {string} string "Product list retrieved successfully"
// @Failure 400 {string} string "Bad request"
// @Router /user/products/filter [get]
func (sc *UserHandler) SortByFilter(c *gin.Context) {
	strminPrize := c.Query("minprize")
	strmaxPrize := c.Query("maxprize")
	strcategory := c.Query("category")
	size := c.Query("size")
	minPrize, _ := strconv.Atoi(strminPrize)
	maxPrize, _ := strconv.Atoi(strmaxPrize)
	category, _ := strconv.Atoi(strcategory)
	productlist, err1 := sc.ProductUseCase.ExecuteProductFilter(size, minPrize, maxPrize, category)
	if err1 != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err1.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"products": productlist})
}

// ApplyCoupon godoc
// @Summary Apply coupon to user's cart
// @Description Applies a coupon to the authenticated user's cart based on the provided coupon code.
// @ID apply-coupon
// @Accept multipart/form-data
// @Tags User Coupon
// @Produce json
// @Param code formData string true "Coupon code to be applied"
// @Success 200 {string} string "Total offer prize and Coupon applied successfully"
// @Failure 400 {string} string "Bad request"
// @Router /user/cart/coupon [post]
func (sc *UserHandler) ApplyCoupon(c *gin.Context) {
	userID, _ := c.Get("userId")
	userid := userID.(int)
	code := c.PostForm("code")

	totaloffer, err := sc.CartUSeCase.ExecuteApplyCoupon(userid, code)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"offer prize": totaloffer, "offer": "applied succesfully"})
}

// AvailableCoupons godoc
// @Summary Retrieve available coupons
// @Description Retrieves a list of available coupons.
// @ID get-available-coupons
// @Tags User Coupon
// @Produce json
// @Success 200 {string} string "List of available coupons"
// @Failure 400 {string} string "Bad request"
// @Router /user/coupons [get]
func (sc *UserHandler) AvailableCoupons(c *gin.Context) {
	couponlist, err := sc.ProductUseCase.ExecuteAvailableCoupons()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"Available Coupons": couponlist})
}
