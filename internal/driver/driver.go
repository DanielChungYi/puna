package driver

import (
	"database/sql"
	"time"

	"github.com/DanielChungYi/puna/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type DB struct {
	GORM *gorm.DB
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
	gormDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	// Test connection
	testDB(gormDB)

	// Set connection pool settings
	gormDB.SetMaxOpenConns(20)                  // Max open connections
	gormDB.SetMaxIdleConns(10)                  // Max idle connections
	gormDB.SetConnMaxLifetime(30 * time.Minute) // Max lifetime of a connection
	gormDB.SetConnMaxIdleTime(10 * time.Minute) // Max idle time before closing connection
	mydb := DB{GORM: db}
	return &mydb, nil
}

// NewDataBase creates a new database for the applications
func NewDataBase(dsn string) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	return db, nil
}

func (db *DB) RunMigrations() error {
	return db.GORM.AutoMigrate(
		&models.User{},
		&models.Court{},
		&models.Reservation{},
		&models.Restriction{},
	)
}

// testDB tries to ping the database
func testDB(db *sql.DB) error {
	err := db.Ping()
	if err != nil {
		return err
	}

	return nil
}
