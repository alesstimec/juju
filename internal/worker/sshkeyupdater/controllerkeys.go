// Copyright 2026 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package sshkeyupdater

import (
	"context"
	"maps"
	"slices"

	"github.com/juju/errors"
	"github.com/juju/utils/v4/ssh"

	"github.com/juju/juju/domain/controller"
)

// controllerKeyComment is the comment the controller's own key carries in
// the authorized_keys file. It matches what cloud-init writes on fresh
// machines, where every Juju managed key is prefixed with "Juju:".
const controllerKeyComment = ssh.JujuCommentPrefix + controller.ControllerSSHKeyComment

// reconcileControllerKeys ensures the controller's own SSH keys are present
// in the machine's authorized_keys file, and that controller keys belonging
// to a controller that no longer manages this machine are removed.
//
// Only keys carrying the controller key comment are ever removed. Keys with
// any other comment, including keys added by hand, are left untouched. A
// wanted key that is already present under a different comment is treated as
// present rather than duplicated, since the file rejects duplicate
// fingerprints.
func (a *AuthWorker) reconcileControllerKeys(ctx context.Context) error {
	keys, err := a.config.FacadeClient.AuthorisedKeys(ctx, a.config.MachineTag)
	if err != nil {
		return errors.Annotate(err, "getting controller authorised keys")
	}

	// wanted maps each controller key fingerprint to the line that should
	// appear in the file.
	wanted := make(map[string]string, len(keys))
	for _, key := range keys {
		fingerprint, _, err := ssh.KeyFingerprint(key)
		if err != nil {
			a.config.Logger.Warningf(ctx, "ignoring invalid controller ssh key %q: %v", key, err)
			continue
		}
		wanted[fingerprint] = ssh.EnsureJujuComment(key)
	}

	existing, err := ssh.ListKeysFromFile(SSHUser, authorizedKeysFile, ssh.FullKeys)
	if err != nil {
		return errors.Annotate(err, "listing authorized keys")
	}

	var stale []string
	for _, key := range existing {
		fingerprint, comment, err := ssh.KeyFingerprint(key)
		if err != nil {
			a.config.Logger.Warningf(ctx, "ignoring invalid authorized key %q: %v", key, err)
			continue
		}
		if _, ok := wanted[fingerprint]; ok {
			delete(wanted, fingerprint)
			continue
		}
		if isControllerKeyComment(comment) {
			stale = append(stale, fingerprint)
		}
	}

	// Remove before adding: the file also rejects duplicate comments, and a
	// stale controller key carries the same comment as the wanted one.
	if len(stale) > 0 {
		if err := ssh.DeleteKeysFromFile(SSHUser, authorizedKeysFile, stale); err != nil {
			return errors.Annotate(err, "removing stale controller ssh keys")
		}
		a.config.Logger.Infof(ctx, "removed %d stale controller ssh key(s)", len(stale))
	}
	if len(wanted) > 0 {
		toAdd := slices.Sorted(maps.Values(wanted))
		if err := ssh.AddKeys(SSHUser, toAdd...); err != nil {
			return errors.Annotate(err, "adding controller ssh keys")
		}
		a.config.Logger.Infof(ctx, "added %d controller ssh key(s)", len(toAdd))
	}
	return nil
}

// isControllerKeyComment reports whether the comment identifies a key as the
// controller's own key, with or without the Juju comment prefix.
func isControllerKeyComment(comment string) bool {
	return comment == controllerKeyComment ||
		comment == controller.ControllerSSHKeyComment
}
