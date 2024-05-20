package maker

import (
	"github.com/team-ide/go-tool/elasticsearch"
	"teamide/pkg/maker/modelers"
)

func NewEsCompiler(config *modelers.ConfigEsModel) *Component {
	component := &Component{}
	return component
}

func NewComponentEs(config *modelers.ConfigEsModel) (res *ComponentEs, err error) {
	ser, err := elasticsearch.New(&elasticsearch.Config{
		Url:      config.Url,
		Username: config.Username,
		Password: config.Password,
	})
	if err != nil {
		return
	}
	res = &ComponentEs{
		ser: ser,
	}
	return
}

type ComponentEs struct {
	ser elasticsearch.IService
}
