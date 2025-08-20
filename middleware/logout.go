package middleware

import (
	"fmt"
	"net/http"

	"github.com/fastbill/go-httperrors"
)

func (m *Middleware) getLogout(w http.ResponseWriter, r *http.Request) error {
	logoutChallenge := r.URL.Query().Get("logout_challenge")
	if logoutChallenge == "" {
		return httperrors.New(http.StatusBadRequest, "'logout_challenge' missing")
	}

	logoutRequest, _, err := m.oryClient.OAuth2API.GetOAuth2LogoutRequest(m.oryAuthedContext).
		LogoutChallenge(logoutChallenge).
		Execute()
	if err != nil {
		return fmt.Errorf("could not retrieve logout request: %w", err)
	}

	err = m.renderer.RenderLogoutPage(w, r, logoutRequest)
	if err != nil {
		return err
	}

	return nil
}

func (m *Middleware) postLogout(w http.ResponseWriter, r *http.Request) error {
	logoutChallenge := r.FormValue("challenge")
	if logoutChallenge == "" {
		return httperrors.New(http.StatusBadRequest, "'logout_challenge' missing")
	}

	submit := r.FormValue("submit")

	logoutRequest, _, err := m.oryClient.OAuth2API.GetOAuth2LogoutRequest(m.oryAuthedContext).
		LogoutChallenge(logoutChallenge).
		Execute()
	if err != nil {
		return fmt.Errorf("could not retrieve logout request: %w", err)
	}

	if submit == "accept" {
		oauth2RedirectTo, _, err := m.oryClient.OAuth2API.AcceptOAuth2LogoutRequest(m.oryAuthedContext).
			LogoutChallenge(logoutChallenge).
			Execute()
		if err != nil {
			return fmt.Errorf("could not accept logout request: %w", err)
		}
		http.Redirect(w, r, oauth2RedirectTo.GetRedirectTo(), http.StatusFound)
		return nil
	} else if submit == "reject" {
		_, err := m.oryClient.OAuth2API.RejectOAuth2LogoutRequest(m.oryAuthedContext).
			LogoutChallenge(logoutChallenge).
			Execute()
		if err != nil {
			return fmt.Errorf("could not reject logout request: %w", err)
		}
		http.Redirect(w, r, logoutRequest.Client.GetClientUri(), http.StatusFound)
		return nil
	}

	return nil
}
