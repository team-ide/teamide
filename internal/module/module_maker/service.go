package module_maker

import (
	"errors"
	"teamide/pkg/maker"
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
	app       *maker.Application
	isStopped bool
}

func (this_ *Service) init() (err error) {
	if this_.Dir == "" {
		err = errors.New("dir is empty")
		return
	}
	this_.app = maker.Load(this_.Dir)
	return
}

func (this_ *Service) Close() {
	this_.isStopped = true
	return
}
