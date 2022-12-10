package module_kafka

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"teamide/internal/base"
	"teamide/internal/module/module_toolbox"
	"teamide/pkg/kafka"
	"teamide/pkg/util"
)

type api struct {
	toolboxService *module_toolbox.ToolboxService
}

func NewApi(toolboxService *module_toolbox.ToolboxService) *api {
	return &api{
		toolboxService: toolboxService,
	}
}

var (
	Power                 = base.AppendPower(&base.PowerAction{Action: "kafka", Text: "Kafka", ShouldLogin: true, StandAlone: true})
	infoPower             = base.AppendPower(&base.PowerAction{Action: "info", Text: "Kafka信息", ShouldLogin: true, StandAlone: true, Parent: Power})
	topicsPower           = base.AppendPower(&base.PowerAction{Action: "topics", Text: "Kafka信息", ShouldLogin: true, StandAlone: true, Parent: Power})
	commitPower           = base.AppendPower(&base.PowerAction{Action: "commit", Text: "Kafka信息", ShouldLogin: true, StandAlone: true, Parent: Power})
	pullPower             = base.AppendPower(&base.PowerAction{Action: "pull", Text: "Kafka信息", ShouldLogin: true, StandAlone: true, Parent: Power})
	pushPower             = base.AppendPower(&base.PowerAction{Action: "push", Text: "Kafka信息", ShouldLogin: true, StandAlone: true, Parent: Power})
	resetPower            = base.AppendPower(&base.PowerAction{Action: "reset", Text: "Kafka信息", ShouldLogin: true, StandAlone: true, Parent: Power})
	deleteTopicPower      = base.AppendPower(&base.PowerAction{Action: "deleteTopic", Text: "Kafka信息", ShouldLogin: true, StandAlone: true, Parent: Power})
	createTopicPower      = base.AppendPower(&base.PowerAction{Action: "createTopic", Text: "Kafka信息", ShouldLogin: true, StandAlone: true, Parent: Power})
	createPartitionsPower = base.AppendPower(&base.PowerAction{Action: "createPartitions", Text: "Kafka信息", ShouldLogin: true, StandAlone: true, Parent: Power})
	deleteRecordsPower    = base.AppendPower(&base.PowerAction{Action: "deleteRecords", Text: "Kafka信息", ShouldLogin: true, StandAlone: true, Parent: Power})
	closePower            = base.AppendPower(&base.PowerAction{Action: "close", Text: "Kafka信息", ShouldLogin: true, StandAlone: true, Parent: Power})
)

func (this_ *api) GetApis() (apis []*base.ApiWorker) {
	apis = append(apis, &base.ApiWorker{Apis: []string{"kafka/info"}, Power: infoPower, Do: this_.info})
	apis = append(apis, &base.ApiWorker{Apis: []string{"kafka/topics"}, Power: topicsPower, Do: this_.topics})
	apis = append(apis, &base.ApiWorker{Apis: []string{"kafka/commit"}, Power: commitPower, Do: this_.commit})
	apis = append(apis, &base.ApiWorker{Apis: []string{"kafka/pull"}, Power: pullPower, Do: this_.pull})
	apis = append(apis, &base.ApiWorker{Apis: []string{"kafka/push"}, Power: pushPower, Do: this_.push})
	apis = append(apis, &base.ApiWorker{Apis: []string{"kafka/reset"}, Power: resetPower, Do: this_.reset})
	apis = append(apis, &base.ApiWorker{Apis: []string{"kafka/deleteTopic"}, Power: deleteTopicPower, Do: this_.deleteTopic})
	apis = append(apis, &base.ApiWorker{Apis: []string{"kafka/createTopic"}, Power: createTopicPower, Do: this_.createTopic})
	apis = append(apis, &base.ApiWorker{Apis: []string{"kafka/createPartitions"}, Power: createPartitionsPower, Do: this_.createPartitions})
	apis = append(apis, &base.ApiWorker{Apis: []string{"kafka/deleteRecords"}, Power: deleteRecordsPower, Do: this_.deleteRecords})
	apis = append(apis, &base.ApiWorker{Apis: []string{"kafka/close"}, Power: closePower, Do: this_.close})

	return
}

func (this_ *api) getConfig(requestBean *base.RequestBean, c *gin.Context) (config *kafka.Config, err error) {
	config = &kafka.Config{}
	err = this_.toolboxService.BindConfig(requestBean, c, config)
	if err != nil {
		return
	}
	return
}

func getService(kafkaConfig kafka.Config) (res *kafka.SaramaService, err error) {
	key := "kafka-" + kafkaConfig.Address
	if kafkaConfig.Username != "" {
		key += "-" + util.GetMd5String(key+kafkaConfig.Username)
	}
	if kafkaConfig.Password != "" {
		key += "-" + util.GetMd5String(key+kafkaConfig.Password)
	}
	if kafkaConfig.CertPath != "" {
		key += "-" + util.GetMd5String(key+kafkaConfig.CertPath)
	}
	var service util.Service
	service, err = util.GetService(key, func() (res util.Service, err error) {
		var s *kafka.SaramaService
		s, err = kafka.CreateKafkaService(kafkaConfig)
		if err != nil {
			util.Logger.Error("getKafkaService error", zap.Any("key", key), zap.Error(err))
			if s != nil {
				s.Stop()
			}
			return
		}
		_, err = s.GetTopics()
		if err != nil {
			util.Logger.Error("getKafkaService error", zap.Any("key", key), zap.Error(err))
			if s != nil {
				s.Stop()
			}
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

type BaseRequest struct {
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

func (this_ *api) info(requestBean *base.RequestBean, c *gin.Context) (res interface{}, err error) {
	config, err := this_.getConfig(requestBean, c)
	if err != nil {
		return
	}
	service, err := getService(*config)
	if err != nil {
		return
	}

	res, err = service.Info()
	if err != nil {
		return
	}
	return
}

func (this_ *api) topics(requestBean *base.RequestBean, c *gin.Context) (res interface{}, err error) {
	config, err := this_.getConfig(requestBean, c)
	if err != nil {
		return
	}
	service, err := getService(*config)
	if err != nil {
		return
	}

	res, err = service.GetTopics()
	if err != nil {
		return
	}
	return
}

func (this_ *api) commit(requestBean *base.RequestBean, c *gin.Context) (res interface{}, err error) {
	config, err := this_.getConfig(requestBean, c)
	if err != nil {
		return
	}
	service, err := getService(*config)
	if err != nil {
		return
	}

	request := &BaseRequest{}
	if !base.RequestJSON(request, c) {
		return
	}

	err = service.MarkOffset(request.GroupId, request.Topic, request.Partition, request.Offset)
	if err != nil {
		return
	}
	return
}

func (this_ *api) pull(requestBean *base.RequestBean, c *gin.Context) (res interface{}, err error) {
	config, err := this_.getConfig(requestBean, c)
	if err != nil {
		return
	}
	service, err := getService(*config)
	if err != nil {
		return
	}

	request := &BaseRequest{}
	if !base.RequestJSON(request, c) {
		return
	}

	res, err = service.Pull(request.GroupId, []string{request.Topic}, request.PullSize, request.PullTimeout, request.KeyType, request.ValueType)
	if err != nil {
		return
	}
	return
}

func (this_ *api) push(requestBean *base.RequestBean, c *gin.Context) (res interface{}, err error) {
	config, err := this_.getConfig(requestBean, c)
	if err != nil {
		return
	}
	service, err := getService(*config)
	if err != nil {
		return
	}

	request := &kafka.Message{}
	if !base.RequestJSON(request, c) {
		return
	}

	err = service.Push(request)
	if err != nil {
		return nil, err
	}
	return
}

func (this_ *api) reset(requestBean *base.RequestBean, c *gin.Context) (res interface{}, err error) {
	config, err := this_.getConfig(requestBean, c)
	if err != nil {
		return
	}
	service, err := getService(*config)
	if err != nil {
		return
	}

	request := &BaseRequest{}
	if !base.RequestJSON(request, c) {
		return
	}

	err = service.ResetOffset(request.GroupId, request.Topic, request.Partition, request.Offset)
	if err != nil {
		return
	}
	return
}

func (this_ *api) deleteTopic(requestBean *base.RequestBean, c *gin.Context) (res interface{}, err error) {
	config, err := this_.getConfig(requestBean, c)
	if err != nil {
		return
	}
	service, err := getService(*config)
	if err != nil {
		return
	}

	request := &BaseRequest{}
	if !base.RequestJSON(request, c) {
		return
	}

	err = service.DeleteTopic(request.Topic)
	if err != nil {
		return
	}
	return
}

func (this_ *api) createTopic(requestBean *base.RequestBean, c *gin.Context) (res interface{}, err error) {
	config, err := this_.getConfig(requestBean, c)
	if err != nil {
		return
	}
	service, err := getService(*config)
	if err != nil {
		return
	}

	request := &BaseRequest{}
	if !base.RequestJSON(request, c) {
		return
	}

	err = service.CreateTopic(request.Topic, request.NumPartitions, request.ReplicationFactor)
	if err != nil {
		return
	}
	return
}

func (this_ *api) createPartitions(requestBean *base.RequestBean, c *gin.Context) (res interface{}, err error) {
	config, err := this_.getConfig(requestBean, c)
	if err != nil {
		return
	}
	service, err := getService(*config)
	if err != nil {
		return
	}

	request := &BaseRequest{}
	if !base.RequestJSON(request, c) {
		return
	}

	err = service.CreatePartitions(request.Topic, request.Count)
	if err != nil {
		return
	}
	return
}

func (this_ *api) deleteRecords(requestBean *base.RequestBean, c *gin.Context) (res interface{}, err error) {
	config, err := this_.getConfig(requestBean, c)
	if err != nil {
		return
	}
	service, err := getService(*config)
	if err != nil {
		return
	}

	request := &BaseRequest{}
	if !base.RequestJSON(request, c) {
		return
	}

	partitionOffsets := make(map[int32]int64)
	partitionOffsets[request.Partition] = request.Offset
	err = service.DeleteRecords(request.Topic, partitionOffsets)
	if err != nil {
		return
	}
	return
}

func (this_ *api) close(requestBean *base.RequestBean, c *gin.Context) (res interface{}, err error) {

	return
}
