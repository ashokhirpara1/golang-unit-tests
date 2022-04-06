package others

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"
	"unit-tests/configuration"
	"unit-tests/prettylogs"

	"github.com/getsentry/sentry-go"
	"github.com/gorilla/mux"

	"gopkg.in/launchdarkly/go-sdk-common.v2/lduser"
	"gopkg.in/launchdarkly/go-sdk-common.v2/ldvalue"
	ld "gopkg.in/launchdarkly/go-server-sdk.v5"
)

type OtherStore struct {
	config configuration.Handler
	Log    prettylogs.Handler
}

func GetOther(config configuration.Handler, dLog prettylogs.Handler) OtherStorage {
	return &OtherStore{
		config: config,
		Log:    dLog,
	}
}

// Function to initialize sentry
func (c *OtherStore) InitializeSentry(dsn string, dLog *prettylogs.Handler) error {
	err := sentry.Init(sentry.ClientOptions{
		Dsn: dsn,
		// Enable printing of SDK debug messages.
		// Useful when getting started or trying to figure something out.
		Debug: true,
		BeforeSend: func(event *sentry.Event, hint *sentry.EventHint) *sentry.Event {
			// Here you can inspect/modify events before they are sent.
			// Returning nil drops the event.
			dLog.Info(fmt.Sprintf("BeforeSend event [%s]", event.EventID))
			return event
		},
	})
	if err != nil {
		return errors.New("sentry initialization failed")
	}
	return nil
}

// Middleware for closing all the connections / free up the resources. [will be called for each request]
func (c *OtherStore) CloseConnections(ldClient *ld.LDClient) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer ldClient.Flush()
			next.ServeHTTP(w, r)
		})
	}
}

// Function to set the LaunchDarkly Client
func (c *OtherStore) SetLaunchDarklyClient() {
	key := c.config.Clients.LD.Config.LD_SDK_KEY
	if key == "" {
		c.config.Clients.LD.LDClient = nil
	}
	client, err := ld.MakeClient(key, 5*time.Second)
	if err != nil {
		c.config.Clients.LD.LDClient = nil
	}
	c.config.Clients.LD.LDClient = client
}

// Evaluate
func (c *OtherStore) Evaluate(ldc *ld.LDClient, key string, defaultVal bool, data map[string]interface{}) bool {
	user, err := CreateUserWithCustomParams(data)
	if err != nil {
		log.Printf("Error Encountered : %s", err.Error())
		return defaultVal
	}
	flag, _ := ldc.BoolVariation(key, user, defaultVal)
	return flag
}

// Bool Variation with Anonymous User
func (c *OtherStore) AnonymousBool(ldc *ld.LDClient, key string, defaultVal bool) bool {
	flag, _ := ldc.BoolVariation(key, GetAnonymousUser(), defaultVal)
	return flag
}

// Middleware that wil recover from almost all kinds of panics and sends custom response
func (c *OtherStore) RecoveryMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				sentry.CurrentHub().Recover(err)
				sentry.Flush(time.Second * 5)

				jsonBody, _ := json.Marshal(map[string]string{
					"message": "something went wrong",
				})

				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusInternalServerError)
				w.Write(jsonBody)
			}
		}()
		next.ServeHTTP(w, r)
	})
}

// Function that returns random string
func GetRandomString() string {
	rand.Seed(time.Now().UnixNano())
	var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	b := make([]rune, 10)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

// Function that will return LaunchDirectly Anonymous User
func GetAnonymousUser() lduser.User {
	// Right now key doesn't matter
	return lduser.NewAnonymousUser("anonymous-user")
}

func CreateUserWithCustomParams(data map[string]interface{}) (lduser.User, error) {
	var (
		user       lduser.UserBuilder = lduser.NewUserBuilder(data["KEY"].(string))
		callOnUser bool               = true
		userObj    lduser.UserBuilderCanMakeAttributePrivate
	)
	for key, val := range data {
		if key != "KEY" {
			if callOnUser {
				switch val.(type) {
				case string:
					userObj = user.Custom(key, ldvalue.String(val.(string)))
				case int, int8, int16, int32, int64:
					userObj = user.Custom(key, ldvalue.Int(val.(int)))
				case float32, float64:
					userObj = user.Custom(key, ldvalue.Float64(val.(float64)))
				case bool:
					userObj = user.Custom(key, ldvalue.Bool(val.(bool)))
				case []string:
					builder := ldvalue.ArrayBuild()
					for _, v := range val.([]string) {
						builder.Add(ldvalue.String(v))
					}
					userObj = user.Custom(key, builder.Build())
				case []int, []int8, []int16, []int32, []int64:
					builder := ldvalue.ArrayBuild()
					for _, v := range val.([]int) {
						builder.Add(ldvalue.Int(v))
					}
					userObj = user.Custom(key, builder.Build())
				case []float32, []float64:
					builder := ldvalue.ArrayBuild()
					for _, v := range val.([]float64) {
						builder.Add(ldvalue.Float64(v))
					}
					userObj = user.Custom(key, builder.Build())
				default:
					return user.Build(), errors.New("Unsupported Type Passed")
				}
				callOnUser = false
			} else {
				switch val.(type) {
				case string:
					userObj = userObj.Custom(key, ldvalue.String(val.(string)))
				case int, int8, int16, int32, int64:
					userObj = userObj.Custom(key, ldvalue.Int(val.(int)))
				case float32, float64:
					userObj = userObj.Custom(key, ldvalue.Float64(val.(float64)))
				case bool:
					userObj = userObj.Custom(key, ldvalue.Bool(val.(bool)))
				case []string:
					builder := ldvalue.ArrayBuild()
					for _, v := range val.([]string) {
						builder.Add(ldvalue.String(v))
					}
					userObj = userObj.Custom(key, builder.Build())
				case []int, []int8, []int16, []int32, []int64:
					builder := ldvalue.ArrayBuild()
					for _, v := range val.([]int) {
						builder.Add(ldvalue.Int(v))
					}
					userObj = userObj.Custom(key, builder.Build())
				case []float32, []float64:
					builder := ldvalue.ArrayBuild()
					for _, v := range val.([]float64) {
						builder.Add(ldvalue.Float64(v))
					}
					userObj = userObj.Custom(key, builder.Build())
				default:
					return userObj.Build(), errors.New("Unsupported Type Passed")
				}
			}
		}
	}
	if userObj != nil {
		return userObj.Build(), nil
	}
	return user.Build(), nil
}
