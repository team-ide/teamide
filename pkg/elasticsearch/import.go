package elasticsearch

import (
	"errors"
	"teamide/pkg/data_engine"
	"teamide/pkg/util"
)

type ImportTask struct {
	*Task
	ImportType  string                `json:"importType,omitempty"`
	Count       int                   `json:"count,omitempty"`
	BatchNumber int                   `json:"batchNumber,omitempty"`
	Id          string                `json:"id,omitempty"`
	ColumnList  []*StrategyDataColumn `json:"columnList"`
}

type StrategyDataColumn struct {
	Name  string `json:"name,omitempty"`
	Value string `json:"value,omitempty"`
}

func (this_ *ImportTask) do() (err error) {
	if this_.ImportType == "strategy" {
		err = this_.doStrategy()
		if err != nil {
			return
		}
	}
	return
}
func (this_ *ImportTask) doStrategy() (err error) {

	if this_.IndexName == "" {
		err = errors.New("必须配置indexName")
		return
	}
	if this_.Count <= 0 {
		return
	}
	if this_.needStop() {
		return
	}

	task := &data_engine.StrategyTask{}

	taskStrategyData := &data_engine.StrategyData{}

	task.StrategyDataList = append(task.StrategyDataList, taskStrategyData)

	taskStrategyData.Count = this_.Count
	taskStrategyData.AddField(&data_engine.StrategyDataField{
		Name:  "id",
		Value: this_.Id,
	})
	for _, column := range this_.ColumnList {
		if column.Name == "" {
			continue
		}
		taskStrategyData.AddField(&data_engine.StrategyDataField{
			Name:  column.Name,
			Value: column.Value,
		})
	}

	task.OnError = func(onErr error) {
		err = onErr
	}

	this_.DataCount += this_.Count

	this_.taskList = append(this_.taskList, task)

	var docs []*InsertDoc
	var batchNumber = this_.BatchNumber
	if batchNumber <= 0 {
		batchNumber = 200
	}

	task.OnData = func(onData map[string]interface{}) (err error) {

		if this_.needStop() {
			return
		}

		doc := &InsertDoc{
			IndexName: this_.IndexName,
		}
		doc.Id, err = util.GetStringValue(onData["id"])
		if err != nil {
			return
		}

		this_.DataReadyCount++

		doc.Doc = onData

		docs = append(docs, doc)
		var size = len(docs)
		if size >= batchNumber {
			_, err = this_.Service.BatchInsertNotWait(docs)
			docs = []*InsertDoc{}
			if err != nil {
				this_.DataErrorCount += size
				return
			}
			this_.DataSuccessCount += size
		}

		return
	}

	task.OnEnd = func() {
		var size = len(docs)
		_, err = this_.Service.BatchInsertNotWait(docs)
		docs = []*InsertDoc{}
		if err != nil {
			this_.DataErrorCount += size
			return
		}
		this_.DataSuccessCount += size
	}

	task.Start()

	return
}
