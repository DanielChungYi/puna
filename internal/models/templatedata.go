package models

import (
	"time"

	forms "github.com/DanielChungYi/puna/internal/form"
)

// TemplateData holds data sent from handlers to templates
type TemplateData struct {
	StringMap map[string]string
	IntMap    map[string]int
	FloatMap  map[string]float32
	Data      map[string]interface{}
	CSRFToken string
	Flash     string
	Warning   string
	Error     string
	Form      *forms.Form
}

type Reservation struct {
	ID           int       `gorm:"primaryKey"`
	ResStartTime time.Time `gorm:"not null"`       // Start date and time of the reservation
	ResEndTime   time.Time `gorm:"not null"`       // End date and time of the reservation
	CreatedAt    time.Time `gorm:"autoCreateTime"` // Auto-generated creation timestamp
}
