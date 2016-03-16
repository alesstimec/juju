// Copyright 2016 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package rsyslogger_test

import (
	"crypto/tls"
	"time"

	"github.com/juju/names"
	"github.com/juju/syslog"
	jtesting "github.com/juju/testing"
	jc "github.com/juju/testing/checkers"
	gc "gopkg.in/check.v1"

	"github.com/juju/juju/agent"
	"github.com/juju/juju/api"
	"github.com/juju/juju/api/logreader"
	"github.com/juju/juju/juju/testing"
	"github.com/juju/juju/state"
	"github.com/juju/juju/testing/factory"
	"github.com/juju/juju/worker"
	"github.com/juju/juju/worker/rsyslogger"
)

// worstCase is used for timeouts when timing out
// will fail the test. Raising this value should
// not affect the overall running time of the tests
// unless they fail.
const worstCase = 5 * time.Second

type RsysloggerSuite struct {
	testing.JujuConnSuite

	reader  *logreader.API
	machine *state.Machine
}

var _ = gc.Suite(&RsysloggerSuite{})

func (s *RsysloggerSuite) SetUpTest(c *gc.C) {
	s.JujuConnSuite.SetUpTest(c)

	nonce := "some-nonce"
	machine, password := s.Factory.MakeMachineReturningPassword(c,
		&factory.MachineParams{Nonce: nonce})
	apiInfo := s.APIInfo(c)
	apiInfo.Tag = machine.Tag()
	apiInfo.Password = password
	apiInfo.Nonce = nonce
	apiConn, err := api.Open(apiInfo, api.DefaultDialOpts())
	c.Assert(err, jc.ErrorIsNil)

	s.reader = logreader.NewAPI(apiConn)
	c.Assert(s.reader, gc.NotNil)
	s.machine = machine
}

type mockConfig struct {
	agent.Config
	c   *gc.C
	tag names.Tag
}

func (mock *mockConfig) Tag() names.Tag {
	return mock.tag
}

func agentConfig(c *gc.C, tag names.Tag) *mockConfig {
	return &mockConfig{c: c, tag: tag}
}

func (s *RsysloggerSuite) makeLogger(c *gc.C) (worker.Worker, *mockConfig) {
	config := agentConfig(c, s.machine.Tag())
	w, err := rsyslogger.NewRsysWorker(s.reader, config)
	c.Assert(err, jc.ErrorIsNil)
	return w, config
}

func (s *RsysloggerSuite) TestRunStop(c *gc.C) {
	loggingWorker, _ := s.makeLogger(c)
	c.Assert(worker.Stop(loggingWorker), gc.IsNil)
}

func (s *RsysloggerSuite) TestRunWorker(c *gc.C) {
	writer := &mockSyslogger{}
	cleanup := jtesting.PatchValue(
		rsyslogger.DialSyslog,
		func(network, raddr string,
			priority syslog.Priority,
			tag string,
			tlsCfg *tls.Config,
		) (rsyslogger.SysLogger, error) {
			return writer, nil
		})
	defer cleanup()
}

type mockSyslogger struct {
	jtesting.Stub
}

func (m *mockSyslogger) Crit(msg string) error {
	m.AddCall("Crit", msg)
	return m.NextErr()
}

func (m *mockSyslogger) Err(msg string) error {
	m.AddCall("Err", msg)
	return m.NextErr()
}

func (m *mockSyslogger) Warning(msg string) error {
	m.AddCall("Warning", msg)
	return m.NextErr()
}

func (m *mockSyslogger) Notice(msg string) error {
	m.AddCall("Notice", msg)
	return m.NextErr()
}

func (m *mockSyslogger) Info(msg string) error {
	m.AddCall("Info", msg)
	return m.NextErr()
}

func (m *mockSyslogger) Debug(msg string) error {
	m.AddCall("Debug", msg)
	return m.NextErr()
}

func (m *mockSyslogger) Write(msg []byte) (int, error) {
	m.AddCall("Write", msg)
	return len(msg), m.NextErr()
}
