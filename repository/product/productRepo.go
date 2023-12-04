package repository

import (
	"errors"
	"fmt"
	"project/domain/entity"
	"time"

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
	result := pn.db.Where(&entity.Product{Name: name}).First(&prodname)
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
	result := cn.db.Where("name=?", name).First(&prodname)
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

	err := cn.db.Delete(&entity.Category{}, Id).Error
	if err != nil {
		return errors.New("Coudnt delete")
	}
	return nil

}
func (pr *ProductRepository) CreateInventory(inventory *entity.Inventory) error {
	return pr.db.Create(inventory).Error
}

func (pr *ProductRepository) CreateCoupon(coupon *entity.Coupon) error {
	if err := pr.db.Create(coupon).Error; err != nil {
		return err
	}
	return nil
}

func (pr *ProductRepository) GetAllCoupons() (*[]entity.Coupon, error) {
	var coupon []entity.Coupon
	currenttime := time.Now()
	err := pr.db.Where("validuntil > ?", currenttime).Find(&coupon).Error
	if err != nil {
		return nil, err
	}
	return &coupon, nil
}

func (pr *ProductRepository) GetCouponByCode(code string) (*entity.Coupon, error) {
	coupon := &entity.Coupon{}
	err := pr.db.Where("code=?", code).First(coupon).Error
	if err != nil {
		return nil, err
	}
	return coupon, nil
}

func (pr *ProductRepository) UpdateCouponCount(coupon *entity.Coupon) error {
	return pr.db.Save(coupon).Error
}

func (pr *ProductRepository) UpdateCouponUsage(coupon *entity.UsedCoupon) error {
	if err := pr.db.Create(coupon).Error; err != nil {
		return err
	}
	return nil
}

func (pr *ProductRepository) CreateOffer(offer *entity.Offer) error {
	if err := pr.db.Create(offer).Error; err != nil {
		return err
	}
	return nil
}

func (pr *ProductRepository) GetOfferByPrize(Prize int) (*[]entity.Offer, error) {
	offers := &[]entity.Offer{}
	err := pr.db.Where("min_prize=?", Prize).Find(&offers).Error
	if err != nil {
		return nil, err
	} else if offers == nil {
		return nil, err
	}
	return offers, nil
}

func (pr *ProductRepository) DecreaseProductQuantity(product *entity.Inventory) error {
	exisitingproduct := &entity.Inventory{}
	err := pr.db.Where("product_category=? AND product_id=?", product.ProductCategory, product.ProductId).First(exisitingproduct).Error
	if err != nil {
		return err
	}
	if exisitingproduct.Quantity == 0 {
		return errors.New("out of stock")
	}
	newQuantity := exisitingproduct.Quantity - product.Quantity

	if newQuantity < 0 {
		return fmt.Errorf("There is only %d quantity avialable", exisitingproduct.Quantity)
	}

	err = pr.db.Model(exisitingproduct).Update("quantity", newQuantity).Error
	if err != nil {
		return err
	}
	return nil
}

func (pr *ProductRepository) GetCouponByCategory(category string) (*entity.Coupon, error) {
	coupon := &entity.Coupon{}
	err := pr.db.Where("category=?", category).First(coupon).Error
	if err != nil {
		return nil, err
	}
	return coupon, nil
}

func (pr *ProductRepository) DeleteCoupon(code string) error {
	err := pr.db.Where("code=?", code).Delete(&entity.Coupon{}).Error
	if err != nil {
		return errors.New("coudnt delete")
	}
	return nil
}

func (ar *ProductRepository) GetProductsBySearch(offset, limit int, search string) ([]entity.Product, error) {
	var products []entity.Product

	err := ar.db.Select("id, name, price, category, image_url, size").Where("name LIKE ?", search+"%").Offset(offset).Limit(limit).Find(&products).Error
	if err != nil {
		return nil, errors.New("record not found")
	}
	return products, nil
}

func (ar *ProductRepository) GetProductsByCategory(offset, limit, id int) ([]entity.Product, error) {
	var product []entity.Product

	err := ar.db.Where("category=? AND removed =?", id, false).Offset(offset).Limit(limit).Find(&product).Error
	if err != nil {
		return nil, errors.New("record not found")
	}
	return product, nil
}

func (ar *ProductRepository) GetProductsByFilter(minPrize, maxPrize, category int, size string) ([]entity.Product, error) {
	var products []entity.Product

	query := ar.db

	if size != "" {
		query = query.Where("size=?", size)
	}
	if minPrize > 0 {
		query = query.Where("price >= ?", minPrize)
	}
	if maxPrize > 0 {
		query = query.Where("price <= ?", maxPrize)
	}
	if category > 0 {
		query = query.Where("category = ?", category)
	}
	query = query.Where("removed= ?", false)
	err := query.Find(&products).Error
	if err != nil {
		return nil, err
	}
	return products, nil
}

func (pr *ProductRepository) GetAllOffers() ([]entity.Offer, error) {
	var offer []entity.Offer
	currenttime := time.Now()
	err := pr.db.Where("valid_until > ?", currenttime).Find(&offer).Error
	if err != nil {
		return nil, errors.New("record not found")
	}
	return offer, nil

}

func (ar *ProductRepository) GetProductsByCategoryoffer(id int) ([]entity.Product, error) {
	var product []entity.Product

	err := ar.db.Where("category=? AND removed =?", id, false).Find(&product).Error
	if err != nil {
		return nil, errors.New("record not found")
	}
	return product, nil
}
