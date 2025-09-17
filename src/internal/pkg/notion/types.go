package notion

import (
	"fmt"
	"time"
)

// APIError represents an error returned by the Notion API
type APIError struct {
	Object           string `json:"object"`
	Status           int    `json:"status"`
	Code             string `json:"code"`
	Message          string `json:"message"`
	APIErrorType     string `json:"error"`
	ErrorDescription string `json:"error_description"`
}

func (e *APIError) Error() string {
	// Prefer OAuth error format if available
	if e.APIErrorType != "" {
		return fmt.Sprintf("Notion API error (%s): %s", e.APIErrorType, e.ErrorDescription)
	}
	return fmt.Sprintf("Notion API error (%s): %s", e.Code, e.Message)
}

// OAuth Types

// OAuthTokenRequest represents the request to exchange code for token
type OAuthTokenRequest struct {
	GrantType   string `json:"grant_type"`
	Code        string `json:"code"`
	RedirectURI string `json:"redirect_uri"`
}

// OAuthTokenResponse represents the OAuth token response
type OAuthTokenResponse struct {
	AccessToken   string `json:"access_token"`
	TokenType     string `json:"token_type"`
	BotID         string `json:"bot_id"`
	WorkspaceID   string `json:"workspace_id"`
	WorkspaceName string `json:"workspace_name"`
	WorkspaceIcon string `json:"workspace_icon"`
	Owner         User   `json:"owner"`
}

// User represents a Notion user
type User struct {
	Object    string  `json:"object"`
	ID        string  `json:"id"`
	Type      string  `json:"type"`
	Name      string  `json:"name"`
	AvatarURL *string `json:"avatar_url"`
	Person    *Person `json:"person,omitempty"`
	Bot       *Bot    `json:"bot,omitempty"`
}

// Person represents person details for a user
type Person struct {
	Email string `json:"email"`
}

// Bot represents bot-specific information
type Bot struct {
	Owner BotOwner `json:"owner"`
}

// BotOwner represents the owner of a bot
type BotOwner struct {
	Type string `json:"type"`
	User User   `json:"user"`
}

// Database Types

// Database represents a Notion database
type Database struct {
	Object         string              `json:"object"`
	ID             string              `json:"id"`
	CreatedTime    time.Time           `json:"created_time"`
	LastEditedTime time.Time           `json:"last_edited_time"`
	Title          []RichText          `json:"title"`
	Description    []RichText          `json:"description"`
	Icon           *Icon               `json:"icon"`
	Cover          *Cover              `json:"cover"`
	Properties     map[string]Property `json:"properties"`
	Parent         Parent              `json:"parent"`
	URL            string              `json:"url"`
	Archived       bool                `json:"archived"`
	IsInline       bool                `json:"is_inline"`
}

// DatabaseQueryRequest represents a database query request
type DatabaseQueryRequest struct {
	Filter      *Filter `json:"filter,omitempty"`
	Sorts       []Sort  `json:"sorts,omitempty"`
	StartCursor string  `json:"start_cursor,omitempty"`
	PageSize    int     `json:"page_size,omitempty"`
}

// DatabaseQueryResponse represents a database query response
type DatabaseQueryResponse struct {
	Object     string `json:"object"`
	Results    []Page `json:"results"`
	NextCursor string `json:"next_cursor"`
	HasMore    bool   `json:"has_more"`
	Type       string `json:"type"`
}

// Page Types

// Page represents a Notion page
type Page struct {
	Object         string                   `json:"object"`
	ID             string                   `json:"id"`
	CreatedTime    time.Time                `json:"created_time"`
	LastEditedTime time.Time                `json:"last_edited_time"`
	CreatedBy      User                     `json:"created_by"`
	LastEditedBy   User                     `json:"last_edited_by"`
	Cover          *Cover                   `json:"cover"`
	Icon           *Icon                    `json:"icon"`
	Parent         Parent                   `json:"parent"`
	Archived       bool                     `json:"archived"`
	Properties     map[string]PropertyValue `json:"properties"`
	URL            string                   `json:"url"`
}

// CreatePageRequest represents a request to create a new page
type CreatePageRequest struct {
	Parent     Parent                   `json:"parent"`
	Properties map[string]PropertyValue `json:"properties"`
	Children   []Block                  `json:"children,omitempty"`
	Icon       *Icon                    `json:"icon,omitempty"`
	Cover      *Cover                   `json:"cover,omitempty"`
}

// UpdatePageRequest represents a request to update a page
type UpdatePageRequest struct {
	Properties map[string]PropertyValue `json:"properties"`
	Archived   *bool                    `json:"archived,omitempty"`
	Icon       *Icon                    `json:"icon,omitempty"`
	Cover      *Cover                   `json:"cover,omitempty"`
}

// Common Types

// Parent represents a page or database parent
type Parent struct {
	Type       string `json:"type"`
	PageID     string `json:"page_id,omitempty"`
	DatabaseID string `json:"database_id,omitempty"`
	Workspace  bool   `json:"workspace,omitempty"`
}

// Property represents a database property
type Property struct {
	ID       string    `json:"id"`
	Type     string    `json:"type"`
	Name     string    `json:"name"`
	Title    *struct{} `json:"title,omitempty"`
	RichText *struct{} `json:"rich_text,omitempty"`
	Number   *struct {
		Format string `json:"format"`
	} `json:"number,omitempty"`
	Select *struct {
		Options []SelectOption `json:"options"`
	} `json:"select,omitempty"`
	MultiSelect *struct {
		Options []SelectOption `json:"options"`
	} `json:"multi_select,omitempty"`
	Date        *struct{} `json:"date,omitempty"`
	People      *struct{} `json:"people,omitempty"`
	Files       *struct{} `json:"files,omitempty"`
	Checkbox    *struct{} `json:"checkbox,omitempty"`
	URL         *struct{} `json:"url,omitempty"`
	Email       *struct{} `json:"email,omitempty"`
	PhoneNumber *struct{} `json:"phone_number,omitempty"`
	Formula     *struct {
		Expression string `json:"expression"`
	} `json:"formula,omitempty"`
	Relation *struct {
		DatabaseID         string `json:"database_id"`
		SyncedPropertyName string `json:"synced_property_name,omitempty"`
		SyncedPropertyID   string `json:"synced_property_id,omitempty"`
	} `json:"relation,omitempty"`
	Rollup *struct {
		RelationPropertyName string `json:"relation_property_name"`
		RelationPropertyID   string `json:"relation_property_id"`
		RollupPropertyName   string `json:"rollup_property_name"`
		RollupPropertyID     string `json:"rollup_property_id"`
		Function             string `json:"function"`
	} `json:"rollup,omitempty"`
	CreatedTime    *struct{} `json:"created_time,omitempty"`
	CreatedBy      *struct{} `json:"created_by,omitempty"`
	LastEditedTime *struct{} `json:"last_edited_time,omitempty"`
	LastEditedBy   *struct{} `json:"last_edited_by,omitempty"`
}

// PropertyValue represents a property value on a page
type PropertyValue struct {
	ID   string `json:"id,omitempty"`
	Type string `json:"type"`

	// Value types
	Title          []RichText     `json:"title,omitempty"`
	RichText       []RichText     `json:"rich_text,omitempty"`
	Number         *float64       `json:"number,omitempty"`
	Select         *SelectOption  `json:"select,omitempty"`
	MultiSelect    []SelectOption `json:"multi_select,omitempty"`
	Date           *DateValue     `json:"date,omitempty"`
	People         []User         `json:"people,omitempty"`
	Files          []File         `json:"files,omitempty"`
	Checkbox       *bool          `json:"checkbox,omitempty"`
	URL            *string        `json:"url,omitempty"`
	Email          *string        `json:"email,omitempty"`
	PhoneNumber    *string        `json:"phone_number,omitempty"`
	Formula        *FormulaValue  `json:"formula,omitempty"`
	Relation       []Relation     `json:"relation,omitempty"`
	Rollup         *RollupValue   `json:"rollup,omitempty"`
	CreatedTime    *time.Time     `json:"created_time,omitempty"`
	CreatedBy      *User          `json:"created_by,omitempty"`
	LastEditedTime *time.Time     `json:"last_edited_time,omitempty"`
	LastEditedBy   *User          `json:"last_edited_by,omitempty"`
}

// RichText represents rich text content
type RichText struct {
	Type string `json:"type"`
	Text *struct {
		Content string `json:"content"`
		Link    *Link  `json:"link"`
	} `json:"text,omitempty"`
	Mention     *Mention    `json:"mention,omitempty"`
	Equation    *Equation   `json:"equation,omitempty"`
	Annotations Annotations `json:"annotations"`
	PlainText   string      `json:"plain_text"`
	Href        string      `json:"href,omitempty"`
}

// Annotations represents text formatting
type Annotations struct {
	Bold          bool   `json:"bold"`
	Italic        bool   `json:"italic"`
	Strikethrough bool   `json:"strikethrough"`
	Underline     bool   `json:"underline"`
	Code          bool   `json:"code"`
	Color         string `json:"color"`
}

// Link represents a hyperlink
type Link struct {
	URL string `json:"url"`
}

// Mention represents a mention in rich text
type Mention struct {
	Type string `json:"type"`
	User *User  `json:"user,omitempty"`
	Page *struct {
		ID string `json:"id"`
	} `json:"page,omitempty"`
	Database *struct {
		ID string `json:"id"`
	} `json:"database,omitempty"`
	Date *DateValue `json:"date,omitempty"`
}

// Equation represents an equation in rich text
type Equation struct {
	Expression string `json:"expression"`
}

// SelectOption represents an option in a select property
type SelectOption struct {
	ID    string `json:"id,omitempty"`
	Name  string `json:"name"`
	Color string `json:"color,omitempty"`
}

// DateValue represents a date property value
type DateValue struct {
	Start    string  `json:"start"`
	End      *string `json:"end"`
	TimeZone *string `json:"time_zone"`
}

// File represents a file in a files property
type File struct {
	Name string `json:"name"`
	Type string `json:"type"`
	File *struct {
		URL        string     `json:"url"`
		ExpiryTime *time.Time `json:"expiry_time"`
	} `json:"file,omitempty"`
	External *struct {
		URL string `json:"url"`
	} `json:"external,omitempty"`
}

// FormulaValue represents a computed formula value
type FormulaValue struct {
	Type    string     `json:"type"`
	String  *string    `json:"string,omitempty"`
	Number  *float64   `json:"number,omitempty"`
	Boolean *bool      `json:"boolean,omitempty"`
	Date    *DateValue `json:"date,omitempty"`
}

// Relation represents a relation to another page
type Relation struct {
	ID string `json:"id"`
}

// RollupValue represents a rollup property value
type RollupValue struct {
	Type   string        `json:"type"`
	Number *float64      `json:"number,omitempty"`
	Date   *DateValue    `json:"date,omitempty"`
	Array  []interface{} `json:"array,omitempty"`
}

// Icon represents an icon (emoji or file)
type Icon struct {
	Type     string `json:"type"`
	Emoji    string `json:"emoji,omitempty"`
	File     *File  `json:"file,omitempty"`
	External *struct {
		URL string `json:"url"`
	} `json:"external,omitempty"`
}

// Cover represents a cover image
type Cover struct {
	Type     string `json:"type"`
	File     *File  `json:"file,omitempty"`
	External *struct {
		URL string `json:"url"`
	} `json:"external,omitempty"`
}

// Block represents a content block
type Block struct {
	Object string `json:"object"`
	Type   string `json:"type"`
	// Block-specific content would be added here
}

// Filter represents query filters
type Filter struct {
	Property string `json:"property,omitempty"`
	// Filter conditions would be added here based on property type
	And []Filter `json:"and,omitempty"`
	Or  []Filter `json:"or,omitempty"`
}

// Sort represents query sorting
type Sort struct {
	Property  string `json:"property,omitempty"`
	Timestamp string `json:"timestamp,omitempty"`
	Direction string `json:"direction"`
}
