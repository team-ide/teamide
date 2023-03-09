package module_node

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/team-ide/go-tool/util"
	"teamide/pkg/base"
)

type netProxyMonitorDataRequest struct {
	IdList []string `json:"idList,omitempty"`
}

type netProxyMonitorDataResponse struct {
	NetProxyMonitorDataList []*NetProxyMonitorData `json:"netProxyMonitorDataList,omitempty"`
}

type NetProxyMonitorData struct {
	Id               string             `json:"id,omitempty"`
	InnerMonitorData *MonitorDataFormat `json:"innerMonitorData,omitempty"`
	OuterMonitorData *MonitorDataFormat `json:"outerMonitorData,omitempty"`
}

func (this_ *NodeApi) netProxyMonitorData(_ *base.RequestBean, c *gin.Context) (res interface{}, err error) {

	request := &netProxyMonitorDataRequest{}
	if !base.RequestJSON(request, c) {
		return
	}
	response := &netProxyMonitorDataResponse{}

	for _, id := range request.IdList {
		netProxyInfo := this_.NodeService.nodeContext.getNetProxyModelByCode(id)
		if netProxyInfo == nil {
			continue
		}
		innerLineNodeIdList := this_.NodeService.nodeContext.GetNodeLineTo(netProxyInfo.InnerServerId)
		outerLineNodeIdList := this_.NodeService.nodeContext.GetNodeLineTo(netProxyInfo.OuterServerId)

		one := &NetProxyMonitorData{
			Id: id,
		}
		if len(innerLineNodeIdList) > 0 {
			one.InnerMonitorData = ToMonitorDataFormat(this_.NodeService.nodeContext.GetServer().GetNetProxyInnerMonitorData(innerLineNodeIdList, id))
		}
		if len(outerLineNodeIdList) > 0 {
			one.OuterMonitorData = ToMonitorDataFormat(this_.NodeService.nodeContext.GetServer().GetNetProxyOuterMonitorData(outerLineNodeIdList, id))
		}
		response.NetProxyMonitorDataList = append(response.NetProxyMonitorDataList, one)
	}

	res = response
	return
}

type NetProxyListRequest struct {
}

type NetProxyListResponse struct {
	NetProxyList []*NetProxyModel `json:"netProxyList,omitempty"`
}

func (this_ *NodeApi) netProxyList(requestBean *base.RequestBean, c *gin.Context) (res interface{}, err error) {

	request := &NetProxyListRequest{}
	if !base.RequestJSON(request, c) {
		return
	}
	response := &NetProxyListResponse{}

	var netProxyModelList = this_.NodeService.nodeContext.getUserNetProxyModelList(requestBean.JWT.UserId)
	response.NetProxyList = netProxyModelList

	res = response
	return
}

type NetProxyInsertRequest struct {
	*NetProxyModel
}

type NetProxyInsertResponse struct {
}

func (this_ *NodeApi) netProxyInsert(requestBean *base.RequestBean, c *gin.Context) (res interface{}, err error) {

	request := &NetProxyInsertRequest{}
	if !base.RequestJSON(request, c) {
		return
	}
	response := &NetProxyInsertResponse{}

	netProxy := request.NetProxyModel
	netProxy.UserId = requestBean.JWT.UserId
	netProxy.Code = util.GetUUID()

	err = this_.NodeService.nodeContext.formatNetProxy(netProxy)
	if err != nil {
		return
	}

	_, err = this_.NodeService.InsertNetProxy(netProxy)
	if err != nil {
		return
	}
	netProxy, err = this_.NodeService.GetNetProxy(netProxy.NetProxyId)
	if err != nil {
		return
	}
	if netProxy == nil {
		err = errors.New("代理数据插入失败")
		return
	}
	this_.NodeService.nodeContext.onAddNetProxyModel(netProxy)

	res = response
	return
}

type NetProxyUpdateRequest struct {
	*NetProxyModel
}

type NetProxyUpdateResponse struct {
}

func (this_ *NodeApi) netProxyUpdate(_ *base.RequestBean, c *gin.Context) (res interface{}, err error) {

	request := &NetProxyUpdateRequest{}
	if !base.RequestJSON(request, c) {
		return
	}
	response := &NetProxyUpdateResponse{}

	_, err = this_.NodeService.UpdateNetProxy(request.NetProxyModel)
	if err != nil {
		return
	}

	res = response
	return
}

func (this_ *NodeApi) netProxyUpdateOption(_ *base.RequestBean, c *gin.Context) (res interface{}, err error) {

	request := &NetProxyUpdateRequest{}
	if !base.RequestJSON(request, c) {
		return
	}
	response := &NetProxyUpdateResponse{}

	_, err = this_.NodeService.UpdateNetProxyOption(request.NetProxyModel)
	if err != nil {
		return
	}

	res = response
	return
}

func (this_ *NodeApi) netProxyEnable(requestBean *base.RequestBean, c *gin.Context) (res interface{}, err error) {

	request := &NetProxyDeleteRequest{}
	if !base.RequestJSON(request, c) {
		return
	}
	response := &NetProxyDeleteResponse{}

	_, err = this_.NodeService.EnableNetProxy(request.NetProxyModel.NetProxyId, requestBean.JWT.UserId)
	if err != nil {
		return
	}

	res = response
	return
}

func (this_ *NodeApi) netProxyDisable(requestBean *base.RequestBean, c *gin.Context) (res interface{}, err error) {

	request := &NetProxyDeleteRequest{}
	if !base.RequestJSON(request, c) {
		return
	}
	response := &NetProxyDeleteResponse{}

	_, err = this_.NodeService.DisableNetProxy(request.NetProxyModel.NetProxyId, requestBean.JWT.UserId)
	if err != nil {
		return
	}

	res = response
	return
}

type NetProxyDeleteRequest struct {
	*NetProxyModel
}

type NetProxyDeleteResponse struct {
}

func (this_ *NodeApi) netProxyDelete(requestBean *base.RequestBean, c *gin.Context) (res interface{}, err error) {

	request := &NetProxyDeleteRequest{}
	if !base.RequestJSON(request, c) {
		return
	}
	response := &NetProxyDeleteResponse{}

	_, err = this_.NodeService.DeleteNetProxy(request.NetProxyModel.NetProxyId, requestBean.JWT.UserId)
	if err != nil {
		return
	}

	res = response
	return
}
