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

func (postgresRepo *postgresRepository) GetFeeds() ([]domain.Feeds, error) {
	var feeds []domain.Feeds

	rows, err := postgresRepo.db.Query("SELECT * FROM feeds")
	if err != nil {
		return nil, fmt.Errorf("Error: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var feed domain.Feeds

		if err := rows.Scan(&feed.ID, &feed.Name, &feed.Url, &feed.Created, &feed.Updated); err != nil {
			return nil, fmt.Errorf("Error: %v", err)
		}
		feeds = append(feeds, feed)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("Error: %v", err)
	}

	return feeds, nil
}

func (postgresRepo *postgresRepository) WriteArticles(articles []domain.RSSItem, feed domain.Feeds) error {
	tx, err := postgresRepo.db.Begin()
	if err != nil {
		return fmt.Errorf("Error Transaction Begin: %v", err)
	}
	defer tx.Rollback()

	for i := range articles {
		query := `INSERT INTO articles (title, link, description, published_at, feed_id) VALUES($1, $2, $3, $4, $5)`

		_, err := tx.Exec(query, articles[i].Title, articles[i].Link, articles[i].Description, articles[i].PubDate, feed.ID)
		if err != nil {
			return fmt.Errorf("Error to write article: %v", err)
		}
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("Error Transaction Commit: %v", err)
	}

	return nil
}

func (postgresRepo *postgresRepository) GetArticles(feed domain.Feeds) ([]domain.Articles, error) {
	var articles []domain.Articles

	rows, err := postgresRepo.db.Query("SELECT title, published_at FROM articles WHERE feed_id = $1", feed.ID)
	if err != nil {
		return nil, fmt.Errorf("Error: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var article domain.Articles

		if err := rows.Scan(&article.Title, &article.Published); err != nil {
			return nil, fmt.Errorf("Error: %v", err)
		}
		articles = append(articles, article)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("Error: %v", err)
	}

	return articles, nil
}
