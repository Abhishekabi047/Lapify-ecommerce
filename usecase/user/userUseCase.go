package usecase

import (
	"errors"
	"fmt"
	"regexp"

	"project/delivery/models"
	"project/domain/entity"
	"project/domain/utils"
	repository "project/repository/user"

	"github.com/go-playground/validator/v10"
	"golang.org/x/crypto/bcrypt"
)

type UserUseCase struct {
	userRepo *repository.UserRepository
}

func NewUser(userRepo *repository.UserRepository) *UserUseCase {
	return &UserUseCase{userRepo: userRepo}
}

func (us *UserUseCase) ExecuteSignup(user entity.User) (*entity.User, error) {
	validate := validator.New()

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
	if err := validate.Struct(newUser); err != nil {

		errors := err.(validator.ValidationErrors)
		for _, e := range errors {
			fmt.Printf("Field: %s, Tag: %s, Value: %s\n", e.Field(), e.Tag(), e.Param())
		}
		return nil, err
	}

	if !isValidName(newUser.Name) {
		return nil, errors.New("Invalid name format")
	}
	err1 := us.userRepo.Create(newUser)
	if err1 != nil {
		return nil, errors.New("Error creating user")
	}
	return newUser, nil
}
func isValidName(name string) bool {
	alphaRegex := regexp.MustCompile("^[a-zA-Z]+$")
	return alphaRegex.MatchString(name)
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

func (uu *UserUseCase) ExecuteAddAddress(address *entity.UserAddress) error {
	validate := validator.New()
	if err := validate.Struct(address); err != nil {
		if _, ok := err.(*validator.InvalidValidationError); ok {
			return err
		}
		errors := err.(validator.ValidationErrors)
		errorMsg := "Validation failed: "
		for _, e := range errors {
			switch e.Tag() {
			case "required":
				errorMsg += fmt.Sprintf("%s is required; ", e.Field())
			case "alpha":
				errorMsg += fmt.Sprintf("%s should contain only alphabetic characters; ", e.Field())
			default:
				errorMsg += fmt.Sprintf("%s has an invalid value; ", e.Field())
			}
		}
		return fmt.Errorf(errorMsg)
	}

	err := uu.userRepo.CreateAddress(address)
	if err != nil {
		return err
	}
	return nil

}

func (uu *UserUseCase) ExecuteEditProfile(user entity.User, userid int) error {
	// validate := validator.New()
	// if err := validate.Struct(user); err != nil {
	// 	if _, ok := err.(*validator.InvalidValidationError); ok {
	// 		return err
	// 	}
	// 	errors := err.(validator.ValidationErrors)
	// 	errorMsg := "Validation failed: "
	// 	for _, e := range errors {
	// 		switch e.Tag() {
	// 		case "required":
	// 			errorMsg += fmt.Sprintf("%s is required; ", e.Field())
	// 		case "alpha":
	// 			errorMsg += fmt.Sprintf("%s should contain only alphabetic characters; ", e.Field())
	// 		case "email":
	// 			errorMsg += fmt.Sprintf("%s should be valid email; ", e.Field())
	// 		default:
	// 			errorMsg += fmt.Sprintf("%s has an invalid value; ", e.Field())
	// 		}
	// 	}
	// 	return fmt.Errorf(errorMsg)
	// }
	user.Id = userid
	err := uu.userRepo.Update(&user)
	if err != nil {
		return errors.New("useer updation failed")
	}
	return nil
}

func (uu *UserUseCase) ExecuteShowUserDetails(userid int) (*entity.User, *entity.UserAddress, error) {
	user, err := uu.userRepo.GetById(userid)
	if err != nil {
		return nil, nil, err
	}
	address, err1 := uu.userRepo.GetAddressById(userid)
	if err1 != nil {
		return nil, nil, err1
	}
	if user != nil && address != nil {
		return user, address, nil
	} else {
		return nil, nil, errors.New("user with this id not found")
	}
}

func (uu *UserUseCase) ExecuteChangePassword(userid int) (string, error) {
	var otpkey entity.OtpKey
	user, err := uu.userRepo.GetById(userid)
	if err != nil {
		return "", err
	}
	key, err1 := utils.SendOtp(user.Phone)
	if err1 != nil {
		return "", err1
	} else {
		otpkey.Key = key
		otpkey.Phone = user.Phone
		err := uu.userRepo.CreateOtpKey(otpkey.Key, otpkey.Phone)
		if err != nil {
			return "", nil
		}
		return key, nil
	}
}

func (uu *UserUseCase) ExecuteOtpValidationPassword(password string, otp string, userid int) error {
	user, err := uu.userRepo.GetById(userid)
	if err != nil {
		return err
	}
	err = utils.CheckOtp(user.Phone, otp)
	if err != nil {
		return err
	}
	hashedpassword, err1 := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	user.Password = string(hashedpassword)
	err1 = uu.userRepo.Update(user)
	if err1 != nil {
		return errors.New("password changing failed")

	}
	return nil

}

func (uu *UserUseCase) ExecuteEditAddress(usaddress entity.UserAddress, id int, useraddress string) error {
	validate := validator.New()
	if err := validate.Struct(usaddress); err != nil {
		if _, ok := err.(*validator.InvalidValidationError); ok {
			return err
		}
		errors := err.(validator.ValidationErrors)
		errorMsg := "Validation failed: "
		for _, e := range errors {
			switch e.Tag() {
			case "required":
				errorMsg += fmt.Sprintf("%s is required; ", e.Field())
			case "alpha":
				errorMsg += fmt.Sprintf("%s should contain only alphabetic characters; ", e.Field())
			case "len":
				errorMsg += fmt.Sprintf("%s should have a length of %s; ", e.Field(), e.Param())
			case "numeric":
				errorMsg += fmt.Sprintf("%s should contain only numeric characters; ", e.Field())
			default:
				errorMsg += fmt.Sprintf("%s has an invalid value; ", e.Field())
			}
		}
		return fmt.Errorf(errorMsg)
	}
	exisitingaddress, err := uu.userRepo.GetAddressByType(id, useraddress)
	if err != nil {
		return err
	}

	exisitingaddress.User_id = id
	exisitingaddress.Address = usaddress.Address
	exisitingaddress.State = usaddress.State
	exisitingaddress.Country = usaddress.Country
	exisitingaddress.Pin = usaddress.Pin
	exisitingaddress.Type = usaddress.Type

	err1 := uu.userRepo.UpdateAddress(exisitingaddress)
	if err1 != nil {
		return err1
	}
	return nil
}

func (uu *UserUseCase) ExecuteDeleteAddress(id int, addtype string) error {
	address, err := uu.userRepo.GetAddressByType(id, addtype)
	if err != nil {
		return err
	}
	err1 := uu.userRepo.DeleteAddress(address.Id)
	if err1 != nil {
		return err1
	}
	return nil
}
