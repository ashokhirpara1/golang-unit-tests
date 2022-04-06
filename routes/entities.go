package routes

// Query parameter options to filter records
type Inputs struct {
	Type             string `json:"type"` // marketing_drip or pulse
	SiteID           int    `json:"site_id"`
	UserID           int    `json:"user_id,omitempty"`        // for marketing_drip
	DistributorID    string `json:"distributor_id,omitempty"` // for pulse
	NotificationSlug string `json:"notification_slug,omitempty"`
	MarketinDrip     string `json:"marketing_drip,omitempty"` // type/category of marketing drip notifications
	Days             int    `json:"days,omitempty"`           // default 0 >= today
	Success          int    `json:"success,omitempty"`
	Disabled         int    `json:"disabled,omitempty"`
}
