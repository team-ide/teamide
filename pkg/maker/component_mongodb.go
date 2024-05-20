package maker

import (
	"github.com/team-ide/go-tool/mongodb"
	"teamide/pkg/maker/modelers"
)

func NewMongodbCompiler(config *modelers.ConfigMongodbModel) *Component {
	component := &Component{}
	return component
}

func NewComponentMongodb(config *modelers.ConfigMongodbModel) (res *ComponentMongodb, err error) {
	ser, err := mongodb.New(&mongodb.Config{
		Address:  config.Address,
		Username: config.Username,
		Password: config.Password,
	})
	if err != nil {
		return
	}
	res = &ComponentMongodb{
		ser: ser,
	}
	return
}

type ComponentMongodb struct {
	ser mongodb.IService
}
