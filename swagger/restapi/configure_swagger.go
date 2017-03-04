package restapi

import (
	"crypto/tls"
	"net/http"

	errors "github.com/go-openapi/errors"
	runtime "github.com/go-openapi/runtime"
	middleware "github.com/go-openapi/runtime/middleware"
	graceful "github.com/tylerb/graceful"

	"github.com/hiromaily/go-gin-wrapper/swagger/restapi/operations"
	"github.com/hiromaily/go-gin-wrapper/swagger/restapi/operations/j_w_t"
	"github.com/hiromaily/go-gin-wrapper/swagger/restapi/operations/users"
)

// This file is safe to edit. Once it exists it will not be overwritten

//go:generate swagger generate server --target .. --name swagger --spec ../swagger.yaml

func configureFlags(api *operations.SwaggerAPI) {
	// api.CommandLineOptionsGroups = []swag.CommandLineOptionsGroup{ ... }
}

func configureAPI(api *operations.SwaggerAPI) http.Handler {
	// configure the api here
	api.ServeError = errors.ServeError

	// Set your custom logger if needed. Default one is log.Printf
	// Expected interface func(string, ...interface{})
	//
	// Example:
	// s.api.Logger = log.Printf

	api.JSONConsumer = runtime.JSONConsumer()

	api.JSONProducer = runtime.JSONProducer()

	api.UsersDeleteUsersIDHandler = users.DeleteUsersIDHandlerFunc(func(params users.DeleteUsersIDParams) middleware.Responder {
		return middleware.NotImplemented("operation users.DeleteUsersID has not yet been implemented")
	})
	api.UsersGetUsersHandler = users.GetUsersHandlerFunc(func(params users.GetUsersParams) middleware.Responder {
		return middleware.NotImplemented("operation users.GetUsers has not yet been implemented")
	})
	api.UsersGetUsersIDHandler = users.GetUsersIDHandlerFunc(func(params users.GetUsersIDParams) middleware.Responder {
		return middleware.NotImplemented("operation users.GetUsersID has not yet been implemented")
	})
	api.UsersGetUsersIdsHandler = users.GetUsersIdsHandlerFunc(func(params users.GetUsersIdsParams) middleware.Responder {
		return middleware.NotImplemented("operation users.GetUsersIds has not yet been implemented")
	})
	api.JWTPostJwtHandler = j_w_t.PostJwtHandlerFunc(func(params j_w_t.PostJwtParams) middleware.Responder {
		return middleware.NotImplemented("operation j_w_t.PostJwt has not yet been implemented")
	})
	api.UsersPostUsersHandler = users.PostUsersHandlerFunc(func(params users.PostUsersParams) middleware.Responder {
		return middleware.NotImplemented("operation users.PostUsers has not yet been implemented")
	})
	api.UsersPutUsersIDHandler = users.PutUsersIDHandlerFunc(func(params users.PutUsersIDParams) middleware.Responder {
		return middleware.NotImplemented("operation users.PutUsersID has not yet been implemented")
	})

	api.ServerShutdown = func() {}

	return setupGlobalMiddleware(api.Serve(setupMiddlewares))
}

// The TLS configuration before HTTPS server starts.
func configureTLS(tlsConfig *tls.Config) {
	// Make all necessary changes to the TLS configuration here.
}

// As soon as server is initialized but not run yet, this function will be called.
// If you need to modify a config, store server instance to stop it individually later, this is the place.
// This function can be called multiple times, depending on the number of serving schemes.
// scheme value will be set accordingly: "http", "https" or "unix"
func configureServer(s *graceful.Server, scheme string) {
}

// The middleware configuration is for the handler executors. These do not apply to the swagger.json document.
// The middleware executes after routing but before authentication, binding and validation
func setupMiddlewares(handler http.Handler) http.Handler {
	return handler
}

// The middleware configuration happens before anything, this middleware also applies to serving the swagger.json document.
// So this is a good place to plug in a panic handling middleware, logging and metrics
func setupGlobalMiddleware(handler http.Handler) http.Handler {
	return handler
}
