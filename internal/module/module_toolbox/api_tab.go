package module_toolbox

import (
	"github.com/gin-gonic/gin"
	"teamide/internal/base"
)

type OpenTabRequest struct {
	*ToolboxOpenTabModel
}

type OpenTabResponse struct {
	Tab *ToolboxOpenTabModel `json:"tab,omitempty"`
}

func (this_ *ToolboxApi) openTab(requestBean *base.RequestBean, c *gin.Context) (res interface{}, err error) {

	request := &OpenTabRequest{}
	if !base.RequestJSON(request, c) {
		return
	}
	response := &OpenTabResponse{}

	toolboxOpenTab := request.ToolboxOpenTabModel
	toolboxOpenTab.UserId = requestBean.JWT.UserId

	_, err = this_.ToolboxService.OpenTab(toolboxOpenTab)
	if err != nil {
		return
	}

	response.Tab = toolboxOpenTab

	res = response
	return
}

type QueryOpenTabsRequest struct {
	*ToolboxOpenTabModel
}

type QueryOpenTabsResponse struct {
	Tabs []*ToolboxOpenTabModel `json:"tabs,omitempty"`
}

func (this_ *ToolboxApi) queryOpenTabs(requestBean *base.RequestBean, c *gin.Context) (res interface{}, err error) {

	request := &QueryOpenTabsRequest{}
	if !base.RequestJSON(request, c) {
		return
	}
	response := &QueryOpenTabsResponse{}

	response.Tabs, err = this_.ToolboxService.QueryOpenTabs(request.OpenId)
	if err != nil {
		return
	}

	res = response
	return
}

type CloseTabRequest struct {
	*ToolboxOpenTabModel
}

type CloseTabResponse struct {
}

func (this_ *ToolboxApi) closeTab(requestBean *base.RequestBean, c *gin.Context) (res interface{}, err error) {

	request := &CloseTabRequest{}
	if !base.RequestJSON(request, c) {
		return
	}
	response := &CloseTabResponse{}

	_, err = this_.ToolboxService.CloseTab(request.TabId)
	if err != nil {
		return
	}

	res = response
	return
}

type UpdateOpenTabExtendRequest struct {
	*ToolboxOpenTabModel
}

type UpdateOpenTabExtendResponse struct {
}

func (this_ *ToolboxApi) updateOpenTabExtend(requestBean *base.RequestBean, c *gin.Context) (res interface{}, err error) {

	request := &UpdateOpenTabExtendRequest{}
	if !base.RequestJSON(request, c) {
		return
	}
	response := &UpdateOpenTabExtendResponse{}

	_, err = this_.ToolboxService.UpdateOpenTabExtend(request.ToolboxOpenTabModel)
	if err != nil {
		return
	}

	res = response
	return
}
