package transaction

import (
	"context"

	"gorm.io/gorm"

	"src/internal/pkg/txctx"
)

// Manager wraps a GORM DB to provide TransactionManager behavior.
type Manager struct {
	db *gorm.DB
}

// NewManager creates a new transaction manager backed by GORM.
func NewManager(db *gorm.DB) *Manager {
	return &Manager{db: db}
}

// WithinTransaction executes fn within a database transaction and passes the transactional DB via context.
func (m *Manager) WithinTransaction(ctx context.Context, fn func(ctx context.Context) error) error {
	return m.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		ctxWithTx := txctx.WithTx(ctx, tx)
		return fn(ctxWithTx)
	})
}
