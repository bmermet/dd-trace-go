// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2016-2020 Datadog, Inc.

// Package httprouter provides functions to trace the julienschmidt/httprouter package (https://github.com/julienschmidt/httprouter).
package httprouter // import "github.com/bmermet/dd-trace-go/contrib/julienschmidt/httprouter"

import (
	"math"
	"net/http"
	"strings"

	"github.com/bmermet/dd-trace-go/contrib/internal/httputil"
	"github.com/bmermet/dd-trace-go/ddtrace/ext"
	"github.com/bmermet/dd-trace-go/ddtrace/tracer"

	"github.com/julienschmidt/httprouter"
)

// Router is a traced version of httprouter.Router.
type Router struct {
	*httprouter.Router
	config *routerConfig
}

// New returns a new router augmented with tracing.
func New(opts ...RouterOption) *Router {
	cfg := new(routerConfig)
	defaults(cfg)
	for _, fn := range opts {
		fn(cfg)
	}
	if !math.IsNaN(cfg.analyticsRate) {
		cfg.spanOpts = append(cfg.spanOpts, tracer.Tag(ext.EventSampleRate, cfg.analyticsRate))
	}
	cfg.spanOpts = append(cfg.spanOpts, tracer.Measured())
	return &Router{httprouter.New(), cfg}
}

// ServeHTTP implements http.Handler.
func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	// get the resource associated to this request
	route := req.URL.Path
	_, ps, _ := r.Router.Lookup(req.Method, route)
	for _, param := range ps {
		route = strings.Replace(route, param.Value, ":"+param.Key, 1)
	}
	resource := req.Method + " " + route
	httputil.TraceAndServe(r.Router, w, req, r.config.serviceName, resource, nil, r.config.spanOpts...)
}
