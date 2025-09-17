package notion

import (
	"fmt"
)

// Pages provides methods for working with Notion pages
type Pages struct {
	client *Client
}

// NewPages creates a new Pages client
func NewPages(opts ...ClientOption) *Pages {
	return &Pages{
		client: NewClient(opts...),
	}
}

// Create creates a new page
func (p *Pages) Create(accessToken string, request *CreatePageRequest) (*Page, error) {
	resp, err := p.client.makeRequest("POST", "/pages", request, accessToken, authTypeBearer)
	if err != nil {
		return nil, fmt.Errorf("failed to create page: %w", err)
	}

	var page Page
	if err := p.client.handleResponse(resp, &page); err != nil {
		return nil, fmt.Errorf("failed to parse create page response: %w", err)
	}

	return &page, nil
}

// Retrieve retrieves a page by ID
func (p *Pages) Retrieve(accessToken, pageID string) (*Page, error) {
	endpoint := fmt.Sprintf("/pages/%s", pageID)

	resp, err := p.client.makeRequest("GET", endpoint, nil, accessToken, authTypeBearer)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve page: %w", err)
	}

	var page Page
	if err := p.client.handleResponse(resp, &page); err != nil {
		return nil, fmt.Errorf("failed to parse page response: %w", err)
	}

	return &page, nil
}

// Update updates an existing page
func (p *Pages) Update(accessToken, pageID string, request *UpdatePageRequest) (*Page, error) {
	endpoint := fmt.Sprintf("/pages/%s", pageID)

	resp, err := p.client.makeRequest("PATCH", endpoint, request, accessToken, authTypeBearer)
	if err != nil {
		return nil, fmt.Errorf("failed to update page: %w", err)
	}

	var page Page
	if err := p.client.handleResponse(resp, &page); err != nil {
		return nil, fmt.Errorf("failed to parse update page response: %w", err)
	}

	return &page, nil
}

// GetPropertyItem retrieves a specific property item from a page
func (p *Pages) GetPropertyItem(accessToken, pageID, propertyID string) (*PropertyValue, error) {
	endpoint := fmt.Sprintf("/pages/%s/properties/%s", pageID, propertyID)

	resp, err := p.client.makeRequest("GET", endpoint, nil, accessToken, authTypeBearer)
	if err != nil {
		return nil, fmt.Errorf("failed to get property item: %w", err)
	}

	var property PropertyValue
	if err := p.client.handleResponse(resp, &property); err != nil {
		return nil, fmt.Errorf("failed to parse property response: %w", err)
	}

	return &property, nil
}
