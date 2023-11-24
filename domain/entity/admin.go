package entity

import "gorm.io/gorm"

type Admin struct {
	gorm.Model `json:"-"`
	AdminName  string `json:"adminname"`
	Email      string `json:"email"`
	Phone      string `json:"phone"`
	Password   string `json:"password"`
	Role       string `json:"role"`
	IsActive   bool   `json:"isactive"`
}

// type AdminDashboard struct {
// 	TotalUsers        int `json:"totalusers"`
// 	NewUsers          int `json:"newusers"`
// 	Totalproducts     int `json:"totalproducts"`
// 	StockLessProducts int `json:"stocklessproducts"`
// 	TotalOrders       int `json:"totalorders"`
// }
type AdminDashboard struct {
	TotalUsers        int    `json:"totalusers"`
	NewUsers          int    `json:"newusers"`
	TotalProducts     int    `json:"totalproducts"`
	StocklessProducts int `json:"stocklessrproducts"`
	TotalOrders       int    `json:"totalorders"`
	AverageOrderValue int    `json:"averageordervalue"`
	PendingOrders     int    `json:"pendingorders"`
	ReturnOrders      int    `json:"returnorders"`
	TotalRevenue      int    `json:"totalrevenue"`
	
}