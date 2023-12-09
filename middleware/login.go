package middleware

import (
	"fmt"
	"github.com/fastbill/go-httperrors"
	"github.com/gorilla/csrf"
	ory "github.com/ory/client-go"
	"html/template"
	"net/http"
	"net/url"
)

type ProviderInfo struct {
	Name string
}

type LoginTemplateData struct {
	Provider           ProviderInfo
	Request            *ory.OAuth2LoginRequest
	ForgotPasswordText string
	ForgotPasswordUri  string
	CSRFToken          string
	CSRFTemplateField  template.HTML
	ErrorMessage       string
}

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

	err = m.renderPage(w, "login.gohtml",
		LoginTemplateData{
			Provider: ProviderInfo{
				Name: m.cfg.Name,
			},
			Request:            loginRequest,
			ForgotPasswordText: m.cfg.ForgotPasswordText,
			ForgotPasswordUri:  m.cfg.ForgotPasswordUri,
			CSRFToken:          csrf.Token(r),
			CSRFTemplateField:  csrf.TemplateField(r),
			ErrorMessage:       errorMessage,
		},
	)
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

	email := r.FormValue("email")
	password := r.FormValue("password")
	remember := r.FormValue("remember") != ""

	_, redirected, err := m.handleLogin(w, r, loginChallenge)
	if err != nil {
		return err
	}
	if redirected {
		return nil
	}

	successfulSignIn, err := m.provider.Validate(email, password)
	if err != nil {
		return err
	}
	if !successfulSignIn {
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
			Subject:     email,
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
