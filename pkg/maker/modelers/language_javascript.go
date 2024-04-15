package modelers

import "strings"

type LanguageJavascriptModel struct {
	Dir          string `json:"dir,omitempty"`          // 常量文件
	ConstantDir  string `json:"constantDir,omitempty"`  // 常量文件
	ConstantName string `json:"constantName,omitempty"` // 常量文件
}

func (this_ *LanguageJavascriptModel) GetDir() string {
	if this_.Dir != "" {
		if !strings.HasSuffix(this_.Dir, "/") {
			this_.Dir = this_.Dir + "/"
		}
		return this_.Dir
	}
	return "src/"
}

func (this_ *LanguageJavascriptModel) GetConstantDir() string {
	if this_.ConstantDir != "" {
		return this_.ConstantDir
	}
	return this_.GetDir() + "constant/"
}

func (this_ *LanguageJavascriptModel) GetConstantName() string {
	if this_.ConstantName != "" {
		return this_.ConstantName
	}
	return "constant.js"
}

func init() {
	addDocTemplate(&docTemplate{
		Name:    TypeLanguageJavascriptName,
		Comment: "语言-JavaScript",
		Fields: []*docTemplateField{
			{Name: "dir", Comment: "目录"},
			{Name: "constantDir", Comment: "常量文件目录"},
			{Name: "constantName", Comment: "常量文件名称"},
		},
	})
}
