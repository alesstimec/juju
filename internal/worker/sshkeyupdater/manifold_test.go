// Copyright 2026 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package sshkeyupdater

import (
	"testing"

	"github.com/juju/tc"
	"github.com/juju/worker/v5"

	"github.com/juju/juju/api/base"
	coreerrors "github.com/juju/juju/core/errors"
	coressh "github.com/juju/juju/core/ssh"
	loggertesting "github.com/juju/juju/internal/logger/testing"
	"github.com/juju/juju/internal/testhelpers"
)

type manifoldSuite struct {
	testhelpers.IsolationSuite
}

func TestManifoldSuite(t *testing.T) {
	tc.Run(t, &manifoldSuite{})
}

func (s *manifoldSuite) newConfig(c *tc.C, modifier func(*ManifoldConfig)) ManifoldConfig {
	cfg := ManifoldConfig{
		AgentName:       "agent",
		APICallerName:   "api-caller",
		Logger:          loggertesting.WrapCheckLog(c),
		NewWorker:       func(WorkerConfig) (worker.Worker, error) { return nil, nil },
		NewFacadeClient: func(base.APICaller) FacadeClient { return nil },
	}
	modifier(&cfg)
	return cfg
}

func (s *manifoldSuite) TestValidate(c *tc.C) {
	c.Check(s.newConfig(c, func(*ManifoldConfig) {}).Validate(), tc.ErrorIsNil)

	c.Check(s.newConfig(c, func(cfg *ManifoldConfig) { cfg.AgentName = "" }).Validate(), tc.ErrorIs, coreerrors.NotValid)
	c.Check(s.newConfig(c, func(cfg *ManifoldConfig) { cfg.APICallerName = "" }).Validate(), tc.ErrorIs, coreerrors.NotValid)
	c.Check(s.newConfig(c, func(cfg *ManifoldConfig) { cfg.Logger = nil }).Validate(), tc.ErrorIs, coreerrors.NotValid)
	c.Check(s.newConfig(c, func(cfg *ManifoldConfig) { cfg.NewWorker = nil }).Validate(), tc.ErrorIs, coreerrors.NotValid)
	c.Check(s.newConfig(c, func(cfg *ManifoldConfig) { cfg.NewFacadeClient = nil }).Validate(), tc.ErrorIs, coreerrors.NotValid)
}

func (s *manifoldSuite) TestInputsAndOutput(c *tc.C) {
	manifold := Manifold(s.newConfig(c, func(*ManifoldConfig) {}))
	c.Check(manifold.Inputs, tc.SameContents, []string{"agent", "api-caller"})
	c.Check(manifold.Output, tc.NotNil)
}

func (s *manifoldSuite) TestOutputExtractsEphemeralKeysUpdater(c *tc.C) {
	in := &AuthWorker{}

	var updater coressh.EphemeralKeysUpdater
	err := Output(in, &updater)
	c.Assert(err, tc.ErrorIsNil)
	c.Check(updater, tc.Equals, in)
}

func (s *manifoldSuite) TestOutputWrongInputType(c *tc.C) {
	var updater coressh.EphemeralKeysUpdater
	err := Output(&stubWorker{}, &updater)
	c.Check(err, tc.ErrorMatches, `expected \*AuthWorker, got .*`)
}

func (s *manifoldSuite) TestOutputWrongOutputType(c *tc.C) {
	in := &AuthWorker{}

	var wrong worker.Worker
	err := Output(in, &wrong)
	c.Check(err, tc.ErrorMatches, `expected \*coressh.EphemeralKeysUpdater, got .*`)
}

type stubWorker struct{}

func (*stubWorker) Kill()       {}
func (*stubWorker) Wait() error { return nil }
