package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/babydevbr/golang/internal/app"
	"github.com/babydevbr/golang/internal/domain"
	"github.com/babydevbr/golang/internal/infra/gateway"
	"github.com/babydevbr/golang/internal/infra/handler"
	"github.com/babydevbr/golang/internal/infra/repository"
	"github.com/babydevbr/golang/pkg/env"
	"github.com/babydevbr/golang/pkg/krypto"
	"github.com/babydevbr/golang/pkg/middleware"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	"github.com/shaj13/go-guardian/v2/auth/strategies/jwt"
	"github.com/urfave/negroni"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	log.Println("Start API")

	dbURL := env.Get("DATABASE_URL", "host=localhost port=5432 user=postgres dbname=postgres password=postgres sslmode=disable")

	db, err := provideDB(dbURL, true)
	if err != nil {
		log.Panicf("database step %s", err.Error())
	}

	userRepo := repository.NeWUserGorm(db)
	userUseCase := app.NewUserUseCase(userRepo, &krypto.Hash{}, validator.New())

	// go-guardian
	keeper := jwt.StaticSecret{
		ID:        env.Get("SECRET_ID", "secret-id"),
		Secret:    []byte(env.Get("SECRET", "secret")),
		Algorithm: jwt.HS256,
	}

	guard := gateway.NewGuardian(keeper, userUseCase, &krypto.Hash{})

	r := mux.NewRouter()
	n := negroni.New(
		negroni.NewLogger(),
		negroni.HandlerFunc(middleware.Cors),
		negroni.HandlerFunc(guard.AuthMiddleware),
	)

	userHandler := handler.NewUserRestHandler(userUseCase, guard)
	userHandler.SetUserRoutes(r.PathPrefix("/users").Subrouter(), *n)

	r.HandleFunc("/", handlerHi)
	http.Handle("/", r)

	err = http.ListenAndServe(":"+env.Get("PORT", "5001"), r)
	if err != nil {
		log.Fatal(err.Error())
	}
}

func handlerHi(w http.ResponseWriter, r *http.Request) {
	msg := "Hello Baby Dev"
	log.Println(msg)
	_, _ = w.Write([]byte(msg))
}

func provideDB(dbURL string, migrate bool) (*gorm.DB, error) {
	dialector := postgres.Open(dbURL)

	db, err := gorm.Open(dialector, &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect database err: %w", err)
	}

	if !migrate {
		return db, nil
	}

	if dbc := db.Exec(`CREATE EXTENSION IF NOT EXISTS "uuid-ossp"`); dbc.Error != nil {
		return nil, fmt.Errorf("failed to migrate database err: %w", err)
	}

	err = db.AutoMigrate(domain.User{})
	if err != nil {
		return nil, fmt.Errorf("failed to migrate database err: %w", err)
	}

	return db, nil
}
