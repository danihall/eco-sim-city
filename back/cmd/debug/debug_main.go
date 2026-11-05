// This file allows to check & test internal database requests with `curlie` by setting fake endpoints.
package main

import (
	"context"
	"fmt"
	"log/slog"
	"mime"
	"net/http"
	"os"

	adaptershttp "git.sr.ht/~danihall/foss-eco-simcity-hex-back/internal/adapters/http"
	"git.sr.ht/~danihall/foss-eco-simcity-hex-back/internal/adapters/http/handlers"
	"git.sr.ht/~danihall/foss-eco-simcity-hex-back/internal/adapters/postgres"
	"git.sr.ht/~danihall/foss-eco-simcity-hex-back/internal/app"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/joho/godotenv/autoload"
	"github.com/lmittmann/tint"
)

func main() {
	slog_handler := slog.New(tint.NewHandler(os.Stdout, &tint.Options{
		Level:      slog.LevelDebug,
		TimeFormat: "15:04:05.000",
		AddSource:  true,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if a.Key == slog.LevelKey {
				level := a.Value.Any().(slog.Level)
				if level == slog.LevelDebug {
					return slog.String(slog.LevelKey, "\033[38;2;175;95;255mDBG\033[0m")
				}
			}
			return a
		},
	}))
	slog.SetDefault(slog_handler)

	mime.AddExtensionType(".js", "application/javascript")

	ctx := context.Background()

	db_params := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable pool_max_conns=40 pool_min_conns=10",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
	)
	app_port := fmt.Sprintf(":%s", os.Getenv("APP_PORT"))

	pool, err := pgxpool.New(ctx, db_params)
	if err != nil {
		slog.Error("could not initiate connection to database")
		panic(err)
	}
	defer pool.Close()

	err = pool.Ping(ctx)
	if err != nil {
		slog.Error("could not connect to database")
		panic(err)
	}

	store := postgres.New(pool)

	err = store.InitDatabase(ctx)
	if err != nil {
		panic(err)
	}

	slog.Info("created & populated tables")

	app := app.Boot(store)
	err = app.StartAllCitiesSimulation(ctx)
	if err != nil {
		panic(err)
	}

	mux := http.NewServeMux()

	// static routes and assets
	mux.Handle("GET /", handlers.PageIndex())
	mux.Handle("GET /assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("./assets/"))))
	mux.Handle("GET /public/", http.StripPrefix("/public/", http.FileServer(http.Dir("./public/"))))
	// API
	mux.Handle("POST /city", handlers.CreateFullCity(app))
	mux.Handle("DELETE /city/{id}", handlers.DeleteCity(app))
	mux.Handle("GET /city/{id}/all-tiles", handlers.GetCityTilesArtifacts(app))
	mux.Handle("PATCH /city/{id}/zoning/update", handlers.UpdateZone(app))
	mux.Handle("PATCH /city/{id}/zoning/revert", handlers.RevertZone(app))
	mux.Handle("PATCH /city/{id}/construction/build", handlers.BuildArtifact(app))
	mux.Handle("PATCH /city/{id}/construction/build-pathway", handlers.BuildPathway(app))
	mux.Handle("PATCH /city/{id}/construction/destroy", handlers.DestroyArtifact(app))
	mux.Handle("POST /city/{id}/ordinances/enact", handlers.EnactOrdinance(app))
	mux.Handle("DELETE /city/{id}/ordinances/revoke", handlers.RevokeOrdinance(app))
	mux.Handle("GET /city/{id}/traffic", handlers.GetCityTraffic(app))
	// SSE
	mux.Handle("GET /city/{id}/sse", handlers.ServerSideEvents(app))

	slog.Info("launching server", "port", app_port[1:])
	err = http.ListenAndServe(app_port, adaptershttp.StrictURLMiddleware(mux))
	if err != nil {
		slog.Error("could not launch server")
		panic(err)
	}
}

// debug
// mux.Handle("GET /debug/city/{id}/tiles", DebugGetTiles(queries))
// mux.Handle("GET /debug/city/{id}/tiles-by-anchor", DebugGetTilesByAnchor(queries))
// mux.Handle("GET /debug/city/{id}/quadrant-tiles", DebugGetQuadrantTiles(queries))
