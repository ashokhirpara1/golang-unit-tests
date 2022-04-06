package others

import (
	"net/http"
	"unit-tests/prettylogs"

	"github.com/gorilla/mux"
	ld "gopkg.in/launchdarkly/go-server-sdk.v5"
)

type OtherStorage interface {
	InitializeSentry(dsn string, dLog *prettylogs.Handler) error
	CloseConnections(ldClient *ld.LDClient) mux.MiddlewareFunc
	SetLaunchDarklyClient()
	Evaluate(ldc *ld.LDClient, key string, defaultVal bool, data map[string]interface{}) bool
	AnonymousBool(ldc *ld.LDClient, key string, defaultVal bool) bool
	RecoveryMiddleware(next http.Handler) http.Handler
}
