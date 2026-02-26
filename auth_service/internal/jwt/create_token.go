package jwt

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func CreateToken(id uint64) (string, error) {
	claims := jwt.RegisteredClaims{Subject: fmt.Sprintf("%v", id), Issuer: "auth_service", ExpiresAt: &jwt.NumericDate{Time: time.Now().Add(time.Hour * 12)}}
	tokenUnsigned := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)

	block, _ := pem.Decode([]byte(os.Getenv("AUTH_SERVICE_PRIVATE_KEY")))
	if block == nil {
		return "", fmt.Errorf("failed to parse PEM block containing the private key")
	}

	key, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return "", err
	}

	return tokenUnsigned.SignedString(key.(*rsa.PrivateKey))
}
