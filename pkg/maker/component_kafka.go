package maker

import (
	"github.com/team-ide/go-tool/kafka"
	"teamide/pkg/maker/modelers"
)

func NewKafkaCompiler(config *modelers.ConfigKafkaModel) *Component {
	component := &Component{}
	return component
}

func NewComponentKafka(config *modelers.ConfigKafkaModel) (res *ComponentKafka, err error) {
	ser, err := kafka.New(&kafka.Config{
		Address:  config.Address,
		Username: config.Username,
		Password: config.Password,
	})
	if err != nil {
		return
	}
	res = &ComponentKafka{
		ser: ser,
	}
	return
}

type ComponentKafka struct {
	ser kafka.IService
}
