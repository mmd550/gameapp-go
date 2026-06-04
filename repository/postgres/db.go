package postgres

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Config struct {
	Host     string
	Port     int
	Username string
	Password string
	DBName   string
	// SSLMOde controls the SSL connection behavior
	// Values: "disable", "allow", "prefer", "require", "verify-ca", "verify-full"
	// Normally is set to "disable" for local dev and "require" | "verify-full" for production
	SSLMode string
}

type connectionsConfig struct {
	maxOpenConnections  int
	maxIdleConnections  int
	connectionMaxLifetime int
}

type PostgresDB struct {
	db *gorm.DB
}

func (m *PostgresDB) Conn() *gorm.DB {
	return m.db
}

func New(cfg Config, gormConfig *gorm.Config) *PostgresDB {
	dsn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.Username, cfg.Password, cfg.DBName, cfg.SSLMode,
	)

	db, err := gorm.Open(postgres.Open(dsn), gormConfig)
	if err != nil {
		panic("failed to connect database")
	}

	sqlDB, err := db.DB()
	if err != nil {
		panic("failed to get underlying sql.DB")
	}

	connConfig := getEnvValues()

	sqlDB.SetMaxOpenConns(connConfig.maxOpenConnections)
	sqlDB.SetMaxIdleConns(connConfig.maxIdleConnections)
	sqlDB.SetConnMaxLifetime(time.Duration(connConfig.connectionMaxLifetime))


	return &PostgresDB{db}
}

func (p *PostgresDB) Migrate() error {
	return p.db.AutoMigrate(
		&User{},
	)
}

func getEnvValues() connectionsConfig {
	maxOpen, err := strconv.Atoi(os.Getenv("DB_MAX_OPEN_CONNS"))
	if err != nil {
		maxOpen = 25
	}

	maxIdle, err := strconv.Atoi(os.Getenv("DB_MAX_IDLE_CONNS"))
	if err != nil {
		maxIdle = 10
	}

	connLifetime, err := time.ParseDuration(os.Getenv("DB_CONN_MAX_LIFETIME"))
	if err != nil {
		connLifetime = 30 * time.Minute
	}

	return connectionsConfig{
		maxOpenConnections:  maxOpen,
		maxIdleConnections:  maxIdle,
		connectionMaxLifetime: int(connLifetime),
	}
}
