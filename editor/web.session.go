package main

import (
	"github.com/gin-gonic/gin"
)

type SessionRequest struct {
}

type SessionResponse struct {
}

func doApiSession(path string, c *gin.Context) (res interface{}, err error) {

	response := &SessionResponse{}

	res = response
	return
}
