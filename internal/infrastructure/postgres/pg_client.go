package postgres

import (
	"errors"
	"fmt"
	"github.com/payment-gateway/internal/config"
	"github.com/payment-gateway/internal/domain/model"
	"github.com/payment-gateway/internal/domain/port"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
	"time"
)

type PgClient struct {
	gormDB *gorm.DB
}

func (pg *PgClient) GetSessionById(sessionID string) (*model.Session, *model.DomainError) {
	var sessionEntity SessionEntity
	if err := pg.gormDB.Preload("TransactionEntities").Where("id = ?", sessionID).First(&sessionEntity).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, &model.DomainError{
				Category: model.FunctionalError,
				Type:     model.SessionNotFound,
				Message:  "Session not found",
			}
		}
		return nil, &model.DomainError{Category: model.TechnicalError, RootCause: err, Message: "error when searching for sessionEntity: " + sessionID}
	}

	session := sessionEntity.toModel()
	return &session, nil
}

func (pg *PgClient) IsSessionUnique(sessionID string) (bool, *model.DomainError) {
	var count int64
	if err := pg.gormDB.Model(&SessionEntity{}).Where("id = ?", sessionID).Count(&count).Error; err != nil {
		return false, &model.DomainError{Category: model.TechnicalError, RootCause: err, Message: "error when checking is session unique"}
	}
	return count == 0, nil
}

func (pg *PgClient) CreateSession(session model.Session) *model.DomainError {
	sessionEntity := mapToSessionEntity(session)
	if err := pg.gormDB.Create(&sessionEntity).Error; err != nil {
		return &model.DomainError{
			Category:  model.TechnicalError,
			RootCause: err,
			Message:   "failed to create a new session",
		}
	}
	return nil
}

func (pg *PgClient) UpdateSession(session model.Session) *model.DomainError {
	sessionEntity := mapToSessionEntity(session)
	if err := pg.gormDB.Updates(&sessionEntity).Error; err != nil {
		return &model.DomainError{
			Category:  model.TechnicalError,
			RootCause: err,
			Message:   "failed to create a new session",
		}
	}
	return nil
}

func NewPostgresClient(config config.PostgresConfig) (port.IPaymentRepository, error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=%s",
		config.Host, config.User, config.Password, config.Database, config.Port, config.SSLMode)
	var gormDB *gorm.DB
	var err error
	// According to the env file, activate or not the database logs
	if !config.LogQueries {
		gormDB, err = gorm.Open(postgres.Open(dsn))
	} else {
		// Enable logs
		logConfig := logger.New(
			log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
			logger.Config{
				SlowThreshold:             time.Second,
				LogLevel:                  logger.Info,
				IgnoreRecordNotFoundError: false,
				ParameterizedQueries:      false,
				Colorful:                  true,
			},
		)
		gormDB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{Logger: logConfig})
	}
	if err != nil {
		return nil, err
	}

	// Upload table schema
	err = gormDB.Debug().AutoMigrate(&SessionEntity{}, &TransactionEntity{})
	if err != nil {
		return nil, err
	}
	return &PgClient{gormDB}, nil
}
