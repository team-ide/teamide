package db

import (
	"fmt"
	"server/base"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

func init() {
}

type MysqlService struct {
	config DatabaseConfig
	db     *sqlx.DB
}

func CreateMysqlService(config DatabaseConfig) (service *MysqlService, err error) {
	service = &MysqlService{
		config: config,
	}
	err = service.init()
	return
}

func (service *MysqlService) init() (err error) {
	url := fmt.Sprint(service.config.Username, ":", service.config.Password, "@tcp(", service.config.Host, ":", service.config.Port, ")/", service.config.Database, "?charset=utf8mb4&loc=Local&parseTime=true")
	var db *sqlx.DB
	db, err = sqlx.Open("mysql", url)
	db.SetMaxIdleConns(10)
	db.SetMaxOpenConns(100)
	service.db = db
	return
}

func (service *MysqlService) Open() (err error) {
	err = service.db.Ping()
	return
}

func (service *MysqlService) Close() (err error) {
	err = service.db.Close()
	return
}

func (service *MysqlService) Query(sqlParam SqlParam, newBean func() interface{}) (list []interface{}, err error) {
	rows, err := service.db.Query(sqlParam.Sql, sqlParam.Params...)
	if err != nil {
		fmt.Println("Query sql error , sql:", sqlParam.Sql)
		fmt.Println("Query sql error , params:", base.ToJSON(sqlParam.Params))
		return
	}
	list, err = ResultToBeans(rows, newBean)
	if err != nil {
		return
	}
	rows.Close()
	return
}

func (service *MysqlService) Count(sqlParam SqlParam) (count int64, err error) {
	rows, err := service.db.Query(sqlParam.Sql, sqlParam.Params...)
	if err != nil {
		fmt.Println("Count sql error , sql:", sqlParam.Sql)
		fmt.Println("Count sql error , params:", base.ToJSON(sqlParam.Params))
		return
	}
	if rows.Next() {
		err = rows.Scan(&count)
		if err != nil {
			return
		}
	}
	rows.Close()
	return
}

func (service *MysqlService) Insert(sqlParam SqlParam) (rowsAffected int64, err error) {

	result, err := service.db.Exec(sqlParam.Sql, sqlParam.Params...)
	if err != nil {
		fmt.Println("Insert sql error , sql:", sqlParam.Sql)
		fmt.Println("Insert sql error , params:", base.ToJSON(sqlParam.Params))
		return
	}
	rowsAffected, err = result.RowsAffected()
	if err != nil {
		return
	}
	return
}

func (service *MysqlService) Update(sqlParam SqlParam) (rowsAffected int64, err error) {
	result, err := service.db.Exec(sqlParam.Sql, sqlParam.Params...)
	if err != nil {
		fmt.Println("Update sql error , sql:", sqlParam.Sql)
		fmt.Println("Update sql error , params:", base.ToJSON(sqlParam.Params))
		return
	}
	rowsAffected, err = result.RowsAffected()
	if err != nil {
		return
	}
	return
}

func (service *MysqlService) Exec(sqlParam SqlParam) (rowsAffected int64, err error) {
	result, err := service.db.Exec(sqlParam.Sql, sqlParam.Params...)
	if err != nil {
		fmt.Println("Exec sql error , sql:", sqlParam.Sql)
		fmt.Println("Exec sql error , params:", base.ToJSON(sqlParam.Params))
		return
	}
	rowsAffected, err = result.RowsAffected()
	if err != nil {
		return
	}
	return
}

func (service *MysqlService) Delete(sqlParam SqlParam) (rowsAffected int64, err error) {
	result, err := service.db.Exec(sqlParam.Sql, sqlParam.Params...)
	if err != nil {
		fmt.Println("Delete sql error , sql:", sqlParam.Sql)
		fmt.Println("Delete sql error , params:", base.ToJSON(sqlParam.Params))
		return
	}
	rowsAffected, err = result.RowsAffected()
	if err != nil {
		return
	}
	return
}
