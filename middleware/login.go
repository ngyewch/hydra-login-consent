package middleware

import (
	"fmt"
	"github.com/fastbill/go-httperrors"
	ory "github.com/ory/client-go"
	"net/http"
	"net/url"
)

func (m *Middleware) getLogin(w http.ResponseWriter, r *http.Request) error {
	loginChallenge := r.URL.Query().Get("login_challenge")
	if loginChallenge == "" {
		return httperrors.New(http.StatusBadRequest, "'login_challenge' missing")
	}

	loginRequest, redirected, err := m.handleLogin(w, r, loginChallenge)
	if err != nil {
		return err
	}
	if redirected {
		return nil
	}

	errorMessage := r.URL.Query().Get("error_message")

	err = m.renderer.RenderLoginPage(w, r, loginRequest, errorMessage)
	if err != nil {
		return err
	}

	return nil
}

func (m *Middleware) postLogin(w http.ResponseWriter, r *http.Request) error {
	loginChallenge := r.FormValue("challenge")
	if loginChallenge == "" {
		return httperrors.New(http.StatusBadRequest, "'login_challenge' missing")
	}

	remember := r.FormValue("remember") != ""

	_, redirected, err := m.handleLogin(w, r, loginChallenge)
	if err != nil {
		return err
	}
	if redirected {
		return nil
	}

	subject, err := m.handler.HandleLogin(r)
	if err != nil {
		return err
	}
	if subject == "" {
		redirectUrl, err := url.Parse(r.URL.String())
		if err != nil {
			return err
		}
		values := redirectUrl.Query()
		values.Set("error_message", "Incorrect sign in")
		redirectUrl.RawQuery = values.Encode()
		http.Redirect(w, r, redirectUrl.String(), http.StatusFound)
		return nil
	}

	oauth2RedirectTo, _, err := m.oryClient.OAuth2API.AcceptOAuth2LoginRequest(m.oryAuthedContext).
		LoginChallenge(loginChallenge).
		AcceptOAuth2LoginRequest(ory.AcceptOAuth2LoginRequest{
			Subject:     subject,
			Remember:    ory.PtrBool(remember),
			RememberFor: ory.PtrInt64(3600),
		}).
		Execute()
	if err != nil {
		return err
	}

	http.Redirect(w, r, oauth2RedirectTo.GetRedirectTo(), http.StatusFound)

	return nil
}

func (m *Middleware) handleLogin(w http.ResponseWriter, r *http.Request, loginChallenge string) (*ory.OAuth2LoginRequest, bool, error) {
	loginRequest, _, err := m.oryClient.OAuth2API.GetOAuth2LoginRequest(m.oryAuthedContext).
		LoginChallenge(loginChallenge).
		Execute()
	if err != nil {
		return nil, false, fmt.Errorf("could not retrieve login request: %w", err)
	}

	if loginRequest.Skip {
		oauth2RedirectTo, _, err := m.oryClient.OAuth2API.AcceptOAuth2LoginRequest(m.oryAuthedContext).
			LoginChallenge(loginChallenge).
			AcceptOAuth2LoginRequest(ory.AcceptOAuth2LoginRequest{
				Subject: loginRequest.GetSubject(),
			}).
			Execute()
		if err != nil {
			return nil, false, fmt.Errorf("could not accept login request: %w", err)
		}
		http.Redirect(w, r, oauth2RedirectTo.GetRedirectTo(), http.StatusFound)
		return loginRequest, true, nil
	}

	return loginRequest, false, nil
}
