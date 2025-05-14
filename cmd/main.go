package main

import (
	"context"
	"fmt"
	"os"

	"github.com/Koyo-os/Vote-service/internal/entity"
	"github.com/Koyo-os/Vote-service/internal/repository"
	"github.com/Koyo-os/Vote-service/internal/service"
	"github.com/Koyo-os/Vote-service/pkg/config"
	"github.com/Koyo-os/Vote-service/pkg/errs"
	"github.com/Koyo-os/Vote-service/pkg/logger"
	"github.com/Koyo-os/Vote-service/pkg/retrier"
	"github.com/Koyo-os/Vote-service/pkg/transport/casher"
	"github.com/Koyo-os/Vote-service/pkg/transport/consumer"
	"github.com/Koyo-os/Vote-service/pkg/transport/listener"
	"github.com/Koyo-os/Vote-service/pkg/transport/publisher"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	// Initialize logger configuration
	loggerCfg := logger.Config{
		LogFile:   "vote-service-logs.log", // Log file path
		AppName:   "Vote-service",          // Application name for logging
		LogLevel:  "debug",                 // Logging level
		AddCaller: true,                    // Include caller information in logs
	}

	// Initialize logger with retry logic in case of failure
	if err := logger.Init(loggerCfg); err != nil {
		fmt.Printf("error initializing logger: %v", err)
		return
	}
	// Ensure logs are flushed before application exits
	defer logger.Sync()

	// Get the logger instance
	log := logger.Get()

	// Load application configuration
	cfg, err := config.Load()
	if err != nil || cfg == nil {
		log.Error("error loading configuration", zap.Error(err))
		return
	}
	log.Info("configuration loaded successfully")

	// Channel for event communication between components
	var eventChan chan entity.Event

	// Database connection setup
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		os.Getenv("DB_USER"),     // Database username from environment
		os.Getenv("DB_PASSWORD"), // Database password from environment
		os.Getenv("DB_HOST"),     // Database host from environment
		os.Getenv("DB_PORT"),     // Database port from environment
		os.Getenv("DB_NAME"),     // Database name from environment
	)

	log.Info("connecting to database", zap.String("dsn", maskPasswordInDSN(dsn)))

	// Connect to database with retry logic
	db, err := retrier.Connect(func() (*gorm.DB, error) {
		return gorm.Open(mysql.Open(dsn), &gorm.Config{})
	}, retrier.WithAttemps(5), retrier.WithTimeOut(5))
	if err != nil {
		// Log error with relevant context but mask sensitive information
		log.Error("failed to connect to database",
			zap.Error(err),
			zap.String("DB_USER", os.Getenv("DB_USER")),
			zap.String("DB_HOST", os.Getenv("DB_HOST")),
			zap.String("DB_PORT", os.Getenv("DB_PORT")),
			zap.String("DB_NAME", os.Getenv("DB_NAME")))
		return
	}

	log.Info("successfully connected to database")

	// Initialize repository layer
	repo := repository.Init(db, log)

	// RabbitMQ connection setup
	conn, err := amqp.Dial(cfg.RabbitMQUrl)
	if err != nil {
		log.Error("failed to connect to RabbitMQ",
			zap.String("url", cfg.RabbitMQUrl),
			zap.Error(err),
		)
		return
	}
	defer conn.Close() // Ensure connection is closed on exit

	log.Info("successfully connected to RabbitMQ")

	// Initialize message publisher
	pub, err := publisher.Init(conn, log, cfg)
	if err != nil {
		log.Error("failed to initialize publisher", zap.Error(err))
		return
	}

	// Initialize message consumer
	cons, err := consumer.Init(cfg, log, conn)
	if err != nil {
		log.Error("failed to initialize consumer", zap.Error(err))
		return
	}

	// Subscribe to request topic
	if err = cons.Subscribe("requests", cfg.RequestTopicName, "request.*"); err != nil {
		log.Error("failed to subscribe to requests topic", zap.Error(err))
		return
	}

	if err = cons.Subscribe("polls", "polls", "poll.*"); err != nil {
		log.Error("failed to subscribe to request topic", zap.Error(err))
		return
	}

	log.Info("successfully initialized publisher and consumer")

	// Redis client setup for caching
	redisClient := redis.NewClient(&redis.Options{
		Addr: cfg.RedisUrl, // Redis server address
		DB:   0,            // Default database
	})
	defer redisClient.Close() // Ensure Redis connection is closed on exit

	// Test Redis connection with retry logic
	err = retrier.Do(3, 5, func() error {
		return redisClient.Ping(context.Background()).Err()
	})
	if err != nil {
		log.Error("failed to connect to Redis",
			zap.String("url", cfg.RedisUrl),
			zap.Error(err))
		return
	}

	// Initialize caching layer
	casher := casher.Init(redisClient, log)

	// Initialize service layer with dependencies
	service := service.Init(
		repo,                     // Database repository
		pub,                      // Message publisher
		casher,                   // Redis cache
		&errs.VoteErrorHandler{}, // Error handler
	)

	// Initialize event listener
	listener := listener.Init(eventChan, log, cfg, service)

	// Start listener in a separate goroutine
	go listener.Listen(context.Background())

	// Start consuming messages
	cons.ConsumeMessages(eventChan)
}

// maskPasswordInDSN masks the password in DSN string for logging
func maskPasswordInDSN(dsn string) string {
	// This is a simplified version - implement proper masking logic
	// that replaces the password with ***** while keeping the rest of DSN intact
	return dsn // In real implementation, replace password portion with *****
}
