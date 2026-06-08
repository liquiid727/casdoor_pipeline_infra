package config

import (
	"errors"
	"fmt"
	"strings"
)

var (
	ErrSiteNotFound = errors.New("site not found")
	ErrSiteInactive = errors.New("site inactive")
)

type Resolver interface {
	ResolveByHost(host string) (*SiteAuthConfig, error)
}

type StaticResolver struct {
	sitesByDomain map[string]*SiteAuthConfig
}

func NewStaticResolver(sites []SiteAuthConfig) *StaticResolver {
	sitesByDomain := map[string]*SiteAuthConfig{}
	for i := range sites {
		site := sites[i]
		addSiteDomain(sitesByDomain, site.Domain, &site)
		for _, domain := range site.OtherDomains {
			addSiteDomain(sitesByDomain, domain, &site)
		}
	}
	return &StaticResolver{sitesByDomain: sitesByDomain}
}

func (r *StaticResolver) ResolveByHost(host string) (*SiteAuthConfig, error) {
	domain := normalizeHost(host)
	site, ok := r.sitesByDomain[domain]
	if !ok {
		return nil, fmt.Errorf("%w: %s", ErrSiteNotFound, domain)
	}
	if site.Status == "Inactive" {
		return nil, fmt.Errorf("%w: %s", ErrSiteInactive, domain)
	}
	return site, nil
}

func addSiteDomain(sitesByDomain map[string]*SiteAuthConfig, domain string, site *SiteAuthConfig) {
	domain = normalizeHost(domain)
	if domain != "" {
		sitesByDomain[domain] = site
	}
}

func normalizeHost(host string) string {
	host = strings.ToLower(strings.TrimSpace(host))
	if i := strings.Index(host, ":"); i >= 0 {
		host = host[:i]
	}
	if strings.HasPrefix(host, "www.") {
		host = strings.TrimPrefix(host, "www.")
	}
	return host
}
