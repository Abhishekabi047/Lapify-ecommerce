package utils

import (
	"errors"
	"fmt"

	// "log"
	// "os"

	// "github.com/joho/godotenv"
	"github.com/twilio/twilio-go"
	openapi "github.com/twilio/twilio-go/rest/verify/v2"
)

var (
	TWILIO_ACCOUNT_SID string
	TWILIO_AUTH_TOKEN  string
	VERIFY_SERVICE_SID string
	client             *twilio.RestClient
)

// func init() {
// err := godotenv.Load()
// if err != nil {
// 	log.Fatal("Error loading .env file")
// }
// 	TWILIO_ACCOUNT_SID = os.Getenv("KEY1")
// 	TWILIO_AUTH_TOKEN = os.Getenv("KEY2")
// 	VERIFY_SERVICE_SID = os.Getenv("KEY3")
// 	client = twilio.NewRestClientWithParams(twilio.ClientParams{
// 		Username: TWILIO_ACCOUNT_SID,
// 		Password: TWILIO_AUTH_TOKEN,
// 	})

// }

func SendOtp(phone string) (string, error) {

	TWILIO_ACCOUNT_SID = "AC1cd1b8f006dcfa198971917f6c6c0ca6"
	TWILIO_AUTH_TOKEN = "d1d53cf7ba005b917745ddf2b54581cd"
	VERIFY_SERVICE_SID = "VAe34d88e9c26c376010bba3daf29304cb"
	client = twilio.NewRestClientWithParams(twilio.ClientParams{
		Username: TWILIO_ACCOUNT_SID,
		Password: TWILIO_AUTH_TOKEN,
	})

	to := "+91" + phone
	params := &openapi.CreateVerificationParams{}
	params.SetTo(to)
	params.SetChannel("sms")
	resp, err := client.VerifyV2.CreateVerification(VERIFY_SERVICE_SID, params)
	if err != nil {
		fmt.Println(err.Error())
		return "", errors.New("failed to generate otp")
	} else {
		fmt.Printf("verification code send '%s \n'", *resp.Sid)
		return *resp.Sid, nil
	}
}

func CheckOtp(phone, code string) error {

	TWILIO_ACCOUNT_SID = "AC1cd1b8f006dcfa198971917f6c6c0ca6"
	TWILIO_AUTH_TOKEN = "d1d53cf7ba005b917745ddf2b54581cd"
	VERIFY_SERVICE_SID = "VAe34d88e9c26c376010bba3daf29304cb"
	client = twilio.NewRestClientWithParams(twilio.ClientParams{
		Username: TWILIO_ACCOUNT_SID,
		Password: TWILIO_AUTH_TOKEN,
	})
	if code == "" {
		return errors.New("OTP code is empty")
	}
	to := "+91" + phone

	params := &openapi.CreateVerificationCheckParams{}

	params.SetTo(to)
	params.SetCode(code)
	fmt.Print(code)
	resp, err := client.VerifyV2.CreateVerificationCheck(VERIFY_SERVICE_SID, params)
	if err != nil {
		fmt.Println(err.Error())
		return errors.New("invalid otp")
	} else if *resp.Status == "approved" {
		return nil
	} else {
		return errors.New("invalid otp")
	}
}
