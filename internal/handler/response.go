package handler

import "github.com/gin-gonic/gin"

type Response struct {
	Data   interface{}       `json:"data,omitempty"`
	Errors map[string]string `json:"errors,omitempty"`
	Meta   interface{}       `json:"meta,omitempty"`
}

func JSONOK(c *gin.Context, data interface{}) {
	c.JSON(200, Response{Data: data})
}

func JSONCreated(c *gin.Context, data interface{}) {
	c.JSON(201, Response{Data: data})
}

func JSONError(c *gin.Context, code int, errs map[string]string) {
	c.JSON(code, Response{Errors: errs})
}

func JSONErrorMsg(c *gin.Context, code int, msg string) {
	c.JSON(code, Response{Errors: map[string]string{"error": msg}})
}
