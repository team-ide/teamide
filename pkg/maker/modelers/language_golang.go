package modelers

import "strings"

type LanguageGolangModel struct {
	ElementNode
	Dir          string `json:"dir,omitempty"` // 常量文件
	ModuleName   string `json:"moduleName,omitempty"`
	ConstantPath string `json:"constantPath,omitempty"` // 常量文件
	GoVersion    string `json:"goVersion,omitempty"`
}

func (this_ *LanguageGolangModel) GetModuleName() string {
	if this_.ConstantPath != "" {
		return this_.ConstantPath
	}
	return "app"
}

func (this_ *LanguageGolangModel) GetGoVersion() string {
	if this_.GoVersion != "" {
		return this_.GoVersion
	}
	return "1.18"
}

func (this_ *LanguageGolangModel) GetConstantDir(dir string) string {
	if this_.ConstantPath == "" {
		this_.ConstantPath = "constant/"
	} else {
		if !strings.HasSuffix(this_.ConstantPath, "/") {
			this_.ConstantPath += "/"
		}
	}
	return dir + this_.ConstantPath
}

func init() {
	addDocTemplate(&docTemplate{
		Name:    TypeLanguageGolangName,
		Comment: "语言-Golang",
		Fields: []*docTemplateField{
			{Name: "dir", Comment: "目录"},
			{Name: "moduleName", Comment: "module名称"},
			{Name: "constantPath", Comment: "常量目录路径"},
		},
	})
}
