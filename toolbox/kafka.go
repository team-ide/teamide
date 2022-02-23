package toolbox

import (
	"context"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/Shopify/sarama"
)

func GetKafkaWorker() *Worker {
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

	return worker_
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

func getKafkaService(address string) (res *KafkaService, err error) {
	key := "kafka-" + address
	var service Service
	service, err = GetService(key, func() (res Service, err error) {
		var s *KafkaService
		s, err = CreateKafkaService(address)
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
	res = service.(*KafkaService)
	return
}

func CreateKafkaService(address string) (*KafkaService, error) {
	service := &KafkaService{
		address: address,
	}
	err := service.init()
	return service, err
}

//注册处理器在线信息等
type KafkaService struct {
	address     string
	lastUseTime int64
}

func (this_ *KafkaService) init() (err error) {
	return
}

func (this_ *KafkaService) GetWaitTime() int64 {
	return 10 * 60 * 1000
}

func (this_ *KafkaService) GetLastUseTime() int64 {
	return this_.lastUseTime
}

func (this_ *KafkaService) Stop() {

}

func (this_ *KafkaService) GetServers() []string {
	var servers []string
	if strings.Contains(this_.address, ",") {
		servers = strings.Split(this_.address, ",")
	} else if strings.Contains(this_.address, ";") {
		servers = strings.Split(this_.address, ";")
	} else {
		servers = []string{this_.address}
	}
	return servers
}

func (this_ *KafkaService) getClient() (saramaClient sarama.Client, err error) {
	SaramaConfig := sarama.NewConfig()
	SaramaConfig.Consumer.Return.Errors = true
	SaramaConfig.Consumer.Offsets.Initial = sarama.OffsetOldest
	SaramaConfig.Consumer.MaxWaitTime = time.Second * 1
	addrs := strings.Split(this_.address, ",")
	saramaClient, err = sarama.NewClient(addrs, SaramaConfig)
	if err != nil {
		return
	}
	return
}

func (service *KafkaService) GetTopics() (topics []string, err error) {
	var saramaClient sarama.Client
	saramaClient, err = service.getClient()
	if err != nil {
		return
	}
	defer saramaClient.Close()
	topics, err = saramaClient.Topics()
	return
}

func (this_ *KafkaService) Pull(groupId string, topics []string) (msgs []*sarama.ConsumerMessage, err error) {
	var saramaClient sarama.Client
	saramaClient, err = this_.getClient()
	if err != nil {
		return
	}
	defer saramaClient.Close()
	group, err := sarama.NewConsumerGroupFromClient(groupId, saramaClient)
	if err != nil {
		return
	}
	handler := &consumerGroupHandler{}
	go func() {
		ctx := context.Background()
		err = group.Consume(ctx, topics, handler)
		if err != nil {
			fmt.Println("group.Consume error:", err)
		}
	}()
	startTime := GetNowTime()
	for {
		time.Sleep(200 * time.Millisecond)
		nowTime := GetNowTime()
		if len(handler.msgs) > 0 && nowTime-startTime > 2*1000 {
			break
		}
		if nowTime-startTime > 5*1000 {
			break
		}
	}
	// go func() {
	err = group.Close()
	if err != nil {
		fmt.Println("group.Close error:", err)
		return
	}
	// }()
	msgs = handler.msgs
	return
}

type consumerGroupHandler struct {
	msgs []*sarama.ConsumerMessage
}

func (consumerGroupHandler) Setup(_ sarama.ConsumerGroupSession) error   { return nil }
func (consumerGroupHandler) Cleanup(_ sarama.ConsumerGroupSession) error { return nil }
func (handler *consumerGroupHandler) ConsumeClaim(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	chanMessages := claim.Messages()
	for msg := range chanMessages {
		handler.msgs = append(handler.msgs, msg)
	}
	return nil
}

func (this_ *KafkaService) MarkOffset(groupId string, topic string, partition int32, offset int64) (err error) {
	var saramaClient sarama.Client
	saramaClient, err = this_.getClient()
	if err != nil {
		return
	}
	defer saramaClient.Close()
	offsetManager, err := sarama.NewOffsetManagerFromClient(groupId, saramaClient)
	if err != nil {
		return
	}
	partitionOffsetManager, err := offsetManager.ManagePartition(topic, partition)
	if err != nil {
		return
	}
	partitionOffsetManager.MarkOffset(offset, "")
	err = offsetManager.Close()
	return
}

func (this_ *KafkaService) ResetOffset(groupId string, topic string, partition int32, offset int64) (err error) {
	var saramaClient sarama.Client
	saramaClient, err = this_.getClient()
	if err != nil {
		return
	}
	defer saramaClient.Close()
	offsetManager, err := sarama.NewOffsetManagerFromClient(groupId, saramaClient)
	if err != nil {
		return
	}
	partitionOffsetManager, err := offsetManager.ManagePartition(topic, partition)
	if err != nil {
		return
	}
	defer offsetManager.Close()
	partitionOffsetManager.ResetOffset(offset, "")
	return
}

func (this_ *KafkaService) CreatePartitions(topic string, count int32) (err error) {
	var saramaClient sarama.Client
	saramaClient, err = this_.getClient()
	if err != nil {
		return
	}
	defer saramaClient.Close()
	admin, err := sarama.NewClusterAdminFromClient(saramaClient)
	if err != nil {
		return
	}

	defer admin.Close()

	err = admin.CreatePartitions(topic, count, nil, false)

	return
}

func (this_ *KafkaService) DeleteTopic(topic string) (err error) {
	var saramaClient sarama.Client
	saramaClient, err = this_.getClient()
	if err != nil {
		return
	}
	defer saramaClient.Close()
	admin, err := sarama.NewClusterAdminFromClient(saramaClient)
	if err != nil {
		return
	}

	defer admin.Close()

	err = admin.DeleteTopic(topic)

	return
}

func (this_ *KafkaService) DeleteConsumerGroup(groupId string) (err error) {
	var saramaClient sarama.Client
	saramaClient, err = this_.getClient()
	if err != nil {
		return
	}
	defer saramaClient.Close()
	admin, err := sarama.NewClusterAdminFromClient(saramaClient)
	if err != nil {
		return
	}

	defer admin.Close()

	err = admin.DeleteConsumerGroup(groupId)

	return
}

func (this_ *KafkaService) DeleteRecords(topic string, partitionOffsets map[int32]int64) (err error) {
	var saramaClient sarama.Client
	saramaClient, err = this_.getClient()
	if err != nil {
		return
	}
	defer saramaClient.Close()
	admin, err := sarama.NewClusterAdminFromClient(saramaClient)
	if err != nil {
		return
	}

	defer admin.Close()

	err = admin.DeleteRecords(topic, partitionOffsets)

	return
}

//创建生产者
func (this_ *KafkaService) NewSyncProducer() (sarama.SyncProducer, error) {

	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	config.Producer.Timeout = 3
	var err error
	syncProducer, err := sarama.NewSyncProducer(this_.GetServers(), config)
	if err != nil {
		return nil, err
	}
	return syncProducer, nil
}

//推送消息到kafka
func (this_ *KafkaService) Push(msg *sarama.ProducerMessage) (err error) {
	syncProducer, err := this_.NewSyncProducer()
	if err != nil {
		return
	}
	defer syncProducer.Close()
	msg.Timestamp = time.Now()
	_, _, err = syncProducer.SendMessage(msg)
	return err
}
