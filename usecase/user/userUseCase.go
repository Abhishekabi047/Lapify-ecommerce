package usecase

import (
	"errors"
	"fmt"

	"project/delivery/models"
	"project/domain/entity"
	"project/domain/utils"
	repository "project/repository/user"

	"golang.org/x/crypto/bcrypt"
)

type UserUseCase struct {
	userRepo *repository.UserRepository
}

func NewUser(userRepo *repository.UserRepository) *UserUseCase {
	return &UserUseCase{userRepo: userRepo}
}

func (us *UserUseCase) ExecuteSignup(user entity.User) (*entity.User, error) {
	email, err := us.userRepo.GetByEmail(user.Email)
	if err != nil {
		return nil, errors.New("error with server")
	}
	if email != nil {
		return nil, errors.New("user with email already exists")
	}
	phone, err := us.userRepo.GetByPhone(user.Phone)
	if err != nil {
		return nil, errors.New("error eith server")
	}
	if phone != nil {
		return nil, errors.New("user with phoone already exists")
	}
	hashedpassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	newUser := &entity.User{
		Name:     user.Name,
		Email:    user.Email,
		Phone:    user.Phone,
		Password: string(hashedpassword),
	}

	err1 := us.userRepo.Create(newUser)
	if err1 != nil {
		return nil, errors.New("Error creating user")
	}
	return newUser, nil
}

// func (uu *UserUseCase) ExecuteSignupOtp(phone string) (string,error){
// 	result,err:=uu.userRepo.GetByPhone(phone)
// 	if err != nil{
// 		return "",err
// 	}
// 	if result == nil{
// 		return "",errors.New("user with phone not found")
// 	}
// 	key,err1:= utils.SendOtp(phone)
// 	if err1 != nil{
// 		return "",nil
// 	}else{
// 		err=uu.userRepo.CreateOtpKey(key,phone)
// 		if err != nil{
// 			return "",err
// 		}
// 		return key,nil
// 	}

// }
func (uu *UserUseCase) ExecuteSignupWithOtp(user models.Signup) (string, error) {
	var otpKey entity.OtpKey
	email, err := uu.userRepo.GetByEmail(user.Email)
	if err != nil {
		return "", errors.New("error with server")
	}
	if email != nil {
		return "", errors.New("user with this email already exists")
	}
	phone, err := uu.userRepo.GetByPhone(user.Phone)
	if err != nil {
		return "", errors.New("error with server")
	}
	if phone != nil {
		return "", errors.New("user with this phone no already exists")
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	user.Password = string(hashedPassword)
	key, err := utils.SendOtp(user.Phone)
	if err != nil {
		return "", err
	} else {
		err = uu.userRepo.CreateSignup(&user)
		otpKey.Key = key
		otpKey.Phone = user.Phone
		err = uu.userRepo.CreateOtpKey(key, user.Phone)
		if err != nil {
			return "", err
		}
		return key, nil
	}
}

func (uu *UserUseCase) ExecuteSignupOtpValidation(key string, otp string) error {
	result, err := uu.userRepo.GetByKey(key)

	if err != nil {
		return errors.New("error in key")
	}
	fmt.Printf("GetByKey Result: %+v\n", result)
	user, err := uu.userRepo.GetSignupByPhone(result.Phone)
	if err != nil {
		return errors.New("error in phone")
	}
	err = utils.CheckOtp(result.Phone, otp)
	if err != nil {
		return err
	} else {   
		newUser := &entity.User{
			Name:     user.Name,
			Email:    user.Email,
			Phone:    user.Phone,
			Password: user.Password,
		}
		err1 := uu.userRepo.Create(newUser)
		if err1 != nil {
			return errors.New("error while crearting user")
		} else {
			return nil
		}

	}
}

func (uu *UserUseCase) ExecuteLoginWithPassword(phone, password string) (int, error) {

	user, err := uu.userRepo.GetByPhone(phone)
	if err != nil {
		return 0, err
	}
	if user == nil {
		return 0, errors.New("user with this phone not found")
	}

	permission, err := uu.userRepo.CheckPermission(user)
	if err != nil {
		return 0, err
	}
	if permission == false {
		return 0, errors.New("permission denied")
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return 0, errors.New("invalid password")
	} else {
		return user.Id, nil
	}
}

func (u *UserUseCase) ExecuteLogin(phone string) (string, error) {
	var otpKey entity.OtpKey
	result, err := u.userRepo.GetByPhone(phone)
	if err != nil {
		return "", err
	}
	if result == nil {
		return "", errors.New("no user with this phone found")
	}
	permission, err := u.userRepo.CheckPermission(result)
	if err != nil {
		return "", err
	}
	if permission == false {
		return "", errors.New("permission denied")
	}
	key, err := utils.SendOtp(phone)
	if err != nil {
		return "", err
	} else {
		otpKey.Key = key
		otpKey.Phone = phone
		err = u.userRepo.CreateOtpKey(key, phone)
		if err != nil {
			return "", err
		}
		return key, nil
	}

}

func (uu *UserUseCase) ExecuteOtpValidation(key, otp string) (*entity.User, error) {
	result, err := uu.userRepo.GetByKey(key)
	if err != nil {
		return nil, err
	}
	user, err := uu.userRepo.GetByPhone(result.Phone)
	if err != nil {
		return nil, err
	}
	err1 := utils.CheckOtp(result.Phone, otp)
	if err1 != nil {
		return nil, err
	}
	return user, nil
}
