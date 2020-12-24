package auth

import (
	"net/http"

	"gihub.com/jastribl/balancedot/apps/api"
)

type Auth interface {
	UnAuthorizedRequest(next api.Handler) api.Handler
	AuthorizedRequest(next api.Handler, allowedRoles ...entities.RoleType) api.Handler
	Login(w api.ResponseWriter, r *http.Request, uuid string) error
	Logout(w api.ResponseWriter, r *http.Request) error
	Close()
}
