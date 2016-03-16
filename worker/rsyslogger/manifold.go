// Copyright 2016 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package rsyslogger

import (
	"github.com/juju/juju/agent"
	"github.com/juju/juju/api/base"
	"github.com/juju/juju/api/logreader"
	"github.com/juju/juju/worker"
	"github.com/juju/juju/worker/dependency"
	"github.com/juju/juju/worker/util"
)

// ManifoldConfig defines the names of the manifolds on which a
// Manifold will depend.
type ManifoldConfig util.PostUpgradeManifoldConfig

// Manifold returns a dependency manifold that runs a logger
// worker, using the resource names defined in the supplied config.
func Manifold(config ManifoldConfig) dependency.Manifold {
	return util.PostUpgradeManifold(util.PostUpgradeManifoldConfig(config), newWorker)
}

// newWorker trivially wraps NewLogger to specialise a PostUpgradeManifold.
var newWorker = func(a agent.Agent, apiCaller base.APICaller) (worker.Worker, error) {
	currentConfig := a.CurrentConfig()
	logreaderFacade := logreader.NewAPI(apiCaller)
	return NewRsysWorker(logreaderFacade, currentConfig)
}
