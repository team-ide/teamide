package main

import (
	"teamide/application/base"

	"github.com/gin-gonic/gin"
)

type HttpResponse struct {
	Code  string      `json:"code"`
	Msg   string      `json:"msg,omitempty"`
	Value interface{} `json:"value,omitempty"`
}

func toWebResponseJSON(value interface{}, err error) *HttpResponse {
	response := &HttpResponse{
		Code: "0",
		Msg:  "成功",
	}
	if err != nil {
		response.Msg = err.Error()
		baseErr, baseErrOk := err.(*base.ErrorBase)
		if baseErrOk {
			response.Code = baseErr.Code
			response.Msg = baseErr.Msg
		} else {
			response.Code = "-1"
		}
	} else {
		response.Value = value
	}
	return response
}

func RequestJSON(data interface{}, c *gin.Context) bool {
	err := c.BindJSON(data)
	if err != nil {
		c.JSON(200, toWebResponseJSON(nil, err))
		return false
	}
	return true
}
