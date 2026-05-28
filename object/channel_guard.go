package object

import (
	"fmt"

	"github.com/casdoor/casdoor/util"
)

func getActiveChannelRelationDependencies(organizationName string) ([]*ChannelRelation, error) {
	relations := []*ChannelRelation{}
	err := ormer.Engine.Where("status = ? and (parent_organization = ? or channel_organization = ?)", OrganizationStatusActive, organizationName, organizationName).Find(&relations)
	if err != nil {
		return nil, err
	}

	for _, relation := range relations {
		NormalizeChannelRelation(relation)
	}

	return relations, nil
}

func getActiveUserChannelBindingDependencies(organizationName string) ([]*UserChannelBinding, error) {
	bindings := []*UserChannelBinding{}
	err := ormer.Engine.Where("status = ? and (root_organization = ? or channel_organization = ?)", ChannelBindingStatusActive, organizationName, organizationName).Find(&bindings)
	if err != nil {
		return nil, err
	}

	for _, binding := range bindings {
		NormalizeUserChannelBinding(binding)
	}

	return bindings, nil
}

func ensureNoActiveChannelDependencies(organizationName, action string) error {
	relations, err := getActiveChannelRelationDependencies(organizationName)
	if err != nil {
		return err
	}
	if len(relations) > 0 {
		return fmt.Errorf("cannot %s organization %s while active channel relation exists: %s", action, organizationName, util.GetId(relations[0].Owner, relations[0].Name))
	}

	bindings, err := getActiveUserChannelBindingDependencies(organizationName)
	if err != nil {
		return err
	}
	if len(bindings) > 0 {
		return fmt.Errorf("cannot %s organization %s while active user channel binding exists: %s", action, organizationName, util.GetId(bindings[0].Owner, bindings[0].Name))
	}

	return nil
}

func ValidateOrganizationChannelMutation(oldOrganization, newOrganization *Organization) error {
	if oldOrganization == nil || newOrganization == nil {
		return nil
	}

	NormalizeOrganizationChannelFields(oldOrganization)
	NormalizeOrganizationChannelFields(newOrganization)

	requiresGuard := false
	action := "change"

	if oldOrganization.Status != OrganizationStatusDisabled && newOrganization.Status == OrganizationStatusDisabled {
		requiresGuard = true
		action = "disable"
	}

	if oldOrganization.OrganizationType != newOrganization.OrganizationType ||
		oldOrganization.ParentOrganization != newOrganization.ParentOrganization ||
		oldOrganization.ChannelMode != newOrganization.ChannelMode {
		requiresGuard = true
		if action == "change" {
			action = "change channel structure for"
		}
	}

	if !requiresGuard {
		return nil
	}

	return ensureNoActiveChannelDependencies(oldOrganization.Name, action)
}
