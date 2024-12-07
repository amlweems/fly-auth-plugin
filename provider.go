package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"strconv"
	"time"
)

const (
	tokenAddr    = "/.fly/api"
	tokenURL     = "http://localhost/v1/tokens/oidc"
	jwtTokenType = "urn:ietf:params:oauth:token-type:jwt"
)

func boolptr(v bool) *bool { return &v }

// Token matches the executable response expected by
// cloud.google.com/go/auth/credentials/externalaccount.
type TokenResponse struct {
	Version        int    `json:"version,omitempty"`
	Success        *bool  `json:"success,omitempty"`
	TokenType      string `json:"token_type,omitempty"`
	ExpirationTime int64  `json:"expiration_time,omitempty"`
	IDToken        string `json:"id_token,omitempty"`
	Code           string `json:"code,omitempty"`
	Message        string `json:"message,omitempty"`
}

// TokenProvider fetches tokens.
type TokenProvider interface {
	Token(ctx context.Context, aud string) TokenResponse
}

// Provider is a externalaccount.TokenProvider used to fetch fly.io OIDC tokens.
type Provider struct {
	client *http.Client
}

// NewProvider creates a provider.
func NewProvider() TokenProvider {
	return &Provider{
		client: &http.Client{
			Transport: &http.Transport{
				DialContext: func(ctx context.Context, _, _ string) (net.Conn, error) {
					var d net.Dialer
					return d.DialContext(ctx, "unix", tokenAddr)
				},
			},
		},
	}
}

// Token fetches an OIDC token from the local fly.io unix domain socket.
func (stp *Provider) Token(ctx context.Context, aud string) (r TokenResponse) {
	r.Version = 1
	r.Success = boolptr(false)
	r.TokenType = jwtTokenType
	r.ExpirationTime = time.Now().Add(600 * time.Second).Unix()

	body, _ := json.Marshal(map[string]string{"aud": aud})
	req, _ := http.NewRequestWithContext(ctx, http.MethodPost, tokenURL, bytes.NewReader(body))

	resp, err := stp.client.Do(req)
	if err != nil {
		r.Message = fmt.Sprintf("failed to request token: %s", err)
		r.Code = "400"
		return
	}

	if resp.StatusCode != http.StatusOK {
		responseBody, _ := io.ReadAll(resp.Body)
		r.Code = strconv.Itoa(resp.StatusCode)
		r.Message = string(responseBody)
		return
	}

	token, err := io.ReadAll(resp.Body)
	if err != nil {
		r.Message = fmt.Sprintf("failed to read token: %s", err)
		r.Code = "400"
		return
	}

	r.IDToken = string(token)
	r.Success = boolptr(true)
	return
}
