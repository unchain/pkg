package xapi

import (
	"net/http"
	"net/url"

	"github.com/pkg/errors"
)

var UserAgent = "xapi"

type Client struct {
	BaseURL   *url.URL
	Client    *http.Client
	UserAgent string
}

// NewClient returns a new generic API client. If a nil httpClient is
// provided, http.DefaultClient will be used.
func NewClient(baseURLString string, client *http.Client) (*Client, error) {
	baseURL, err := url.Parse(baseURLString)

	if err != nil {
		return nil, errors.Wrap(err, "Could not parse url from the config")
	}

	if client == nil {
		client = http.DefaultClient
	}

	return &Client{
		BaseURL:   baseURL,
		Client:    client,
		UserAgent: UserAgent,
	}, nil
}
