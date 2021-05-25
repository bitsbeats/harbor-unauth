package middleware

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"regexp"
)

// routes
/*

  These are all the known urls for the API
  https://docs.docker.com/registry/spec/api/

  /v2/_catalog
  /v2/<name>/blobs/<digest>
  /v2/<name>/blobs/<digest>
  /v2/<name>/blobs/uploads/
  /v2/<name>/blobs/uploads/?mount=<digest>&from=<repository  name>
  /v2/<name>/blobs/uploads/<uuid>
  /v2/<name>/blobs/uploads/<uuid>
  /v2/<name>/blobs/uploads/<uuid>
  /v2/<name>/blobs/uploads/<uuid>?digest=<digest>
  /v2/<name>/blobs/uploads/<uuid>?digest=<digest>
  /v2/<name>/manifests/<reference>
  /v2/<name>/manifests/<reference>
  /v2/<name>/manifests/<reference>
  /v2/<name>/manifests/<reference>
  /v2/<name>/tags/list
  /v2/<name>/tags/list?n=<integer>
*/
var projectMatch = regexp.MustCompile(`^/v2/([^/]*)/.*`)

type (
	UnauthMiddleware struct {
		upstream      *url.URL
		tokenProvider TokenProvider
	}

	TokenProvider interface {
		GetPushPullToken(project string) (string, error)
		GetCatalogToken() (string, error)
	}
)

func NewUnauthMiddleware(upstream *url.URL, tokenProvider TokenProvider) *UnauthMiddleware {
	return &UnauthMiddleware{
		upstream:      upstream,
		tokenProvider: tokenProvider,
	}
}

func (um *UnauthMiddleware) Middleware() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token, err := um.getToken(r)
			if err != nil {
				log.Printf("unable to inject auth into %q: %s", r.RequestURI, err)
			} else {
				bearer := fmt.Sprintf("Bearer %s", token)
				r.Header.Add("Authorization", bearer)
			}
			next.ServeHTTP(w, r)
		})
	}
}

func (um *UnauthMiddleware) getToken(r *http.Request) (string, error) {
	switch r.RequestURI {
	case "/v2":
		fallthrough
	case "/v2/":
		fallthrough
	case "/v2/_catalog":
		fallthrough
	case "/v2/_catalog/":
		return um.tokenProvider.GetCatalogToken()
	}

	match := projectMatch.FindStringSubmatch(r.RequestURI)
	if len(match) != 2 {
		return "", fmt.Errorf("url does not match token injector")
	}
	project := match[1]
	return um.tokenProvider.GetPushPullToken(project)

}
