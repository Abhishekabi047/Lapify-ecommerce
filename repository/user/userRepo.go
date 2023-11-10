package repository

import (
	"errors"
	"log"
	"project/delivery/models"
	"project/domain/entity"

	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db}
}

func (ur *UserRepository) GetById(id int) (*entity.User, error) {
	var user entity.User
	result := ur.db.Find(&user, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, result.Error
	}
	return &user, nil
}

func (ur *UserRepository) GetByEmail(email string) (*entity.User, error) {
	var user entity.User
	result := ur.db.Where(&entity.User{Email: email}).First(&user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, result.Error
	}
	return &user, nil
}

func (ur *UserRepository) GetByPhone(phone string) (*entity.User, error) {
	var user entity.User
	result := ur.db.Where(&entity.User{Phone: phone}).First(&user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, result.Error
	}
	return &user, nil
}

func (ur *UserRepository) CheckPermission(user *entity.User) (bool, error) {
	result := ur.db.Where(&entity.User{Phone: user.Phone}).First(&user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return false, nil
		}
		return false, result.Error
	}
	permission := user.Permission
	return permission, nil
}
func (ur *UserRepository) CreateSignup(user *models.Signup) error {
	return ur.db.Create(&user).Error
}

func (ur *UserRepository) GetSignupByPhone(phone string) (*models.Signup, error) {
	var user models.Signup
	result := ur.db.Where(&models.Signup{Phone: phone}).First(&user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, result.Error
	}
	return &user, nil
}

func (ur *UserRepository) Create(user *entity.User) error {
	return ur.db.Create(user).Error
}

func (ur *UserRepository) Update(user *entity.User) error {
	return ur.db.Updates(user).Error
}

func (ur *UserRepository) Delete(user *entity.User) error {
	return ur.db.Delete(user).Error
}
func (ur *UserRepository) CreateOtpKey(key, phone string) error {
	var otpkey entity.OtpKey
	otpkey.Key = key
	otpkey.Phone = phone
	// return ar.db.Create(otpkey).Error
	if err := ur.db.Create(&otpkey).Error; err != nil {
		log.Printf("Error creating OtpKey: %v", err)
		return err
	}
	return nil
}
func (ur *UserRepository) GetByKey(key string) (*entity.OtpKey, error) {
	var otpKey entity.OtpKey
	result := ur.db.Where(&entity.OtpKey{Key: key}).Order("created_at DESC").First(&otpKey)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, result.Error
	}
	return &otpKey, nil
}
