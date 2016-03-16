// Copyright 2016 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package rsyslogger

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"time"

	"github.com/juju/errors"
	"github.com/juju/loggo"
	"github.com/juju/syslog"

	"github.com/juju/juju/agent"
	"github.com/juju/juju/api/logreader"
	"github.com/juju/juju/apiserver/params"
	"github.com/juju/juju/worker"
)

var (
	logger     = loggo.GetLogger("juju.worker.rSysLogger")
	dialSyslog = func(network, raddr string, priority syslog.Priority, tag string, tlsCfg *tls.Config) (SysLogger, error) {
		return syslog.Dial(network, raddr, priority, tag, tlsCfg)
	}

	longWait = time.Second
)

type SysLogger interface {
	Crit(string) error
	Err(string) error
	Warning(string) error
	Notice(string) error
	Info(string) error
	Debug(string) error
	Write([]byte) (int, error)
}

// NewRsysWorker returns a worker.Worker that uses the notify watcher returned
// from the setup.
func NewRsysWorker(api *logreader.API, agentConfig agent.Config) (worker.Worker, error) {
	logger := &rsysLogger{
		api:         api,
		agentConfig: agentConfig,
	}
	return worker.NewSimpleWorker(logger.loop), nil
}

// RsysLogger uses the api/logreader to tail the log
// collection and forwards all log messages to a
// rsyslog server.
type rsysLogger struct {
	api         *logreader.API
	agentConfig agent.Config
	writer      SysLogger
}

func (r *rsysLogger) loop(stop <-chan struct{}) error {
	configWatcher, err := r.api.WatchRsyslogConfig(r.agentConfig.Tag())
	if err != nil {
		return errors.Trace(err)
	}

	stopChannel := make(chan struct{})
	runningLogger := false
	for {
		select {
		case <-stop:
			select {
			case stopChannel <- struct{}{}:
			default:
			}
			return nil
		case <-configWatcher.Changes():
			if runningLogger {
				select {
				case stopChannel <- struct{}{}:
				case <-time.After(longWait):
					logger.Errorf("failed to stop the read loop: time out")
				}
				runningLogger = false
			}
			writer, err := newSysLogger(r.api, r.agentConfig)
			if err == nil {
				go r.readLoop(writer, stopChannel)
				runningLogger = true
			} else {
				logger.Errorf("failed to create a syslog writer: %v", err)
			}
		}
	}
}

func (r *rsysLogger) readLoop(writer SysLogger, stop <-chan struct{}) error {
	logReader, err := r.api.LogReader()
	if err != nil {
		return errors.Trace(err)
	}
	defer logReader.Close()

	logChannel := logReader.ReadLogs()
	for {
		select {
		case <-stop:
			return nil
		case rec := <-logChannel:
			if rec.Error != nil {
				logger.Errorf("received error: %v", rec.Error.Message)
				continue
			}
			var err error
			switch rec.Level {
			case loggo.CRITICAL:
				err = r.writer.Crit(formatLogRecord(rec))
			case loggo.ERROR:
				err = r.writer.Err(formatLogRecord(rec))
			case loggo.WARNING:
				err = r.writer.Warning(formatLogRecord(rec))
			case loggo.INFO:
				err = r.writer.Info(formatLogRecord(rec))
			case loggo.DEBUG:
				err = r.writer.Debug(formatLogRecord(rec))
			default:
			}
			if err != nil {
				logger.Errorf("failed to forward the log entry: %v", err)
			}
		}
	}
}

func formatLogRecord(r params.LogRecordResult) string {
	return fmt.Sprintf("%s: %s %s %s %s %s\n",
		r.ModelUUID,
		formatTime(r.Time),
		r.Level.String(),
		r.Module,
		r.Location,
		r.Message,
	)
}

func formatTime(t time.Time) string {
	return t.In(time.UTC).Format("2006-01-02 15:04:05")
}

func newSysLogger(api *logreader.API, agentConfig agent.Config) (SysLogger, error) {
	url, err := api.RsyslogURLConfig(agentConfig.Tag())
	if err != nil {
		return nil, errors.Trace(err)
	}

	caCert, err := api.RsyslogCACertConfig(agentConfig.Tag())
	if err != nil {
		return nil, errors.Trace(err)
	}

	if url != "" && caCert != "" {
		caCertPool := x509.NewCertPool()
		ok := caCertPool.AppendCertsFromPEM([]byte(caCert))
		if !ok {
			return nil, errors.Errorf("failed to parse the ca certificate")
		}

		writer, err := dialSyslog("tcp", url, syslog.LOG_LOCAL0, "juju-syslog", &tls.Config{RootCAs: caCertPool})
		if err != nil {
			return nil, errors.Trace(err)
		}
		return writer, nil
	}
	return nil, errors.Errorf("configured: url %v, ca cert %v", url != "", caCert != "")
}
