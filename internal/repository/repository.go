package repository

import "github.com/DanielChungYi/puna/internal/models"

type DatabaseRepo interface {
	RunMigrate(res models.Reservation) error

	InsertReservation(res models.Reservation) error
}
