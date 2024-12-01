package mysql

import (
	"database/sql"
	"fmt"
	"log"
	"net/url"
	"time"
)

const (
	// DefaultMaxOpen is default value for max open connection
	DefaultMaxOpen = 10
	// DefaultMaxIdle is default value for max idle connection
	DefaultMaxIdle = 10
	// DefaultMaxLifetime is default value for max connection lifetime
	DefaultMaxLifetime = 3 * time.Minute
	// DriverMySQL is driver name for mysql
	DriverMySQL = "mysql"
)

type Config struct {
	Host        string
	Port        string
	User        string
	Password    string
	Name        string
	maxOpen     int
	maxIdle     int
	maxLifetime time.Duration
	maxIdleTime time.Duration
	dsn         string
	driverName  string
	err         error
}

// DB return new sql db
func DB(config *Config, options ...Option) (*sql.DB, error) {
	defaults(config)

	for _, option := range options {
		option(config)
	}

	if config.err != nil {
		return nil, config.err
	}

	var db *sql.DB
	var err error

	if db, err = sql.Open(config.driverName, config.dsn); err != nil {
		log.Printf("error while opening mysql DB: %v", err)
		return nil, err
	}

	db.SetMaxOpenConns(config.maxOpen)
	db.SetMaxIdleConns(config.maxIdle)
	db.SetConnMaxLifetime(config.maxLifetime)
	// db.SetConnMaxIdleTime(config.maxIdleTime)

	return db, nil
}

func mysqlDSN(config *Config, parseTime bool, location string) string {
	connection := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", config.User, config.Password, config.Host, config.Port, config.Name)
	val := url.Values{}

	if parseTime {
		val.Add("parseTime", "1")
	}
	if len(location) > 0 {
		val.Add("loc", location)
	}

	if len(val) == 0 {
		return connection
	}
	return fmt.Sprintf("%s?%s", connection, val.Encode())
}

type Option func(*Config)

func defaults(config *Config) {
	config.maxOpen = DefaultMaxOpen
	config.maxIdle = DefaultMaxIdle
	config.maxLifetime = DefaultMaxLifetime
	config.driverName = DriverMySQL
	config.dsn = mysqlDSN(config, true, "Asia/Jakarta")
}

func WithConnection(maxOpen, maxIdle int, maxLifetime, maxIdleTime time.Duration) Option {
	return func(config *Config) {
		if maxOpen > 0 {
			config.maxOpen = maxOpen
		}
		if maxIdle > 0 {
			config.maxIdle = maxIdle
		}
		if maxLifetime > 0 {
			config.maxLifetime = maxLifetime
		}
		if maxIdleTime > 0 {
			config.maxIdleTime = maxIdleTime
		}
	}
}

func WithMysql(serverName string,
	parseTime bool, location string) Option {

	return func(config *Config) {
		config.driverName = DriverMySQL
		config.dsn = mysqlDSN(config, parseTime, location)
	}
}
