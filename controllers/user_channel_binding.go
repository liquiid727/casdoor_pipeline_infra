package controllers

import (
	"encoding/json"

	"github.com/beego/beego/v2/core/utils/pagination"
	"github.com/casdoor/casdoor/object"
	"github.com/casdoor/casdoor/util"
)

func (c *ApiController) GetUserChannels() {
	user, ok := c.RequireSignedInUser()
	if !ok {
		return
	}

	isGlobalAdmin := user.IsGlobalAdmin()
	rootOrganization := c.Ctx.Input.Query("rootOrganization")
	userName := c.Ctx.Input.Query("user")
	owner := c.Ctx.Input.Query("owner")
	limit := c.Ctx.Input.Query("pageSize")
	page := c.Ctx.Input.Query("p")
	field := c.Ctx.Input.Query("field")
	value := c.Ctx.Input.Query("value")
	sortField := c.Ctx.Input.Query("sortField")
	sortOrder := c.Ctx.Input.Query("sortOrder")

	if !isGlobalAdmin && !user.IsAdmin {
		rootOrganization = user.Owner
		userName = user.Name
	} else if !isGlobalAdmin && rootOrganization == "" {
		rootOrganization = user.Owner
	}

	if limit == "" || page == "" {
		bindings, err := object.GetUserChannelBindings(owner, rootOrganization, userName)
		if err != nil {
			c.ResponseError(err.Error())
			return
		}
		c.ResponseOk(bindings)
		return
	}

	limitInt := util.ParseInt(limit)
	count, err := object.GetUserChannelBindingCount(owner, field, value)
	if err != nil {
		c.ResponseError(err.Error())
		return
	}

	paginator := pagination.NewPaginator(c.Ctx.Request, limitInt, count)
	bindings, err := object.GetPaginationUserChannelBindings(owner, rootOrganization, userName, paginator.Offset(), limitInt, field, value, sortField, sortOrder)
	if err != nil {
		c.ResponseError(err.Error())
		return
	}
	c.ResponseOk(bindings, paginator.Nums())
}

func (c *ApiController) GetUserChannel() {
	_, ok := c.RequireAdmin()
	if !ok {
		return
	}

	id := c.Ctx.Input.Query("id")
	binding, err := object.GetUserChannelBinding(id)
	if err != nil {
		c.ResponseError(err.Error())
		return
	}
	c.ResponseOk(binding)
}

func (c *ApiController) BindUserChannel() {
	adminOrganization, ok := c.RequireAdmin()
	if !ok {
		return
	}

	var binding object.UserChannelBinding
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &binding); err != nil {
		c.ResponseError(err.Error())
		return
	}

	if adminOrganization != "" {
		binding.RootOrganization = adminOrganization
	}

	c.Data["json"] = wrapActionResponse(object.AddUserChannelBinding(&binding))
	c.ServeJSON()
}

func (c *ApiController) UpdateUserChannel() {
	adminOrganization, ok := c.RequireAdmin()
	if !ok {
		return
	}

	id := c.Ctx.Input.Query("id")
	var binding object.UserChannelBinding
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &binding); err != nil {
		c.ResponseError(err.Error())
		return
	}

	if adminOrganization != "" {
		binding.RootOrganization = adminOrganization
	}

	c.Data["json"] = wrapActionResponse(object.UpdateUserChannelBinding(id, &binding))
	c.ServeJSON()
}

func (c *ApiController) UnbindUserChannel() {
	adminOrganization, ok := c.RequireAdmin()
	if !ok {
		return
	}

	var binding object.UserChannelBinding
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &binding); err != nil {
		c.ResponseError(err.Error())
		return
	}

	if adminOrganization != "" && binding.RootOrganization != "" && binding.RootOrganization != adminOrganization {
		c.ResponseError(c.T("auth:Unauthorized operation"))
		return
	}

	c.Data["json"] = wrapActionResponse(object.DeleteUserChannelBinding(&binding))
	c.ServeJSON()
}
