package module

import (
	"github.com/gin-gonic/gin"
	"github.com/team-ide/go-tool/db"
	"github.com/team-ide/go-tool/util"
	"regexp"
	"strings"
	"teamide/internal/context"
	"teamide/internal/module/module_toolbox"
	"teamide/pkg/base"
)

type DataRequest struct {
	Origin   string `json:"origin,omitempty"`
	Pathname string `json:"pathname,omitempty"`
	Text     string `json:"text,omitempty"`
}

type DataResponse struct {
	Url          string           `json:"url"`
	Api          string           `json:"api"`
	FilesUrl     string           `json:"filesUrl"`
	IsServer     bool             `json:"isServer"`
	ClientKey    string           `json:"clientKey"`
	ClientTabKey string           `json:"clientTabKey"`
	Setting      *context.Setting `json:"setting"`

	ToolboxTypes             []*module_toolbox.ToolboxType      `json:"toolboxTypes"`
	SqlConditionalOperations []*db.SqlConditionalOperation      `json:"sqlConditionalOperations"`
	DatabaseTypes            []*db.DatabaseType                 `json:"databaseTypes"`
	QuickCommandTypes        []*module_toolbox.QuickCommandType `json:"quickCommandTypes"`
}

func (this_ *Api) apiData(requestBean *base.RequestBean, c *gin.Context) (res interface{}, err error) {
	path := requestBean.Path[0:strings.LastIndex(requestBean.Path, "api/")]
	request := &DataRequest{}
	if !base.RequestJSON(request, c) {
		return
	}
	response := &DataResponse{}

	pathname := request.Pathname
	re, _ := regexp.Compile("/+")
	pathname = re.ReplaceAllLiteralString(pathname, "/")

	path = strings.TrimSuffix(path, "/")
	pathname = strings.TrimSuffix(pathname, "/")
	pathname = strings.TrimSuffix(pathname, path)

	if !strings.HasSuffix(pathname, "/") {
		pathname += "/"
	}

	response.Url = request.Origin + pathname
	response.Api = response.Url + "api/"
	response.FilesUrl = response.Api + "files/"
	response.IsServer = this_.IsServer

	response.ToolboxTypes = module_toolbox.GetToolboxTypes()
	response.SqlConditionalOperations = db.GetSqlConditionalOperations()
	response.DatabaseTypes = db.DatabaseTypes
	response.QuickCommandTypes = module_toolbox.GetQuickCommandTypes()

	response.Setting = this_.Setting

	if requestBean.ClientKey == "" {
		response.ClientKey = util.GetUUID()
	} else {
		response.ClientKey = requestBean.ClientKey
	}
	if requestBean.ClientTabKey == "" {
		response.ClientTabKey = util.GetUUID()
	} else {
		response.ClientTabKey = requestBean.ClientTabKey
	}

	res = response
	return
}

func (this_ *Api) apiShowPlaintext(requestBean *base.RequestBean, c *gin.Context) (res interface{}, err error) {
	request := &DataRequest{}
	if !base.RequestJSON(request, c) {
		return
	}

	res = this_.toolboxService.DecryptOptionAttr(request.Text)
	return
}
