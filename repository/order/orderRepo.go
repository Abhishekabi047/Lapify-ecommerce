package order

import (
	"errors"
	"project/domain/entity"
	"time"

	"gorm.io/gorm"
)

type OrderRepository struct {
	db *gorm.DB
}

func NewOrderRepository(db *gorm.DB) *OrderRepository {
	return &OrderRepository{db}
}

func (or *OrderRepository) Create(order *entity.Order) (int, error) {
	if err := or.db.Create(order).Error; err != nil {
		return 0, err
	}
	return int(order.ID), nil
}

func (or *OrderRepository) GetOrderById(orderid int) (*entity.Order, error) {
	var order entity.Order
	result := or.db.Where("id=?", orderid).First(&order)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("order not found")
		}
		return nil, errors.New("order not found")
	}
	return &order, nil
}

func (or *OrderRepository) Update(order *entity.Order) error {
	return or.db.Save(order).Error
}

func (or *OrderRepository) GetAllOrders(userid, offset, limit int) ([]entity.Order, error) {
	var order []entity.Order
	result := or.db.Offset(offset).Limit(limit).Where("user_id=?", userid).Find(&order)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("order not found")
		}
		return nil, errors.New("order not found")

	}
	return order, nil
}
func (or *OrderRepository) CreateOrderItems(orderitems []entity.OrderItem) error {
	if err := or.db.Create(orderitems).Error; err != nil {
		return err
	}
	return nil
}

func (or *OrderRepository) GetByStatus(offset, limit int, status string) ([]entity.Order, error) {
	var order []entity.Order
	result := or.db.Offset(offset).Limit(limit).Where("status=?", status).Find(&order)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("order not found")
		}
		return nil, errors.New("order not found")
	}
	return order, nil
}
func (or *OrderRepository) CreateInvoice(invoice *entity.Invoice) (*entity.Invoice, error) {
	if err := or.db.Create(invoice).Error; err != nil {
		return nil, errors.New("eroor creating invoice")
	}
	return invoice, nil
}

func (or *OrderRepository) GetAllOrderList(offset, limit int) ([]entity.Order, error) {
	var order []entity.Order
	result := or.db.Offset(offset).Limit(limit).Find(&order)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, result.Error
		}
		return nil, result.Error
	}
	return order, nil
}

func (or *OrderRepository) GetByRazorId(razorId string) (*entity.Order, error) {
	var order entity.Order
	result := or.db.Where("payment_id=?", razorId).First(&order)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("order not found")
		}
	}
	return &order, nil
}

func (or *OrderRepository) GetByDate(startdate, enddate time.Time) (*entity.SalesReport, error) {
	var order []entity.Order
	var report entity.SalesReport

	if err := or.db.Model(&order).Where("created_at BETWEEN ? AND ? AND status =?", startdate, enddate, "confirmed").Select("SUM(total) as total_sales").Scan(&report).Error; err != nil {
		return nil, err
	}
	if err := or.db.Model(&order).Where("created_at BETWEEN ? AND ? AND status =?", startdate, enddate, "confirmed").Count(&report.TotalOrders).Error; err != nil {
		return nil, err
	}
	if err := or.db.Model(&order).Where("created_at BETWEEN ? AND ? AND status =?", startdate, enddate, "confirmed").Select("AVG(total) as average_order").Scan(&report).Error; err != nil {
		return nil, err
	}

	return &report, nil

}
func (or *OrderRepository) GetByPaymentMethod(startdate, enddate time.Time, paymentmethod string) (*entity.SalesReport, error) {
	var order []entity.Order
	var report entity.SalesReport

	if err := or.db.Model(&order).Where("created_at BETWEEN ? AND ? AND status =? AND payment_method=?", startdate, enddate, "confirmed", paymentmethod).Select("SUM(total) as total_sales").Scan(&report).Error; err != nil {
		return nil, err
	}
	if err := or.db.Model(&order).Where("created_at BETWEEN ? AND ? AND status =? AND payment_method=?", startdate, enddate, "confirmed", paymentmethod).Count(&report.TotalOrders).Error; err != nil {
		return nil, err
	}
	if err := or.db.Model(&order).Where("created_at BETWEEN ? AND ? AND status =? AND payment_method=?", startdate, enddate, "confirmed", paymentmethod).Select("AVG(total) as average_order").Scan(&report).Error; err != nil {
		return nil, err
	}

	return &report, nil

}
func (or *OrderRepository) SavePayment(charge *entity.Charge) (err error) {
	if err := or.db.Create(charge).Error; err != nil {
		return err
	}
	return nil
}

func (or *OrderRepository) UpdateInvoice(invoice *entity.Invoice) error {
	return or.db.Save(&invoice).Error
}
