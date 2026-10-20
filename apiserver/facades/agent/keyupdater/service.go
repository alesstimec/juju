// Copyright 2024 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package keyupdater

import (
	"context"
)

// KeyUpdaterService is the interface for retrieving the controller's
// authorised keys.
type KeyUpdaterService interface {
	// GetControllerAuthorisedKeys returns the controller's own authorised
	// keys that must be present on every machine the controller manages.
	GetControllerAuthorisedKeys(context.Context) ([]string, error)
}
