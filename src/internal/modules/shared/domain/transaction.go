package domain

import "context"

// TransactionManager provides transaction support for use cases
type TransactionManager interface {
	WithinTransaction(ctx context.Context, fn func(ctx context.Context) error) error
}

// NoopTransactionManager implements TransactionManager without actual transactions
// Useful for testing or when transactions are not needed
type NoopTransactionManager struct{}

func NewNoopTransactionManager() TransactionManager {
	return &NoopTransactionManager{}
}

func (m *NoopTransactionManager) WithinTransaction(ctx context.Context, fn func(ctx context.Context) error) error {
	return fn(ctx)
}
