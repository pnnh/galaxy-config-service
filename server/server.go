package server

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/pnnh/multiverse-server/server/middleware"

	"github.com/pnnh/multiverse-server/server/helpers"

	"github.com/pnnh/multiverse-server/server/handlers/pages"

	"github.com/pnnh/multiverse-server/server/handlers"

	"github.com/pnnh/multiverse-server/server/auth"

	"github.com/pnnh/multiverse-server/config"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

type WebServer struct {
	router     *gin.Engine
	middleware *middleware.ServerMiddleware
}

func NewWebServer(smw *middleware.ServerMiddleware) (*WebServer, error) {
	router := gin.Default()

	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	server := &WebServer{
		router:     router,
		middleware: smw}
	router.SetFuncMap(helpers.FuncMap())
	router.LoadHTMLGlob("docker/templates/**/*.mst")

	corsDomain := []string{"https://multiverse.direct", "https://www.multiverse.direct"}

	if config.Debug() {
		corsDomain = append(corsDomain, "https://debug.multiverse.direct")
	}

	router.Use(cors.New(cors.Config{
		AllowOrigins: corsDomain,
		AllowMethods:     []string{"PUT", "PATCH", "POST", "GET"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	store := cookie.NewStore([]byte("secret"))
	router.Use(sessions.Sessions("mysession", store))

	return server, nil
}

func (s *WebServer) Init() error {
	indexHandler := pages.NewIndexHandler(s.middleware)
	s.router.GET("/", indexHandler.Query)

	s.router.Static("/assets", "./docker/assets")
	authHandler := &webauthnHandler{middleware: s.middleware}
	s.router.GET("/register/begin/:username", authHandler.BeginRegistration)
	s.router.POST("/register/finish/:username", authHandler.FinishRegistration)
	s.router.GET("/login/begin/:username", authHandler.BeginLogin)
	s.router.POST("/login/finish/:username", authHandler.FinishLogin)

	accountHandler := handlers.NewAccountHandler(s.middleware)

	s.router.GET("/login", accountHandler.LoginByWebAuthn)

	auth.InitOAuth2(s.router, s.middleware)

	return nil
}

func (s *WebServer) Start() error {
	port := os.Getenv("PORT")
	if len(port) < 1 {
		port = "8000"
	}
	var handler http.Handler = s

	serv := &http.Server{
		Addr:           ":" + port,
		Handler:        handler,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	if err := serv.ListenAndServe(); err != nil {
		return fmt.Errorf("服务出错停止: %w", err)
	}
	return nil
}

func (s *WebServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	//if config.Debug() && strings.HasPrefix(r.URL.Path, "/blog/") {
	//	devBlogHandler(w, r)
	//	return
	//}
	s.router.ServeHTTP(w, r)
}
