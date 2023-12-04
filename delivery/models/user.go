package models

type EditUser struct {
	Name  string `json:"name" binding:"required" `
	Email string `json:"email" binding:"required"`
}

type Signup struct {
	Name     string `json:"name" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Phone    string `json:"phone" validate:"required,numeric,len=10"`
	Password string `json:"password" validate:"required,min=8"`
	ReferalCode string `json:"referalcode"`

}
type CombinedOrderDetails struct {
	OrderId       string  `json:"order_id"`
	Amount        float64 `json:"amount"`
	OrderStatus   string  `json:"order_status"`
	PaymentStatus bool    `json:"payment_status"`
	Name          string  `json:"name"`
	Email         string  `json:"email"`
	Phone         string  `json:"phone"`
	State         string  `json:"state" validate:"required"`
	Pin           string  `json:"pin" validate:"required"`
	Street        string  `json:"street"`
	City          string  `json:"city"`
	Address       string  `json:"address"`
}
