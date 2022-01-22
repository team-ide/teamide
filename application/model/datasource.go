package model

type DatasourceDatabase struct {
	Name         string `json:"name,omitempty"`                             // 名称，同一个应用中唯一
	Comment      string `json:"comment,omitempty" yaml:"comment,omitempty"` // 注释说明
	Type         string `json:"type,omitempty" yaml:"type,omitempty"`
	Host         string `json:"host,omitempty" yaml:"host,omitempty"`
	Port         int    `json:"port,omitempty" yaml:"port,omitempty"`
	Sid          string `json:"sid,omitempty" yaml:"sid,omitempty"`
	Database     string `json:"database,omitempty" yaml:"database,omitempty"`
	Username     string `json:"username,omitempty" yaml:"username,omitempty"`
	Password     string `json:"password,omitempty" yaml:"password,omitempty"`
	CharacterSet string `json:"characterSet,omitempty" yaml:"characterSet,omitempty"`
	Collate      string `json:"collate,omitempty" yaml:"collate,omitempty"`
}

type DatasourceRedis struct {
	Name    string `json:"name,omitempty"`                             // 名称，同一个应用中唯一
	Comment string `json:"comment,omitempty" yaml:"comment,omitempty"` // 注释说明

	Address string `json:"address,omitempty" yaml:"address,omitempty"`
	Auth    string `json:"auth,omitempty" yaml:"auth,omitempty"`
	Prefix  string `json:"prefix,omitempty" yaml:"prefix,omitempty"` // 前缀
}

type DatasourceKafka struct {
	Name    string `json:"name,omitempty"`                             // 名称，同一个应用中唯一
	Comment string `json:"comment,omitempty" yaml:"comment,omitempty"` // 注释说明

	Address string `json:"address,omitempty" yaml:"address,omitempty"`
	Prefix  string `json:"prefix,omitempty" yaml:"prefix,omitempty"` // 前缀
}

type DatasourceZookeeper struct {
	Name    string `json:"name,omitempty"`                             // 名称，同一个应用中唯一
	Comment string `json:"comment,omitempty" yaml:"comment,omitempty"` // 注释说明

	Address   string `json:"address,omitempty" yaml:"address,omitempty"`
	Namespace string `json:"namespace,omitempty" yaml:"namespace,omitempty"` // 命名空间
}

func TextToDatasourceDatabase(namePath string, text string) (model *DatasourceDatabase, err error) {
	var name string
	model = &DatasourceDatabase{}
	name, err = TextToModel(namePath, text, model)
	if err != nil {
		return
	}
	if name == "default" {
		name = ""
	}
	model.Name = (name)
	return
}

func TextToDatasourceRedis(namePath string, text string) (model *DatasourceRedis, err error) {
	var name string
	model = &DatasourceRedis{}
	name, err = TextToModel(namePath, text, model)
	if err != nil {
		return
	}
	if name == "default" {
		name = ""
	}
	model.Name = (name)
	return
}

func TextToDatasourceKafka(namePath string, text string) (model *DatasourceKafka, err error) {
	var name string
	model = &DatasourceKafka{}
	name, err = TextToModel(namePath, text, model)
	if err != nil {
		return
	}
	if name == "default" {
		name = ""
	}
	model.Name = (name)
	return
}

func TextToDatasourceZookeeper(namePath string, text string) (model *DatasourceZookeeper, err error) {
	var name string
	model = &DatasourceZookeeper{}
	name, err = TextToModel(namePath, text, model)
	if err != nil {
		return
	}
	if name == "default" {
		name = ""
	}
	model.Name = (name)
	return
}
