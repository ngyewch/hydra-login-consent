package middleware

import (
	"bytes"
	"net/http"
)

func (m *Middleware) renderPage(w http.ResponseWriter, templateName string, templateData any) error {
	buf := bytes.NewBuffer(nil)
	err := m.templates.Lookup(templateName).
		Execute(buf, templateData)
	if err != nil {
		return err
	}
	_, err = w.Write(buf.Bytes())
	if err != nil {
		return err
	}
	return nil
}
