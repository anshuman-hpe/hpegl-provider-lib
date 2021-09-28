package httpclient

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/hewlettpackard/hpegl-provider-lib/pkg/token/identitytoken"
	"github.com/hewlettpackard/hpegl-provider-lib/pkg/token/issuertoken"
	tokenutil "github.com/hewlettpackard/hpegl-provider-lib/pkg/token/token-util"
)

type Client struct {
	identityServiceURL  string
	httpClient          tokenutil.HttpClient
	vendedServiceClient bool
}

// New creates a new identity Client object
func New(identityServiceURL string, vendedServiceClient bool) *Client {
	client := &http.Client{Timeout: 10 * time.Second}
	identityServiceURL = strings.TrimRight(identityServiceURL, "/")
	return &Client{
		identityServiceURL:  identityServiceURL,
		httpClient:          client,
		vendedServiceClient: vendedServiceClient,
	}
}

func (c *Client) GenerateToken(ctx context.Context, tenantID, clientID, clientSecret string) (string, error) {
	if c.vendedServiceClient {
		token, err := issuertoken.GenerateToken(ctx, tenantID, clientID, clientSecret, c.identityServiceURL, c.httpClient)
		return token, err
	} else {
		token, err := identitytoken.GenerateToken(ctx, tenantID, clientID, clientSecret, c.identityServiceURL, c.httpClient)
		return token, err
	}
}