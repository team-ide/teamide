package module_toolbox

import (
	"github.com/gin-gonic/gin"
	"teamide/internal/base"
)

type InsertGroupRequest struct {
	*ToolboxGroupModel
}

type InsertGroupResponse struct {
}

func (this_ *ToolboxApi) insertGroup(requestBean *base.RequestBean, c *gin.Context) (res interface{}, err error) {

	request := &InsertGroupRequest{}
	if !base.RequestJSON(request, c) {
		return
	}
	response := &InsertGroupResponse{}

	bean := request.ToolboxGroupModel
	bean.UserId = requestBean.JWT.UserId

	_, err = this_.ToolboxService.InsertGroup(bean)
	if err != nil {
		return
	}

	res = response
	return
}

type UpdateGroupRequest struct {
	*ToolboxGroupModel
}

type UpdateGroupResponse struct {
}

func (this_ *ToolboxApi) updateGroup(requestBean *base.RequestBean, c *gin.Context) (res interface{}, err error) {

	request := &UpdateGroupRequest{}
	if !base.RequestJSON(request, c) {
		return
	}
	response := &UpdateGroupResponse{}

	_, err = this_.ToolboxService.UpdateGroup(request.ToolboxGroupModel)
	if err != nil {
		return
	}

	res = response
	return
}

type DeleteGroupRequest struct {
	*ToolboxGroupModel
}

type DeleteGroupResponse struct {
}

func (this_ *ToolboxApi) deleteGroup(requestBean *base.RequestBean, c *gin.Context) (res interface{}, err error) {

	request := &DeleteGroupRequest{}
	if !base.RequestJSON(request, c) {
		return
	}
	response := &DeleteGroupResponse{}

	_, err = this_.ToolboxService.DeleteGroup(request.ToolboxGroupModel.GroupId)
	if err != nil {
		return
	}

	res = response
	return
}
