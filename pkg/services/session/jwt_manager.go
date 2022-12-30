package session

import (
	"crypto/ed25519"
	"crypto/rand"
	"fmt"

	"github.com/golang-jwt/jwt/v4"
)

type JWTManager struct {
	secretKey ed25519.PrivateKey
	publicKey ed25519.PublicKey
}

func NewJWTManager() *JWTManager {
	pubKey, privKey, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		panic(err)
	}

	return &JWTManager{
		secretKey: privKey,
		publicKey: pubKey,
	}
}

type UserClaims struct {
	jwt.RegisteredClaims
	UserID string `json:"userID"`
}

func (manager *JWTManager) Generate(userID string) (string, error) {
	claims := UserClaims{
		RegisteredClaims: jwt.RegisteredClaims{},
		UserID:           userID,
	}
	t := jwt.NewWithClaims(&jwt.SigningMethodEd25519{}, claims)
	return t.SignedString(manager.secretKey)
}

func (manager *JWTManager) Verify(accessToken string) (*UserClaims, error) {
	token, err := jwt.ParseWithClaims(
		accessToken,
		&UserClaims{},
		func(token *jwt.Token) (interface{}, error) {
			_, ok := token.Method.(*jwt.SigningMethodEd25519)
			if !ok {
				return nil, fmt.Errorf("unexpected token signing method")
			}

			return manager.publicKey, nil
		},
	)
	if err != nil || !token.Valid {
		return nil, fmt.Errorf("invalid token: %w", err)
	}

	claims, ok := token.Claims.(*UserClaims)
	if !ok {
		return nil, fmt.Errorf("invalid token claims")
	}

	return claims, nil
}
