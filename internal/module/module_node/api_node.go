package module_node

import (
	"github.com/gin-gonic/gin"
	"teamide/internal/base"
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

	_, err = this_.NodeService.Insert(node)
	if err != nil {
		return
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
