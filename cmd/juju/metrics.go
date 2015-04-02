// Copyright 2015 Canonical Ltd. All rights reserved.

package main

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/juju/errors"
	"gopkg.in/macaroon-bakery.v0/httpbakery"
)

// TODO (alesstimec, mattyw): Add a registration uuid to the metricsRegistrarFunc.
type metricsRegistrarFunc func(environmentUUID, charmURL, serviceName string, client *http.Client, visitWebPage func(*url.URL) error) ([]byte, error)

var (
	registerMetrics    metricsRegistrarFunc = nilMetricsRegistrar
	registerMetricsURL                      = ""

	_ metricsRegistrarFunc = nilMetricsRegistrar
	_ metricsRegistrarFunc = httpMetricsRegistrar
)

func nilMetricsRegistrar(_, _, _ string, _ *http.Client, _ func(*url.URL) error) ([]byte, error) {
	return []byte{}, nil
}

type metricRegistrationPost struct {
	EnvironmentUUID string `json:"env-uuid"`
	CharmURL        string `json:"charm-url"`
	ServiceName     string `json:"service-name"`
}

func httpMetricsRegistrar(environmentUUID, charmURL, serviceName string, client *http.Client, visitWebPage func(*url.URL) error) ([]byte, error) {
	registerURL, err := url.Parse(registerMetricsURL)
	if err != nil {
		return nil, errors.Trace(err)
	}

	registrationPost := metricRegistrationPost{
		EnvironmentUUID: environmentUUID,
		CharmURL:        charmURL,
		ServiceName:     serviceName,
	}

	buff := &bytes.Buffer{}
	encoder := json.NewEncoder(buff)
	err = encoder.Encode(registrationPost)
	if err != nil {
		return nil, errors.Trace(err)
	}

	req, err := http.NewRequest("POST", registerURL.String(), nil)
	if err != nil {
		return nil, errors.Trace(err)
	}
	req.Header.Set("Content-Type", "application/json")

	bodyGetter := httpbakery.SeekerBody(bytes.NewReader(buff.Bytes()))

	response, err := httpbakery.DoWithBody(client, req, bodyGetter, visitWebPage)
	if err != nil {
		return nil, errors.Trace(err)
	}
	defer discardClose(response)

	if response.StatusCode != http.StatusOK {
		return nil, errors.Errorf("failed to register metrics: http response is %d", response.StatusCode)
	}

	data, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, errors.Trace(err)
	}

	return data, nil
}

// discardClose reads any remaining data from the response body and closes it.
func discardClose(response *http.Response) {
	if response == nil || response.Body == nil {
		return
	}
	io.Copy(ioutil.Discard, response.Body)
	response.Body.Close()
}
