package handlers

import (
	"net/http"
	"project/delivery/middleware"
	"project/delivery/models"
	"project/domain/entity"
	Productusecase "project/usecase/product"
	usecase "project/usecase/user"
	"strconv"

	"github.com/gin-gonic/gin"
	
)

type UserHandler struct {
	UserUseCase    *usecase.UserUseCase
	ProductUseCase *Productusecase.ProductUseCase
}

func NewUserhandler(UserUseCase *usecase.UserUseCase, ProductUseCase *Productusecase.ProductUseCase) *UserHandler {
	return &UserHandler{UserUseCase, ProductUseCase}
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
		middleware.CreateToken(userId, phone, "user", c)
		c.JSON(http.StatusOK, gin.H{"message": "user logged in succesfully and cookie stored"})
	}

}



func (po *UserHandler) Products(c *gin.Context) {
	pagestr := c.DefaultQuery("page", "1")
	limitstr := c.DefaultQuery("limit", "5")
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
	c.JSON(http.StatusOK, gin.H{"products": prodlist})
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
