package middleware

import (
	"fmt"
	"github.com/fastbill/go-httperrors"
	"github.com/gorilla/csrf"
	ory "github.com/ory/client-go"
	"html/template"
	"net/http"
)

type ConsentTemplateData struct {
	Provider          ProviderInfo
	Request           *ory.OAuth2ConsentRequest
	CSRFToken         string
	CSRFTemplateField template.HTML
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

	err = m.renderPage(w, "consent.gohtml",
		ConsentTemplateData{
			Provider: ProviderInfo{
				Name: m.cfg.Name,
			},
			Request:           consentRequest,
			CSRFToken:         csrf.Token(r),
			CSRFTemplateField: csrf.TemplateField(r),
		},
	)
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
		err = m.provider.PopulateClaims(consentRequest, idToken)
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
		err = m.provider.PopulateClaims(consentRequest, idToken)
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
