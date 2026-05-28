package object

import (
	"fmt"

	"github.com/casdoor/casdoor/util"
	"github.com/xorm-io/core"
)

type ChannelRelation struct {
	Owner       string `xorm:"varchar(100) notnull pk" json:"owner"`
	Name        string `xorm:"varchar(100) notnull pk" json:"name"`
	CreatedTime string `xorm:"varchar(100)" json:"createdTime"`

	ParentOrganization  string            `xorm:"varchar(100)" json:"parentOrganization"`
	ChannelOrganization string            `xorm:"varchar(100)" json:"channelOrganization"`
	ChannelMode         string            `xorm:"varchar(100)" json:"channelMode"`
	SettlementRule      string            `xorm:"varchar(100)" json:"settlementRule"`
	SettlementRatio     float64           `json:"settlementRatio"`
	SettlementAccount   string            `xorm:"varchar(100)" json:"settlementAccount"`
	Status              string            `xorm:"varchar(100)" json:"status"`
	Metadata            map[string]string `xorm:"json" json:"metadata"`
}

func NormalizeChannelRelation(relation *ChannelRelation) {
	if relation == nil {
		return
	}

	if relation.Status == "" {
		relation.Status = OrganizationStatusActive
	}
	if relation.ChannelMode == "" {
		relation.ChannelMode = ChannelModeSubTenant
	}
	if relation.Metadata == nil {
		relation.Metadata = map[string]string{}
	}
}

func ValidateChannelRelation(relation *ChannelRelation) error {
	if relation == nil {
		return nil
	}

	NormalizeChannelRelation(relation)
	if relation.ParentOrganization == "" {
		return fmt.Errorf("parentOrganization is required")
	}
	if relation.ChannelOrganization == "" {
		return fmt.Errorf("channelOrganization is required")
	}
	if relation.ParentOrganization == relation.ChannelOrganization {
		return fmt.Errorf("parentOrganization and channelOrganization must be different")
	}
	if relation.ChannelMode != ChannelModeSharedAccount && relation.ChannelMode != ChannelModeSubTenant {
		return fmt.Errorf("unsupported channelMode: %s", relation.ChannelMode)
	}
	if relation.Status != OrganizationStatusActive && relation.Status != OrganizationStatusDisabled {
		return fmt.Errorf("unsupported status: %s", relation.Status)
	}

	parentOrg, err := getOrganization("admin", relation.ParentOrganization)
	if err != nil {
		return err
	}
	if parentOrg == nil {
		return fmt.Errorf("parent organization does not exist: %s", relation.ParentOrganization)
	}

	channelOrg, err := getOrganization("admin", relation.ChannelOrganization)
	if err != nil {
		return err
	}
	if channelOrg == nil {
		return fmt.Errorf("channel organization does not exist: %s", relation.ChannelOrganization)
	}

	NormalizeOrganizationChannelFields(channelOrg)
	if channelOrg.OrganizationType != OrganizationTypeChannel {
		return fmt.Errorf("channel organization must have organizationType=channel_org")
	}
	if channelOrg.ParentOrganization != relation.ParentOrganization {
		return fmt.Errorf("channel organization parent mismatch: expected %s, got %s", relation.ParentOrganization, channelOrg.ParentOrganization)
	}
	if channelOrg.ChannelMode != relation.ChannelMode {
		return fmt.Errorf("channel organization mode mismatch: expected %s, got %s", relation.ChannelMode, channelOrg.ChannelMode)
	}

	return nil
}

func GetChannelRelationCount(owner, field, value string) (int64, error) {
	session := GetSession(owner, -1, -1, field, value, "", "")
	return session.Count(&ChannelRelation{})
}

func GetChannelRelations(owner, parentOrganization string) ([]*ChannelRelation, error) {
	relations := []*ChannelRelation{}
	session := ormer.Engine.Desc("created_time")
	hasCondition := false
	if owner != "" {
		session = session.Where("owner = ?", owner)
		hasCondition = true
	}
	if parentOrganization != "" {
		if hasCondition {
			session = session.And("parent_organization = ?", parentOrganization)
		} else {
			session = session.Where("parent_organization = ?", parentOrganization)
		}
	}
	err := session.Find(&relations)
	return relations, err
}

func GetPaginationChannelRelations(owner, parentOrganization string, offset, limit int, field, value, sortField, sortOrder string) ([]*ChannelRelation, error) {
	relations := []*ChannelRelation{}
	session := GetSession(owner, offset, limit, field, value, sortField, sortOrder)
	if parentOrganization != "" {
		session = session.Where("parent_organization = ?", parentOrganization)
	}
	err := session.Find(&relations)
	return relations, err
}

func getChannelRelation(owner, name string) (*ChannelRelation, error) {
	if owner == "" || name == "" {
		return nil, nil
	}

	relation := ChannelRelation{Owner: owner, Name: name}
	existed, err := ormer.Engine.Get(&relation)
	if err != nil {
		return nil, err
	}
	if existed {
		NormalizeChannelRelation(&relation)
		return &relation, nil
	}

	return nil, nil
}

func GetChannelRelation(id string) (*ChannelRelation, error) {
	owner, name, err := util.GetOwnerAndNameFromIdWithError(id)
	if err != nil {
		return nil, err
	}
	return getChannelRelation(owner, name)
}

func AddChannelRelation(relation *ChannelRelation) (bool, error) {
	if relation.Owner == "" {
		relation.Owner = "admin"
	}
	if relation.CreatedTime == "" {
		relation.CreatedTime = util.GetCurrentTime()
	}
	if err := ValidateChannelRelation(relation); err != nil {
		return false, err
	}

	affected, err := ormer.Engine.Insert(relation)
	if err != nil {
		return false, err
	}
	return affected != 0, nil
}

func UpdateChannelRelation(id string, relation *ChannelRelation) (bool, error) {
	owner, name, err := util.GetOwnerAndNameFromIdWithError(id)
	if err != nil {
		return false, err
	}
	oldRelation, err := getChannelRelation(owner, name)
	if err != nil {
		return false, err
	}
	if oldRelation == nil {
		return false, nil
	}

	if err = ValidateChannelRelation(relation); err != nil {
		return false, err
	}

	affected, err := ormer.Engine.ID(core.PK{owner, name}).AllCols().Update(relation)
	if err != nil {
		return false, err
	}
	return affected != 0, nil
}

func DeleteChannelRelation(relation *ChannelRelation) (bool, error) {
	affected, err := ormer.Engine.ID(core.PK{relation.Owner, relation.Name}).Delete(&ChannelRelation{})
	if err != nil {
		return false, err
	}
	return affected != 0, nil
}

func GetActiveChannelRelation(parentOrganization, channelOrganization string) (*ChannelRelation, error) {
	relation := &ChannelRelation{}
	has, err := ormer.Engine.Where("parent_organization = ? and channel_organization = ? and status = ?", parentOrganization, channelOrganization, OrganizationStatusActive).Get(relation)
	if err != nil {
		return nil, err
	}
	if !has {
		return nil, nil
	}
	NormalizeChannelRelation(relation)
	return relation, nil
}
