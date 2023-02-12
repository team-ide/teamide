package module_toolbox

import (
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"strconv"
	"teamide/internal/base"
	"teamide/pkg/db"
	"teamide/pkg/elasticsearch"
	"teamide/pkg/form"
	"teamide/pkg/kafka"
	"teamide/pkg/redis"
	"teamide/pkg/ssh"
	"teamide/pkg/zookeeper"
)

// EncryptOptionAttr 加密属性
func (this_ *ToolboxService) EncryptOptionAttr(str string) (res string) {
	if str == "" {
		return
	}
	// 如果是加密字符串，则直接返回
	if this_.Decryption.IsEncrypt(str) {
		res = str
		return
	} else {
		res, _ = this_.Decryption.Encrypt(str)
		return
	}
}

// DecryptOptionAttr 解密属性
func (this_ *ToolboxService) DecryptOptionAttr(str string) (res string) {
	if str == "" {
		return
	}
	// 如果不是加密字符串，则直接返回
	if !this_.Decryption.IsEncrypt(str) {
		res = str
		return
	}
	res, _ = this_.Decryption.Decrypt(str)
	if res == "" {
		res = str
	}
	return
}

// FormatOption 执行
func (this_ *ToolboxService) FormatOption(toolboxData *ToolboxModel) (err error) {
	if toolboxData.Option == "" {
		return
	}
	if toolboxData.ToolboxType == "" {
		return
	}
	toolboxType := GetToolboxType(toolboxData.ToolboxType)
	if toolboxType == nil {
		err = errors.New("不支持的工具类型[" + toolboxData.ToolboxType + "]")
		return
	}
	optionBytes := []byte(toolboxData.Option)
	optionMap := map[string]interface{}{}
	err = json.Unmarshal(optionBytes, &optionMap)
	if err != nil {
		return
	}

	switch toolboxType {
	case databaseWorker_:
		if optionMap["password"] != nil {
			str, ok := optionMap["password"].(string)
			if ok {
				optionMap["password"] = this_.EncryptOptionAttr(str)
			} else {
				delete(optionMap, "password")
			}
		}
		break
	case redisWorker_:
		if optionMap["auth"] != nil {
			str, ok := optionMap["auth"].(string)
			if ok {
				optionMap["auth"] = this_.EncryptOptionAttr(str)
			} else {
				delete(optionMap, "auth")
			}
		}
		break
	case zookeeperWorker_:
		if optionMap["password"] != nil {
			str, ok := optionMap["password"].(string)
			if ok {
				optionMap["password"] = this_.EncryptOptionAttr(str)
			} else {
				delete(optionMap, "password")
			}
		}
		break
	case elasticsearchWorker_:
		if optionMap["password"] != nil {
			str, ok := optionMap["password"].(string)
			if ok {
				optionMap["password"] = this_.EncryptOptionAttr(str)
			} else {
				delete(optionMap, "password")
			}
		}
		break
	case kafkaWorker_:
		if optionMap["password"] != nil {
			str, ok := optionMap["password"].(string)
			if ok {
				optionMap["password"] = this_.EncryptOptionAttr(str)
			} else {
				delete(optionMap, "password")
			}
		}
		break
	case otherWorker_:
		break
	case sshWorker_:
		if optionMap["password"] != nil {
			str, ok := optionMap["password"].(string)
			if ok {
				optionMap["password"] = this_.EncryptOptionAttr(str)
			} else {
				delete(optionMap, "password")
			}
		}
		break
	}

	optionBytes, err = json.Marshal(optionMap)
	if err != nil {
		return
	}
	toolboxData.Option = string(optionBytes)
	return
}

func (this_ *ToolboxService) GetSSHConfig(option string) (config *ssh.Config, err error) {
	optionBytes := []byte(option)
	err = json.Unmarshal(optionBytes, &config)
	if err != nil {
		return
	}
	config.Password = this_.DecryptOptionAttr(config.Password)
	if config.PublicKey != "" {
		config.PublicKey = this_.GetFilesFile(config.PublicKey)
	}
	return
}

type BindConfigRequest struct {
	ToolboxId   int64  `json:"toolboxId,omitempty"`
	ToolboxType string `json:"toolboxType,omitempty"`
}

func (this_ *ToolboxService) BindConfig(requestBean *base.RequestBean, c *gin.Context, config interface{}) (err error) {
	bindConfigRequest := &BindConfigRequest{}
	if !base.RequestJSON(bindConfigRequest, c) {
		return
	}

	find, err := this_.Get(bindConfigRequest.ToolboxId)
	if err != nil {
		return
	}
	if find != nil && find.UserId != 0 {
		if requestBean.JWT == nil || find.UserId != requestBean.JWT.UserId {
			err = errors.New("工具[" + find.Name + "]不属于当前用户，无法操作")
			return
		}
	}
	if find == nil {
		find = this_.GetOtherToolbox(bindConfigRequest.ToolboxId)
	}
	option := ""
	if find != nil {
		option = find.Option
	}

	optionBytes := []byte(option)

	if len(optionBytes) > 0 {
		optionData := map[string]interface{}{}
		e := json.Unmarshal(optionBytes, &optionData)
		if e == nil {
			strV, strVOk := optionData["port"].(string)
			if strVOk {
				if strV == "" {
					delete(optionData, "port")
				} else {
					optionData["port"], _ = strconv.Atoi(strV)
				}
				optionBytes, _ = json.Marshal(optionData)
			}
		}
	}

	err = json.Unmarshal(optionBytes, &config)
	if err != nil {
		return
	}
	switch conf := config.(type) {
	case *db.DatabaseConfig:
		conf.Password = this_.DecryptOptionAttr(conf.Password)
		break
	case *redis.Config:
		if conf.CertPath != "" {
			conf.CertPath = this_.GetFilesFile(conf.CertPath)
		}
		conf.Auth = this_.DecryptOptionAttr(conf.Auth)
		break
	case *zookeeper.Config:
		conf.Password = this_.DecryptOptionAttr(conf.Password)
		break
	case *elasticsearch.Config:
		if conf.CertPath != "" {
			conf.CertPath = this_.GetFilesFile(conf.CertPath)
		}
		conf.Password = this_.DecryptOptionAttr(conf.Password)
		break
	case *kafka.Config:
		if conf.CertPath != "" {
			conf.CertPath = this_.GetFilesFile(conf.CertPath)
		}
		conf.Password = this_.DecryptOptionAttr(conf.Password)
		break
	}
	if err != nil {
		return
	}
	return
}

type ToolboxType struct {
	Name       string                `json:"name,omitempty"`
	Text       string                `json:"text,omitempty"`
	Icon       string                `json:"icon,omitempty"`
	Comment    string                `json:"comment,omitempty"`
	ConfigForm *form.Form            `json:"configForm,omitempty"`
	OtherForm  map[string]*form.Form `json:"otherForm,omitempty"`
}

var (
	toolboxTypes         = &[]*ToolboxType{}
	databaseWorker_      = databaseWorker()
	sshWorker_           = sshWorker()
	redisWorker_         = redisWorker()
	zookeeperWorker_     = zookeeperWorker()
	elasticsearchWorker_ = elasticsearchWorker()
	kafkaWorker_         = kafkaWorker()
	otherWorker_         = otherWorker()
)

func init() {
	*toolboxTypes = append(*toolboxTypes, databaseWorker_)
	*toolboxTypes = append(*toolboxTypes, sshWorker_)
	*toolboxTypes = append(*toolboxTypes, redisWorker_)
	*toolboxTypes = append(*toolboxTypes, zookeeperWorker_)
	*toolboxTypes = append(*toolboxTypes, elasticsearchWorker_)
	*toolboxTypes = append(*toolboxTypes, kafkaWorker_)
	*toolboxTypes = append(*toolboxTypes, otherWorker_)
}

func GetToolboxTypes() (res []*ToolboxType) {
	res = *toolboxTypes
	return
}

func GetToolboxType(name string) (res *ToolboxType) {
	for _, one := range *toolboxTypes {
		if one.Name == name {
			res = one
		}
	}
	return
}

func databaseWorker() *ToolboxType {

	worker_ := &ToolboxType{
		Name: "database",
		Text: "Database",
		ConfigForm: &form.Form{
			Fields: []*form.Field{
				{
					Label: "类型", Name: "type", Type: "select", DefaultValue: "mysql",
					Options: []*form.Option{
						{Text: "MySql", Value: "mysql"},
						{Text: "Sqlite", Value: "sqlite"},
						{Text: "达梦", Value: "dameng"},
						{Text: "金仓", Value: "kingbase"},
						{Text: "神通", Value: "shentong"},
						{Text: "Oracle", Value: "oracle"},
						{Text: "Postgresql", Value: "postgresql"},
						{Text: "Odbc", Value: "odbc"},
					},
					Rules: []*form.Rule{
						{Required: true, Message: "数据库类型不能为空"},
					},
				},
				{
					Label: "Host（127.0.0.1）", Name: "host", DefaultValue: "127.0.0.1", VIf: `type != 'sqlite' && type != 'odbc'`,
					Rules: []*form.Rule{
						{Required: true, Message: "数据库连接地址不能为空"},
					},
				},
				{
					Label: "Port（3306）", Name: "port", IsNumber: true, DefaultValue: 3306, VIf: `type != 'sqlite' && type != 'odbc'`,
					Rules: []*form.Rule{
						{Required: true, Message: "数据库连接端口不能为空"},
					},
				},
				{Label: "Username", Name: "username", VIf: `type != 'sqlite'`},
				{Label: "Password", Name: "password", Type: "password", VIf: `type != 'sqlite'`},
				{Label: "Database", Name: "database", VIf: `type == 'mysql'`},
				{Label: "SID", Name: "sid", VIf: `type == 'oracle'`,
					Rules: []*form.Rule{
						{Required: true, Message: "SID不能为空"},
					},
				},
				{Label: "DbName", Name: "dbName", VIf: `type == 'kingbase' || type == 'shentong' || type == 'postgresql'`,
					Rules: []*form.Rule{
						{Required: true, Message: "dbName径不能为空"},
					},
				},
				{Label: "数据库文件路径", Name: "databasePath", VIf: `type == 'sqlite'`,
					Rules: []*form.Rule{
						{Required: true, Message: "数据库文件路径不能为空"},
					},
				},
				{Label: "OdbcName", Name: "odbcName", VIf: `type == 'odbc'`,
					Rules: []*form.Rule{
						{Required: true, Message: "OdbcName不能为空"},
					},
				},
				{Label: "Odbc参数（key1=value1;key2=value2;）", Name: "odbcParams", VIf: `type == 'odbc'`},
				{Label: "OdbcDialectName", Name: "odbcDialectName", Type: "select", VIf: `type == 'odbc'`,
					Options: []*form.Option{
						{Text: "默认", Value: ""},
						{Text: "MySql", Value: "mysql"},
						{Text: "Sqlite", Value: "sqlite"},
						{Text: "达梦", Value: "dameng"},
						{Text: "金仓", Value: "kingbase"},
						{Text: "神通", Value: "shentong"},
						{Text: "Oracle", Value: "oracle"},
						{Text: "Postgresql", Value: "postgresql"},
					},
				},
			},
		},
	}

	return worker_
}

func elasticsearchWorker() *ToolboxType {
	worker_ := &ToolboxType{
		Name: "elasticsearch",
		Text: "Elasticsearch",
		ConfigForm: &form.Form{
			Fields: []*form.Field{
				{
					Label: "连接地址（http://127.0.0.1:9200）", Name: "url", DefaultValue: "http://127.0.0.1:9200",
					Rules: []*form.Rule{
						{Required: true, Message: "连接地址不能为空"},
					},
				},
				{Label: "用户名", Name: "username"},
				{Label: "密码", Name: "password"},
				{Label: "Cert", Name: "certPath", Type: "file", Placeholder: "请上传Cert"},
			},
		},
		OtherForm: map[string]*form.Form{
			"index": {
				Fields: []*form.Field{
					{
						Label: "IndexName（索引）", Name: "indexName", DefaultValue: "index_xxx",
						Rules: []*form.Rule{
							{Required: true, Message: "索引不能为空"},
						},
					},
					{
						Label: "结构", Name: "mapping", Type: "json", DefaultValue: map[string]interface{}{
							"settings": map[string]interface{}{
								"number_of_shards":   1,
								"number_of_replicas": 0,
							},
							"mappings": map[string]interface{}{
								"properties": map[string]interface{}{
									"title": map[string]interface{}{
										"type": "text",
									},
								},
							},
						},
						Rules: []*form.Rule{
							{Required: true, Message: "结构不能为空"},
						},
					},
				},
			},
		},
	}

	return worker_
}

func kafkaWorker() *ToolboxType {
	worker_ := &ToolboxType{
		Name: "kafka",
		Text: "Kafka",
		ConfigForm: &form.Form{
			Fields: []*form.Field{
				{Label: "连接地址（127.0.0.1:9092）", Name: "address", DefaultValue: "127.0.0.1:9092",
					Rules: []*form.Rule{
						{Required: true, Message: "连接地址不能为空"},
					},
				},
				{Label: "用户名", Name: "username"},
				{Label: "密码", Name: "password"},
				{Label: "Cert", Name: "certPath", Type: "file", Placeholder: "请上传Cert"},
			},
		},
		OtherForm: map[string]*form.Form{
			"topic": {
				Fields: []*form.Field{
					{
						Label: "Topic（主题）", Name: "topic", DefaultValue: "topic_xxx",
						Rules: []*form.Rule{
							{Required: true, Message: "主题不能为空"},
						},
					},
					{
						Label: "Partitions（分区）", Name: "numPartitions", DefaultValue: 1, IsNumber: true,
						Rules: []*form.Rule{
							{Required: true, Message: "分区不能为空"},
						},
					},
					{
						Label: "ReplicationFactor（分区副本）", Name: "replicationFactor", DefaultValue: 1, IsNumber: true,
						Rules: []*form.Rule{
							{Required: true, Message: "分区副本不能为空"},
						},
					},
				},
			},
			"push": {
				Fields: []*form.Field{
					{
						Label: "Topic（主题）", Name: "topic", DefaultValue: "topic_xxx",
						Rules: []*form.Rule{
							{Required: true, Message: "主题不能为空"},
						},
					},
					{
						Label: "Header", Name: "headers", Type: "list",
						Fields: []*form.Field{
							{Label: "Header Key", Name: "key"},
							{Label: "Header Value", Name: "value"},
						},
					},
					{
						Label: "KeyType", Name: "keyType", DefaultValue: "string", Type: "select",
						Options: []*form.Option{
							{Text: "String", Value: "string"},
							{Text: "Long（int64）", Value: "long"},
						},
						Rules: []*form.Rule{
							{Required: true, Message: "KeyType不能为空"},
						},
					},
					{
						Label: "Key", Name: "key",
					},
					{
						Label: "ValueType", Name: "valueType", DefaultValue: "string", Type: "select",
						Options: []*form.Option{
							{Text: "String", Value: "string"},
							{Text: "Long（int64）", Value: "long"},
						},
						Rules: []*form.Rule{
							{Required: true, Message: "ValueType不能为空"},
						},
					},
					{
						Label: "Value", Name: "value", Type: "textarea",
						Rules: []*form.Rule{
							{Required: true, Message: "Value不能为空"},
						},
					},
					{
						Label: "ValueJSON预览", Name: "valueView", BindName: "value", Type: "jsonView",
					},
				},
			},
		},
	}

	return worker_
}

func sshWorker() *ToolboxType {
	worker_ := &ToolboxType{
		Name: "ssh",
		Text: "SSH",
		ConfigForm: &form.Form{
			Fields: []*form.Field{
				{
					Label: "类型", Name: "type", Type: "select", DefaultValue: "tcp",
					Options: []*form.Option{
						{Text: "TCP", Value: "tcp"},
					},
					Rules: []*form.Rule{
						{Required: true, Message: "SSH类型不能为空"},
					},
				},
				{
					Label: "连接地址（127.0.0.1:22）", Name: "address", DefaultValue: "127.0.0.1:22",
					Rules: []*form.Rule{
						{Required: true, Message: "连接地址不能为空"},
					},
				},
				{Label: "Username", Name: "username"},
				{Label: "Password", Name: "password", Type: "password"},
				{Label: "PublicKey", Name: "publicKey", Type: "file", Placeholder: "请上传PublicKey文件"},
			},
		},
	}

	return worker_
}

func redisWorker() *ToolboxType {
	worker_ := &ToolboxType{
		Name: "redis",
		Text: "Redis",
		ConfigForm: &form.Form{
			Fields: []*form.Field{
				{Label: "连接地址（127.0.0.1:6379）", Name: "address", DefaultValue: "127.0.0.1:6379",
					Rules: []*form.Rule{
						{Required: true, Message: "连接地址不能为空"},
					},
				},
				{Label: "用户名", Name: "username"},
				{Label: "密码", Name: "auth", Type: "password"},
				{Label: "Cert", Name: "certPath", Type: "file", Placeholder: "请上传Cert"},
			},
		},
	}

	return worker_
}

func zookeeperWorker() *ToolboxType {
	worker_ := &ToolboxType{
		Name: "zookeeper",
		Text: "Zookeeper",
		ConfigForm: &form.Form{
			Fields: []*form.Field{
				{
					Label: "连接地址（127.0.0.1:2181）", Name: "address", DefaultValue: "127.0.0.1:2181",
					Rules: []*form.Rule{
						{Required: true, Message: "连接地址不能为空"},
					},
				},
				{Label: "Username", Name: "username"},
				{Label: "Password", Name: "password", Type: "password"},
			},
		},
	}

	return worker_
}

func otherWorker() *ToolboxType {
	worker_ := &ToolboxType{
		Name: "other",
		Text: "其它",
		ConfigForm: &form.Form{
			Fields: []*form.Field{},
		},
	}

	return worker_
}
