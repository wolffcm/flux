package http_test

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/wolffcm/flux"
	_ "github.com/wolffcm/flux/builtin"
	"github.com/wolffcm/flux/codes"
	"github.com/wolffcm/flux/dependencies/url"
	"github.com/wolffcm/flux/internal/errors"
	"github.com/wolffcm/flux/semantic"
	"github.com/wolffcm/flux/values"
)

func addFail(scope values.Scope) {
	scope.Set("fail", values.NewFunction(
		"fail",
		semantic.NewFunctionPolyType(semantic.FunctionPolySignature{
			Return: semantic.Bool,
		}),
		func(ctx context.Context, args values.Object) (values.Value, error) {
			return nil, errors.New(codes.Aborted, "fail")
		},
		false,
	))
}

func TestGet(t *testing.T) {
	var req *http.Request

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, request *http.Request) {
		req = request
		var err error
		if err != nil {
			w.WriteHeader(500)
			w.Write([]byte(err.Error()))
		}
		w.WriteHeader(204)
	}))
	defer ts.Close()

	script := fmt.Sprintf(`
import "experimental/http"

status = http.get(url:"%s/path/a/b/c", headers: {x:"a",y:"b",z:"c"})
`, ts.URL)

	ctx := flux.NewDefaultDependencies().Inject(context.Background())
	if _, _, err := flux.Eval(ctx, script, addFail); err != nil {
		t.Fatal("evaluation of http.get failed: ", err)
	}
	if want, got := "/path/a/b/c", req.URL.Path; want != got {
		t.Errorf("unexpected url want: %q got: %q", want, got)
	}
	if want, got := "GET", req.Method; want != got {
		t.Errorf("unexpected method want: %q got: %q", want, got)
	}
	header := make(http.Header)
	header.Set("x", "a")
	header.Set("y", "b")
	header.Set("z", "c")
	header.Set("Accept-Encoding", "gzip")
	header.Set("User-Agent", "Go-http-client/1.1")
	if !cmp.Equal(header, req.Header) {
		t.Errorf("unexpected header -want/+got\n%s", cmp.Diff(header, req.Header))
	}

}

func TestGet_ValidationFail(t *testing.T) {
	script := `
import "experimental/http"

http.get(url:"http://127.1.1.1/path/a/b/c", headers: {x:"a",y:"b",z:"c"})
`

	deps := flux.NewDefaultDependencies()
	deps.Deps.HTTPClient = http.DefaultClient
	deps.Deps.URLValidator = url.PrivateIPValidator{}
	ctx := deps.Inject(context.Background())
	_, _, err := flux.Eval(ctx, script, addFail)
	if err == nil {
		t.Fatal("expected failure")
	}
	if !strings.Contains(err.Error(), "url is not valid") {
		t.Errorf("unexpected cause of failure, got err: %v", err)
	}
}

func TestGet_Timeout(t *testing.T) {

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, request *http.Request) {
		var err error
		if err != nil {
			w.WriteHeader(500)
			w.Write([]byte(err.Error()))
		}
		// Sleep for 1s
		time.Sleep(time.Second)
		w.WriteHeader(204)
	}))
	defer ts.Close()

	script := fmt.Sprintf(`
import "experimental/http"

resp = http.get(url:"%s/path/a/b/c", timeout: 10ms)
`, ts.URL)

	ctx := flux.NewDefaultDependencies().Inject(context.Background())
	_, _, err := flux.Eval(ctx, script, addFail)
	if err == nil {
		t.Fatal("expected timeout failure")
	}
	if !strings.Contains(err.Error(), "context deadline exceeded") {
		t.Errorf("unexpected cause of failure, got err: %v", err)
	}
}
