package modelers

//
//import (
//	"encoding/json"
//	"errors"
//	"fmt"
//	"github.com/team-ide/go-tool/util"
//	"go.uber.org/zap"
//	"testing"
//)
//
//func TestDocService(t *testing.T) {
//	defer func() {
//		if e := recover(); e != nil {
//			util.Logger.Error("TestDocService error", zap.Any("error", e))
//		}
//	}()
//
//	model := &ServiceModel{
//		Name: "user/insert",
//		Note: `用户结构体`,
//		Steps: []interface{}{
//			&StepRedisModel{
//				Redis: "get",
//				Key:   "user-1",
//			},
//		},
//	}
//
//	text, err := ServiceToText(model)
//	if err != nil {
//		util.Logger.Error("ServiceToText error", zap.Any("model", model), zap.Error(err))
//		return
//	}
//
//	util.Logger.Info("ServiceToText success")
//
//	fmt.Println(text)
//
//	model, err = TextToService(text)
//	if err != nil {
//		util.Logger.Error("TextToService error", zap.Any("text", text), zap.Error(err))
//		return
//	}
//	util.Logger.Info("TextToService success", zap.Any("model", model))
//	bs, _ := json.Marshal(model)
//
//	fmt.Println(string(bs))
//
//	newText, err := ServiceToText(model)
//	if err != nil {
//		util.Logger.Error("ServiceToText error", zap.Any("model", model), zap.Error(err))
//		return
//	}
//	if text != newText {
//		err = errors.New("text not eq new text")
//		fmt.Println("text")
//		fmt.Println(text)
//		fmt.Println("new text")
//		fmt.Println(newText)
//		util.Logger.Error("text eq new text error", zap.Any("model", model), zap.Error(err))
//	}
//}
//
//func TestDocStruct(t *testing.T) {
//
//	model := &StructModel{
//		Name: "user",
//		Note: `用户结构体`,
//		Fields: []*StructField{
//			{
//				Name:    "name",
//				Comment: "用户名称",
//			},
//			{
//				Name:    "account",
//				Comment: "账号",
//			},
//			{
//				Name: "phone",
//			},
//			{
//				Name: "photo",
//			},
//			{
//				Name: "age",
//			},
//		},
//	}
//
//	text, err := StructToText(model)
//	if err != nil {
//		util.Logger.Error("StructToText error", zap.Any("model", model), zap.Error(err))
//		return
//	}
//
//	util.Logger.Info("StructToText success")
//
//	fmt.Println(text)
//
//	model, err = TextToStruct(text)
//	if err != nil {
//		util.Logger.Error("TextToStruct error", zap.Any("text", text), zap.Error(err))
//		return
//	}
//	util.Logger.Info("TextToStruct success", zap.Any("model", model))
//	bs, _ := json.Marshal(model)
//
//	fmt.Println(string(bs))
//
//	newText, err := StructToText(model)
//	if err != nil {
//		util.Logger.Error("StructToText error", zap.Any("model", model), zap.Error(err))
//		return
//	}
//	if text != newText {
//		err = errors.New("text not eq new text")
//		fmt.Println("text")
//		fmt.Println(text)
//		fmt.Println("new text")
//		fmt.Println(newText)
//		util.Logger.Error("text eq new text error", zap.Any("model", model), zap.Error(err))
//	}
//}
//
//var _ = `
//{
//	"node": {
//		"Kind": 1,
//		"Style": 0,
//		"Tag": "",
//		"Value": "",
//		"Anchor": "",
//		"Alias": null,
//		"Content": [{
//			"Kind": 4,
//			"Style": 0,
//			"Tag": "!!map",
//			"Value": "",
//			"Anchor": "",
//			"Alias": null,
//			"Content": [{
//				"Kind": 8,
//				"Style": 0,
//				"Tag": "!!str",
//				"Value": "name",
//				"Anchor": "",
//				"Alias": null,
//				"Content": null,
//				"HeadComment": "",
//				"LineComment": "",
//				"FootComment": "",
//				"Line": 1,
//				"Column": 1
//			}, {
//				"Kind": 8,
//				"Style": 0,
//				"Tag": "!!str",
//				"Value": "这是名称",
//				"An    chor": "",
//				"Alias": null,
//				"Content": null,
//				"HeadComment": "",
//				"LineComment": "",
//				"FootComment": "",
//				"Line": 1,
//				"Column": 7
//			}],
//			"HeadComment": "",
//			"LineComment": "",
//			"FootComment": "",
//			"Line": 1,
//			"Column": 1
//		}],
//		"HeadComment": "",
//		"LineComment": "",
//		"FootComment": "",
//		"Line": 1,
//		"Column": 1
//	}
//}
//`
