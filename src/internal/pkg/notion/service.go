package notion

// Service provides a unified interface to the Notion API
type Service struct {
	OAuth     *OAuth
	Databases *Databases
	Pages     *Pages
	client    *Client
}

// ServiceConfig contains configuration for the Notion service
type ServiceConfig struct {
	ClientID     string
	ClientSecret string
	RedirectURI  string
	APIVersion   string
}

// NewService creates a new Notion API service
func NewService(config ServiceConfig, opts ...ClientOption) *Service {
	// Apply API version if provided
	if config.APIVersion != "" {
		opts = append(opts, WithAPIVersion(config.APIVersion))
	}

	return &Service{
		OAuth:     NewOAuth(config.ClientID, config.ClientSecret, config.RedirectURI, opts...),
		Databases: NewDatabases(opts...),
		Pages:     NewPages(opts...),
		client:    NewClient(opts...),
	}
}

// GetAuthorizationURL is a convenience method for OAuth authorization URL generation
func (s *Service) GetAuthorizationURL(state string) string {
	return s.OAuth.GetAuthorizationURL(state)
}

// ExchangeCodeForToken is a convenience method for OAuth token exchange
func (s *Service) ExchangeCodeForToken(code string) (*OAuthTokenResponse, error) {
	return s.OAuth.ExchangeCodeForToken(code)
}

// GetCurrentUser is a convenience method for getting current user
func (s *Service) GetCurrentUser(accessToken string) (*User, error) {
	return s.OAuth.GetCurrentUser(accessToken)
}

// QueryDatabase is a convenience method for querying databases
func (s *Service) QueryDatabase(accessToken, databaseID string, request *DatabaseQueryRequest) (*DatabaseQueryResponse, error) {
	return s.Databases.Query(accessToken, databaseID, request)
}

// CreatePage is a convenience method for creating pages
func (s *Service) CreatePage(accessToken string, request *CreatePageRequest) (*Page, error) {
	return s.Pages.Create(accessToken, request)
}

// UpdatePage is a convenience method for updating pages
func (s *Service) UpdatePage(accessToken, pageID string, request *UpdatePageRequest) (*Page, error) {
	return s.Pages.Update(accessToken, pageID, request)
}

// RetrievePage is a convenience method for retrieving pages
func (s *Service) RetrievePage(accessToken, pageID string) (*Page, error) {
	return s.Pages.Retrieve(accessToken, pageID)
}

// RetrieveDatabase is a convenience method for retrieving databases
func (s *Service) RetrieveDatabase(accessToken, databaseID string) (*Database, error) {
	return s.Databases.Retrieve(accessToken, databaseID)
}

// ListDatabases is a convenience method for listing databases
func (s *Service) ListDatabases(accessToken string, startCursor string, pageSize int) (*DatabaseQueryResponse, error) {
	return s.Databases.List(accessToken, startCursor, pageSize)
}
