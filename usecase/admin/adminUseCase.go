package usecase

import (
	"errors"
	repository "project/repository/admin"

	"project/domain/entity"
	"project/domain/utils"
)

type AdminUseCase struct {
	adminRepo *repository.AdminRepository
}

func NewAdmin(adminRepo *repository.AdminRepository) *AdminUseCase {
	return &AdminUseCase{adminRepo: adminRepo}
}

func (au *AdminUseCase) ExecuteAdminLoginWithPassword(Email, Password string) (int, error) {
	admin, err := au.adminRepo.GetByEmail(Email)
	if err != nil {
		return 0, err
	}
	if admin == nil {
		return 0, errors.New("User doesnt exists")
	}
	if admin.Password != Password {
		return 0, errors.New("Invalid password")
	} else {
		return int(admin.ID), nil
	}
}

func (au *AdminUseCase) ExecutAdminLogin(Phone string) error {
	result, err := au.adminRepo.GetByPhone(Phone)
	if err != nil {
		return err
	}
	if result == nil {
		return errors.New("admin with this phone doesnt exist")
	}
	key, err1 := utils.SendOtp(Phone)
	if err1 != nil {
		return err
	} else {
		err = au.adminRepo.CreateOtpKey(key, Phone)
		if err != nil {
			return nil
		}
		return nil
	}
}

func (au *AdminUseCase) ExecuteOtpValidation(Phone, otp string) (*entity.Admin, error) {
	result, err := au.adminRepo.GetByPhone(Phone)
	if err != nil {
		return nil, err
	}
	err1 := utils.CheckOtp(Phone, otp)
	if err1 != nil {
		return nil, err
	}
	return result, nil
}

func (au *AdminUseCase) ExecuteUsersList(page, limit int) ([]entity.User, error) {
	offset := (page - 1) * limit
	userlist, err := au.adminRepo.GetAllUsers(offset, limit)
	if err != nil {
		return nil, err
	}
	return userlist, nil
}
func (tp *AdminUseCase) ExecuteTogglePermission(id int) error {
	result, err := tp.adminRepo.GetById(id)
	if err != nil {
		return err
	}
	result.Permission = !result.Permission
	err1 := tp.adminRepo.Update(result)
	if err1 != nil {
		return errors.New("user permission toggle failed")
	}
	return nil
}

func (ad *AdminUseCase) ExecuteAdminDashBoard()(*entity.AdminDashboard,error) {
	totalusers,newusers,err:=ad.adminRepo.GetUsers()
	if err != nil{
		return nil,errors.New("error fetching user count")
	}

	totalProducts,stocklessProducts,err:=ad.adminRepo.GetProducts()
	if err != nil{
		return nil,errors.New("error fetching products")
	}

	pendingOrders,returnedOrders,err:=ad.adminRepo.GetOrderByStatus()
	if err != nil{
		return nil,errors.New("error fetching order")
	}

	totalOrders,averageOrdervalue,err:=ad.adminRepo.GetOrders()
	if err != nil {
		return nil,errors.New("error fetching order ")
	}

	totalrevenue,err:=ad.adminRepo.GetRevenue()
	if err != nil{
		return nil,errors.New("error fetching revenue")
	}

	dashboardResponse := entity.AdminDashboard{
		TotalUsers: totalusers,
		NewUsers: newusers,
		TotalProducts: totalProducts,
		StocklessProducts: stocklessProducts,
		TotalOrders: totalOrders,
		AverageOrderValue: averageOrdervalue,
		PendingOrders: pendingOrders,
		ReturnOrders: returnedOrders,
		TotalRevenue: totalrevenue,
	}
	return &dashboardResponse,nil
}