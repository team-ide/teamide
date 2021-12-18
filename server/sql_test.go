package main

import (
	"encoding/json"
	"fmt"
	"server/model/sqlModel"
	"testing"
)

func TestSelect(t *testing.T) {
	model := &sqlModel.Select{
		Table: "TM_USER",
		Columns: []*sqlModel.SelectColumn{
			{Name: "userId"},
			{Name: "name"},
			{Name: "email"},
		},
		Wheres: []*sqlModel.Where{
			{Name: "account"},
			{Piece: true,
				Wheres: []*sqlModel.Where{
					{Name: "name"},
					{Name: "email"},
				},
			},
		},
		UnionSelects: []*sqlModel.Select{
			{
				Table: "TM_USER",
				Columns: []*sqlModel.SelectColumn{
					{Name: "userId"},
					{Name: "name"},
					{Name: "email"},
				},
				Wheres: []*sqlModel.Where{
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
	model := &sqlModel.Insert{
		Table: "TM_USER",
		Columns: []*sqlModel.InsertColumn{
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
	model := &sqlModel.Update{
		Table: "TM_USER",
		Columns: []*sqlModel.UpdateColumn{
			{Name: "c1"},
			{Name: "c2"},
		},
		Wheres: []*sqlModel.Where{
			{Name: "c3"},
			{Piece: true,
				Wheres: []*sqlModel.Where{
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
	model := &sqlModel.Delete{
		Table: "TM_USER",
		Wheres: []*sqlModel.Where{
			{Name: "c3"},
			{Piece: true,
				Wheres: []*sqlModel.Where{
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
