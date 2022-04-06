package routes

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
	"unit-tests/configuration"
	"unit-tests/internal"
	"unit-tests/mongo"
	mockmongo "unit-tests/mongo/mock"
	mockothers "unit-tests/others/mock"
	"unit-tests/prettylogs"

	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var router *mux.Router

func getMockRouter(db mongo.MongoStorage, mcothers *mockothers.MockOtherStorage) *mux.Router {
	// Get application specific configurations
	config := configuration.Get()

	// Initialize structured logs
	logs := prettylogs.Get()

	ctlr := internal.InitController(db, logs, config)

	// Initialize handler functions
	handler := initHandlers(ctlr, logs, config, mcothers)

	return getRouter(handler)
}

func bypassCloseConnectionMiddelware() mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			next.ServeHTTP(w, r)
		})
	}
}

func TestSearchNotifications(t *testing.T) {
	notifications := randomNotifications()
	ctlr := gomock.NewController(t)

	others := mockothers.NewMockOtherStorage(ctlr)
	others.EXPECT().
		InitializeSentry(gomock.Any(), gomock.Any()).
		Return(nil)
	others.EXPECT().
		RecoveryMiddleware(gomock.Any()).DoAndReturn(bypassCloseConnectionMiddelware()).Times(1)
	others.EXPECT().
		CloseConnections(gomock.Any()).Return(bypassCloseConnectionMiddelware()).Times(1)
	others.EXPECT().
		Evaluate(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
		Times(1).
		Return(true)

	store := mockmongo.NewMockMongoStorage(ctlr)
	store.EXPECT().
		GetLatestNotifications().
		Times(1).
		Return(notifications, nil)

	router = getMockRouter(store, others)

	req, err := http.NewRequest(http.MethodGet, "/mns/search", nil)
	require.NoError(t, err)

	values := req.URL.Query()
	values.Add("type", "ntype")
	values.Add("client_id", "xxx-xxx-xxx-xxx")
	req.URL.RawQuery = values.Encode()

	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, req)
	require.Equal(t, http.StatusOK, recorder.Code)
	decodeResponseBody(t, recorder.Body, notifications)
}

func randomNotifications() []mongo.Notification {
	notification := mongo.Notification{
		ID:                      primitive.NewObjectID(),
		SiteID:                  123,
		SiteSlug:                "xxxxx",
		UserID:                  123,
		TaskID:                  "ar-124",
		DistributorID:           "124",
		NotificationSlug:        "ar-123-xxxxx",
		Success:                 true,
		Disabled:                false,
		Created:                 time.Now(),
		NotificationAPIResponse: "",
	}
	return append([]mongo.Notification{}, notification)
}

func decodeResponseBody(t *testing.T, resBody *bytes.Buffer, notifications []mongo.Notification) {
	var respData Response

	body, err := ioutil.ReadAll(resBody)
	require.NoError(t, err)

	err = json.Unmarshal(body, &respData)
	require.NoError(t, err)

	require.Equal(t, notifications, respData.Data)
}
