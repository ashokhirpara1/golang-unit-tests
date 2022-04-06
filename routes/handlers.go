package routes

import (
	"net/http"
	"strings"
	"unit-tests/configuration"
	"unit-tests/internal"
	"unit-tests/others"
	"unit-tests/prettylogs"
)

// handler - holds routes struct
type handler struct {
	ctlr   *internal.Ctlr
	Log    *prettylogs.Handler
	config configuration.Handler
	others others.OtherStorage
}

// initHandlers - intializing handler functions
func initHandlers(ctlr *internal.Ctlr, prettylogs *prettylogs.Handler, config configuration.Handler, mcothers others.OtherStorage) *handler {
	handler := handler{ctlr: ctlr, Log: prettylogs, config: config, others: mcothers}

	return &handler
}

// API middleware to check a valid client id
func (c *handler) apiMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Request Query Parameter client_id check
		clientID := r.URL.Query().Get("client_id")
		if clientID == "" {
			sendResponse(w, r, http.StatusBadRequest, "Missing required parameter client_id", nil)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (c *handler) isRouteEnabled(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var data map[string]interface{} = map[string]interface{}{
			"KEY":       "analytics-" + strings.ReplaceAll(strings.TrimLeft(r.URL.Path, "/"), "/", "-"), // dls-stagnant-users
			"route":     r.URL.Path,
			"client-id": r.URL.Query().Get("client_id"),
		}
		if !c.others.Evaluate(c.config.Clients.LD.LDClient, c.config.Clients.LD.Config.Flags.ROUTE_FLAG, false, data) {
			sendResponse(w, r, http.StatusNotFound, "Not Found", nil)
			return
		}
		next.ServeHTTP(w, r)
	})
}

// SearchNotifications
func (c *handler) SearchNotifications(w http.ResponseWriter, r *http.Request) {
	defer c.Log.Exit(c.Log.Enter())

	// marketing_drip or pulse
	nType := r.URL.Query().Get("type")
	if nType == "" {
		sendResponse(w, r, http.StatusBadRequest, "invalid value for type field", nil)
		return
	}

	notifications, err := c.ctlr.GetLatestNotifications()
	if err != nil {
		sendResponse(w, r, http.StatusAccepted, err.Error(), nil)
		return
	}

	sendResponse(w, r, http.StatusOK, "", notifications)
}
