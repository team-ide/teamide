package modelers

type AppModel struct {
	ElementNode
	Db           *ConfigDbModel                 `json:"db,omitempty"`
	DbOther      map[string]*ConfigDbModel      `json:"dbOther,omitempty"`
	Redis        *ConfigRedisModel              `json:"redis,omitempty"`
	RedisOther   map[string]*ConfigRedisModel   `json:"redisOther,omitempty"`
	Zk           *ConfigZkModel                 `json:"zk,omitempty"`
	ZkOther      map[string]*ConfigZkModel      `json:"zkOther,omitempty"`
	Kafka        *ConfigKafkaModel              `json:"kafka,omitempty"`
	KafkaOther   map[string]*ConfigKafkaModel   `json:"kafkaOther,omitempty"`
	Es           *ConfigEsModel                 `json:"es,omitempty"`
	EsOther      map[string]*ConfigEsModel      `json:"esOther,omitempty"`
	Mongodb      *ConfigMongodbModel            `json:"mongodb,omitempty"`
	MongodbOther map[string]*ConfigMongodbModel `json:"mongodbOther,omitempty"`
	Other        map[string]any                 `json:"other,omitempty"`
	Text         string                         `json:"-"`
}

func init() {
	addDocTemplate(&docTemplate{
		Name:    TypeAppName,
		Comment: "应用",
		Fields:  []*docTemplateField{},
	})
}
