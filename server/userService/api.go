package userService

import "github.com/gin-gonic/gin"

func index(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "Wellcome Index",
	})
}

func BindApi(apiCache map[string]func(c *gin.Context)) {
	apiCache["user/index"] = index
}
