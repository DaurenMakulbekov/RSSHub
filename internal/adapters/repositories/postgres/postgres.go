package postgres

import (
	"database/sql"
	"fmt"
	"os"

	"RSSHub/internal/infrastructure/config"
)

type postgresRepository struct {
	db *sql.DB
}

func NewPostgresRepository(config *config.DB) *postgresRepository {
	url := fmt.Sprintf("user=%s password=%s host=%s port=%s database=%s sslmode=disable",
		config.User, config.Password, config.Host, config.Port, config.Name,
	)

	db, err := sql.Open("pgx", url)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to create connection pool: %v\n", err)
	}

	return &postgresRepository{
		db: db,
	}
}
