package xapi_test

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/unchainio/pkg/xapi"
)

const (
	baseURL     = "https://api.example.com/"
	baseURLPath = "/api-v3"
)

// setup sets up a test HTTP server along with a github.Client that is
// configured to talk to that test server. Tests should register handlers on
// mux which provide mock responses for the API method being tested.
func setup(tb testing.TB) (client *xapi.Client, mux *http.ServeMux, serverURL string, teardown func()) {
	// mux is the HTTP request multiplexer used with the test server.
	mux = http.NewServeMux()

	// We want to ensure that tests catch mistakes where the endpoint URL is
	// specified as absolute rather than relative. It only makes a difference
	// when there's a non-empty base URL path. So, use that. See issue #752.
	apiHandler := http.NewServeMux()
	apiHandler.Handle(baseURLPath+"/", http.StripPrefix(baseURLPath, mux))
	apiHandler.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		fmt.Fprintln(os.Stderr, "FAIL: Client.BaseURL path prefix is not preserved in the request URL:")
		fmt.Fprintln(os.Stderr)
		fmt.Fprintln(os.Stderr, "\t"+req.URL.String())
		fmt.Fprintln(os.Stderr)
		fmt.Fprintln(os.Stderr, "\tDid you accidentally use an absolute endpoint URL rather than relative?")
		fmt.Fprintln(os.Stderr, "\tSee https://github.com/google/go-github/issues/752 for information.")
		http.Error(w, "Client.BaseURL path prefix is not preserved in the request URL.", http.StatusInternalServerError)
	})

	// server is a test HTTP server used to provide mock API responses.
	server := httptest.NewServer(apiHandler)

	// client is the GitHub client being tested and is
	// configured to use test server.
	client, err := xapi.NewClient(server.URL+baseURLPath+"/", nil)
	assert.NoError(tb, err)

	testURL, err := url.Parse(server.URL + baseURLPath + "/")
	assert.NoError(tb, err)

	client.BaseURL = testURL

	return client, mux, server.URL, server.Close
}

func testMethod(t *testing.T, r *http.Request, want string) {
	if got := r.Method; got != want {
		t.Errorf("Request method: %v, want %v", got, want)
	}
}

func TestNewClient(t *testing.T) {
	baseURL := "https://custom-url/"

	c, err := xapi.NewClient(baseURL, nil)
	assert.NoError(t, err)

	if got, want := c.BaseURL.String(), baseURL; got != want {
		t.Errorf("NewClient BaseURL is %v, want %v", got, want)
	}
}

func TestBasicAuthTransport(t *testing.T) {
	client, mux, serverURL, teardown := setup(t)
	defer teardown()

	username, password, otp := "u", "p", "123456"

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		u, p, ok := r.BasicAuth()
		if !ok {
			t.Errorf("request does not contain basic auth credentials")
		}
		if u != username {
			t.Errorf("request contained basic auth username %q, want %q", u, username)
		}
		if p != password {
			t.Errorf("request contained basic auth password %q, want %q", p, password)
		}
	})

	tp := &xapi.BasicAuthTransport{
		Username: username,
		Password: password,
		OTP:      otp,
	}
	basicAuthClient, err := xapi.NewClient(serverURL+baseURLPath+"/", tp.Client())
	assert.NoError(t, err)

	basicAuthClient.BaseURL = client.BaseURL
	req, _ := basicAuthClient.NewRequest("GET", ".", nil)
	basicAuthClient.Do(context.Background(), req, nil)
}

func TestBasicAuthTransport_transport(t *testing.T) {
	// default GetTransport
	tp := &xapi.BasicAuthTransport{}
	if tp.GetTransport() != http.DefaultTransport {
		t.Errorf("Expected http.DefaultTransport to be used.")
	}

	// custom GetTransport
	tp = &xapi.BasicAuthTransport{
		Transport: &http.Transport{},
	}
	if tp.GetTransport() == http.DefaultTransport {
		t.Errorf("Expected custom Transport to be used.")
	}
}

func TestNewRequest(t *testing.T) {
	c, err := xapi.NewClient(baseURL, nil)
	assert.NoError(t, err)

	inURL, outURL := "/foo", baseURL+"foo"
	inBody, outBody := map[string]string{
		"field1": "val1",
	}, `{"field1":"val1"}`+"\n"
	req, err := c.NewRequest("GET", inURL, inBody)
	assert.NoError(t, err)

	// test that relative URL was expanded
	if got, want := req.URL.String(), outURL; got != want {
		t.Errorf("NewRequest(%q) URL is %v, want %v", inURL, got, want)
	}

	// test that body was JSON encoded
	body, _ := ioutil.ReadAll(req.Body)
	if got, want := string(body), outBody; got != want {
		t.Errorf("NewRequest(%q) Body is %v, want %v", inBody, got, want)
	}

	// test that default user-agent is attached to the request
	if got, want := req.Header.Get("User-Agent"), c.UserAgent; got != want {
		t.Errorf("NewRequest() User-Agent is %v, want %v", got, want)
	}
}

func TestDo(t *testing.T) {
	client, mux, _, teardown := setup(t)
	defer teardown()

	type foo struct {
		A string
	}

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `{"A":"a"}`)
	})

	req, err := client.NewRequest("GET", ".", nil)
	assert.NoError(t, err)

	body := new(foo)
	_, err = client.Do(context.Background(), req, body)
	assert.NoError(t, err)

	want := &foo{"a"}
	if !reflect.DeepEqual(body, want) {
		t.Errorf("Response body = %v, want %v", body, want)
	}
}

func TestDo_httpError(t *testing.T) {
	client, mux, _, teardown := setup(t)
	defer teardown()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Bad Request", 400)
	})

	req, _ := client.NewRequest("GET", ".", nil)
	resp, err := client.Do(context.Background(), req, nil)

	if err == nil {
		t.Fatal("Expected HTTP 400 error, got no error.")
	}
	if resp.StatusCode != 400 {
		t.Errorf("Expected HTTP 400 error, got %d status code.", resp.StatusCode)
	}
}

// Test handling of an error caused by the internal http client's Do()
// function. A redirect loop is pretty unlikely to occur within the GitHub
// API, but does allow us to exercise the right code path.
func TestDo_redirectLoop(t *testing.T) {
	client, mux, _, teardown := setup(t)
	defer teardown()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, baseURLPath, http.StatusFound)
	})

	req, _ := client.NewRequest("GET", ".", nil)
	_, err := client.Do(context.Background(), req, nil)

	if err == nil {
		t.Error("Expected error to be returned.")
	}
	if err, ok := err.(*url.Error); !ok {
		t.Errorf("Expected a URL error; got %#v.", err)
	}
}

func TestDo_noContent(t *testing.T) {
	client, mux, _, teardown := setup(t)
	defer teardown()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	})

	var body json.RawMessage

	req, _ := client.NewRequest("GET", ".", nil)
	_, err := client.Do(context.Background(), req, &body)
	if err != nil {
		t.Fatalf("Do returned unexpected error: %v", err)
	}
}

func TestSanitizeURL(t *testing.T) {
	tests := []struct {
		in, want string
	}{
		{"/?a=b", "/?a=b"},
		{"/?a=b&client_secret=secret", "/?a=b&client_secret=REDACTED"},
		{"/?a=b&client_id=id&client_secret=secret", "/?a=b&client_id=id&client_secret=REDACTED"},
	}

	for _, tt := range tests {
		inURL, _ := url.Parse(tt.in)
		want, _ := url.Parse(tt.want)

		if got := xapi.SanitizeURL(inURL); !reflect.DeepEqual(got, want) {
			t.Errorf("SanitizeURL(%v) returned %v, want %v", tt.in, got, want)
		}
	}
}
