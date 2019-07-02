package auth

import (
	"crypto/rsa"
	"github.com/wexel-nath/wexel-auth/pkg/logger"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/wexel-nath/wexel-auth/pkg/config"
	"github.com/wexel-nath/wexel-auth/pkg/user"
)

var(
	privateKey *rsa.PrivateKey
)

func Configure() {
	// todo: get PEM rsa key
	pem := []byte("")
	key, err := jwt.ParseRSAPrivateKeyFromPEM(pem)
	if err != nil {
		logger.Error(err)
		return
	}

	privateKey = key
}

func SignUser(userModel user.User) (string, error) {
	timestamp := time.Now().Unix()

	claims := struct{
		User user.User `json:"user"`
		// todo: add user permissions

		jwt.StandardClaims
	}{
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
