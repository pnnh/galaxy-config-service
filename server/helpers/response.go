package helpers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pnnh/multiverse-server/server/protocols"
	"github.com/sirupsen/logrus"
)

func ResponseCode(gctx *gin.Context, code int) {
	message := protocols.CodeMessage(code)
	ResponseCodeMessageData(gctx, code, message, nil)
}

func ResponseCodeMessage(gctx *gin.Context, code int, message string) {
	ResponseCodeMessageData(gctx, code, message, nil)
}

func ResponseCodeMessageData(gctx *gin.Context, code int, message string, data interface{}) {
	jsonBody := gin.H{"code": code}
	if len(message) > 0 {
		jsonBody["message"] = message
	}
	if data != nil {
		jsonBody["data"] = data
	}
	gctx.JSON(http.StatusOK, jsonBody)
}

func ResponseCodeError(gctx *gin.Context, code int, err error) {
	if err != nil {
		logrus.Errorln("ResponseCodeError", gctx.FullPath(), err)
	}
	message := protocols.CodeMessage(code)
	ResponseCodeMessageData(gctx, code, message, nil)
}

func ResponseCodeMessageError(gctx *gin.Context, code int, message string, err error) {
	if err != nil {
		logrus.Errorln("ResponseCodeMessageError", gctx.FullPath(), message, err)
	}
	ResponseCodeMessageData(gctx, code, message, nil)
}

func ResponseMessageError(gctx *gin.Context, message string, err error) {
	if err != nil {
		logrus.Errorln("ResponseMessageError", message, err)
	}
	ResponseCodeMessageData(gctx, protocols.CodeError, message, nil)
}
