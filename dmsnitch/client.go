package dmsnitch

import (
	"context"
	"errors"
	"net/http"

	"github.com/suzuki-shunsuke/go-httpclient/httpclient"
)

type Client struct {
	Client httpclient.Client
}

func NewClient(apiKey string) Client {
	client := httpclient.New("https://api.deadmanssnitch.com")
	client.SetRequest = func(req *http.Request) error {
		req.SetBasicAuth(apiKey, ":")
		req.Header.Add("Content-Type", "application/json")
		req.Close = true
		return nil
	}
	return Client{
		Client: client,
	}
}

var ErrTokenIsRequired = errors.New("token is required")

func (client Client) Get(ctx context.Context, token string) (Snitch, *http.Response, error) {
	if token == "" {
		return Snitch{}, nil, ErrTokenIsRequired
	}
	snitch := Snitch{}
	resp, err := client.Client.Call(ctx, httpclient.CallParams{
		Method:       "GET",
		Path:         "/v1/snitches/" + token,
		ResponseBody: &snitch,
	})
	return snitch, resp, err
}

func (client Client) Post(ctx context.Context, snitch Snitch) (Snitch, *http.Response, error) {
	responseBody := Snitch{}
	resp, err := client.Client.Call(ctx, httpclient.CallParams{
		Method:       "POST",
		Path:         "/v1/snitches",
		RequestBody:  snitch,
		ResponseBody: &responseBody,
	})
	return responseBody, resp, err
}

func (client Client) Patch(ctx context.Context, token string, snitch Snitch) (*http.Response, error) {
	if token == "" {
		return nil, ErrTokenIsRequired
	}
	resp, err := client.Client.Call(ctx, httpclient.CallParams{
		Method:      "PATCH",
		Path:        "/v1/snitches/" + token,
		RequestBody: snitch,
	})
	return resp, err
}

func (client Client) Delete(ctx context.Context, token string) (*http.Response, error) {
	if token == "" {
		return nil, ErrTokenIsRequired
	}
	resp, err := client.Client.Call(ctx, httpclient.CallParams{
		Method: "DELETE",
		Path:   "/v1/snitches/" + token,
	})
	return resp, err
}
