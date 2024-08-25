package middleware

import (
	"fmt"
	"github.com/fastbill/go-httperrors"
	ory "github.com/ory/client-go"
	"net/http"
)

func (m *Middleware) ConsentHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		err := m.getConsent(w, r)
		if err != nil {
			m.handleError(w, r, err)
			return
		}
	case http.MethodPost:
		err := m.postConsent(w, r)
		if err != nil {
			m.handleError(w, r, err)
			return
		}
	default:
		http.Error(w, fmt.Sprintf("Method %s not allowed", r.Method), http.StatusMethodNotAllowed)
	}
}

func (m *Middleware) getConsent(w http.ResponseWriter, r *http.Request) error {
	consentChallenge := r.URL.Query().Get("consent_challenge")
	if consentChallenge == "" {
		return httperrors.New(http.StatusBadRequest, "'consent_challenge' missing")
	}

	consentRequest, redirected, err := m.handleConsent(w, r, consentChallenge)
	if err != nil {
		return err
	}
	if redirected {
		return nil
	}

	err = m.renderer.RenderConsentPage(w, r, consentRequest)
	if err != nil {
		return err
	}

	return nil
}

func (m *Middleware) postConsent(w http.ResponseWriter, r *http.Request) error {
	consentChallenge := r.FormValue("challenge")
	if consentChallenge == "" {
		return httperrors.New(http.StatusBadRequest, "'consent_challenge' missing")
	}

	submit := r.FormValue("submit")

	consentRequest, redirected, err := m.handleConsent(w, r, consentChallenge)
	if err != nil {
		return err
	}
	if redirected {
		return nil
	}

	if submit == "accept" {
		idToken := make(map[string]interface{})
		session := &ory.AcceptOAuth2ConsentRequestSession{
			IdToken: idToken,
		}
		err = m.handler.PopulateClaims(consentRequest, idToken)
		if err != nil {
			return err
		}

		oauth2RedirectTo, _, err := m.oryClient.OAuth2API.AcceptOAuth2ConsentRequest(m.oryAuthedContext).
			ConsentChallenge(consentChallenge).
			AcceptOAuth2ConsentRequest(ory.AcceptOAuth2ConsentRequest{
				GrantScope:               consentRequest.RequestedScope,
				GrantAccessTokenAudience: consentRequest.RequestedAccessTokenAudience,
				Session:                  session,
			}).
			Execute()
		if err != nil {
			return err
		}

		http.Redirect(w, r, oauth2RedirectTo.GetRedirectTo(), http.StatusFound)
		return nil
	} else if submit == "reject" {
		oauth2RedirectTo, _, err := m.oryClient.OAuth2API.RejectOAuth2ConsentRequest(m.oryAuthedContext).
			ConsentChallenge(consentChallenge).
			RejectOAuth2Request(ory.RejectOAuth2Request{
				Error:            ory.PtrString("access_denied"),
				ErrorDescription: ory.PtrString("The resource owner denied the request"),
			}).
			Execute()
		if err != nil {
			return err
		}
		http.Redirect(w, r, oauth2RedirectTo.GetRedirectTo(), http.StatusFound)
		return nil
	}

	return nil
}

func (m *Middleware) handleConsent(w http.ResponseWriter, r *http.Request, consentChallenge string) (*ory.OAuth2ConsentRequest, bool, error) {
	consentRequest, _, err := m.oryClient.OAuth2API.GetOAuth2ConsentRequest(m.oryAuthedContext).
		ConsentChallenge(consentChallenge).
		Execute()
	if err != nil {
		return nil, false, fmt.Errorf("could not retrieve consent request: %w", err)
	}

	if consentRequest.GetSkip() || consentRequest.Client.GetSkipConsent() {
		idToken := make(map[string]interface{})
		session := &ory.AcceptOAuth2ConsentRequestSession{
			IdToken: idToken,
		}
		err = m.handler.PopulateClaims(consentRequest, idToken)
		if err != nil {
			return nil, false, fmt.Errorf("could not populate claims: %w", err)
		}

		oauth2RedirectTo, _, err := m.oryClient.OAuth2API.AcceptOAuth2ConsentRequest(m.oryAuthedContext).
			ConsentChallenge(consentChallenge).
			AcceptOAuth2ConsentRequest(ory.AcceptOAuth2ConsentRequest{
				GrantScope:               consentRequest.RequestedScope,
				GrantAccessTokenAudience: consentRequest.RequestedAccessTokenAudience,
				Session:                  session,
			}).Execute()
		if err != nil {
			return nil, false, fmt.Errorf("could not accept consent request: %w", err)
		}

		http.Redirect(w, r, oauth2RedirectTo.GetRedirectTo(), http.StatusFound)
		return consentRequest, true, nil
	}

	return consentRequest, false, nil
}
