package domain

import (
	"fmt"

	"github.com/google/uuid"
)

// IDGenerator provides unique ID generation for domain entities
type IDGenerator interface {
	NewID(prefix string) string
}

// UUIDGenerator implements IDGenerator using UUID v4
type UUIDGenerator struct{}

func NewUUIDGenerator() IDGenerator {
	return &UUIDGenerator{}
}

func (g *UUIDGenerator) NewID(prefix string) string {
	id := uuid.New().String()
	if prefix != "" {
		return fmt.Sprintf("%s_%s", prefix, id)
	}
	return id
}
