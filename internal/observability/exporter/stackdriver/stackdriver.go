//
// Copyright (C) 2019-2020 Vdaas.org Vald team ( kpango, rinx, kmrmt )
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

// Package stackdriver provides a stackdriver exporter.
package stackdriver

import (
	"context"

	traceapi "cloud.google.com/go/trace/apiv2"
	"contrib.go.opencensus.io/exporter/stackdriver"
	"go.opencensus.io/trace"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
)

type Stackdriver interface {
	Start(ctx context.Context) error
	Stop(ctx context.Context)
}

type exporter struct {
	exporter *stackdriver.Exporter

	monitoringEnabled bool
	tracingEnabled    bool

	*stackdriver.Options
}

func New(opts ...Option) (s Stackdriver, err error) {
	e := new(exporter)
	e.Options = new(stackdriver.Options)

	for _, opt := range append(defaultOpts, opts...) {
		err = opt(e)
		if err != nil {
			return nil, err
		}
	}

	return e, nil
}

func (e *exporter) Start(ctx context.Context) (err error) {
	e.Options.Context = ctx

	// TODO: debugging use
	creds, err := google.FindDefaultCredentials(ctx, traceapi.DefaultAuthScopes()...)
	if err != nil {
		return err
	}

	err = WithMonitoringClientOptions(option.WithCredentials(creds))(e)
	if err != nil {
		return err
	}

	err = WithTraceClientOptions(option.WithCredentials(creds))(e)
	if err != nil {
		return err
	}

	e.exporter, err = stackdriver.NewExporter(*e.Options)
	if err != nil {
		return err
	}

	if e.monitoringEnabled {
		err = e.exporter.StartMetricsExporter()
		if err != nil {
			return err
		}
	}

	if e.tracingEnabled {
		trace.RegisterExporter(e.exporter)
	}

	return nil
}

func (e *exporter) Stop(ctx context.Context) {
	if e.exporter != nil {
		if e.monitoringEnabled {
			e.exporter.StopMetricsExporter()
		}

		e.exporter.Flush()
	}
}
