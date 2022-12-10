package module_log

import (
	"github.com/gin-gonic/gin"
	"teamide/internal/base"
	"teamide/internal/context"
)

type Api struct {
	*context.ServerContext
	LogService *LogService
}

func NewApi(LogService *LogService) *Api {
	return &Api{
		ServerContext: LogService.ServerContext,
		LogService:    LogService,
	}
}

var (
	// 用户 权限

	// Power 用户基本 权限
	Power          = base.AppendPower(&base.PowerAction{Action: "log", Text: "日志", ShouldLogin: true, StandAlone: true})
	queryPagePower = base.AppendPower(&base.PowerAction{Action: "queryPage", Text: "用户日志查询", Parent: Power, ShouldLogin: true, StandAlone: true})
	cleanPower     = base.AppendPower(&base.PowerAction{Action: "clean", Text: "清理用户日志", Parent: Power, ShouldLogin: true, StandAlone: true})
)

func (this_ *Api) GetApis() (apis []*base.ApiWorker) {
	apis = append(apis, &base.ApiWorker{Power: queryPagePower, Do: this_.queryPage})
	apis = append(apis, &base.ApiWorker{Power: cleanPower, Do: this_.clean})

	return
}

type QueryPageRequest struct {
	*LogPage
	Action string `json:"action"`
}

type QueryPageResponse struct {
	*LogPage
}

func (this_ *Api) queryPage(requestBean *base.RequestBean, c *gin.Context) (res interface{}, err error) {

	request := &QueryPageRequest{}
	if !base.RequestJSON(request, c) {
		return
	}
	response := &QueryPageResponse{}

	log := &LogModel{
		UserId: requestBean.JWT.UserId,
		Action: request.Action,
	}
	err = this_.LogService.QueryPage(log, request.LogPage)
	if err != nil {
		return
	}
	response.LogPage = request.LogPage
	res = response
	return
}

func (this_ *Api) clean(requestBean *base.RequestBean, c *gin.Context) (res interface{}, err error) {

	log := &LogModel{
		UserId: requestBean.JWT.UserId,
	}
	err = this_.LogService.clean(log)
	if err != nil {
		return
	}
	return
}
