package api

import "net/http"

type Options struct {
	Password string

	LoginURL string
}
type APIIface interface {
	DoGetRequest(requestURL string) (response, error)
}
type API struct {
	Options Options
	Client  http.Client
}

func New(option Options) APIIface {
	return API{
		Options: option,
		Client: http.Client{
			Transport: &MyJWTTransport{
				transport: http.DefaultTransport,
				password:  option.Password,
				loginURL:  option.LoginURL,
			},
		},
	}
}
