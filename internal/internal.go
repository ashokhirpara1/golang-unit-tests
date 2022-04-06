package internal

import (
	"unit-tests/configuration"
	"unit-tests/mongo"
	"unit-tests/prettylogs"
)

// Ctlr is the big boss
type Ctlr struct {
	mongoDB mongo.MongoStorage
	config  configuration.Handler
	Log     *prettylogs.Handler
}

// InitController - return new Ctlr
func InitController(mongoDB mongo.MongoStorage, log *prettylogs.Handler, config configuration.Handler) *Ctlr {
	return &Ctlr{
		mongoDB: mongoDB,
		config:  config,
		Log:     log,
	}
}

func (c *Ctlr) GetLatestNotifications() ([]mongo.Notification, error) {
	return c.mongoDB.GetLatestNotifications()
}
