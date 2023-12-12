package handlers

import (
	"fmt"
	"io"
	"net/http"
	"project/delivery/middleware"
	"project/domain/entity"
	usecase "project/usecase/admin"
	product "project/usecase/product"

	"strconv"

	"github.com/gin-gonic/gin"
)

type AdminHandler struct {
	AdminUseCase   *usecase.AdminUseCase
	ProductUseCase *product.ProductUseCase
}

func NewAdminHandler(AdminUsecase *usecase.AdminUseCase, ProductUsecase *product.ProductUseCase) *AdminHandler {
	return &AdminHandler{AdminUsecase, ProductUsecase}
}

// @Summary Admin Login with Password
// @Description Authenticate admin using email and password and generate an authentication token.
// @ID admin-login
// @Tags Admin
// @Accept json
// @Produce json
// @Param			admin	body		entity.AdminLogin	true	"Admin Data"
// @Success 200 {object} string "message": "Admin logged in successfully"
// @Failure 400 {object} string "error": "Empty request body"
// @Router /admin/login [post]
func (uh *AdminHandler) AdminLoginWithPassword(c *gin.Context) {
	var payload map[string]interface{}
	if err := c.ShouldBindJSON(&payload); err != nil {
		if err == io.EOF {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Empty request body"})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error ": "invalid request payload"})
		return
	}
	email, _ := payload["email"].(string)
	password, _ := payload["password"].(string)

	adminId, err := uh.AdminUseCase.ExecuteAdminLoginWithPassword(email, password)
	if err != nil {
		fmt.Printf("Authentication failed: %v\n", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Authentication failed", "details": err.Error()})
		return
	} else {
		middleware.CreateToken(adminId, email, "admin", c)
	}
	c.JSON(http.StatusOK, gin.H{"message": "Admin logged in succesfully"})
}

// @Summary List Users
// @Description Get a paginated list of users.
// @ID list-users
// @Accept json
// @Tags Admin User Management
// @Produce json
// @Param page query int false "Page number (default is 1)"
// @Param limit query int false "Number of users per page (default is 5)"
// @Success 200 {object} entity.ListUsersResponse
// @Failure 400 {object} entity.ErrorResponse
// @Router /admin/users [get]
func (ul *AdminHandler) UsersList(c *gin.Context) {
	pagestr := c.DefaultQuery("page", "1")
	page, err := strconv.Atoi(pagestr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid page parameter"})
		return
	}
	limitstr := c.DefaultQuery("limit", "5")
	limit, err := strconv.Atoi(limitstr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid limit parameter"})
		return
	}
	userlist, err1 := ul.AdminUseCase.ExecuteUsersList(page, limit)
	if err1 != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user list nor found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"users": userlist})
}

// @Summary Toggle User Permission
// @Description Toggle the permission of a user by providing the user's ID.
// @ID toggle-user-permission
// @Accept json
// @Tags Admin User Management
// @Produce json
// @Param id path int true "User ID" minimum(1) format(int32)
// @Success 200 {string} string "success: User permission toggled successfully"
// @Failure 400 {string} string "error: Invalid user ID"
// @Failure 401 {string} string "error: User not found"
// @Router /admin/users/toggle-permission/{id} [put]
func (tp *AdminHandler) TogglePermission(c *gin.Context) {
	ids := c.Param("id")
	fmt.Println("ids:", ids)
	id, err := strconv.Atoi(ids)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "str conv failed"})
		return
	}

	err1 := tp.AdminUseCase.ExecuteTogglePermission(id)
	if err1 != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": "user permission toogled"})
	return
}

// @Summary Create a new category
// @Description Create a new category by providing the category details.
// @ID create-category
// @Accept json
// @Tags Admin Category Management
// @Produce json
// @Param category body entity.Category true "Category details"
// @Success 200 {object} string "success": "Category added successfully" entity.Category
// @Failure 400 {object} string "error": "Invalid input" entity.ErrorResponse
// @Router /admin/categories [post]
func (ct *AdminHandler) CreateCategory(c *gin.Context) {
	var input entity.Category
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	category := entity.Category{
		Name:        input.Name,
		Description: input.Description,
	}

	if input.Name == "" || input.Description == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "category name and description cannot be empty"})
		return
	}

	newcat, err := ct.ProductUseCase.ExecuteCreateCategory(category)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	catlist, err := ct.ProductUseCase.ExecuteGetCategoryId(newcat)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": "category added succesfully", "category ": catlist})
}

// EditCategory godoc
// @Summary Edit a category
// @Description Edit a category based on the provided JSON data
// @ID editCategory
// @Tags Admin Category Management
// @Accept json
// @Produce json
// @Param id path int true "Category ID" Format(int64)
// @Param category body entity.Category true "Category object to be edited"
// @Success 200 {object} string "success": "product edited successfully", "edited category": entity.Category
// @Failure 400 {object} string "error": "str conversion failed"
// @Failure 400 {object} string "error": "JSON binding failed"
// @Failure 400 {object} string "error": "editing category failed"
// @Router /admin/categories/{id} [put]
func (et *AdminHandler) EditCategory(c *gin.Context) {
	var category entity.Category
	id := c.Param("id")
	Id, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "str conversion failed"})
		return
	}
	if err := c.ShouldBindJSON(&category); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err1 := et.ProductUseCase.ExecuteEditCategory(category, Id)
	if err1 != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err1.Error()})
		return
	}
	catlist, err := et.ProductUseCase.ExecuteGetCategoryId(Id)

	c.JSON(http.StatusOK, gin.H{"success": "product edited succesfully", "edited category": catlist})
}

// DeleteCategory godoc
// @Summary Delete a category
// @Description Delete an existing category based on the provided ID
// @ID deleteCategory
// @Tags Admin Category Management
// @Accept json
// @Produce json
// @Param id path int true "Category ID" Format(int64)
// @Success 200 {object} string "success": "Category deleted successfully"
// @Failure 400 {object} string "error": "Invalid category ID"
// @Failure 400 {object} string "error": "Failed to delete category"
// @Router /admin/categories/{id} [delete]
func (et *AdminHandler) DeleteCategory(c *gin.Context) {
	id := c.Param("id")
	Id, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "str conversion failed"})
	}

	err1 := et.ProductUseCase.ExecuteDeleteCategory(Id)
	if err1 != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err1.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": "category deleted successfully"})
}

// @Summary Get all categories
// @Description Retrieve a list of all categories
// @Tags Admin Category Management
// @Accept json
// @Produce json
// @Success 200 {object} entity.Category "List of categories"
// @Failure 400 {object} entity.ErrorResponse "Bad Request"
// @Router /admin/categories [get]
func (et *AdminHandler) AllCategory(c *gin.Context) {
	categorylist, err := et.ProductUseCase.ExecuteGetAllCategory()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"Categories": categorylist})
}

// CreateProduct godoc
// @Summary Create a new product
// @Description Create a new product with details, including image upload
// @ID createProduct
// @Tags Admin Product Management
// @Accept multipart/form-data
// @Produce json
// @Param name formData string true "Product Name"
// @Param price formData number true "Product Price"
// @Param size formData string true "Product Size"
// @Param category formData int true "Category ID"
// @Param description formData string true "Product Description"
// @Param specification formData string true "Product Specification"
// @Param image formData file true "Product Image"
// @Param imageURL formData string false "Product Image URL"
// @Param quantity formData int true "Product Quantity"
// @Success 200 {string} string "Product added successfully" "products":entity.products
// @Failure 400 {string} string "Invalid input data"
// @Failure 400 {string} string "Failed to get category"
// @Failure 400 {string} string "Failed to create product"
// @Failure 400 {string} string "Failed to create product details and rollback"
// @Failure 400 {string} string "Failed to create product details"
// @Failure 400 {string} string "Failed to create inventory"
// @Failure 400 {string} string "Failed to get product by ID"
// @Router /admin/products [post]
func (cp *AdminHandler) CreateProduct(c *gin.Context) {
	var input entity.ProductInput
	if err := c.ShouldBind(&input); err != nil {
		fmt.Println("err", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	image, _ := c.FormFile("image")
	category, err := cp.ProductUseCase.ExecuteGetCategory(entity.Category{ID: input.Category})
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	product := &entity.Product{
		Name:     input.Name,
		Price:    input.Price,
		Size:     input.Size,
		Category: category,
	}
	productId, err := cp.ProductUseCase.ExecuteCreateProduct(*product, image)
	if err != nil {

		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	} else {
		ProductDetails := entity.ProductDetails{
			ProductID:     productId,
			Description:   input.Description,
			Specification: input.Specification,
		}
		err := cp.ProductUseCase.ExecuteCreateProductDetails(ProductDetails)
		if err != nil {
			err2 := cp.ProductUseCase.ExecuteDeleteProductAdd(productId)
			if err2 != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err2.Error()})
				return
			}
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		inventory := entity.Inventory{
			ProductId:       productId,
			ProductCategory: category,
			Quantity:        input.Quantity,
		}
		err = cp.ProductUseCase.ExecuteCreateInventory(inventory)
		if err != nil {
			err2 := cp.ProductUseCase.ExecuteDeleteProductAdd(productId)
			if err2 != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err2.Error()})
				return
			}
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

	}
	prod, err := cp.ProductUseCase.ExecuteGetProductById(productId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": "product added succesfully", "product": prod})

}

// EditProduct godoc
// @Summary Edit a product
// @Description Edit an existing product based on the provided JSON data
// @ID editProduct
// @Tags Admin Product Management
// @Accept json
// @Produce json
// @Param id path int true "Product ID" Format(int64)
// @Param category body entity.Product true "Product object to be edited"
// @Success 200 {string} string "Product edit success" "product":entity.product
// @Failure 400 {string} string "String conversion error"
// @Failure 400 {string} string "JSON binding failed"
// @Failure 400 {string} string "Product edit failed"
// @Router /admin/products/{id} [put]
func (ep *AdminHandler) EditProduct(c *gin.Context) {
	var product entity.Product
	ids := c.Param("id")
	id, err := strconv.Atoi(ids)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "string conv err"})
		return
	}
	if err := c.ShouldBindJSON(&product); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err1 := ep.ProductUseCase.ExecuteEditProduct(product, id)
	if err1 != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err1.Error()})
		return
	}
	prod, err := ep.ProductUseCase.ExecuteGetProductById(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err1.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"succes": "product edit success", "product": prod})
}

// DeleteProduct godoc
// @Summary Delete a product
// @Description Delete an existing product based on the provided ID
// @Tags Admin Product Management
// @Tags Admin Product Management
// @Accept json
// @Produce json
// @Param id path int true "Product ID" Format(int64)
// @Success 200 {string} string "Product deleted"
// @Failure 400 {string} string "String conversion failed"
// @Failure 400 {string} string "Product not found"
// @Router /admin/products/{id} [delete]
func (dp *AdminHandler) DeleteProduct(c *gin.Context) {
	ids := c.Param("id")
	id, err := strconv.Atoi(ids)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "str conv failed"})
	}
	err1 := dp.ProductUseCase.ExecuteDeleteProduct(id)
	if err1 != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "product not found"})
	} else {
		c.JSON(http.StatusOK, gin.H{"succes": "product deleted"})
	}
}

// @Summary Get a list of products for admin
// @Description Retrieve a list of products for the admin dashboard.
// @ID get-admin-products
// @Tags Admin Product Management
// @Produce json
// @Success 200 {array} models.ProductWithQuantityResponse
// @Failure 401 {string} string "Unauthorized"
// @Router /admin/products [get]
func (pl *AdminHandler) AdminProductlist(c *gin.Context) {
	pagestr := c.DefaultQuery("page", "1")
	page, err := strconv.Atoi(pagestr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "string conversion failed"})
		return
	}
	limitstr := c.DefaultQuery("limit", "5")
	limit, err := strconv.Atoi(limitstr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "string conversion failed"})
		return
	}
	productlist, err1 := pl.ProductUseCase.ExecuteProductList(page, limit)
	if err1 != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err1.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"products": productlist})
}

// Home godoc
// @Summary Get home information
// @Description Get information about available options and the admin dashboard
// @ID home
// @Tags Admin
// @Produce json
// @Success 200 {string} string "Options: SalesReport - User Management - Product Management - Order Management. Dashboard information is also included." "Dashboard": entity.AdminDashboard
// @Failure 400 {string} string "Failed to retrieve dashboard information"
// @Router /admin/home [get]
func (ad *AdminHandler) Home(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"options": "SalesReport - User Mangement - Product Management -Order Management"})
	dashboardresponse, err := ad.AdminUseCase.ExecuteAdminDashBoard()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"dashboard": dashboardresponse})
}

// AddCoupon godoc
// @Summary Add a new coupon
// @Description Add a new coupon with the provided JSON data
// @ID addCoupon
// @Tags Admin Coupon Management
// @Accept json
// @Produce json
// @Param category body entity.Coupon true "Coupon details"
// @Success 200 {string} string "Coupon added successfully" "coupon": entity.coupon
// @Failure 400 {string} string "JSON binding failed"
// @Failure 400 {string} string "Coupon code already exists"
// @Failure 400 {string} string "Failed to add coupon"
// @Router /admin/coupons [post]
func (ad *AdminHandler) AddCoupon(c *gin.Context) {
	var coupon entity.Coupon

	if err := c.ShouldBindJSON(&coupon); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	exisitingCode, _ := ad.ProductUseCase.ExecuteGetCouponByCode(coupon.Code)

	if exisitingCode != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "code already exists"})
		return
	}

	err := ad.ProductUseCase.ExecuteAddCoupon(&coupon)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	coup, err := ad.ProductUseCase.ExecuteGetCouponByCode(coupon.Code)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "coupon added succesfully", "coupon": coup})
}

// AllCoupons godoc
// @Summary Get all available coupons
// @Description Get a list of all available coupons
// @ID allCoupons
// @Tags Admin Coupon Management
// @Produce json
// @Success 200 {string} string "List of available coupons" entity.coupon
// @Failure 400 {string} string "Failed to retrieve available coupons"
// @Router /admin/coupons [get]
func (ad *AdminHandler) AllCoupons(c *gin.Context) {
	couponlist, err := ad.ProductUseCase.ExecuteAvailableCoupons()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errro": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"avialable coupons": couponlist})
}

// DeleteCoupon godoc
// @Summary Delete a coupon
// @Description Delete an existing coupon based on the provided code
// @ID deleteCoupon
// @Tags Admin Coupon Management
// @Accept multipart/form-data
// @Produce json
// @Param code formData string true "Coupon code to be deleted"
// @Success 200 {string} string "Coupon successfully deleted"
// @Failure 400 {string} string "Failed to delete coupon"
// @Router /admin/coupons [delete]
func (ad *AdminHandler) DeleteCoupon(c *gin.Context) {
	code := c.PostForm("code")
	err := ad.ProductUseCase.ExecuteDeleteCoupon(code)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "coupon succesfully deleted"})
}

// AddProductOffer godoc
// @Summary Add an offer to a product
// @Description Add an offer to a product based on the provided product ID and offer
// @ID addProductOffer
// @Tags Admin Offer Management
// @Accept multipart/form-data
// @Produce json
// @Param productid formData string true "Product ID to add the offer to"
// @Param offer formData string true "Offer to be added to the product"
// @Success 200 {string} string "Product offer added successfully"
// @Failure 400 {string} string "Failed to convert string to integer"
// @Failure 400 {string} string "Failed to add product offer"
// @Router /admin/product/offer [post]
func (ad *AdminHandler) AddProductOffer(c *gin.Context) {
	strpro := c.PostForm("productid")
	stroffer := c.PostForm("offer")
	productid, err := strconv.Atoi(strpro)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "str conv failed"})
		return
	}
	offer, err := strconv.Atoi(stroffer)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "str conv failed"})
		return
	}

	prod, err1 := ad.ProductUseCase.ExecuteAddProductOffer(productid, offer)
	if err1 != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err1.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"offer added ": prod})
}

// AddCategoryOffer godoc
// @Summary Add an offer to a category
// @Description Add an offer to a category based on the provided category ID and offer
// @ID addCategoryOffer
// @Tags Admin Offer Management
// @Accept multipart/form-data
// @Produce json
// @Param categoryid formData string true "form" "Category ID to add the offer to"
// @Param offer formData string true "form" "Offer to be added to the category"
// @Success 200 {string} string "offer added: Category offer added successfully"
// @Failure 400 {string} string "error: Failed to convert string to integer"
// @Failure 400 {string} string "error: Failed to add category offer"
// @Router /admin/category/offer [post]
func (ad *AdminHandler) AddCategoryOffer(c *gin.Context) {
	strcat := c.PostForm("categoryid")
	stroffer := c.PostForm("offer")
	categoryid, err := strconv.Atoi(strcat)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "str1 conv failed"})
		return
	}
	offer, err := strconv.Atoi(stroffer)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "str conv failed"})
		return
	}
	productlist, err1 := ad.ProductUseCase.ExecuteCategoryOffer(categoryid, offer)
	if err1 != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err1.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"offer addded": productlist})
}

// StocklessProducts godoc
// @Summary Get a list of stockless products
// @Description Retrieve a list of products with zero stock
// @ID stocklessProducts
// @Tags Admin Product Management
// @Produce json
// @Success 200 {string} string "List of stockless products: []entity.Product"
// @Failure 400 {string} string "error: Failed to retrieve stockless products"
// @Router /admin/stockless/products [get]
func (au *AdminHandler) StocklessProducts(c *gin.Context) {
	prod, err := au.AdminUseCase.ExecuteStocklessProducts()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"Stockless Products": prod})
}

// SearchUsers godoc
// @Summary Search users based on criteria
// @Description Retrieve a list of users based on search criteria, paginated with optional limit
// @ID searchUsers
// @Tags Admin User Management
// @Produce json
// @Param page query string false "Page number for pagination (default: 1)"
// @Param limit query string false "Limit the number of users per page (default: 5)"
// @Param search query string true "Search criteria to filter users"
// @Success 200 {string} string "userlist: []entity.User"
// @Failure 400 {string} string "error: Failed to convert string to integer"
// @Failure 400 {string} string "error: Invalid search criteria"
// @Failure 400 {string} string "error: Failed to retrieve userlist"
// @Router /admin/search/users [get]
func (or *AdminHandler) SearchUsers(c *gin.Context) {
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
	userlist, err := or.AdminUseCase.ExecuteUserSearch(page, limit, search)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"userlist": userlist})
}

// @Summary Add stock to a product
// @Description Add stock to a product based on the provided ID and quantity
// @Tags Admin Product Management
// @Accept json
// @Produce json
// @Param id path int true "Product ID"
// @Param request body entity.Inventory true "Stock details to be added"
// @Success 200 {object} entity.Inventory "Updated inventory"
// @Failure 400 {string} string "Bad Request"
// @Router /admin/products/stocks/{id} [put]
func (or *AdminHandler) AddStock(c *gin.Context) {
	var inventory *entity.Inventory
	strid := c.Param("id")
	id, err1 := strconv.Atoi(strid)
	if err1 != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "str conv failed"})
		return
	}
	if err := c.ShouldBindJSON(&inventory); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	quantity := int(inventory.Quantity)
	inventory, err := or.ProductUseCase.ExecuteAddStock(id, quantity)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"Updated result :": inventory})
}

// Logout godoc
// @Summary Logs out the Admin
// @Description Deletes the authentication token cookie to log the admin out
// @Tags Admin
// @Produce json
// @Success 200 {string} string "Admin logged out successfully"
// @Failure 400 {string} string "cookie delete failed"
// @Router /admin/logout [post]
func (cu *AdminHandler) Logout(c *gin.Context) {
	err := middleware.DeleteToken(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errror": "cookie delete failed"})
		return
	} else {
		c.JSON(http.StatusOK, gin.H{"message": "Admin logged out succesfully"})
	}
}
