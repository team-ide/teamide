package toolbox

import (
	"fmt"
	"go.uber.org/zap"
	"teamide/pkg/db"
	"teamide/pkg/util"
	"time"
)

type saveDataListTask struct {
	Key             string                 `json:"key,omitempty"`
	Database        string                 `json:"database,omitempty"`
	Table           string                 `json:"table,omitempty"`
	ColumnList      []*db.TableColumnModel `json:"columnList,omitempty"`
	generateParam   *db.GenerateParam
	UpdateList      []map[string]interface{} `json:"-"`
	UpdateWhereList []map[string]interface{} `json:"-"`
	InsertList      []map[string]interface{} `json:"-"`
	DeleteList      []map[string]interface{} `json:"-"`
	SaveBatchNumber int                      `json:"saveBatchNumber,omitempty"`
	DataCount       int                      `json:"dataCount"`
	ReadyDataCount  int                      `json:"readyDataCount"`
	SaveSuccess     int                      `json:"saveSuccess"`
	SaveError       int                      `json:"saveError"`
	IsEnd           bool                     `json:"isEnd,omitempty"`
	StartTime       time.Time                `json:"startTime,omitempty"`
	EndTime         time.Time                `json:"endTime,omitempty"`
	Error           string                   `json:"error,omitempty"`
	UseTime         int64                    `json:"useTime"`
	IsStop          bool                     `json:"isStop"`
	service         DatabaseService
}

func (this_ *saveDataListTask) Stop() {
	this_.IsStop = true
}

func (this_ *saveDataListTask) Start() (err error) {
	this_.StartTime = time.Now()
	defer func() {
		if err := recover(); err != nil {
			Logger.Error("根据保存数据异常", zap.Any("error", err))
			this_.Error = fmt.Sprint(err)
		}
		this_.EndTime = time.Now()
		this_.IsEnd = true
		this_.UseTime = util.GetTimeTime(this_.EndTime) - util.GetTimeTime(this_.StartTime)
	}()
	this_.DataCount += len(this_.InsertList)
	this_.DataCount += len(this_.UpdateList)
	this_.DataCount += len(this_.DeleteList)

	saveBatchNumber := this_.SaveBatchNumber
	if saveBatchNumber <= 0 {
		saveBatchNumber = 10
	}

	var sqlList []string
	var valuesList [][]interface{}

	sqlList, valuesList, err = this_.SaveDataListSql()
	if err != nil {
		return
	}
	res, err := this_.service.Execs(sqlList, valuesList)
	if err != nil {
		this_.SaveError += len(sqlList)
		return
	}
	this_.SaveSuccess += int(res)
	return
}

func (this_ *saveDataListTask) SaveDataListSql() (sqlList []string, valuesList [][]interface{}, err error) {

	var sqlList_ []string
	var valuesList_ [][]interface{}
	if len(this_.InsertList) > 0 {
		sqlList_, valuesList_, err = db.DataListInsertSql(this_.generateParam, this_.Database, this_.Table, this_.ColumnList, this_.InsertList)
		if err != nil {
			return
		}
		sqlList = append(sqlList, sqlList_...)
		valuesList = append(valuesList, valuesList_...)
	}
	if len(this_.UpdateList) > 0 {
		sqlList_, valuesList_, err = db.DataListUpdateSql(this_.generateParam, this_.Database, this_.Table, this_.ColumnList, this_.UpdateList, this_.UpdateWhereList)
		if err != nil {
			return
		}
		sqlList = append(sqlList, sqlList_...)
		valuesList = append(valuesList, valuesList_...)
	}
	if len(this_.DeleteList) > 0 {
		sqlList_, valuesList_, err = db.DataListDeleteSql(this_.generateParam, this_.Database, this_.Table, this_.ColumnList, this_.DeleteList)
		if err != nil {
			return
		}
		sqlList = append(sqlList, sqlList_...)
		valuesList = append(valuesList, valuesList_...)
	}
	return
}
