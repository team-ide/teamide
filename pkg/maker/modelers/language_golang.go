package modelers

import "strings"

type LanguageGolangModel struct {
	ElementNode
	Dir          string `json:"dir,omitempty"`          // 常量文件
	ConstantDir  string `json:"constantDir,omitempty"`  // 常量文件
	ConstantName string `json:"constantName,omitempty"` // 常量文件
}

func (this_ *LanguageGolangModel) GetDir() string {
	if this_.Dir != "" {
		if !strings.HasSuffix(this_.Dir, "/") {
			this_.Dir = this_.Dir + "/"
		}
		return this_.Dir
	}
	return "src/"
}

func (this_ *LanguageGolangModel) GetConstantDir() string {
	if this_.ConstantDir != "" {
		return this_.ConstantDir
	}
	return this_.GetDir() + "constant/"
}

func (this_ *LanguageGolangModel) GetConstantName() string {
	if this_.ConstantName != "" {
		return this_.ConstantName
	}
	return "constant.js"
}

func init() {
	addDocTemplate(&docTemplate{
		Name:    TypeLanguageGolangName,
		Comment: "语言-Golang",
		Fields: []*docTemplateField{
			{Name: "dir", Comment: "目录"},
			{Name: "constantDir", Comment: "常量文件目录"},
			{Name: "constantName", Comment: "常量文件名称"},
		},
	})
}
