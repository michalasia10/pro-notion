package notion

import (
	"fmt"
)

// Databases provides methods for working with Notion databases
type Databases struct {
	client *Client
}

// NewDatabases creates a new Databases client
func NewDatabases(opts ...ClientOption) *Databases {
	return &Databases{
		client: NewClient(opts...),
	}
}

// Query queries a database and returns matching pages
func (d *Databases) Query(accessToken, databaseID string, request *DatabaseQueryRequest) (*DatabaseQueryResponse, error) {
	endpoint := fmt.Sprintf("/databases/%s/query", databaseID)

	resp, err := d.client.makeRequest("POST", endpoint, request, accessToken)
	if err != nil {
		return nil, fmt.Errorf("failed to query database: %w", err)
	}

	var queryResp DatabaseQueryResponse
	if err := d.client.handleResponse(resp, &queryResp); err != nil {
		return nil, fmt.Errorf("failed to parse database query response: %w", err)
	}

	return &queryResp, nil
}

// Retrieve retrieves a database by ID
func (d *Databases) Retrieve(accessToken, databaseID string) (*Database, error) {
	endpoint := fmt.Sprintf("/databases/%s", databaseID)

	resp, err := d.client.makeRequest("GET", endpoint, nil, accessToken)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve database: %w", err)
	}

	var database Database
	if err := d.client.handleResponse(resp, &database); err != nil {
		return nil, fmt.Errorf("failed to parse database response: %w", err)
	}

	return &database, nil
}

// List lists all databases accessible to the integration
func (d *Databases) List(accessToken string, startCursor string, pageSize int) (*DatabaseQueryResponse, error) {
	endpoint := "/search"

	request := map[string]interface{}{
		"filter": map[string]string{
			"value":    "database",
			"property": "object",
		},
	}

	if startCursor != "" {
		request["start_cursor"] = startCursor
	}

	if pageSize > 0 {
		request["page_size"] = pageSize
	}

	resp, err := d.client.makeRequest("POST", endpoint, request, accessToken)
	if err != nil {
		return nil, fmt.Errorf("failed to list databases: %w", err)
	}

	var listResp DatabaseQueryResponse
	if err := d.client.handleResponse(resp, &listResp); err != nil {
		return nil, fmt.Errorf("failed to parse database list response: %w", err)
	}

	return &listResp, nil
}
