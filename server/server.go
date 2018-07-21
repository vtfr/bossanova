package server

import (
	"encoding/json"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/kelseyhightower/envconfig"
	"github.com/vtfr/bossanova/service"
	"github.com/vtfr/bossanova/store"
)

func Start() error {
	var config struct {
		Address       string   `default:":8080"`
		MongoAddrs    []string `default:":27017"`
		MongoUsername string
		MongoPassword string
		Production    bool `default:"false"`
	}

	// read configuration
	if err := envconfig.Process("bossanova", &config); err != nil {
		return err
	}

	// connect to store
	store, err := store.NewMongoStore(config.MongoAddrs, config.MongoUsername,
		config.MongoPassword)
	if err != nil {
		return err
	}
	defer store.Close()

	authentication := service.NewAuthenticator(store, []byte("test"))

	authorization, err := ReadAuthorizationConfig()
	if err != nil {
		return err
	}

	// starts the server
	return http.ListenAndServe(config.Address, Route(store, authentication, authorization))
}

// ReadAuthorizationConfig read authorization configuration file
func ReadAuthorizationConfig() (service.Authorizator, error) {
	perms := make(map[string][]string)

	file, err := os.Open("permissions.json")
	if err != nil {
		return nil, err
	}
	defer file.Close()

	if err := json.NewDecoder(file).Decode(&perms); err != nil {
		return nil, err
	}

	return service.NewAuthorizator(perms), nil
}

// Route
func Route(store store.Store,
	authentication service.Authenticator,
	authorization service.Authorizator) http.Handler {
	r := gin.Default()

	r.Use(StoreMiddleware(store))
	r.Use(AuthenticationMiddleware(authentication))
	r.Use(ErrorMiddleware())

	api := r.Group("/api")

	// login
	api.POST("/login", login(authentication))

	// boards
	api.GET("/boards", HasAuthorization(authorization, "boards.read"), getBoards)
	api.POST("/boards", HasAuthorization(authorization, "boards.create"), createBoard)
	api.GET("/boards/:uri", HasAuthorization(authorization, "boards.read"), getBoard)
	api.PUT("/boards/:uri", HasAuthorization(authorization, "boards.update"), updateBoard)
	api.DELETE("/boards/:uri", HasAuthorization(authorization, "boards.create"), deleteBoard)

	// users
	api.GET("/users", HasAuthorization(authorization, "users.read"), getUsers)
	api.POST("/users", HasAuthorization(authorization, "users.create"), createUser)
	api.GET("/users/:username", HasAuthorization(authorization, "users.read"), getUser)
	api.PUT("/users/:username", HasAuthorization(authorization, "users.update"), updateUser)
	api.DELETE("/users/:username", HasAuthorization(authorization, "users.create"), deleteUser)

	// post
	api.GET("/post/:id", HasAuthorization(authorization, "post.read"), getPost)
	api.DELETE("/post/:id", HasAuthorization(authorization, "post.delete"), deletePost)
	api.PUT("/post/:id", HasAuthorization(authorization, "post.update"), updatePost)

	// thread
	api.GET("/thread/:id", HasAuthorization(authorization, "post.read"), getThread)

	// post creation
	api.POST("/post", post(authorization))

	return r
}
