package postgres

import (
	"database/sql"
	"fmt"
	"os"

	"RSSHub/internal/core/domain"
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

func (postgresRepo *postgresRepository) AddFeed(feed domain.Feeds) error {
	query := `INSERT INTO feeds (name, url) VALUES($1, $2)`

	_, err := postgresRepo.db.Exec(query, feed.Name, feed.Url)
	if err != nil {
		return fmt.Errorf("Error: %v", err)
	}

	return nil
}
