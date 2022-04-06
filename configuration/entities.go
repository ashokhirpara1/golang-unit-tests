package configuration

import (
	ld "gopkg.in/launchdarkly/go-server-sdk.v5"
)

// Struct for all clients
type Clients struct {
	LD     LD
	Sentry Sentry
}

// Struct for Sentry
type Sentry struct {
	Config SentryConfig
}

// Struct for Launch Darkly
type LD struct {
	LDClient *ld.LDClient
	Config   LDConfig
}

// Struct for all the Launch Darkly Configurations
type LDConfig struct {
	LD_SDK_KEY string
	Flags      struct {
		SENTRY     string
		ROUTE_FLAG string
	}
}

// Struct for all the Sentry Configurations
type SentryConfig struct {
	SENTRY_DSN string
}
