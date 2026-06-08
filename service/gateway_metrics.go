package service

import (
	"github.com/liquiid727/pipeline-auth/object"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	gatewayLoginStartTotal = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "gateway_login_start_total",
		Help: "Total OAuth login starts initiated by the gateway.",
	}, []string{"site", "domain", "authApplication"})
	gatewayCallbackSuccessTotal = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "gateway_callback_success_total",
		Help: "Total successful OAuth callbacks handled by the gateway.",
	}, []string{"site", "domain", "authApplication"})
	gatewayCallbackFailureTotal = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "gateway_callback_failure_total",
		Help: "Total failed OAuth callbacks handled by the gateway.",
	}, []string{"site", "domain", "authApplication", "reason"})
	gatewayTokenInvalidTotal = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "gateway_token_invalid_total",
		Help: "Total invalid access tokens seen by the gateway.",
	}, []string{"site", "domain", "authApplication"})
	gatewayInactiveSiteBlockTotal = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "gateway_inactive_site_block_total",
		Help: "Total requests blocked because a site is inactive.",
	}, []string{"site", "domain"})
	gatewayUpstreamErrorTotal = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "gateway_upstream_error_total",
		Help: "Total upstream proxy errors seen by the gateway.",
	}, []string{"site", "domain", "upstream"})
)

func observeGatewayLoginStart(site *object.Site) {
	gatewayLoginStartTotal.WithLabelValues(siteLabel(site), domainLabel(site), authApplicationLabel(site)).Inc()
}

func observeGatewayCallbackSuccess(site *object.Site) {
	gatewayCallbackSuccessTotal.WithLabelValues(siteLabel(site), domainLabel(site), authApplicationLabel(site)).Inc()
}

func observeGatewayCallbackFailure(site *object.Site, reason string) {
	gatewayCallbackFailureTotal.WithLabelValues(siteLabel(site), domainLabel(site), authApplicationLabel(site), reason).Inc()
}

func observeGatewayTokenInvalid(site *object.Site) {
	gatewayTokenInvalidTotal.WithLabelValues(siteLabel(site), domainLabel(site), authApplicationLabel(site)).Inc()
}

func observeGatewayInactiveSiteBlock(site *object.Site) {
	gatewayInactiveSiteBlockTotal.WithLabelValues(siteLabel(site), domainLabel(site)).Inc()
}

func observeGatewayUpstreamError(site *object.Site, upstream string) {
	gatewayUpstreamErrorTotal.WithLabelValues(siteLabel(site), domainLabel(site), upstream).Inc()
}

func siteLabel(site *object.Site) string {
	if site == nil {
		return ""
	}
	return site.GetId()
}

func domainLabel(site *object.Site) string {
	if site == nil {
		return ""
	}
	return site.Domain
}

func authApplicationLabel(site *object.Site) string {
	if site == nil {
		return ""
	}
	return site.AuthApplication
}
