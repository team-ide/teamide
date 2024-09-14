package module_toolbox

import (
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/team-ide/go-tool/db"
	"github.com/team-ide/go-tool/elasticsearch"
	"github.com/team-ide/go-tool/kafka"
	"github.com/team-ide/go-tool/mongodb"
	"github.com/team-ide/go-tool/redis"
	"github.com/team-ide/go-tool/util"
	"github.com/team-ide/go-tool/zookeeper"
	"strconv"
	"sync"
	"teamide/pkg/base"
	"teamide/pkg/form"
	"teamide/pkg/maker"
	"teamide/pkg/ssh"
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
func (this_ *ToolboxService) FormatOption(toolboxData *ToolboxModel, decrypt bool) (err error) {
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
	// 使用JSONDecodeUseNumber 防止精度丢失
	err = util.JSONDecodeUseNumber(optionBytes, &optionMap)
	if err != nil {
		return
	}
	for k, v := range optionMap {
		bV, isBool := v.(bool)
		sV, isString := v.(string)
		if v == nil ||
			(isString && sV == "") ||
			(isBool && !bV) {
			delete(optionMap, k)
		}
	}
	switch toolboxType {
	case databaseWorker_:
		if optionMap["password"] != nil {
			str, ok := optionMap["password"].(string)
			if ok {
				if decrypt {
					optionMap["password"] = this_.DecryptOptionAttr(str)
				} else {
					optionMap["password"] = this_.EncryptOptionAttr(str)
				}
			} else {
				delete(optionMap, "password")
			}
		}
		break
	case redisWorker_:
		if optionMap["auth"] != nil {
			str, ok := optionMap["auth"].(string)
			if ok {
				if decrypt {
					optionMap["auth"] = this_.DecryptOptionAttr(str)
				} else {
					optionMap["auth"] = this_.EncryptOptionAttr(str)
				}
			} else {
				delete(optionMap, "auth")
			}
		}
		break
	case zookeeperWorker_:
		if optionMap["password"] != nil {
			str, ok := optionMap["password"].(string)
			if ok {
				if decrypt {
					optionMap["password"] = this_.DecryptOptionAttr(str)
				} else {
					optionMap["password"] = this_.EncryptOptionAttr(str)
				}
			} else {
				delete(optionMap, "password")
			}
		}
		break
	case elasticsearchWorker_:
		if optionMap["password"] != nil {
			str, ok := optionMap["password"].(string)
			if ok {
				if decrypt {
					optionMap["password"] = this_.DecryptOptionAttr(str)
				} else {
					optionMap["password"] = this_.EncryptOptionAttr(str)
				}
			} else {
				delete(optionMap, "password")
			}
		}
		break
	case kafkaWorker_:
		if optionMap["password"] != nil {
			str, ok := optionMap["password"].(string)
			if ok {
				if decrypt {
					optionMap["password"] = this_.DecryptOptionAttr(str)
				} else {
					optionMap["password"] = this_.EncryptOptionAttr(str)
				}
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
				if decrypt {
					optionMap["password"] = this_.DecryptOptionAttr(str)
				} else {
					optionMap["password"] = this_.EncryptOptionAttr(str)
				}
			} else {
				delete(optionMap, "password")
			}
		}
		break
	case mongodbWorker_:
		if optionMap["password"] != nil {
			str, ok := optionMap["password"].(string)
			if ok {
				if decrypt {
					optionMap["password"] = this_.DecryptOptionAttr(str)
				} else {
					optionMap["password"] = this_.EncryptOptionAttr(str)
				}
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

func (this_ *ToolboxService) GetSSHConfig(option string) (config *ssh.Config, sshConfig *ssh.Config, err error) {
	config = &ssh.Config{}
	sshConfig, err = this_.BindConfigByOption(option, config, nil)
	return
}

type BindConfigRequest struct {
	ToolboxToTest string `json:"toolboxToTest,omitempty"`
	ToolboxId     int64  `json:"toolboxId,omitempty"`
	ToolboxType   string `json:"toolboxType,omitempty"`
}

func (this_ *ToolboxService) CheckToolboxPower(requestBean *base.RequestBean, toolboxModel *ToolboxModel) (err error) {

	if toolboxModel.UserId != 0 {
		// 如果是 开放的 则都可以操作
		if toolboxModel.Visibility == visibilityOpen {

		} else {
			// 验证创建者是否是当前用户
			if requestBean.JWT == nil || toolboxModel.UserId != requestBean.JWT.UserId {
				err = errors.New("工具[" + toolboxModel.Name + "]不属于当前用户，无法操作")
				return
			}
		}

	}
	return
}

func (this_ *ToolboxService) CheckPower(requestBean *base.RequestBean) (err error) {
	if v := requestBean.GetExtend("toolboxModel"); v != nil {
		find := v.(*ToolboxModel)
		err = this_.CheckToolboxPower(requestBean, find)
		if err != nil {
			return
		}
	}
	return
}

func (this_ *ToolboxService) initExtent(requestBean *base.RequestBean, c *gin.Context) (err error) {
	if requestBean.GetExtend("toolboxModel") != nil {
		return
	}
	bindConfigRequest := &BindConfigRequest{}
	if !base.RequestJSON(bindConfigRequest, c) {
		return
	}
	if bindConfigRequest.ToolboxToTest == "1" {
		toolboxModel := &ToolboxModel{}
		if !base.RequestJSON(toolboxModel, c) {
			return
		}
		if toolboxModel.Option == "" {
			err = errors.New("toolbox info is null")
			return
		}
		requestBean.SetExtend("toolboxModel", toolboxModel)
		return
	}

	find, err := this_.Get(bindConfigRequest.ToolboxId)
	if err != nil {
		return
	}
	if find == nil {
		find = this_.GetOtherToolbox(bindConfigRequest.ToolboxId)
	}
	if find != nil {
		requestBean.SetExtend("toolboxModel", find)
	}
	return
}

func (this_ *ToolboxService) BindConfigById(toolboxId int64, config interface{}) (sshConfig *ssh.Config, err error) {
	find, err := this_.Get(toolboxId)
	if err != nil {
		return
	}
	var option string
	if find != nil {
		option = find.Option
	}
	sshConfig, err = this_.BindConfigByOption(option, config, nil)
	return
}

func (this_ *ToolboxService) BindConfigByOption(option string, config interface{}, sshOptions []string) (sshConfig *ssh.Config, err error) {

	sshConfig = nil

	optionBytes := []byte(option)

	if len(optionBytes) > 0 {
		optionData := map[string]interface{}{}
		// 使用JSONDecodeUseNumber 防止精度丢失
		e := util.JSONDecodeUseNumber(optionBytes, &optionData)
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
		if optionData["sshToolboxId"] != nil {
			s := util.GetStringValue(optionData["sshToolboxId"])
			if s != "" {
				sshToolboxId, _ := strconv.ParseInt(s, 10, 64)
				if sshToolboxId > 0 {
					var sshToolbox *ToolboxModel
					sshToolbox, err = this_.Get(sshToolboxId)
					if err != nil {
						err = errors.New("ssh toolbox get error:" + err.Error())
						return
					}
					if sshToolbox != nil {
						sshConfig, _, err = this_.GetSSHConfig(sshToolbox.Option)
						if err != nil {
							err = errors.New("ssh toolbox config error:" + err.Error())
							return
						}
						//this_.Logger.Info("BindConfig find sshConfig", zap.Any("sshConfig", sshConfig))
					}
				}
			}
		}
	}

	err = json.Unmarshal(optionBytes, &config)
	if err != nil {
		return
	}
	switch conf := config.(type) {
	case *db.Config:
		if conf.TlsRootCert != "" {
			conf.TlsRootCert = this_.GetFilesFile(conf.TlsRootCert)
		}
		if conf.TlsClientCert != "" {
			conf.TlsClientCert = this_.GetFilesFile(conf.TlsClientCert)
		}
		if conf.TlsClientKey != "" {
			conf.TlsClientKey = this_.GetFilesFile(conf.TlsClientKey)
		}
		conf.Password = this_.DecryptOptionAttr(conf.Password)
		break
	case *ssh.Config:
		if conf.PublicKey != "" {
			conf.PublicKey = this_.GetFilesFile(conf.PublicKey)
		}
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
	case *mongodb.Config:
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
func (this_ *ToolboxService) BindConfig(requestBean *base.RequestBean, c *gin.Context, config interface{}) (sshConfig *ssh.Config, err error) {
	err = this_.initExtent(requestBean, c)
	if err != nil {
		return
	}

	err = this_.CheckPower(requestBean)
	if err != nil {
		return
	}
	option := ""
	if v := requestBean.GetExtend("toolboxModel"); v != nil {
		find := v.(*ToolboxModel)
		option = find.Option
	}

	sshConfig, err = this_.BindConfigByOption(option, config, nil)

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
	toolboxTypesLock     = &sync.Mutex{}
	toolboxTypesInit     bool
	databaseWorker_      = databaseWorker()
	sshWorker_           = sshWorker()
	redisWorker_         = redisWorker()
	zookeeperWorker_     = zookeeperWorker()
	elasticsearchWorker_ = elasticsearchWorker()
	kafkaWorker_         = kafkaWorker()
	mongodbWorker_       = mongodbWorker()
	netConnWorker_       = netConn()

	thriftWorker_ = thriftWorker()
	makerWorker_  = makerWorker()

	otherWorker_ = otherWorker()
)

func init() {
	initToolboxTypes()
}
func initToolboxTypes() {
	if toolboxTypesInit {
		return
	}
	toolboxTypesInit = true
	*toolboxTypes = append(*toolboxTypes, databaseWorker_)
	*toolboxTypes = append(*toolboxTypes, sshWorker_)
	*toolboxTypes = append(*toolboxTypes, redisWorker_)
	*toolboxTypes = append(*toolboxTypes, zookeeperWorker_)
	*toolboxTypes = append(*toolboxTypes, elasticsearchWorker_)
	*toolboxTypes = append(*toolboxTypes, kafkaWorker_)
	*toolboxTypes = append(*toolboxTypes, mongodbWorker_)
	*toolboxTypes = append(*toolboxTypes, netConnWorker_)
	*toolboxTypes = append(*toolboxTypes, thriftWorker_)
	if maker.HasMaker {
		*toolboxTypes = append(*toolboxTypes, makerWorker_)
	}
	//*toolboxTypes = append(*toolboxTypes, otherWorker_)
}

func GetToolboxTypes() (res []*ToolboxType) {
	toolboxTypesLock.Lock()
	defer toolboxTypesLock.Unlock()
	initToolboxTypes()
	res = *toolboxTypes
	return
}

func GetToolboxType(name string) (res *ToolboxType) {
	for _, one := range *toolboxTypes {
		if one.Name == name {
			res = one
			return
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
						{Text: "OpenGauss", Value: "opengauss"},
						{Text: "GBase", Value: "gbase"},
						{Text: "Odbc", Value: "odbc"},
					},
					Rules: []*form.Rule{
						{Required: true, Message: "数据库类型不能为空"},
					},
				},
				{
					Label: "SSH隧道", Name: "sshToolboxId", Type: "select", VIf: `type == 'mysql' || type == 'kingbase' || type == 'postgresql' || type == 'opengauss'`,
					OptionsName: "sshToolboxOptions",
					Rules:       []*form.Rule{},
				},
				{
					Label: "Host（127.0.0.1）", Name: "host", DefaultValue: "127.0.0.1", VIf: `type != 'sqlite' && type != 'odbc' && type != 'gbase'`,
					Rules: []*form.Rule{
						{Required: true, Message: "数据库连接地址不能为空"},
					},
					Col: 12,
				},
				{
					Label: "Port（3306）", Name: "port", IsNumber: true, DefaultValue: 3306, VIf: `type != 'sqlite' && type != 'odbc' && type != 'gbase'`,
					Rules: []*form.Rule{
						{Required: true, Message: "数据库连接端口不能为空"},
					},
					Col: 12,
				},
				{Label: "Username", Name: "username", VIf: `type != 'sqlite' && type != 'odbc' && type != 'gbase'`, Col: 12},
				{Label: "Password", Name: "password", Type: "password", VIf: `type != 'sqlite' && type != 'odbc' && type != 'gbase'`, Col: 12, ShowPlaintextBtn: true},
				{Label: "Database", Name: "database", VIf: `type == 'mysql'`},
				{Label: "SID", Name: "sid", VIf: `type == 'oracle'`,
					Rules: []*form.Rule{
						{Required: true, Message: "SID不能为空"},
					},
				},
				{Label: "DbName", Name: "dbName", VIf: `type == 'kingbase' || type == 'shentong' || type == 'postgresql' || type == 'opengauss'`,
					Rules: []*form.Rule{
						{Message: "dbName径不能为空"},
					},
				},
				{Label: "数据库文件路径", Name: "databasePath", VIf: `type == 'sqlite'`,
					Rules: []*form.Rule{
						{Required: true, Message: "数据库文件路径不能为空"},
					},
				},
				{Label: "OdbcDsn", Name: "odbcDsn", VIf: `type == 'odbc' || type == 'gbase'`,
					Rules: []*form.Rule{
						{Required: true, Message: "OdbcDsn不能为空"},
					},
				},
				{Label: "OdbcDialectName", Name: "odbcDialectName", Type: "select", VIf: `type == 'odbc'`,
					Options: []*form.Option{
						{Text: "默认", Value: ""},
						{Text: "MySql", Value: "mysql"},
						{Text: "Sqlite", Value: "sqlite"},
						{Text: "达梦", Value: "dameng"},
						{Text: "金仓", Value: "kingbase"},
						{Text: "神通", Value: "shentong"},
						{Text: "Oracle", Value: "oracle"},
						{Text: "GBase", Value: "gbase"},
						{Text: "Postgresql", Value: "postgresql"},
						{Text: "OpenGauss", Value: "opengauss"},
					},
				},
				{Label: "TLS", Name: "tlsConfig", Type: "select", VIf: `type == 'mysql'`, DefaultValue: "", Placeholder: "不配置",
					Options: []*form.Option{
						{Text: "不配置", Value: ""},
						{Text: "Skip Verify", Value: "skip-verify"},
						{Text: "Preferred", Value: "preferred"},
						{Text: "自定义", Value: "custom"},
					},
				},
				{Label: "TLS RootCert", Name: "tlsRootCert", Type: "file", VIf: `type == 'mysql' && tlsConfig == 'custom'`},
				{Label: "TLS Client Cert", Name: "tlsClientCert", Type: "file", VIf: `type == 'mysql' && tlsConfig == 'custom'`},
				{Label: "TLS Client Key", Name: "tlsClientKey", Type: "file", VIf: `type == 'mysql' && tlsConfig == 'custom'`},
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
				{Label: "用户名", Name: "username", Col: 12},
				{Label: "密码", Name: "password", Col: 12, ShowPlaintextBtn: true},
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
				{Label: "用户名", Name: "username", Col: 12},
				{Label: "密码", Name: "password", Col: 12, ShowPlaintextBtn: true},
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
			"commit": {
				Fields: []*form.Field{
					{
						Label: "Group（消费组）", Name: "groupId", DefaultValue: "test-group",
						Rules: []*form.Rule{
							{Required: true, Message: "Group不能为空"},
						},
					},
					{
						Label: "Topic（主题）", Name: "topic", DefaultValue: "topic_xxx",
						Rules: []*form.Rule{
							{Required: true, Message: "主题不能为空"},
						},
					},
					{
						Label: "Partition（分区）", Name: "partition", DefaultValue: 0, IsNumber: true,
						Rules: []*form.Rule{
							{Required: true, Message: "分区不能为空"},
						},
					},
					{
						Label: "Offset", Name: "offset", DefaultValue: 0, IsNumber: true,
						Rules: []*form.Rule{
							{Required: true, Message: "Offset不能为空"},
						},
					},
				},
			},
			"deleteRecord": {
				Fields: []*form.Field{
					{
						Label: "Topic（主题）", Name: "topic", DefaultValue: "topic_xxx",
						Rules: []*form.Rule{
							{Required: true, Message: "主题不能为空"},
						},
					},
					{
						Label: "Partition（分区）", Name: "partition", DefaultValue: 0, IsNumber: true,
						Rules: []*form.Rule{
							{Required: true, Message: "分区不能为空"},
						},
					},
					{
						Label: "Offset", Name: "offset", DefaultValue: 0, IsNumber: true,
						Rules: []*form.Rule{
							{Required: true, Message: "Offset不能为空"},
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
					Label: "SSH隧道", Name: "sshToolboxId", Type: "select",
					OptionsName: "sshToolboxOptions",
					Rules:       []*form.Rule{},
					Col:         12,
				},
				{
					Label: "连接地址（127.0.0.1:22）", Name: "address", DefaultValue: "127.0.0.1:22",
					Rules: []*form.Rule{
						{Required: true, Message: "连接地址不能为空"},
					},
					Col: 12,
				},
				{Label: "Username", Name: "username", Col: 9},
				{Label: "Password（密码或密钥文件密码）", Name: "password", Type: "password", Col: 9, ShowPlaintextBtn: true},

				{Label: `连接超时时间（秒）`, Name: "timeout", IsNumber: true, Col: 6, DefaultValue: 5},

				{Label: "空闲自动发送（防止会话超时）", Name: "idleSendOpen", Type: "switch", Col: 8, DefaultValue: false},
				{Label: `发送间隔（秒）`, Name: "idleSendTime", IsNumber: true, Col: 8, DefaultValue: 60, VIf: "idleSendOpen == true"},
				{Label: `发送字符（^C：Ctrl+C、\n：回车）`, Name: "idleSendChar", Col: 8, DefaultValue: "^C", VIf: "idleSendOpen == true"},

				{Label: "PrivateKey（通常跳板机需要的密钥文件）", Name: "publicKey", Type: "file", Placeholder: "请上传PrivateKey文件"},
				{Label: "连接后执行命令(回车执行多条，sleep 5，表示等待5秒执行下一条)", Name: "command", Type: "textarea"},
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
				{
					Label: "SSH隧道", Name: "sshToolboxId", Type: "select",
					OptionsName: "sshToolboxOptions",
					Rules:       []*form.Rule{},
					Col:         12,
				},
				{Label: "连接地址（127.0.0.1:6379）", Name: "address", DefaultValue: "127.0.0.1:6379",
					Rules: []*form.Rule{
						{Required: true, Message: "连接地址不能为空"},
					},
					Col: 12,
				},
				{Label: "用户名", Name: "username", Col: 12},
				{Label: "密码", Name: "auth", Type: "password", Col: 12, ShowPlaintextBtn: true},
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
					Label: "SSH隧道", Name: "sshToolboxId", Type: "select",
					OptionsName: "sshToolboxOptions",
					Rules:       []*form.Rule{},
				},
				{
					Label: "连接地址（127.0.0.1:2181）", Name: "address", DefaultValue: "127.0.0.1:2181",
					Rules: []*form.Rule{
						{Required: true, Message: "连接地址不能为空"},
					},
				},
				{Label: "Username", Name: "username", Col: 12},
				{Label: "Password", Name: "password", Type: "password", Col: 12, ShowPlaintextBtn: true},
			},
		},
	}

	return worker_
}

func mongodbWorker() *ToolboxType {
	worker_ := &ToolboxType{
		Name: "mongodb",
		Text: "Mongodb",
		ConfigForm: &form.Form{
			Fields: []*form.Field{
				{
					Label: "连接地址（127.0.0.1:27017）", Name: "address", DefaultValue: "127.0.0.1:27017",
					Rules: []*form.Rule{
						{Required: true, Message: "连接地址不能为空"},
					},
				},
				{Label: "Username", Name: "username", Col: 12},
				{Label: "Password", Name: "password", Type: "password", Col: 12, ShowPlaintextBtn: true},
				{Label: "Cert", Name: "certPath", Type: "file", Placeholder: "请上传Cert"},
			},
		},
	}

	return worker_
}

func netConn() *ToolboxType {
	worker_ := &ToolboxType{
		Name: "connection",
		Text: "网络（TCP）",
		ConfigForm: &form.Form{
			Fields: []*form.Field{
				{
					Label: "SSH隧道", Name: "sshToolboxId", Type: "select",
					OptionsName: "sshToolboxOptions",
					Rules:       []*form.Rule{},
					Col:         12,
				},
				{
					Label: "连接地址（127.0.0.1:6379）", Name: "address", DefaultValue: "127.0.0.1:6379",
					Rules: []*form.Rule{
						{Required: true, Message: "连接地址不能为空"},
					},
					Col: 12,
				},
				{Label: "超时时间（毫秒）", Name: "timeout", Col: 12, IsNumber: true, DefaultValue: 5000},
			},
		},
	}

	return worker_
}
func thriftWorker() *ToolboxType {
	worker_ := &ToolboxType{
		Name: "thrift",
		Text: "Thrift",
		ConfigForm: &form.Form{
			Fields: []*form.Field{
				{
					Label: "Thrift文件目录", Name: "thriftDir",
					Rules: []*form.Rule{
						{Required: true, Message: "Thrift文件目录不能为空"},
					},
				},
			},
		},
	}

	return worker_
}

func makerWorker() *ToolboxType {
	worker_ := &ToolboxType{
		Name: "maker",
		Text: "模型转代码",
		ConfigForm: &form.Form{
			Fields: []*form.Field{
				{
					Label: "目录（请配置一个目录）", Name: "dir",
					Rules: []*form.Rule{
						{Required: true, Message: "目录不能为空"},
					},
				},
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
