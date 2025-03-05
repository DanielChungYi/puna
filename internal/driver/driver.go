package driver

import (
	"database/sql"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type DB struct {
	SQL *sql.DB
}

var dbConn = &DB{}

const maxOpenDBConn = 10
const maxIdleDBconn = 5
const maxDbLifetime = 5 * time.Minute

// Connect to SQL
func ConnectSQL(dsn string) (*DB, error) {
	db, err := NewDataBase(dsn)
	if err != nil {
		panic(err)
	}

	// Get underlying *sql.DB
	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	// Set connection pool settings
	sqlDB.SetMaxOpenConns(20)                  // Max open connections
	sqlDB.SetMaxIdleConns(10)                  // Max idle connections
	sqlDB.SetConnMaxLifetime(30 * time.Minute) // Max lifetime of a connection
	sqlDB.SetConnMaxIdleTime(10 * time.Minute) // Max idle time before closing connection
}

// NewDataBase creates a new database for the applications
func NewDataBase(dsn string) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	return db, nil
}

// testDB tries to ping the database
func testDB(db *sql.DB) error {
	err := db.Ping()
	if err != nil {
		return err
	}

	return nil
}
