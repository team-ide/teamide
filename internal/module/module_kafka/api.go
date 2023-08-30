package module_kafka

import (
	"github.com/gin-gonic/gin"
	"github.com/team-ide/go-tool/kafka"
	"github.com/team-ide/go-tool/util"
	"go.uber.org/zap"
	"teamide/internal/module/module_toolbox"
	"teamide/pkg/base"
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
	test                  = base.AppendPower(&base.PowerAction{Action: "test", Text: "Kafka测试", ShouldLogin: true, StandAlone: true, Parent: Power})
	infoPower             = base.AppendPower(&base.PowerAction{Action: "info", Text: "Kafka信息", ShouldLogin: true, StandAlone: true, Parent: Power})
	topicsPower           = base.AppendPower(&base.PowerAction{Action: "topics", Text: "Kafka Topic查询", ShouldLogin: true, StandAlone: true, Parent: Power})
	topicPower            = base.AppendPower(&base.PowerAction{Action: "topic", Text: "Kafka Topic查询", ShouldLogin: true, StandAlone: true, Parent: Power})
	commitPower           = base.AppendPower(&base.PowerAction{Action: "commit", Text: "Kafka提交", ShouldLogin: true, StandAlone: true, Parent: Power})
	pullPower             = base.AppendPower(&base.PowerAction{Action: "pull", Text: "Kafka拉取", ShouldLogin: true, StandAlone: true, Parent: Power})
	pushPower             = base.AppendPower(&base.PowerAction{Action: "push", Text: "Kafka推送", ShouldLogin: true, StandAlone: true, Parent: Power})
	resetPower            = base.AppendPower(&base.PowerAction{Action: "reset", Text: "Kafka信息", ShouldLogin: true, StandAlone: true, Parent: Power})
	deleteTopicPower      = base.AppendPower(&base.PowerAction{Action: "deleteTopic", Text: "Kafka删除Topic", ShouldLogin: true, StandAlone: true, Parent: Power})
	createTopicPower      = base.AppendPower(&base.PowerAction{Action: "createTopic", Text: "Kafka创建Topic", ShouldLogin: true, StandAlone: true, Parent: Power})
	createPartitionsPower = base.AppendPower(&base.PowerAction{Action: "createPartitions", Text: "Kafka创建分区", ShouldLogin: true, StandAlone: true, Parent: Power})
	deleteRecordsPower    = base.AppendPower(&base.PowerAction{Action: "deleteRecords", Text: "Kafka删除记录", ShouldLogin: true, StandAlone: true, Parent: Power})
	topicDescribe         = base.AppendPower(&base.PowerAction{Action: "topicDescribe", Text: "Topic详情", ShouldLogin: true, StandAlone: true, Parent: Power})

	group              = base.AppendPower(&base.PowerAction{Action: "group", Text: "Kafka组", ShouldLogin: true, StandAlone: true, Parent: Power})
	groupList          = base.AppendPower(&base.PowerAction{Action: "list", Text: "组列表", ShouldLogin: true, StandAlone: true, Parent: group})
	groupDescribe      = base.AppendPower(&base.PowerAction{Action: "describe", Text: "组详情", ShouldLogin: true, StandAlone: true, Parent: group})
	groupOffsets       = base.AppendPower(&base.PowerAction{Action: "offsets", Text: "组Offsets", ShouldLogin: true, StandAlone: true, Parent: group})
	groupDeleteOffsets = base.AppendPower(&base.PowerAction{Action: "deleteOffsets", Text: "删除组Offsets", ShouldLogin: true, StandAlone: true, Parent: group})
	groupDelete        = base.AppendPower(&base.PowerAction{Action: "delete", Text: "删除组", ShouldLogin: true, StandAlone: true, Parent: group})

	closePower = base.AppendPower(&base.PowerAction{Action: "close", Text: "Kafka关闭", ShouldLogin: true, StandAlone: true, Parent: Power})
)

func (this_ *api) GetApis() (apis []*base.ApiWorker) {
	apis = append(apis, &base.ApiWorker{Power: test, Do: this_.test})
	apis = append(apis, &base.ApiWorker{Power: infoPower, Do: this_.info})
	apis = append(apis, &base.ApiWorker{Power: topicsPower, Do: this_.topics})
	apis = append(apis, &base.ApiWorker{Power: topicPower, Do: this_.topic})
	apis = append(apis, &base.ApiWorker{Power: commitPower, Do: this_.commit})
	apis = append(apis, &base.ApiWorker{Power: pullPower, Do: this_.pull})
	apis = append(apis, &base.ApiWorker{Power: pushPower, Do: this_.push})
	apis = append(apis, &base.ApiWorker{Power: resetPower, Do: this_.reset})
	apis = append(apis, &base.ApiWorker{Power: deleteTopicPower, Do: this_.deleteTopic})
	apis = append(apis, &base.ApiWorker{Power: createTopicPower, Do: this_.createTopic})
	apis = append(apis, &base.ApiWorker{Power: createPartitionsPower, Do: this_.createPartitions})
	apis = append(apis, &base.ApiWorker{Power: deleteRecordsPower, Do: this_.deleteRecords})
	apis = append(apis, &base.ApiWorker{Power: topicDescribe, Do: this_.topicDescribe})

	apis = append(apis, &base.ApiWorker{Power: groupList, Do: this_.groupList})
	apis = append(apis, &base.ApiWorker{Power: groupDescribe, Do: this_.groupDescribe})
	apis = append(apis, &base.ApiWorker{Power: groupOffsets, Do: this_.groupOffsets})
	apis = append(apis, &base.ApiWorker{Power: groupDeleteOffsets, Do: this_.groupDeleteOffsets})
	apis = append(apis, &base.ApiWorker{Power: groupDelete, Do: this_.groupDelete})

	apis = append(apis, &base.ApiWorker{Power: closePower, Do: this_.close})

	return
}

func (this_ *api) getConfig(requestBean *base.RequestBean, c *gin.Context) (config *kafka.Config, err error) {
	config = &kafka.Config{}
	_, err = this_.toolboxService.BindConfig(requestBean, c, config)
	if err != nil {
		return
	}
	return
}

func getService(kafkaConfig *kafka.Config) (res kafka.IService, err error) {
	key := "kafka-" + kafkaConfig.Address
	if kafkaConfig.Username != "" {
		key += "-" + base.GetMd5String(key+kafkaConfig.Username)
	}
	if kafkaConfig.Password != "" {
		key += "-" + base.GetMd5String(key+kafkaConfig.Password)
	}
	if kafkaConfig.CertPath != "" {
		key += "-" + base.GetMd5String(key+kafkaConfig.CertPath)
	}
	var serviceInfo *base.ServiceInfo
	serviceInfo, err = base.GetService(key, func() (res *base.ServiceInfo, err error) {
		var s kafka.IService
		s, err = kafka.New(kafkaConfig)
		if err != nil {
			util.Logger.Error("getKafkaService error", zap.Any("key", key), zap.Error(err))
			if s != nil {
				s.Close()
			}
			return
		}
		_, err = s.GetTopic("toolbox-kafka-test-topic", -2)
		if err != nil {
			util.Logger.Error("getKafkaService error", zap.Any("key", key), zap.Error(err))
			if s != nil {
				s.Close()
			}
			return
		}
		res = &base.ServiceInfo{
			WaitTime:    10 * 60 * 1000,
			LastUseTime: util.GetNowMilli(),
			Service:     s,
			Stop:        s.Close,
		}
		return
	})
	if err != nil {
		return
	}
	res = serviceInfo.Service.(kafka.IService)
	serviceInfo.SetLastUseTime()
	return
}

type BaseRequest struct {
	GroupId           string `json:"groupId"`
	Topic             string `json:"topic"`
	Time              int64  `json:"time"`
	PullSize          int    `json:"pullSize"`
	PullTimeout       int    `json:"pullTimeout"`
	Partition         int32  `json:"partition"`
	NumPartitions     int32  `json:"numPartitions"`
	ReplicationFactor int16  `json:"replicationFactor"`

	TopicPartitions map[string][]int32 `json:"topicPartitions"`

	Offset    int64  `json:"offset"`
	Count     int32  `json:"count"`
	KeyType   string `json:"keyType"`
	ValueType string `json:"valueType"`
}

func (this_ *api) test(requestBean *base.RequestBean, c *gin.Context) (res interface{}, err error) {
	config, err := this_.getConfig(requestBean, c)
	if err != nil {
		return
	}
	_, err = getService(config)
	if err != nil {
		return
	}

	return
}

func (this_ *api) info(requestBean *base.RequestBean, c *gin.Context) (res interface{}, err error) {
	config, err := this_.getConfig(requestBean, c)
	if err != nil {
		return
	}
	service, err := getService(config)
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
	service, err := getService(config)
	if err != nil {
		return
	}

	res, err = service.GetTopics()
	if err != nil {
		return
	}
	return
}

func (this_ *api) topic(requestBean *base.RequestBean, c *gin.Context) (res interface{}, err error) {
	config, err := this_.getConfig(requestBean, c)
	if err != nil {
		return
	}
	service, err := getService(config)
	if err != nil {
		return
	}

	request := &BaseRequest{}
	if !base.RequestJSON(request, c) {
		return
	}

	res, err = service.GetTopic(request.Topic, request.Time)
	if err != nil {
		return
	}
	return
}

func (this_ *api) topicDescribe(requestBean *base.RequestBean, c *gin.Context) (res interface{}, err error) {
	config, err := this_.getConfig(requestBean, c)
	if err != nil {
		return
	}
	service, err := getService(config)
	if err != nil {
		return
	}

	request := &BaseRequest{}
	if !base.RequestJSON(request, c) {
		return
	}

	list, err := service.DescribeTopics([]string{request.Topic})
	if err != nil {
		return
	}
	if len(list) > 0 {
		res = list[0]
	} else {
		res = nil
	}
	return
}

func (this_ *api) commit(requestBean *base.RequestBean, c *gin.Context) (res interface{}, err error) {
	config, err := this_.getConfig(requestBean, c)
	if err != nil {
		return
	}
	service, err := getService(config)
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
	service, err := getService(config)
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
	service, err := getService(config)
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
	service, err := getService(config)
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
	service, err := getService(config)
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
	service, err := getService(config)
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
	service, err := getService(config)
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
	service, err := getService(config)
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

func (this_ *api) groupList(requestBean *base.RequestBean, c *gin.Context) (res interface{}, err error) {
	config, err := this_.getConfig(requestBean, c)
	if err != nil {
		return
	}
	service, err := getService(config)
	if err != nil {
		return
	}

	request := &BaseRequest{}
	if !base.RequestJSON(request, c) {
		return
	}

	res, err = service.ListConsumerGroups()
	if err != nil {
		return
	}
	return
}

func (this_ *api) groupOffsets(requestBean *base.RequestBean, c *gin.Context) (res interface{}, err error) {
	config, err := this_.getConfig(requestBean, c)
	if err != nil {
		return
	}
	service, err := getService(config)
	if err != nil {
		return
	}

	request := &BaseRequest{}
	if !base.RequestJSON(request, c) {
		return
	}

	res, err = service.ListConsumerGroupOffsets(request.GroupId, request.TopicPartitions)
	if err != nil {
		return
	}
	return
}

func (this_ *api) groupDeleteOffsets(requestBean *base.RequestBean, c *gin.Context) (res interface{}, err error) {
	config, err := this_.getConfig(requestBean, c)
	if err != nil {
		return
	}
	service, err := getService(config)
	if err != nil {
		return
	}

	request := &BaseRequest{}
	if !base.RequestJSON(request, c) {
		return
	}

	err = service.DeleteConsumerGroupOffset(request.GroupId, request.Topic, request.Partition)
	if err != nil {
		return
	}
	return
}

func (this_ *api) groupDelete(requestBean *base.RequestBean, c *gin.Context) (res interface{}, err error) {
	config, err := this_.getConfig(requestBean, c)
	if err != nil {
		return
	}
	service, err := getService(config)
	if err != nil {
		return
	}

	request := &BaseRequest{}
	if !base.RequestJSON(request, c) {
		return
	}

	err = service.DeleteConsumerGroup(request.GroupId)
	if err != nil {
		return
	}
	return
}

func (this_ *api) groupDescribe(requestBean *base.RequestBean, c *gin.Context) (res interface{}, err error) {
	config, err := this_.getConfig(requestBean, c)
	if err != nil {
		return
	}
	service, err := getService(config)
	if err != nil {
		return
	}

	request := &BaseRequest{}
	if !base.RequestJSON(request, c) {
		return
	}

	list, err := service.DescribeConsumerGroups([]string{request.GroupId})
	if err != nil {
		return
	}
	if len(list) > 0 {
		res = list[0]
	} else {
		res = nil
	}
	return
}

func (this_ *api) close(requestBean *base.RequestBean, c *gin.Context) (res interface{}, err error) {

	return
}
