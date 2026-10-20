// Copyright 2013 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package sshkeyupdater

import (
	"context"

	"github.com/juju/errors"
	"github.com/juju/names/v6"
	"github.com/juju/utils/v4/ssh"
	"github.com/juju/worker/v5"
	"github.com/juju/worker/v5/catacomb"
	gossh "golang.org/x/crypto/ssh"

	"github.com/juju/juju/core/logger"
	coremachineauthentication "github.com/juju/juju/core/machineauthentication"
)

// The user name used to ssh into Juju nodes.
// Override for testing.
var SSHUser = "ubuntu"

// authorizedKeysFile is the name of the ssh authorized_keys file that the
// worker manages, relative to the user's .ssh directory.
const authorizedKeysFile = "authorized_keys"

// opType identifies the kind of ephemeral key operation carried by a request.
type opType int

const (
	addOp opType = iota
	removeOp
)

// ephemeralRequest is used to pass ephemeral key add/remove operations into the
// worker loop, so that all writes to the authorized_keys file are serialised
// through a single goroutine. The done channel is buffered so the loop never
// blocks replying to a caller that has stopped waiting.
type ephemeralRequest struct {
	op      opType
	key     gossh.PublicKey
	comment string
	done    chan error
}

// makeAddRequest creates a request to add the supplied ephemeral key, tagged
// with the given comment for later removal.
func makeAddRequest(key gossh.PublicKey, comment string) ephemeralRequest {
	return ephemeralRequest{
		op:      addOp,
		key:     key,
		comment: comment,
		done:    make(chan error, 1),
	}
}

// makeRemoveRequest creates a request to remove the supplied ephemeral key.
func makeRemoveRequest(key gossh.PublicKey) ephemeralRequest {
	return ephemeralRequest{
		op:   removeOp,
		key:  key,
		done: make(chan error, 1),
	}
}

// FacadeClient is the subset of the KeyUpdater facade the worker needs.
type FacadeClient interface {
	// AuthorisedKeys returns the controller's authorised keys that must be
	// present on the machine identified by tag.
	AuthorisedKeys(ctx context.Context, tag names.MachineTag) ([]string, error)
}

// WorkerConfig holds the dependencies of the sshkeyupdater worker.
type WorkerConfig struct {
	// Logger is the logger to use for the worker.
	Logger logger.Logger
	// MachineTag is the tag of the machine this agent runs on.
	MachineTag names.MachineTag
	// FacadeClient is used to fetch the controller's authorised keys.
	FacadeClient FacadeClient
}

// Validate validates the worker configuration.
func (c WorkerConfig) Validate() error {
	if c.Logger == nil {
		return errors.NotValidf("nil Logger")
	}
	if c.MachineTag.Id() == "" {
		return errors.NotValidf("empty MachineTag")
	}
	if c.FacadeClient == nil {
		return errors.NotValidf("nil FacadeClient")
	}
	return nil
}

// AuthWorker keeps the controller's own SSH keys present in the machine's
// authorized_keys file and serialises ephemeral SSH key updates for reverse
// tunnels.
type AuthWorker struct {
	catacomb catacomb.Catacomb
	config   WorkerConfig

	// requests carries ephemeral key add/remove operations into the loop.
	requests chan ephemeralRequest
}

// NewWorker returns a worker that reconciles the controller's SSH keys on
// start and then manages ephemeral SSH keys.
func NewWorker(config WorkerConfig) (worker.Worker, error) {
	if err := config.Validate(); err != nil {
		return nil, errors.Trace(err)
	}
	w := &AuthWorker{
		config:   config,
		requests: make(chan ephemeralRequest),
	}
	if err := catacomb.Invoke(catacomb.Plan{
		Name: "authentication",
		Site: &w.catacomb,
		Work: w.loop,
	}); err != nil {
		return nil, errors.Trace(err)
	}
	return w, nil
}

// Kill is part of the worker.Worker interface.
func (a *AuthWorker) Kill() {
	a.catacomb.Kill(nil)
}

// Wait is part of the worker.Worker interface.
func (a *AuthWorker) Wait() error {
	return a.catacomb.Wait()
}

// AddEphemeralKey adds an ephemeral key to the authorized_keys file. The
// supplied comment is used to identify the key for later removal. The write is
// performed by the worker loop.
func (a *AuthWorker) AddEphemeralKey(key gossh.PublicKey, comment string) error {
	return a.enqueue(makeAddRequest(key, comment))
}

// RemoveEphemeralKey removes an ephemeral key from the authorized_keys file.
// The write is performed by the worker loop.
func (a *AuthWorker) RemoveEphemeralKey(ephemeralKey gossh.PublicKey) error {
	return a.enqueue(makeRemoveRequest(ephemeralKey))
}

// enqueue sends an ephemeral key request into the worker loop and waits for it
// to be applied, returning any error. If the worker is dying it returns an
// error indicating so.
func (a *AuthWorker) enqueue(req ephemeralRequest) error {
	select {
	case a.requests <- req:
	case <-a.catacomb.Dying():
		return errors.Trace(coremachineauthentication.ErrSShKeyUpdaterWorkerDying)
	}

	select {
	case err := <-req.done:
		return errors.Trace(err)
	case <-a.catacomb.Dying():
		return errors.Trace(coremachineauthentication.ErrSShKeyUpdaterWorkerDying)
	}
}

func (a *AuthWorker) loop() error {
	ctx, cancel := a.scopedContext()
	defer cancel()

	// The controller's keys are reconciled once per start. The worker is
	// restarted whenever the agent reconnects to a controller, which covers
	// the post-migration case where the machine must learn the keys of the
	// controller that now manages it.
	if err := a.reconcileControllerKeys(ctx); err != nil {
		return errors.Trace(err)
	}

	for {
		select {
		case <-a.catacomb.Dying():
			return a.catacomb.ErrDying()
		case req := <-a.requests:
			// Ephemeral key requests are non-fatal: reply the result to the
			// caller and keep serving. A bad request must not bring down the
			// worker or interrupt other tunnel key operations.
			req.done <- a.handleEphemeralRequest(req)
		}
	}
}

// handleEphemeralRequest applies a single ephemeral key add or remove
// operation. It runs on the loop goroutine, so it is serialised against other
// ephemeral requests.
func (a *AuthWorker) handleEphemeralRequest(req ephemeralRequest) error {
	switch req.op {
	case addOp:
		keyWithComment := ensureJujuEphemeralComment(req.key, req.comment)
		if err := ssh.AddKeys(SSHUser, keyWithComment); err != nil {
			return errors.Trace(err)
		}
		return nil
	case removeOp:
		fingerprint := gossh.FingerprintLegacyMD5(req.key)
		// Use DeleteKeysFromFile rather than DeleteKeys: the latter refuses to
		// remove the final key in the file, which would leave a torn-down
		// tunnel's ephemeral key authorised until the worker restarts.
		if err := ssh.DeleteKeysFromFile(SSHUser, authorizedKeysFile, []string{fingerprint}); err != nil {
			return errors.Trace(err)
		}
		return nil
	default:
		return errors.Errorf("unknown op %v", req.op)
	}
}

// scopedContext returns a context that is cancelled when the worker is dying.
func (a *AuthWorker) scopedContext() (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithCancel(context.Background())
	return a.catacomb.Context(ctx), cancel
}
