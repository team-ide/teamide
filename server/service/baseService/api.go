package baseService

import (
	"server/base"
	"server/component"

	"github.com/gin-gonic/gin"
)

func NewQueryPageApiWorker(
	Api string,
	Power *base.PowerAction,
	NewRequestBean func() interface{},
	NewBean func() interface{},
	GetSqlParam func(requestBean interface{}) (sqlParam base.SqlParam, err error),
) *base.ApiWorker {
	return ApiQueryPage{Api, Power, NewRequestBean, NewBean, GetSqlParam}.ApiWorker()
}

type ApiQueryPage struct {
	Api            string
	Power          *base.PowerAction
	NewRequestBean func() interface{}
	NewBean        func() interface{}
	GetSqlParam    func(requestBean interface{}) (sqlParam base.SqlParam, err error)
}

func (service ApiQueryPage) ApiWorker() (apiWorker *base.ApiWorker) {
	apiWorker = &base.ApiWorker{
		Apis:  []string{service.Api},
		Power: service.Power,
		Do:    service.Do,
	}
	return
}

func (service ApiQueryPage) Do(request *base.RequestBean, c *gin.Context) (res interface{}, err error) {
	var requestBean interface{}
	if service.NewRequestBean != nil {
		requestBean = service.NewRequestBean()
		if requestBean != nil {
			err = c.BindJSON(requestBean)
			if err != nil {
				return
			}
		}
	}
	var sqlParam base.SqlParam
	sqlParam, err = service.GetSqlParam(requestBean)
	if err != nil {
		return
	}
	res, err = component.DB.QueryPage(sqlParam, base.NewUserEntityInterface)
	if err != nil {
		return
	}
	return
}

func NewQueryListApiWorker(
	Api string,
	Power *base.PowerAction,
	NewRequestBean func() interface{},
	NewBean func() interface{},
	GetSqlParam func(requestBean interface{}) (sqlParam base.SqlParam, err error),
) *base.ApiWorker {
	return ApiQueryList{Api, Power, NewRequestBean, NewBean, GetSqlParam}.ApiWorker()
}

type ApiQueryList struct {
	Api            string
	Power          *base.PowerAction
	NewRequestBean func() interface{}
	NewBean        func() interface{}
	GetSqlParam    func(requestBean interface{}) (sqlParam base.SqlParam, err error)
}

func (service ApiQueryList) ApiWorker() (apiWorker *base.ApiWorker) {
	apiWorker = &base.ApiWorker{
		Apis:  []string{service.Api},
		Power: service.Power,
		Do:    service.Do,
	}
	return
}

func (service ApiQueryList) Do(request *base.RequestBean, c *gin.Context) (res interface{}, err error) {
	var requestBean interface{}
	if service.NewRequestBean != nil {
		requestBean = service.NewRequestBean()
		if requestBean != nil {
			err = c.BindJSON(requestBean)
			if err != nil {
				return
			}
		}
	}
	var sqlParam base.SqlParam
	sqlParam, err = service.GetSqlParam(requestBean)
	if err != nil {
		return
	}
	res, err = component.DB.Query(sqlParam, base.NewUserEntityInterface)
	if err != nil {
		return
	}
	return
}

func NewQueryOneApiWorker(
	Api string,
	Power *base.PowerAction,
	NewRequestBean func() interface{},
	NewBean func() interface{},
	GetSqlParam func(requestBean interface{}) (sqlParam base.SqlParam, err error),
) *base.ApiWorker {
	return ApiQueryOne{Api, Power, NewRequestBean, NewBean, GetSqlParam}.ApiWorker()
}

type ApiQueryOne struct {
	Api            string
	Power          *base.PowerAction
	NewRequestBean func() interface{}
	NewBean        func() interface{}
	GetSqlParam    func(requestBean interface{}) (sqlParam base.SqlParam, err error)
}

func (service ApiQueryOne) ApiWorker() (apiWorker *base.ApiWorker) {
	apiWorker = &base.ApiWorker{
		Apis:  []string{service.Api},
		Power: service.Power,
		Do:    service.Do,
	}
	return
}

func (service ApiQueryOne) Do(request *base.RequestBean, c *gin.Context) (res interface{}, err error) {
	var requestBean interface{}
	if service.NewRequestBean != nil {
		requestBean = service.NewRequestBean()
		if requestBean != nil {
			err = c.BindJSON(requestBean)
			if err != nil {
				return
			}
		}
	}
	var sqlParam base.SqlParam
	sqlParam, err = service.GetSqlParam(requestBean)
	if err != nil {
		return
	}
	var list []interface{}
	list, err = component.DB.Query(sqlParam, base.NewUserEntityInterface)
	if err != nil {
		return
	}
	if len(list) > 0 {
		res = list[0]
	}
	return
}

func NewExecApiWorker(
	Api string,
	Power *base.PowerAction,
	NewRequestBean func() interface{},
	GetSqlParam func(requestBean interface{}) (sqlParam base.SqlParam, err error),
) *base.ApiWorker {
	return ApiExec{Api, Power, NewRequestBean, GetSqlParam}.ApiWorker()
}

type ApiExec struct {
	Api            string
	Power          *base.PowerAction
	NewRequestBean func() interface{}
	GetSqlParam    func(requestBean interface{}) (sqlParam base.SqlParam, err error)
}

func (service ApiExec) ApiWorker() (apiWorker *base.ApiWorker) {
	apiWorker = &base.ApiWorker{
		Apis:  []string{service.Api},
		Power: service.Power,
		Do:    service.Do,
	}
	return
}

func (service ApiExec) Do(request *base.RequestBean, c *gin.Context) (res interface{}, err error) {
	var requestBean interface{}
	if service.NewRequestBean != nil {
		requestBean = service.NewRequestBean()
		if requestBean != nil {
			err = c.BindJSON(requestBean)
			if err != nil {
				return
			}
		}
	}
	var sqlParam base.SqlParam
	sqlParam, err = service.GetSqlParam(requestBean)
	if err != nil {
		return
	}
	res, err = component.DB.Exec(sqlParam)
	if err != nil {
		return
	}
	return
}

func NewInsertApiWorker(
	Api string,
	Power *base.PowerAction,
	NewRequestBean func() interface{},
	GetTableBean func(requestBean interface{}) (table string, bean interface{}, err error),
) *base.ApiWorker {
	return ApiInsert{Api, Power, NewRequestBean, GetTableBean}.ApiWorker()
}

type ApiInsert struct {
	Api            string
	Power          *base.PowerAction
	NewRequestBean func() interface{}
	GetTableBean   func(requestBean interface{}) (table string, bean interface{}, err error)
}

func (service ApiInsert) ApiWorker() (apiWorker *base.ApiWorker) {
	apiWorker = &base.ApiWorker{
		Apis:  []string{service.Api},
		Power: service.Power,
		Do:    service.Do,
	}
	return
}

func (service ApiInsert) Do(request *base.RequestBean, c *gin.Context) (res interface{}, err error) {
	var requestBean interface{}
	if service.NewRequestBean != nil {
		requestBean = service.NewRequestBean()
		if requestBean != nil {
			err = c.BindJSON(requestBean)
			if err != nil {
				return
			}
		}
	}
	var table string
	var bean interface{}
	table, bean, err = service.GetTableBean(requestBean)
	if err != nil {
		return
	}
	err = component.DB.InsertBean(table, bean)
	if err != nil {
		return
	}
	return
}

func NewUpdateApiWorker(
	Api string,
	Power *base.PowerAction,
	NewRequestBean func() interface{},
	GetTableBean func(requestBean interface{}) (table string, keys []string, bean interface{}, err error),
) *base.ApiWorker {
	return ApiUpdate{Api, Power, NewRequestBean, GetTableBean}.ApiWorker()
}

type ApiUpdate struct {
	Api            string
	Power          *base.PowerAction
	NewRequestBean func() interface{}
	GetTableBean   func(requestBean interface{}) (table string, keys []string, bean interface{}, err error)
}

func (service ApiUpdate) ApiWorker() (apiWorker *base.ApiWorker) {
	apiWorker = &base.ApiWorker{
		Apis:  []string{service.Api},
		Power: service.Power,
		Do:    service.Do,
	}
	return
}

func (service ApiUpdate) Do(request *base.RequestBean, c *gin.Context) (res interface{}, err error) {
	var requestBean interface{}
	if service.NewRequestBean != nil {
		requestBean = service.NewRequestBean()
		if requestBean != nil {
			err = c.BindJSON(requestBean)
			if err != nil {
				return
			}
		}
	}
	var table string
	var keys []string
	var bean interface{}
	table, keys, bean, err = service.GetTableBean(requestBean)
	if err != nil {
		return
	}
	err = component.DB.UpdateBean(table, keys, bean)
	if err != nil {
		return
	}
	return
}
