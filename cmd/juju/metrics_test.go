// Copyright 2015 Canonical Ltd. All rights reserved.

package main

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"

	jujutesting "github.com/juju/testing"
	gc "gopkg.in/check.v1"
)

var _ = gc.Suite(&metricsSuite{})

type metricsSuite struct {
	handler *testMetricsRegistrationHandler
	server  *httptest.Server
}

func (s *metricsSuite) SetUpTest(c *gc.C) {
	s.handler = &testMetricsRegistrationHandler{}
	s.server = httptest.NewServer(s.handler)
}

func (s *metricsSuite) TearDownTest(c *gc.C) {
	s.server.Close()
}

func (s *metricsSuite) TestNilMetricsRegistrar(c *gc.C) {
	data, err := nilMetricsRegistrar("registration uuid", "environment uuid", "charm url", "service name", &http.Client{}, func(*url.URL) error { return nil })
	c.Assert(err, gc.IsNil)
	c.Assert(data, gc.DeepEquals, []byte{})
}

func (s *metricsSuite) TestHttpMetricsRegistrar(c *gc.C) {
	cleanup := jujutesting.PatchValue(&registerMetricsURL, s.server.URL)
	defer cleanup()

	data, err := httpMetricsRegistrar("registration uuid", "environment uuid", "charm url", "service name", &http.Client{}, func(*url.URL) error { return nil })
	c.Assert(err, gc.IsNil)
	c.Assert(data, gc.DeepEquals, []byte("hello metrics"))
	c.Assert(s.handler.registrationCalls, gc.HasLen, 1)
	c.Assert(s.handler.registrationCalls[0], gc.DeepEquals, metricRegistrationPost{RegistrationUUID: "registration uuid", EnvironmentUUID: "environment uuid", CharmURL: "charm url", ServiceName: "service name"})
}

type testMetricsRegistrationHandler struct {
	registrationCalls []metricRegistrationPost
}

// ServeHTTP implements http.Handler.
func (c *testMetricsRegistrationHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	if req.Method != "POST" {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var registrationPost metricRegistrationPost
	decoder := json.NewDecoder(req.Body)
	err := decoder.Decode(&registrationPost)
	if err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	data := bytes.NewBuffer([]byte("hello metrics"))
	io.Copy(w, data)

	c.registrationCalls = append(c.registrationCalls, registrationPost)
}
