package mongo

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func getDay(day int) time.Time {
	t := time.Now().AddDate(0, 0, day)
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
}

func (d *DB) GetLatestNotifications() ([]Notification, error) {
	defer d.Log.Exit(d.Log.Enter())

	notifications := []Notification{}
	collection := d.DocumentDB.Collection("collection-name")
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(60)*time.Second)
	defer cancel()

	opts := options.Find()
	opts.SetSort(bson.D{{"created", -1}})

	setElements := bson.D{{"created", bson.M{"$gte": getDay(25)}}}

	sortCursor, err := collection.Find(ctx, setElements, opts)
	if err != nil {
		return nil, err
	}

	if err = sortCursor.All(ctx, &notifications); err != nil {
		return nil, err
	}

	return notifications, nil
}
