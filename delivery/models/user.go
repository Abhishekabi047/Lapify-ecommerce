package models

type EditUser struct {
	Name  string `json:"name" binding:"required" `
	Email string `json:"email" binding:"required"`
}

type Signup struct {
	Name     string `json:"name" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Phone    string `json:"phone" validate:"required,phone"`
	Password string `json:"password" validate:"required,min=8"`
}
