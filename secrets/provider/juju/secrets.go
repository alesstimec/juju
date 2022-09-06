// Copyright 2021 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package juju

import (
	"github.com/juju/juju/core/secrets"
	"github.com/juju/juju/secrets/provider"
)

const (
	// Store is the name of the Juju secrets store.
	Store = "juju"
)

// NewProvider returns a Juju secrets provider.
func NewProvider() provider.SecretStoreProvider {
	return jujuProvider{}
}

type jujuProvider struct {
}

// Initialise is not used.
func (p jujuProvider) Initialise(m provider.Model) error {
	return nil
}

// CleanupSecrets is not used.
func (p jujuProvider) CleanupSecrets(m provider.Model, removed []*secrets.URI) error {
	return nil
}

// StoreConfig returns nil config params since the Juju store saves
// secret content to the Juju database.
func (p jujuProvider) StoreConfig(m provider.Model, admin bool, owned []*secrets.URI, read []*secrets.URI) (*provider.StoreConfig, error) {
	return &provider.StoreConfig{StoreType: Store}, nil
}

// NewStore returns a nil store since the Juju store saves
// secret content to the Juju database.
func (jujuProvider) NewStore(*provider.StoreConfig) (provider.SecretsStore, error) {
	return nil, nil
}
