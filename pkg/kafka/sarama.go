package kafka

import (
	"context"
	"encoding/binary"
	"fmt"
	"github.com/Shopify/sarama"
	"sort"
	"strconv"
	"strings"
	"teamide/pkg/util"
	"time"
)

// Config kafka配置
type Config struct {
	Address string `json:"address"`
}

// CreateKafkaService 创建kafka服务
func CreateKafkaService(config Config) (*SaramaService, error) {
	service := &SaramaService{
		Address: config.Address,
	}
	err := service.Init()
	return service, err
}

// SaramaService 注册处理器在线信息等
type SaramaService struct {
	Address     string
	lastUseTime int64
}

func (this_ *SaramaService) Init() (err error) {
	return
}

func (this_ *SaramaService) GetWaitTime() int64 {
	return 10 * 60 * 1000
}

func (this_ *SaramaService) GetLastUseTime() int64 {
	return this_.lastUseTime
}

func (this_ *SaramaService) Stop() {

}

func (this_ *SaramaService) GetServers() []string {
	var servers []string
	if strings.Contains(this_.Address, ",") {
		servers = strings.Split(this_.Address, ",")
	} else if strings.Contains(this_.Address, ";") {
		servers = strings.Split(this_.Address, ";")
	} else {
		servers = []string{this_.Address}
	}
	return servers
}

func (this_ *SaramaService) getClient() (saramaClient sarama.Client, err error) {
	SaramaConfig := sarama.NewConfig()
	SaramaConfig.Consumer.Return.Errors = true
	SaramaConfig.Consumer.Offsets.Initial = sarama.OffsetOldest
	SaramaConfig.Consumer.MaxWaitTime = time.Second * 1
	adders := strings.Split(this_.Address, ",")
	saramaClient, err = sarama.NewClient(adders, SaramaConfig)
	if err != nil {
		_ = saramaClient.Close()
		return
	}
	return
}
func closeSaramaClient(saramaClient sarama.Client) {
	_ = saramaClient.Close()
}
func closeClusterAdmin(clusterAdmin sarama.ClusterAdmin) {
	_ = clusterAdmin.Close()
}

func (this_ *SaramaService) GetTopics() (topics []string, err error) {
	var saramaClient sarama.Client
	saramaClient, err = this_.getClient()
	if err != nil {
		return
	}
	defer closeSaramaClient(saramaClient)
	topics, err = saramaClient.Topics()
	if err != nil {
		return
	}

	sort.Strings(topics)
	return
}

func (this_ *SaramaService) Pull(groupId string, topics []string, PullSize int, PullTimeout int, keyType, valueType string) (msgList []*Message, err error) {
	if PullSize <= 0 {
		PullSize = 10
	}
	if PullTimeout <= 0 {
		PullTimeout = 1000
	}
	var saramaClient sarama.Client
	saramaClient, err = this_.getClient()
	if err != nil {
		return
	}
	defer closeSaramaClient(saramaClient)
	group, err := sarama.NewConsumerGroupFromClient(groupId, saramaClient)
	if err != nil {
		return
	}
	handler := &consumerGroupHandler{
		size: PullSize,
	}
	go func() {
		ctx := context.Background()
		err = group.Consume(ctx, topics, handler)

		if err != nil {
			fmt.Println("group.Consume error:", err)
		}
	}()
	startTime := util.GetNowTime()
	for {
		time.Sleep(100 * time.Millisecond)
		nowTime := util.GetNowTime()
		if handler.appended || nowTime-startTime >= int64(PullTimeout) {
			break
		}
	}
	err = group.Close()
	if err != nil {
		fmt.Println("group.Close error:", err)
		return
	}
	for _, one := range handler.messages {
		var msg *Message
		msg, err = ConsumerMessageToMessage(keyType, valueType, one)
		if err != nil {
			return
		}
		msgList = append(msgList, msg)
	}
	return
}

type consumerGroupHandler struct {
	messages []*sarama.ConsumerMessage
	appended bool
	size     int
}

func (consumerGroupHandler) Setup(_ sarama.ConsumerGroupSession) error   { return nil }
func (consumerGroupHandler) Cleanup(_ sarama.ConsumerGroupSession) error { return nil }
func (handler *consumerGroupHandler) ConsumeClaim(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	if sess == nil {
		return nil
	}
	chanMessages := claim.Messages()
	for msg := range chanMessages {
		handler.messages = append(handler.messages, msg)
		if len(handler.messages) >= handler.size {
			break
		}
	}
	handler.appended = true
	return nil
}

func (this_ *SaramaService) MarkOffset(groupId string, topic string, partition int32, offset int64) (err error) {
	var saramaClient sarama.Client
	saramaClient, err = this_.getClient()
	if err != nil {
		return
	}
	defer closeSaramaClient(saramaClient)
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

func (this_ *SaramaService) ResetOffset(groupId string, topic string, partition int32, offset int64) (err error) {
	var saramaClient sarama.Client
	saramaClient, err = this_.getClient()
	if err != nil {
		return
	}
	defer closeSaramaClient(saramaClient)
	offsetManager, err := sarama.NewOffsetManagerFromClient(groupId, saramaClient)
	if err != nil {
		return
	}
	partitionOffsetManager, err := offsetManager.ManagePartition(topic, partition)
	if err != nil {
		return
	}
	defer func() {
		_ = offsetManager.Close()
	}()
	partitionOffsetManager.ResetOffset(offset, "")
	return
}

func (this_ *SaramaService) CreatePartitions(topic string, count int32) (err error) {
	var saramaClient sarama.Client
	saramaClient, err = this_.getClient()
	if err != nil {
		return
	}
	defer closeSaramaClient(saramaClient)
	admin, err := sarama.NewClusterAdminFromClient(saramaClient)
	if err != nil {
		return
	}

	defer closeClusterAdmin(admin)

	err = admin.CreatePartitions(topic, count, nil, false)

	return
}

func (this_ *SaramaService) CreateTopic(topic string, numPartitions int32, replicationFactor int16) (err error) {
	var saramaClient sarama.Client
	saramaClient, err = this_.getClient()
	if err != nil {
		return
	}
	defer closeSaramaClient(saramaClient)
	admin, err := sarama.NewClusterAdminFromClient(saramaClient)
	if err != nil {
		return
	}

	defer closeClusterAdmin(admin)
	if numPartitions <= 0 {
		numPartitions = 1
	}
	if replicationFactor <= 0 {
		replicationFactor = 1
	}
	detail := &sarama.TopicDetail{
		NumPartitions:     numPartitions,
		ReplicationFactor: replicationFactor,
	}
	err = admin.CreateTopic(topic, detail, false)

	return
}

func (this_ *SaramaService) DeleteTopic(topic string) (err error) {
	var saramaClient sarama.Client
	saramaClient, err = this_.getClient()
	if err != nil {
		return
	}
	defer closeSaramaClient(saramaClient)
	admin, err := sarama.NewClusterAdminFromClient(saramaClient)
	if err != nil {
		return
	}

	defer closeClusterAdmin(admin)

	err = admin.DeleteTopic(topic)

	return
}

func (this_ *SaramaService) DeleteConsumerGroup(groupId string) (err error) {
	var saramaClient sarama.Client
	saramaClient, err = this_.getClient()
	if err != nil {
		return
	}
	defer closeSaramaClient(saramaClient)
	admin, err := sarama.NewClusterAdminFromClient(saramaClient)
	if err != nil {
		return
	}

	defer closeClusterAdmin(admin)

	err = admin.DeleteConsumerGroup(groupId)

	return
}

func (this_ *SaramaService) DeleteRecords(topic string, partitionOffsets map[int32]int64) (err error) {
	var saramaClient sarama.Client
	saramaClient, err = this_.getClient()
	if err != nil {
		return
	}
	defer closeSaramaClient(saramaClient)
	admin, err := sarama.NewClusterAdminFromClient(saramaClient)
	if err != nil {
		return
	}

	defer closeClusterAdmin(admin)

	err = admin.DeleteRecords(topic, partitionOffsets)

	return
}

//NewSyncProducer 创建生产者
func (this_ *SaramaService) NewSyncProducer() (sarama.SyncProducer, error) {

	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	config.Producer.Timeout = 3
	var err error
	syncProducer, err := sarama.NewSyncProducer(this_.GetServers(), config)
	if err != nil {
		_ = syncProducer.Close()
		return nil, err
	}
	return syncProducer, nil
}

type ProducerMessage struct {
	*sarama.ProducerMessage
}

// Push 推送消息到kafka
func (this_ *SaramaService) Push(msg *Message) (err error) {
	producerMessage, err := MessageToProducerMessage(msg)
	if err != nil {
		return
	}
	syncProducer, err := this_.NewSyncProducer()
	if err != nil {
		return
	}
	defer func() {
		_ = syncProducer.Close()
	}()

	_, _, err = syncProducer.SendMessage(producerMessage)
	return err
}

func MessageToProducerMessage(msg *Message) (producerMessage *sarama.ProducerMessage, err error) {
	var key sarama.Encoder
	var value sarama.Encoder
	if msg.Key != "" {
		if strings.ToLower(msg.KeyType) == "long" {
			longV, err := strconv.ParseInt(msg.Key, 10, 64)
			if err != nil {
				return nil, err
			}
			var bytes = make([]byte, 8)
			binary.BigEndian.PutUint64(bytes, uint64(longV))
			key = sarama.ByteEncoder(bytes)
		} else {
			key = sarama.ByteEncoder(msg.Key)
		}
	}
	if msg.Value != "" {
		if strings.ToLower(msg.ValueType) == "long" {
			longV, err := strconv.ParseInt(msg.Value, 10, 64)
			if err != nil {
				return nil, err
			}
			var bytes = make([]byte, 8)
			binary.BigEndian.PutUint64(bytes, uint64(longV))
			value = sarama.ByteEncoder(bytes)
		} else {
			value = sarama.ByteEncoder(msg.Value)
		}
	}

	producerMessage = &sarama.ProducerMessage{}
	producerMessage.Topic = msg.Topic
	producerMessage.Key = key
	producerMessage.Value = value
	if msg.Timestamp == nil || (*msg.Timestamp).IsZero() {
		producerMessage.Timestamp = time.Now()
	} else {
		producerMessage.Timestamp = *msg.Timestamp
	}
	if msg.Partition != nil {
		producerMessage.Partition = *msg.Partition
	}
	if msg.Offset != nil {
		producerMessage.Offset = *msg.Offset
	}
	if msg.Headers != nil {
		for _, one := range msg.Headers {
			producerMessage.Headers = append(producerMessage.Headers, sarama.RecordHeader{
				Key:   []byte(one.Key),
				Value: []byte(one.Value),
			})
		}
	}
	return
}

func ConsumerMessageToMessage(keyType string, valueType string, consumerMessage *sarama.ConsumerMessage) (msg *Message, err error) {
	var key string
	var value string

	if consumerMessage.Key != nil && len(consumerMessage.Key) > 0 {
		if len(consumerMessage.Key) == 8 {
			Uint64Key := binary.BigEndian.Uint64(consumerMessage.Key)
			int64Key := int64(Uint64Key)
			if int64Key >= 0 {
				key = strconv.FormatInt(int64Key, 10)
			}
		}
		if key == "" {
			key = string(consumerMessage.Key)
		}
	}
	if consumerMessage.Value != nil && len(consumerMessage.Value) > 0 {
		if len(consumerMessage.Value) == 8 {
			Uint64Value := binary.BigEndian.Uint64(consumerMessage.Value)
			int64Value := int64(Uint64Value)
			if int64Value >= 0 {
				value = strconv.FormatInt(int64Value, 10)
			}
		}
		if value == "" {
			value = string(consumerMessage.Value)
		}
	}
	msg = &Message{
		Key:       key,
		Value:     value,
		Topic:     consumerMessage.Topic,
		Partition: &consumerMessage.Partition,
		Offset:    &consumerMessage.Offset,
	}
	if consumerMessage.Headers != nil {
		for _, header := range consumerMessage.Headers {
			msg.Headers = append(msg.Headers, MessageHeader{Key: string(header.Key), Value: string(header.Value)})
		}
	}
	return
}

type MessageHeader struct {
	Key   string `json:"key,omitempty"`
	Value string `json:"value,omitempty"`
}

type Message struct {
	KeyType   string          `json:"keyType,omitempty"`
	Key       string          `json:"key,omitempty"`
	ValueType string          `json:"valueType,omitempty"`
	Value     string          `json:"value,omitempty"`
	Topic     string          `json:"topic,omitempty"`
	Partition *int32          `json:"partition,omitempty"`
	Offset    *int64          `json:"offset,omitempty"`
	Headers   []MessageHeader `json:"headers,omitempty"`
	Timestamp *time.Time      `json:"timestamp,omitempty"`
}
