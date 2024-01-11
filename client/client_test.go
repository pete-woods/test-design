package client_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/circleci/backplane-go-x/x/routerx"
	"github.com/circleci/backplane-go-x/x/testing/httprecorder"
	"github.com/circleci/backplane-go-x/x/testing/httprecorder/ginrecorder"
	"github.com/circleci/backplane-go/httpserver/router"
	"github.com/circleci/backplane-go/testing/testcontext"
	"github.com/gin-gonic/gin"
	"gotest.tools/v3/assert"
	"gotest.tools/v3/assert/cmp"

	"github.com/pete-woods/test-design/client"
)

func TestClient_Foo(t *testing.T) {
	ctx := testcontext.Background()

	rec := httprecorder.New()
	r := routerx.New(ctx, router.Config{})
	r.Use(ginrecorder.Middleware(ctx, rec))
	r.POST("/foo/:id", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "ok"})
	})

	srv := httptest.NewServer(r)
	t.Cleanup(srv.Close)

	c := client.NewClient(ctx, srv.URL)

	t.Run("Make request", func(t *testing.T) {
		err := c.Foo(ctx, "thingy", client.FooParams{Name: "the name"})
		assert.Assert(t, err)
	})

	t.Run("Check request format", func(t *testing.T) {
		body, err := json.Marshal(client.FooParams{Name: "the name"})
		assert.Assert(t, err)

		assert.Check(t, cmp.DeepEqual(rec.LastRequest(), &httprecorder.Request{
			Method: "POST",
			URL: url.URL{
				Path: "/foo/thingy",
			},
			Header: map[string][]string{
				"Content-Type": {"application/json; charset=utf-8"},
			},
			Body: append(body, '\n'),
		}, httprecorder.OnlyHeaders("Content-Type")))
	})

}
