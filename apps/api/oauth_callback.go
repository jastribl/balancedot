package api

import (
	"net/http"
)

// OauthCallback is the callback for splitwise Oauth
func (m *App) OauthCallback(w ResponseWriter, r *Request) WriterResponse {
	if r.URL.Query().Get("state") != m.config.Splitwise.State {
		return w.SendError("Invalid state in callback", http.StatusInternalServerError)
	}

	m.config.Splitwise.AuthCodeChan <- r.URL.Query().Get("code")
	return w.SendSimpleMessage("Processing Token...\nPlease close this tab and return to your application.")
}
