package module_node

import (
	"errors"
	"fmt"
	"go.uber.org/zap"
	"strings"
	"teamide/internal/context"
	"teamide/internal/module/module_id"
	"time"
)

// NewNodeService 根据库配置创建NodeService
func NewNodeService(ServerContext *context.ServerContext) (res *NodeService) {

	idService := module_id.NewIDService(ServerContext)

	res = &NodeService{
		ServerContext: ServerContext,
		idService:     idService,
	}
	return
}

// NodeService 工具箱服务
type NodeService struct {
	*context.ServerContext
	idService   *module_id.IDService
	nodeContext *NodeContext
}

// Get 查询单个
func (this_ *NodeService) Get(nodeId int64) (res *NodeModel, err error) {
	res = &NodeModel{}

	sql := `SELECT * FROM ` + TableNode + ` WHERE nodeId=? `
	find, err := this_.DatabaseWorker.QueryOne(sql, []interface{}{nodeId}, res)
	if err != nil {
		this_.Logger.Error("Get Error", zap.Error(err))
		return
	}

	if !find {
		res = nil
	}
	res.ConnServerIdList = GetStringList(res.ConnServerIds)
	res.HistoryConnServerIdList = GetStringList(res.HistoryConnServerIds)
	return
}

// Query 查询
func (this_ *NodeService) Query(node *NodeModel) (res []*NodeModel, err error) {

	var values []interface{}
	sql := `SELECT * FROM ` + TableNode + ` WHERE deleted=2 `
	if node.UserId != 0 {
		sql += " AND userId = ?"
		values = append(values, node.UserId)
	}
	if node.Name != "" {
		sql += " AND name like ?"
		values = append(values, fmt.Sprint("%", node.Name, "%"))
	}
	sql += " ORDER BY NAME ASC "

	err = this_.DatabaseWorker.Query(sql, values, &res)
	if err != nil {
		this_.Logger.Error("Query Error", zap.Error(err))
		return
	}
	for _, one := range res {
		one.ConnServerIdList = GetStringList(one.ConnServerIds)
		one.HistoryConnServerIdList = GetStringList(one.HistoryConnServerIds)
	}

	return
}

// CheckNodeNameExist 查询
func (this_ *NodeService) CheckNodeNameExist(name string) (res bool, err error) {

	sql := `SELECT COUNT(1) FROM ` + TableNode + ` WHERE deleted=2 AND (name = ?)`

	count, err := this_.DatabaseWorker.Count(sql, []interface{}{name})
	if err != nil {
		this_.Logger.Error("CheckNodeNameExist Error", zap.Error(err))
		return
	}

	res = count > 0

	return
}

// CheckNodeServerIdExist 查询
func (this_ *NodeService) CheckNodeServerIdExist(nodeServerId string) (res bool, err error) {

	sql := `SELECT COUNT(1) FROM ` + TableNode + ` WHERE deleted=2 AND (serverId = ?)`

	count, err := this_.DatabaseWorker.Count(sql, []interface{}{nodeServerId})
	if err != nil {
		this_.Logger.Error("CheckNodeServerIdExist Error", zap.Error(err))
		return
	}

	res = count > 0

	return
}

// Insert 新增
func (this_ *NodeService) Insert(node *NodeModel) (rowsAffected int64, err error) {

	checked, err := this_.CheckNodeNameExist(node.Name)
	if err != nil {
		return
	}
	if checked {
		err = errors.New(fmt.Sprint("节点名称[", node.Name, "]已存在"))
		return
	}
	checked, err = this_.CheckNodeServerIdExist(node.ServerId)
	if err != nil {
		return
	}
	if checked {
		err = errors.New(fmt.Sprint("节点服务[", node.ServerId, "]已存在"))
		return
	}
	if node.NodeId == 0 {
		node.NodeId, err = this_.idService.GetNextID(module_id.IDTypeNode)
		if err != nil {
			return
		}
	}
	if node.CreateTime.IsZero() {
		node.CreateTime = time.Now()
	}

	var columns = "nodeId, serverId, name, comment, bindAddress, bindToken, connAddress, connToken, connServerIds, historyConnServerIds, option, isLocal, userId, createTime"
	var values = "?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?"

	sql := `INSERT INTO ` + TableNode + `(` + columns + `) VALUES (` + values + `) `

	rowsAffected, err = this_.DatabaseWorker.Exec(sql, []interface{}{
		node.NodeId,
		node.ServerId,
		node.Name,
		node.Comment,
		node.BindAddress,
		node.BindToken,
		node.ConnAddress,
		node.ConnToken,
		node.ConnServerIds,
		node.HistoryConnServerIds,
		node.Option,
		node.IsLocal,
		node.UserId,
		node.CreateTime,
	})
	if err != nil {
		this_.Logger.Error("Insert Error", zap.Error(err))
		return
	}

	return
}

// Update 更新
func (this_ *NodeService) Update(node *NodeModel) (rowsAffected int64, err error) {
	if node.Name != "" || node.ServerId != "" {
		if node.Name != "" {
			sql := `SELECT COUNT(1) FROM ` + TableNode + ` WHERE deleted=2 AND (nodeId != ? AND name = ?)`
			var count int64
			count, err = this_.DatabaseWorker.Count(sql, []interface{}{node.NodeId, node.Name})
			if err != nil {
				return
			}
			if count > 0 {
				err = errors.New(fmt.Sprint("节点名称[", node.Name, "]已存在"))
				return
			}
		}

		if node.ServerId != "" {
			sql := `SELECT COUNT(1) FROM ` + TableNode + ` WHERE deleted=2 AND (nodeId != ? AND serverId = ?)`

			var count int64
			count, err = this_.DatabaseWorker.Count(sql, []interface{}{node.NodeId, node.ServerId})
			if err != nil {
				return
			}
			if count > 0 {
				err = errors.New(fmt.Sprint("节点服务[", node.ServerId, "]已存在"))
				return
			}
		}
	}

	var values []interface{}

	sql := `UPDATE ` + TableNode + ` SET `

	sql += "updateTime=?,"
	values = append(values, time.Now())

	if node.ServerId != "" {
		sql += "serverId=?,"
		values = append(values, node.ServerId)
	}
	if node.Name != "" {
		sql += "name=?,"
		values = append(values, node.Name)
	}
	if node.Comment != "" {
		sql += "comment=?,"
		values = append(values, node.Comment)
	}

	if node.BindAddress != "" {
		sql += "bindAddress=?,"
		values = append(values, node.BindAddress)
	}

	if node.BindToken != "" {
		sql += "bindToken=?,"
		values = append(values, node.BindToken)
	}

	if node.ConnAddress != "" {
		sql += "connAddress=?,"
		values = append(values, node.ConnAddress)
	}

	if node.ConnToken != "" {
		sql += "connToken=?,"
		values = append(values, node.ConnToken)
	}

	if node.ConnServerIds != "" {
		sql += "connServerIds=?,"
		values = append(values, node.ConnServerIds)
	}

	sql = strings.TrimSuffix(sql, ",")

	sql += " WHERE nodeId=? "
	values = append(values, node.NodeId)

	rowsAffected, err = this_.DatabaseWorker.Exec(sql, values)
	if err != nil {
		this_.Logger.Error("Update Error", zap.Error(err))
		return
	}

	node, err = this_.Get(node.NodeId)
	if err != nil {
		return
	}
	if node != nil {
		this_.nodeContext.onUpdateNodeModel(node)
	}
	return
}

// UpdateOption 更新
func (this_ *NodeService) UpdateOption(node *NodeModel) (rowsAffected int64, err error) {

	var values []interface{}

	sql := `UPDATE ` + TableNode + ` SET `

	sql += "updateTime=?,"
	values = append(values, time.Now())

	sql += "option=?,"
	values = append(values, node.Option)

	sql = strings.TrimSuffix(sql, ",")

	sql += " WHERE nodeId=? "
	values = append(values, node.NodeId)

	rowsAffected, err = this_.DatabaseWorker.Exec(sql, values)
	if err != nil {
		this_.Logger.Error("UpdateOption Error", zap.Error(err))
		return
	}

	var find = this_.nodeContext.getNodeModel(node.NodeId)
	if find != nil {
		find.Option = node.Option
	}
	//this_.nodeContext.onUpdateNodeModel(node)
	return
}

// UpdateConnServerIds 更新
func (this_ *NodeService) UpdateConnServerIds(nodeId int64, connServerIds string) (rowsAffected int64, err error) {

	var values []interface{}

	sql := `UPDATE ` + TableNode + ` SET `

	sql += "updateTime=?,"
	values = append(values, time.Now())

	sql += "connServerIds=?,"
	values = append(values, connServerIds)

	sql = strings.TrimSuffix(sql, ",")

	sql += " WHERE nodeId=? "
	values = append(values, nodeId)

	rowsAffected, err = this_.DatabaseWorker.Exec(sql, values)
	if err != nil {
		this_.Logger.Error("UpdateConnServerIds Error", zap.Error(err))
		return
	}

	this_.nodeContext.onUpdateNodeConnServerIds(nodeId, connServerIds)
	return
}

// UpdateHistoryConnServerIds 更新
func (this_ *NodeService) UpdateHistoryConnServerIds(nodeId int64, historyConnServerIds string) (rowsAffected int64, err error) {

	var values []interface{}

	sql := `UPDATE ` + TableNode + ` SET `

	sql += "updateTime=?,"
	values = append(values, time.Now())

	sql += "historyConnServerIds=?,"
	values = append(values, historyConnServerIds)

	sql = strings.TrimSuffix(sql, ",")

	sql += " WHERE nodeId=? "
	values = append(values, nodeId)

	rowsAffected, err = this_.DatabaseWorker.Exec(sql, values)
	if err != nil {
		this_.Logger.Error("UpdateHistoryConnServerIds Error", zap.Error(err))
		return
	}

	this_.nodeContext.onUpdateNodeHistoryConnServerIds(nodeId, historyConnServerIds)
	return
}

// Enable 更新
func (this_ *NodeService) Enable(nodeId int64, _ int64) (rowsAffected int64, err error) {

	sql := `UPDATE ` + TableNode + ` SET enabled=?,updateTime=? WHERE nodeId=? `
	rowsAffected, err = this_.DatabaseWorker.Exec(sql, []interface{}{1, time.Now(), nodeId})
	if err != nil {
		this_.Logger.Error("Enabled Error", zap.Error(err))
		return
	}

	this_.nodeContext.onEnableNodeModel(nodeId)
	return
}

// Disable 更新
func (this_ *NodeService) Disable(nodeId int64, _ int64) (rowsAffected int64, err error) {

	sql := `UPDATE ` + TableNode + ` SET enabled=?,updateTime=? WHERE nodeId=? `
	rowsAffected, err = this_.DatabaseWorker.Exec(sql, []interface{}{2, time.Now(), nodeId})
	if err != nil {
		this_.Logger.Error("Disable Error", zap.Error(err))
		return
	}

	this_.nodeContext.onDisableNodeModel(nodeId)
	return
}

// Delete 更新
func (this_ *NodeService) Delete(nodeId int64, userId int64) (rowsAffected int64, err error) {

	sql := `UPDATE ` + TableNode + ` SET deleted=?,deleteUserId=?,deleteTime=? WHERE nodeId=? `
	rowsAffected, err = this_.DatabaseWorker.Exec(sql, []interface{}{1, userId, time.Now(), nodeId})
	if err != nil {
		this_.Logger.Error("Delete Error", zap.Error(err))
		return
	}

	this_.nodeContext.onRemoveNodeModel(nodeId)
	return
}
