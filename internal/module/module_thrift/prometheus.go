package module_thrift

import (
	"bufio"
	"io"
	"net/http"
	"strings"
)

type prometheusData struct {
	Time         int64  `json:"time"`
	Text         string `json:"text"`
	SummaryCount int64  `json:"summaryCount"`
	SummarySum   int64  `json:"summarySum"` // 毫秒
}

type prometheusDataCollect struct {
	*BaseRequest
	data []*prometheusData
}

func (this_ *prometheusDataCollect) start() {

}

func (this_ *prometheusDataCollect) collect() {

	//nowTime := util.GetNow()

	res, err := http.Get(this_.PrometheusMetricsAddress)
	if err != nil {
		return
	}
	defer func() { _ = res.Body.Close() }()
	if res.StatusCode != 200 {
		return
	}

	lines, err := ReadLine(res.Body)
	if err != nil {
		return
	}

	for _, line := range lines {
		if strings.HasPrefix(line, this_.PrometheusSummaryCountMatch) {

		} else if strings.HasPrefix(line, this_.PrometheusSummarySumMatch) {

		}
	}

}

// ReadLine 逐行读取文件
func ReadLine(rd io.Reader) (lines []string, err error) {
	buf := bufio.NewReader(rd)
	var line string
	for {
		line, err = buf.ReadString('\n')
		line = strings.TrimSpace(line)
		if err != nil {
			if err == io.EOF { //读取结束，会报EOF
				err = nil
				return
			}
			return nil, err
		}
		lines = append(lines, line)
	}
	return
}
