// This file is safe to edit. Once it exists it will not be overwritten

package restapi

import (
	"context"
	"crypto/tls"
	"github.com/cegorah/blockchain_info/auth"
	"github.com/cegorah/blockchain_info/cache"
	"github.com/cegorah/blockchain_info/handlers"
	"github.com/cegorah/blockchain_info/repo"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"net/http"
	"strings"
	"time"

	"github.com/cegorah/blockchain_info/restapi/operations"
	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
)

//go:generate swagger generate server --target ../../blockchain_info --name BlockchainInfo --spec ../docs/swagger.json --principal interface{}

func configureFlags(api *operations.BlockchainInfoAPI) {
	// api.CommandLineOptionsGroups = []swag.CommandLineOptionsGroup{ ... }
}

func configureAPI(api *operations.BlockchainInfoAPI) http.Handler {
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

	// Applies when the "Authorization" header is set

	api.BearerAuth = auth.TokenValidate
	lg := logrus.New()

	api.Logger = lg.Printf

	// Set your custom authorizer if needed. Default one is security.Authorized()
	// Expected interface runtime.Authorizer
	//
	// Example:
	// api.APIAuthorizer = security.Authorized()
	serviceInitTimeout := viper.GetInt("services.init_timeout")
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(serviceInitTimeout)*time.Second)
	defer cancel()
	psClient, err := repo.NewPsClient(ctx, viper.GetString("PSQL_DSN"))
	if db := viper.GetBool("debug"); db {
		lg.SetLevel(logrus.DebugLevel)
	}
	if err != nil {
		lg.Fatalf("psClient initialization error: %s", err)
	}
	cacheServer := cache.NewRedisServer(map[string]interface{}{
		"addr":     viper.GetString("services.redis.connection_string"),
		"username": viper.GetString("redis_username"),
		"password": viper.GetString("redis_password"),
		"db":       viper.GetInt("services.redis.db"),
		"ttl":      viper.GetInt("services.redis.ttl"),
	})

	qt := viper.GetInt("services.query_timeout")

	api.BlockchainInfoGetBlockHandler = &handlers.BlockInfoImpl{
		DefaultTimeoutSecond: qt,
		DBClient:             &psClient,
		Cache:                &cacheServer,
		Logger:               lg,
	}
	api.BlockchainInfoGetTxInfoHandler = &handlers.TxInfoImpl{
		DefaultTimeoutSecond: qt,
		DBClient:             &psClient,
		Cache:                &cacheServer,
		Logger:               lg,
	}

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
// scheme value will be set accordingly: "http", "https" or "unix".
func configureServer(s *http.Server, scheme, addr string) {
}

// The middleware configuration is for the handler executors. These do not apply to the swagger.json document.
// The middleware executes after routing but before authentication, binding and validation.
func setupMiddlewares(handler http.Handler) http.Handler {
	return handler
}

// The middleware configuration happens before anything, this middleware also applies to serving the swagger.json document.
// So this is a good place to plug in a panic handling middleware, logging and metrics.
func setupGlobalMiddleware(handler http.Handler) http.Handler {
	return uiMiddleware(handler)
}

func uiMiddleware(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Shortcut helpers for swagger-ui
		if r.URL.Path == "/swagger-ui" || r.URL.Path == "/api/help" {
			http.Redirect(w, r, "/swagger-ui/", http.StatusFound)
			return
		}
		// Serving ./swagger-ui/
		if strings.Index(r.URL.Path, "/swagger-ui/") == 0 {
			http.StripPrefix("/swagger-ui/", http.FileServer(http.Dir("swagger-ui"))).ServeHTTP(w, r)
			return
		}
		handler.ServeHTTP(w, r)
	})
}
