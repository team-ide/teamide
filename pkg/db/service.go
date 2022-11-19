package db

import (
	"fmt"
	"github.com/team-ide/go-dialect/dialect"
	"github.com/team-ide/go-dialect/worker"
	"teamide/pkg/util"
	"time"
)

func CreateService(config *DatabaseConfig) (service *Service, err error) {
	service = &Service{
		config: config,
	}
	service.lastUseTime = util.GetNowTime()
	err = service.init()
	return
}

type SqlParam struct {
	Sql    string        `json:"sql,omitempty"`
	Params []interface{} `json:"params,omitempty"`
}

type Service struct {
	config         *DatabaseConfig
	lastUseTime    int64
	DatabaseWorker *DatabaseWorker
}

func (this_ *Service) init() (err error) {
	this_.DatabaseWorker, err = NewDatabaseWorker(this_.config)
	if err != nil {
		return
	}
	return
}

func (this_ *Service) GetDatabaseWorker() *DatabaseWorker {
	return this_.DatabaseWorker
}

func (this_ *Service) GetWaitTime() int64 {
	return 10 * 60 * 1000
}

func (this_ *Service) GetLastUseTime() int64 {
	return this_.lastUseTime
}

func (this_ *Service) SetLastUseTime() {
	this_.lastUseTime = util.GetNowTime()
}

func (this_ *Service) Stop() {
	if this_.DatabaseWorker != nil {
		_ = this_.DatabaseWorker.Close()
	}
}

func (this_ *Service) OwnersSelect(param *dialect.ParamModel) (owners []*dialect.OwnerModel, err error) {
	owners, err = worker.OwnersSelect(this_.DatabaseWorker.db, this_.DatabaseWorker.Dialect, param)
	return
}

func (this_ *Service) TablesSelect(param *dialect.ParamModel, ownerName string) (tables []*dialect.TableModel, err error) {
	tables, err = worker.TablesSelect(this_.DatabaseWorker.db, this_.DatabaseWorker.Dialect, param, ownerName)
	return
}

func (this_ *Service) TableDetail(param *dialect.ParamModel, ownerName string, tableName string) (tableDetail *dialect.TableModel, err error) {
	tableDetail, err = worker.TableDetail(this_.DatabaseWorker.db, this_.DatabaseWorker.Dialect, param, ownerName, tableName, true)
	return
}

func (this_ *Service) OwnerCreate(param *dialect.ParamModel, owner *dialect.OwnerModel) (created bool, err error) {
	created, err = worker.OwnerCreate(this_.DatabaseWorker.db, this_.DatabaseWorker.Dialect, param, owner)
	return
}

func (this_ *Service) OwnerDelete(param *dialect.ParamModel, ownerName string) (deleted bool, err error) {
	deleted, err = worker.OwnerDelete(this_.DatabaseWorker.db, this_.DatabaseWorker.Dialect, param, ownerName)
	return
}

func (this_ *Service) TableCreate(param *dialect.ParamModel, ownerName string, table *dialect.TableModel) (err error) {
	err = worker.TableCreate(this_.DatabaseWorker.db, this_.DatabaseWorker.Dialect, param, ownerName, table)
	return
}

func (this_ *Service) TableDelete(param *dialect.ParamModel, ownerName string, tableName string) (err error) {
	err = worker.TableDelete(this_.DatabaseWorker.db, this_.DatabaseWorker.Dialect, param, ownerName, tableName)
	return
}

type DataListResult struct {
	Sql      string                   `json:"sql"`
	Total    int                      `json:"total"`
	Params   []interface{}            `json:"params"`
	DataList []map[string]interface{} `json:"dataList"`
}

func (this_ *Service) DataList(param *dialect.ParamModel, ownerName string, tableName string, columnList []*dialect.ColumnModel, whereList []*dialect.Where, orderList []*dialect.Order, pageSize int, pageNo int) (dataListResult DataListResult, err error) {

	sql, values, err := this_.DatabaseWorker.Dialect.DataListSelectSql(param, ownerName, tableName, columnList, whereList, orderList)
	if err != nil {
		return
	}

	page := worker.NewPage()
	page.PageSize = pageSize
	page.PageNo = pageNo
	listMap, err := this_.DatabaseWorker.QueryMapPage(sql, values, page)
	if err != nil {
		return
	}
	for _, one := range listMap {
		for k, v := range one {
			if v == nil {
				continue
			}
			switch tV := v.(type) {
			case time.Time:
				if tV.IsZero() {
					one[k] = nil
				} else {
					one[k] = util.GetTimeTime(tV)
				}
			default:
				one[k] = fmt.Sprint(tV)
			}
		}
	}
	dataListResult.Sql = this_.DatabaseWorker.PackPageSql(sql, page.PageSize, page.PageNo)
	dataListResult.Params = values
	dataListResult.Total = page.TotalCount
	dataListResult.DataList = listMap
	return
}

func (this_ *Service) Execs(sqlList []string, paramsList [][]interface{}) (res int64, err error) {
	res, err = this_.DatabaseWorker.Execs(sqlList, paramsList)
	if err != nil {
		return
	}
	return
}
