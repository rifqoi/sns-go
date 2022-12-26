package postgres

import (
	"context"
	"fmt"
	"log"
	"net/url"
	"os"

	"github.com/jackc/pgx/v4/pgxpool"
	snsgo "github.com/rifqoi/sns-go"
	"github.com/rifqoi/sns-go/postgres/dbsqlc"
)

func convertGender(g dbsqlc.Gender) snsgo.Gender {
	switch g {
	case dbsqlc.GenderMale:
		return snsgo.Male
	case dbsqlc.GenderFemale:
		return snsgo.Female
	}

	return snsgo.GenderUndefined
}

func newGender(g snsgo.Gender) dbsqlc.Gender {
	switch g {
	case snsgo.Male:
		return dbsqlc.GenderMale
	case snsgo.Female:
		return dbsqlc.GenderFemale
	}

	return dbsqlc.GenderMale
}

func ConnectPostgres() (*pgxpool.Pool, error) {
	dbHost := os.Getenv("DB_HOST")
	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASS")
	dbName := os.Getenv("DB_NAME")
	dbPort := os.Getenv("DB_PORT")

	dsn := url.URL{
		Scheme: "postgres",
		User:   url.UserPassword(dbUser, dbPass),
		Host:   fmt.Sprintf("%s:%s", dbHost, dbPort),
		Path:   dbName,
	}

	pool, err := pgxpool.Connect(context.Background(), dsn.String())
	if err != nil {
		log.Fatal(err)
	}

	err = pool.Ping(context.Background())
	if err != nil {
		log.Fatal("Cannot ping the database", err)
	}

	return pool, nil
}
