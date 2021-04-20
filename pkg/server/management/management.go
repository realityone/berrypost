package management

import (
	"github.com/gin-gonic/gin"
)

type Management struct {
	*gin.Engine
}

func New() *Management {
	m := &Management{}
	return m
}
func (m Management) intro(ctx *gin.Context) {
	introSchema := struct {
		Name    string `json:"name"`
		Version string `json:"version"`
	}{
		Name:    "berrypost-management-server",
		Version: "0.0.1",
	}
	ctx.JSON(200, introSchema)
}

func (m Management) SetupRoute(in gin.IRouter) {
	in.GET("/api/_intro", m.intro)
	in.GET("/hello", func(ctx *gin.Context) {
		ctx.JSON(200, map[string]string{
			"data": "hello",
		})
	})
}
