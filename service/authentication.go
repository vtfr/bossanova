package service

import (
	"errors"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/vtfr/bossanova/model"
)

// Authenticator handles all authentication roles
//go:generate mockgen -destination=../mocks/authentication_servcice.go -package=mocks github.com/vtfr/bossanova/service Authenticator
type Authenticator interface {
	// AuthenticateToken authenticates a token and returns it's user
	AuthenticateToken(tokenString string) (*model.User, error)

	// CreateToken creates a new JWT token
	CreateToken(user *model.User) (string, error)
}

// authenticator is the default implementation for Authenticator
// using JWT
type authenticator struct {
	users  UserGetter
	secret []byte
}

// UserGetter gets an user by their username
type UserGetter interface {
	GetUser(username string) (*model.User, error)
}

// NewAuthenticator creates a new Authenticator using the
// default implementation.
func NewAuthenticator(users UserGetter, secret []byte) Authenticator {
	return &authenticator{
		users:  users,
		secret: secret,
	}
}

// GetUserFromRequest gets an user from a request's Authorization token
func (auth *authenticator) AuthenticateToken(tokenString string) (*model.User, error) {
	token, err := jwt.ParseWithClaims(tokenString, &jwt.StandardClaims{},
		func(token *jwt.Token) (interface{}, error) {
			return auth.secret, nil
		})
	if err != nil {
		return nil, err
	}

	claims := token.Claims.(*jwt.StandardClaims)
	return auth.users.GetUser(claims.Subject)
}

// CreateToken creates a new JWT token
func (auth *authenticator) CreateToken(user *model.User) (string, error) {
	now := time.Now()
	return jwt.NewWithClaims(jwt.SigningMethodHS256, &jwt.StandardClaims{
		ExpiresAt: now.AddDate(0, 1, 0).Unix(),
		IssuedAt:  now.Unix(),
		Subject:   user.Username,
	}).SignedString(auth.secret)
}

// ExtractToken extracts a token from it's bearer prefix and returns it's raw
// value. Returns an error if no valid token is found
func ExtractToken(headerValue string) (string, error) {
	parts := strings.SplitN(headerValue, " ", 2)
	if strings.ToLower(parts[0]) == "bearer" && len(parts[1]) != 0 {
		return parts[1], nil
	}

	return "", errors.New("invalid bearer prefix")
}
