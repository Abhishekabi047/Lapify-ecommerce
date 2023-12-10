package entity

import (
	"time"

	"gorm.io/gorm"
)

type Product struct {
	gorm.Model `json:"-"`
	ID         int    `gorm:"primarykey"`
	Name       string `json:"name" validate:"required" form:"name"`
	Price      int    `json:"price" validate:"required,number" form:"price"`
	OfferPrize int    `json:"offerprice"  `
	Size       string `json:"size" validate:"required" form:"size"`
	Removed    bool   `json:"removed"`
	Category   int    `form:"category" gorm:"foreignKey:ID;references:ID" validate:"required,numeric"`
	ImageURL   string `json:"imageurl" `
}

type ProductDetails struct {
	gorm.Model    `json:"-"`
	ProductID     int    `json:"productid"`
	Description   string `json:"description" validate:"required" form:"description"`
	Specification string `json:"specification" validate:"required" form:"specification"`
}

type ProductInput struct {
	Product
	ProductDetails
	Inventory
}

type Inventory struct {
	gorm.Model      `json:"-"`
	ProductId       int
	Quantity        int `validate:"required,numeric" form:"quantity"`
	ProductCategory int
}
type Category struct {
	gorm.Model  `json:"-"`
	ID          int    `gorm:"primarykey"`
	Name        string `json:"name" validate:"required,alpha"`
	Description string `json:"description" validate:"required"`
}

type Coupon struct {
	gorm.Model `json:"-"`
	Id         int       `json:"id" `
	Code       string    `json:"code" validate:"required,max=8"   `
	Type       string    `json:"type" validate:"required,alpha"`
	Amount     int       `json:"amount" validate:"required,numeric,positive"`
	ValidFrom  time.Time `json:"-"`
	Validuntil time.Time `json:"valid_until"`
	UsageLimit int       `json:"usage_limit" validate:"required,numeric"`
	UsedCount  int       `json:"usedcount"`
}

type Offer struct {
	gorm.Model `json:"-"`
	Id         int       `json:"-" gorm:"primarykey"`
	Name       string    `json:"name"`
	Type       string    `json:"type"`
	Amount     int       `json:"amount"`
	MinPrice   int       `json:"minprice"`
	ValidFrom  time.Time `json:"-"`
	ValidUntil time.Time `json:"valid_until"`
	UsageLimit int       `json:"usage_limit"`
	UsedCount  int       `json:"-"`
	Category   int       `json:"category"`
	ProductId  int       `json:"product_id"`
}

type UsedCoupon struct {
	UserId     int    `json:"userid"`
	CouponCode string `json:"couponcode"`
}
