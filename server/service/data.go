package service

import (
	"regexp"
	"strings"
	"teamide/application/model"
	"teamide/server/base"

	"github.com/gin-gonic/gin"
)

type DataRequest struct {
	Origin   string `json:"origin,omitempty"`
	Pathname string `json:"pathname,omitempty"`
}

type DataResponse struct {
	Url           string                `json:"url,omitempty"`
	Api           string                `json:"api,omitempty"`
	IsStandAlone  bool                  `json:"isStandAlone,omitempty"`
	ColumnTypes   []*model.ColumnType   `json:"columnTypes,omitempty"`
	DataTypes     []*model.DataType     `json:"dataTypes,omitempty"`
	IndexTypes    []*model.IndexType    `json:"indexTypes,omitempty"`
	ModelTypes    []*model.ModelType    `json:"modelTypes,omitempty"`
	DataPlaces    []*model.DataPlace    `json:"dataPlaces,omitempty"`
	DatabaseTypes []*model.DatabaseType `json:"databaseTypes,omitempty"`
}

func apiData(requestBean *base.RequestBean, c *gin.Context) (res interface{}, err error) {
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
	response.IsStandAlone = base.IS_STAND_ALONE
	response.ColumnTypes = model.COLUMN_TYPES
	response.DataTypes = model.DATA_TYPES
	response.ModelTypes = model.MODEL_TYPES
	response.IndexTypes = model.INDEX_TYPES
	response.DataPlaces = model.DATA_PLACES
	response.DatabaseTypes = model.DATABASE_TYPES

	res = response
	return
}
