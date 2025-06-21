package models

import (
	"time"
)

type User struct {
	ID          uint   `gorm:"primaryKey"`
	Name        string `gorm:"size:15;not null"`
	Email       string `gorm:"uniqueIndex;not null"`
	Password    string `gorm:"not null"`
	AccessLevel int    `gorm:"default:1"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type Court struct {
	ID        uint   `gorm:"primaryKey"`
	CourtName string `gorm:"size:100;not null;unique"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Restriction struct {
	ID              uint   `gorm:"primaryKey"`
	RestrictionName string `gorm:"size:100;not null;unique"`
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

type Reservation struct {
	ID           uint      `gorm:"primaryKey"`
	ResStartTime time.Time `gorm:"not null"`
	ResEndTime   time.Time `gorm:"not null"`
	CourtID      uint      `gorm:"not null"` // foreign key
	Court        Court     `gorm:"foreignKey:CourtID"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
