package client

import (
	"context"
	"net/http"

	hc "github.com/circleci/backplane-go-x/x/httpclient"
)

type Client struct {
	client *hc.Client
}

func NewClient(ctx context.Context, baseURL string) *Client {
	return &Client{
		client: hc.New(ctx, hc.Config{
			Name:       "my-client",
			BaseURL:    baseURL,
			AcceptType: hc.JSON,
		}),
	}
}

type FooParams struct {
	Name string `json:"name"`
}

func (c *Client) Foo(ctx context.Context, id string, params FooParams) error {
	return c.client.Call(ctx, hc.NewRequest(http.MethodPost, "/foo/%s",
		hc.RouteParams(id),
		hc.Body(params),
	))
}
