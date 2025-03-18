package utils

import (
	"time"

	jwt "github.com/golang-jwt/jwt/v5"

	mErr "github.com/itmrchow/todolist-users/internal/errors"
)

// GenerateToken generates a JWT token for a user
func GenerateToken(userID string, secretKey string, issuer string, expireAt int) (tokenStr string, err error) {
	now := time.Now()

	registeredClaims := jwt.RegisteredClaims{
		IssuedAt:  jwt.NewNumericDate(now),
		NotBefore: jwt.NewNumericDate(now),
		ExpiresAt: jwt.NewNumericDate(now.Add(time.Hour * time.Duration(expireAt))),
		Issuer:    issuer,
		Subject:   userID,
		Audience:  []string{userID},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS512, registeredClaims)

	tokenStr, err = token.SignedString([]byte(secretKey))
	if err != nil {
		return "", err
	}

	return
}

// ValidateToken validates a JWT token
func ValidateToken(tokenStr string, secretKey string, issuer string) (userID string, err error) {

	if tokenStr == "" {
		return "", &mErr.Err401Unauthorized
	}

	// parse token
	tokenClaims, err := jwt.ParseWithClaims(tokenStr, &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	}, jwt.WithLeeway(5*time.Second))

	if err != nil {
		return
	}

	claims, ok := tokenClaims.Claims.(*jwt.RegisteredClaims)
	if !ok {
		return "", &mErr.Err401Unauthorized
	}

	if claims.Subject == "" {
		return "", &mErr.Err401Unauthorized
	}

	if claims.Issuer != issuer {
		return "", &mErr.Err401Unauthorized
	}

	return
}
