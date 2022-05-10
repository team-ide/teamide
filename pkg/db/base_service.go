package db

import (
	"context"
	"gitee.com/chunanyong/zorm"
	"go.uber.org/zap"
)

var (
	baseContext = context.Background()
)

// DatabaseBase 基础操作
type DatabaseBase struct {
	config DatabaseConfig
	dbDao  *zorm.DBDao
}

func (this_ *DatabaseBase) GetContext() (ctx context.Context) {
	ctx, _ = this_.dbDao.BindContextDBConnection(baseContext)
	return ctx
}

func (this_ *DatabaseBase) GetConfig() (config DatabaseConfig) {
	config = this_.config
	return
}

func (this_ *DatabaseBase) Open() (err error) {
	finder := zorm.NewFinder()
	finder.Append("SELECT 1")

	var count int
	_, err = zorm.QueryRow(this_.GetContext(), finder, &count)
	return
}

func (this_ *DatabaseBase) Close() (err error) {
	err = this_.dbDao.CloseDB()
	return
}

func (this_ *DatabaseBase) Exec(sql string, params []interface{}) (rowsAffected int64, err error) {

	rowsAffected, err = this_.Execs([]string{sql}, [][]interface{}{params})

	if err != nil {
		return
	}
	return
}

func (this_ *DatabaseBase) Execs(sqlList []string, paramsList [][]interface{}) (rowsAffected int64, err error) {

	_, err = zorm.Transaction(this_.GetContext(), func(ctx context.Context) (interface{}, error) {

		var err error
		for index, _ := range sqlList {
			sql := sqlList[index]
			params := paramsList[index]
			finder := zorm.NewFinder()
			finder.Append(sql, params...)

			var updated int
			updated, err = zorm.UpdateFinder(ctx, finder)
			if err != nil {
				Logger.Error("Exec Error", zap.Any("sql", sql), zap.Any("params", params), zap.Error(err))
				return nil, err
			}
			rowsAffected += int64(updated)
		}

		//如果返回的err不是nil,事务就会回滚
		return nil, err
	})

	if err != nil {
		return
	}
	return
}

func (this_ *DatabaseBase) Count(sql string, params []interface{}) (count int64, err error) {

	finder := zorm.NewFinder()
	finder.Append(sql, params...)

	count, err = this_.FinderCount(finder)

	if err != nil {
		Logger.Error("Count Error", zap.Any("sql", sql), zap.Any("params", params), zap.Error(err))
		return
	}
	return
}

func (this_ *DatabaseBase) FinderCount(finder *zorm.Finder) (count int64, err error) {

	_, err = zorm.QueryRow(this_.GetContext(), finder, &count)

	if err != nil {
		return
	}
	return
}

func (this_ *DatabaseBase) QueryOne(sql string, params []interface{}, one interface{}) (find bool, err error) {

	finder := zorm.NewFinder()
	finder.Append(sql, params...)

	find, err = this_.FinderQueryOne(finder, one)

	if err != nil {
		Logger.Error("Query Error", zap.Any("sql", sql), zap.Any("params", params), zap.Error(err))
		return
	}

	return
}

func (this_ *DatabaseBase) FinderQueryOne(finder *zorm.Finder, one interface{}) (find bool, err error) {

	find, err = zorm.QueryRow(this_.GetContext(), finder, one)

	if err != nil {
		return
	}

	return
}

func (this_ *DatabaseBase) Query(sql string, params []interface{}, list interface{}) (err error) {

	finder := zorm.NewFinder()
	finder.Append(sql, params...)

	err = this_.FinderQuery(finder, list)

	if err != nil {
		Logger.Error("Query Error", zap.Any("sql", sql), zap.Any("params", params), zap.Error(err))
		return
	}

	return
}

func (this_ *DatabaseBase) FinderQuery(finder *zorm.Finder, list interface{}) (err error) {

	err = zorm.Query(this_.GetContext(), finder, list, nil)

	if err != nil {
		return
	}

	return
}

func (this_ *DatabaseBase) QueryMap(sql string, params []interface{}) (list []map[string]interface{}, err error) {

	finder := zorm.NewFinder()
	finder.Append(sql, params...)

	list, err = this_.FinderQueryMap(finder)

	if err != nil {
		Logger.Error("QueryMap Error", zap.Any("sql", sql), zap.Any("params", params), zap.Error(err))
		return
	}

	return
}

func (this_ *DatabaseBase) FinderQueryMap(finder *zorm.Finder) (list []map[string]interface{}, err error) {

	list, err = zorm.QueryMap(this_.GetContext(), finder, nil)

	if err != nil {
		return
	}

	return
}

func (this_ *DatabaseBase) QueryPage(sql string, params []interface{}, list interface{}, page *zorm.Page) (err error) {

	finder := zorm.NewFinder()
	finder.Append(sql, params...)

	err = this_.FinderQueryPage(finder, list, page)

	if err != nil {
		Logger.Error("QueryPage Error", zap.Any("sql", sql), zap.Any("params", params), zap.Error(err))
		return
	}

	return
}

func (this_ *DatabaseBase) FinderQueryPage(finder *zorm.Finder, list interface{}, page *zorm.Page) (err error) {

	err = zorm.Query(this_.GetContext(), finder, list, page)

	if err != nil {
		return
	}

	return
}

func (this_ *DatabaseBase) QueryMapPage(sql string, params []interface{}, page *zorm.Page) (list []map[string]interface{}, err error) {

	finder := zorm.NewFinder()
	finder.Append(sql, params...)

	list, err = this_.FinderQueryMapPage(finder, page)

	if err != nil {
		Logger.Error("QueryMap Error", zap.Any("sql", sql), zap.Any("params", params), zap.Error(err))
		return
	}

	return
}

func (this_ *DatabaseBase) FinderQueryMapPage(finder *zorm.Finder, page *zorm.Page) (list []map[string]interface{}, err error) {

	list, err = zorm.QueryMap(this_.GetContext(), finder, page)

	if err != nil {
		return
	}

	return
}
