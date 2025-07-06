package driver

import (
	"database/sql"
	"fmt"
	"log"
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
	err := db.GORM.AutoMigrate(
		&models.User{},
		&models.Court{},
		&models.Reservation{},
		&models.Restriction{},
	)

	if err != nil {
		return err
	}

	// Add composite index manually
	createIndex := `
		CREATE INDEX IF NOT EXISTS idx_reservation_court_date_hours
		ON reservations (court_id, booking_date, start_hour, end_hour);
	`
	if err := db.GORM.Exec(createIndex).Error; err != nil {
		log.Println("⚠️ Failed to create composite index:", err)
		return err
	}

	// Seed 12 courts
	if err := SeedCourts(db.GORM); err != nil {
		log.Fatal("seeding courts failed:", err)
	}

	return nil
}

func SeedCourts(db *gorm.DB) error {
	// Check if courts already exist
	var count int64
	if err := db.Model(&models.Court{}).Count(&count).Error; err != nil {
		return err
	}

	if count > 0 {
		log.Println("✅ Courts already seeded")
		return nil
	}

	// Create 12 courts
	var courts []models.Court
	for i := 1; i <= 12; i++ {
		courts = append(courts, models.Court{
			CourtName: "Court " + itoa(i),
		})
	}

	if err := db.Create(&courts).Error; err != nil {
		return err
	}

	log.Println("✅ Successfully seeded 12 courts")
	return nil
}

// Helper function to convert int to string
func itoa(i int) string {
	return fmt.Sprintf("%d", i)
}

// testDB tries to ping the database
func testDB(db *sql.DB) error {
	err := db.Ping()
	if err != nil {
		return err
	}

	return nil
}
