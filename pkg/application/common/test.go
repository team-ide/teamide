package common

type TestResult struct {
	Error        error       `json:"error,omitempty"`
	Infos        []*TestInfo `json:"infos,omitempty"`
	Count        int         `json:"count,omitempty"`
	SuccessCount int         `json:"successCount,omitempty"`
	ErrorCount   int         `json:"errorCount,omitempty"`
}

type TestInfo struct {
	ThreadIndex int         `json:"threadIndex,omitempty"`
	ThreadName  string      `json:"threadName,omitempty"`
	ForIndex    int         `json:"forIndex,omitempty"`
	ForName     string      `json:"forName,omitempty"`
	Result      interface{} `json:"result,omitempty"`
	Error       error       `json:"error,omitempty"`
}
