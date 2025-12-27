package txctx

import (
	"context"

	"gorm.io/gorm"
)

type txKey struct{}

// WithTx stores the current transaction in the context so downstream repos can reuse it.
func WithTx(ctx context.Context, tx *gorm.DB) context.Context {
	return context.WithValue(ctx, txKey{}, tx)
}

// FromContext extracts a transactional *gorm.DB if present.
func FromContext(ctx context.Context) *gorm.DB {
	if ctx == nil {
		return nil
	}
	tx, _ := ctx.Value(txKey{}).(*gorm.DB)
	return tx
}
