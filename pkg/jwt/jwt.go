package jwt

import (
	"github.com/wexel-nath/jwt"
	"github.com/wexel-nath/wexel-auth/pkg/config"
	"github.com/wexel-nath/wexel-auth/pkg/logger"
	"net/http"
)

type User jwt.User

var (
	signer        *jwt.Signer
	authenticator *jwt.Authenticator

	ErrExpiredToken = jwt.ErrExpiredToken
)

func Configure() {
	var err error
	signer, err = jwt.NewSigner(
		config.GetJwtIssuer(),
		config.GetJwtExpiry(),
		config.GetPrivateKeyPath(),
	)
	if err != nil {
		logger.Error(err)
	}

	authenticator, err = jwt.NewAuthenticator(config.GetPublicKeyPath())
	if err != nil {
		logger.Error(err)
	}
}

func Sign(user User) (string, error) {
	return signer.Sign(jwt.User(user))
}

func Authenticate(r *http.Request) (User, error) {
	user, err := authenticator.Authenticate(r)
	return user.(User), err
}
