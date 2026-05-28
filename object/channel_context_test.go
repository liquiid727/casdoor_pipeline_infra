package object

import (
	"strings"
	"testing"

	"github.com/golang-jwt/jwt/v5"
)

func TestNormalizeOrganizationChannelFields(t *testing.T) {
	organization := &Organization{Name: "root-org"}
	NormalizeOrganizationChannelFields(organization)

	if organization.OrganizationType != OrganizationTypeRoot {
		t.Fatalf("expected default organization type %q, got %q", OrganizationTypeRoot, organization.OrganizationType)
	}
	if organization.Status != OrganizationStatusActive {
		t.Fatalf("expected default status %q, got %q", OrganizationStatusActive, organization.Status)
	}
	if organization.ParentOrganization != "" {
		t.Fatalf("expected root organization to clear parent, got %q", organization.ParentOrganization)
	}
	if organization.ChannelMode != "" {
		t.Fatalf("expected root organization to clear channel mode, got %q", organization.ChannelMode)
	}

	channelOrganization := &Organization{
		Name:               "channel-a",
		OrganizationType:   OrganizationTypeChannel,
		ParentOrganization: "root-org",
	}
	NormalizeOrganizationChannelFields(channelOrganization)
	if channelOrganization.ChannelMode != ChannelModeSubTenant {
		t.Fatalf("expected default channel mode %q, got %q", ChannelModeSubTenant, channelOrganization.ChannelMode)
	}
}

func TestValidateOrganizationChannelFields(t *testing.T) {
	tests := []struct {
		name         string
		organization *Organization
		wantErr      string
	}{
		{
			name: "valid root organization",
			organization: &Organization{
				Name:             "root-org",
				OrganizationType: OrganizationTypeRoot,
				Status:           OrganizationStatusActive,
			},
		},
		{
			name: "valid shared-account channel organization",
			organization: &Organization{
				Name:               "channel-a",
				OrganizationType:   OrganizationTypeChannel,
				ParentOrganization: "root-org",
				ChannelMode:        ChannelModeSharedAccount,
				Status:             OrganizationStatusActive,
			},
		},
		{
			name: "channel organization requires parent",
			organization: &Organization{
				Name:             "channel-a",
				OrganizationType: OrganizationTypeChannel,
			},
			wantErr: "parentOrganization is required",
		},
		{
			name: "channel organization cannot parent itself",
			organization: &Organization{
				Name:               "channel-a",
				OrganizationType:   OrganizationTypeChannel,
				ParentOrganization: "channel-a",
				ChannelMode:        ChannelModeSubTenant,
			},
			wantErr: "parentOrganization cannot be the organization itself",
		},
		{
			name: "invalid organization type",
			organization: &Organization{
				Name:             "channel-a",
				OrganizationType: "unknown",
			},
			wantErr: "unsupported organizationType",
		},
		{
			name: "invalid organization status",
			organization: &Organization{
				Name:             "root-org",
				OrganizationType: OrganizationTypeRoot,
				Status:           "archived",
			},
			wantErr: "unsupported organization status",
		},
		{
			name: "invalid channel mode",
			organization: &Organization{
				Name:               "channel-a",
				OrganizationType:   OrganizationTypeChannel,
				ParentOrganization: "root-org",
				ChannelMode:        "hybrid",
			},
			wantErr: "unsupported channelMode",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateOrganizationChannelFields(tt.organization)
			if tt.wantErr == "" && err != nil {
				t.Fatalf("expected no error, got %v", err)
			}
			if tt.wantErr != "" {
				if err == nil {
					t.Fatalf("expected error containing %q, got nil", tt.wantErr)
				}
				if got := err.Error(); got == "" || !strings.Contains(got, tt.wantErr) {
					t.Fatalf("expected error containing %q, got %q", tt.wantErr, got)
				}
			}
		})
	}
}

func TestGetRootOrganizationName(t *testing.T) {
	rootOrganization := &Organization{Name: "root-org"}
	if got := GetRootOrganizationName(rootOrganization); got != "root-org" {
		t.Fatalf("expected root organization name %q, got %q", "root-org", got)
	}

	channelOrganization := &Organization{
		Name:               "channel-a",
		OrganizationType:   OrganizationTypeChannel,
		ParentOrganization: "root-org",
		ChannelMode:        ChannelModeSharedAccount,
	}
	if got := GetRootOrganizationName(channelOrganization); got != "root-org" {
		t.Fatalf("expected parent organization name %q, got %q", "root-org", got)
	}
}

func TestChannelContextValue(t *testing.T) {
	channelContext := &ChannelContext{
		ChannelOrganization: "channel-a",
		RootOrganization:    "root-org",
		ChannelMode:         ChannelModeSharedAccount,
	}

	if got := channelContextValue(channelContext, "channel"); got != "channel-a" {
		t.Fatalf("expected channel value %q, got %q", "channel-a", got)
	}
	if got := channelContextValue(channelContext, "root"); got != "root-org" {
		t.Fatalf("expected root value %q, got %q", "root-org", got)
	}
	if got := channelContextValue(channelContext, "mode"); got != ChannelModeSharedAccount {
		t.Fatalf("expected mode value %q, got %q", ChannelModeSharedAccount, got)
	}
	if got := channelContextValue(nil, "channel"); got != "" {
		t.Fatalf("expected nil channel context to return empty string, got %q", got)
	}
}

func TestClaimsChannelPropagation(t *testing.T) {
	claims := Claims{
		User: &User{
			Owner:         "root-org",
			Name:          "alice",
			Id:            "user-id",
			DisplayName:   "Alice",
			Email:         "alice@example.com",
			EmailVerified: true,
			Phone:         "13800138000",
		},
		TokenType:        "access-token",
		Nonce:            "nonce-1",
		Scope:            "openid profile email",
		Channel:          "channel-a",
		RootOrganization: "root-org",
		ChannelMode:      ChannelModeSharedAccount,
		SigninMethod:     "password",
		Provider:         "built-in",
		Azp:              "client-id",
		RegisteredClaims: jwt.RegisteredClaims{
			Subject: "user-id",
		},
	}

	shortClaims := getShortClaims(claims)
	if shortClaims.Channel != "channel-a" || shortClaims.RootOrganization != "root-org" || shortClaims.ChannelMode != ChannelModeSharedAccount {
		t.Fatalf("short claims lost channel context: %+v", shortClaims)
	}

	claimsWithoutThirdIdp := getClaimsWithoutThirdIdp(claims)
	if claimsWithoutThirdIdp.Channel != "channel-a" || claimsWithoutThirdIdp.RootOrganization != "root-org" || claimsWithoutThirdIdp.ChannelMode != ChannelModeSharedAccount {
		t.Fatalf("claims without third idp lost channel context: %+v", claimsWithoutThirdIdp)
	}

	standardClaims := getStandardClaims(claims)
	if standardClaims.Channel != "channel-a" || standardClaims.RootOrganization != "root-org" || standardClaims.ChannelMode != ChannelModeSharedAccount {
		t.Fatalf("standard claims lost channel context: %+v", standardClaims)
	}

	customClaims := getClaimsCustom(claims, []string{"Channel", "RootOrganization", "ChannelMode"}, nil)
	if got := customClaims["channel"]; got != "channel-a" {
		t.Fatalf("expected custom channel claim %q, got %#v", "channel-a", got)
	}
	if got := customClaims["rootOrganization"]; got != "root-org" {
		t.Fatalf("expected custom rootOrganization claim %q, got %#v", "root-org", got)
	}
	if got := customClaims["channelMode"]; got != ChannelModeSharedAccount {
		t.Fatalf("expected custom channelMode claim %q, got %#v", ChannelModeSharedAccount, got)
	}
}

func TestGetUserInfoIncludesChannelContext(t *testing.T) {
	user := &User{
		Id:    "user-id",
		Email: "alice@example.com",
	}
	channelContext := &ChannelContext{
		ChannelOrganization: "channel-a",
		RootOrganization:    "root-org",
		ChannelMode:         ChannelModeSharedAccount,
	}

	userInfo, err := GetUserInfo(user, "email", "client-id", "https://example.com", channelContext)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if userInfo.Channel != "channel-a" {
		t.Fatalf("expected channel %q, got %q", "channel-a", userInfo.Channel)
	}
	if userInfo.RootOrganization != "root-org" {
		t.Fatalf("expected root organization %q, got %q", "root-org", userInfo.RootOrganization)
	}
	if userInfo.ChannelMode != ChannelModeSharedAccount {
		t.Fatalf("expected channel mode %q, got %q", ChannelModeSharedAccount, userInfo.ChannelMode)
	}
}

func TestNormalizeChannelRelationAndUserChannelBinding(t *testing.T) {
	relation := &ChannelRelation{}
	NormalizeChannelRelation(relation)
	if relation.Status != OrganizationStatusActive {
		t.Fatalf("expected default relation status %q, got %q", OrganizationStatusActive, relation.Status)
	}
	if relation.ChannelMode != ChannelModeSubTenant {
		t.Fatalf("expected default relation mode %q, got %q", ChannelModeSubTenant, relation.ChannelMode)
	}
	if relation.Metadata == nil {
		t.Fatal("expected relation metadata to be initialized")
	}

	binding := &UserChannelBinding{}
	NormalizeUserChannelBinding(binding)
	if binding.Status != ChannelBindingStatusActive {
		t.Fatalf("expected default binding status %q, got %q", ChannelBindingStatusActive, binding.Status)
	}
	if binding.Roles == nil || binding.Permissions == nil {
		t.Fatal("expected binding roles and permissions to be initialized")
	}
}
