package config

type SiteAuthConfig struct {
	SiteID          string
	Owner           string
	Name            string
	Domain          string
	OtherDomains    []string
	Status          string
	AuthApplication string
	Auth            AuthConfig
	Upstream        UpstreamConfig
}

type AuthConfig struct {
	Endpoint     string
	ClientID     string
	ClientSecret string
	Certificate  string
	Organization string
	Application  string
}

type UpstreamConfig struct {
	URL     string
	Headers map[string]string
}
