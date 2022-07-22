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

func (this_ *NodeApi) list(requestBean *base.RequestBean, c *gin.Context) (res interface{}, err error) {

	request := &ListRequest{}
	if !base.RequestJSON(request, c) {
		return
	}
	response := &ListResponse{}

	response.NodeList, err = this_.NodeService.Query(&NodeModel{
		UserId: requestBean.JWT.UserId,
	})
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

	var nodeInfo *NodeInfo
	if !node.IsROOT() {
		if request.ParentServerId == "" {
			err = errors.New("父节点ID丢失")
			return
		}
		nodeInfo = this_.NodeService.nodeContext.getNodeInfo(request.ParentServerId)
		if nodeInfo == nil {
			err = errors.New("父节点[" + request.ParentServerId + "]不存在")
			return
		}
		if nodeInfo.Info == nil {
			err = errors.New("父节点[" + request.ParentServerId + "]节点服务不存在")
			return
		}
		if nodeInfo.NodeModel == nil {
			err = errors.New("父节点[" + request.ParentServerId + "]节点数据不存在")
			return
		}
	}

	_, err = this_.NodeService.Insert(node)
	if err != nil {
		return
	}
	this_.NodeService.nodeContext.onAddNodeModel(node)
	if nodeInfo != nil && nodeInfo.Info != nil && nodeInfo.NodeModel != nil {
		var connNodeIdList []string
		if nodeInfo.NodeModel.ConnServerIds != "" {
			_ = json.Unmarshal([]byte(nodeInfo.NodeModel.ConnServerIds), &connNodeIdList)
		}
		if util.ContainsString(connNodeIdList, node.ServerId) < 0 {
			connNodeIdList = append(connNodeIdList, node.ServerId)
		}
		bs, _ := json.Marshal(connNodeIdList)
		if bs != nil {
			_, err = this_.NodeService.UpdateConnServerIds(nodeInfo.NodeModel.NodeId, string(bs))
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
