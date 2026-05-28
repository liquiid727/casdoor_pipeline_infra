package object

import "fmt"

func ResolveChannelContextForLogin(application *Application, user *User, requestedChannel string) (*ChannelContext, error) {
	if user == nil {
		return nil, nil
	}

	userOrg, err := GetOrganizationByUser(user)
	if err != nil {
		return nil, err
	}
	if userOrg == nil {
		return nil, nil
	}
	NormalizeOrganizationChannelFields(userOrg)

	// Sub-tenant users carry their channel context from their own organization.
	if userOrg.IsSubTenantChannel() {
		if application != nil && application.Organization != "" && application.Organization != user.Owner {
			return nil, fmt.Errorf("sub-tenant user cannot sign in to an application outside its organization")
		}
		return &ChannelContext{
			ChannelOrganization: userOrg.Name,
			RootOrganization:    GetRootOrganizationName(userOrg),
			ChannelMode:         ChannelModeSubTenant,
		}, nil
	}

	if requestedChannel == "" {
		bindings, err := GetActiveUserChannelBindings(user.Owner, user.Name)
		if err != nil {
			return nil, err
		}
		switch len(bindings) {
		case 0:
			return nil, nil
		case 1:
			requestedChannel = bindings[0].ChannelOrganization
		default:
			return nil, fmt.Errorf("channel selection required")
		}
	}

	channelOrg, err := getOrganization("admin", requestedChannel)
	if err != nil {
		return nil, err
	}
	if channelOrg == nil {
		return nil, fmt.Errorf("channel organization does not exist: %s", requestedChannel)
	}
	NormalizeOrganizationChannelFields(channelOrg)

	if !channelOrg.IsChannelOrganization() {
		return nil, fmt.Errorf("requested channel is not a channel organization: %s", requestedChannel)
	}
	if !channelOrg.IsActive() {
		return nil, fmt.Errorf("requested channel organization is disabled: %s", requestedChannel)
	}

	relation, err := GetActiveChannelRelation(channelOrg.ParentOrganization, channelOrg.Name)
	if err != nil {
		return nil, err
	}
	if relation == nil {
		return nil, fmt.Errorf("active channel relation does not exist for: %s", requestedChannel)
	}

	if channelOrg.IsSubTenantChannel() {
		if user.Owner != channelOrg.Name {
			return nil, fmt.Errorf("user does not belong to the requested sub-tenant channel")
		}
		return &ChannelContext{
			ChannelOrganization: channelOrg.Name,
			RootOrganization:    relation.ParentOrganization,
			ChannelMode:         ChannelModeSubTenant,
		}, nil
	}

	if relation.ParentOrganization != user.Owner {
		return nil, fmt.Errorf("user does not belong to the root organization for this channel")
	}

	binding, err := GetActiveUserChannelBinding(user.Owner, user.Name, channelOrg.Name)
	if err != nil {
		return nil, err
	}
	if binding == nil {
		return nil, fmt.Errorf("user is not bound to the requested channel")
	}

	if application != nil && application.Organization != "" && application.Organization != user.Owner {
		return nil, fmt.Errorf("shared-account channel login must use an application bound to the root organization")
	}

	return &ChannelContext{
		ChannelOrganization: channelOrg.Name,
		RootOrganization:    user.Owner,
		ChannelMode:         ChannelModeSharedAccount,
	}, nil
}
