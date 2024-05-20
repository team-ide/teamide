package maker

import (
	"github.com/team-ide/go-tool/zookeeper"
	"teamide/pkg/maker/modelers"
)

func NewZkCompiler(config *modelers.ConfigZkModel) *Component {
	component := &Component{}
	return component
}

func NewComponentZk(config *modelers.ConfigZkModel) (res *ComponentZk, err error) {
	ser, err := zookeeper.New(&zookeeper.Config{
		Address:  config.Address,
		Username: config.Username,
		Password: config.Password,
	})
	if err != nil {
		return
	}
	res = &ComponentZk{
		ser: ser,
	}
	return
}

type ComponentZk struct {
	ser zookeeper.IService
}
