package authorizationserver

import (
	"net/http"

	"github.com/sirupsen/logrus"
)

func IntrospectionEndpoint(rw http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	mySessionData := newSession("")
	ir, err := oauth2.NewIntrospectionRequest(ctx, req, mySessionData)
	if err != nil {
		logrus.Printf("Error occurred in NewIntrospectionRequest: %+v", err)
		oauth2.WriteIntrospectionError(ctx, rw, err)
		return
	}

	oauth2.WriteIntrospectionResponse(ctx, rw, ir)
}
