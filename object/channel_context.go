package object

import "fmt"

const (
	OrganizationTypeRoot    = "root_org"
	OrganizationTypeChannel = "channel_org"

	ChannelModeSharedAccount = "shared_account"
	ChannelModeSubTenant     = "sub_tenant"

	OrganizationStatusActive   = "active"
	OrganizationStatusDisabled = "disabled"

	ChannelBindingStatusActive   = "active"
	ChannelBindingStatusDisabled = "disabled"
)

type ChannelContext struct {
	ChannelOrganization string `json:"channelOrganization"`
	RootOrganization    string `json:"rootOrganization"`
	ChannelMode         string `json:"channelMode"`
}

func NormalizeOrganizationChannelFields(organization *Organization) {
	if organization == nil {
		return
	}

	if organization.OrganizationType == "" {
		organization.OrganizationType = OrganizationTypeRoot
	}
	if organization.Status == "" {
		organization.Status = OrganizationStatusActive
	}

	if organization.OrganizationType != OrganizationTypeChannel {
		organization.ParentOrganization = ""
		organization.ChannelMode = ""
		return
	}

	if organization.ChannelMode == "" {
		organization.ChannelMode = ChannelModeSubTenant
	}
}

func ValidateOrganizationChannelFields(organization *Organization) error {
	if organization == nil {
		return nil
	}

	NormalizeOrganizationChannelFields(organization)

	switch organization.OrganizationType {
	case OrganizationTypeRoot, OrganizationTypeChannel:
	default:
		return fmt.Errorf("unsupported organizationType: %s", organization.OrganizationType)
	}

	if organization.Status != OrganizationStatusActive && organization.Status != OrganizationStatusDisabled {
		return fmt.Errorf("unsupported organization status: %s", organization.Status)
	}

	if organization.OrganizationType == OrganizationTypeChannel {
		if organization.ParentOrganization == "" {
			return fmt.Errorf("parentOrganization is required for channel organizations")
		}
		if organization.ParentOrganization == organization.Name {
			return fmt.Errorf("parentOrganization cannot be the organization itself")
		}
		if organization.ChannelMode != ChannelModeSharedAccount && organization.ChannelMode != ChannelModeSubTenant {
			return fmt.Errorf("unsupported channelMode: %s", organization.ChannelMode)
		}
	}

	return nil
}

func (organization *Organization) IsActive() bool {
	if organization == nil {
		return false
	}

	NormalizeOrganizationChannelFields(organization)
	return organization.Status != OrganizationStatusDisabled
}

func (organization *Organization) IsChannelOrganization() bool {
	if organization == nil {
		return false
	}

	NormalizeOrganizationChannelFields(organization)
	return organization.OrganizationType == OrganizationTypeChannel
}

func (organization *Organization) IsSharedAccountChannel() bool {
	return organization.IsChannelOrganization() && organization.ChannelMode == ChannelModeSharedAccount
}

func (organization *Organization) IsSubTenantChannel() bool {
	return organization.IsChannelOrganization() && organization.ChannelMode == ChannelModeSubTenant
}

func GetRootOrganizationName(organization *Organization) string {
	if organization == nil {
		return ""
	}

	NormalizeOrganizationChannelFields(organization)
	if organization.IsChannelOrganization() && organization.ParentOrganization != "" {
		return organization.ParentOrganization
	}

	return organization.Name
}
