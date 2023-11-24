package repository

import (
	"errors"
	"log"
	"project/domain/entity"
	"time"

	"gorm.io/gorm"
)

type AdminRepository struct {
	db *gorm.DB
}

func NewAdminRepository(db *gorm.DB) *AdminRepository {
	return &AdminRepository{db}
}

func (ac *AdminRepository) Create(admin *entity.Admin) error {
	return ac.db.Create(admin).Error
}

func (ar *AdminRepository) GetByPhone(phone string) (*entity.Admin, error) {
	var admin entity.Admin
	result := ar.db.Where(&entity.Admin{Phone: phone}).First(&admin)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, result.Error
	}
	return &admin, nil
}

func (ar *AdminRepository) GetByEmail(email string) (*entity.Admin, error) {
	var admin entity.Admin
	result := ar.db.Where(&entity.Admin{Email: email}).First(&admin)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, result.Error
	}
	return &admin, nil
}

func (ar *AdminRepository) GetById(id int) (*entity.User, error) {
	var user entity.User
	result := ar.db.Where(&entity.User{Id: id}).First(&user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, result.Error
	}
	return &user, nil
}

func (ar *AdminRepository) Update(user *entity.User) error {
	return ar.db.Save(user).Error
}

func (ar *AdminRepository) CreateOtpKey(key, phone string) error {
	var otpkey entity.OtpKey
	otpkey.Key = key
	otpkey.Phone = phone

	if err := ar.db.Create(&otpkey).Error; err != nil {
		log.Printf("Error creating OtpKey: %v", err)
		return err
	}
	return nil
}
func (ar *AdminRepository) GetAllUsers(offset, limit int) ([]entity.User, error) {
	var users []entity.User
	err := ar.db.Offset(offset).Limit(limit).Find(&users).Error
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (ar *AdminRepository) GetUsers() (int, int, error) {
	var totalUsers, newUsers int64
	if err := ar.db.Model(&entity.User{}).Count(&totalUsers).Error; err != nil {
		return 0, 0, err
	}
	if err := ar.db.Model(&entity.User{}).Where("created_at >= ?", time.Now().AddDate(0, 0, -7)).Count(&newUsers).Error; err != nil {
		return 0, 0, err
	}
	return int(totalUsers), int(newUsers), nil
}

func (ar *AdminRepository) GetProducts() (int,int,error) {
	var totalproducts int64
	var stocklessProducts int64
	if err:=ar.db.Model(&entity.Product{}).Where("removed =?",false).Count(&totalproducts).Error;err != nil{
		return 0,0,err
	}
	if err := ar.db.Where(&entity.Inventory{}).Where("quantity=?",0).Count(&stocklessProducts).Error;err != nil{
		return 0,0,err
	}
	return int(totalproducts),int(stocklessProducts),nil
}

func (ar *AdminRepository) GetOrders() (int,int,error){
	var totalorders int64
	var totalamount int64

	if err:=ar.db.Model(&entity.Order{}).Count(&totalorders).Error;err != nil{
		return 0,0,err
	}
	if err :=ar.db.Model(&entity.Order{}).Select("AVG(total)").Row().Scan(&totalamount);err != nil{
		return 0,0,err
	}
	return int(totalorders),int(totalamount),nil
}
func (ar *AdminRepository) GetOrderByStatus() (int,int,error){
	var pendingorder,returnedorder int64
	if err:= ar.db.Model(&entity.Order{}).Where("status =?","pending").Count(&pendingorder).Error;err != nil{
		return 0,0,err
	}
	if err := ar.db.Model(&entity.Order{}).Where("status= ?","return").Count(&returnedorder).Error;err !=nil{
		return 0,0,err
	}
	return int(pendingorder),int(returnedorder),nil
}
func (ar *AdminRepository) GetRevenue() (int,error){
	var totalrevenue int64

	if err:=ar.db.Model(&entity.Order{}).Where("SUM(total)").Row().Scan(&totalrevenue);err !=nil{
		return 0,err
	}
	return int(totalrevenue),nil
}