package middleware

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"strings"

	"github.com/bitsbeats/harbor-unauth/core"
)

type (
	CIDRCheck struct {
		allowList  []*net.IPNet
		proxyCount int
	}
)

func NewCIDRCheck(config *core.Config) (*CIDRCheck, error) {
	allowList := make([]*net.IPNet, len(config.AllowList))
	for i, cidrString := range config.AllowList {
		_, cidr, err := net.ParseCIDR(cidrString)
		if err != nil {
			return nil, fmt.Errorf("unable to parse cidr %q: %w", cidrString, err)
		}
		allowList[i] = cidr
	}
	return &CIDRCheck{
		allowList:  allowList,
		proxyCount: config.ProxyCount,
	}, nil
}

func (c *CIDRCheck) Middleware() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			clientIP := c.getRequestRemoteIP(r)

			if !c.validateIP(clientIP) {
				w.WriteHeader(401)
				_, _ = fmt.Fprint(w, "unauthorized\n")
				log.Printf("unauthorized access from %v", clientIP)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func (c *CIDRCheck) validateIP(ip net.IP) bool {
	for _, cidr := range c.allowList {
		if cidr.Contains(ip) {
			return true
		}
	}
	return false
}

func (c *CIDRCheck) getRequestRemoteIP(r *http.Request) net.IP {
	remoteAddr, _, _ := net.SplitHostPort(r.RemoteAddr)
	clientIP := net.ParseIP(remoteAddr)

	forwarded := r.Header.Get("X-Forwarded-For")
	if c.proxyCount > 0 && forwarded != "" {
		proxies := strings.Split(forwarded, ", ")

		// check if we got through all supplied proxies
		proxyCount := len(proxies)
		if proxyCount < c.proxyCount {
			log.Printf("too few proxies (%d): %s", proxyCount, forwarded)
			return clientIP
		}

		proxy := proxies[c.proxyCount - 1]
		clientIP = net.ParseIP(proxy)

	}

	return clientIP
}
