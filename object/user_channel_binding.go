package object

import (
	"fmt"

	"github.com/casdoor/casdoor/util"
	"github.com/xorm-io/core"
)

type UserChannelBinding struct {
	Owner       string `xorm:"varchar(100) notnull pk" json:"owner"`
	Name        string `xorm:"varchar(100) notnull pk" json:"name"`
	CreatedTime string `xorm:"varchar(100)" json:"createdTime"`

	User                string   `xorm:"varchar(100)" json:"user"`
	RootOrganization    string   `xorm:"varchar(100)" json:"rootOrganization"`
	ChannelOrganization string   `xorm:"varchar(100)" json:"channelOrganization"`
	Roles               []string `xorm:"json" json:"roles"`
	Permissions         []string `xorm:"json" json:"permissions"`
	Status              string   `xorm:"varchar(100)" json:"status"`
}

func NormalizeUserChannelBinding(binding *UserChannelBinding) {
	if binding == nil {
		return
	}
	if binding.Status == "" {
		binding.Status = ChannelBindingStatusActive
	}
	if binding.Roles == nil {
		binding.Roles = []string{}
	}
	if binding.Permissions == nil {
		binding.Permissions = []string{}
	}
}

func ValidateUserChannelBinding(binding *UserChannelBinding) error {
	if binding == nil {
		return nil
	}

	NormalizeUserChannelBinding(binding)
	if binding.User == "" {
		return fmt.Errorf("user is required")
	}
	if binding.RootOrganization == "" {
		return fmt.Errorf("rootOrganization is required")
	}
	if binding.ChannelOrganization == "" {
		return fmt.Errorf("channelOrganization is required")
	}
	if binding.RootOrganization == binding.ChannelOrganization {
		return fmt.Errorf("rootOrganization and channelOrganization must be different")
	}
	if binding.Status != ChannelBindingStatusActive && binding.Status != ChannelBindingStatusDisabled {
		return fmt.Errorf("unsupported status: %s", binding.Status)
	}

	user, err := getUser(binding.RootOrganization, binding.User)
	if err != nil {
		return err
	}
	if user == nil {
		return fmt.Errorf("user does not exist under root organization: %s/%s", binding.RootOrganization, binding.User)
	}

	channelOrg, err := getOrganization("admin", binding.ChannelOrganization)
	if err != nil {
		return err
	}
	if channelOrg == nil {
		return fmt.Errorf("channel organization does not exist: %s", binding.ChannelOrganization)
	}
	NormalizeOrganizationChannelFields(channelOrg)
	if !channelOrg.IsSharedAccountChannel() {
		return fmt.Errorf("user channel bindings only support shared_account channel organizations")
	}
	if channelOrg.ParentOrganization != binding.RootOrganization {
		return fmt.Errorf("channel organization parent mismatch: expected %s, got %s", binding.RootOrganization, channelOrg.ParentOrganization)
	}

	relation, err := GetActiveChannelRelation(binding.RootOrganization, binding.ChannelOrganization)
	if err != nil {
		return err
	}
	if relation == nil {
		return fmt.Errorf("active channel relation does not exist for %s -> %s", binding.RootOrganization, binding.ChannelOrganization)
	}

	return nil
}

func GetUserChannelBindingCount(owner, field, value string) (int64, error) {
	session := GetSession(owner, -1, -1, field, value, "", "")
	return session.Count(&UserChannelBinding{})
}

func GetUserChannelBindings(owner, rootOrganization, userName string) ([]*UserChannelBinding, error) {
	bindings := []*UserChannelBinding{}
	session := ormer.Engine.Desc("created_time")
	hasCondition := false
	if owner != "" {
		session = session.Where("owner = ?", owner)
		hasCondition = true
	}
	if rootOrganization != "" {
		if hasCondition {
			session = session.And("root_organization = ?", rootOrganization)
		} else {
			session = session.Where("root_organization = ?", rootOrganization)
			hasCondition = true
		}
	}
	if userName != "" {
		if hasCondition {
			session = session.And("user = ?", userName)
		} else {
			session = session.Where("user = ?", userName)
		}
	}
	err := session.Find(&bindings)
	return bindings, err
}

func GetPaginationUserChannelBindings(owner, rootOrganization, userName string, offset, limit int, field, value, sortField, sortOrder string) ([]*UserChannelBinding, error) {
	bindings := []*UserChannelBinding{}
	session := GetSession(owner, offset, limit, field, value, sortField, sortOrder)
	if rootOrganization != "" {
		session = session.Where("root_organization = ?", rootOrganization)
	}
	if userName != "" {
		session = session.And("user = ?", userName)
	}
	err := session.Find(&bindings)
	return bindings, err
}

func getUserChannelBinding(owner, name string) (*UserChannelBinding, error) {
	if owner == "" || name == "" {
		return nil, nil
	}

	binding := UserChannelBinding{Owner: owner, Name: name}
	existed, err := ormer.Engine.Get(&binding)
	if err != nil {
		return nil, err
	}
	if existed {
		NormalizeUserChannelBinding(&binding)
		return &binding, nil
	}

	return nil, nil
}

func GetUserChannelBinding(id string) (*UserChannelBinding, error) {
	owner, name, err := util.GetOwnerAndNameFromIdWithError(id)
	if err != nil {
		return nil, err
	}
	return getUserChannelBinding(owner, name)
}

func AddUserChannelBinding(binding *UserChannelBinding) (bool, error) {
	if binding.Owner == "" {
		binding.Owner = "admin"
	}
	if binding.CreatedTime == "" {
		binding.CreatedTime = util.GetCurrentTime()
	}
	if err := ValidateUserChannelBinding(binding); err != nil {
		return false, err
	}

	affected, err := ormer.Engine.Insert(binding)
	if err != nil {
		return false, err
	}
	return affected != 0, nil
}

func UpdateUserChannelBinding(id string, binding *UserChannelBinding) (bool, error) {
	owner, name, err := util.GetOwnerAndNameFromIdWithError(id)
	if err != nil {
		return false, err
	}
	oldBinding, err := getUserChannelBinding(owner, name)
	if err != nil {
		return false, err
	}
	if oldBinding == nil {
		return false, nil
	}

	if err = ValidateUserChannelBinding(binding); err != nil {
		return false, err
	}

	affected, err := ormer.Engine.ID(core.PK{owner, name}).AllCols().Update(binding)
	if err != nil {
		return false, err
	}
	return affected != 0, nil
}

func DeleteUserChannelBinding(binding *UserChannelBinding) (bool, error) {
	affected, err := ormer.Engine.ID(core.PK{binding.Owner, binding.Name}).Delete(&UserChannelBinding{})
	if err != nil {
		return false, err
	}
	return affected != 0, nil
}

func GetActiveUserChannelBindings(rootOrganization, userName string) ([]*UserChannelBinding, error) {
	bindings := []*UserChannelBinding{}
	err := ormer.Engine.Where("root_organization = ? and user = ? and status = ?", rootOrganization, userName, ChannelBindingStatusActive).Asc("created_time").Find(&bindings)
	if err != nil {
		return nil, err
	}
	for _, binding := range bindings {
		NormalizeUserChannelBinding(binding)
	}
	return bindings, nil
}

func GetActiveUserChannelBinding(rootOrganization, userName, channelOrganization string) (*UserChannelBinding, error) {
	binding := &UserChannelBinding{}
	has, err := ormer.Engine.Where("root_organization = ? and user = ? and channel_organization = ? and status = ?", rootOrganization, userName, channelOrganization, ChannelBindingStatusActive).Get(binding)
	if err != nil {
		return nil, err
	}
	if !has {
		return nil, nil
	}
	NormalizeUserChannelBinding(binding)
	return binding, nil
}
