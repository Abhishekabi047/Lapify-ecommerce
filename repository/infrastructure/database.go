package infrastructure

import (
	"fmt"
	"project/config"
	"project/delivery/models"
	"project/domain/entity"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB
var dsn string
var Dbco = "host=localhost user=postgres dbname=ecom password=2589 port=5432 sslmode=disable"

func ConnectDb(config config.DataBase) (*gorm.DB, error) {
	psqlInfo := fmt.Sprintf("host=%s user=%s dbname=%s port=%s password=%s", config.DBHost, config.DBUser, config.DBName, config.DBPort, config.DBPassword)
	db, err := gorm.Open(postgres.Open(psqlInfo), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to db : %w", err)
	}
	DB = db
	DB.AutoMigrate(&entity.Admin{}, &entity.OtpKey{}, &models.Signup{}, &entity.User{}, &entity.Category{}, &entity.Product{}, &entity.ProductDetails{}, &entity.ProductInput{}, &entity.Inventory{}, &entity.CartItem{}, &entity.Cart{}, &entity.WishList{}, &entity.Order{}, &entity.OrderItem{}, &entity.UserAddress{}, &entity.Coupon{}, &entity.UsedCoupon{}, &entity.Offer{}, &entity.Inventory{}, &entity.Invoice{})
	return db, nil
}
