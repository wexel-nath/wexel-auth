package auth

import (
	"crypto/rsa"
	"io/ioutil"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/wexel-nath/wexel-auth/pkg/config"
	"github.com/wexel-nath/wexel-auth/pkg/logger"
	"github.com/wexel-nath/wexel-auth/pkg/user"
)

var(
	privateKey *rsa.PrivateKey
	publicKey  *rsa.PublicKey
)

type Claims struct {
	User user.User `json:"user"`
	// todo: add user permissions

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

func SignUser(userModel user.User) (string, error) {
	timestamp := time.Now().Unix()

	claims := Claims{
		User: userModel,
		// user permissions

		StandardClaims: jwt.StandardClaims{
			Issuer:    "wexel-auth.herokuapp.com",
			IssuedAt:  timestamp,
			ExpiresAt: timestamp + config.GetJwtExpiry(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	return token.SignedString(privateKey)
}

func publicKeyFunc(token *jwt.Token) (interface{}, error) {
	return publicKey, nil
}

func Verify(tokenString string) error {
	claims := Claims{}

	token, err := jwt.ParseWithClaims(tokenString, &claims, publicKeyFunc)
	if err != nil {
		return err
	}

	logger.Info("claims: %v", token.Claims)
	return nil
}
