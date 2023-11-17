package repository

import (
	"project/domain/entity"
)

type UserInterface interface {
	GetByID(id uint) (*entity.User, error)
	GetByEmail(email string) (*entity.User, error)
	Create(user *entity.User) error
	Update(user *entity.User) error
	Delete(id uint) error
}
