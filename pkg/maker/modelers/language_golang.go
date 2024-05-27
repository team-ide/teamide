package modelers

import "strings"

type LanguageGolangModel struct {
	ElementNode
	Dir          string `json:"dir,omitempty"`
	ModuleName   string `json:"moduleName,omitempty"`
	GoVersion    string `json:"goVersion,omitempty"`
	CommonPath   string `json:"commonPath,omitempty"`
	CommonPack   string `json:"commonPack,omitempty"`
	ConstantPath string `json:"constantPath,omitempty"`
	ConstantPack string `json:"constantPack,omitempty"`
	ErrorPath    string `json:"errorPath,omitempty"`
	ErrorPack    string `json:"errorPack,omitempty"`
	StructPath   string `json:"structPath,omitempty"`
	StructPack   string `json:"structPack,omitempty"`
	FuncPath     string `json:"funcPath,omitempty"`
	FuncPack     string `json:"funcPack,omitempty"`
	DaoPath      string `json:"daoPath,omitempty"`
	DaoPack      string `json:"daoPack,omitempty"`
	ServicePath  string `json:"servicePath,omitempty"`
	ServicePack  string `json:"servicePack,omitempty"`

	ComponentDbPath string `json:"componentDbPath,omitempty"`
	ComponentDbPack string `json:"componentDbPack,omitempty"`

	ComponentRedisPath string `json:"componentRedisPath,omitempty"`
	ComponentRedisPack string `json:"componentRedisPack,omitempty"`
}

func (this_ *LanguageGolangModel) GetModuleName() string {
	if this_.ModuleName != "" {
		return this_.ModuleName
	}
	return "app"
}

func (this_ *LanguageGolangModel) GetGoVersion() string {
	if this_.GoVersion != "" {
		return this_.GoVersion
	}
	return "1.18"
}

func (this_ *LanguageGolangModel) GetCommonDir(dir string) string {
	return GetDir(dir, this_.GetCommonPath())
}

func (this_ *LanguageGolangModel) GetCommonPath() string {
	return GetPath(&this_.CommonPath, "common/")
}

func (this_ *LanguageGolangModel) GetCommonPack() string {
	return GetPack(&this_.CommonPack, "common")
}

func (this_ *LanguageGolangModel) GetCommonImport() string {
	return this_.GetPackImport(this_.GetCommonPath(), this_.GetCommonPack())
}

func (this_ *LanguageGolangModel) GetConstantDir(dir string) string {
	return GetDir(dir, this_.GetConstantPath())
}

func (this_ *LanguageGolangModel) GetConstantPath() string {
	return GetPath(&this_.ConstantPath, "constant/")
}

func (this_ *LanguageGolangModel) GetConstantPack() string {
	return GetPack(&this_.ConstantPack, "constant")
}

func (this_ *LanguageGolangModel) GetConstantImport() string {
	return this_.GetPackImport(this_.GetConstantPath(), this_.GetConstantPack())
}

func (this_ *LanguageGolangModel) GetErrorDir(dir string) string {
	return GetDir(dir, this_.GetErrorPath())
}

func (this_ *LanguageGolangModel) GetErrorPath() string {
	return GetPath(&this_.ErrorPath, "exception/")
}

func (this_ *LanguageGolangModel) GetErrorPack() string {
	return GetPack(&this_.ErrorPack, "exception")
}
func (this_ *LanguageGolangModel) GetErrorImport() string {
	return this_.GetPackImport(this_.GetErrorPath(), this_.GetErrorPack())
}

func (this_ *LanguageGolangModel) GetStructDir(dir string) string {
	return GetDir(dir, this_.GetStructPath())
}

func (this_ *LanguageGolangModel) GetStructPath() string {
	return GetPath(&this_.StructPath, "bean/")
}

func (this_ *LanguageGolangModel) GetStructPack() string {
	return GetPack(&this_.StructPack, "bean")
}
func (this_ *LanguageGolangModel) GetStructImport() string {
	return this_.GetPackImport(this_.GetStructPath(), this_.GetStructPack())
}

func (this_ *LanguageGolangModel) GetFuncDir(dir string) string {
	return GetDir(dir, this_.GetFuncPath())
}

func (this_ *LanguageGolangModel) GetFuncPath() string {
	return GetPath(&this_.FuncPath, "tool/")
}

func (this_ *LanguageGolangModel) GetFuncPack() string {
	return GetPack(&this_.FuncPack, "tool")
}
func (this_ *LanguageGolangModel) GetFuncImport() string {
	return this_.GetPackImport(this_.GetFuncPath(), this_.GetFuncPack())
}

func (this_ *LanguageGolangModel) GetDaoDir(dir string) string {
	return GetDir(dir, this_.GetDaoPath())
}

func (this_ *LanguageGolangModel) GetDaoPath() string {
	return GetPath(&this_.DaoPath, "dao/")
}

func (this_ *LanguageGolangModel) GetDaoPack() string {
	return GetPack(&this_.DaoPack, "dao")
}
func (this_ *LanguageGolangModel) GetDaoImport() string {
	return this_.GetPackImport(this_.GetDaoPath(), this_.GetDaoPack())
}

func (this_ *LanguageGolangModel) GetServiceDir(dir string) string {
	return GetDir(dir, this_.GetServicePath())
}

func (this_ *LanguageGolangModel) GetServicePath() string {
	return GetPath(&this_.ServicePath, "service/")
}

func (this_ *LanguageGolangModel) GetServicePack() string {
	return GetPack(&this_.ServicePack, "service")
}

func (this_ *LanguageGolangModel) GetServiceImport() string {
	return this_.GetPackImport(this_.GetServicePath(), this_.GetServicePack())
}

func (this_ *LanguageGolangModel) GetComponentDbDir(dir string) string {
	return GetDir(dir, this_.GetComponentDbPath())
}

func (this_ *LanguageGolangModel) GetComponentDbPath() string {
	return GetPath(&this_.ComponentDbPath, "component_db/")
}

func (this_ *LanguageGolangModel) GetComponentDbPack() string {
	return GetPack(&this_.ComponentDbPack, "component_db")
}

func (this_ *LanguageGolangModel) GetComponentDbImport() string {
	return this_.GetPackImport(this_.GetComponentDbPath(), this_.GetComponentDbPack())
}

func (this_ *LanguageGolangModel) GetComponentRedisDir(dir string) string {
	return GetDir(dir, this_.GetComponentRedisPath())
}

func (this_ *LanguageGolangModel) GetComponentRedisPath() string {
	return GetPath(&this_.ComponentRedisPath, "component_redis/")
}

func (this_ *LanguageGolangModel) GetComponentRedisPack() string {
	return GetPack(&this_.ComponentRedisPack, "component_redis")
}

func (this_ *LanguageGolangModel) GetComponentRedisImport() string {
	return this_.GetPackImport(this_.GetComponentRedisPath(), this_.GetComponentRedisPack())
}

func (this_ *LanguageGolangModel) GetPackImport(path string, pack string) string {
	path = strings.TrimPrefix(path, "/")
	path = strings.TrimSuffix(path, "/")
	moduleName := this_.GetModuleName()
	dot := strings.LastIndex(path, "/")
	if dot > 0 {
		moduleName += path[:dot]
	}
	return moduleName + "/" + pack
}

func GetDir(dir string, path string) string {
	return dir + path
}

func GetPath(name *string, defaultPath string) string {
	if *name == "" {
		*name = defaultPath
	} else {
		if !strings.HasSuffix(*name, "/") {
			*name += "/"
		}
	}
	return *name
}

func GetPack(name *string, defaultPack string) string {
	if *name == "" {
		*name = defaultPack
	}

	return *name
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
