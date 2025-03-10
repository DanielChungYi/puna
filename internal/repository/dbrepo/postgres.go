package dbrepo

import (
	"log"

	"github.com/DanielChungYi/puna/internal/models"
)

func (m *postgresDBRepo) RunMigrate(res models.Reservation) error {
	err := m.DB.AutoMigrate(res)
	if err != nil {
		log.Println("❌ Fail to migrate:", err)
		return err
	}

	return nil
}

func (m *postgresDBRepo) InsertReservation(res models.Reservation) error {

	result := m.DB.Create(&res)
	if result.Error != nil {
		log.Println("❌ Fail to insert:", result.Error)
		return result.Error
	}

	return nil
}
