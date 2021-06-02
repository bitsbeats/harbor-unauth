package authloader

import (
	"github.com/bitsbeats/harbor-unauth/core"
)

type (
	AuthLoader struct {
		auth     core.Auth
		projects []string
	}
)

func NewAuthLoader(config *core.Config) *AuthLoader {
	return &AuthLoader{
		config.Auth,
		config.Projects,
	}
}

func (a *AuthLoader) Lookup() (core.Auth, bool) {
	return a.auth, true
}

func (a *AuthLoader) Projects() []string {
	return a.projects
}
