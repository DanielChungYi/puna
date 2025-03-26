package handlers

import (
	"encoding/json"
	"fmt"
	"time"

	"log"
	"math/rand"
	"net/http"

	"github.com/DanielChungYi/puna/internal/config"
	"github.com/DanielChungYi/puna/internal/driver"
	forms "github.com/DanielChungYi/puna/internal/form"
	"github.com/DanielChungYi/puna/internal/models"
	"github.com/DanielChungYi/puna/internal/render"
	"github.com/DanielChungYi/puna/internal/repository"
	"github.com/DanielChungYi/puna/internal/repository/dbrepo"
)

// Repo the repository used by the handlers
var Repo *Repository

// Repository is the repository type
type Repository struct {
	App *config.AppConfig
	DB  repository.DatabaseRepo
}

// NewRepo creates a new repository
func NewRepo(a *config.AppConfig, db *driver.DB) *Repository {
	return &Repository{
		App: a,
		DB:  dbrepo.NewPostgresRepo(a, db.GORM),
	}
}

// NewHandlers sets the repository for the handlers
func NewHandlers(r *Repository) {
	Repo = r
}

// Home is the handler for the home page
func (m *Repository) Home(w http.ResponseWriter, r *http.Request) {
	remoteIP := r.RemoteAddr
	m.App.Session.Put(r.Context(), "remote_ip", remoteIP)

	render.RenderTemplate(w, r, "home.page.tmpl", &models.TemplateData{})
}

// Reservation renders the make a reservation page and displays form
func (m *Repository) Reservation(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, r, "make-reservation.page.tmpl", &models.TemplateData{
		Form: forms.New(nil),
	})
}

// Court information
func (m *Repository) CourtInfo(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, r, "court-info.page.tmpl", &models.TemplateData{})
}

// Availability renders the search availability page
func (m *Repository) Availability(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, r, "search-availability.page.tmpl", &models.TemplateData{})
}

type jsonResponse struct {
	Date string `json:"date"`
	Time string `json:"time"`
}

// PostAvailability handles post
func (m *Repository) PostAvailability(w http.ResponseWriter, r *http.Request) {
	selectedDate := r.PostForm.Get("date")
	selectedTime := r.PostForm.Get("timeslot")

	// Dump all the post data
	for key, values := range r.PostForm {
		for _, value := range values {
			fmt.Printf("%s = %s\n", key, value)
		}
	}

	// Place the DB lookup logic here and send it back to client
	// Convertsion
	rDate, err := time.Parse("2006-01-02", selectedDate)
	if err != nil {
		fmt.Println("Error parsing date:", err)
		w.WriteHeader(400)
		return
	}
	rTime, err := time.Parse("3:04 PM", selectedTime)
	if err != nil {
		fmt.Println("Error parsing time:", err)
		w.WriteHeader(400)
		w.Write([]byte("{ \"message\": \"test failure\" }"))
		return
	}
	// Combine the parsed date and time into a single time.Time
	sDate := time.Date(rDate.Year(), rDate.Month(), rDate.Day(),
		rTime.Hour(), rTime.Minute(), rTime.Second(), rTime.Nanosecond(),
		time.Local)
	eDate := sDate.Add(time.Hour)

	// Insert Reservation
	reservation := models.Reservation{
		ID:           rand.Intn(1000),
		ResStartTime: sDate,
		ResEndTime:   eDate,
		CreatedAt:    time.Now(),
	}
	m.DB.RunMigrate(reservation)
	m.DB.InsertReservation(reservation)

	// Response
	resp := jsonResponse{
		Date: selectedDate,
		Time: selectedTime,
	}

	out, err := json.MarshalIndent(resp, "", "     ")
	if err != nil {
		log.Println(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(201)
	w.Write(out)

	fmt.Println("Respon JSON:", string(out))

	//w.Write([]byte(fmt.Sprintf("Date is %s and time is %s", data, time)))
}

// Contact renders the contact page
func (m *Repository) Contact(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, r, "contact.page.tmpl", &models.TemplateData{})
}
