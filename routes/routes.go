package routes

import (
	"context"
	"fmt"
	"net/http"
	"unit-tests/configuration"
	"unit-tests/internal"
	"unit-tests/others"
	"unit-tests/prettylogs"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/awslabs/aws-lambda-go-api-proxy/core"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"

	muxproxy "github.com/awslabs/aws-lambda-go-api-proxy/gorillamux"
)

var adapter *muxproxy.GorillaMuxAdapter

// InitRoutes - init routs and handlers
func InitRoutes(ctlr *internal.Ctlr, prettylogs *prettylogs.Handler, config configuration.Handler) {

	others := others.GetOther(config, *prettylogs)

	others.SetLaunchDarklyClient()

	// Initialize handler functions
	handler := initHandlers(ctlr, prettylogs, config, others)

	router := getRouter(handler)

	adapter = muxproxy.New(router)
	lambda.Start(lambdaHttpHandler)
}

func lambdaHttpHandler(ctx context.Context, req core.SwitchableAPIGatewayRequest) (*core.SwitchableAPIGatewayResponse, error) {
	c, err := adapter.ProxyWithContext(ctx, req)
	return c, err
}

func getRouter(handler *handler) *mux.Router {

	clients := handler.config.Clients

	// Creates a http server
	router := mux.NewRouter()

	subRouter := router.PathPrefix("/mns").Subrouter()

	subRouter.HandleFunc("/search", handler.SearchNotifications).Methods(http.MethodGet)

	subRouter.Use(handler.apiMiddleware)
	subRouter.Use(handler.isRouteEnabled)

	cors := handlers.CORS(
		handlers.AllowedHeaders([]string{"Origin", "X-Requested-With", "Access-Control-Request-Headers", "Access-Control-Request-Method", "Content-Type", "Accept"}),
		handlers.AllowedOrigins([]string{"*"}),
		handlers.AllowCredentials(),
	)
	router.Use(cors)

	if err := handler.others.InitializeSentry(clients.Sentry.Config.SENTRY_DSN, handler.Log); err != nil {
		handler.Log.Info(fmt.Sprintf("Sentry Initialization failed %v", err))
	} else {
		subRouter.Use(handler.others.RecoveryMiddleware)
	}

	subRouter.Use(handler.others.CloseConnections(clients.LD.LDClient))

	return router
}
