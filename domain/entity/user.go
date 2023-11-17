package entity

import "gorm.io/gorm"

type User struct {
	gorm.Model `json:"-"`
	Id         int    `gorm:"primarykey"  json:"-"`
	Name       string `json:"name" validate:"required,alpha" `
	Email      string `json:"email" validate:"required,email"`
	Phone      string `json:"phone" validate:"required"`
	Password   string `json:"password" validate:"required,min=8"`
	IsBlocked  bool   `gorm:"not null;default:true" json:"-"`
	Permission bool   `gorm:"not null;default:true" json:"-"`
}

type UserAddress struct {
	gorm.Model `json:"-"`
	Id         int    `gorm:"primarykey" json:"id"`
	User_id    int    `json:"-"`
	Address    string `json:"address" validate:"required"`
	State      string `json:"state" validate:"required,alpha"`
	Country    string `json:"country" validate:"required,alpha"`
	Pin        string `json:"pin" validate:"required,numeric,len=6"`
	Type       string `json:"type" validate:"required,alpha"`
}

type Login struct {
	Email    string `json:"email" bson:"email" binding:"required"`
	Password string `json:"pasword" bson:"password" binding:"required"`
}

type OtpKey struct {
	gorm.Model
	Key   string `json:"key"`
	Phone string `json:"phone"`
}
