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

	fmt.Println(email, password)

	adminId, err := uh.AdminUseCase.ExecuteAdminLoginWithPassword(email, password)
	if err != nil {
		fmt.Printf("Authentication failed: %v\n", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Authentication failed", "details": err.Error()})
		return
	} else {
		middleware.CreateToken(adminId, email, "admin", c)
	}
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

	_, err := ct.ProductUseCase.ExecuteCreateCategory(category)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": "category added succesfully"})
}
func (et *AdminHandler) EditCategory(c *gin.Context) {
	var category entity.Category
	id := c.Param("id")
	Id, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "str conversion failed"})
	}
	if err := c.ShouldBindJSON(&category); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
	err1 := et.ProductUseCase.ExecuteEditCategory(category, Id)
	if err1 != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err1.Error()})
	}
	c.JSON(http.StatusOK, gin.H{"success": "product edited succesfully"})
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
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
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
	productId, err := cp.ProductUseCase.ExecuteCreateProduct(*product)
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
		c.JSON(http.StatusBadRequest, gin.H{"error": "eror editing"})
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
