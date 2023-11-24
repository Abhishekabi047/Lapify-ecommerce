package utils

import (
	"crypto/hmac"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/hex"
	"errors"
)

func RazorPaymentVerification(sign ,orderId,PaymentId string) error{
	signature:=sign
	secret:="R59k58EhgS48BaauF22urj5A"
	data := orderId + "|" + PaymentId
	h:=hmac.New(sha256.New,[]byte(secret))
	_,err:=h.Write([]byte(data))
	if err != nil{
		panic(err)
	}
	sha:= hex.EncodeToString(h.Sum(nil))
	if subtle.ConstantTimeCompare([]byte(sha),[]byte(signature)) != 1{
		return errors.New("Payment failed")
	}
	return nil
}