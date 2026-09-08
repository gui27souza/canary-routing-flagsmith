package testutil

import "github.com/gin-gonic/gin"

func MockRouter(
	method, endpoint string, handlers ...gin.HandlerFunc,
) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Handle(method, endpoint, handlers...)
	return r
}
