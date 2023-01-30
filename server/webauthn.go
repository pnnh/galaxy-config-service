package server

import (
	"encoding/json"
	"fmt" 
	"net/http"
	"strings"

	"github.com/pnnh/multiverse-server/server/protocols"

	"github.com/pnnh/multiverse-server/server/models"

	"github.com/pnnh/multiverse-server/server/middleware"

	"github.com/pnnh/multiverse-server/server/helpers"

	"github.com/pnnh/multiverse-server/config"

	"github.com/duo-labs/webauthn.io/session"
	"github.com/duo-labs/webauthn/protocol"
	"github.com/duo-labs/webauthn/webauthn"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

var webAuthn *webauthn.WebAuthn
var sessionStore *session.Store

// Your initialization function
func init() {
	webauthnConfig := &webauthn.Config{
		RPDisplayName: "Multiverse",                                          // Display Name for your site
		RPID:          "multiverse.direct",                                // Generally the FQDN for your site
		RPOrigin:      "https://multiverse.direct",                        // The origin URL for WebAuthn requests
		RPIcon:        "https://multiverse.direct/static/images/logo.png", // Optional icon URL for your site
	}
	if config.Debug() {
		webauthnConfig.RPID = "debug.multiverse.direct"
		webauthnConfig.RPOrigin = "http://debug.multiverse.direct"
		webauthnConfig.Debug = true
	}
	var err error
	webAuthn, err = webauthn.New(webauthnConfig)
	if err != nil {
		logrus.Fatalln("webauthn初始化出错: %w", err)
	}

	sessionStore, err = session.NewStore()
	if err != nil {
		logrus.Fatal("failed to create session store:", err)
	}
}

type webauthnHandler struct {
	middleware *middleware.ServerMiddleware
}

func (s *webauthnHandler) BeginRegistration(gctx *gin.Context) {

	username := gctx.Param("username")
	if len(username) < 1 {
		helpers.ResponseCode(gctx, protocols.CodeInvalidParameter)
		return
	}

	// get user
	model, err := models.GetAccount(username)
	// user doesn't exist, create new user
	if err != nil {
		helpers.ResponseCodeMessageError(gctx, protocols.CodeError, "GetAccount error", err)
		return
	}
	if model == nil {
		displayName := strings.Split(username, "@")[0]
		model = models.NewAccountModel(username, displayName)

		if err = models.PutAccount(model); err != nil {
			helpers.ResponseMessageError(gctx, "PutAccount error", err)
			return
		}
	}

	registerOptions := func(credCreationOpts *protocol.PublicKeyCredentialCreationOptions) {
		credCreationOpts.CredentialExcludeList = model.CredentialExcludeList()
	}

	// generate PublicKeyCredentialCreationOptions, session data
	options, sessionData, err := webAuthn.BeginRegistration(
		model,
		registerOptions,
	)

	if err != nil {
		helpers.ResponseMessageError(gctx, "参数有误2", err)
		return
	}

	// store session data as marshaled JSON
	err = sessionStore.SaveWebauthnSession("registration", sessionData, gctx.Request, gctx.Writer)
	if err != nil {
		helpers.ResponseMessageError(gctx, "参数有误3", err)
		return
	}

	jsonResponse(gctx.Writer, options, http.StatusOK)
}

func (s *webauthnHandler) FinishRegistration(gctx *gin.Context) {

	username := gctx.Param("username")
	if len(username) < 1 {
		helpers.ResponseMessageError(gctx, "参数有误a", nil)
		return
	}

	// get user
	user, err := models.GetAccount(username)
	// user doesn't exist
	if err != nil {
		helpers.ResponseMessageError(gctx, "参数有误5", err)
		return
	}
	if user == nil {
		helpers.ResponseMessageError(gctx, fmt.Sprintf("GetAccount结果为空: %s", username), nil)
		return
	}

	// load the session data
	sessionData, err := sessionStore.GetWebauthnSession("registration", gctx.Request)
	if err != nil {
		helpers.ResponseMessageError(gctx, "参数有误6", err)
		return
	}

	// bodyBytes, err := ioutil.ReadAll(gctx.Request.Body)
	// bodyString := string(bodyBytes)
	// logrus.Debugln("bodyString", bodyString)

	credential, err := webAuthn.FinishRegistration(user, sessionData, gctx.Request)
	if err != nil {
		helpers.ResponseMessageError(gctx, "参数有误37", err)
		return
	}

	user.AddCredential(*credential)

	err = models.UpdateAccountCredentials(user)
	if err != nil {
		helpers.ResponseMessageError(gctx, "UpdateAccountCredentials: %w", err)
		return
	}

	jsonResponse(gctx.Writer, "Registration Success", http.StatusOK)
}

func (s *webauthnHandler) BeginLogin(gctx *gin.Context) {

	username := gctx.Param("username")
	if len(username) < 1 {
		helpers.ResponseMessageError(gctx, "参数有误b", nil)
		return
	}

	// get user
	user, err := models.GetAccount(username)

	// user doesn't exist
	if err != nil {
		helpers.ResponseMessageError(gctx, "参数有误316", err)
		return
	}

	if user == nil {
		helpers.ResponseCode(gctx, protocols.CodeAccountNotExists)
		return
	}

	// generate PublicKeyCredentialRequestOptions, session data
	options, sessionData, err := webAuthn.BeginLogin(user)
	if err != nil {
		helpers.ResponseMessageError(gctx, "参数有误39", err)
		return
	}

	// store session data as marshaled JSON
	err = sessionStore.SaveWebauthnSession("authentication", sessionData, gctx.Request, gctx.Writer)
	if err != nil {
		helpers.ResponseMessageError(gctx, "参数有误310", err)
		return
	}

	jsonResponse(gctx.Writer, options, http.StatusOK)
}

func (s *webauthnHandler) FinishLogin(gctx *gin.Context) {

	username := gctx.Param("username")
	if len(username) < 1 {
		helpers.ResponseMessageError(gctx, "参数有误", nil)
		return
	}

	// get user
	user, err := models.GetAccount(username)

	// user doesn't exist
	if err != nil {
		helpers.ResponseMessageError(gctx, "参数有误312", err)
		return
	}

	// load the session data
	sessionData, err := sessionStore.GetWebauthnSession("authentication", gctx.Request)
	if err != nil {
		helpers.ResponseMessageError(gctx, "参数有误314", err)
		return
	}

	// in an actual implementation, we should perform additional checks on
	// the returned 'credential', i.e. check 'credential.Authenticator.CloneWarning'
	// and then increment the credentials counter
	_, err = webAuthn.FinishLogin(user, sessionData, gctx.Request)
	if err != nil {
		helpers.ResponseMessageError(gctx, "参数有误315", err)
		return
	}
	session := sessions.Default(gctx)

	session.Set("authuser", username)
	session.Save()

	// handle successful login
	jsonResponse(gctx.Writer, "Login Success", http.StatusOK)
}

func jsonResponse(w http.ResponseWriter, d interface{}, c int) {
	dj, err := json.Marshal(d)
	if err != nil {
		http.Error(w, "Error creating JSON response", http.StatusInternalServerError)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(c)
	fmt.Fprintf(w, "%s", dj)
}
