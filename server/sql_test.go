package main

import (
	"encoding/json"
	"fmt"
	"server/model/modelSql"
	"testing"
)

func TestSelect(t *testing.T) {
	model := &modelSql.Select{
		Table: "TM_USER",
		Columns: []*modelSql.SelectColumn{
			{Name: "userId"},
			{Name: "name"},
			{Name: "email"},
		},
		Wheres: []*modelSql.Where{
			{Name: "account"},
			{Piece: true,
				Wheres: []*modelSql.Where{
					{Name: "name"},
					{Name: "email"},
				},
			},
		},
		UnionSelects: []*modelSql.Select{
			{
				Table: "TM_USER",
				Columns: []*modelSql.SelectColumn{
					{Name: "userId"},
					{Name: "name"},
					{Name: "email"},
				},
				Wheres: []*modelSql.Where{
					{Name: "name"},
					{Name: "email"},
				},
			},
		},
	}

	data := make(map[string]interface{})
	data["id"] = 1
	data["name"] = "户名1"
	data["email"] = 20
	data["account"] = "户名1"
	sqlParam, err := model.GetSqlParam(data)
	if err != nil {
		panic(err)
	}
	fmt.Println("Sql  :", sqlParam.Sql)
	fmt.Println("Param:", ToJSON(sqlParam.Params))
	fmt.Println("Exec :", sqlParam.ToExecSql())
}

func TestInser(t *testing.T) {
	model := &modelSql.Insert{
		Table: "TM_USER",
		Columns: []*modelSql.InsertColumn{
			{Name: "id", AutoIncrement: true},
			{Name: "name"},
			{Name: "age"},
		},
	}

	data := make(map[string]interface{})
	data["name"] = "户名1"
	data["age"] = 18
	sqlParam, err := model.GetSqlParam(data)
	if err != nil {
		panic(err)
	}
	fmt.Println("Sql  :", sqlParam.Sql)
	fmt.Println("Param:", ToJSON(sqlParam.Params))
	fmt.Println("Exec :", sqlParam.ToExecSql())
}

func TestUpdate(t *testing.T) {
	model := &modelSql.Update{
		Table: "TM_USER",
		Columns: []*modelSql.UpdateColumn{
			{Name: "c1"},
			{Name: "c2"},
		},
		Wheres: []*modelSql.Where{
			{Name: "c3"},
			{Piece: true,
				Wheres: []*modelSql.Where{
					{Name: "c1"},
					{Name: "c2"},
				},
			},
		},
	}

	data := make(map[string]interface{})
	data["id"] = 1
	data["c1"] = "户名1"
	data["c3"] = 20
	data["c2"] = "户名1"
	sqlParam, err := model.GetSqlParam(data)
	if err != nil {
		panic(err)
	}
	fmt.Println("Sql  :", sqlParam.Sql)
	fmt.Println("Param:", ToJSON(sqlParam.Params))
	fmt.Println("Exec :", sqlParam.ToExecSql())
}

func TestDelete(t *testing.T) {
	model := &modelSql.Delete{
		Table: "TM_USER",
		Wheres: []*modelSql.Where{
			{Name: "c3"},
			{Piece: true,
				Wheres: []*modelSql.Where{
					{Name: "c1"},
					{Name: "c2"},
				},
			},
		},
	}

	data := make(map[string]interface{})
	data["id"] = 1
	data["c1"] = "户名1"
	data["c3"] = 20
	data["c2"] = "户名1"
	sqlParam, err := model.GetSqlParam(data)
	if err != nil {
		panic(err)
	}
	fmt.Println("Sql  :", sqlParam.Sql)
	fmt.Println("Param:", ToJSON(sqlParam.Params))
	fmt.Println("Exec :", sqlParam.ToExecSql())
}

func ToJSON(data interface{}) string {
	if data != nil {
		bs, _ := json.Marshal(data)
		if bs != nil {
			return string(bs)
		}
	}
	return ""
}

func ToBean(bytes []byte, req interface{}) (err error) {
	err = json.Unmarshal(bytes, req)
	return
}
