package auth

import (
	"crypto/rsa"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
)

var(
	jwtIssuer  string
	jwtExpiry  int64
	publicKey  *rsa.PublicKey
	privateKey *rsa.PrivateKey

	ErrExpiredToken = errors.New("expired jwt needs refresh")
)

type User struct {
	UserID      int64               `json:"user_id"`
	FirstName   string              `json:"first_name"`
	LastName    string              `json:"last_name"`
	Email       string              `json:"email"`
	Username    string              `json:"username"`
	Permissions map[string][]string `json:"permissions"`
}

type claims struct {
	User User `json:"user"`
	jwt.StandardClaims
}

type Config struct {
	JwtIssuer      string
	JwtExpiry      int64
	PublicKeyPath  string
	PrivateKeyPath string
}

func Configure(config Config) error {
	if config.PublicKeyPath == "" {
		return errors.New("public key path must be specified")
	}
	publicKeyFile, err := ioutil.ReadFile(config.PublicKeyPath)
	if err == nil {
		publicKey, err = jwt.ParseRSAPublicKeyFromPEM(publicKeyFile)
	}
	if err != nil {
		return err
	}

	jwtIssuer = config.JwtIssuer
	jwtExpiry = config.JwtExpiry

	if config.PrivateKeyPath == "" {
		// don't attempt to initialize privateKey
		return nil
	}

	privateKeyFile, err := ioutil.ReadFile(config.PrivateKeyPath)
	if err == nil {
		privateKey, err = jwt.ParseRSAPrivateKeyFromPEM(privateKeyFile)
	}
	return err
}

func buildStandardClaims() jwt.StandardClaims {
	timestamp := time.Now().Unix()

	return jwt.StandardClaims{
		Issuer:    jwtIssuer,
		IssuedAt:  timestamp,
		ExpiresAt: timestamp + jwtExpiry,
	}
}

// Sign is used to create a signed token
func Sign(user User) (string, error) {
	if privateKey == nil {
		return "", errors.New("private key has not been initialized")
	}
	if jwtIssuer == "" {
		return "", errors.New("jwt issuer must be specified")
	}
	if jwtIssuer == "" {
		return "", errors.New("jwt expiry must be specified")
	}

	c := claims{
		User:           user,
		StandardClaims: buildStandardClaims(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, c)
	return token.SignedString(privateKey)
}

func GetAuthenticatedUser(r *http.Request) (User, error) {
	token, err := getTokenFromRequest(r)
	if err != nil {
		return User{}, nil
	}

	return authenticate(token)
}

func GetAuthorizedUser(r *http.Request, service string, capability string) (User, error) {
	token, err := getTokenFromRequest(r)
	if err != nil {
		return User{}, err
	}

	return authorize(token, service, capability)
}

func getTokenFromRequest(r *http.Request) (string, error) {
	value := r.Header.Get("Authorization")
	s := strings.Split(value, "Bearer")
	if len(s) != 2 {
		return "", errors.New("authorization header is not of the form: Bearer <token>")
	}

	return strings.TrimSpace(s[1]), nil
}

func publicKeyFunc(_ *jwt.Token) (interface{}, error) {
	return publicKey, nil
}

func authenticate(token string) (User, error) {
	var c claims
	_, err := jwt.ParseWithClaims(token, &c, publicKeyFunc)
	if err != nil {
		v, ok := err.(*jwt.ValidationError)
		if ok && v.Errors == jwt.ValidationErrorExpired {
			err = ErrExpiredToken
		}
	}
	return c.User, err
}

func authorize(token string, service string, capability string) (User, error) {
	user, err := authenticate(token)
	if err != nil {
		return user, err
	}

	servicePermissions, ok := user.Permissions[service]
	if !ok {
		return user, fmt.Errorf("not authorized for service %s", service)
	}

	if capability == "" {
		return user, nil
	}

	allPermissions, ok := user.Permissions["ALL"]
	if !ok {
		allPermissions = []string{}
	}

	servicePermissions = append(servicePermissions, allPermissions...)
	if !inStringSlice(capability, servicePermissions) {
		return user, fmt.Errorf("missing capability %s", capability)
	}
	return user, nil
}

func inStringSlice(s string, ss []string) bool {
	for _, element :=  range ss {
		if element == s {
			return true
		}
	}
	return false
}
