package repository

import (
	"errors"
	"project/domain/entity"

	"gorm.io/gorm"
)

type ProductRepository struct {
	db *gorm.DB
}

func NewProductRepository(db *gorm.DB) *ProductRepository {
	return &ProductRepository{db}
}

func (pr *ProductRepository) GetAllProducts(offset, limit int) ([]entity.Product, error) {
	var Products []entity.Product
	err := pr.db.Offset(offset).Limit(limit).Where("removed=?", false).Find(&Products).Error
	if err != nil {
		return nil, err
	}
	return Products, nil
}

func (au *ProductRepository) GetProductDetailsById(id int) (*entity.ProductDetails, error) {
	var productDetails entity.ProductDetails
	result := au.db.Where("product_id=?", id).Find(&productDetails)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, result.Error
		}
		return nil, result.Error
	}
	return &productDetails, nil
}

func (pr *ProductRepository) GetProductById(id int) (*entity.Product, error) {
	var product entity.Product
	result := pr.db.First(&product, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, result.Error
		}
		return nil, result.Error
	}
	return &product, nil
}

func (pn *ProductRepository) GetProductByName(name string) error {
	var prodname entity.Product
	result := pn.db.Where(&entity.Product{Name: name}).Find(&prodname)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return result.Error
		}
		return result.Error
	}
	return nil

}
func (ct *ProductRepository) CreateProduct(product *entity.Product) (int, error) {
	if err := ct.db.Create(product).Error; err != nil {
		return 0, err
	}
	return product.ID, nil
}

func (up *ProductRepository) UpdateProduct(product *entity.Product) error {
	return up.db.Save(product).Error
}
func (dp *ProductRepository) DeleteProduct(product *entity.Product) error {
	return dp.db.Delete(product).Error
}
func (up *ProductRepository) CreateProductDetails(details *entity.ProductDetails) error {
	return up.db.Create(details).Error
}

func (cc *ProductRepository) CreateCategory(category *entity.Category) (int, error) {
	if err := cc.db.Create(category).Error; err != nil {
		return 0, errors.New("error creating category")
	}
	return category.ID, nil
}

func (uc *ProductRepository) UpdateCategory(category *entity.Category) error {
	return uc.db.Save(category).Error
}

func (cn *ProductRepository) GetCategoryByName(name string) error {
	var prodname entity.Category
	result := cn.db.Where(&entity.Category{Name: name}).Find(&prodname)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return result.Error
		}
		return errors.New("failed to get category by name: %v")
	}
	return nil

}
func (cr *ProductRepository) GetCategoryById(id int) (*entity.Category, error) {
	var category entity.Category
	result := cr.db.First(&category, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, result.Error
		}
		return nil, result.Error
	}
	return &category, nil
}

func (cn *ProductRepository) DeleteCategory(Id int) error {

    err := cn.db.Table("categories").Where("id = ?", Id).Delete(&entity.Category{})
	if err != nil{
		return errors.New("Coudnt delete")
	}
	return nil
    
}
func (pr *ProductRepository) CreateInventory(inventory *entity.Inventory) error {
	return pr.db.Create(inventory).Error
}