package token

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/bitsbeats/harbor-unauth/core"
)

type (
	TokenProvider struct {
		url        *url.URL
		authLoader AuthLoader
	}

	tokenResponse struct {
		Token string `json:"token"`
	}

	AuthLoader interface {
		Lookup(project string) (auth core.Auth, ok bool)
	}
)

func NewTokenProvider(url *url.URL, authLoader AuthLoader) *TokenProvider {
	return &TokenProvider{url, authLoader}
}

func (t *TokenProvider) GetCatalogToken() (string, error) {
	return t.getUrlFor("", "registry:catalog:*")
}

func (t *TokenProvider) GetPushPullToken(project string) (string, error) {
	scope := fmt.Sprintf("repository:%s/:push", project)
	return t.getUrlFor(project, scope)
}

func (t *TokenProvider) getUrlFor(project, scope string) (string, error) {
	url := fmt.Sprintf(
		"%s://%s/service/token?service=harbor-registry&scope=%s",
		t.url.Scheme, t.url.Host, scope,
	)

	auth, ok := t.authLoader.Lookup(project)
	if !ok {
		return "", fmt.Errorf("no access data for %q found", project)
	}

	token, err := t.getToken(url, auth.User, auth.Password)
	if err != nil {
		return "", fmt.Errorf("unable to load token: %w", err)
	}
	return token, nil
}

func (t *TokenProvider) getToken(url, user, password string) (string, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", fmt.Errorf("unable to create request: %w", err)
	}
	req.SetBasicAuth(user, password)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("unable to handle request: %w", err)
	}
	defer resp.Body.Close()

	tr := &tokenResponse{}
	err = json.NewDecoder(resp.Body).Decode(tr)
	if err != nil {
		return "", fmt.Errorf("unable to parse json: %w", err)
	}
	return tr.Token, nil
}
