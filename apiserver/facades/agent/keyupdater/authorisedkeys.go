// Copyright 2013 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package keyupdater

import (
	"context"
	"fmt"

	"github.com/juju/names/v6"

	"github.com/juju/juju/apiserver/common"
	apiservererrors "github.com/juju/juju/apiserver/errors"
	"github.com/juju/juju/rpc/params"
)

// KeyUpdaterAPI implements the KeyUpdater interface and is the concrete
// implementation of the API endpoint. It is consumed by the machine agent's
// sshkeyupdater worker to keep the controller's own SSH keys present on the
// machine, so that the controller can still reach it after a model migration
// or during disaster recovery. User keys are never served: users reach
// machines through the controller's SSH proxy.
type KeyUpdaterAPI struct {
	getCanRead        common.GetAuthFunc
	keyUpdaterService KeyUpdaterService
}

// newKeyUpdaterAPI constructs a new [KeyUpdaterAPI] for use in the Juju facade
// model.
func newKeyUpdaterAPI(
	getCanRead common.GetAuthFunc,
	keyUpdaterService KeyUpdaterService,
) *KeyUpdaterAPI {
	return &KeyUpdaterAPI{
		getCanRead:        getCanRead,
		keyUpdaterService: keyUpdaterService,
	}
}

// AuthorisedKeys reports the controller's authorised ssh keys for the
// specified machines. The keys are the same for every machine in the
// controller; the machine tag is used for authorisation only.
// The following param error codes can be expected:
// - [params.CodeTagInvalid] When a tag provided does not parse and is
// considered invalid.
// - [params.CodeTagKindNotSupported] When a tag has been supplied that is not a
// machine tag.
// - [params.CodeUnathorized] When the caller does not have permissions to get
// the authorised keys for a requested tag.
func (api *KeyUpdaterAPI) AuthorisedKeys(
	ctx context.Context,
	arg params.Entities,
) (params.StringsResults, error) {
	if len(arg.Entities) == 0 {
		return params.StringsResults{}, nil
	}
	results := make([]params.StringsResult, len(arg.Entities))

	canRead, err := api.getCanRead(ctx)
	if err != nil {
		return params.StringsResults{}, fmt.Errorf(
			"checking can read for authorised keys: %w",
			err,
		)
	}

	// The controller keys do not vary per machine, so fetch them once and
	// only after at least one entity has passed authorisation.
	var keys []string
	fetched := false

	for i, entity := range arg.Entities {
		tag, err := names.ParseTag(entity.Tag)
		if err != nil {
			results[i].Error = apiservererrors.ParamsErrorf(
				params.CodeTagInvalid,
				"cannot parse tag %q: %s",
				entity.Tag,
				err.Error(),
			)
			continue
		}

		if tag.Kind() != names.MachineTagKind {
			results[i].Error = apiservererrors.ParamsErrorf(
				params.CodeTagKindNotSupported,
				"tag %q unsupported, can only accept tags of kind %q",
				tag, names.MachineTagKind,
			)
			continue
		}

		if !canRead(tag) {
			results[i].Error = apiservererrors.ParamsErrorf(
				params.CodeUnauthorized,
				"no permission to read authorised keys for %q",
				tag,
			)
			continue
		}

		if !fetched {
			keys, err = api.keyUpdaterService.GetControllerAuthorisedKeys(ctx)
			if err != nil {
				// We don't understand this error. At this stage we consider
				// it an internal server error and bail out of the call
				// completely.
				return params.StringsResults{}, fmt.Errorf(
					"cannot get controller authorised keys: %w", err,
				)
			}
			fetched = true
		}

		results[i].Result = keys
	}

	return params.StringsResults{
		Results: results,
	}, nil
}
