package middleware

import (
	"github.com/pnnh/multiverse-server/services/templs"
)

type ServerMiddleware struct {
	Templs *templs.Service
}
