package dbrepo

import (
	"errors"
	"log"

	"github.com/DanielChungYi/puna/internal/models"
	"golang.org/x/crypto/bcrypt"
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

func (m *postgresDBRepo) Authenticate(email, testPassword string) (int, string, error) {
	var user models.User

	result := m.DB.Where("email = ?", email).First(&user)
	if result.Error != nil {
		return 0, "", result.Error
	}

	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(testPassword))
	if err == bcrypt.ErrMismatchedHashAndPassword {
		return 0, "", errors.New("incorrect password")
	} else if err != nil {
		return 0, "", err
	}

	return int(user.ID), user.Email, nil
}
