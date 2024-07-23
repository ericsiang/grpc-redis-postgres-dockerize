package db

import (
	"database/sql"
	"fmt"
	"grpc-redis-postgres/proto"

	_ "github.com/lib/pq"
)

type DatabaseClient interface {
	GetUser(id string) (*proto.User, error)
	CreateUser(user *proto.User) (*proto.User, error)
}

type databaseClient struct {
	db *sql.DB
}

// NewDatabaseClient creates a new database client.
func NewDatabaseClient(dsn string) (DatabaseClient, error) {
	fmt.Println("dsn:", dsn)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("error connecting to PostgreSQL: %w", err)
	}
	// if err := db.Ping(); err != nil {
	// 	return nil, fmt.Errorf("error pinging PostgreSQL: %w", err)
	// }
	return &databaseClient{db: db}, nil
}

// GetUser retrieves a user from the database.
func (c *databaseClient) GetUser(id string) (*proto.User, error) {
	var user proto.User
	err := c.db.QueryRow("SELECT id, name, email FROM users WHERE id = ?", id).Scan(&user.Id, &user.Name, &user.Email)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// CreateUser creates a new user in the database.
func (c *databaseClient) CreateUser(user *proto.User) (*proto.User, error) {
	var lastInsertId int64 = 0
	err := c.db.QueryRow("INSERT INTO users (name, email) VALUES ($1, $2) RETURNING id", user.Name, user.Email).Scan(&lastInsertId)
	if err != nil {
		return nil, err
	}

	user.Id = lastInsertId
	return user, nil
}
