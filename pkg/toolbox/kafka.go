package toolbox

import (
	"encoding/json"
	"teamide/pkg/form"
	"teamide/pkg/kafka"
)

func kafkaWorker() *Worker {
	worker_ := &Worker{
		Name: "kafka",
		Text: "Kafka",
		Work: kafkaWork,
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

func getKafkaService(kafkaConfig kafka.Config) (res *kafka.SaramaService, err error) {
	key := "kafka-" + kafkaConfig.Address
	var service Service
	service, err = GetService(key, func() (res Service, err error) {
		var s *kafka.SaramaService
		s, err = kafka.CreateKafkaService(kafkaConfig)
		if err != nil {
			return
		}
		_, err = s.GetTopics()
		if err != nil {
			return
		}
		res = s
		return
	})
	if err != nil {
		return
	}
	res = service.(*kafka.SaramaService)
	return
}

type KafkaBaseRequest struct {
	GroupId           string `json:"groupId"`
	Topic             string `json:"topic"`
	PullSize          int    `json:"pullSize"`
	PullTimeout       int    `json:"pullTimeout"`
	Partition         int32  `json:"partition"`
	NumPartitions     int32  `json:"numPartitions"`
	ReplicationFactor int16  `json:"replicationFactor"`

	Offset    int64  `json:"offset"`
	Count     int32  `json:"count"`
	KeyType   string `json:"keyType"`
	ValueType string `json:"valueType"`
}

func kafkaWork(work string, config map[string]interface{}, data map[string]interface{}) (res map[string]interface{}, err error) {

	var kafkaConfig kafka.Config
	var configBS []byte
	configBS, err = json.Marshal(config)
	if err != nil {
		return
	}
	err = json.Unmarshal(configBS, &kafkaConfig)
	if err != nil {
		return
	}

	var service *kafka.SaramaService
	service, err = getKafkaService(kafkaConfig)
	if err != nil {
		return
	}

	dataBS, err := json.Marshal(data)
	if err != nil {
		return
	}
	request := &KafkaBaseRequest{}
	err = json.Unmarshal(dataBS, request)
	if err != nil {
		return
	}

	res = map[string]interface{}{}

	switch work {
	case "topics":
		var topics []string
		topics, err = service.GetTopics()
		if err != nil {
			return
		}
		res["topics"] = topics
	case "commit":
		err = service.MarkOffset(request.GroupId, request.Topic, request.Partition, request.Offset)
		if err != nil {
			return
		}
	case "pull":
		var msgList []*kafka.Message
		msgList, err = service.Pull(request.GroupId, []string{request.Topic}, request.PullSize, request.PullTimeout, request.KeyType, request.ValueType)
		if err != nil {
			return
		}
		res["msgList"] = msgList
	case "push":

		msg := &kafka.Message{}
		err = json.Unmarshal(dataBS, msg)
		if err != nil {
			return
		}
		err = service.Push(msg)
		if err != nil {
			return nil, err
		}
	case "reset":
		err = service.ResetOffset(request.GroupId, request.Topic, request.Partition, request.Offset)
		if err != nil {
			return
		}
	case "deleteTopic":
		err = service.DeleteTopic(request.Topic)
		if err != nil {
			return
		}
	case "createTopic":
		err = service.CreateTopic(request.Topic, request.NumPartitions, request.ReplicationFactor)
		if err != nil {
			return
		}
	case "createPartitions":
		err = service.CreatePartitions(request.Topic, request.Count)
		if err != nil {
			return
		}
	case "deleteRecords":
		partitionOffsets := make(map[int32]int64)
		partitionOffsets[request.Partition] = request.Offset
		err = service.DeleteRecords(request.Topic, partitionOffsets)
		if err != nil {
			return
		}
	}
	return
}
