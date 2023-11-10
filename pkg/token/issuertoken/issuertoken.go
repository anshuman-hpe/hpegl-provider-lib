// (C) Copyright 2021-2023 Hewlett Packard Enterprise Development LP

package issuertoken

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"io"
	"net/http"
	"strings"
	"time"

	tokenutil "github.com/hewlettpackard/hpegl-provider-lib/pkg/token/token-util"
)

const (
	retryLimit = 3
)

type GenerateTokenInput struct {
	TenantID     string `json:"tenant_id"`
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	GrantType    string `json:"grant_type"`
}

type TokenResponse struct {
	TokenType       string    `json:"token_type"`
	AccessToken     string    `json:"access_token"`
	RefreshToken    string    `json:"refresh_token"`
	Expiry          time.Time `json:"expiry"`
	ExpiresIn       int       `json:"expires_in"`
	Scope           string    `json:"scope"`
	AccessTokenOnly bool      `json:"accessTokenOnly"`
}

func GenerateToken(
	ctx context.Context,
	tenantID,
	clientID,
	clientSecret string,
	identityServiceURL string,
	httpClient tokenutil.HttpClient,
) (string, error) {
	//params := url.Values{}
	//params.Add("client_id", clientID)
	//params.Add("client_secret", clientSecret)
	//params.Add("grant_type", "client_credentials")
	//params.Add("scope", "hpe-tenant")
	params := GenerateTokenInput{
		TenantID:     tenantID,
		ClientID:     clientID,
		ClientSecret: clientSecret,
		GrantType:    "client_credentials",
	}
	b, err := json.Marshal(params)
	if err != nil {
		return "", err
	}

	url := fmt.Sprintf("%s/identity/v1/token", identityServiceURL)
	tflog.Info(ctx, "anshuman, debug_logs")
	tflog.Info(ctx, string(b))
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, strings.NewReader(string(b)))
	//req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, strings.NewReader(params.Encode()))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := tokenutil.DoRetries(func() (*http.Response, error) {
		return httpClient.Do(req)
	}, retryLimit)
	if err != nil {
		return "", fmt.Errorf("network error in post to get token")
	}
	defer resp.Body.Close()

	err = tokenutil.ManageHTTPErrorCodes(resp, clientID)
	if err != nil {
		return "", err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var token TokenResponse

	err = json.Unmarshal(body, &token)
	if err != nil {
		return "", err
	}

	return token.AccessToken, nil
}
