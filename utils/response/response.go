package response

import (
	"net/http"

	"github.com/cwloo/gonet/utils/result"
	"github.com/gin-gonic/gin"
)

func Result(code int, msg string, req, data any, c *gin.Context) {
	c.JSON(http.StatusOK, result.R{
		Code:   code,
		ErrMsg: msg,
		Req:    req,
		Data:   data,
	})
}

func BadRequest(c *gin.Context) {
	c.JSON(http.StatusBadRequest, gin.H{
		"code":   http.StatusBadRequest,
		"errmsg": "bad request"})
}

func Ok(req, data any, c *gin.Context) {
	Result(0, "ok", req, data, c)
}

func OkMsg(msg string, req, data any, c *gin.Context) {
	Result(0, msg, req, data, c)
}

func Err(req, data any, c *gin.Context) {
	Result(1, "error", req, data, c)
}

func ErrMsg(msg string, req, data any, c *gin.Context) {
	Result(1, msg, req, data, c)
}
