package main

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/bitsbeats/harbor-unauth/authloader"
	"github.com/bitsbeats/harbor-unauth/config"
	"github.com/bitsbeats/harbor-unauth/middleware"
	"github.com/bitsbeats/harbor-unauth/token"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("unable to load config: %s", err)
	}
	upstream, err := url.Parse(cfg.URL)
	if err != nil {
		log.Fatalf("unable to parse url: %s", err)
	}

	// setup app
	cidrCheck, err := middleware.NewCIDRCheck(cfg)
	if err != nil {
		log.Fatalf("unable to parse cidrs: %s", err)
	}
	authLoader := authloader.NewAuthLoader(cfg)
	tokenProvider := token.NewTokenProvider(upstream, authLoader)
	proxy := httputil.NewSingleHostReverseProxy(upstream)
	unauthMiddleware := middleware.NewUnauthMiddleware(upstream, tokenProvider)

	// setup mux
	mux := http.NewServeMux()
	mux.Handle("/", proxy)
	handler := middleware.Register(
		mux,
		middleware.Logger,
		cidrCheck.Middleware(),
		unauthMiddleware.Middleware(),
	)

	// run
	err = http.ListenAndServe(":5000", handler)
	if err != nil {
		log.Fatalf("unable to listen: %s", err)
	}

}
