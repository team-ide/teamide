package module

import (
	"github.com/gin-gonic/gin"
	"regexp"
	"strings"
	"teamide/internal/base"
	"teamide/internal/module/module_toolbox"
	"teamide/pkg/db"
)

type DataRequest struct {
	Origin   string `json:"origin,omitempty"`
	Pathname string `json:"pathname,omitempty"`
}

type DataResponse struct {
	Url      string `json:"url"`
	Api      string `json:"api"`
	FilesUrl string `json:"filesUrl"`
	IsServer bool   `json:"isServer"`

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

	res = response
	return
}
