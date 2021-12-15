package base

import "github.com/gin-gonic/gin"

type ApiWorker struct {
	Api     string
	Do      func(request *RequestBean, c *gin.Context) (res interface{}, err error)
	DoOther func(request *RequestBean, c *gin.Context)
}
