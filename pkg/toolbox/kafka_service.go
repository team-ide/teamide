package toolbox

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/Shopify/sarama"
)

func getKafkaService(kafkaConfig KafkaConfig) (res *KafkaService, err error) {
	key := "kafka-" + kafkaConfig.Address
	var service Service
	service, err = GetService(key, func() (res Service, err error) {
		var s *KafkaService
		s, err = CreateKafkaService(kafkaConfig)
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

func CreateKafkaService(kafkaConfig KafkaConfig) (*KafkaService, error) {
	service := &KafkaService{
		address: kafkaConfig.Address,
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
	adders := strings.Split(this_.address, ",")
	saramaClient, err = sarama.NewClient(adders, SaramaConfig)
	if err != nil {
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

func (this_ *KafkaService) GetTopics() (topics []string, err error) {
	var saramaClient sarama.Client
	saramaClient, err = this_.getClient()
	if err != nil {
		return
	}
	defer closeSaramaClient(saramaClient)
	topics, err = saramaClient.Topics()
	return
}
func (this_ *KafkaService) Pull(groupId string, topics []string, PullSize int, PullTimeout int) (msgs []*sarama.ConsumerMessage, err error) {
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
	startTime := GetNowTime()
	for {
		time.Sleep(100 * time.Millisecond)
		nowTime := GetNowTime()
		if handler.appended || nowTime-startTime >= int64(PullTimeout) {
			break
		}
	}
	err = group.Close()
	if err != nil {
		fmt.Println("group.Close error:", err)
		return
	}
	msgs = handler.messages
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

func (this_ *KafkaService) MarkOffset(groupId string, topic string, partition int32, offset int64) (err error) {
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

func (this_ *KafkaService) ResetOffset(groupId string, topic string, partition int32, offset int64) (err error) {
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
	defer closeSaramaClient(saramaClient)
	admin, err := sarama.NewClusterAdminFromClient(saramaClient)
	if err != nil {
		return
	}

	defer closeClusterAdmin(admin)

	err = admin.CreatePartitions(topic, count, nil, false)

	return
}

func (this_ *KafkaService) CreateTopic(topic string, numPartitions int32, replicationFactor int16) (err error) {
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

func (this_ *KafkaService) DeleteTopic(topic string) (err error) {
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

func (this_ *KafkaService) DeleteConsumerGroup(groupId string) (err error) {
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

func (this_ *KafkaService) DeleteRecords(topic string, partitionOffsets map[int32]int64) (err error) {
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
