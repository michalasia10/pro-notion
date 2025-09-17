package notion

import (
	"encoding/base64"
	"fmt"
	"log"
	"net/url"
)

// OAuth provides OAuth 2.0 methods for Notion
type OAuth struct {
	client       *Client
	clientID     string
	clientSecret string
	redirectURI  string
}

// NewOAuth creates a new OAuth client
func NewOAuth(clientID, clientSecret, redirectURI string, opts ...ClientOption) *OAuth {
	return &OAuth{
		client:       NewClient(opts...),
		clientID:     clientID,
		clientSecret: clientSecret,
		redirectURI:  redirectURI,
	}
}

// GetAuthorizationURL generates the OAuth authorization URL
func (o *OAuth) GetAuthorizationURL(state string) string {
	baseURL := "https://api.notion.com/v1/oauth/authorize"

	params := url.Values{}
	params.Add("client_id", o.clientID)
	params.Add("response_type", "code")
	params.Add("owner", "user")
	params.Add("redirect_uri", o.redirectURI)
	if state != "" {
		params.Add("state", state)
	}
	return fmt.Sprintf("%s?%s", baseURL, params.Encode())
}

// ExchangeCodeForToken exchanges an authorization code for an access token
func (o *OAuth) ExchangeCodeForToken(code string) (*OAuthTokenResponse, error) {
	tokenReq := OAuthTokenRequest{
		GrantType:   "authorization_code",
		Code:        code,
		RedirectURI: o.redirectURI,
	}
	accessToken := base64.StdEncoding.EncodeToString(fmt.Appendf(nil, "%s:%s", o.clientID, o.clientSecret))

	resp, err := o.client.makeRequest("POST", "/oauth/token", tokenReq, accessToken, authTypeBasic)
	if err != nil {
		return nil, fmt.Errorf("failed to exchange code for token: %w", err)
	}

	var tokenResp OAuthTokenResponse
	log.Println("resp", resp.Body)
	if err := o.client.handleResponse(resp, &tokenResp); err != nil {
		log.Printf("failed to parse token response: %v", err)
		return nil, fmt.Errorf("failed to parse token response: %w", err)
	}

	return &tokenResp, nil
}

// GetCurrentUser retrieves the current user information
func (o *OAuth) GetCurrentUser(accessToken string) (*User, error) {
	resp, err := o.client.makeRequest("GET", "/users/me", nil, accessToken, authTypeBearer)
	if err != nil {
		return nil, fmt.Errorf("failed to get current user: %w", err)
	}

	var user User
	if err := o.client.handleResponse(resp, &user); err != nil {
		return nil, fmt.Errorf("failed to parse user response: %w", err)
	}

	return &user, nil
}
