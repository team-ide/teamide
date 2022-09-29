package model

import (
	"encoding/json"
	"fmt"
	"go.uber.org/zap"
	"teamide/pkg/util"
	"testing"
)

func TestDoc(t *testing.T) {

	model := &StructModel{
		Name: "user",
		Note: `用户结构体`,
		Fields: []*StructField{
			{
				Name:    "name",
				Comment: "用户名称",
			},
			{
				Name:    "account",
				Comment: "账号",
			},
			{
				Name: "phone",
			},
			{
				Name: "photo",
			},
			{
				Name: "age",
			},
		},
	}

	text, err := StructToText(model)
	if err != nil {
		util.Logger.Error("StructToText error", zap.Any("model", model), zap.Error(err))
		return
	}

	util.Logger.Info("StructToText success")

	fmt.Println(text)

	model, err = TextToStruct(text)
	if err != nil {
		util.Logger.Error("TextToStruct error", zap.Any("text", text), zap.Error(err))
		return
	}
	util.Logger.Info("TextToStruct success", zap.Any("model", model))
	bs, _ := json.Marshal(model)

	fmt.Println(string(bs))
	//data := map[string]interface{}{
	//	"name": "这是名称",
	//}
	//
	//bs, _ := yaml.Marshal(data)
	//
	//var node = &yaml.Node{}
	//_ = yaml.Unmarshal(bs, node)
	//
	//util.Logger.Info("unmarshal to node", zap.Any("node", node))
	//
	//bs, _ = yaml.Marshal(node)
	//util.Logger.Info("node marshal to data", zap.Any("data", string(bs)))
	//var waitGroupForStop sync.WaitGroup
	//waitGroupForStop.Add(1)
	//waitGroupForStop.Wait()
}

var _ = `
{
	"node": {
		"Kind": 1,
		"Style": 0,
		"Tag": "",
		"Value": "",
		"Anchor": "",
		"Alias": null,
		"Content": [{
			"Kind": 4,
			"Style": 0,
			"Tag": "!!map",
			"Value": "",
			"Anchor": "",
			"Alias": null,
			"Content": [{
				"Kind": 8,
				"Style": 0,
				"Tag": "!!str",
				"Value": "name",
				"Anchor": "",
				"Alias": null,
				"Content": null,
				"HeadComment": "",
				"LineComment": "",
				"FootComment": "",
				"Line": 1,
				"Column": 1
			}, {
				"Kind": 8,
				"Style": 0,
				"Tag": "!!str",
				"Value": "这是名称",
				"An    chor": "",
				"Alias": null,
				"Content": null,
				"HeadComment": "",
				"LineComment": "",
				"FootComment": "",
				"Line": 1,
				"Column": 7
			}],
			"HeadComment": "",
			"LineComment": "",
			"FootComment": "",
			"Line": 1,
			"Column": 1
		}],
		"HeadComment": "",
		"LineComment": "",
		"FootComment": "",
		"Line": 1,
		"Column": 1
	}
}
`
