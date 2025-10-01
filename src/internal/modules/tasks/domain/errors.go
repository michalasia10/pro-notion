package domain

import "errors"

// ErrTaskNotFound is returned when a task cannot be located in persistence.
var (
	ErrTaskNotFound            = errors.New("task not found")
	ErrInvalidProjectID        = errors.New("invalid project ID")
	ErrInvalidNotionDatabaseID = errors.New("invalid notion database ID")
	ErrInvalidNotionPageID     = errors.New("invalid notion page ID")
	ErrInvalidIDGenerator      = errors.New("invalid ID generator")
)
