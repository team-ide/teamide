package module_toolbox

import (
	"errors"
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

	if request.OpenId != 0 {
		var find *ToolboxOpenModel
		find, err = this_.ToolboxService.GetOpen(request.OpenId)
		if err != nil {
			return
		}
		if find != nil && find.UserId != 0 {
			if requestBean.JWT == nil || find.UserId != requestBean.JWT.UserId {
				err = errors.New("不属于当前用户，无法操作")
				return
			}
		}
	}

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

func (this_ *ToolboxApi) queryOpenTabs(_ *base.RequestBean, c *gin.Context) (res interface{}, err error) {

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

	if request.TabId != 0 {
		var find *ToolboxOpenTabModel
		find, err = this_.ToolboxService.GetOpenTab(request.TabId)
		if err != nil {
			return
		}
		if find != nil && find.UserId != 0 {
			if requestBean.JWT == nil || find.UserId != requestBean.JWT.UserId {
				err = errors.New("不属于当前用户，无法操作")
				return
			}
		}
	}
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

	if request.TabId != 0 {
		var find *ToolboxOpenTabModel
		find, err = this_.ToolboxService.GetOpenTab(request.TabId)
		if err != nil {
			return
		}
		if find != nil && find.UserId != 0 {
			if requestBean.JWT == nil || find.UserId != requestBean.JWT.UserId {
				err = errors.New("不属于当前用户，无法操作")
				return
			}
		}
	}

	_, err = this_.ToolboxService.UpdateOpenTabExtend(request.ToolboxOpenTabModel)
	if err != nil {
		return
	}

	res = response
	return
}
