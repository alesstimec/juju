// Copyright 2013 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package keyupdater

import (
	"context"
	stdtesting "testing"

	"github.com/canonical/gomock/gomock"
	"github.com/juju/names/v6"
	"github.com/juju/tc"

	"github.com/juju/juju/apiserver/common"
	apiservertesting "github.com/juju/juju/apiserver/testing"
	"github.com/juju/juju/internal/errors"
	"github.com/juju/juju/rpc/params"
)

type authorisedKeysSuite struct {
	authorizer        apiservertesting.FakeAuthorizer
	keyUpdaterService *MockKeyUpdaterService
	machineTag        names.MachineTag
}

func TestAuthorisedKeysSuite(t *stdtesting.T) {
	tc.Run(t, &authorisedKeysSuite{})
}

func (s *authorisedKeysSuite) setupMocks(c *tc.C) *gomock.Controller {
	ctrl := gomock.NewController(c)
	s.keyUpdaterService = NewMockKeyUpdaterService(ctrl)
	return ctrl
}

func (s *authorisedKeysSuite) SetUpTest(c *tc.C) {
	s.machineTag = names.NewMachineTag("0")

	// The default auth is as a machine agent.
	s.authorizer = apiservertesting.FakeAuthorizer{
		Tag: s.machineTag,
	}
}

func (s *authorisedKeysSuite) getCanRead(context.Context) (common.AuthFunc, error) {
	return s.authorizer.AuthOwner, nil
}

// TestAuthorisedKeysForNone is asserting that if we ask for authorised keys
// for zero machines we get back an empty result with no errors and no call
// to the service.
func (s *authorisedKeysSuite) TestAuthorisedKeysForNone(c *tc.C) {
	defer s.setupMocks(c).Finish()
	endPoint := newKeyUpdaterAPI(s.getCanRead, s.keyUpdaterService)

	results, err := endPoint.AuthorisedKeys(c.Context(), params.Entities{})
	c.Assert(err, tc.ErrorIsNil)
	c.Assert(results.Results, tc.HasLen, 0)
}

// TestAuthorisedKeys is asserting that the caller can get back the
// controller's authorised keys for the authenticated machine. For any other
// machines that the caller is not authenticated for we get back unauthorised
// errors, and the service is only consulted once.
func (s *authorisedKeysSuite) TestAuthorisedKeys(c *tc.C) {
	defer s.setupMocks(c).Finish()
	endPoint := newKeyUpdaterAPI(s.getCanRead, s.keyUpdaterService)

	args := params.Entities{
		Entities: []params.Entity{
			{Tag: "machine-40"},
			{Tag: s.machineTag.String()},
			{Tag: "machine-42"},
		},
	}

	s.keyUpdaterService.EXPECT().GetControllerAuthorisedKeys(gomock.Any()).
		Return([]string{"key1", "key2"}, nil).Times(1)

	result, err := endPoint.AuthorisedKeys(c.Context(), args)
	c.Check(err, tc.ErrorIsNil)
	c.Assert(result, tc.DeepEquals, params.StringsResults{
		Results: []params.StringsResult{
			{Error: &params.Error{
				Code:    params.CodeUnauthorized,
				Message: "no permission to read authorised keys for \"machine-40\"",
			}},
			{Result: []string{"key1", "key2"}},
			{Error: &params.Error{
				Code:    params.CodeUnauthorized,
				Message: "no permission to read authorised keys for \"machine-42\"",
			}},
		},
	})
}

// TestAuthorisedKeysUnauthorisedOnly is asserting that when no entity passes
// authorisation the service is never consulted.
func (s *authorisedKeysSuite) TestAuthorisedKeysUnauthorisedOnly(c *tc.C) {
	defer s.setupMocks(c).Finish()
	endPoint := newKeyUpdaterAPI(s.getCanRead, s.keyUpdaterService)

	args := params.Entities{
		Entities: []params.Entity{
			{Tag: "machine-40"},
		},
	}

	result, err := endPoint.AuthorisedKeys(c.Context(), args)
	c.Check(err, tc.ErrorIsNil)
	c.Assert(result, tc.DeepEquals, params.StringsResults{
		Results: []params.StringsResult{
			{Error: &params.Error{
				Code:    params.CodeUnauthorized,
				Message: "no permission to read authorised keys for \"machine-40\"",
			}},
		},
	})
}

// TestAuthorisedKeysForNonMachineEntity is asserting that if we try and get
// authorised keys for a non machine entity we get back a
// [params.CodeTagKindNotSupported] error.
func (s *authorisedKeysSuite) TestAuthorisedKeysForNonMachineEntity(c *tc.C) {
	defer s.setupMocks(c).Finish()
	endPoint := newKeyUpdaterAPI(s.getCanRead, s.keyUpdaterService)

	args := params.Entities{
		Entities: []params.Entity{
			{Tag: names.NewUnitTag("ubuntu/1").String()},
		},
	}

	result, err := endPoint.AuthorisedKeys(c.Context(), args)
	c.Check(err, tc.ErrorIsNil)
	c.Assert(result, tc.DeepEquals, params.StringsResults{
		Results: []params.StringsResult{
			{Error: &params.Error{
				Code:    params.CodeTagKindNotSupported,
				Message: "tag \"unit-ubuntu-1\" unsupported, can only accept tags of kind \"machine\"",
			}},
		},
	})
}

// TestAuthorisedKeysForInvalidTag is asserting that an unparseable tag
// yields a [params.CodeTagInvalid] error for that entity.
func (s *authorisedKeysSuite) TestAuthorisedKeysForInvalidTag(c *tc.C) {
	defer s.setupMocks(c).Finish()
	endPoint := newKeyUpdaterAPI(s.getCanRead, s.keyUpdaterService)

	args := params.Entities{
		Entities: []params.Entity{
			{Tag: "not-a-tag"},
		},
	}

	result, err := endPoint.AuthorisedKeys(c.Context(), args)
	c.Check(err, tc.ErrorIsNil)
	c.Assert(result.Results, tc.HasLen, 1)
	c.Check(result.Results[0].Error.Code, tc.Equals, params.CodeTagInvalid)
}

// TestAuthorisedKeysServiceError is asserting that a service failure fails
// the whole call rather than being reported per entity.
func (s *authorisedKeysSuite) TestAuthorisedKeysServiceError(c *tc.C) {
	defer s.setupMocks(c).Finish()
	endPoint := newKeyUpdaterAPI(s.getCanRead, s.keyUpdaterService)

	args := params.Entities{
		Entities: []params.Entity{
			{Tag: s.machineTag.String()},
		},
	}

	boom := errors.New("boom")
	s.keyUpdaterService.EXPECT().GetControllerAuthorisedKeys(gomock.Any()).
		Return(nil, boom)

	_, err := endPoint.AuthorisedKeys(c.Context(), args)
	c.Check(err, tc.ErrorIs, boom)
}
