package server

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/kelseyhightower/envconfig"
	"github.com/sirupsen/logrus"
	"github.com/vtfr/bossanova/service"
	"github.com/vtfr/bossanova/store"
)

// Start creates a new Bossanova server with the configuration retrieved from
// environment variable and runs it
func Start() error {
	var config struct {
		Address      string `default:":8080"`
		MongoAddress string `default:"mongodb://localhost:27017"`
		Production   bool   `default:"false"`
	}

	// read configuration
	logrus.Infoln("Reading environment variable configuration")
	if err := envconfig.Process("bossanova", &config); err != nil {
		return err
	}

	// connect to store
	logrus.Infoln("Attempting to connect to MongoDB store")
	store, err := store.NewMongoStore(config.MongoAddress)
	if err != nil {
		return err
	}
	defer store.Close()

	logrus.Infoln("Creating Authentication service")
	authentication := service.NewAuthenticator(store, []byte("test"))

	logrus.Infoln("Creating Authorization service")
	authorization, err := newAuthorizator()
	if err != nil {
		return err
	}

	// routing
	logrus.Infoln("Configuring routing")
	handler := Route(store, authentication, authorization)

	// starts the server
	logrus.Infoln("Starting server at", config.Address)
	return http.ListenAndServe(config.Address, handler)
}

// newAuthorizator read authorization configuration file
func newAuthorizator() (service.Authorizator, error) {
	file, err := os.Open("permissions.json")
	if err != nil {
		return nil, err
	}
	defer file.Close()

	return service.NewAuthorizatorFromFile(file)
}

// Route configures all necessary middlewares and routes needed
func Route(store store.Store,
	authentication service.Authenticator,
	authorization service.Authorizator) http.Handler {
	r := gin.Default()

	r.Use(ErrorMiddleware())
	r.Use(StoreMiddleware(store))
	r.Use(AuthenticationMiddleware(authentication))

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
	api.POST("/post", HasAuthorization(authorization, "posting.post"), post(authorization))

	return r
}
