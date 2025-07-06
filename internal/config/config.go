package config

import (
	"html/template"
	"log"

	"github.com/alexedwards/scs/v2"
)

// AppConfig holds the application config
type AppConfig struct {
	UseCache      bool
	TemplateCache map[string]*template.Template
	InfoLog       *log.Logger
	InProduction  bool
	Session       *scs.SessionManager
	BookingRules  BookingConfig
}

type BookingConfig struct {
	OpenHour        int
	CloseHour       int
	MaxBookingHours int
}

func NewDefaultBookingConfig() BookingConfig {
	return BookingConfig{
		OpenHour:        8,
		CloseHour:       23,
		MaxBookingHours: 4,
	}
}
