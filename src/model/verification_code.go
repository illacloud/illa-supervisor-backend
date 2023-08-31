package model

import (
	"errors"
	"fmt"
	"math/rand"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/illacloud/illa-supervisor-backend/src/utils/config"
)

type VCodeClaims struct {
	Email string `json:"email"`
	Code  string `json:"code"`
	Usage string `json:"usage"`
	jwt.RegisteredClaims
}

func GenerateAndSendVerificationCode(email, usage string) (string, error) {
	// generate random code
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	vCode := fmt.Sprintf("%06v", rnd.Int31n(1000000))

	// send
	if err := SendVerificationEmail(email, vCode, usage); err != nil {
		return "", err
	}

	claims := &VCodeClaims{
		Email: email,
		Code:  vCode,
		Usage: usage,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer: "ILLA",
			ExpiresAt: &jwt.NumericDate{
				Time: time.Now().Add(time.Minute * 15),
			},
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	conf := config.GetInstance()
	codeToken, err := token.SignedString([]byte(conf.GetSecretKey()))
	if err != nil {
		return "", err
	}

	return codeToken, nil
}

func ValidateVerificationCode(vCode, codeToken, email, usage string) (bool, error) {
	vCodeClaims := &VCodeClaims{}
	token, err := jwt.ParseWithClaims(codeToken, vCodeClaims, func(token *jwt.Token) (interface{}, error) {
		conf := config.GetInstance()
		return []byte(conf.GetSecretKey()), nil
	})
	if err != nil {
		return false, err
	}

	claims, ok := token.Claims.(*VCodeClaims)
	if !(ok && claims.Usage == usage) {
		return false, errors.New("invalid verification token")
	}
	if !(claims.Code == vCode && claims.Email == email) {
		return false, errors.New("verification code wrong")
	}
	return true, nil
}
