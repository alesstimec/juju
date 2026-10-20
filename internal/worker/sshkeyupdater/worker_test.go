// Copyright 2013 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package sshkeyupdater_test

import (
	"context"
	"testing"

	"github.com/juju/names/v6"
	"github.com/juju/tc"
	"github.com/juju/utils/v4/ssh"
	sshtesting "github.com/juju/utils/v4/ssh/testing"
	"github.com/juju/worker/v5"
	"github.com/juju/worker/v5/workertest"
	gossh "golang.org/x/crypto/ssh"

	coreerrors "github.com/juju/juju/core/errors"
	coremachineauthentication "github.com/juju/juju/core/machineauthentication"
	coressh "github.com/juju/juju/core/ssh"
	"github.com/juju/juju/internal/errors"
	loggertesting "github.com/juju/juju/internal/logger/testing"
	coretesting "github.com/juju/juju/internal/testing"
	"github.com/juju/juju/internal/worker/sshkeyupdater"
)

var (
	// controllerKey is how the controller's key is returned by the facade:
	// the raw controller config value, with its bare comment.
	controllerKey = sshtesting.ValidKeyOne.Key + " juju-system-key"
	// controllerKeyLine is how the controller's key must appear in the
	// authorized_keys file, matching what cloud-init writes.
	controllerKeyLine = sshtesting.ValidKeyOne.Key + " Juju:juju-system-key"
	// staleControllerKeyLine is a previous controller's key as written by
	// cloud-init, which must be removed on reconcile.
	staleControllerKeyLine = sshtesting.ValidKeyTwo.Key + " Juju:juju-system-key"
	// userKeyLine is a key added outside of Juju, which must never be
	// touched.
	userKeyLine = sshtesting.ValidKeyThree.Key + " alice@example.com"
)

type fakeFacadeClient struct {
	keys []string
	err  error
}

func (f *fakeFacadeClient) AuthorisedKeys(context.Context, names.MachineTag) ([]string, error) {
	return f.keys, f.err
}

type workerSuite struct {
	coretesting.FakeJujuXDGDataHomeSuite
}

func TestWorkerSuite(t *testing.T) {
	tc.Run(t, &workerSuite{})
}

func (s *workerSuite) SetUpTest(c *tc.C) {
	s.FakeJujuXDGDataHomeSuite.SetUpTest(c)
	s.PatchValue(&sshkeyupdater.SSHUser, "")
}

func (s *workerSuite) newWorker(c *tc.C, client sshkeyupdater.FacadeClient) worker.Worker {
	w, err := sshkeyupdater.NewWorker(sshkeyupdater.WorkerConfig{
		Logger:       loggertesting.WrapCheckLog(c),
		MachineTag:   names.NewMachineTag("0"),
		FacadeClient: client,
	})
	c.Assert(err, tc.ErrorIsNil)
	return w
}

func (s *workerSuite) listKeys(c *tc.C) []string {
	keys, err := ssh.ListKeys(sshkeyupdater.SSHUser, ssh.FullKeys)
	c.Assert(err, tc.ErrorIsNil)
	return keys
}

func (s *workerSuite) TestValidateConfig(c *tc.C) {
	cfg := sshkeyupdater.WorkerConfig{
		Logger:       loggertesting.WrapCheckLog(c),
		MachineTag:   names.NewMachineTag("0"),
		FacadeClient: &fakeFacadeClient{},
	}
	c.Check(cfg.Validate(), tc.ErrorIsNil)

	bad := cfg
	bad.Logger = nil
	c.Check(bad.Validate(), tc.ErrorIs, coreerrors.NotValid)

	bad = cfg
	bad.MachineTag = names.MachineTag{}
	c.Check(bad.Validate(), tc.ErrorIs, coreerrors.NotValid)

	bad = cfg
	bad.FacadeClient = nil
	c.Check(bad.Validate(), tc.ErrorIs, coreerrors.NotValid)
}

// TestReconcileAddsControllerKey ensures that on start the controller's key
// is written to an empty authorized_keys file with the Juju comment prefix.
func (s *workerSuite) TestReconcileAddsControllerKey(c *tc.C) {
	w := s.newWorker(c, &fakeFacadeClient{keys: []string{controllerKey}})
	// The reconcile runs before the worker loop services anything else, so
	// once the worker has stopped cleanly the file reflects the result.
	workertest.CleanKill(c, w)

	c.Check(s.listKeys(c), tc.DeepEquals, []string{controllerKeyLine})
}

// TestReconcileReplacesStaleControllerKey ensures that a key left behind by
// a previous controller is removed, the new controller's key is added, and
// keys that are not Juju controller keys are left untouched.
func (s *workerSuite) TestReconcileReplacesStaleControllerKey(c *tc.C) {
	c.Assert(ssh.AddKeys(sshkeyupdater.SSHUser, staleControllerKeyLine, userKeyLine), tc.ErrorIsNil)

	w := s.newWorker(c, &fakeFacadeClient{keys: []string{controllerKey}})
	workertest.CleanKill(c, w)

	c.Check(s.listKeys(c), tc.DeepEquals, []string{userKeyLine, controllerKeyLine})
}

// TestReconcileIsIdempotent ensures that a machine already carrying the
// controller's key, as written by cloud-init, is left exactly as it is.
func (s *workerSuite) TestReconcileIsIdempotent(c *tc.C) {
	c.Assert(ssh.AddKeys(sshkeyupdater.SSHUser, userKeyLine, controllerKeyLine), tc.ErrorIsNil)

	w := s.newWorker(c, &fakeFacadeClient{keys: []string{controllerKey}})
	workertest.CleanKill(c, w)

	c.Check(s.listKeys(c), tc.DeepEquals, []string{userKeyLine, controllerKeyLine})
}

// TestReconcileKeyPresentUnderOtherComment ensures that a wanted key already
// present under a different comment is neither duplicated nor rewritten.
func (s *workerSuite) TestReconcileKeyPresentUnderOtherComment(c *tc.C) {
	other := sshtesting.ValidKeyOne.Key + " added-by-hand"
	c.Assert(ssh.AddKeys(sshkeyupdater.SSHUser, other), tc.ErrorIsNil)

	w := s.newWorker(c, &fakeFacadeClient{keys: []string{controllerKey}})
	workertest.CleanKill(c, w)

	c.Check(s.listKeys(c), tc.DeepEquals, []string{other})
}

// TestReconcileNoControllerKeys ensures that when the controller has no keys
// (for example a Kubernetes controller) stale controller keys are removed and
// nothing is added.
func (s *workerSuite) TestReconcileNoControllerKeys(c *tc.C) {
	c.Assert(ssh.AddKeys(sshkeyupdater.SSHUser, staleControllerKeyLine, userKeyLine), tc.ErrorIsNil)

	w := s.newWorker(c, &fakeFacadeClient{})
	workertest.CleanKill(c, w)

	c.Check(s.listKeys(c), tc.DeepEquals, []string{userKeyLine})
}

// TestReconcileFacadeErrorStopsWorker ensures that a failure to fetch the
// controller's keys is fatal to the worker, so that the dependency engine
// restarts it rather than leaving the machine without the key.
func (s *workerSuite) TestReconcileFacadeErrorStopsWorker(c *tc.C) {
	c.Assert(ssh.AddKeys(sshkeyupdater.SSHUser, userKeyLine), tc.ErrorIsNil)

	w := s.newWorker(c, &fakeFacadeClient{err: errors.New("boom")})
	err := workertest.CheckKilled(c, w)
	c.Check(err, tc.ErrorMatches, "getting controller authorised keys: boom")

	c.Check(s.listKeys(c), tc.DeepEquals, []string{userKeyLine})
}

func (s *workerSuite) TestAddAndRemoveEphemeralKey(c *tc.C) {
	w := s.newWorker(c, &fakeFacadeClient{})
	defer workertest.CleanKill(c, w)

	updater, ok := w.(coressh.EphemeralKeysUpdater)
	c.Assert(ok, tc.IsTrue)
	key, _, _, _, err := gossh.ParseAuthorizedKey([]byte(sshtesting.ValidKeyOne.Key))
	c.Assert(err, tc.ErrorIsNil)

	c.Assert(updater.AddEphemeralKey(key, "tunnel-0"), tc.ErrorIsNil)
	c.Check(s.listKeys(c), tc.DeepEquals, []string{sshtesting.ValidKeyOne.Key + " Juju:Ephemeral:tunnel-0"})

	c.Assert(updater.RemoveEphemeralKey(key), tc.ErrorIsNil)
	c.Check(s.listKeys(c), tc.HasLen, 0)
}

// TestRemoveLastEphemeralKey ensures that removing the only authorized key
// succeeds and leaves the authorized_keys file empty.
func (s *workerSuite) TestRemoveLastEphemeralKey(c *tc.C) {
	w := s.newWorker(c, &fakeFacadeClient{})
	defer workertest.CleanKill(c, w)

	updater := w.(coressh.EphemeralKeysUpdater)
	key, _, _, _, err := gossh.ParseAuthorizedKey([]byte(sshtesting.ValidKeyOne.Key))
	c.Assert(err, tc.ErrorIsNil)

	c.Assert(updater.AddEphemeralKey(key, "tunnel-0"), tc.ErrorIsNil)
	c.Assert(updater.RemoveEphemeralKey(key), tc.ErrorIsNil)

	c.Check(s.listKeys(c), tc.HasLen, 0)
}

func (s *workerSuite) TestEphemeralKeyOpsReturnWorkerDying(c *tc.C) {
	w := s.newWorker(c, &fakeFacadeClient{})
	updater := w.(coressh.EphemeralKeysUpdater)
	workertest.CleanKill(c, w)

	key, _, _, _, err := gossh.ParseAuthorizedKey([]byte(sshtesting.ValidKeyOne.Key))
	c.Assert(err, tc.ErrorIsNil)
	c.Check(updater.AddEphemeralKey(key, "tunnel-0"), tc.ErrorIs, coremachineauthentication.ErrSShKeyUpdaterWorkerDying)
	c.Check(updater.RemoveEphemeralKey(key), tc.ErrorIs, coremachineauthentication.ErrSShKeyUpdaterWorkerDying)
}
