package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/vtfr/bossanova/common"
	"github.com/vtfr/bossanova/model"
	"github.com/vtfr/bossanova/service"
	"github.com/vtfr/bossanova/store"
)

const (
	storeKey    = "store"
	configKey   = "config"
	usernameKey = "username"
	userKey     = "user"
	expiresKey  = "expires"
	tokenKey    = "token"
)

// ErrorMiddleware handles error processing
func ErrorMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		// If errors
		if len(c.Errors) != 0 {
			err := c.Errors.Last().Err
			var detailed *common.DetailedError
			if er, ok := err.(*common.DetailedError); ok {
				detailed = er
			} else {
				detailed = common.NewDetailedError("internal-server-error",
					err.Error(), http.StatusInternalServerError)
			}

			c.JSON(detailed.Status, gin.H{
				"error": detailed,
			})
		}
	}
}

// StoreMiddleware creates a new connection for each request providing a fast
// and parallelable store backend access
func StoreMiddleware(store store.Store) gin.HandlerFunc {
	return func(c *gin.Context) {
		st := store.Clone()
		defer st.Close()

		c.Set(storeKey, st)
		c.Next()
	}
}

// GetStore returns a store saved in a context
func GetStore(c *gin.Context) store.Store {
	return c.MustGet(storeKey).(store.Store)
}

// AuthenticationMiddleware authenticates an user and sets it's context value
func AuthenticationMiddleware(auth service.Authenticator) gin.HandlerFunc {
	return func(c *gin.Context) {
		// if no token is found, simply ignore
		if headerValue := c.GetHeader("Authentication"); headerValue != "" {
			token, err := service.ExtractToken(headerValue)
			if err != nil {
				c.Error(err)
				return
			}

			// attempts to parse user
			user, err := auth.AuthenticateToken(token)
			if err != nil {
				c.Error(err)
				return
			}

			// stores token and user
			c.Set(tokenKey, token)
			c.Set(userKey, user)
		}

		c.Next()
	}
}

// GetUser returns the current context user
func GetUser(c *gin.Context) *model.User {
	if user, ok := c.Value(userKey).(*model.User); ok {
		return user
	}

	return nil
}

// HasAuthorization verifies if a given user has enough authorization to perform
// an action
func HasAuthorization(auth service.Authorizator, action string) gin.HandlerFunc {
	return func(c *gin.Context) {
		if !auth.IsAuthorized(GetUser(c), action) {
			c.AbortWithStatus(http.StatusForbidden)
		}
	}

}
