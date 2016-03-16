// Copyright 2016 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package rsyslogger

import (
	"github.com/juju/loggo"
	"github.com/juju/names"

	"github.com/juju/juju/apiserver/common"
	"github.com/juju/juju/apiserver/params"
	"github.com/juju/juju/state"
	"github.com/juju/juju/state/watcher"
)

var logger = loggo.GetLogger("juju.apiserver.rsyslogger")

func init() {
	common.RegisterStandardFacade("RsyslogConfig", 1, NewRsyslogConfigAPI)
}

// RsyslogConfigWatcher defines the methods on the logger API end point.
type RsyslogConfigWatcher interface {
	WatchRsyslogConfig(args params.Entities) params.NotifyWatchResults
	RsyslogURLConfig(args params.Entities) params.StringResults
	RsyslogCACertConfig(args params.Entities) params.StringResults
}

// RsyslogConfigAPI implements the RsyslogConfigWatcher interface and is the concrete
// implementation of the api end point.
type RsyslogConfigAPI struct {
	state      *state.State
	resources  *common.Resources
	authorizer common.Authorizer
}

var _ RsyslogConfigWatcher = (*RsyslogConfigAPI)(nil)

// NewRsyslogConfigWatcher creates a new server-side rsyslogger API end point.
func NewRsyslogConfigAPI(
	st *state.State,
	resources *common.Resources,
	authorizer common.Authorizer,
) (*RsyslogConfigAPI, error) {
	if !authorizer.AuthMachineAgent() {
		return nil, common.ErrPerm
	}
	return &RsyslogConfigAPI{state: st, resources: resources, authorizer: authorizer}, nil
}

// WatchRsyslogConfig starts a watcher to track changes to the rsyslog config
// for the agents specified.  Unfortunately the current infrastructure makes
// watching parts of the config non-trivial, so currently any change to the
// config will cause the watcher to notify the client.
func (api *RsyslogConfigAPI) WatchRsyslogConfig(arg params.Entities) params.NotifyWatchResults {
	result := make([]params.NotifyWatchResult, len(arg.Entities))
	for i, entity := range arg.Entities {
		tag, err := names.ParseTag(entity.Tag)
		if err != nil {
			result[i].Error = common.ServerError(err)
			continue
		}
		err = common.ErrPerm
		if api.authorizer.AuthOwner(tag) {
			watch := api.state.WatchForModelConfigChanges()
			// Consume the initial event. Technically, API calls to Watch
			// 'transmit' the initial event in the Watch response. But
			// NotifyWatchers have no state to transmit.
			if _, ok := <-watch.Changes(); ok {
				result[i].NotifyWatcherId = api.resources.Register(watch)
				err = nil
			} else {
				err = watcher.EnsureErr(watch)
			}
		}
		result[i].Error = common.ServerError(err)
	}
	return params.NotifyWatchResults{Results: result}
}

// RsyslogURLConfig reports the rsyslog url config for the agents specified.
func (api *RsyslogConfigAPI) RsyslogURLConfig(arg params.Entities) params.StringResults {
	if len(arg.Entities) == 0 {
		return params.StringResults{}
	}
	results := make([]params.StringResult, len(arg.Entities))
	config, configErr := api.state.ModelConfig()
	url, urlOK := config.RsyslogURL()
	for i, entity := range arg.Entities {
		tag, err := names.ParseTag(entity.Tag)
		if err != nil {
			results[i].Error = common.ServerError(err)
			continue
		}
		err = common.ErrPerm
		if api.authorizer.AuthOwner(tag) {
			if configErr == nil && urlOK {
				results[i].Result = url
				err = nil
			} else if !urlOK {
				results[i].Result = ""
				err = nil
			} else {
				err = configErr
			}
		}
		results[i].Error = common.ServerError(err)
	}
	return params.StringResults{Results: results}
}

// RsyslogCACertConfig reports the rsyslog ca cert config for the agents specified.
func (api *RsyslogConfigAPI) RsyslogCACertConfig(arg params.Entities) params.StringResults {
	if len(arg.Entities) == 0 {
		return params.StringResults{}
	}
	results := make([]params.StringResult, len(arg.Entities))
	config, configErr := api.state.ModelConfig()
	caCert, certOK := config.RsyslogCACert()
	for i, entity := range arg.Entities {
		tag, err := names.ParseTag(entity.Tag)
		if err != nil {
			results[i].Error = common.ServerError(err)
			continue
		}
		err = common.ErrPerm
		if api.authorizer.AuthOwner(tag) {
			if configErr == nil && certOK {
				results[i].Result = caCert
				err = nil
			} else if !certOK {
				results[i].Result = ""
				err = nil
			} else {
				err = configErr
			}
		}
		results[i].Error = common.ServerError(err)
	}
	return params.StringResults{Results: results}
}
