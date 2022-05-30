package module_toolbox

import (
	"github.com/gin-gonic/gin"
	"teamide/internal/base"
)

type QueryQuickCommandRequest struct {
	*ToolboxQuickCommandModel
}

type QueryQuickCommandResponse struct {
	QuickCommands []*ToolboxQuickCommandModel `json:"quickCommands,omitempty"`
}

func (this_ *ToolboxApi) queryQuickCommand(requestBean *base.RequestBean, c *gin.Context) (res interface{}, err error) {

	request := &QueryQuickCommandRequest{}
	if !base.RequestJSON(request, c) {
		return
	}
	response := &QueryQuickCommandResponse{}

	bean := request.ToolboxQuickCommandModel
	if bean == nil {
		bean = &ToolboxQuickCommandModel{}
	}
	bean.UserId = requestBean.JWT.UserId

	response.QuickCommands, err = this_.ToolboxService.QueryQuickCommand(bean)
	if err != nil {
		return
	}

	res = response
	return
}

type InsertQuickCommandRequest struct {
	*ToolboxQuickCommandModel
}

type InsertQuickCommandResponse struct {
}

func (this_ *ToolboxApi) insertQuickCommand(requestBean *base.RequestBean, c *gin.Context) (res interface{}, err error) {

	request := &InsertQuickCommandRequest{}
	if !base.RequestJSON(request, c) {
		return
	}
	response := &InsertQuickCommandResponse{}

	bean := request.ToolboxQuickCommandModel
	bean.UserId = requestBean.JWT.UserId

	_, err = this_.ToolboxService.InsertQuickCommand(bean)
	if err != nil {
		return
	}

	res = response
	return
}

type UpdateQuickCommandRequest struct {
	*ToolboxQuickCommandModel
}

type UpdateQuickCommandResponse struct {
}

func (this_ *ToolboxApi) updateQuickCommand(requestBean *base.RequestBean, c *gin.Context) (res interface{}, err error) {

	request := &UpdateQuickCommandRequest{}
	if !base.RequestJSON(request, c) {
		return
	}
	response := &UpdateQuickCommandResponse{}

	_, err = this_.ToolboxService.UpdateQuickCommand(request.ToolboxQuickCommandModel)
	if err != nil {
		return
	}

	res = response
	return
}

type DeleteQuickCommandRequest struct {
	*ToolboxQuickCommandModel
}

type DeleteQuickCommandResponse struct {
}

func (this_ *ToolboxApi) deleteQuickCommand(requestBean *base.RequestBean, c *gin.Context) (res interface{}, err error) {

	request := &DeleteQuickCommandRequest{}
	if !base.RequestJSON(request, c) {
		return
	}
	response := &DeleteQuickCommandResponse{}

	_, err = this_.ToolboxService.DeleteQuickCommand(request.ToolboxQuickCommandModel.QuickCommandId)
	if err != nil {
		return
	}

	res = response
	return
}
