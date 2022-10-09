package model

type StepHttpModel struct {
	*StepModel

	Http string `json:"http,omitempty"` // HTTP操作
	Url  string `json:"url,omitempty"`  // HTTP地址
}
