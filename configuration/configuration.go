package configuration

import (
	"flag"
	"fmt"
	"os"
)

// Handler holds data necessary for configuring application
type Handler struct {
	AWS     AWS     `json:"aws"`
	General General `json:"general"`
	MongoDB MongoDB `json:"mongodb"`
	Clients Clients `json:"-"`
}

// MongoDB holds data for database configuring
type MongoDB struct {
	Template          string
	Name              string
	Username          string
	Password          string
	ClusterEndpoint   string
	ReadPreference    string
	CAFilePath        string
	DBPort            string
	ConnectionTimeout int
	QueryTimeout      int
}

// General - holds general data for configuration
type General struct {
	AppStage           string
	HTTPRequestTimeout int
}

type AWS struct {
	MN_QUEUE_URL          string
	SSM_SERVICE_BASE_PATH string
}

const (
	HTTP_REQUEST_TIMEOUT = 90
	MAX_OPEN_CONNS       = 25
	MAX_IDLE_CONNS       = 25
	dbTemplate           = "mongodb://%s:%s@%s:%s/%s?ssl=true&ssl_ca_certs=bundle.pem&replicaSet=rs0&readPreference=%s"
	readPreference       = "secondaryPreferred"
	CA_FILE_PATH         = "bundle.pem"
	DB_CONN_TIMEOUT      = 60
	QUERY_TIMEOUT        = 60
)

// Get configurations
func Get() (config Handler) {

	flag.StringVar(&config.AWS.MN_QUEUE_URL, "MNQueueUrl", os.Getenv("MN_QUEUE_URL"), "AWS Queue Url")
	flag.StringVar(&config.AWS.SSM_SERVICE_BASE_PATH, "ssmServiceBasePath", os.Getenv("SSM_SERVICE_BASE_PATH"), "Basic search path for SSM")

	flag.IntVar(&config.General.HTTPRequestTimeout, "httprequesttimeout", HTTP_REQUEST_TIMEOUT, "http request timeout")
	flag.StringVar(&config.General.AppStage, "AppStage", os.Getenv("APP_STAGE"), "Application environment")

	flag.StringVar(&config.MongoDB.Username, "Username", os.Getenv("DB_MASTER_USER"), "database username")
	flag.StringVar(&config.MongoDB.Password, "Password", os.Getenv("DB_MASTER_PASSWORD"), "database password")
	flag.StringVar(&config.MongoDB.Template, "Template", dbTemplate, "database template")
	flag.StringVar(&config.MongoDB.ClusterEndpoint, "ClusterEndpoint", os.Getenv("DB_CLUSTER_ENDPOINT"), "database cluster endpoint")
	flag.StringVar(&config.MongoDB.DBPort, "DBPort", os.Getenv("DB_CLUSTER_PORT"), "database port")
	flag.StringVar(&config.MongoDB.Name, "Name", os.Getenv("DB_NAME"), "database name")
	flag.StringVar(&config.MongoDB.ReadPreference, "ReadPreference", readPreference, "database read preference")
	flag.StringVar(&config.MongoDB.CAFilePath, "CAFilePath", CA_FILE_PATH, "database secret file")
	flag.IntVar(&config.MongoDB.ConnectionTimeout, "ConnectionTimeout", DB_CONN_TIMEOUT, "database connection timeout")
	flag.IntVar(&config.MongoDB.QueryTimeout, "QueryTimeout", QUERY_TIMEOUT, "database query timeout")

	config.Clients.Sentry.Config.SENTRY_DSN = os.Getenv("SENTRY_DSN_KEY_NAME")
	config.Clients.LD.Config.LD_SDK_KEY = os.Getenv("LD_SDK_KEY_NAME")

	config.Clients.LD.Config.Flags.ROUTE_FLAG = "backend-crm-mobile-notifications-evt-analytics-route"

	return config
}

// GetMock set env for mock configurations
func GetMock() (config Handler) {

	os.Setenv("MN_QUEUE_URL", "mn_queue_url")
	os.Setenv("BASE_DOMAIN", "127.0.0.1")
	os.Setenv("CLIENT_ID", "xxx-xxx-xxx-xxx")
	os.Setenv("APP_STAGE", "dev")

	os.Setenv("_LAMBDA_SERVER_PORT", "55555")
	os.Setenv("AWS_LAMBDA_RUNTIME_API", "127.0.0.1")

	return Get()
}

// GetMongoDBConnStr handler
func (c *Handler) GetMongoDBConnStr() (connectionURI string) {

	connectionURI = fmt.Sprintf(c.MongoDB.Template, c.MongoDB.Username, c.MongoDB.Password, c.MongoDB.ClusterEndpoint, c.MongoDB.DBPort, c.MongoDB.Name, c.MongoDB.ReadPreference)

	return
}
