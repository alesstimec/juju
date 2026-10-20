// Copyright 2013 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package keyupdater

import (
	"context"

	"github.com/juju/names/v6"

	"github.com/juju/juju/api/base"
	apiservererrors "github.com/juju/juju/apiserver/errors"
	"github.com/juju/juju/internal/errors"
	"github.com/juju/juju/rpc/params"
)

// Option is a function that can be used to configure a Client.
type Option = base.Option

// WithTracer returns an Option that configures the Client to use the
// supplied tracer.
var WithTracer = base.WithTracer

// Client provides access to the KeyUpdater agent facade.
type Client struct {
	facade base.FacadeCaller
}

// NewClient returns a new KeyUpdater facade client.
func NewClient(caller base.APICaller, options ...Option) *Client {
	return &Client{facade: base.NewFacadeCaller(caller, "KeyUpdater", options...)}
}

// AuthorisedKeys returns the controller's authorised ssh keys that must be
// present on the machine identified by tag.
func (c *Client) AuthorisedKeys(ctx context.Context, tag names.MachineTag) ([]string, error) {
	var results params.StringsResults
	args := params.Entities{
		Entities: []params.Entity{{Tag: tag.String()}},
	}
	if err := c.facade.FacadeCall(ctx, "AuthorisedKeys", args, &results); err != nil {
		return nil, errors.Capture(err)
	}
	if len(results.Results) != 1 {
		return nil, errors.Errorf("expected 1 result, got %d", len(results.Results))
	}
	result := results.Results[0]
	if result.Error != nil {
		return nil, apiservererrors.RestoreError(result.Error)
	}
	return result.Result, nil
}
