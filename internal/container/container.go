package container

import (
	"database/sql"

	"github.com/Jeielsantosdev/Work_processamento_assincrono/config"
	"github.com/Jeielsantosdev/Work_processamento_assincrono/internal/domain/interfaces"
	"github.com/Jeielsantosdev/Work_processamento_assincrono/internal/infra/auth"
	"github.com/Jeielsantosdev/Work_processamento_assincrono/internal/infra/blockchain"
	"github.com/Jeielsantosdev/Work_processamento_assincrono/internal/infra/crypto"
	"github.com/Jeielsantosdev/Work_processamento_assincrono/internal/infra/database"
	"github.com/Jeielsantosdev/Work_processamento_assincrono/internal/infra/notification"
	"github.com/Jeielsantosdev/Work_processamento_assincrono/internal/infra/queue"
	"github.com/Jeielsantosdev/Work_processamento_assincrono/internal/infra/repository"
	"github.com/Jeielsantosdev/Work_processamento_assincrono/internal/usecase"
	"github.com/Jeielsantosdev/Work_processamento_assincrono/pkg/logger"
)

// Container contém todas as dependências injetadas
type Container struct {
	Config config.Config
	DB     *sql.DB
	Logger logger.Logger

	// Repositories
	TransactionRepo  interfaces.TransactionRepository
	UserRepo         interfaces.UserRepository
	NotificationRepo interfaces.NotificationRepository
	BlockchainRepo   interfaces.BlockchainRepository
	AuditLogRepo     interfaces.AuditLogRepository

	// Services
	PasswordService     interfaces.PasswordService
	AuthService         interfaces.AuthenticationService
	BlockchainService   interfaces.BlockchainService
	NotificationService interfaces.NotificationService
	QueueService        interfaces.QueueService

	// Use Cases
	CreateTransactionUC     *usecase.CreateTransactionUseCase
	ProcessTransactionUC    *usecase.ProcessTransactionUseCase
	GetTransactionHistoryUC *usecase.GetTransactionHistoryUseCase
	AuthenticateUserUC      *usecase.AuthenticateUserUseCase
}

// New cria um novo container com todas as dependências
func New(cfg config.Config) (*Container, error) {
	// Logger
	log := logger.NewLogrusLogger(cfg.Logger.Level)

	// Database
	db, err := database.NewDB(&database.Config{
		Driver:           cfg.Database.Driver,
		ConnectionString: cfg.Database.ConnectionString,
		MaxOpenConns:     cfg.Database.MaxOpenConns,
		MaxIdleConns:     cfg.Database.MaxIdleConns,
		MaxConnLifetime:  cfg.Database.MaxConnLifetime,
	})
	if err != nil {
		return nil, err
	}

	// Run migrations
	if err := database.RunMigrations(db); err != nil {
		return nil, err
	}

	// Repositories
	transactionRepo := repository.NewTransactionRepositorySQL(db)
	userRepo := repository.NewUserRepositorySQL(db)
	notificationRepo := repository.NewNotificationRepositorySQL(db)
	blockchainRepo := repository.NewBlockchainRepositorySQL(db)
	auditLogRepo := repository.NewAuditLogRepositorySQL(db)

	// Services
	passwordService := crypto.NewBcryptPasswordService(0)
	authService := auth.NewJWTAuthService(cfg.JWT.SecretKey, cfg.JWT.ExpiresIn)
	blockchainService := blockchain.NewSimpleBlockchainService(blockchainRepo)
	notificationService := notification.NewConsoleNotificationService(notificationRepo)

	var queueService interfaces.QueueService
	switch cfg.Queue.Type {
	case "rabbitmq":
		queueService = queue.NewRabbitMQQueueService(cfg.Queue.ConnectionString, cfg.Queue.QueueName)
	default:
		queueService = queue.NewInMemoryQueueService(1000)
	}

	// Use Cases
	createTransactionUC := usecase.NewCreateTransactionUseCase(transactionRepo, queueService, auditLogRepo)
	processTransactionUC := usecase.NewProcessTransactionUseCase(transactionRepo, blockchainService, notificationService, auditLogRepo)
	getTransactionHistoryUC := usecase.NewGetTransactionHistoryUseCase(transactionRepo, auditLogRepo)
	authenticateUserUC := usecase.NewAuthenticateUserUseCase(userRepo, passwordService, authService)

	return &Container{
		Config:                  cfg,
		DB:                      db,
		Logger:                  log,
		TransactionRepo:         transactionRepo,
		UserRepo:                userRepo,
		NotificationRepo:        notificationRepo,
		BlockchainRepo:          blockchainRepo,
		AuditLogRepo:            auditLogRepo,
		PasswordService:         passwordService,
		AuthService:             authService,
		BlockchainService:       blockchainService,
		NotificationService:     notificationService,
		QueueService:            queueService,
		CreateTransactionUC:     createTransactionUC,
		ProcessTransactionUC:    processTransactionUC,
		GetTransactionHistoryUC: getTransactionHistoryUC,
		AuthenticateUserUC:      authenticateUserUC,
	}, nil
}

// Close fecha as conexões
func (c *Container) Close() error {
	return c.DB.Close()
}
