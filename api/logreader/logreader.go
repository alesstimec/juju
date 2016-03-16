// Copyright 2016 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

// Package logreader implements the API for
// retrieving log messages from the API server.
package logreader

import (
	"fmt"
	"io"
	"net/url"

	"github.com/juju/errors"
	"github.com/juju/names"
	"launchpad.net/tomb"

	"github.com/juju/juju/api/base"
	apiwatcher "github.com/juju/juju/api/watcher"
	"github.com/juju/juju/apiserver/params"
	"github.com/juju/juju/watcher"
)

// LogReader is the interface that allows reading
// log messages transmitted by the server.
type LogReader interface {
	// Read returns a channel that can be used to receive log
	// messages.
	ReadLogs() chan params.LogRecordResult

	io.Closer
}

// API provides access to the LogReader API.
type API struct {
	facade    base.FacadeCaller
	connector base.StreamConnector
}

// NewAPI creates a new client-side logsender API.
func NewAPI(api base.APICaller) *API {
	facade := base.NewFacadeCaller(api, "RsyslogConfig")
	return &API{facade: facade, connector: api}
}

func (api *API) WatchRsyslogConfig(agentTag names.Tag) (watcher.NotifyWatcher, error) {
	var results params.NotifyWatchResults
	args := params.Entities{
		Entities: []params.Entity{{Tag: agentTag.String()}},
	}
	err := api.facade.FacadeCall("WatchRsyslogConfig", args, &results)
	if err != nil {
		// TODO: Not directly tested
		return nil, err
	}
	if len(results.Results) != 1 {
		// TODO: Not directly tested
		return nil, errors.Errorf("expected 1 result, got %d", len(results.Results))
	}
	result := results.Results[0]
	if result.Error != nil {
		//  TODO: Not directly tested
		return nil, result.Error
	}
	w := apiwatcher.NewNotifyWatcher(api.facade.RawAPICaller(), result)
	return w, nil
}

func (api *API) RsyslogURLConfig(agentTag names.Tag) (string, error) {
	var results params.StringResults
	args := params.Entities{
		Entities: []params.Entity{{Tag: agentTag.String()}},
	}
	err := api.facade.FacadeCall("RsyslogURLConfig", args, &results)
	if err != nil {
		// TODO: Not directly tested
		return "", err
	}
	if len(results.Results) != 1 {
		// TODO: Not directly tested
		return "", errors.Errorf("expected 1 result, got %d", len(results.Results))
	}
	result := results.Results[0]
	if err := result.Error; err != nil {
		return "", err
	}
	return result.Result, nil
}

func (api *API) RsyslogCACertConfig(agentTag names.Tag) (string, error) {
	var results params.StringResults
	args := params.Entities{
		Entities: []params.Entity{{Tag: agentTag.String()}},
	}
	err := api.facade.FacadeCall("RsyslogCACertConfig", args, &results)
	if err != nil {
		// TODO: Not directly tested
		return "", err
	}
	if len(results.Results) != 1 {
		// TODO: Not directly tested
		return "", errors.Errorf("expected 1 result, got %d", len(results.Results))
	}
	result := results.Results[0]
	if err := result.Error; err != nil {
		return "", err
	}
	return result.Result, nil
}

// LogReader returns a  structure that implements
// the LogReader interface,
// which must be closed when finished with.
func (api *API) LogReader() (LogReader, error) {
	conn, err := api.connector.ConnectStream("/log", url.Values{"jsonFormat": []string{"true"}})
	if err != nil {
		return nil, errors.Annotatef(err, "cannot connect to /log")
	}
	return &reader{conn: conn}, nil
}

type reader struct {
	conn base.Stream
	tomb tomb.Tomb
}

func (r *reader) ReadLogs() chan params.LogRecordResult {
	channel := make(chan params.LogRecordResult, 1)
	go func() {
		defer r.tomb.Done()
		r.tomb.Kill(r.loop(channel))
	}()
	return channel
}

func (r *reader) loop(channel chan params.LogRecordResult) error {
	for {
		var record params.LogRecordResult
		err := r.conn.ReadJSON(&record.LogRecord)
		if err != nil {
			record.Error = &params.Error{Message: fmt.Sprintf("failed to read JSON: %v", err.Error())}
		}
		select {
		case <-r.tomb.Dying():
			return tomb.ErrDying
		case channel <- record:
		default:
			return errors.Errorf("failed to send log record")
		}
	}
}

func (r *reader) Close() error {
	r.tomb.Kill(nil)
	return r.conn.Close()
}
