package router

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"io"
)

// GinRouter - Wrapper on gin that implements router interface
type GinRouter struct {
	gin    *gin.Engine
	logger *zap.Logger
	port   string
}

func New(logger *zap.Logger, port string) *GinRouter {
	r := gin.Default()
	return &GinRouter{logger: logger, gin: r, port: port}
}

func (r *GinRouter) Get(path string, handler func(string) (int, map[string]any)) {
	param := ""
	for i := range path {
		if path[i] == ':' {
			param = path[i+1:]
			break
		}
	}

	r.gin.GET(path, func(c *gin.Context) {
		p := ""
		if param != "" {
			p = c.Param(param)
		}

		code, resp := handler(p)
		c.JSON(code, resp)
	})
}

func (r *GinRouter) GetFile(path string, handler func(string) (int, string)) {
	param := ""
	for i := range path {
		if path[i] == ':' {
			param = path[i+1:]
			break
		}
	}

	r.gin.GET(path, func(c *gin.Context) {
		p := ""
		if param != "" {
			p = c.Param(param)
		}

		code, str := handler(p)
		if code != 200 {
			c.JSON(code, str)
		} else {
			c.File(str)
		}
	})
}

func (r *GinRouter) Post(path string, handler func([]byte) (int, map[string]any)) {
	r.gin.POST(path, func(c *gin.Context) {
		rawBody, err := io.ReadAll(c.Request.Body)
		if err != nil {
			r.logger.Error("Unable to read body", zap.Error(err), zap.Any("request", c.Request))
			c.AbortWithStatusJSON(500, gin.H{"error": err})
		}
		code, resp := handler(rawBody)
		c.JSON(code, resp)
	})
}

func (r *GinRouter) Run() {
	r.logger.Info("Gin Router start")
	r.gin.Run(":" + r.port)
}
