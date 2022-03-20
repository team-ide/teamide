package module

import (
	"github.com/gin-gonic/gin"
	"regexp"
	"strings"
	base2 "teamide/internal/server/base"
	model2 "teamide/pkg/application/model"
)

type DataRequest struct {
	Origin   string `json:"origin,omitempty"`
	Pathname string `json:"pathname,omitempty"`
}

type DataResponse struct {
	Url           string                 `json:"url,omitempty"`
	Api           string                 `json:"api,omitempty"`
	IsStandAlone  bool                   `json:"isStandAlone,omitempty"`
	ColumnTypes   []*model2.ColumnType   `json:"columnTypes,omitempty"`
	DataTypes     []*model2.DataType     `json:"dataTypes,omitempty"`
	IndexTypes    []*model2.IndexType    `json:"indexTypes,omitempty"`
	ModelTypes    []*model2.ModelType    `json:"modelTypes,omitempty"`
	DataPlaces    []*model2.DataPlace    `json:"dataPlaces,omitempty"`
	DatabaseTypes []*model2.DatabaseType `json:"databaseTypes,omitempty"`
}

func apiData(requestBean *base2.RequestBean, c *gin.Context) (res interface{}, err error) {
	path := requestBean.Path[0:strings.LastIndex(requestBean.Path, "api/")]
	request := &DataRequest{}
	if !base2.RequestJSON(request, c) {
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
	response.IsStandAlone = base2.IsStandAlone
	response.ColumnTypes = model2.COLUMN_TYPES
	response.DataTypes = model2.DATA_TYPES
	response.ModelTypes = model2.MODEL_TYPES
	response.IndexTypes = model2.INDEX_TYPES
	response.DataPlaces = model2.DATA_PLACES
	response.DatabaseTypes = model2.DATABASE_TYPES

	res = response
	return
}
