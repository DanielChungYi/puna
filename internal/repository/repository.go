package repository

import "github.com/DanielChungYi/puna/internal/models"

type DatabaseRepo interface {
	RunMigrate() error

	InsertReservation(res models.Reservation) error
	CreateAccount(Name, email, plainPassword string) (int, error)
	Authenticate(email, testPassword string) (int, string, error)
}
