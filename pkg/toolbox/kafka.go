package toolbox

import (
	"encoding/binary"
	"encoding/json"
	"strconv"
	"strings"
	"teamide/pkg/form"

	"github.com/Shopify/sarama"
)

func init() {
	worker_ := &Worker{
		Name: "kafka",
		Text: "Kafka",
		Work: kafkaWork,
		ConfigForm: &form.Form{
			Fields: []*form.Field{
				{Label: "连接地址（127.0.0.1:9092）", Name: "url", DefaultValue: "127.0.0.1:9092",
					Rules: []*form.Rule{
						{Required: true, Message: "连接地址不能为空"},
					},
				},
			},
		},
	}

	AddWorker(worker_)
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

	Headers []KafkaMessageHeader `json:"headers"`
	Key     string               `json:"key"`
	Value   string               `json:"value"`
}

type KafkaMessageHeader struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type KafkaMessage struct {
	Key       interface{}          `json:"key"`
	Value     interface{}          `json:"value"`
	Topic     string               `json:"topic"`
	Partition int32                `json:"partition"`
	Offset    int64                `json:"offset"`
	Headers   []KafkaMessageHeader `json:"headers"`
}

type KafkaConfig struct {
	Address string `json:"address"`
}

func kafkaWork(work string, config map[string]interface{}, data map[string]interface{}) (res map[string]interface{}, err error) {

	var kafkaConfig KafkaConfig
	var bs []byte
	bs, err = json.Marshal(config)
	if err != nil {
		return
	}
	err = json.Unmarshal(bs, &kafkaConfig)
	if err != nil {
		return
	}

	var service *KafkaService
	service, err = getKafkaService(kafkaConfig)
	if err != nil {
		return
	}

	bs, err = json.Marshal(data)
	if err != nil {
		return
	}
	request := &KafkaBaseRequest{}
	err = json.Unmarshal(bs, request)
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
		kafkaMsgs, err := service.Pull(request.GroupId, []string{request.Topic}, request.PullSize, request.PullTimeout)
		if err != nil {
			return nil, err
		}
		msgs := []KafkaMessage{}
		for _, kafkaMsg := range kafkaMsgs {
			var key interface{}
			var value interface{}
			if strings.ToLower(request.KeyType) == "string" {
				key = sarama.StringEncoder(kafkaMsg.Key)
			} else if strings.ToLower(request.KeyType) == "long" {
				if len(kafkaMsg.Key) == 8 {
					key = uint64(binary.BigEndian.Uint64(kafkaMsg.Key))
				} else {
					key = sarama.StringEncoder(kafkaMsg.Key)
				}
			} else {
				key = sarama.ByteEncoder(kafkaMsg.Key)
			}
			if strings.ToLower(request.ValueType) == "string" {
				value = sarama.StringEncoder(kafkaMsg.Value)
			} else if strings.ToLower(request.ValueType) == "long" {
				if len(kafkaMsg.Value) == 8 {
					value = uint64(binary.BigEndian.Uint64(kafkaMsg.Value))
				} else {
					value = sarama.StringEncoder(kafkaMsg.Value)
				}
			} else {
				value = sarama.ByteEncoder(kafkaMsg.Value)
			}
			msg := KafkaMessage{
				Key:       key,
				Value:     value,
				Topic:     kafkaMsg.Topic,
				Partition: kafkaMsg.Partition,
				Offset:    kafkaMsg.Offset,
			}
			if kafkaMsg.Headers != nil {
				for _, header := range kafkaMsg.Headers {
					msg.Headers = append(msg.Headers, KafkaMessageHeader{Key: string(header.Key), Value: string(header.Value)})
				}
			}
			msgs = append(msgs, msg)
		}
		res["msgs"] = msgs
	case "push":
		var key sarama.Encoder
		var value sarama.Encoder
		if request.Key != "" {
			if strings.ToLower(request.KeyType) == "string" {
				key = sarama.StringEncoder(request.Key)
			} else if strings.ToLower(request.KeyType) == "long" {
				longV, err := strconv.ParseInt(request.Key, 10, 64)
				if err != nil {
					return nil, err
				}
				uintV := uint64(longV)
				bytes := make([]byte, 8)
				binary.BigEndian.PutUint64(bytes, uintV)
				key = sarama.ByteEncoder(bytes)
			} else {
				key = sarama.ByteEncoder(request.Key)
			}
		}
		if request.Value != "" {
			if strings.ToLower(request.ValueType) == "string" {
				value = sarama.StringEncoder(request.Value)
			} else if strings.ToLower(request.ValueType) == "long" {
				longV, err := strconv.ParseInt(request.Value, 10, 64)
				if err != nil {
					return nil, err
				}
				uintV := uint64(longV)
				bytes := make([]byte, 8)
				binary.BigEndian.PutUint64(bytes, uintV)
				value = sarama.ByteEncoder(bytes)
			} else {
				value = sarama.ByteEncoder(request.Value)
			}
		}

		kafkaMsg := &sarama.ProducerMessage{}
		kafkaMsg.Topic = request.Topic
		kafkaMsg.Key = key
		kafkaMsg.Value = value
		if request.Partition >= 0 {
			kafkaMsg.Partition = request.Partition
		}
		if request.Offset >= 0 {
			kafkaMsg.Offset = request.Offset
		}
		if request.Headers != nil {
			for _, one := range request.Headers {
				kafkaMsg.Headers = append(kafkaMsg.Headers, sarama.RecordHeader{
					Key:   []byte(one.Key),
					Value: []byte(one.Value),
				})
			}
		}

		err = service.Push(kafkaMsg)
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
