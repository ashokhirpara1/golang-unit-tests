package mongo

// mockgen -package mongo -destination mongo/mock_mongo.go unit-tests/mongo MongoStorage
// MongoStorage - holds functions which can be used by outside of package
type MongoStorage interface {
	GetLatestNotifications() ([]Notification, error)
}
