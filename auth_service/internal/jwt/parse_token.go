package jwt

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"
	"strconv"

	"github.com/golang-jwt/jwt/v5"
)

func ParseToken(tokenString string) (uint64, error) {
	block, _ := pem.Decode([]byte(os.Getenv("AUTH_SERVICE_PUBLIC_KEY")))
	if block == nil {
		return 0, fmt.Errorf("failed to parse PEM block containing the private key")
	}

	key, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return 0, err
	}
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
		return key.(*rsa.PublicKey), nil
	}, jwt.WithValidMethods([]string{jwt.SigningMethodRS256.Alg()}), jwt.WithIssuer("auth_service"))
	if err != nil && err != jwt.ErrTokenExpired {
		return 0, err
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		idStr, err := claims.GetSubject()
		if err != nil {
			return 0, err
		}
		id, err := strconv.ParseUint(idStr, 10, 64)
		return id, err
	}
	return 0, nil
}
