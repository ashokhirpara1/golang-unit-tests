package mongo

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"time"
	"unit-tests/configuration"
	"unit-tests/prettylogs"

	mongod "go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

//DB - connection instance
type DB struct {
	DocumentDB *mongod.Database
	Log        *prettylogs.Handler
}

// Get - Returns an instance of the database client
func Get(config configuration.Handler, dLog *prettylogs.Handler) (MongoStorage, error) {

	db, err := get(config, dLog)
	if err != nil {
		dLog.DBError("Failed to DB Connection: ", err)
		return nil, err
	}

	return &DB{
		DocumentDB: db,
		Log:        dLog,
	}, nil
}

/* returns an instance of the databse
 *  documentDB is enabled with TLS, hence need to provide the aws combined certificate info when connecting to documentDB
 *  connects to document
 *  pings the database connection to verify the connection was successful
 */
func get(config configuration.Handler, dLog *prettylogs.Handler) (*mongod.Database, error) {

	connectionURI := config.GetMongoDBConnStr()

	tlsConfig, err := getCustomTLSConfig(config.MongoDB.CAFilePath)
	if err != nil {
		return nil, err
	}

	client, err := mongod.NewClient(options.Client().ApplyURI(connectionURI).SetTLSConfig(tlsConfig))
	if err != nil {
		dLog.DBError("Error in creating db client, URI: "+connectionURI, err)
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(config.MongoDB.ConnectionTimeout)*time.Second)
	defer cancel()

	err = client.Connect(ctx)
	if err != nil {
		dLog.DBError("Failed to Connect DB, URI: "+connectionURI, err)
		return nil, err
	}

	// Force a connection to verify our connection string
	err = client.Ping(ctx, nil)
	if err != nil {
		dLog.DBError("Failed to ping DB, URI: "+connectionURI, err)
		return nil, err
	}

	documentDB := client.Database(config.MongoDB.Name)

	return documentDB, nil
}

/*
 * documentDB is enabled with TLS,
 * reads the aws combined cert file and returns the tls configuration
 */
func getCustomTLSConfig(caFile string) (*tls.Config, error) {
	tlsConfig := new(tls.Config)

	currDir, _ := os.Getwd()
	caFile = fmt.Sprintf("%s/%s", currDir, caFile)
	certs, err := ioutil.ReadFile(caFile)

	if err != nil {
		fmt.Printf("Error getCustomTLSConfig.ReadFile : %s", err.Error())
		return tlsConfig, err
	}

	tlsConfig.RootCAs = x509.NewCertPool()
	ok := tlsConfig.RootCAs.AppendCertsFromPEM(certs)

	if !ok {
		err = errors.New("failed parsing pem file")
		fmt.Printf("Error getCustomTLSConfig.AppendCertsFromPEM : %s", err.Error())
		return tlsConfig, err
	}

	return tlsConfig, nil
}
