package entity

import (
	"time"

	"gorm.io/gorm"
)

type Product struct {
	gorm.Model `json:"-"`
	ID         int    `gorm:"primarykey"`
	Name       string `json:"name" validate:"required"`
	Price      int    `json:"price" validate:"required,numeric,positive"`
	Size       string `json:"size" validate:"required"`
	Removed    bool   `json:"removed"`
	Category   int    `gorm:"foreignKey:ID;references:ID" validate:"required,numeric"`
	ImageURL   string `json:"imageurl" validate:"required"`
}

type ProductDetails struct {
	gorm.Model    `json:"-"`
	ProductID     int    `json:"productid"`
	Description   string `json:"description" validate:"required"`
	Specification string `json:"specification" validate:"required"`
}

type ProductInput struct {
	Product
	ProductDetails
	Inventory
}

type Inventory struct {
	gorm.Model      `json:"-"`
	ProductId       int
	Quantity        int `validate:"required,numeric"`
	ProductCategory int
}
type Category struct {
	gorm.Model  `json:"-"`
	ID          int    `gorm:"primarykey"`
	Name        string `json:"name" validate:"required,alpha"`
	Description string `json:"description" validate:"required"`
}

type Coupon struct {
	Id         int       `json:"id"`
	Code       string    `json:"code"`
	Type       string    `json:"type"`
	Amount     int       `json:"amount"`
	ValidFrom  time.Time `json:"-"`
	Validuntil time.Time `json:"valid_until"`
	UsageLimit int       `json:"usage_limit"`
	UsedCount  int       `json:"usedcount"`
	Category   int       `json:"category"`
	Adminid    int       `json:"-"`
}

type Offer struct {
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
	AdminId    int       `json:"-"`
}

type UsedCoupon struct {
	UserId     int    `json:"userid"`
	CouponCode string `json:"couponcode"`
}
