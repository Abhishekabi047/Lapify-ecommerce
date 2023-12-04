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
func (tp *AdminHandler) TogglePermission(c *gin.Context) {
	ids := c.Param("id")
	id, err := strconv.Atoi(ids)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "str conv failed"})
	}

	err1 := tp.AdminUseCase.ExecuteTogglePermission(id)
	if err1 != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not found"})
	}
	c.JSON(http.StatusOK, gin.H{"success": "user permission toogled"})
}

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
		ImageURL: input.ImageURL,
	}
	productId, err := cp.ProductUseCase.ExecuteCreateProduct(*product,image)
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
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
	}
	c.JSON(http.StatusOK, gin.H{"success": "product added succesfully"})

}
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
	c.JSON(http.StatusOK, gin.H{"succes": "product edit success"})
}

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

func (ad *AdminHandler) Home(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"options": "SalesReport - User Mangement - Product Management -Order Management"})
	dashboardresponse, err := ad.AdminUseCase.ExecuteAdminDashBoard()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"dashboard": dashboardresponse})
}

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
	c.JSON(http.StatusOK, gin.H{"message": "coupon added succesfully"})
}

func (ad *AdminHandler) AllCoupons(c *gin.Context) {
	couponlist, err := ad.ProductUseCase.ExecuteAvailableCoupons()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errro": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"avialable coupons": couponlist})
}

func (ad *AdminHandler) DeleteCoupon(c *gin.Context) {
	code := c.PostForm("code")
	err := ad.ProductUseCase.ExecuteDeleteCoupon(code)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "coupon succesfully deleted"})
}

func (ad *AdminHandler) AddOffer(c *gin.Context) {
	var offer entity.Offer

	if err := c.ShouldBindJSON(&offer); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err := ad.ProductUseCase.ExecuteAddOffer(&offer)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "offer added sussefully"})
}

func (ad *AdminHandler) AllOffer(c *gin.Context) {

	offerlist, err := ad.ProductUseCase.ExecuteGetOffers()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"offers": offerlist})
}

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
