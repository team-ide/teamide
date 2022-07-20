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

	return
}

// CheckUserNodeNameExist 查询
func (this_ *NodeService) CheckUserNodeNameExist(name string) (res bool, err error) {

	sql := `SELECT COUNT(1) FROM ` + TableNode + ` WHERE deleted=2 AND (name = ?)`

	count, err := this_.DatabaseWorker.Count(sql, []interface{}{name})
	if err != nil {
		this_.Logger.Error("CheckUserNodeNameExist Error", zap.Error(err))
		return
	}

	res = count > 0

	return
}

// CheckUserNodeServerIdExist 查询
func (this_ *NodeService) CheckUserNodeServerIdExist(nodeServerId string) (res bool, err error) {

	sql := `SELECT COUNT(1) FROM ` + TableNode + ` WHERE deleted=2 AND (serverId = ?)`

	count, err := this_.DatabaseWorker.Count(sql, []interface{}{nodeServerId})
	if err != nil {
		this_.Logger.Error("CheckUserNodeServerIdExist Error", zap.Error(err))
		return
	}

	res = count > 0

	return
}

// Insert 新增
func (this_ *NodeService) Insert(node *NodeModel) (rowsAffected int64, err error) {

	checked, err := this_.CheckUserNodeNameExist(node.Name)
	if checked {
		err = errors.New(fmt.Sprint("节点名称[", node.Name, "]已存在"))
		return
	}
	checked, err = this_.CheckUserNodeServerIdExist(node.ServerId)
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

	var columns = "nodeId, serverId, name, comment, bindAddress, bindToken, connAddress, connToken, parentServerIds, option, isRoot, userId, createTime"
	var values = "?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?"

	sql := `INSERT INTO ` + TableNode + `(` + columns + `) VALUES (` + values + `) `

	rowsAffected, err = this_.DatabaseWorker.Exec(sql, []interface{}{node.NodeId, node.ServerId, node.Name, node.Comment, node.BindAddress, node.BindToken, node.ConnAddress, node.ConnToken, node.ParentServerIds, node.Option, node.IsRoot, node.UserId, node.CreateTime})
	if err != nil {
		this_.Logger.Error("Insert Error", zap.Error(err))
		return
	}

	err = this_.InitContext()
	if err != nil {
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

	if node.ParentServerIds != "" {
		sql += "parentServerIds=?,"
		values = append(values, node.ParentServerIds)
	}

	sql = strings.TrimSuffix(sql, ",")

	sql += " WHERE nodeId=? "
	values = append(values, node.NodeId)

	rowsAffected, err = this_.DatabaseWorker.Exec(sql, values)
	if err != nil {
		this_.Logger.Error("Update Error", zap.Error(err))
		return
	}

	err = this_.InitContext()
	if err != nil {
		return
	}
	return
}

// Delete 更新
func (this_ *NodeService) Delete(nodeId int64, userId int64) (rowsAffected int64, err error) {

	sql := `UPDATE ` + TableNode + ` SET deleted=?,deletedUserId=?,deleteTime=? WHERE nodeId=? `
	rowsAffected, err = this_.DatabaseWorker.Exec(sql, []interface{}{1, userId, time.Now(), nodeId})
	if err != nil {
		this_.Logger.Error("Delete Error", zap.Error(err))
		return
	}

	err = this_.InitContext()
	if err != nil {
		return
	}
	return
}
