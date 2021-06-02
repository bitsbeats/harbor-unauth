package token

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"

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
		Lookup() (auth core.Auth, ok bool)
		Projects() []string
	}
)

func NewTokenProvider(url *url.URL, authLoader AuthLoader) *TokenProvider {
	return &TokenProvider{url, authLoader}
}

func (t *TokenProvider) GetCatalogToken() (string, error) {
	return t.getUrlFor("scope=registry:catalog:*")
}

func (t *TokenProvider) GetToken() (string, error) {
	scopes := t.getScopes()
	return t.getUrlFor(scopes)
}

func (t *TokenProvider) getUrlFor(scopes string) (string, error) {
	url := fmt.Sprintf(
		"%s://%s/service/token?service=harbor-registry&%s",
		t.url.Scheme, t.url.Host, scopes,
	)

	auth, ok := t.authLoader.Lookup()
	if !ok {
		return "", fmt.Errorf("no access data found")
	}

	token, err := t.getToken(url, auth.User, auth.Password)
	if err != nil {
		return "", fmt.Errorf("unable to load token: %w", err)
	}
	return token, nil
}

func (t *TokenProvider) getScopes() string {
	projects := t.authLoader.Projects()
	scopes := make([]string, len(projects))
	for i, project := range projects {
		scopes[i] = fmt.Sprintf("scope=repository:%s/:push", project)
	}
	return strings.Join(scopes, "&")
}

func (t *TokenProvider) getToken(url, user, password string) (string, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", fmt.Errorf("unable to create request: %w", err)
	}
	req.SetBasicAuth(user, password)
	log.Printf("auth: %s:%s", user, password)

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
