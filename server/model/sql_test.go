package model

import (
	"encoding/json"
	"fmt"
	"server/model/sql"
	"testing"
)

func TestSelect(t *testing.T) {
	model := &sql.Select{
		Table: "TM_USER",
		Columns: []*sql.SelectColumn{
			{Name: "userId"},
			{Name: "name"},
			{Name: "email"},
		},
		Wheres: []*sql.Where{
			{Name: "account"},
			{Piece: true,
				Wheres: []*sql.Where{
					{Name: "name"},
					{Name: "email"},
				},
			},
		},
		UnionSelects: []*sql.Select{
			{
				Table: "TM_USER",
				Columns: []*sql.SelectColumn{
					{Name: "userId"},
					{Name: "name"},
					{Name: "email"},
				},
				Wheres: []*sql.Where{
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
	model := &sql.Insert{
		Table: "TM_USER",
		Columns: []*sql.InsertColumn{
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
	model := &sql.Update{
		Table: "TM_USER",
		Columns: []*sql.UpdateColumn{
			{Name: "c1"},
			{Name: "c2"},
		},
		Wheres: []*sql.Where{
			{Name: "c3"},
			{Piece: true,
				Wheres: []*sql.Where{
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
	model := &sql.Delete{
		Table: "TM_USER",
		Wheres: []*sql.Where{
			{Name: "c3"},
			{Piece: true,
				Wheres: []*sql.Where{
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
