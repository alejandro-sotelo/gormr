package gormr

import (
	"gorm.io/gorm"

	"github.com/alejandro-sotelo/gormr/internal/db"
	"github.com/alejandro-sotelo/gormr/internal/repository"
)

// Client is the main entry point for interacting with the gormr sdk.
// It holds the database connection and repository helpers.
type Client struct {
	db   *gorm.DB
	repo *repository.Repository
}

type DBConfig = db.DBConfig

// New creates a new SDK instance: it opens the DB connection and prepares helpers.
func New(cfg DBConfig) (*Client, error) {
	connection, err := db.Connect(cfg)
	if err != nil {
		return nil, err
	}
	return &Client{
		db:   connection,
		repo: repository.New(connection),
	}, nil
}

// Close closes the underlying sql.DB connection.
func (c *Client) Close() error {
	if c == nil || c.db == nil {
		return nil
	}
	sqlDB, err := c.db.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}

// DB returns the underlying *gorm.DB instance.
func (c *Client) DB() *gorm.DB {
	return c.db
}

// Repo returns the repository helper.
func (c *Client) Repo() *repository.Repository {
	return c.repo
}
