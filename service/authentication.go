package service

import (
	"net/http"

	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/dgrijalva/jwt-go/request"
	"github.com/vtfr/bossanova/model"
)

// Authenticator handles all authentication roles
//go:generate mockgen -destination=../mocks/authentication_servcice.go -package=mocks github.com/vtfr/bossanova/service Authenticator
type Authenticator interface {
	// GetUserFromRequest gets an user from a request's Authorization token
	GetUserFromRequest(req *http.Request) (*model.User, error)

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
func (auth *authenticationService) GetUserFromRequest(req *http.Request) (*model.User, error) {
	token, err := request.ParseFromRequestWithClaims(req, request.AuthorizationHeaderExtractor,
		&jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
			return auth.secret, nil
		})
	if err != nil {
		return nil, err
	}

	username := token.Claims.(*jwt.StandardClaims).Subject
	return auth.users.GetUser(username)
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
