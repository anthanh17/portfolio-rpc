package db

import (
	"context"
	"fmt"
	"portfolio-profile-rpc/util"

	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

// Repository pattern
type Store interface {
	Querier
	// Extend more transactions
	CreatePortfolioTx(ctx context.Context, arg CreatePortfolioTxParams) (CreatePortfolioTxResult, error)
}

// SQLStore provides all functions to execute SQL queries and transactions
type SQLStore struct {
	// All individual query functions provided by Queries will be available to Store
	*Queries
	// Extend more transactions
	connPool *pgxpool.Pool
}

// NewStore creates a new store
func NewStore(connPool *pgxpool.Pool) Store {
	return &SQLStore{
		connPool: connPool,
		Queries:  New(connPool),
	}
}

func InitializeUpDB(databaseConfig util.DatabaseConfig, logger *zap.Logger) (Store, func(), error) {
	connectionString := fmt.Sprintf("postgresql://%s:%s@%s:%d/%s?sslmode=disable",
		databaseConfig.Username,
		databaseConfig.Password,
		databaseConfig.Host,
		databaseConfig.Port,
		databaseConfig.Database)

	// Connect database
	connPool, err := pgxpool.New(context.Background(), connectionString)
	if err != nil {
		logger.Info("cannot connect to db")
		return nil, nil, err
	}

	// Create database accessor
	store := NewStore(connPool)

	cleanup := func() {
		connPool.Close()
	}

	return store, cleanup, nil
}
