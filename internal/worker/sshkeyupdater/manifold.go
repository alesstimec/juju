// Copyright 2015 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package sshkeyupdater

import (
	"context"

	"github.com/juju/names/v6"
	"github.com/juju/worker/v5"
	"github.com/juju/worker/v5/dependency"

	"github.com/juju/juju/agent"
	"github.com/juju/juju/api/agent/keyupdater"
	"github.com/juju/juju/api/base"
	coreerrors "github.com/juju/juju/core/errors"
	"github.com/juju/juju/core/logger"
	coressh "github.com/juju/juju/core/ssh"
	"github.com/juju/juju/internal/errors"
)

// ManifoldConfig holds the information necessary to run the sshkeyupdater
// worker in a dependency.Engine.
type ManifoldConfig struct {
	// AgentName is the agent dependency name.
	AgentName string
	// APICallerName is the api caller dependency name.
	APICallerName string
	// Logger is the logger to use for the worker.
	Logger logger.Logger
	// NewWorker creates a new sshkeyupdater worker.
	NewWorker func(WorkerConfig) (worker.Worker, error)
	// NewFacadeClient creates the facade client from an API caller.
	NewFacadeClient func(base.APICaller) FacadeClient
}

// Validate validates the manifold configuration.
func (config ManifoldConfig) Validate() error {
	if config.AgentName == "" {
		return errors.Errorf("empty AgentName").Add(coreerrors.NotValid)
	}
	if config.APICallerName == "" {
		return errors.Errorf("empty APICallerName").Add(coreerrors.NotValid)
	}
	if config.Logger == nil {
		return errors.Errorf("nil Logger").Add(coreerrors.NotValid)
	}
	if config.NewWorker == nil {
		return errors.Errorf("nil NewWorker").Add(coreerrors.NotValid)
	}
	if config.NewFacadeClient == nil {
		return errors.Errorf("nil NewFacadeClient").Add(coreerrors.NotValid)
	}
	return nil
}

// Manifold returns a dependency manifold that runs the sshkeyupdater worker.
// The worker's EphemeralKeysUpdater is exposed as the manifold output so that
// the sshsession worker can inject and remove ephemeral keys for the lifetime
// of a reverse tunnel.
func Manifold(config ManifoldConfig) dependency.Manifold {
	return dependency.Manifold{
		Inputs: []string{
			config.AgentName,
			config.APICallerName,
		},
		Start:  config.start,
		Output: Output,
	}
}

func (config ManifoldConfig) start(ctx context.Context, getter dependency.Getter) (worker.Worker, error) {
	if err := config.Validate(); err != nil {
		return nil, errors.Capture(err)
	}

	var thisAgent agent.Agent
	if err := getter.Get(config.AgentName, &thisAgent); err != nil {
		return nil, errors.Capture(err)
	}

	machineTag, ok := thisAgent.CurrentConfig().Tag().(names.MachineTag)
	if !ok {
		return nil, errors.Errorf("expected a machine tag, got %T", thisAgent.CurrentConfig().Tag())
	}

	var apiCaller base.APICaller
	if err := getter.Get(config.APICallerName, &apiCaller); err != nil {
		return nil, errors.Capture(err)
	}

	w, err := config.NewWorker(WorkerConfig{
		Logger:       config.Logger,
		MachineTag:   machineTag,
		FacadeClient: config.NewFacadeClient(apiCaller),
	})
	if err != nil {
		return nil, errors.Errorf("cannot start ssh key updater worker: %w", err)
	}
	return w, nil
}

// NewFacadeClient returns a FacadeClient backed by the KeyUpdater facade.
func NewFacadeClient(caller base.APICaller) FacadeClient {
	return keyupdater.NewClient(caller)
}

// Output extracts an EphemeralKeysUpdater from the running AuthWorker.
func Output(in worker.Worker, out any) error {
	w, ok := in.(*AuthWorker)
	if !ok {
		return errors.Errorf("expected *AuthWorker, got %T", in)
	}
	switch outPtr := out.(type) {
	case *coressh.EphemeralKeysUpdater:
		*outPtr = w
	default:
		return errors.Errorf("expected *coressh.EphemeralKeysUpdater, got %T", out)
	}
	return nil
}
