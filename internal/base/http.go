package base

import (
	"github.com/gin-gonic/gin/binding"
	"net/http"

	"github.com/gin-gonic/gin"
)

type RequestBean struct {
	JWT    *JWTBean
	Path   string
	Action string
}

var (
	HttpNotResponse = &HttpResponse{}
)

type JWTBean struct {
	Sign    string `json:"sign,omitempty"`
	UserId  int64  `json:"userId,omitempty"`
	Name    string `json:"name,omitempty"`
	Account string `json:"account,omitempty"`
	Time    int64  `json:"time,omitempty"`
	LoginId int64  `json:"loginId,omitempty"`
}

type HttpResponse struct {
	Code string      `json:"code"`
	Msg  string      `json:"msg,omitempty"`
	Data interface{} `json:"data,omitempty"`
}

func RequestJSON(data interface{}, c *gin.Context) bool {
	err := c.ShouldBindBodyWith(data, binding.JSON)
	if err != nil {
		ResponseJSON(nil, err, c)
		return false
	}
	return true
}

func ResponseJSON(data interface{}, err error, c *gin.Context) {
	response := HttpResponse{
		Code: "0",
		Data: data,
	}
	if err != nil {
		response.Msg = err.Error()
		baseErr := ToBaseError(err)
		if baseErr != nil {
			response.Code = baseErr.Code
		} else {
			response.Code = "-1"
		}
	} else {
		response.Data = data
	}
	c.JSON(http.StatusOK, response)
}
