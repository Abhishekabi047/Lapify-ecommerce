package entity

import "gorm.io/gorm"

type Cart struct{
	gorm.Model `json:"-"`
	UserId int `json:"userid"`
	ProductQuantity int `json:"productquantity"`
	TotalPrize int `json:"totalprize"`
	OfferPrize int `json:"offerprize"`
}

type CartItem struct{
	gorm.Model `json:"-"`
	CartId int `json:"cartid"`
	Category int `json:"category"`
	ProductId int  `json:"productid"`
	ProductName string `json:"productname"`
	Quantity int `json:"quantity"`
	Price int `json:"prize"`
}

type WishList struct{
	gorm.Model `json:"-"`
	UserId int `json:"userid"`
	Category int `json:"category"`
	ProductId int `json:"productid"`
	ProductName string `json:"productname"`
	Prize int `json:"prize"`
}