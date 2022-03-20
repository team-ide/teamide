package db

import (
	"fmt"
	"strings"
	"testing"
	"time"
)

func TestSqlite(t *testing.T) {

	config := DatabaseConfig{
		Type:     "sqlite",
		Database: "./TEAM_IDE_TEST",
	}
	all(config)
}

func TestMysql(t *testing.T) {
	config := DatabaseConfig{
		Type:     "mysql",
		Database: "TEAM_IDE_TEST",
		Host:     "mysql",
		Port:     3306,
		Username: "root",
		Password: "123456",
	}
	all(config)
}

func all(config DatabaseConfig) {
	createDatabase(config)
	openAndClose(config)
	createTable(config)
	insert(config)
	update(config)
	query(config)
}

func openAndClose(config DatabaseConfig) {
	databaseWorker, err := NewDatabaseWorker(config)
	if err != nil {
		panic(err)
	}
	err = databaseWorker.Open()
	if err != nil {
		panic(err)
	}
	err = databaseWorker.Close()
	if err != nil {
		panic(err)
	}
}

func createDatabase(config DatabaseConfig) {
	var err error
	switch strings.ToLower(config.Type) {
	case "mysql":
		database := config.Database
		config.Database = ""
		databaseWorker, err := NewDatabaseWorker(config)
		if err != nil {
			panic(err)
		}
		_, err = databaseWorker.Exec("CREATE DATABASE IF NOT EXISTS `"+database+"` CHARACTER SET 'utf8mb4' COLLATE 'utf8mb4_general_ci'", []interface{}{})
		break
	case "sqlite":
		databaseWorker, err := NewDatabaseWorker(config)
		if err != nil {
			panic(err)
		}
		err = databaseWorker.Open()
		break
	}
	if err != nil {
		panic(err)
	}
}

func createTable(config DatabaseConfig) {
	var err error
	databaseWorker, err := NewDatabaseWorker(config)
	_, err = databaseWorker.Exec(`
CREATE TABLE TM_USER (
	userId bigint(20) NOT NULL,
	name varchar(100) NOT NULL,
	createTime datetime NOT NULL,
	updateTime datetime DEFAULT NULL,
	PRIMARY KEY (userId)
);
`, []interface{}{})
	if err != nil {
		panic(err)
	}
}

func insert(config DatabaseConfig) {
	var err error
	databaseWorker, err := NewDatabaseWorker(config)
	_, err = databaseWorker.Exec(`
INSERT INTO TM_USER(userId, name, createTime)
VALUES (?,?,?)
`, []interface{}{1, "张三", time.Now()})
	if err != nil {
		panic(err)
	}
	_, err = databaseWorker.Exec(`
INSERT INTO TM_USER(userId, name, createTime)
VALUES (?,?,?)
`, []interface{}{2, "李四", time.Now()})
	if err != nil {
		panic(err)
	}
}

func update(config DatabaseConfig) {
	var err error
	databaseWorker, err := NewDatabaseWorker(config)
	_, err = databaseWorker.Exec(`
UPDATE TM_USER SET name=?,updateTime=? WHERE userId=?
`, []interface{}{"张三（修改）", time.Now(), 1})
	if err != nil {
		panic(err)
	}
}

func query(config DatabaseConfig) {

	var err error
	databaseWorker, err := NewDatabaseWorker(config)
	list, err := databaseWorker.Query(`
select userId, name, createTime, updateTime from TM_USER
`, []interface{}{}, map[string]string{
		"userId":     "int64",
		"name":       "string",
		"createTime": "time",
		"updateTime": "time",
	})
	if err != nil {
		panic(err)
	}

	for _, one := range list {
		//bs, err := json.Marshal(one)
		//if err != nil {
		//	panic(err)
		//}
		for key, value := range one {
			fmt.Println(key, ":", value)
		}
	}

}
