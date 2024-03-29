package module_thrift

import (
	"encoding/json"
	"fmt"
	"github.com/team-ide/go-tool/metric"
)

func toMarkdown(requestMd5 string, taskList []map[string]interface{}) (content string) {

	var groupList []*[]map[string]interface{}
	groupCache := map[string]*[]map[string]interface{}{}
	for _, one := range taskList {
		if one["requestMd5"] == nil {
			continue
		}
		requestMd5_ := one["requestMd5"].(string)
		if requestMd5 != "" && requestMd5 != requestMd5_ {
			continue
		}
		group := groupCache[requestMd5_]
		if group == nil {
			group = &[]map[string]interface{}{}
			groupCache[requestMd5] = group
			groupList = append(groupList, group)
		}
		*group = append(*group, one)
	}

	content += fmt.Sprintf("# 测试结果  \n\n")
	for index, group := range groupList {

		content += groupToMarkdown(index, *group)
	}
	return
}

func groupToMarkdown(index int, group []map[string]interface{}) (content string) {
	if len(group) == 0 {
		return
	}
	var bs []byte
	bs, _ = json.Marshal(group[0]["request"])
	request := &BaseRequest{}
	_ = json.Unmarshal(bs, request)

	content += fmt.Sprintf("## 测试组-%d  \n\n", index+1)
	content += fmt.Sprintf("#### 接口信息  \n\n")
	content += fmt.Sprintf("* 服务名称：%s  \n", request.ServiceName)
	content += fmt.Sprintf("* 方法名称：%s  \n", request.MethodName)

	content += fmt.Sprintf("\n")

	content += fmt.Sprintf("#### 测试信息  \n\n")
	content += fmt.Sprintf("* 线程数：%d  \n", request.Worker)
	if request.Frequency > 0 {
		content += fmt.Sprintf("* 执行次数：%d  \n", request.Frequency)
	} else {
		content += fmt.Sprintf("* 执行时长：%d  \n", request.Duration)
	}
	content += fmt.Sprintf("* 测试地址：%s  \n", request.ServerAddress)
	content += fmt.Sprintf("* 超时时长：%d  \n", request.Timeout)
	content += fmt.Sprintf("* ProtocolFactory类型：%s  \n", request.ProtocolFactory)
	content += fmt.Sprintf("* Buffered：%v  \n", request.Buffered)
	content += fmt.Sprintf("* Framed：%v  \n", request.Framed)

	content += fmt.Sprintf("\n")
	for i, arg := range request.Args {
		content += fmt.Sprintf("* 参数-%d：  \n\n", i+1)
		content += fmt.Sprintf("```json\n")
		content += arg
		content += fmt.Sprintf("\n")
		content += fmt.Sprintf("```\n\n")
	}

	content += fmt.Sprintf("\n\n")
	content += fmt.Sprintf("#### 测试记录  \n\n")
	content += fmt.Sprintf("* 任务用时：任务的开始时间~结束时间耗时； \n")
	content += fmt.Sprintf("* 执行用时：单个线程执行用时累计，取最大；（这里的用时是调用接口耗时，去除了额外开销，所以执行用时小于任务执行时间，两者相差越大，则表示额外开销越多） \n")
	content += fmt.Sprintf("* 累计用时：所有执行用时累计 \n")
	content += fmt.Sprintf("* TPS：总次数 / 任务用时 \n")

	content += fmt.Sprintf("\n")

	var cs []*metric.Count
	for _, task := range group {
		bs, _ = json.Marshal(task["metric"])
		count := &metric.Count{}
		_ = json.Unmarshal(bs, count)
		cs = append(cs, count)
	}
	content += metric.MarkdownTable(cs, &metric.Options{
		AddHtmlFormat: true,
		WarnUseTime:   1000,
	})
	content += fmt.Sprintf("\n\n")
	return
}

type tS struct {
	Size int64
	Unit string
}

var (
	tList = []*tS{
		{Size: 1000 * 60 * 60 * 24, Unit: "天"},
		{Size: 1000 * 60 * 60, Unit: "时"},
		{Size: 1000 * 60, Unit: "分"},
		{Size: 1000, Unit: "秒"},
	}
)

func toTime(size int64) (v string) {

	var timeV = size

	for _, s := range tList {
		if timeV >= s.Size {
			tV := timeV / s.Size
			timeV -= tV * s.Size
			v += fmt.Sprintf("%d%s", tV, s.Unit)
		}
	}
	if timeV > 0 {
		v += fmt.Sprintf("%d%s", timeV, "毫秒")
	}
	return
}
