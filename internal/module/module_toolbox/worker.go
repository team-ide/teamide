package module_toolbox

import (
	"encoding/json"
	"errors"
	"teamide/pkg/db"
	"teamide/pkg/elasticsearch"
	"teamide/pkg/form"
	"teamide/pkg/kafka"
	"teamide/pkg/redis"
	"teamide/pkg/ssh"
	"teamide/pkg/toolbox"
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
		res = this_.Decryption.Encrypt(str)
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
	res = this_.Decryption.Decrypt(str)
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
	toolboxWorker := GetWorker(toolboxData.ToolboxType)
	if toolboxWorker == nil {
		err = errors.New("不支持的工具类型[" + toolboxData.ToolboxType + "]")
		return
	}
	optionBytes := []byte(toolboxData.Option)
	optionMap := map[string]interface{}{}
	err = json.Unmarshal(optionBytes, &optionMap)
	if err != nil {
		return
	}

	switch toolboxWorker {
	case databaseWorker_:
		str, ok := optionMap["password"].(string)
		if ok {
			optionMap["password"] = this_.EncryptOptionAttr(str)
		} else {
			delete(optionMap, "password")
		}
		break
	case redisWorker_:
		str, ok := optionMap["auth"].(string)
		if ok {
			optionMap["auth"] = this_.EncryptOptionAttr(str)
		} else {
			delete(optionMap, "auth")
		}
		break
	case zookeeperWorker_:
		break
	case elasticsearchWorker_:
		break
	case kafkaWorker_:
		break
	case otherWorker_:
		break
	case sshWorker_:
		str, ok := optionMap["password"].(string)
		if ok {
			optionMap["password"] = this_.EncryptOptionAttr(str)
		} else {
			delete(optionMap, "password")
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

// Work 执行
func (this_ *ToolboxService) Work(toolboxId int64, work string, data map[string]interface{}) (res interface{}, err error) {

	toolboxData, err := this_.Get(toolboxId)
	if err != nil {
		return
	}
	if toolboxData == nil {
		toolboxData = this_.GetOtherToolbox(toolboxId)
	}
	if toolboxData == nil {
		err = errors.New("工具配置丢失")
		return
	}

	toolboxWorker := GetWorker(toolboxData.ToolboxType)
	if toolboxWorker == nil {
		err = errors.New("不支持的工具类型[" + toolboxData.ToolboxType + "]")
		return
	}

	optionBytes := []byte(toolboxData.Option)

	switch toolboxWorker {
	case databaseWorker_:
		var config *db.DatabaseConfig
		err = json.Unmarshal(optionBytes, &config)
		if err != nil {
			return
		}
		config.Password = this_.DecryptOptionAttr(config.Password)
		res, err = toolbox.DatabaseWork(work, config, data)
		break
	case redisWorker_:
		var config *redis.Config
		err = json.Unmarshal(optionBytes, &config)
		if err != nil {
			return
		}
		config.Auth = this_.DecryptOptionAttr(config.Auth)
		res, err = toolbox.RedisWork(work, config, data)
		break
	case zookeeperWorker_:
		var config *zookeeper.Config
		err = json.Unmarshal(optionBytes, &config)
		if err != nil {
			return
		}
		res, err = toolbox.ZKWork(work, config, data)
		break
	case elasticsearchWorker_:
		var config *elasticsearch.Config
		err = json.Unmarshal(optionBytes, &config)
		if err != nil {
			return
		}
		if config.CertPath != "" {
			config.CertPath = this_.GetFilesFile(config.CertPath)
		}
		res, err = toolbox.ESWork(work, config, data)
		break
	case kafkaWorker_:
		var config *kafka.Config
		err = json.Unmarshal(optionBytes, &config)
		if err != nil {
			return
		}
		res, err = toolbox.KafkaWork(work, config, data)
		break
	case otherWorker_:
		var config *toolbox.OtherConfig
		err = json.Unmarshal(optionBytes, &config)
		if err != nil {
			return
		}
		res, err = toolbox.OtherWork(work, config, data)
		break
	case sshWorker_:
		var config *ssh.Config
		err = json.Unmarshal(optionBytes, &config)
		if err != nil {
			return
		}
		config.Password = this_.DecryptOptionAttr(config.Password)
		if config.PublicKey != "" {
			config.PublicKey = this_.GetFilesFile(config.PublicKey)
		}
		res, err = toolbox.SSHWork(work, config, data)
		break
	}
	if err != nil {
		return
	}

	return
}

type Worker struct {
	Name       string                `json:"name,omitempty"`
	Text       string                `json:"text,omitempty"`
	Icon       string                `json:"icon,omitempty"`
	Comment    string                `json:"comment,omitempty"`
	ConfigForm *form.Form            `json:"configForm,omitempty"`
	OtherForm  map[string]*form.Form `json:"otherForm,omitempty"`
}
type WorkerConfig interface {
	GetConfigKey() string
}

var (
	workers              = &[]*Worker{}
	databaseWorker_      = databaseWorker()
	sshWorker_           = sshWorker()
	redisWorker_         = redisWorker()
	zookeeperWorker_     = zookeeperWorker()
	elasticsearchWorker_ = elasticsearchWorker()
	kafkaWorker_         = kafkaWorker()
	otherWorker_         = otherWorker()
)

func init() {
	*workers = append(*workers, databaseWorker_)
	*workers = append(*workers, sshWorker_)
	*workers = append(*workers, redisWorker_)
	*workers = append(*workers, zookeeperWorker_)
	*workers = append(*workers, elasticsearchWorker_)
	*workers = append(*workers, kafkaWorker_)
	*workers = append(*workers, otherWorker_)
}

func GetWorkers() (res []*Worker) {
	res = *workers
	return
}

func GetWorker(name string) (res *Worker) {
	for _, one := range *workers {
		if one.Name == name {
			res = one
		}
	}
	return
}

func databaseWorker() *Worker {

	worker_ := &Worker{
		Name: "database",
		Text: "Database",
		ConfigForm: &form.Form{
			Fields: []*form.Field{
				{
					Label: "类型", Name: "type", Type: "select", DefaultValue: "mysql",
					Options: []*form.Option{
						{Text: "MySql", Value: "mysql"},
					},
					Rules: []*form.Rule{
						{Required: true, Message: "数据库类型不能为空"},
					},
				},
				{
					Label: "Host（127.0.0.1）", Name: "host", DefaultValue: "127.0.0.1",
					Rules: []*form.Rule{
						{Required: true, Message: "数据库连接地址不能为空"},
					},
				},
				{
					Label: "Port（3306）", Name: "port", IsNumber: true, DefaultValue: "3306",
					Rules: []*form.Rule{
						{Required: true, Message: "数据库连接端口不能为空"},
					},
				},
				{Label: "Username", Name: "username"},
				{Label: "Password", Name: "password", Type: "password"},
			},
		},
	}

	return worker_
}

func elasticsearchWorker() *Worker {
	worker_ := &Worker{
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

func kafkaWorker() *Worker {
	worker_ := &Worker{
		Name: "kafka",
		Text: "Kafka",
		ConfigForm: &form.Form{
			Fields: []*form.Field{
				{Label: "连接地址（127.0.0.1:9092）", Name: "address", DefaultValue: "127.0.0.1:9092",
					Rules: []*form.Rule{
						{Required: true, Message: "连接地址不能为空"},
					},
				},
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

func sshWorker() *Worker {
	worker_ := &Worker{
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

func redisWorker() *Worker {
	worker_ := &Worker{
		Name: "redis",
		Text: "Redis",
		ConfigForm: &form.Form{
			Fields: []*form.Field{
				{Label: "连接地址（127.0.0.1:6379）", Name: "address", DefaultValue: "127.0.0.1:6379",
					Rules: []*form.Rule{
						{Required: true, Message: "连接地址不能为空"},
					},
				},
				{Label: "密码", Name: "auth", Type: "password"},
			},
		},
	}

	return worker_
}

func zookeeperWorker() *Worker {
	worker_ := &Worker{
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
			},
		},
	}

	return worker_
}

func otherWorker() *Worker {
	worker_ := &Worker{
		Name: "other",
		Text: "其它",
		ConfigForm: &form.Form{
			Fields: []*form.Field{},
		},
	}

	return worker_
}
