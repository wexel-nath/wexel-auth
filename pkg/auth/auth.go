package auth

import (
	"crypto/rsa"
	"errors"
	"io/ioutil"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/wexel-nath/wexel-auth/pkg/config"
	"github.com/wexel-nath/wexel-auth/pkg/logger"
	"github.com/wexel-nath/wexel-auth/pkg/permission"
	"github.com/wexel-nath/wexel-auth/pkg/user"
)

var(
	privateKey *rsa.PrivateKey
	publicKey  *rsa.PublicKey

	ErrExpiredToken = errors.New("expired jwt needs refresh")
)

type claims struct {
	User        user.User                  `json:"user"`
	Permissions permission.UserPermissions `json:"permissions"`

	jwt.StandardClaims
}

func Configure() {
	privateKeyFile, err := ioutil.ReadFile("keys/test.private.pem")
	if err == nil {
		privateKey, err = jwt.ParseRSAPrivateKeyFromPEM(privateKeyFile)
	}
	if err != nil {
		logger.Error(err)
	}

	publicKeyFile, err := ioutil.ReadFile("keys/test.public.pem")
	if err == nil {
		publicKey, err = jwt.ParseRSAPublicKeyFromPEM(publicKeyFile)
	}
	if err != nil {
		logger.Error(err)
	}
}

func SignUser(userModel user.User, permissions permission.UserPermissions) (string, error) {
	timestamp := time.Now().Unix()

	c := claims{
		User:        userModel,
		Permissions: permissions,

		// default claims
		StandardClaims: jwt.StandardClaims{
			Issuer:    config.GetJwtIssuer(),
			IssuedAt:  timestamp,
			ExpiresAt: timestamp + config.GetJwtExpiry(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, c)
	return token.SignedString(privateKey)
}

func publicKeyFunc(token *jwt.Token) (interface{}, error) {
	return publicKey, nil
}

func Verify(tokenString string) (user.User, error) {
	var c claims
	_, err := jwt.ParseWithClaims(tokenString, &c, publicKeyFunc)
	if err != nil {
		v, ok := err.(*jwt.ValidationError)
		if ok && v.Errors == jwt.ValidationErrorExpired {
			err = ErrExpiredToken
		}
	}

	return c.User, err
}
