package middleware

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

var secretkey = []byte("123455666")

func UserRetreiveCookie(c *gin.Context) {
	valid := ValidToken(c)

	if valid == false {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not loged in"})
		c.Abort()
	} else {
		userId, phone, role, err := RetreiveToken(c)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "occured while rereiving"})
			c.Abort()
		} else {
			c.Set("userId", userId)
			c.Set("phonenumber", phone)
		}
		if role != "user" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Not a user"})
			c.Abort()
		} else {
			c.Next()
		}
	}

}

func AdminRetreiveToken(c *gin.Context) {
	valid := ValidToken(c)
	if valid == false {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "admin dont have a token"})
		c.Abort()
	} else {
		Userid, Phone, role, err := RetreiveToken(c)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unable to retreive token"})
			c.Abort()
		} else {
			c.Set("UserId", Userid)
			c.Set("Phonenumber", Phone)
		}
		if role != "admin" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "not admin"})
			c.Abort()
		} else {
			c.Next()
		}
	}
}

func CreateToken(userId int, useremail string, role string, c *gin.Context) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userId": userId,
		"email":  useremail,
		"role":   role,
	})
	tokenstring, err := token.SignedString([]byte("12345678"))

	if err == nil {
		fmt.Println("token created")
	}

	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("Authorise", tokenstring, 3600, "", "", false, true)
}

func ValidToken(c *gin.Context) bool {
	cookie, _ := c.Cookie("Authorise")
	if cookie == "" {
		fmt.Println("cookie not found")
		return false
	} else {
		return true
	}
}
func RetreiveToken(c *gin.Context) (int, int, string, error) {
	cookie, _ := c.Cookie("Authorise")
	if cookie == "" {
		return 0, 0, "", errors.New("cookie not found")
	} else {
		token, err := jwt.Parse(cookie, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte("12345678"), nil
		})
		if err != nil {
			return 0, 0, "", err
		}
		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			var userId, userPhone int
			var role string
		
			if id, exists := claims["userid"]; exists {
				userId, _ = id.(int)
			}
		
			if phone, exists := claims["phone"]; exists {
				userPhone, _ = phone.(int)
			}
		
			if r, exists := claims["role"]; exists {
				role, _ = r.(string)
			}
		
			return userId, userPhone, role, nil
		} else {
			return 0, 0, "", fmt.Errorf("invalid token")
		}
		
		// if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		// 	userId := claims["userid"].(int)
		// 	userPhone := claims["phone"].(int)
		// 	role := claims["role"].(string)
		// 	return userId, userPhone, role, nil
		// } else {
		// 	return 0, 0, "", fmt.Errorf("invalid token")
		// }
	}
}
func DeleteToken(c *gin.Context) error {
	c.SetCookie("Authorise", "", 0, "", "", true, true)
	fmt.Println("cookie deleted")
	return nil
}
