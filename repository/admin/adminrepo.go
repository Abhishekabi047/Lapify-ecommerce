package repository

import (
	"errors"
	"log"
	"project/domain/entity"

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
func (ar *AdminRepository) GetAllUsers(offset,limit int) ([]entity.User,error){
	var users []entity.User
	err:=ar.db.Offset(offset).Limit(limit).Find(&users).Error
	if err != nil{
		return nil,err
	}
	return users,nil
}