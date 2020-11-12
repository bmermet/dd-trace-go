// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2016-2020 Datadog, Inc.

package redigo // import "github.com/bmermet/dd-trace-go/contrib/gomodule/redigo"

import (
	"math"

	"github.com/bmermet/dd-trace-go/internal"
)

type dialConfig struct {
	serviceName   string
	analyticsRate float64
}

// DialOption represents an option that can be passed to Dial.
type DialOption func(*dialConfig)

func defaults(cfg *dialConfig) {
	cfg.serviceName = "redis.conn"
	// cfg.analyticsRate = globalconfig.AnalyticsRate()
	if internal.BoolEnv("DD_TRACE_REDIGO_ANALYTICS_ENABLED", false) {
		cfg.analyticsRate = 1.0
	} else {
		cfg.analyticsRate = math.NaN()
	}
}

// WithServiceName sets the given service name for the dialled connection.
func WithServiceName(name string) DialOption {
	return func(cfg *dialConfig) {
		cfg.serviceName = name
	}
}

// WithAnalytics enables Trace Analytics for all started spans.
func WithAnalytics(on bool) DialOption {
	return func(cfg *dialConfig) {
		if on {
			cfg.analyticsRate = 1.0
		} else {
			cfg.analyticsRate = math.NaN()
		}
	}
}

// WithAnalyticsRate sets the sampling rate for Trace Analytics events
// correlated to started spans.
func WithAnalyticsRate(rate float64) DialOption {
	return func(cfg *dialConfig) {
		if rate >= 0.0 && rate <= 1.0 {
			cfg.analyticsRate = rate
		} else {
			cfg.analyticsRate = math.NaN()
		}
	}
}
