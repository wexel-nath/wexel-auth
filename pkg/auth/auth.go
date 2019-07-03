package auth

import (
	"crypto/rsa"
	"fmt"
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

type claims struct {
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

	c := claims{
		User: userModel,
		// user permissions

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
	token, err := jwt.ParseWithClaims(tokenString, &claims{}, publicKeyFunc)
	if err != nil {
		return user.User{}, err
	}

	c, ok := token.Claims.(*claims)
	if !ok {
		return user.User{}, fmt.Errorf("error casting token.Claims[%v] to Claims", token)
	}

	logger.Info("claims: %v", c)
	return c.User, nil
}
