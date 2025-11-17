package middlewares

import (
	"fmt"
	"net/http"
	"strings"
)

type HPPOptions struct {
	CheckQuery               bool
	CheckBody                bool
	CheckBodyForContentType  string
	Whitelist                []string
}

func HPP(options HPPOptions) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			if options.CheckBody && r.Method == http.MethodPost && isCorrectContentType(r, options.CheckBodyForContentType) {
				// Trata os parametros para evitar duplicações/sobrecarga
				filterBodyParams(r, options.Whitelist)
			}

			if options.CheckQuery && r.URL.Query() != nil {
				// Trata os parametros para evitar duplicações/sobrecarga
				filterQueryParams(r, options.Whitelist)
			}

			next.ServeHTTP(w, r)
		})
	}
}


func isCorrectContentType(r *http.Request, contentType string) bool {
	return strings.Contains(r.Header.Get("Content-Type"), contentType)
}

func filterBodyParams(r *http.Request, whitelist []string) {
	err := r.ParseForm()
	if err != nil {
		fmt.Println("Erro ao parsear body:", err)
		return
	}

	for key, values := range r.Form {
		if len(values) > 1 {
			r.Form.Set(key, values[0])
		}

		if !isWhitelisted(key, whitelist) {
			delete(r.Form, key)
		}
	}
}

func filterQueryParams(r *http.Request, whitelist []string) {
	query := r.URL.Query()

	for key, values := range query {
		if len(values) > 1 {
			query.Set(key, values[0])
		}

		if !isWhitelisted(key, whitelist) {
			query.Del(key)
		}
	}
	r.URL.RawQuery = query.Encode()
}

func isWhitelisted(param string, whitelist []string) bool {
	for _, allowed := range whitelist {
		if param == allowed {
			return true
		}
	}
	return false
}
