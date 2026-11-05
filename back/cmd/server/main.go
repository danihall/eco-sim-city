package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"git.sr.ht/~danihall/foss-eco-simcity-hex-back/internal/adapters/http/handlers"
	"git.sr.ht/~danihall/foss-eco-simcity-hex-back/internal/adapters/postgres"
	"git.sr.ht/~danihall/foss-eco-simcity-hex-back/internal/app"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/joho/godotenv/autoload"
	_ "go.uber.org/automaxprocs"
)

func main() {
	slog_handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelError})
	slog.SetDefault(slog.New(slog_handler))

	ctx := context.Background()

	db_params := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable pool_max_conns=48 pool_min_conns=10",
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

	app := app.Boot(store)

	mux := http.NewServeMux()
	mux.Handle("POST /city", handlers.CreateFullCity(app))
	mux.Handle("DELETE /city/{id}", handlers.DeleteCity(app))

	err = http.ListenAndServe(app_port, mux)
	if err != nil {
		slog.Error("could not launch server")
		panic(err)
	}
}
