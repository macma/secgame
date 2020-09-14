// This file is safe to edit. Once it exists it will not be overwritten

package restapi

import (
	"crypto/tls"
	"fmt"
	"net/http"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/runtime/middleware"

	"game/db"
	"game/reactjs"
	"game/restapi/operations"
)

//go:generate swagger generate server --target ../../game --name ClpsecApp --spec ../swagger.json --principal interface{}

func configureFlags(api *operations.ClpsecAppAPI) {
	// api.CommandLineOptionsGroups = []swag.CommandLineOptionsGroup{ ... }
}

func configureAPI(api *operations.ClpsecAppAPI) http.Handler {
	// configure the api here
	api.ServeError = errors.ServeError

	// Set your custom logger if needed. Default one is log.Printf
	// Expected interface func(string, ...interface{})
	//
	// Example:
	// api.Logger = log.Printf

	api.UseSwaggerUI()
	// To continue using redoc as your UI, uncomment the following line
	// api.UseRedoc()

	api.JSONConsumer = runtime.JSONConsumer()

	api.JSONProducer = runtime.JSONProducer()

	// if api.OrangePressedHandler == nil {
	api.OrangePressedHandler = operations.OrangePressedHandlerFunc(func(params operations.OrangePressedParams) middleware.Responder {
		str := reactjs.ProcessClickForSwagger(0)
		res := operations.OrangePressedOKBody{Message: str}

		return &operations.OrangePressedOK{Payload: &res}
	})
	api.BluePressedHandler = operations.BluePressedHandlerFunc(func(params operations.BluePressedParams) middleware.Responder {
		str := reactjs.ProcessClickForSwagger(1)
		res := operations.BluePressedOKBody{Message: str}

		return &operations.BluePressedOK{Payload: &res}
	})
	api.QueryHandler = operations.QueryHandlerFunc(func(params operations.QueryParams) middleware.Responder {
		val := params.GqlFields.GqlField

		str := reactjs.QueryForSwagger(*val)
		res := operations.QueryOKBody{Message: str}

		return &operations.QueryOK{Payload: &res}
	})
	// }

	api.PreServerShutdown = func() {}

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
func configureServer(s *http.Server, scheme, addr string) {
}

// The middleware configuration is for the handler executors. These do not apply to the swagger.json document.
// The middleware executes after routing but before authentication, binding and validation
func setupMiddlewares(handler http.Handler) http.Handler {

	// var dbconn *sql.DB
	// dbconn =
	db.InitDB()
	fmt.Println("this is init of middle ware")
	return handler
}

// The middleware configuration happens before anything, this middleware also applies to serving the swagger.json document.
// So this is a good place to plug in a panic handling middleware, logging and metrics
func setupGlobalMiddleware(handler http.Handler) http.Handler {
	go reactjs.ServeWS()
	fmt.Println("this is init of global middleware")
	go reactjs.ServeHTML()
	return handler
}
