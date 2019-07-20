package jwt

import (
	"net/http"

	"github.com/wexel-nath/authrouter"
	"github.com/wexel-nath/wexel-auth/pkg/config"
	"github.com/wexel-nath/wexel-auth/pkg/logger"
)

var (
	signer        *authrouter.Signer
	authenticator *authrouter.Authenticator

	ErrExpiredToken = authrouter.ErrExpiredToken
)

func Configure() {
	var err error
	signer, err = authrouter.NewSigner(
		config.GetJwtIssuer(),
		config.GetJwtExpiry(),
		config.GetPrivateKeyPath(),
	)
	if err != nil {
		logger.Error(err)
	}

	authenticator, err = authrouter.NewAuthenticator(config.GetPublicKeyPath())
	if err != nil {
		logger.Error(err)
	}
}

func Sign(user authrouter.User) (string, error) {
	return signer.Sign(authrouter.User(user))
}

func Authenticate(r *http.Request) (authrouter.User, error) {
	user, err := authenticator.Authenticate(r)
	return user, err
}
