package module_toolbox

import (
	"errors"
	"github.com/gin-gonic/gin"
	"teamide/pkg/base"
)

type CountRequest struct {
}

type CountResponse struct {
	Count int64 `json:"count,omitempty"`
}

func (this_ *ToolboxApi) count(requestBean *base.RequestBean, c *gin.Context) (res interface{}, err error) {

	request := &CountRequest{}
	if !base.RequestJSON(request, c) {
		return
	}
	response := &CountResponse{}

	response.Count, err = this_.ToolboxService.Count(&ToolboxModel{
		UserId: requestBean.JWT.UserId,
	})
	if err != nil {
		return
	}

	response.Count += int64(len(Others))

	res = response
	return
}

type ListRequest struct {
	ToolboxType string `json:"toolboxType,omitempty"`
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
		UserId:      requestBean.JWT.UserId,
		ToolboxType: request.ToolboxType,
	})
	if err != nil {
		return
	}

	//if request.ToolboxType == "" {
	//	response.ToolboxList = append(response.ToolboxList, Others...)
	//}

	res = response
	return
}

type GetRequest struct {
	ToolboxId int64 `json:"toolboxId,omitempty"`
}

func (this_ *ToolboxApi) get(requestBean *base.RequestBean, c *gin.Context) (res interface{}, err error) {

	request := &GetRequest{}
	if !base.RequestJSON(request, c) {
		return
	}

	res, err = this_.ToolboxService.Get(request.ToolboxId)
	if err != nil {
		return
	}

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

	if request.ToolboxId != 0 {
		var find *ToolboxModel
		find, err = this_.ToolboxService.Get(request.ToolboxId)
		if err != nil {
			return
		}
		if find != nil && find.UserId != 0 {
			if requestBean.JWT == nil || find.UserId != requestBean.JWT.UserId {
				err = errors.New("工具[" + find.Name + "]不属于当前用户，无法操作")
				return
			}
		}
	}

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

	if request.ToolboxId != 0 {
		var find *ToolboxModel
		find, err = this_.ToolboxService.Get(request.ToolboxId)
		if err != nil {
			return
		}
		if find != nil && find.UserId != 0 {
			if requestBean.JWT == nil || find.UserId != requestBean.JWT.UserId {
				err = errors.New("工具[" + find.Name + "]不属于当前用户，无法操作")
				return
			}
		}
	}

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

	if request.ToolboxId != 0 {
		var find *ToolboxModel
		find, err = this_.ToolboxService.Get(request.ToolboxId)
		if err != nil {
			return
		}
		if find != nil && find.UserId != 0 {
			if requestBean.JWT == nil || find.UserId != requestBean.JWT.UserId {
				err = errors.New("工具[" + find.Name + "]不属于当前用户，无法操作")
				return
			}
		}
	}

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

	if request.ToolboxId != 0 {
		var find *ToolboxModel
		find, err = this_.ToolboxService.Get(request.ToolboxId)
		if err != nil {
			return
		}
		if find != nil && find.UserId != 0 {
			if requestBean.JWT == nil || find.UserId != requestBean.JWT.UserId {
				err = errors.New("工具[" + find.Name + "]不属于当前用户，无法操作")
				return
			}
		}
	}

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

	if request.ToolboxId != 0 {
		var find *ToolboxModel
		find, err = this_.ToolboxService.Get(request.ToolboxId)
		if err != nil {
			return
		}
		if find != nil && find.UserId != 0 {
			if requestBean.JWT == nil || find.UserId != requestBean.JWT.UserId {
				err = errors.New("工具[" + find.Name + "]不属于当前用户，无法操作")
				return
			}
		}
	}

	toolboxOpen := request.ToolboxOpenModel
	toolboxOpen.UserId = requestBean.JWT.UserId

	_, err = this_.ToolboxService.Open(toolboxOpen)
	if err != nil {
		return
	}

	response.Open, err = this_.ToolboxService.GetOpen(toolboxOpen.OpenId)
	if err != nil {
		return
	}

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

func (this_ *ToolboxApi) getOpen(_ *base.RequestBean, c *gin.Context) (res interface{}, err error) {

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

	_, err = this_.ToolboxService.UpdateOpenExtend(request.ToolboxOpenModel)
	if err != nil {
		return
	}

	res = response
	return
}

type UpdateOpenSequenceRequest struct {
	*ToolboxOpenModel
}

type UpdateOpenSequenceResponse struct {
}

func (this_ *ToolboxApi) UpdateOpenSequence(requestBean *base.RequestBean, c *gin.Context) (res interface{}, err error) {

	request := &UpdateOpenSequenceRequest{}
	if !base.RequestJSON(request, c) {
		return
	}
	response := &UpdateOpenSequenceResponse{}

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

	_, err = this_.ToolboxService.UpdateOpenSequence(request.ToolboxOpenModel)
	if err != nil {
		return
	}

	res = response
	return
}
