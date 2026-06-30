package server

import (
	"net"
	"net/http"
	"net/netip"
	"strings"

	"github.com/lwmacct/260630-go-hsr-auth/pkg/auth"
)

type requestContextMiddleware struct {
	trustedProxies []netip.Prefix
}

func newRequestContextMiddleware(trustedProxies []string) requestContextMiddleware {
	return requestContextMiddleware{trustedProxies: parseTrustedProxies(trustedProxies)}
}

func (m requestContextMiddleware) Wrap(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		request, ok := m.sessionRequest(r)
		if ok {
			r = r.WithContext(auth.ContextWithRequest(r.Context(), request))
		}
		next.ServeHTTP(w, r)
	})
}

func (m requestContextMiddleware) sessionRequest(r *http.Request) (auth.SessionRequest, bool) {
	ip, ok := m.clientIP(r)
	if !ok {
		return auth.SessionRequest{}, false
	}
	return auth.SessionRequest{
		IP:         ip.String(),
		Host:       r.Host,
		UserAgent:  r.UserAgent(),
		Method:     r.Method,
		Path:       r.URL.Path,
		RemoteAddr: r.RemoteAddr,
	}, true
}

func (m requestContextMiddleware) clientIP(r *http.Request) (netip.Addr, bool) {
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		host = r.RemoteAddr
	}
	remoteIP, ok := parseIP(host)
	if !ok {
		return netip.Addr{}, false
	}
	if len(m.trustedProxies) == 0 || !ipInPrefixes(remoteIP, m.trustedProxies) {
		return remoteIP, true
	}
	if ip, ok := parseXForwardedFor(r.Header.Get("X-Forwarded-For")); ok {
		return ip, true
	}
	if ip, ok := parseIP(r.Header.Get("X-Real-IP")); ok {
		return ip, true
	}
	return remoteIP, true
}

func parseTrustedProxies(values []string) []netip.Prefix {
	prefixes := make([]netip.Prefix, 0, len(values))
	for _, value := range values {
		value = strings.TrimSpace(value)
		if value == "" {
			continue
		}
		if prefix, err := netip.ParsePrefix(value); err == nil {
			prefixes = append(prefixes, prefix)
			continue
		}
		if addr, ok := parseIP(value); ok {
			prefixes = append(prefixes, netip.PrefixFrom(addr, addr.BitLen()))
		}
	}
	return prefixes
}

func ipInPrefixes(ip netip.Addr, prefixes []netip.Prefix) bool {
	for _, prefix := range prefixes {
		if prefix.Contains(ip) {
			return true
		}
	}
	return false
}

func parseXForwardedFor(value string) (netip.Addr, bool) {
	for _, part := range strings.Split(value, ",") {
		if ip, ok := parseIP(part); ok {
			return ip, true
		}
	}
	return netip.Addr{}, false
}

func parseIP(value string) (netip.Addr, bool) {
	value = strings.TrimSpace(strings.Trim(value, `"`))
	if value == "" || strings.EqualFold(value, "unknown") {
		return netip.Addr{}, false
	}
	if ip, err := netip.ParseAddr(value); err == nil {
		return ip, true
	}
	if host, _, err := net.SplitHostPort(value); err == nil {
		if ip, err := netip.ParseAddr(host); err == nil {
			return ip, true
		}
	}
	return netip.Addr{}, false
}
