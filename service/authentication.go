package service

import (
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/vtfr/bossanova/model"
)

// ErrNoToken is returned when no token has been found
var ErrNoToken = errors.New("no token")

// Authenticator handles all authentication roles
//go:generate mockgen -destination=../mocks/authentication_servcice.go -package=mocks github.com/vtfr/bossanova/service Authenticator
type Authenticator interface {
	// AuthenticateToken authenticates a token and returns it's user
	AuthenticateToken(tokenString string) (*model.User, error)

	// CreateToken creates a new JWT token
	CreateToken(user *model.User) (string, error)
}

// authenticationService is the default implementation for Authenticator
// using JWT
type authenticationService struct {
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
	return &authenticationService{
		users:  users,
		secret: secret,
	}
}

// GetUserFromRequest gets an user from a request's Authorization token
func (auth *authenticationService) AuthenticateToken(tokenString string) (*model.User, error) {
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
func (auth *authenticationService) CreateToken(user *model.User) (string, error) {
	now := time.Now()
	return jwt.NewWithClaims(jwt.SigningMethodHS256, &jwt.StandardClaims{
		ExpiresAt: now.AddDate(0, 1, 0).Unix(),
		IssuedAt:  now.Unix(),
		Subject:   user.Username,
	}).SignedString(auth.secret)
}

// ExtractToken extracts a token from it's Authorization headers. Returns an
// error if no token is found
func ExtractToken(req *http.Request) (string, error) {
	const authKey = "Authorization"
	const bearerPrefix = "bearer "

	value := req.Header.Get("Authorization")
	if value == "" {
		return "", ErrNoToken
	}

	// Searches for bearer in the beginning of the string. If found, remove it
	if strings.HasPrefix(strings.ToLower(value), bearerPrefix) {
		return value[len(bearerPrefix):], nil
	}

	return "", errors.New("invalid bearer prefix")
}
