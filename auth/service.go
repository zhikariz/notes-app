package auth

import (
	"errors"
	"fmt"
	. "notes-app/entity"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type Service interface {
	GenerateToken(account Account) (string, error)
	GenerateVerifyEmailToken(account Account) (string, error)
	//GenerateResetPasswordToken(account Account) (string, error)
	ValidateToken(encodedToken string) (*jwt.Token, error)
}

type jwtService struct{}

func NewService() *jwtService {
	return &jwtService{}
}

func (s *jwtService) GenerateToken(account Account) (string, error) {
	claim := jwt.MapClaims{}
	claim["account_id"] = account.ID
	claim["full_name"] = account.FullName
	claim["email"] = account.Email
	claim["role"] = account.Role
	claim["exp"] = time.Now().Add(time.Minute * 15).Unix()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)

	signedToken, err := token.SignedString([]byte(os.Getenv("SECRET_KEY")))

	if err != nil {
		return signedToken, err
	}
	return signedToken, nil
}

func (s *jwtService) GenerateVerifyEmailToken(account Account) (string, error) {
	claim := jwt.MapClaims{}
	claim["email"] = account.Email
	claim["exp"] = time.Now().Add(time.Minute * 2).Unix()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)

	signedToken, err := token.SignedString([]byte(os.Getenv("SECRET_KEY")))

	if err != nil {
		return signedToken, err
	}
	return signedToken, nil
}

func (s *jwtService) ValidateToken(encodedToken string) (*jwt.Token, error) {
	tokenString, err := jwt.Parse(encodedToken, func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)

		if !ok {
			return nil, errors.New("invalid token")
		}

		return []byte(os.Getenv("SECRET_KEY")), nil
	})

	if err != nil {
		return tokenString, err
	}

	claims := tokenString.Claims.(jwt.MapClaims)
	var currentTimeClaims int64 = int64(claims["exp"].(float64))
	timeNow := time.Now().Unix() - int64(time.Minute)*15
	if (currentTimeClaims) > timeNow {
		fmt.Println(currentTimeClaims)
		fmt.Println(timeNow)
	}

	return tokenString, nil

}
