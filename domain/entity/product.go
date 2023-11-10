package entity

import "gorm.io/gorm"

type Product struct {
	gorm.Model `json:"-"`
	ID         int    `gorm:"primarykey"`
	Name       string `json:"name"`
	Price      int    `json:"price"`
	Size       string `json:"size"`
	Removed    bool   `json:"removed"`
	// Category   string `json:"category"`
	Category int    `gorm:"foreignKey:ID;references:ID"`
	ImageURL string `json:"imageurl"`
}

type ProductDetails struct {
	gorm.Model    `json:"-"`
	ProductID     int    `json:"productid"`
	Description   string `json:"description"`
	Specification string `json:"specification"`
}

type ProductInput struct {
	Product
	ProductDetails
	Inventory
}

type Inventory struct {
	gorm.Model      `json:"-"`
	ProductId       int
	Quantity        int
	ProductCategory int
}
type Category struct {
	gorm.Model  `json:"-"`
	ID          int    `gorm:"primarykey"`
	Name        string `json:"name"`
	Description string `json:"description"`
}
