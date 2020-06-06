package jwt

import (
	"net/http"

	"wexel-auth/pkg/config"
	"wexel-auth/pkg/logger"

	"github.com/wexel-nath/auth"
)

var (
	signer        *auth.Signer
	authenticator *auth.Authenticator

	ErrExpiredToken = auth.ErrExpiredToken
)

func Configure() {
	var err error
	signer, err = auth.NewSigner(
		config.GetJwtIssuer(),
		config.GetJwtExpiry(),
		config.GetPrivateKeyPath(),
	)
	if err != nil {
		logger.Error(err)
	}

	authenticator, err = auth.NewAuthenticator(config.GetPublicKeyPath())
	if err != nil {
		logger.Error(err)
	}
}

func Sign(user auth.User) (string, error) {
	return signer.Sign(user)
}

func Authenticate(r *http.Request) (auth.User, error) {
	return authenticator.Authenticate(r)
}

func Authorize(r *http.Request, capability string) (auth.User, error) {
	user, err := Authenticate(r)
	if err != nil {
		return user, err
	}

	return user, user.IsAuthorized(config.GetServiceName(), []string{capability})
}
