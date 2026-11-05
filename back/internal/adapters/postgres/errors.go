package postgres

import (
	"errors"
	"fmt"
	"strings"

	"git.sr.ht/~danihall/foss-eco-simcity-hex-back/internal/domain"
	"git.sr.ht/~danihall/foss-eco-simcity-hex-back/internal/domain/api"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

var errors_map = map[string]error{
	"mutable.cities_23505": api.ErrEMailConflict,
}

func getPostgresError(err error) error {
	var pg_error *pgconn.PgError

	switch {
	case errors.As(err, &pg_error):
		var s strings.Builder
		s.Grow(len(pg_error.SchemaName) + len(pg_error.TableName) + 8)
		s.WriteString(pg_error.SchemaName)
		s.WriteString(".")
		s.WriteString(pg_error.TableName)
		s.WriteString("_")
		s.WriteString(pg_error.Code)

		effective_error, ok := errors_map[s.String()]
		if ok {
			return effective_error
		}

	case errors.Is(err, pgx.ErrNoRows):
		return fmt.Errorf("%w: %v", domain.ErrEmptyResult, pgx.ErrNoRows)
	}

	return fmt.Errorf("postgres: %w", err)
}
