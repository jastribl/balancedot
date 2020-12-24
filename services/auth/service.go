package auth

import (
	"fmt"
	"net/http"
	"os"

	"gihub.com/jastribl/balancedot/apps/api"
	"gihub.com/jastribl/balancedot/repos"
	"github.com/gorilla/sessions"
	"github.com/jinzhu/gorm"
	"gopkg.in/boj/redistore.v1"
)

const (
	authenticatedKey = "authenticated"
	uuidKey          = "uuid"
)

// APIAuth todo
type APIAuth struct {
	cookieName   string
	sessionStore *redistore.RediStore
	db           *gorm.DB
}

// NewAPIAuth todo
func NewAPIAuth(db *gorm.DB, cookieName string, keyPairStrings ...string) (*APIAuth, error) {
	var keyPairs [][]byte
	for _, keyPairString := range keyPairStrings {
		keyPairs = append(keyPairs, []byte(keyPairString))
	}
	redisStore, err := redistore.NewRediStore(
		10,
		"tcp",
		fmt.Sprintf("%s:6379", os.Getenv("REDIS_URL")),
		"",
		keyPairs...,
	)
	if err != nil {
		return nil, err
	}
	return &APIAuth{
		db:           db,
		cookieName:   cookieName,
		sessionStore: redisStore,
	}, nil
}

// Close todo
func (a *APIAuth) Close() {
	a.sessionStore.Close()
	// don't need to close Db since someone else owns it
}

func (a *APIAuth) getSession(r *http.Request) (*sessions.Session, error) {
	return a.sessionStore.Get(r, a.cookieName)
}

// Login sets the current request session to be authenticated.
// The caller of this function is required to perform authentication before calling
// this function to ensure that the user is indeed valid.
func (a *APIAuth) Login(w ResponseWriter, r *http.Request, uuid string) error {
	session, err := a.getSession(r)

	if err != nil {
		return err
	}

	// Set user as authenticated
	session.Values[authenticatedKey] = true
	session.Values[uuidKey] = uuid
	return session.Save(r, w.GetUnderlyingWriter())
}

// Logout revokes the users authentication.
func (a *APIAuth) Logout(w ResponseWriter, r *http.Request) error {
	session, err := a.getSession(r)

	if err != nil {
		return err
	}

	// Revoke users authentication
	session.Options.MaxAge = -1
	session.Values[authenticatedKey] = false
	return session.Save(r, w.GetUnderlyingWriter())
}

// UnAuthorizedRequest todo
func (a *APIAuth) UnAuthorizedRequest(next api.Handler) api.Handler {
	return func(w api.ResponseWriter, r *http.Request) interface{} {
		// Delegate request to the given authorized handle
		next(w, r)
	}
}

// AuthorizedRequest is a wrapper around normal requests.
// It will only allow requests for users that are authorized.
// Otherwise it will return an error
func (a *APIAuth) AuthorizedRequest(next api.Handler) api.Handler {
	m := make(map[entities.RoleType]bool)
	for _, allowedRole := range allowedRoles {
		m[allowedRole] = true
	}
	userRepo := repos.NewUserRepo(a.db)

	return func(w api.ResponseWriter, r *http.Request) interface{} {
		if session, err := a.getSession(r); err == nil {
			authenticated, ok1 := session.Values[authenticatedKey].(bool)
			loggedInUUID, ok2 := session.Values[uuidKey].(string)
			if ok1 && ok2 && authenticated {
				if user, err := userRepo.GetUseByUUID(loggedInUUID); err == nil {
					// Delegate request to the given authorized handle
					// todo: somehow attach the user
					return next(w, r)
				}
			}
		}
		// return an error to the client
		return api.Error{
			Message: http.StatusText(http.StatusUnauthorized),
			Code:    http.StatusUnauthorized,
			Error:   nil,
		}
	}
}
