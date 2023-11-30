package utils

import (
	"github.com/newrelic/go-agent/v3/newrelic"
)

type NewRelic struct {
	config Config
	App    *newrelic.Application
}

//func NewRelicInstance() *NewRelic {
//
//	config := NewConfig()
//
//	app, err := newrelic.NewApplication(
//		newrelic.ConfigAppName(config.NewRelicAppName()),
//		newrelic.ConfigLicense(config.NewRelicLicense()),
//		newrelic.ConfigAppLogForwardingEnabled(true),
//	)
//
//	if err != nil {
//		().Warn("Failed to initialize new relic monitor: " + err.Error())
//	}
//
//	return &NewRelic{
//		config: config,
//		App:    app,
//	}
//}
