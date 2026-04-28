package database

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
)

// Config contém as configurações de banco de dados
type Config struct {
	Driver           string
	ConnectionString string
	MaxOpenConns     int
	MaxIdleConns     int
	MaxConnLifetime  time.Duration
}

// NewDB cria uma nova conexão com o banco de dados
func NewDB(cfg *Config) (*sql.DB, error) {
	db, err := sql.Open(cfg.Driver, cfg.ConnectionString)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Configurar pool de conexões
	if cfg.MaxOpenConns > 0 {
		db.SetMaxOpenConns(cfg.MaxOpenConns)
	}
	if cfg.MaxIdleConns > 0 {
		db.SetMaxIdleConns(cfg.MaxIdleConns)
	}
	if cfg.MaxConnLifetime > 0 {
		db.SetConnMaxLifetime(cfg.MaxConnLifetime)
	}

	// Testar conexão
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return db, nil
}

// RunMigrations executa as migrações do banco de dados
func RunMigrations(db *sql.DB) error {
	migrations := []string{
		createUsersTable,
		createTransactionsTable,
		createNotificationsTable,
		createBlockchainRecordsTable,
		createAuditLogsTable,
	}

	for _, migration := range migrations {
		if _, err := db.Exec(migration); err != nil {
			return fmt.Errorf("failed to run migration: %w", err)
		}
	}

	return nil
}

const (
	createUsersTable = `
	CREATE TABLE IF NOT EXISTS users (
		id TEXT PRIMARY KEY,
		email TEXT UNIQUE NOT NULL,
		username TEXT UNIQUE NOT NULL,
		password TEXT NOT NULL,
		role TEXT NOT NULL,
		is_active BOOLEAN DEFAULT true,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		last_login TIMESTAMP
	);
	`

	createTransactionsTable = `
	CREATE TABLE IF NOT EXISTS transactions (
		id TEXT PRIMARY KEY,
		client_id TEXT NOT NULL,
		source_account TEXT NOT NULL,
		destination_account TEXT NOT NULL,
		amount DECIMAL(15, 2) NOT NULL,
		currency TEXT DEFAULT 'BRL',
		description TEXT,
		status TEXT DEFAULT 'PENDING',
		reason TEXT,
		retry_count INTEGER DEFAULT 0,
		max_retries INTEGER DEFAULT 3,
		blockchain_hash TEXT,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		processed_at TIMESTAMP,
		FOREIGN KEY (client_id) REFERENCES users(id)
	);
	CREATE INDEX IF NOT EXISTS idx_transactions_client_id ON transactions(client_id);
	CREATE INDEX IF NOT EXISTS idx_transactions_status ON transactions(status);
	`

	createNotificationsTable = `
	CREATE TABLE IF NOT EXISTS notifications (
		id TEXT PRIMARY KEY,
		client_id TEXT NOT NULL,
		type TEXT NOT NULL,
		title TEXT NOT NULL,
		message TEXT NOT NULL,
		status TEXT DEFAULT 'PENDING',
		channel TEXT DEFAULT 'email',
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		sent_at TIMESTAMP,
		failure_reason TEXT,
		FOREIGN KEY (client_id) REFERENCES users(id)
	);
	CREATE INDEX IF NOT EXISTS idx_notifications_client_id ON notifications(client_id);
	CREATE INDEX IF NOT EXISTS idx_notifications_status ON notifications(status);
	`

	createBlockchainRecordsTable = `
	CREATE TABLE IF NOT EXISTS blockchain_records (
		id TEXT PRIMARY KEY,
		transaction_id TEXT NOT NULL UNIQUE,
		block_index INTEGER NOT NULL,
		block_hash TEXT NOT NULL,
		previous_hash TEXT NOT NULL,
		timestamp TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		status TEXT DEFAULT 'RECORDED',
		transaction_data TEXT,
		FOREIGN KEY (transaction_id) REFERENCES transactions(id)
	);
	CREATE INDEX IF NOT EXISTS idx_blockchain_records_tx_id ON blockchain_records(transaction_id);
	`

	createAuditLogsTable = `
	CREATE TABLE IF NOT EXISTS audit_logs (
		id TEXT PRIMARY KEY,
		transaction_id TEXT,
		action TEXT NOT NULL,
		actor_id TEXT NOT NULL,
		actor_role TEXT NOT NULL,
		ip_address TEXT,
		timestamp TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		status TEXT DEFAULT 'SUCCESS',
		error_message TEXT,
		FOREIGN KEY (transaction_id) REFERENCES transactions(id),
		FOREIGN KEY (actor_id) REFERENCES users(id)
	);
	CREATE INDEX IF NOT EXISTS idx_audit_logs_tx_id ON audit_logs(transaction_id);
	CREATE INDEX IF NOT EXISTS idx_audit_logs_actor_id ON audit_logs(actor_id);
	`
)
