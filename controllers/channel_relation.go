package controllers

import (
	"encoding/json"

	"github.com/beego/beego/v2/core/utils/pagination"
	"github.com/casdoor/casdoor/object"
	"github.com/casdoor/casdoor/util"
)

func (c *ApiController) GetChannelRelations() {
	adminOrganization, ok := c.RequireAdmin()
	if !ok {
		return
	}

	owner := c.Ctx.Input.Query("owner")
	parentOrganization := c.Ctx.Input.Query("parentOrganization")
	limit := c.Ctx.Input.Query("pageSize")
	page := c.Ctx.Input.Query("p")
	field := c.Ctx.Input.Query("field")
	value := c.Ctx.Input.Query("value")
	sortField := c.Ctx.Input.Query("sortField")
	sortOrder := c.Ctx.Input.Query("sortOrder")

	if adminOrganization != "" && parentOrganization == "" {
		parentOrganization = adminOrganization
	}

	if limit == "" || page == "" {
		relations, err := object.GetChannelRelations(owner, parentOrganization)
		if err != nil {
			c.ResponseError(err.Error())
			return
		}
		c.ResponseOk(relations)
		return
	}

	limitInt := util.ParseInt(limit)
	count, err := object.GetChannelRelationCount(owner, field, value)
	if err != nil {
		c.ResponseError(err.Error())
		return
	}

	paginator := pagination.NewPaginator(c.Ctx.Request, limitInt, count)
	relations, err := object.GetPaginationChannelRelations(owner, parentOrganization, paginator.Offset(), limitInt, field, value, sortField, sortOrder)
	if err != nil {
		c.ResponseError(err.Error())
		return
	}
	c.ResponseOk(relations, paginator.Nums())
}

func (c *ApiController) GetChannelRelation() {
	_, ok := c.RequireAdmin()
	if !ok {
		return
	}

	id := c.Ctx.Input.Query("id")
	relation, err := object.GetChannelRelation(id)
	if err != nil {
		c.ResponseError(err.Error())
		return
	}
	c.ResponseOk(relation)
}

func (c *ApiController) AddChannelRelation() {
	adminOrganization, ok := c.RequireAdmin()
	if !ok {
		return
	}

	var relation object.ChannelRelation
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &relation); err != nil {
		c.ResponseError(err.Error())
		return
	}

	if adminOrganization != "" {
		relation.ParentOrganization = adminOrganization
	}

	c.Data["json"] = wrapActionResponse(object.AddChannelRelation(&relation))
	c.ServeJSON()
}

func (c *ApiController) UpdateChannelRelation() {
	adminOrganization, ok := c.RequireAdmin()
	if !ok {
		return
	}

	id := c.Ctx.Input.Query("id")
	var relation object.ChannelRelation
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &relation); err != nil {
		c.ResponseError(err.Error())
		return
	}

	if adminOrganization != "" {
		relation.ParentOrganization = adminOrganization
	}

	c.Data["json"] = wrapActionResponse(object.UpdateChannelRelation(id, &relation))
	c.ServeJSON()
}

func (c *ApiController) DeleteChannelRelation() {
	adminOrganization, ok := c.RequireAdmin()
	if !ok {
		return
	}

	var relation object.ChannelRelation
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &relation); err != nil {
		c.ResponseError(err.Error())
		return
	}

	if adminOrganization != "" && relation.ParentOrganization != "" && relation.ParentOrganization != adminOrganization {
		c.ResponseError(c.T("auth:Unauthorized operation"))
		return
	}

	c.Data["json"] = wrapActionResponse(object.DeleteChannelRelation(&relation))
	c.ServeJSON()
}
