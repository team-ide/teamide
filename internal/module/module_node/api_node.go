package module_node

import (
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"teamide/internal/base"
	"teamide/pkg/util"
)

type ListRequest struct {
}

type ListResponse struct {
	NodeList []*NodeModel `json:"nodeList,omitempty"`
}

func (this_ *NodeApi) list(_ *base.RequestBean, c *gin.Context) (res interface{}, err error) {

	request := &ListRequest{}
	if !base.RequestJSON(request, c) {
		return
	}
	response := &ListResponse{}

	response.NodeList = this_.NodeService.nodeContext.nodeModelList
	if err != nil {
		return
	}

	res = response
	return
}

type InsertRequest struct {
	*NodeModel
	ParentServerId string `json:"parentServerId,omitempty"`
}

type InsertResponse struct {
}

func (this_ *NodeApi) insert(requestBean *base.RequestBean, c *gin.Context) (res interface{}, err error) {

	request := &InsertRequest{}
	if !base.RequestJSON(request, c) {
		return
	}
	response := &InsertResponse{}

	node := request.NodeModel
	node.UserId = requestBean.JWT.UserId

	var parentNodeModel *NodeModel
	if !node.IsROOT() {
		if request.ParentServerId == "" {
			err = errors.New("父节点ID丢失")
			return
		}
		parentNodeModel = this_.NodeService.nodeContext.getNodeModelByServerId(request.ParentServerId)
		if parentNodeModel == nil {
			err = errors.New("父节点[" + request.ParentServerId + "]不存在")
			return
		}
	}

	_, err = this_.NodeService.Insert(node)
	if err != nil {
		return
	}
	node, err = this_.NodeService.Get(node.NodeId)
	if err != nil {
		return
	}
	if node == nil {
		err = errors.New("节点数据插入失败")
		return
	}
	this_.NodeService.nodeContext.onAddNodeModel(node)
	if parentNodeModel != nil {
		var connNodeIdList []string
		if parentNodeModel.ConnServerIds != "" {
			_ = json.Unmarshal([]byte(parentNodeModel.ConnServerIds), &connNodeIdList)
		}
		if util.ContainsString(connNodeIdList, node.ServerId) < 0 {
			connNodeIdList = append(connNodeIdList, node.ServerId)
		}
		bs, _ := json.Marshal(connNodeIdList)
		if bs != nil {
			_, err = this_.NodeService.UpdateConnServerIds(parentNodeModel.NodeId, string(bs))
			if err != nil {
				return
			}
		}
	}

	res = response
	return
}

type UpdateRequest struct {
	*NodeModel
}

type UpdateResponse struct {
}

func (this_ *NodeApi) update(_ *base.RequestBean, c *gin.Context) (res interface{}, err error) {

	request := &UpdateRequest{}
	if !base.RequestJSON(request, c) {
		return
	}
	response := &UpdateResponse{}

	node := request.NodeModel

	_, err = this_.NodeService.Update(node)
	if err != nil {
		return
	}

	res = response
	return
}

func (this_ *NodeApi) updateOption(_ *base.RequestBean, c *gin.Context) (res interface{}, err error) {

	request := &UpdateRequest{}
	if !base.RequestJSON(request, c) {
		return
	}
	response := &UpdateResponse{}

	node := request.NodeModel

	_, err = this_.NodeService.UpdateOption(node)
	if err != nil {
		return
	}

	res = response
	return
}

func (this_ *NodeApi) enable(requestBean *base.RequestBean, c *gin.Context) (res interface{}, err error) {

	request := &DeleteRequest{}
	if !base.RequestJSON(request, c) {
		return
	}
	response := &DeleteResponse{}

	_, err = this_.NodeService.Enable(request.NodeModel.NodeId, requestBean.JWT.UserId)
	if err != nil {
		return
	}

	res = response
	return
}

func (this_ *NodeApi) disable(requestBean *base.RequestBean, c *gin.Context) (res interface{}, err error) {

	request := &DeleteRequest{}
	if !base.RequestJSON(request, c) {
		return
	}
	response := &DeleteResponse{}

	_, err = this_.NodeService.Disable(request.NodeModel.NodeId, requestBean.JWT.UserId)
	if err != nil {
		return
	}

	res = response
	return
}

type DeleteRequest struct {
	*NodeModel
}

type DeleteResponse struct {
}

func (this_ *NodeApi) delete(requestBean *base.RequestBean, c *gin.Context) (res interface{}, err error) {

	request := &DeleteRequest{}
	if !base.RequestJSON(request, c) {
		return
	}
	response := &DeleteResponse{}

	_, err = this_.NodeService.Delete(request.NodeModel.NodeId, requestBean.JWT.UserId)
	if err != nil {
		return
	}

	res = response
	return
}
