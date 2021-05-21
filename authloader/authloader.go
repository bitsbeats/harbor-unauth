package authloader

import (
	"github.com/bitsbeats/harbor-unauth/core"
)

type (
	AuthLoader struct {
		auths map[string]core.Auth
	}
)

func NewAuthLoader(config *core.Config) *AuthLoader {
	return &AuthLoader{
		config.Auths,
	}
}

func (a *AuthLoader) Lookup(project string) (core.Auth, bool) {
	// if no project define return a random one
	// required for /v2/_catalog
	if project == "" {
		for _, value := range a.auths {
			return value, true
		}
	}

	auth, ok := a.auths[project]
	return auth, ok
}
