package model

type ApiModel struct {
	Name     string       `json:"name,omitempty"`                               // 名称，同一个应用中唯一
	Request  *ApiRequest  `json:"request,omitempty" yaml:"request,omitempty"`   //
	Response *ApiResponse `json:"response,omitempty" yaml:"response,omitempty"` //
}

type ApiRequest struct {
	Method      string `json:"method,omitempty" yaml:"method,omitempty"`           //
	Path        string `json:"path,omitempty" yaml:"path,omitempty"`               //
	ContentType string `json:"contentType,omitempty" yaml:"contentType,omitempty"` //
}

type ApiResponse struct {
	DownloadFile bool   `json:"downloadFile,omitempty" yaml:"downloadFile,omitempty"` //
	OpenFile     bool   `json:"openFile,omitempty" yaml:"openFile,omitempty"`         //
	FileName     string `json:"fileName,omitempty" yaml:"fileName,omitempty"`         //
}
