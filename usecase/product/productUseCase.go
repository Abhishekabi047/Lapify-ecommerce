package usecase

import (
	"errors"
	"fmt"
	"project/domain/entity"
	repository "project/repository/product"

	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

type ProductUseCase struct {
	productRepo *repository.ProductRepository
}

func NewProduct(productRepo *repository.ProductRepository) *ProductUseCase {
	return &ProductUseCase{productRepo: productRepo}
}

func (pu *ProductUseCase) ExecuteProductList(page, limit int) ([]entity.Product, error) {
	offset := (page - 1) * limit
	productlist, err := pu.productRepo.GetAllProducts(offset, limit)
	if err != nil {
		return nil, err
	} else {
		return productlist, nil
	}
}

func (pu *ProductUseCase) ExecuteProductDetails(id int) (*entity.Product, *entity.ProductDetails, error) {
	product, err := pu.productRepo.GetProductById(id)
	if err != nil {
		return nil, nil, err
	}
	productdetails, err := pu.productRepo.GetProductDetailsById(id)
	if err != nil {
		return nil, nil, err
	}
	return product, productdetails, nil
}

func (pu *ProductUseCase) ExecuteCreateProduct(product entity.Product) (int, error) {
	validate := validator.New()
	if err := validate.Struct(product); err != nil {
		if _, ok := err.(*validator.InvalidValidationError); ok {
			return 0, err
		}
		errors := err.(validator.ValidationErrors)
		errorMsg := "Validation failed: "
		for _, e := range errors {
			switch e.Tag() {
			case "required":
				errorMsg += fmt.Sprintf("%s is required; ", e.Field())
			case "numeric":
				errorMsg += fmt.Sprintf("%s should contain only numeric characters; ", e.Field())
			default:
				errorMsg += fmt.Sprintf("%s has an invalid value; ", e.Field())
			}
		}
		return 0, fmt.Errorf(errorMsg)
	}
	err := pu.productRepo.GetProductByName(product.Name)
	if err == nil {
		return 0, errors.New("Product already exists")
	}
	newprod := &entity.Product{
		Name:     product.Name,
		Price:    product.Price,
		Category: product.Category,
		Size:     product.Size,
		ImageURL: product.ImageURL,
	}
	productid, err := pu.productRepo.CreateProduct(newprod)
	if err != nil {
		return 0, err
	} else {
		return productid, nil
	}
}

func (pu *ProductUseCase) ExecuteCreateProductDetails(details entity.ProductDetails) error {
	validate := validator.New()
	if err := validate.Struct(details); err != nil {
		if _, ok := err.(*validator.InvalidValidationError); ok {
			return err
		}
		errors := err.(validator.ValidationErrors)
		errorMsg := "Validation failed: "
		for _, e := range errors {
			switch e.Tag() {
			case "required":
				errorMsg += fmt.Sprintf("%s is required; ", e.Field())
			default:
				errorMsg += fmt.Sprintf("%s has an invalid value; ", e.Field())
			}
		}
		return fmt.Errorf(errorMsg)
	}
	productDetails := &entity.ProductDetails{
		ProductID:     details.ProductID,
		Description:   details.Description,
		Specification: details.Specification,
	}
	err := pu.productRepo.CreateProductDetails(productDetails)
	if err != nil {
		return errors.New("creating details failed")
	} else {
		return nil
	}
}

func (pt *ProductUseCase) ExecuteEditProduct(product entity.Product, id int) error {
	validate := validator.New()
	if err := validate.Struct(product); err != nil {
		if _, ok := err.(*validator.InvalidValidationError); ok {
			return  err
		}
		errors := err.(validator.ValidationErrors)
		errorMsg := "Validation failed: "
		for _, e := range errors {
			switch e.Tag() {
			case "required":
				errorMsg += fmt.Sprintf("%s is required; ", e.Field())
			case "numeric":
				errorMsg += fmt.Sprintf("%s should contain only numeric characters; ", e.Field())
			default:
				errorMsg += fmt.Sprintf("%s has an invalid value; ", e.Field())
			}
		}
		return  fmt.Errorf(errorMsg)
	}
	existingProduct, err := pt.productRepo.GetProductById(id)
	if err != nil {
		return err
	}

	existingProduct.Name = product.Name
	existingProduct.Price = product.Price
	existingProduct.Category = product.Category
	existingProduct.Size = product.Size
	existingProduct.ImageURL = product.ImageURL

	err1 := pt.productRepo.UpdateProduct(existingProduct)
	if err1 != nil {
		return err1
	} else {
		return nil
	}
}

func (de *ProductUseCase) ExecuteDeleteProduct(id int) error {
	result, err := de.productRepo.GetProductById(id)
	if err != nil {
		return err
	}
	result.Removed = !result.Removed
	err1 := de.productRepo.UpdateProduct(result)
	if err1 != nil {
		return errors.New("product deleted")
	}
	return nil
}

func (pu *ProductUseCase) ExecuteCreateCategory(category entity.Category) (int, error) {
	validate := validator.New()
	if err := validate.Struct(category); err != nil {
		if _, ok := err.(*validator.InvalidValidationError); ok {
			return 0, err
		}
		errors := err.(validator.ValidationErrors)
		errorMsg := "Validation failed: "
		for _, e := range errors {
			switch e.Tag() {
			case "required":
				errorMsg += fmt.Sprintf("%s is required; ", e.Field())
			case "alpha":
				errorMsg += fmt.Sprintf("%s should contain only alphabetic characters; ", e.Field())
			default:
				errorMsg += fmt.Sprintf("%s has an invalid value; ", e.Field())
			}
		}
		return 0, fmt.Errorf(errorMsg)
	}
	err := pu.productRepo.GetCategoryByName(category.Name)
	if err == nil {
		return 0, errors.New("category already exists")
	}
	newcat := &entity.Category{
		Name:        category.Name,
		Description: category.Description,
	}
	categoryid, err := pu.productRepo.CreateCategory(newcat)
	if err != nil {
		return 0, errors.New("category not created")
	} else {
		return categoryid, nil
	}
}

func (pt *ProductUseCase) ExecuteEditCategory(category entity.Category, id int) error {
	validate := validator.New()
	if err := validate.Struct(category); err != nil {
		if _, ok := err.(*validator.InvalidValidationError); ok {
			return err
		}
		errors := err.(validator.ValidationErrors)
		errorMsg := "Validation failed: "
		for _, e := range errors {
			switch e.Tag() {
			case "required":
				errorMsg += fmt.Sprintf("%s is required; ", e.Field())
			case "alpha":
				errorMsg += fmt.Sprintf("%s should contain only alphabetic characters; ", e.Field())
			default:
				errorMsg += fmt.Sprintf("%s has an invalid value; ", e.Field())
			}
		}
		return fmt.Errorf(errorMsg)
	}
	existingCat, err := pt.productRepo.GetCategoryById(id)
	if err != nil {
		return err
	}

	existingCat.Name = category.Name
	existingCat.Description = category.Description

	err = pt.productRepo.UpdateCategory(existingCat)
	if err != nil {
		return err
	}

	return nil
}

func (pu *ProductUseCase) ExecuteDeleteCategory(Id int) error {

	category, err := pu.productRepo.GetCategoryById(Id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("category does not exist")
		}
		return err
	}

	err = pu.productRepo.DeleteCategory(category.ID)
	if err != nil {
		return err
	}

	return nil
}
func (pu *ProductUseCase) ExecuteGetCategory(category entity.Category) (int, error) {
	name, err := pu.productRepo.GetCategoryById(category.ID)
	if err != nil {
		return 0, errors.New("error getting category")
	}
	return name.ID, err
}
func (pu *ProductUseCase) ExecuteCreateInventory(inventory entity.Inventory) error {
	validate := validator.New()
	if err := validate.Struct(inventory); err != nil {
		if _, ok := err.(*validator.InvalidValidationError); ok {
			return err
		}
		errors := err.(validator.ValidationErrors)
		errorMsg := "Validation failed: "
		for _, e := range errors {
			switch e.Tag() {
			case "required":
				errorMsg += fmt.Sprintf("%s is required; ", e.Field())
			case "numeric":
				errorMsg += fmt.Sprintf("%s should contain only numeric characters; ", e.Field())
			default:
				errorMsg += fmt.Sprintf("%s has an invalid value; ", e.Field())
			}
		}
		return fmt.Errorf(errorMsg)
	}

	err := pu.productRepo.CreateInventory(&inventory)
	if err != nil {
		return errors.New("Creating inventory failed")
	} else {
		return nil
	}

}

func (p *ProductUseCase) ExecuteAddCoupon(coupon *entity.Coupon) error {
	err := p.productRepo.CreateCoupon(coupon)
	if err != nil {
		return errors.New("creating coupon failed")
	} else {
		return nil
	}
}

func (p *ProductUseCase) ExecuteAddOffer(offer *entity.Offer) error {
	err := p.productRepo.CreateOffer(offer)
	if err != nil {
		return errors.New("error creating offer")
	} else {
		return nil
	}
}

func (p *ProductUseCase) ExecuteAvailableCoupons() (*[]entity.Coupon, error) {
	coupons, err := p.productRepo.GetAllCoupons()
	if err != nil {
		return nil, errors.New(err.Error())
	}
	avialablecoup := []entity.Coupon{}
	for _, coupons := range *coupons {
		if coupons.UsageLimit != coupons.UsedCount {
			avialablecoup = append(avialablecoup, coupons)
		}
	}
	return &avialablecoup, nil
}

func (pu *ProductUseCase) ExecuteGetCategoryId(id int) (*entity.Category, error) {
	cat, err := pu.productRepo.GetCategoryById(id)
	if err != nil {
		return nil, errors.New("error getting category")
	}
	return cat, err
}
