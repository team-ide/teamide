package data_engine

//import (
//	"github.com/team-ide/go-tool/util"
//	"go.uber.org/zap"
//	"testing"
//)
//
//func TestStrategy(t *testing.T) {
//
//	defer func() {
//		if rec := recover(); rec != nil {
//			util.Logger.Error("数据生成异常", zap.Any("error", rec))
//		}
//	}()
//	task := &StrategyTask{}
//	task.OnData = func(data map[string]interface{}) (err error) {
//		//util.Logger.Info("on data_engine", zap.Any("data_engine", data_engine))
//		return
//	}
//	task.OnError = func(err error) {
//		util.Logger.Error("on error", zap.Error(err))
//		return
//	}
//	task.OnEnd = func() {
//		util.Logger.Info("on end")
//		return
//	}
//	task.StrategyDataList = []*StrategyData{
//		{
//			Count: 10,
//			FieldList: []*StrategyDataField{
//				{Name: "名称", Value: `"用户-" + _$index`},
//				{
//					Name:        "好友列表",
//					ValueCount:  2,
//					ValueIsList: true,
//					FieldList: []*StrategyDataField{
//						{Name: "好友名称", Value: `"好友-" + _$value_index`},
//					},
//				},
//				{
//					Name:        "好友名称",
//					ValueCount:  3,
//					ValueIsList: true,
//					Value:       `"好友-" + _$value_index`,
//				},
//			},
//		},
//	}
//	task.Start()
//}
