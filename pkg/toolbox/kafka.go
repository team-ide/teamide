package toolbox

import (
	"encoding/binary"
	"encoding/json"
	"strconv"

	"github.com/Shopify/sarama"
)

func init() {
	worker_ := &Worker{
		Name:    "kafka",
		Text:    "Kafka",
		WorkMap: map[string]func(map[string]interface{}) (map[string]interface{}, error){},
	}
	worker_.WorkMap["topics"] = func(m map[string]interface{}) (map[string]interface{}, error) {
		return kafkaWork("topics", m["config"].(map[string]interface{}), m["data"].(map[string]interface{}))
	}
	worker_.WorkMap["pull"] = func(m map[string]interface{}) (map[string]interface{}, error) {
		return kafkaWork("pull", m["config"].(map[string]interface{}), m["data"].(map[string]interface{}))
	}
	worker_.WorkMap["push"] = func(m map[string]interface{}) (map[string]interface{}, error) {
		return kafkaWork("push", m["config"].(map[string]interface{}), m["data"].(map[string]interface{}))
	}
	worker_.WorkMap["commit"] = func(m map[string]interface{}) (map[string]interface{}, error) {
		return kafkaWork("commit", m["config"].(map[string]interface{}), m["data"].(map[string]interface{}))
	}
	worker_.WorkMap["reset"] = func(m map[string]interface{}) (map[string]interface{}, error) {
		return kafkaWork("reset", m["config"].(map[string]interface{}), m["data"].(map[string]interface{}))
	}
	worker_.WorkMap["deleteTopic"] = func(m map[string]interface{}) (map[string]interface{}, error) {
		return kafkaWork("deleteTopic", m["config"].(map[string]interface{}), m["data"].(map[string]interface{}))
	}
	worker_.WorkMap["createPartitions"] = func(m map[string]interface{}) (map[string]interface{}, error) {
		return kafkaWork("createPartitions", m["config"].(map[string]interface{}), m["data"].(map[string]interface{}))
	}
	worker_.WorkMap["deleteRecords"] = func(m map[string]interface{}) (map[string]interface{}, error) {
		return kafkaWork("deleteRecords", m["config"].(map[string]interface{}), m["data"].(map[string]interface{}))
	}

	AddWorker(worker_)
}

type KafkaBaseRequest struct {
	GroupId   string `json:"groupId"`
	Topic     string `json:"topic"`
	Partition int32  `json:"partition"`
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

func kafkaWork(work string, config map[string]interface{}, data map[string]interface{}) (res map[string]interface{}, err error) {
	var service *KafkaService
	service, err = getKafkaService(config["address"].(string))
	if err != nil {
		return
	}

	var bs []byte
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
		kafkaMsgs, err := service.Pull(request.GroupId, []string{request.Topic})
		if err != nil {
			return nil, err
		}
		msgs := []KafkaMessage{}
		for _, kafkaMsg := range kafkaMsgs {
			var key interface{}
			var value interface{}
			if request.KeyType == "String" {
				key = sarama.StringEncoder(kafkaMsg.Key)
			} else if request.KeyType == "Long" {
				if len(kafkaMsg.Key) == 8 {
					key = uint64(binary.BigEndian.Uint64(kafkaMsg.Key))
				} else {
					key = sarama.StringEncoder(kafkaMsg.Key)
				}
			} else {
				key = sarama.ByteEncoder(kafkaMsg.Key)
			}
			if request.ValueType == "String" {
				value = sarama.StringEncoder(kafkaMsg.Value)
			} else if request.ValueType == "Long" {
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
			if request.KeyType == "String" {
				key = sarama.StringEncoder(request.Key)
			} else if request.KeyType == "Long" {
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
			if request.ValueType == "String" {
				value = sarama.StringEncoder(request.Value)
			} else if request.ValueType == "Long" {
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
