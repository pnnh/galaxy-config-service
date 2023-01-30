package handlers

import (
	"net/http"
	"strings"

	"github.com/pnnh/multiverse-server/server/middleware"
	"github.com/pnnh/multiverse-server/server/models"
	"github.com/sirupsen/logrus"
 
	"github.com/gin-gonic/gin"
)

type accountHandler struct {
	middleware *middleware.ServerMiddleware
}

func (s *accountHandler) LoginByWebAuthn(gctx *gin.Context) {
	
	project, ok := gctx.GetQuery("project")
	if !ok {
		gctx.Status(http.StatusBadRequest)
		return
	}

	configs, err :=	models.SelectConfigs(project)
	if err != nil {
		gctx.Status(http.StatusInternalServerError)
		logrus.Errorln("SelectConfigs: ", err)
		return
	}
	sb := &strings.Builder{}
	for _, v := range configs {
		sb.WriteString(v.Key)
		sb.WriteString("=")
		sb.WriteString(v.Value)
	}

	content := sb.String()

	gctx.String(http.StatusOK, "%s", content)
}

func NewAccountHandler(middleware *middleware.ServerMiddleware) *accountHandler {
	return &accountHandler{
		middleware,
	}
}

 



