package repository

import (
	"time"

	"github.com/DanielChungYi/puna/internal/models"
)

type DatabaseRepo interface {
	RunMigrate() error

	InsertReservation(res models.Reservation) error
	CreateAccount(Name, email, plainPassword string) (int, error)
	Authenticate(email, testPassword string) (int, string, string, error)
	GetReservedCourtIDs(date time.Time, startHour, endHour int) ([]uint, error)
	GetCourtAvailabilityMapByDate(dateStr string) (map[int]int, error)
	GetCourtAvailabilityMapByTime(date time.Time, startHour, endHour int) (map[int]int, error)
}
