package module_http

import (
	"net/url"
	"path"
	"strings"
)

type ContentInfo struct {
	IsBinary     bool   `json:"isBinary,omitempty"`
	IsText       bool   `json:"isText,omitempty"`  // 文本
	IsImage      bool   `json:"isImage,omitempty"` // 图片
	IsVideo      bool   `json:"isVideo,omitempty"` // 视频
	IsAudio      bool   `json:"isAudio,omitempty"` // 音频
	IsJson       bool   `json:"isJson,omitempty"`  // JSON
	IsXml        bool   `json:"isXml,omitempty"`   // XML
	IsPdf        bool   `json:"isPdf,omitempty"`   // PDF
	IsZip        bool   `json:"isZip,omitempty"`   // Zip
	IsGzip       bool   `json:"isGzip,omitempty"`  // Gzip
	IsFile       bool   `json:"isFile,omitempty"`  // 是文件
	IsWord       bool   `json:"isWord,omitempty"`  // Microsoft Word文档
	IsExcel      bool   `json:"isExcel,omitempty"` // Microsoft Excel文档
	IsHtml       bool   `json:"isHtml,omitempty"`
	IsJavascript bool   `json:"isJavascript,omitempty"`
	IsCss        bool   `json:"isCss,omitempty"`
	IsFont       bool   `json:"isFont,omitempty"`
	FileType     string `json:"fileType,omitempty"` // 文件类型
}

func GetContentInfo(contentType string) (info *ContentInfo) {
	info = new(ContentInfo)
	if strings.Contains(contentType, "text/") {
		info.IsText = true
		info.IsFile = true
	}
	if strings.Contains(contentType, "/html") {
		info.IsHtml = true
		info.IsText = true
		//info.IsFile = true
	}
	if strings.Contains(contentType, "/css") {
		info.IsCss = true
		info.IsText = true
		//info.IsFile = true
	}
	if strings.Contains(contentType, "/javascript") {
		info.IsJavascript = true
		info.IsText = true
		//info.IsFile = true
	}
	if strings.Contains(contentType, "font/") {
		info.IsBinary = true
		info.IsFile = true
	}
	if strings.Contains(contentType, "image/") {
		info.IsImage = true
		info.IsBinary = true
		info.IsFile = true
	}
	if strings.Contains(contentType, "video/") {
		info.IsVideo = true
		info.IsBinary = true
		info.IsFile = true
	}
	if strings.Contains(contentType, "audio/") {
		info.IsAudio = true
		info.IsBinary = true
		info.IsFile = true
	}
	if strings.Contains(contentType, "/octet-stream") {
		info.IsBinary = true
		info.IsFile = true
	}
	if strings.Contains(contentType, "/json") {
		info.IsJson = true
		info.IsText = true
		//info.IsFile = true
	}
	if strings.Contains(contentType, "/xml") {
		info.IsXml = true
		info.IsText = true
		//info.IsFile = true
	}
	if strings.Contains(contentType, "/pdf") {
		info.IsPdf = true
		info.IsBinary = true
		info.IsFile = true
	}
	if strings.Contains(contentType, "/zip") {
		info.IsZip = true
		info.IsBinary = true
		info.IsFile = true
	}
	if strings.Contains(contentType, "/x-gzip") {
		info.IsGzip = true
		info.IsBinary = true
		info.IsFile = true
	}
	if strings.Contains(contentType, "/msword") {
		info.IsWord = true
		info.IsBinary = true
		info.IsFile = true
	}
	if strings.Contains(contentType, "/vnd.ms-excel") {
		info.IsExcel = true
		info.IsBinary = true
		info.IsFile = true
	}
	return
}
func GetFileNameFromURL(urlString string) string {
	u, err := url.Parse(urlString)
	if err != nil {
		return ""
	}
	fileName := path.Base(u.Path)
	return fileName
}
