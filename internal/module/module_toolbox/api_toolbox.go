package module_toolbox

import (
	"github.com/gin-gonic/gin"
	"teamide/internal/base"
)

type ListRequest struct {
}

type ListResponse struct {
	ToolboxList []*ToolboxModel `json:"toolboxList,omitempty"`
}

func (this_ *ToolboxApi) list(requestBean *base.RequestBean, c *gin.Context) (res interface{}, err error) {

	request := &ListRequest{}
	if !base.RequestJSON(request, c) {
		return
	}
	response := &ListResponse{}

	response.ToolboxList, err = this_.ToolboxService.Query(&ToolboxModel{
		UserId: requestBean.JWT.UserId,
	})
	if err != nil {
		return
	}
	response.ToolboxList = append(response.ToolboxList, Others...)

	res = response
	return
}

type InsertRequest struct {
	*ToolboxModel
}

type InsertResponse struct {
}

func (this_ *ToolboxApi) insert(requestBean *base.RequestBean, c *gin.Context) (res interface{}, err error) {

	request := &InsertRequest{}
	if !base.RequestJSON(request, c) {
		return
	}
	response := &InsertResponse{}

	toolbox := request.ToolboxModel
	toolbox.UserId = requestBean.JWT.UserId

	_, err = this_.ToolboxService.Insert(toolbox)
	if err != nil {
		return
	}

	res = response
	return
}

type UpdateRequest struct {
	*ToolboxModel
}

type UpdateResponse struct {
}

func (this_ *ToolboxApi) update(requestBean *base.RequestBean, c *gin.Context) (res interface{}, err error) {

	request := &UpdateRequest{}
	if !base.RequestJSON(request, c) {
		return
	}
	response := &UpdateResponse{}

	toolbox := request.ToolboxModel

	_, err = this_.ToolboxService.Update(toolbox)
	if err != nil {
		return
	}

	res = response
	return
}

type MoveGroupRequest struct {
	*ToolboxModel
}

type MoveGroupResponse struct {
}

func (this_ *ToolboxApi) moveGroup(requestBean *base.RequestBean, c *gin.Context) (res interface{}, err error) {

	request := &MoveGroupRequest{}
	if !base.RequestJSON(request, c) {
		return
	}
	response := &MoveGroupResponse{}

	toolbox := request.ToolboxModel

	_, err = this_.ToolboxService.MoveGroup(toolbox)
	if err != nil {
		return
	}

	res = response
	return
}

type RenameRequest struct {
	*ToolboxModel
}

type RenameResponse struct {
}

func (this_ *ToolboxApi) rename(requestBean *base.RequestBean, c *gin.Context) (res interface{}, err error) {

	request := &RenameRequest{}
	if !base.RequestJSON(request, c) {
		return
	}
	response := &RenameResponse{}

	toolbox := request.ToolboxModel

	_, err = this_.ToolboxService.Rename(toolbox.ToolboxId, toolbox.Name)
	if err != nil {
		return
	}

	res = response
	return
}

type DeleteRequest struct {
	*ToolboxModel
}

type DeleteResponse struct {
}

func (this_ *ToolboxApi) delete(requestBean *base.RequestBean, c *gin.Context) (res interface{}, err error) {

	request := &DeleteRequest{}
	if !base.RequestJSON(request, c) {
		return
	}
	response := &DeleteResponse{}

	_, err = this_.ToolboxService.Delete(request.ToolboxModel.ToolboxId)
	if err != nil {
		return
	}

	res = response
	return
}

type OpenRequest struct {
	*ToolboxOpenModel
}

type OpenResponse struct {
	Open *ToolboxOpenModel `json:"open,omitempty"`
}

func (this_ *ToolboxApi) open(requestBean *base.RequestBean, c *gin.Context) (res interface{}, err error) {

	request := &OpenRequest{}
	if !base.RequestJSON(request, c) {
		return
	}
	response := &OpenResponse{}

	toolboxOpen := request.ToolboxOpenModel
	toolboxOpen.UserId = requestBean.JWT.UserId

	_, err = this_.ToolboxService.Open(toolboxOpen)
	if err != nil {
		return
	}

	response.Open = toolboxOpen

	res = response
	return
}

type QueryOpensRequest struct {
}

type QueryOpensResponse struct {
	Opens []*ToolboxOpenModel `json:"opens,omitempty"`
}

func (this_ *ToolboxApi) queryOpens(requestBean *base.RequestBean, c *gin.Context) (res interface{}, err error) {

	request := &QueryOpensRequest{}
	if !base.RequestJSON(request, c) {
		return
	}
	response := &QueryOpensResponse{}

	response.Opens, err = this_.ToolboxService.QueryOpens(requestBean.JWT.UserId)
	if err != nil {
		return
	}

	res = response
	return
}

type GetOpenRequest struct {
	*ToolboxOpenModel
}

type GetOpenResponse struct {
	Open *ToolboxOpenModel `json:"open,omitempty"`
}

func (this_ *ToolboxApi) getOpen(requestBean *base.RequestBean, c *gin.Context) (res interface{}, err error) {

	request := &GetOpenRequest{}
	if !base.RequestJSON(request, c) {
		return
	}
	response := &GetOpenResponse{}

	response.Open, err = this_.ToolboxService.GetOpen(request.OpenId)
	if err != nil {
		return
	}

	res = response
	return
}

type CloseRequest struct {
	*ToolboxOpenModel
}

type CloseResponse struct {
}

func (this_ *ToolboxApi) close(requestBean *base.RequestBean, c *gin.Context) (res interface{}, err error) {

	request := &CloseRequest{}
	if !base.RequestJSON(request, c) {
		return
	}
	response := &CloseResponse{}

	_, err = this_.ToolboxService.Close(request.OpenId)
	if err != nil {
		return
	}

	res = response
	return
}

type UpdateOpenExtendRequest struct {
	*ToolboxOpenModel
}

type UpdateOpenExtendResponse struct {
}

func (this_ *ToolboxApi) updateOpenExtend(requestBean *base.RequestBean, c *gin.Context) (res interface{}, err error) {

	request := &UpdateOpenExtendRequest{}
	if !base.RequestJSON(request, c) {
		return
	}
	response := &UpdateOpenExtendResponse{}

	_, err = this_.ToolboxService.UpdateOpenExtend(request.ToolboxOpenModel)
	if err != nil {
		return
	}

	res = response
	return
}
