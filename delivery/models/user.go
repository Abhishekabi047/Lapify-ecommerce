package models

type EditUser struct{
	Name string `json:"firstname" binding:"required" `
	Email string `json:"email" binding:"required"`
}

type Signup struct{
	Name string `json:"name"`
	Email string `json:"email"`
	Phone string `json:"phone"`
	Password string `json:"password"`
}