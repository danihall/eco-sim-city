package postgres

import (
	"context"
	"embed"
	"fmt"
	"io/fs"

	"git.sr.ht/~danihall/foss-eco-simcity-hex-back/internal/adapters/postgres/sqlc"
	"git.sr.ht/~danihall/foss-eco-simcity-hex-back/internal/domain/core"
	"git.sr.ht/~danihall/foss-eco-simcity-hex-back/internal/domain/core/registry"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/tern/v2/migrate"
)

//go:embed .migrations/*.sql
var migrations_fs embed.FS

type Store struct {
	pool    *pgxpool.Pool
	queries *sqlc.Queries
}

func New(pool *pgxpool.Pool) *Store {
	return &Store{
		pool:    pool,
		queries: sqlc.New(pool),
	}
}

func (s *Store) InitDatabase(ctx context.Context) error {
	err := runMigrations(ctx, s.pool)
	if err != nil {
		return fmt.Errorf("migrations, %w", err)
	}

	transaction, err := s.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("transaction: %w", err)
	}
	defer transaction.Rollback(ctx)

	queries_with_tx := s.queries.WithTx(transaction)

	err = runSeedings(ctx, queries_with_tx)
	if err != nil {
		return fmt.Errorf("seeding: %w", err)
	}

	err = populateCitiesCurrentState(ctx, queries_with_tx)
	if err != nil {
		return fmt.Errorf("populating cities latest state: %w", err)
	}

	err = transaction.Commit(ctx)
	if err != nil {
		return err
	}

	return nil
}

func runMigrations(ctx context.Context, pool *pgxpool.Pool) error {
	connection, err := pool.Acquire(ctx)
	if err != nil {
		return err
	}
	defer connection.Release()

	migrator, err := migrate.NewMigrator(ctx, connection.Conn(), "schema_version")
	if err != nil {
		return err
	}

	migrations_folder, err := fs.Sub(migrations_fs, ".migrations")
	if err != nil {
		return err
	}

	err = migrator.LoadMigrations(migrations_folder)
	if err != nil {
		return err
	}

	return migrator.Migrate(ctx)
}

func runSeedings(
	ctx context.Context,
	queries *sqlc.Queries,
) error {
	for super_kind, kinds := range registry.SuperKinds {
		err := queries.UpsertSuperKind(ctx, sqlc.UpsertSuperKindParams{
			ID:   super_kind,
			Name: super_kind.String(),
		})
		if err != nil {
			return fmt.Errorf("super_kinds: %w", err)
		}

		for _, kind := range kinds {
			err = queries.UpsertKind(ctx, sqlc.UpsertKindParams{
				ID:          kind,
				SuperKindID: super_kind,
				Name:        kind.String(),
			})
			if err != nil {
				return fmt.Errorf("kinds: %w", err)
			}
		}
	}

	for artifact, relations := range registry.ArtifactsRelations {
		err := queries.UpsertArtifact(ctx, sqlc.UpsertArtifactParams{
			ID:          artifact,
			SuperKindID: relations.SuperKind,
			KindID:      relations.Kind,
			Name:        artifact.String(),
		})
		if err != nil {
			return fmt.Errorf("artifacts: %w", err)
		}
	}

	for ordinance := range registry.Ordinances {
		err := queries.UpsertOrdinance(ctx, sqlc.UpsertOrdinanceParams{
			ID:   ordinance,
			Name: ordinance.String(),
		})
		if err != nil {
			return fmt.Errorf("ordinances: %w", err)
		}
	}

	return nil
}

func populateCitiesCurrentState(ctx context.Context, queries *sqlc.Queries) error {
	cities_latest_state, err := queries.GetAllCitiesLatestHistoryEntry(ctx)
	if err != nil {
		return err
	}

	if len(cities_latest_state) == 0 {
		return nil
	}

	params := make([]sqlc.PopulateCitiesCurrentStateParams, len(cities_latest_state))
	for index := range cities_latest_state {
		city_state := &cities_latest_state[index]

		params[index] = sqlc.PopulateCitiesCurrentStateParams{
			CityID:                 city_state.CityID,
			GameTick:               pgtype.Int8{Valid: true, Int64: city_state.GameTick},
			ElectricityConsumption: city_state.ElectricityConsumption,
			WaterConsumption:       city_state.WaterConsumption,
			Ghg1Level:              city_state.Ghg1Level,
			Ghg2Level:              city_state.Ghg2Level,
			Ghg3Level:              city_state.Ghg3Level,
			ChildrenCount:          city_state.ChildrenCount,
			AdultsCount:            city_state.AdultsCount,
			ElderlyCount:           city_state.ElderlyCount,
			ChildrenEducationLevel: city_state.ChildrenEducationLevel,
			AdultsEducationLevel:   city_state.AdultsEducationLevel,
			ElderlyEducationLevel:  city_state.ElderlyEducationLevel,
			ChildrenHealthLevel:    city_state.ChildrenHealthLevel,
			AdultsHealthLevel:      city_state.AdultsHealthLevel,
			ElderlyHealthLevel:     city_state.ElderlyHealthLevel,
		}
	}

	err = queries.ClearCitiesCurrentState(ctx)
	if err != nil {
		return err
	}
	_, err = queries.PopulateCitiesCurrentState(ctx, params)
	if err != nil {
		return err
	}

	return nil
}

// q returns queries bound to the current transaction if one is present in ctx,
// or to the plain pool otherwise. Use for queries with no locking requirement.
func (s *Store) q(ctx context.Context) *sqlc.Queries {
	transaction, ok := ExtractTaskTransaction(ctx)
	if ok {
		return s.queries.WithTx(transaction)
	}

	return s.queries
}

func makeQueryCoordinates(coordinates []core.TilePosition) ([]int16, []int16) {
	qs := make([]int16, len(coordinates))
	rs := make([]int16, len(coordinates))

	for index := range coordinates {
		qs[index] = int16(coordinates[index].Q)
		rs[index] = int16(coordinates[index].R)
	}

	return qs, rs
}

// byTilePosition is to be used as a callback function in slices.SortFunc.
//
// see postgres/cadastre.go
func byTilePosition[T core.PositionGetter](a, b T) int {
	tile_a := a.Position()
	tile_b := b.Position()

	if tile_a.Q != tile_b.Q {
		return int(tile_a.Q) - int(tile_b.Q)
	}
	return int(tile_a.R) - int(tile_b.R)
}
