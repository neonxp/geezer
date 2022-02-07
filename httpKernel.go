package geezer

import (
	"io"
	"net/http"
	"net/url"
	"strings"
)

type HttpKernel struct {
	AppKernel
}

func NewHttpKernel() *HttpKernel {
	return &HttpKernel{
		AppKernel: newKernel(),
	}
}

func (s *HttpKernel) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	u, err := url.ParseRequestURI(r.RequestURI)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	parts := strings.Split(strings.Trim(u.Path, "/"), "/")
	if len(parts) == 0 {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	name := parts[0]
	id := ""
	if len(parts) > 1 {
		id = parts[1]
	}
	method := MethodFind
	switch {
	case r.Method == http.MethodGet && id == "":
		method = MethodFind
	case r.Method == http.MethodGet && id != "":
		method = MethodGet
	case r.Method == http.MethodPost && id == "":
		method = MethodCreate
	case r.Method == http.MethodPost && id != "":
		method = MethodUpdate
	case r.Method == http.MethodPut:
		method = MethodUpdate
	case r.Method == http.MethodPatch:
		method = MethodPatch
	case r.Method == http.MethodDelete:
		method = MethodRemove
	default:
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	params := Params{
		Path:     parts,
		Query:    Values(u.Query()),
		Headers:  Values(r.Header),
		Provider: "http",
	}
	b, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	defer r.Body.Close()
	data := Data(b)

	result, err := s.Call(ctx, method, name, id, data, params)
	if err != nil {
		if err == ErrMethodNotFound {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		// TODO log
		return
	}

	w.Header().Set("Content-Type", result.ContentType())

	switch method {
	case MethodFind, MethodGet, MethodUpdate, MethodPatch:
		w.WriteHeader(http.StatusOK)
	case MethodCreate:
		w.WriteHeader(http.StatusCreated)
	case MethodRemove:
		w.WriteHeader(http.StatusNoContent)
	}
	if method == MethodRemove {
		return
	}
	if err := result.Render(w); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		// TODO log
	}
}
