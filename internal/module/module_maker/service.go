package module_maker

import (
	"errors"
	_ "teamide/pkg/maker/invokers"
	"teamide/pkg/maker/modelers"
)

type Config struct {
	Dir string `json:"dir"`
}

func createService(config *Config) (service *Service, err error) {
	service = &Service{
		Config: config,
	}
	err = service.init()
	return
}

type Service struct {
	*Config
	app       *modelers.Application
	isStopped bool
}

func (this_ *Service) init() (err error) {
	if this_.Dir == "" {
		err = errors.New("dir is empty")
		return
	}
	this_.app = modelers.Load(this_.Dir)
	return
}

func (this_ *Service) Close() {
	this_.isStopped = true
	return
}
