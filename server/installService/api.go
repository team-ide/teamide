package installService

import "github.com/gin-gonic/gin"

func index(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "Wellcome Index",
	})
}

func BindApi(root string, r *gin.Engine) {
	r.GET(root+"install/index", index)
}
