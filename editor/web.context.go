package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"teamide/application/base"
	"teamide/application/model"

	"github.com/gin-gonic/gin"
)

type ContextRequest struct {
	AppName string `json:"appName,omitempty"`
}

func doApiContextGet(path string, c *gin.Context) (res interface{}, err error) {
	request := &ContextRequest{}
	if !RequestJSON(request, c) {
		return
	}
	var context = &model.ModelContext{}
	if base.IsNotEmpty(request.AppName) {

		var bs []byte
		bs, err = ReadFile(appsDir + "/" + request.AppName)
		if err != nil {
			return
		}
		if len(bs) > 0 {
			bs, err = base.AesCBCDecrypt(bs, []byte(K))
			if err != nil {
				return
			}
			err = base.ToBean(bs, context)
			if err != nil {
				return
			}
		}
	}

	res = context
	return
}

type ContextSaveRequest struct {
	AppName string `json:"appName,omitempty"`
	Content string `json:"content,omitempty"`
}

func doApiContextSave(path string, c *gin.Context) (res interface{}, err error) {
	request := &ContextSaveRequest{}
	if !RequestJSON(request, c) {
		return
	}

	filename := appsDir + "/" + request.AppName
	var exists bool
	exists, err = base.PathExists(filename)
	if err != nil {
		return
	}
	if !exists {
		err = errors.New(fmt.Sprint("应用[", request.AppName, "]不存在"))
		return
	}
	var context = &model.ModelContext{}

	err = base.ToBean([]byte(request.Content), context)
	if err != nil {
		return
	}
	var bs []byte
	bs, err = json.Marshal(context)
	if err != nil {
		return
	}
	bs, err = base.AesCBCEncrypt(bs, []byte(K))
	if err != nil {
		return
	}
	err = WriteFile(filename, bs)
	if err != nil {
		return
	}

	res = context
	return
}
