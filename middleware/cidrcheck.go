package middleware

import (
	"fmt"
	"log"
	"net"
	"net/http"

	"github.com/bitsbeats/harbor-unauth/core"
)

type (
	CIDRCheck struct {
		allowList []*net.IPNet
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
		allowList: allowList,
	}, nil
}

func (c *CIDRCheck) Middleware() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			clientIP := getRequestRemoteIP(r)

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

func getRequestRemoteIP(r *http.Request) net.IP {
	remoteAddr, _, _ := net.SplitHostPort(r.RemoteAddr)
	clientIP := net.ParseIP(remoteAddr)
	return clientIP
}
