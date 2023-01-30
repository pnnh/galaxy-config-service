package handlers

import (
	"net/http"
	"github.com/pnnh/multiverse-server/server/middleware"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

type accountHandler struct {
	middleware *middleware.ServerMiddleware
}

func (s *accountHandler) LoginByWebAuthn(gctx *gin.Context) {
	session := sessions.Default(gctx)
	authuser := session.Get("authuser")

	if authuser != nil {
		gctx.HTML(http.StatusOK, "account/logined.mst", gin.H{
			"username": authuser,
		})
		return
	}

	gctx.HTML(http.StatusOK, "account/webauthn.mst", gin.H{})
}

func NewAccountHandler(middleware *middleware.ServerMiddleware) *accountHandler {
	return &accountHandler{
		middleware,
	}
}
