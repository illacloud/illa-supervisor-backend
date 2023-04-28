package authenticator

import (
	"errors"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"

	"github.com/illacloud/illa-supervisor-backend/src/model"
)

type AuthClaims struct {
	User   int       `json:"user"`
	UUID   uuid.UUID `json:"uuid"`
	Random string    `json:"rnd"`
	jwt.RegisteredClaims
}

type Authenticator struct {
	Storage *model.Storage
	Cache   *model.Cache
}

func NewAuthenticator(storage *model.Storage, cache *model.Cache) *Authenticator {
	a := &Authenticator{
		Storage: storage,
		Cache:   cache,
	}
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

func ExtractExpiresAtFromToken(accessToken string) (*jwt.NumericDate, error) {
	authClaims := &AuthClaims{}
	token, err := jwt.ParseWithClaims(accessToken, authClaims, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("ILLA_SECRET_KEY")), nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*AuthClaims)
	if !(ok && token.Valid) {
		return nil, err
	}

	return claims.ExpiresAt, nil
}

func (a *Authenticator) ValidateUser(user *model.User, id int, uid uuid.UUID) (bool, error) {
	// refuse invalied user
	emptyUUID, _ := uuid.Parse("00000000-0000-0000-0000-000000000000")
	if id == 0 || uid == emptyUUID {
		return false, errors.New("invalied user ID or UID.")
	}
	if user.ID != id || user.UID != uid {
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
		var token string
		if len(accessToken) != 1 {
			c.AbortWithStatus(http.StatusUnauthorized)
		} else {
			token = accessToken[0]
		}

		// fetch user
		userID, userUID, extractErr := ExtractUserIDFromToken(token)
		user, err := a.Storage.UserStorage.RetrieveByIDAndUID(userID, userUID)
		if err != nil {
			c.AbortWithStatus(http.StatusUnauthorized)
		}

		// validate
		validAccessToken, validaAccessErr := a.ValidateAccessToken(token)
		validUser, validUserErr := a.ValidateUser(user, userID, userUID)
		expireAtAvaliable, errInValidatteExpireAt := a.DoesAccessTokenExpiredAtAvaliable(user, token)

		if validAccessToken && validUser && expireAtAvaliable && validaAccessErr == nil && extractErr == nil && validUserErr == nil && errInValidatteExpireAt == nil {
			c.Set("userID", userID)
		} else {
			c.AbortWithStatus(http.StatusUnauthorized)
		}
		c.Next()
	}
}

func (a *Authenticator) ManualAuth(accessToken string) (bool, error) {
	// fetch user
	userID, userUID, extractErr := ExtractUserIDFromToken(accessToken)
	user, err := a.Storage.UserStorage.RetrieveByIDAndUID(userID, userUID)
	if err != nil {
		return false, errors.New("auth failed.")
	}
	validAccessToken, validaAccessErr := a.ValidateAccessToken(accessToken)
	validUser, validUserErr := a.ValidateUser(user, userID, userUID)
	expireAtAvaliable, errInValidatteExpireAt := a.DoesAccessTokenExpiredAtAvaliable(user, accessToken)

	if validAccessToken && validUser && expireAtAvaliable && validaAccessErr == nil && extractErr == nil && validUserErr == nil && errInValidatteExpireAt == nil {
		return true, nil
	} else {
		return false, errors.New("auth failed.")
	}
}

func ExtractExpiresAtFromTokenInString(accessToken string) (string, error) {
	// extract now token expiresAt
	expireDate, errInExtract := ExtractExpiresAtFromToken(accessToken)
	if errInExtract != nil {
		return "", errInExtract
	}
	expiresAt := strconv.FormatInt(expireDate.UTC().Unix(), 10)
	return expiresAt, nil
}

// for logout case
func (a *Authenticator) DoesAccessTokenExpiredAtAvaliable(user *model.User, accessToken string) (bool, error) {
	// extract now token expiresAt
	expireDate, errInExtract := ExtractExpiresAtFromToken(accessToken)
	if errInExtract != nil {
		return false, errInExtract
	}
	expiresAt := strconv.FormatInt(expireDate.UTC().Unix(), 10)
	// get history data
	return a.Cache.JWTCache.DoesUserJWTTokenAvaliable(user, expiresAt)
}
