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
	idService *module_id.IDService
}

// FormatOption 格式化配置
func (this_ *NodeService) FormatOption(nodeData *NodeModel) (err error) {
	if nodeData.Option == "" {
		return
	}
	return
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
func (this_ *NodeService) CheckUserNodeNameExist(name string, userId int64) (res bool, err error) {

	sql := `SELECT COUNT(1) FROM ` + TableNode + ` WHERE deleted=2 AND (userId = ? AND name = ?)`

	count, err := this_.DatabaseWorker.Count(sql, []interface{}{userId, name})
	if err != nil {
		this_.Logger.Error("CheckUserNodeNameExist Error", zap.Error(err))
		return
	}

	res = count > 0

	return
}

// CheckUserNodeServerIdExist 查询
func (this_ *NodeService) CheckUserNodeServerIdExist(nodeServerId string, userId int64) (res bool, err error) {

	sql := `SELECT COUNT(1) FROM ` + TableNode + ` WHERE deleted=2 AND (userId = ? AND nodeServerId = ?)`

	count, err := this_.DatabaseWorker.Count(sql, []interface{}{userId, nodeServerId})
	if err != nil {
		this_.Logger.Error("CheckUserNodeServerIdExist Error", zap.Error(err))
		return
	}

	res = count > 0

	return
}

// Insert 新增
func (this_ *NodeService) Insert(node *NodeModel) (rowsAffected int64, err error) {

	checked, err := this_.CheckUserNodeNameExist(node.Name, node.UserId)
	if checked {
		err = errors.New(fmt.Sprint("节点名称[", node.Name, "]已存在"))
		return
	}
	checked, err = this_.CheckUserNodeServerIdExist(node.NodeServerId, node.UserId)
	if checked {
		err = errors.New(fmt.Sprint("节点服务[", node.NodeServerId, "]已存在"))
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

	var columns = "nodeId, nodeServerId, name, comment, option, userId, createTime"
	var values = "?, ?, ?, ?, ?, ?, ?"

	err = this_.FormatOption(node)
	if err != nil {
		return
	}

	sql := `INSERT INTO ` + TableNode + `(` + columns + `) VALUES (` + values + `) `

	rowsAffected, err = this_.DatabaseWorker.Exec(sql, []interface{}{node.NodeId, node.NodeServerId, node.Name, node.Comment, node.Option, node.UserId, node.CreateTime})
	if err != nil {
		this_.Logger.Error("Insert Error", zap.Error(err))
		return
	}

	return
}

// Update 更新
func (this_ *NodeService) Update(node *NodeModel) (rowsAffected int64, err error) {
	if node.Name != "" || node.NodeServerId != "" {
		var old *NodeModel
		old, err = this_.Get(node.NodeId)
		if err != nil {
			return
		}
		if old == nil {
			err = errors.New("节点不存在")
			return
		}
		if node.Name != "" {
			sql := `SELECT COUNT(1) FROM ` + TableNode + ` WHERE deleted=2 AND (nodeId != ? AND userId = ? AND name = ?)`

			var count int64
			count, err = this_.DatabaseWorker.Count(sql, []interface{}{node.NodeId, old.UserId, node.Name})
			if err != nil {
				return
			}
			if count > 0 {
				err = errors.New(fmt.Sprint("节点名称[", node.Name, "]已存在"))
				return
			}
		}

		if node.NodeServerId != "" {
			sql := `SELECT COUNT(1) FROM ` + TableNode + ` WHERE deleted=2 AND (nodeId != ? AND userId = ? AND nodeServerId = ?)`

			var count int64
			count, err = this_.DatabaseWorker.Count(sql, []interface{}{node.NodeId, old.UserId, node.NodeServerId})
			if err != nil {
				return
			}
			if count > 0 {
				err = errors.New(fmt.Sprint("节点服务[", node.NodeServerId, "]已存在"))
				return
			}
		}
	}

	err = this_.FormatOption(node)
	if err != nil {
		return
	}

	var values []interface{}

	sql := `UPDATE ` + TableNode + ` SET `

	sql += "updateTime=?,"
	values = append(values, time.Now())

	if node.NodeServerId != "" {
		sql += "nodeServerId=?,"
		values = append(values, node.NodeServerId)
	}
	if node.Name != "" {
		sql += "name=?,"
		values = append(values, node.Name)
	}
	if node.Comment != "" {
		sql += "comment=?,"
		values = append(values, node.Comment)
	}

	if node.Option != "" {
		sql += "option=?,"
		values = append(values, node.Option)
	}

	sql = strings.TrimSuffix(sql, ",")

	sql += " WHERE nodeId=? "
	values = append(values, node.NodeId)

	rowsAffected, err = this_.DatabaseWorker.Exec(sql, values)
	if err != nil {
		this_.Logger.Error("Update Error", zap.Error(err))
		return
	}

	return
}

// Delete 更新
func (this_ *NodeService) Delete(nodeId int64) (rowsAffected int64, err error) {

	sql := `UPDATE ` + TableNode + ` SET deleted=?,deleteTime=? WHERE nodeId=? `
	rowsAffected, err = this_.DatabaseWorker.Exec(sql, []interface{}{1, time.Now(), nodeId})
	if err != nil {
		this_.Logger.Error("Delete Error", zap.Error(err))
		return
	}

	return
}
