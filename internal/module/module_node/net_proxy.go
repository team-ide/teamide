package module_node

import (
	"errors"
	"fmt"
	"go.uber.org/zap"
	"strings"
	"teamide/internal/module/module_id"
	"time"
)

// GetNetProxy 查询单个
func (this_ *NodeService) GetNetProxy(netProxyId int64) (res *NetProxyModel, err error) {
	res = &NetProxyModel{}

	sql := `SELECT * FROM ` + TableNodeNetProxy + ` WHERE netProxyId=? `
	find, err := this_.DatabaseWorker.QueryOne(sql, []interface{}{netProxyId}, res)
	if err != nil {
		this_.Logger.Error("Get Error", zap.Error(err))
		return
	}

	if !find {
		res = nil
	}
	return
}

// QueryNetProxy 查询
func (this_ *NodeService) QueryNetProxy(netProxy *NetProxyModel) (res []*NetProxyModel, err error) {

	var values []interface{}
	sql := `SELECT * FROM ` + TableNodeNetProxy + ` WHERE deleted=2 `
	if netProxy.UserId != 0 {
		sql += " AND userId = ?"
		values = append(values, netProxy.UserId)
	}
	if netProxy.Name != "" {
		sql += " AND name like ?"
		values = append(values, fmt.Sprint("%", netProxy.Name, "%"))
	}
	sql += " ORDER BY NAME ASC "

	err = this_.DatabaseWorker.Query(sql, values, &res)
	if err != nil {
		this_.Logger.Error("QueryNetProxy Error", zap.Error(err))
		return
	}

	return
}

// CheckNetProxyNameExist 查询
func (this_ *NodeService) CheckNetProxyNameExist(name string) (res bool, err error) {

	sql := `SELECT COUNT(1) FROM ` + TableNodeNetProxy + ` WHERE deleted=2 AND (name = ?)`

	count, err := this_.DatabaseWorker.Count(sql, []interface{}{name})
	if err != nil {
		this_.Logger.Error("CheckNetProxyNameExist Error", zap.Error(err))
		return
	}

	res = count > 0

	return
}

// CheckServerBindAddressExist 查询
func (this_ *NodeService) CheckServerBindAddressExist(serverId string, bindAddress string) (res bool, err error) {
	lastIndex := strings.LastIndex(bindAddress, ":")
	if lastIndex < 0 || lastIndex == len(bindAddress)-1 {
		err = errors.New(fmt.Sprint("网络代理[", bindAddress, "]配置有误"))
		return
	}
	endStr := bindAddress[lastIndex:]

	sql := `SELECT COUNT(1) FROM ` + TableNodeNetProxy + ` WHERE deleted=2 AND (innerServerId = ? AND innerAddress LIKE ?)`

	count, err := this_.DatabaseWorker.Count(sql, []interface{}{serverId, "%" + endStr})
	if err != nil {
		this_.Logger.Error("CheckServerBindAddressExist Error", zap.Error(err))
		return
	}

	res = count > 0

	return
}

// InsertNetProxy 新增
func (this_ *NodeService) InsertNetProxy(netProxy *NetProxyModel) (rowsAffected int64, err error) {
	checked, err := this_.CheckNetProxyNameExist(netProxy.Name)
	if err != nil {
		return
	}
	if checked {
		err = errors.New(fmt.Sprint("网络代理[", netProxy.Name, "]已存在"))
		return
	}
	checked, err = this_.CheckServerBindAddressExist(netProxy.InnerServerId, netProxy.InnerAddress)
	if err != nil {
		return
	}
	if checked {
		err = errors.New(fmt.Sprint("网络代理[", netProxy.InnerAddress, "]冲突"))
		return
	}
	if netProxy.NetProxyId == 0 {
		netProxy.NetProxyId, err = this_.idService.GetNextID(module_id.IDTypeNodeNetProxy)
		if err != nil {
			return
		}
	}
	if netProxy.CreateTime.IsZero() {
		netProxy.CreateTime = time.Now()
	}

	var columns = "netProxyId, name, comment, code, innerServerId, innerType, innerAddress, outerServerId, outerType, outerAddress, lineServerIds, option, userId, createTime"
	var values = "?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?"

	sql := `INSERT INTO ` + TableNodeNetProxy + `(` + columns + `) VALUES (` + values + `) `

	rowsAffected, err = this_.DatabaseWorker.Exec(sql, []interface{}{
		netProxy.NetProxyId,
		netProxy.Name,
		netProxy.Comment,
		netProxy.Code,
		netProxy.InnerServerId,
		netProxy.InnerType,
		netProxy.InnerAddress,
		netProxy.OuterServerId,
		netProxy.OuterType,
		netProxy.OuterAddress,
		netProxy.LineServerIds,
		netProxy.Option,
		netProxy.UserId,
		netProxy.CreateTime,
	})
	if err != nil {
		this_.Logger.Error("InsertNetProxy Error", zap.Error(err))
		return
	}

	return
}

// UpdateNetProxy 更新
func (this_ *NodeService) UpdateNetProxy(netProxy *NetProxyModel) (rowsAffected int64, err error) {

	var values []interface{}

	sql := `UPDATE ` + TableNodeNetProxy + ` SET `

	sql += "updateTime=?,"
	values = append(values, time.Now())

	if netProxy.Name != "" {
		sql += "name=?,"
		values = append(values, netProxy.Name)
	}
	if netProxy.Comment != "" {
		sql += "comment=?,"
		values = append(values, netProxy.Comment)
	}

	sql = strings.TrimSuffix(sql, ",")

	sql += " WHERE netProxyId=? "
	values = append(values, netProxy.NetProxyId)

	rowsAffected, err = this_.DatabaseWorker.Exec(sql, values)
	if err != nil {
		this_.Logger.Error("UpdateNetProxy Error", zap.Error(err))
		return
	}

	netProxy, err = this_.GetNetProxy(netProxy.NetProxyId)
	if err != nil {
		return
	}
	if netProxy != nil {
		this_.nodeContext.onUpdateNetProxyModel(netProxy)
	}
	return
}

// UpdateNetProxyOption 更新
func (this_ *NodeService) UpdateNetProxyOption(netProxy *NetProxyModel) (rowsAffected int64, err error) {

	var values []interface{}

	sql := `UPDATE ` + TableNodeNetProxy + ` SET `

	sql += "updateTime=?,"
	values = append(values, time.Now())

	sql += "option=?,"
	values = append(values, netProxy.Option)

	sql = strings.TrimSuffix(sql, ",")

	sql += " WHERE netProxyId=? "
	values = append(values, netProxy.NetProxyId)

	rowsAffected, err = this_.DatabaseWorker.Exec(sql, values)
	if err != nil {
		this_.Logger.Error("UpdateOption Error", zap.Error(err))
		return
	}

	var find = this_.nodeContext.getNetProxyModel(netProxy.NetProxyId)
	if find != nil {
		find.Option = netProxy.Option
	}
	//this_.nodeContext.onUpdateNetProxyModel(node)
	return
}

// EnableNetProxy 更新
func (this_ *NodeService) EnableNetProxy(netProxyId int64, _ int64) (rowsAffected int64, err error) {

	sql := `UPDATE ` + TableNodeNetProxy + ` SET enabled=?,updateTime=? WHERE netProxyId=? `
	rowsAffected, err = this_.DatabaseWorker.Exec(sql, []interface{}{1, time.Now(), netProxyId})
	if err != nil {
		this_.Logger.Error("EnableNetProxy Error", zap.Error(err))
		return
	}

	this_.nodeContext.onEnableNetProxyModel(netProxyId)
	return
}

// DisableNetProxy 更新
func (this_ *NodeService) DisableNetProxy(netProxyId int64, _ int64) (rowsAffected int64, err error) {

	sql := `UPDATE ` + TableNodeNetProxy + ` SET enabled=?,updateTime=? WHERE netProxyId=? `
	rowsAffected, err = this_.DatabaseWorker.Exec(sql, []interface{}{2, time.Now(), netProxyId})
	if err != nil {
		this_.Logger.Error("DisableNetProxy Error", zap.Error(err))
		return
	}

	this_.nodeContext.onDisableNetProxyModel(netProxyId)
	return
}

// DeleteNetProxy 更新
func (this_ *NodeService) DeleteNetProxy(netProxyId int64, userId int64) (rowsAffected int64, err error) {

	sql := `UPDATE ` + TableNodeNetProxy + ` SET deleted=?,deleteUserId=?,deleteTime=? WHERE netProxyId=? `
	rowsAffected, err = this_.DatabaseWorker.Exec(sql, []interface{}{1, userId, time.Now(), netProxyId})
	if err != nil {
		this_.Logger.Error("DeleteNetProxy Error", zap.Error(err))
		return
	}

	this_.nodeContext.onRemoveNetProxyModel(netProxyId)
	return
}
