package authenticator

import (
	"errors"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"

	"github.com/illacloud/illa-supervisior-backend/src/model"
)

type AuthClaims struct {
	User   int       `json:"user"`
	UUID   uuid.UUID `json:"uuid"`
	Random string    `json:"rnd"`
	jwt.RegisteredClaims
}

type Authenticator struct {
	Storage *model.Storage
}

func NewAuthenticator(storage *model.Storage) *Authenticator {
	a := &Authenticator{}
	a.Storage = storage
	return a
}

func (a *Authenticator) ValidateAccessToken(accessToken string) (bool, error) {
	_, _, err := ExtractUserIDFromToken(accessToken)
	if err != nil {
		return false, err
	}
	return true, nil
}

func ExtractUserIDFromToken(accessToken string) (int, uuid.UUID, error) {
	authClaims := &AuthClaims{}
	token, err := jwt.ParseWithClaims(accessToken, authClaims, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("ILLA_SECRET_KEY")), nil
	})
	if err != nil {
		return 0, uuid.Nil, err
	}

	claims, ok := token.Claims.(*AuthClaims)
	if !(ok && token.Valid) {
		return 0, uuid.Nil, err
	}

	return claims.User, claims.UUID, nil
}

func (a *Authenticator) ValidateUser(id int, uid uuid.UUID) (bool, error) {
	// query datasource
	userRecord, err := a.Storage.UserStorage.RetrieveByIDAndUID(id, uid)
	if err != nil {
		return false, err
	}
	if userRecord.ID != id || userRecord.UID != uid {
		return false, errors.New("no such user")
	}

	return true, nil
}

func CreateAccessToken(id int, uid uuid.UUID) (string, error) {
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	vCode := fmt.Sprintf("%06v", rnd.Int31n(10000))

	claims := &AuthClaims{
		User:   id,
		UUID:   uid,
		Random: vCode,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer: "ILLA",
			ExpiresAt: &jwt.NumericDate{
				Time: time.Now().Add(time.Hour * 24 * 7),
			},
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	accessToken, err := token.SignedString([]byte(os.Getenv("ILLA_SECRET_KEY")))
	if err != nil {
		return "", err
	}
	return accessToken, nil
}

func (a *Authenticator) JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		accessToken := c.Request.Header["Authorization"]
		fmt.Printf("accessToken: %v\n", accessToken)
		var token string
		if len(accessToken) != 1 {
			c.AbortWithStatus(http.StatusUnauthorized)
		} else {
			token = accessToken[0]
		}
		userID, userUID, extractErr := ExtractUserIDFromToken(token)
		fmt.Printf("userID: %v, userUID:%v\n", userID, userUID)
		validAccessToken, validaAccessErr := a.ValidateAccessToken(token)
		fmt.Printf("validaAccessErr: %v\n", validaAccessErr)
		validUser, validUserErr := a.ValidateUser(userID, userUID)
		fmt.Printf("validUserErr: %v\n", validUserErr)

		if validAccessToken && validUser && validaAccessErr == nil && extractErr == nil && validUserErr == nil {
			c.Set("userID", userID)
		} else {
			c.AbortWithStatus(http.StatusUnauthorized)
		}
		c.Next()
	}
}

func (a *Authenticator) ManualAuth(accessToken string) (bool, error) {
	userID, userUID, extractErr := ExtractUserIDFromToken(accessToken)
	validAccessToken, validaAccessErr := a.ValidateAccessToken(accessToken)
	validUser, validUserErr := a.ValidateUser(userID, userUID)

	if validAccessToken && validUser && validaAccessErr == nil && extractErr == nil && validUserErr == nil {
		return true, nil
	} else {
		return false, errors.New("auth failed.")
	}
}
