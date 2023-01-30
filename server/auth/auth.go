package auth

import (
	"github.com/pnnh/multiverse-server/server/middleware"

	"github.com/pnnh/multiverse-server/server/auth/authorizationserver"

	"github.com/gin-gonic/gin"
)

// A valid oauth2 client (check the store) that additionally requests an OpenID Connect id token
// var ClientConf = goauth.Config{
// 	ClientID:     "my-client",
// 	ClientSecret: "foobar",
// 	RedirectURL:  "https://drm.sfx.xyz/callback",
// 	Scopes:       []string{"photos", "openid", "offline"},
// 	Endpoint: goauth.Endpoint{
// 		TokenURL: "https://drm.sfx.xyz/oauth2/token",
// 		AuthURL:  "https://drm.sfx.xyz/oauth2/auth",
// 	},
// }

// The same thing (valid oauth2 client) but for using the client credentials grant
// var appClientConf = clientcredentials.Config{
// 	ClientID:     "my-client",
// 	ClientSecret: "foobar",
// 	Scopes:       []string{"fosite"},
// 	TokenURL:     "https://drm.sfx.xyz/oauth2/token",
// }

// Samle client as above, but using a different secret to demonstrate secret rotation
// var appClientConfRotated = clientcredentials.Config{
// 	ClientID:     "my-client",
// 	ClientSecret: "foobaz",
// 	Scopes:       []string{"fosite"},
// 	TokenURL:     "https://drm.sfx.xyz/oauth2/token",
// }

func InitOAuth2(router *gin.Engine, middleware *middleware.ServerMiddleware) {
	router.GET("/oauth2/auth", func(gctx *gin.Context) {
		authorizationserver.AuthEndpointHtml(gctx)
	})
	router.POST("/oauth2/auth", func(gctx *gin.Context) {
		authorizationserver.AuthEndpointJson(gctx)
	})
	router.POST("/oauth2/token", func(gctx *gin.Context) {
		authorizationserver.TokenEndpoint(gctx.Writer, gctx.Request)
	})
	// router.GET("/oauth2/revoke", func(gctx *gin.Context) {
	// 	authorizationserver.RevokeEndpoint(gctx.Writer, gctx.Request)
	// })
	// router.GET("/oauth2/introspect", func(gctx *gin.Context) {
	// 	authorizationserver.IntrospectionEndpoint(gctx.Writer, gctx.Request)
	// })

	// router.GET("/oauth2/protected", func(gctx *gin.Context) {
	// 	resourceserver.ProtectedEndpoint(appClientConf)(gctx.Writer, gctx.Request)
	// })
}
