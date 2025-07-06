package dbrepo

import (
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/DanielChungYi/puna/internal/models"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func (m *postgresDBRepo) RunMigrate() error {
	err := m.DB.AutoMigrate(
		&models.User{},
		&models.Court{},
		&models.Reservation{},
		&models.Restriction{},
	)
	if err != nil {
		log.Println("❌ Failed to migrate:", err)
		return err
	}
	return nil
}

func (m *postgresDBRepo) InsertReservation(res models.Reservation) error {
	result := m.DB.Create(&res)
	if result.Error != nil {
		log.Println("❌ Failed to insert reservation:", result.Error)
		return result.Error
	}
	return nil
}

func (m *postgresDBRepo) UpdateReservation(res models.Reservation) error {
	result := m.DB.Save(&res)
	if result.Error != nil {
		log.Println("❌ Failed to update reservation:", result.Error)
		return result.Error
	}
	return nil
}

func (m *postgresDBRepo) CreateAccount(Name, email, plainPassword string) (int, error) {
	// Check if user already exists
	var existing models.User
	if err := m.DB.Where("email = ?", email).First(&existing).Error; err == nil {
		return 0, errors.New("email already registered")
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return 0, err
	}

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(plainPassword), bcrypt.DefaultCost)
	if err != nil {
		return 0, err
	}

	// Create the user
	user := models.User{
		Name:        Name,
		Email:       email,
		Password:    string(hashedPassword),
		AccessLevel: 1,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := m.DB.Create(&user).Error; err != nil {
		return 0, err
	}

	return int(user.ID), nil
}

func (m *postgresDBRepo) Authenticate(email, testPassword string) (int, string, string, error) {
	var user models.User

	result := m.DB.Where("email = ?", email).First(&user)
	if result.Error != nil {
		return 0, "", "", result.Error
	}

	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(testPassword))
	if err == bcrypt.ErrMismatchedHashAndPassword {
		return 0, "", "", errors.New("incorrect password")
	} else if err != nil {
		return 0, "", "", err
	}

	return int(user.ID), user.Email, user.Name, nil
}

func (m *postgresDBRepo) GetReservedCourtIDs(date time.Time, startHour, endHour int) ([]uint, error) {
	var reservedCourtIDs []uint
	err := m.DB.Model(&models.Reservation{}).
		Where("booking_date = ? AND start_hour < ? AND end_hour > ?", date, endHour, startHour).
		Pluck("court_id", &reservedCourtIDs).Error
	return reservedCourtIDs, err
}

// GetCourtAvailabilityMap returns a map of available courts per hour for a given date string.
func (m *postgresDBRepo) GetCourtAvailabilityMapByDate(dateStr string) (map[int]int, error) {
	fmt.Println("🔍 [GetCourtAvailabilityMapByDate] Input dateStr:", dateStr)

	// Parse date
	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		fmt.Println("❌ Invalid date format:", err)
		return nil, fmt.Errorf("invalid date format: %v", err)
	}
	fmt.Println("📅 Parsed date (time.Time):", date)

	type result struct {
		Hour            int
		AvailableCourts int
	}
	var availability []result

	query := `
		WITH hours AS (
			SELECT generate_series(8, 22) AS hour
		)
		SELECT 
			h.hour,
			COUNT(c.id) AS available_courts
		FROM 
			hours h
		CROSS JOIN 
			courts c
		WHERE NOT EXISTS (
			SELECT 1
			FROM reservations r
			WHERE r.court_id = c.id
			  AND r.booking_date = ?
			  AND r.start_hour <= h.hour
			  AND r.end_hour > h.hour
		)
		GROUP BY h.hour
		ORDER BY h.hour;
	`

	fmt.Println("🛠️ Executing SQL query...")

	err = m.DB.Raw(query, date).Scan(&availability).Error
	if err != nil {
		fmt.Println("❌ SQL query error:", err)
		return nil, err
	}

	fmt.Println("📊 Raw availability result:", availability)

	availabilityMap := make(map[int]int)
	for h := 8; h <= 22; h++ {
		availabilityMap[h] = 0
	}
	for _, a := range availability {
		availabilityMap[a.Hour] = a.AvailableCourts
	}

	fmt.Println("✅ Final availability map:", availabilityMap)

	return availabilityMap, nil
}

func (m *postgresDBRepo) GetCourtAvailabilityMapByTime(date time.Time, startHour, endHour int) (map[int]int, error) {
	// Truncate time component to ensure we only compare date
	dayOnly := date.Truncate(24 * time.Hour)

	type result struct {
		Hour            int
		AvailableCourts int
	}

	var availability []result

	query := `
		WITH hours AS (
			SELECT generate_series($1, $2 - 1) AS hour
		)
		SELECT 
			h.hour,
			COUNT(c.id) AS available_courts
		FROM 
			hours h
		CROSS JOIN 
			courts c
		WHERE NOT EXISTS (
			SELECT 1
			FROM reservations r
			WHERE r.court_id = c.id
			  AND r.booking_date = $3
			  AND r.start_hour < h.hour + 1
			  AND r.end_hour > h.hour
		)
		GROUP BY h.hour
		ORDER BY h.hour;
	`

	err := m.DB.Raw(query, startHour, endHour, dayOnly).Scan(&availability).Error
	if err != nil {
		return nil, err
	}

	// Convert result to map
	availabilityMap := make(map[int]int)
	for _, a := range availability {
		availabilityMap[a.Hour] = a.AvailableCourts
	}

	return availabilityMap, nil
}
