package postgres

import (
	"context"
	"fmt"

	"git.sr.ht/~danihall/foss-eco-simcity-hex-back/internal/domain/cityhall"
	"github.com/jackc/pgx/v5"
)

var _ cityhall.TaskForceExecutor = (*Store)(nil)

type taskKey struct{}

// ExtractTaskTransaction retrieves the current transaction from ctx, if any.
func ExtractTaskTransaction(ctx context.Context) (pgx.Tx, bool) {
	tx, ok := ctx.Value(taskKey{}).(pgx.Tx)
	return tx, ok
}

func withTaskTransaction(ctx context.Context, tx pgx.Tx) context.Context {
	return context.WithValue(ctx, taskKey{}, tx)
}

func (s *Store) Task(ctx context.Context, action func(task_ctx context.Context) error) error {
	_, already_in_ctx := ExtractTaskTransaction(ctx)
	if already_in_ctx {
		return action(ctx)
	}

	transaction, err := s.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin tx: %w", err)
	}
	defer transaction.Rollback(ctx)

	task_ctx := withTaskTransaction(ctx, transaction)

	err = action(task_ctx)
	if err != nil {
		return err
	}

	err = transaction.Commit(ctx)
	if err != nil {
		return fmt.Errorf("failed to commit tx: %w", err)
	}

	return nil
}
