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
	CreatePortfolioProfileTx(ctx context.Context, arg CreatePortfolioProfileTxParams) (CreatePortfolioProfileTxResult, error)
	UpdatePortfolioProfileTx(ctx context.Context, arg UpdatePortfolioProfileTxParams) (UpdatePortfolioProfileTxResult, error)
	DeletePortfolioProfileTx(ctx context.Context, arg DeletePortfolioProfileTxParams) (DeletePortfolioProfileTxResult, error)
}

// SQLStore provides all functions to execute SQL queries and transactions
type SQLStore struct {
	// All individual query functions provided by Queries will be available to Store
	*Queries
	// Extend more transactions
	connPool *pgxpool.Pool

	// Add logger
	logger *zap.Logger

	secretKeyEncryption string
}

// NewStore creates a new store
func NewStore(connPool *pgxpool.Pool, logger *zap.Logger, secretKeyEncryption string) Store {
	return &SQLStore{
		connPool:            connPool,
		Queries:             New(connPool),
		logger:              logger,
		secretKeyEncryption: secretKeyEncryption,
	}
}

func InitializeUpDB(databaseConfig util.DatabaseConfig, logger *zap.Logger, secretKeyEncryption string) (Store, func(), error) {
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
	store := NewStore(connPool, logger, secretKeyEncryption)

	cleanup := func() {
		connPool.Close()
	}

	return store, cleanup, nil
}
