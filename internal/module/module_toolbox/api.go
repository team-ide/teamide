package module_toolbox

import (
	"github.com/gin-gonic/gin"
	"teamide/internal/base"
)

type ToolboxApi struct {
	ToolboxService *ToolboxService
}

func NewToolboxApi(ToolboxService *ToolboxService) *ToolboxApi {

	return &ToolboxApi{
		ToolboxService: ToolboxService,
	}
}

func (this_ *ToolboxApi) GetApis() (apis []*base.ApiWorker) {
	apis = append(apis, &base.ApiWorker{Apis: []string{"toolbox/context"}, Power: base.PowerToolboxContext, Do: this_.context})
	apis = append(apis, &base.ApiWorker{Apis: []string{"toolbox/insert"}, Power: base.PowerToolboxInsert, Do: this_.insert})
	apis = append(apis, &base.ApiWorker{Apis: []string{"toolbox/update"}, Power: base.PowerToolboxUpdate, Do: this_.update})
	apis = append(apis, &base.ApiWorker{Apis: []string{"toolbox/rename"}, Power: base.PowerToolboxRename, Do: this_.rename})
	apis = append(apis, &base.ApiWorker{Apis: []string{"toolbox/delete"}, Power: base.PowerToolboxDelete, Do: this_.delete})
	apis = append(apis, &base.ApiWorker{Apis: []string{"toolbox/work"}, Power: base.PowerToolboxDelete, Do: this_.work})
	//apis = append(apis, &base.ApiWorker{Apis: []string{"ws/toolbox/ssh/connection"}, Power: base.PowerToolboxSSHConnection, Do: this_.sshConnection})

	return
}

type ContextRequest struct {
}

type ContextResponse struct {
	Context map[string][]*ToolboxModel `json:"context,omitempty"`
}

func (this_ *ToolboxApi) context(requestBean *base.RequestBean, c *gin.Context) (res interface{}, err error) {

	request := &ContextRequest{}
	if !base.RequestJSON(request, c) {
		return
	}
	response := &ContextResponse{}

	list, err := this_.ToolboxService.Query(&ToolboxModel{
		UserId: requestBean.JWT.UserId,
	})
	if err != nil {
		return
	}
	context := make(map[string][]*ToolboxModel)
	for _, one := range list {
		context[one.ToolboxType] = append(context[one.ToolboxType], one)
	}
	response.Context = context
	res = response
	return
}

type InsertRequest struct {
	*ToolboxModel
}

type InsertResponse struct {
}

func (this_ *ToolboxApi) insert(requestBean *base.RequestBean, c *gin.Context) (res interface{}, err error) {

	request := &InsertRequest{}
	if !base.RequestJSON(request, c) {
		return
	}
	response := &InsertResponse{}

	toolbox := request.ToolboxModel
	toolbox.UserId = requestBean.JWT.UserId

	_, err = this_.ToolboxService.Insert(toolbox)
	if err != nil {
		return
	}

	res = response
	return
}

type UpdateRequest struct {
	*ToolboxModel
}

type UpdateResponse struct {
}

func (this_ *ToolboxApi) update(requestBean *base.RequestBean, c *gin.Context) (res interface{}, err error) {

	request := &UpdateRequest{}
	if !base.RequestJSON(request, c) {
		return
	}
	response := &UpdateResponse{}

	toolbox := request.ToolboxModel

	_, err = this_.ToolboxService.Update(toolbox)
	if err != nil {
		return
	}

	res = response
	return
}

type RenameRequest struct {
	*ToolboxModel
}

type RenameResponse struct {
}

func (this_ *ToolboxApi) rename(requestBean *base.RequestBean, c *gin.Context) (res interface{}, err error) {

	request := &RenameRequest{}
	if !base.RequestJSON(request, c) {
		return
	}
	response := &RenameResponse{}

	toolbox := request.ToolboxModel

	_, err = this_.ToolboxService.Rename(toolbox.ToolboxId, toolbox.Name)
	if err != nil {
		return
	}

	res = response
	return
}

type DeleteRequest struct {
	*ToolboxModel
}

type DeleteResponse struct {
}

func (this_ *ToolboxApi) delete(requestBean *base.RequestBean, c *gin.Context) (res interface{}, err error) {

	request := &RenameRequest{}
	if !base.RequestJSON(request, c) {
		return
	}
	response := &RenameResponse{}

	toolbox := request.ToolboxModel

	_, err = this_.ToolboxService.Delete(toolbox.ToolboxId)
	if err != nil {
		return
	}

	res = response
	return
}

type WorkRequest struct {
	ToolboxId int64                  `json:"toolboxId,omitempty"`
	Work      string                 `json:"work,omitempty"`
	Data      map[string]interface{} `json:"data,omitempty"`
}

func (this_ *ToolboxApi) work(requestBean *base.RequestBean, c *gin.Context) (res interface{}, err error) {

	request := &WorkRequest{}
	if !base.RequestJSON(request, c) {
		return
	}

	res, err = this_.ToolboxService.Work(request.ToolboxId, request.Work, request.Data)
	if err != nil {
		return
	}

	return
}
