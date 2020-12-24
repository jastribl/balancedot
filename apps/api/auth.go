package api

import (
	"net/http"

	"gihub.com/jastribl/balancedot/entities"
	"golang.org/x/crypto/bcrypt"
)

type authRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// Register todo
func (m *App) Register(w ResponseWriter, r *Request) WriterResponse {
	var authRequestParams authRequest
	if err := r.DecodeParams(&authRequestParams); err != nil {
		return err
	}

	user := &entities.User{}
	err := m.db.Find(user, "username = ?", authRequestParams.Username).Error
	if err != nil {
		return err
	}

	if user.Username == authRequestParams.Username {
		return Error{
			Message: "Username already taken",
			Error:   nil,
			Code:    http.StatusConflict,
		}
	}

	// Generate "hash" to store from user password
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(authRequestParams.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
		// return w.SendError("Error while registering", http.StatusUnprocessableEntity, err)
	}

	user.Username = authRequestParams.Username
	user.Password = string(passwordHash)
	err = m.db.Create(user).Error
	if err != nil {
		return err
	}
	return w.RenderJSON(map[string]interface{}{
		"success": "ok", // todo: return something different
	})
}

// Login todo
func (m *App) Login(w ResponseWriter, r *Request) WriterResponse {
	var authRequestParams authRequest
	if err := r.DecodeParams(&authRequestParams); err != nil {
		return err
		// return w.SendParseBodyError(err)
	}

	user := &entities.User{}
	if err := m.db.Find(user, "username = ?", authRequestParams.Username).Error; err != nil {
		return err
	}

	if user.Username != authRequestParams.Username {
		return Error{
			Message: "Username not found",
			Error:   nil,
			Code:    http.StatusUnauthorized,
		}
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(authRequestParams.Password)); err != nil {
		return Error{
			Message: "Password does not match",
			Error:   nil,
			Code:    http.StatusUnauthorized,
		}
	}

	if err := m.Auth.Login(w, r, user.UUID.String()); err != nil {
		return Error{
			Message: "Error while logging in",
			Error:   err,
			Code:    http.StatusInternalServerError,
		}
	}

	return w.RenderJSON(map[string]interface{}{
		"success": "ok", // todo: return something different
	})
}

// Logout todo
func (a *App) Logout(w ResponseWriter, r *Request) WriterResponse {
	err := a.Auth.Logout(w, r)
	if err != nil {
		return w.SendError("Error while logging out", http.StatusInternalServerError, err)
	}

	return w.SendSimpleMessage("Successfully logged out")
}
