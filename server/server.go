package server

import (
	"errors"
	_ "log"
	"net/http"
	"net/url"
	"strings"
)

type Server interface {
	ListenAndServe(*http.ServeMux) error
	Address() string
}

func NewServer(u *url.URL, args ...interface{}) (Server, error) {

	var svr Server
	var err error

	scheme := u.Scheme

	switch strings.ToUpper(scheme) {

	case "HTTP":

		svr, err = NewHTTPServer(u, args...)

	case "LAMBDA":

		svr, err = NewLambdaServer(u, args...)

	default:
		return nil, errors.New("Invalid server protocol")
	}

	if err != nil {
		return nil, err
	}

	return svr, nil
}
