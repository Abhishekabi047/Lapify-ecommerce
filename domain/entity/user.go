package entity

import "gorm.io/gorm"


type User struct{
	gorm.Model `json:"-"`
	Id  int `gorm:"primarykey" bson:"_id,omitempty" json:"-"`
	Name string `json:"name" bson:"name" binding:"required"`
	Email string `json:"email" bson:"email" binding:"required"`
	Phone string `jsom:"phone" bson:"phone" binding:"required"`
	Password string `json:"password" bson:"password" binding:"required"`
	IsBlocked bool `gorm:"not null;default:true" json:"-"`
	Permission bool   `gorm:"not null;default:true" json:"-"`
}

type UserAddress struct{
	gorm.Model `json:"-"`
	Id int `gorm:"primarykey" json:"id"`
	User_id int `json:"-"`
	Address string `json:"address"`
	State string `json:"state"`
	Country string `json:"country"`
	Pin int `json:"pin"`
	Contact_number int `json:"contact"`

}

type Login struct{
	Email string `json:"email" bson:"email" binding:"required"`
	Password string `json:"pasword" bson:"password" binding:"required"`
}

type OtpKey struct{
	gorm.Model
	Key string `json:"key"`
	Phone string `json:"phone"`
}