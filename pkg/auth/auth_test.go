package auth

import (
	"crypto/rsa"
	"errors"
	"io/ioutil"
	"testing"

	"github.com/dgrijalva/jwt-go"
	"github.com/stretchr/testify/assert"
)

func TestConfigure(t *testing.T) {
	type want struct {
		err     error
		issuer  string
		expiry  int64
		pubKey  bool
		privKey bool
	}
	tests := map[string]struct {
		config Config
		want
	}{
		"success": {
			config: Config{
				JwtIssuer:      "test.jwt.issuer",
				JwtExpiry:      600,
				PublicKeyPath:  "../../keys/test.public.pem",
				PrivateKeyPath: "../../keys/test.private.pem",
			},
			want: want{
				err:     nil,
				issuer:  "test.jwt.issuer",
				expiry:  600,
				pubKey:  true,
				privKey: true,
			},
		},
		"missing public key": {
			config: Config{
				JwtIssuer:      "test.jwt.issuer",
				JwtExpiry:      600,
				PrivateKeyPath: "../../keys/test.private.pem",
			},
			want: want{
				err:     errors.New("public key path must be specified"),
				issuer:  "",
				expiry:  0,
				pubKey:  false,
				privKey: false,
			},
		},
		"missing issuer and expiry": {
			config: Config{
				PublicKeyPath:  "../../keys/test.public.pem",
				PrivateKeyPath: "../../keys/test.private.pem",
			},
			want: want{
				err:     nil,
				issuer:  "",
				expiry:  0,
				pubKey:  true,
				privKey: true,
			},
		},
		"missing private key": {
			config: Config{
				PublicKeyPath: "../../keys/test.public.pem",
			},
			want: want{
				err:     nil,
				issuer:  "",
				expiry:  0,
				pubKey:  true,
				privKey: false,
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(st *testing.T) {
			resetVars()

			err := Configure(test.config)

			assert.Equal(st, test.err, err)
			assert.Equal(st, test.issuer, jwtIssuer)
			assert.Equal(st, test.expiry, jwtExpiry)

			wantPublicKey, wantPrivateKey := getTestKeys(test.pubKey, test.privKey)
			assert.Equal(st, wantPublicKey, publicKey)
			assert.Equal(st, wantPrivateKey, privateKey)
		})
	}
}

func resetVars() {
	jwtIssuer = ""
	jwtExpiry = 0
	publicKey = nil
	privateKey = nil
}

func getTestKeys(pubKey bool, privKey bool) (*rsa.PublicKey, *rsa.PrivateKey) {
	var wantPublicKey *rsa.PublicKey
	var wantPrivateKey *rsa.PrivateKey

	if pubKey {
		publicKeyFile, _ := ioutil.ReadFile("../../keys/test.public.pem")
		wantPublicKey, _ = jwt.ParseRSAPublicKeyFromPEM(publicKeyFile)
	}
	if privKey {
		privateKeyFile, _ := ioutil.ReadFile("../../keys/test.private.pem")
		wantPrivateKey, _ = jwt.ParseRSAPrivateKeyFromPEM(privateKeyFile)
	}

	return wantPublicKey, wantPrivateKey
}
