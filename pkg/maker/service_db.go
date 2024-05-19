package maker

import (
	"errors"
	"fmt"
	"github.com/team-ide/go-tool/db"
	"github.com/team-ide/go-tool/util"
	"go.uber.org/zap"
	"teamide/pkg/maker/modelers"
)

func (this_ *Invoker) GetDbService() (res *ServiceDb, err error) {
	defer func() {
		if e := recover(); e != nil {
			err = errors.New(fmt.Sprint(e))
			util.Logger.Error("get db service error", zap.Any("error", err))
		}
	}()
	res, err = this_.GetDbServiceByName("")
	return
}

func (this_ *Invoker) GetDbServiceByName(name string) (res *ServiceDb, err error) {
	if name == "" {
		name = "default"
	}
	defer func() {
		if e := recover(); e != nil {
			err = errors.New(fmt.Sprint(e))
			util.Logger.Error("get db service by name error", zap.Any("name", name), zap.Any("error", err))
		}
	}()
	this_.dbServiceCacheLock.Lock()
	defer this_.dbServiceCacheLock.Unlock()
	res, find := this_.dbServiceCache[name]

	if find {
		return
	}

	util.Logger.Info("db service not found,now create service", zap.Any("name", name))

	var config *modelers.ConfigDbModel
	config = this_.GetConfigDb(name)
	if config == nil {
		err = errors.New("config db [" + name + "] is not exist")
		util.Logger.Error("create db service error", zap.Any("error", err))
		return
	}
	ser, err := db.New(&db.Config{
		Type:     config.Type,
		Host:     config.Host,
		Port:     config.Port,
		Username: config.Username,
		Password: config.Password,
		Database: config.Database,
	})
	if err != nil {
		util.Logger.Error("create db service error", zap.Any("error", err))
		return
	}
	res = &ServiceDb{
		ser: ser,
	}
	this_.dbServiceCache[name] = res

	var scriptVar = "db"
	if name != "default" {
		scriptVar = "db_" + name
	}
	err = this_.setScriptVar(scriptVar, res)

	return
}

type ServiceDb struct {
	ser db.IService
}

func (this_ *ServiceDb) ShouldMappingFunc() bool {
	return true
}

func (this_ *ServiceDb) SelectOne(args ...interface{}) (res any, err error) {
	find := map[string]interface{}{
		"userId": 1,
		"name":   "张三",
	}
	res = find
	return
}
