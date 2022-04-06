package mongo

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Notification struct {
	ID                      primitive.ObjectID `json:"_id" bson:"_id"`
	SiteID                  int                `json:"site_id" bson:"site_id,omitempty"`
	SiteSlug                string             `json:"site_slug" bson:"site_slug,omitempty"`
	UserID                  int                `json:"user_id" bson:"user_id,omitempty"`
	TaskID                  string             `json:"task_id" bson:"task_id,omitempty"`
	DistributorID           string             `json:"distributor_id" bson:"distributor_id,omitempty"`
	MarketingDrip           string             `json:"marketing_drip" bson:"marketing_drip,omitempty"`
	NotificationSlug        string             `json:"notification_slug" bson:"notification_slug,omitempty"`
	Success                 bool               `json:"success" bson:"success,omitempty"`
	Disabled                bool               `json:"disabled" bson:"disabled,omitempty"`
	Created                 time.Time          `json:"created" bson:"created,omitempty"`
	NotificationAPIResponse interface{}        `json:"notification_api_response" bson:"notification_api_response,omitempty"`
}
