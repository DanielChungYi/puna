package handlers

import (
	"encoding/json"
	"fmt"

	"log"
	"net/http"

	"github.com/DanielChungYi/puna/internal/config"
	forms "github.com/DanielChungYi/puna/internal/form"
	"github.com/DanielChungYi/puna/internal/models"
	"github.com/DanielChungYi/puna/internal/render"
)

// Repo the repository used by the handlers
var Repo *Repository

// Repository is the repository type
type Repository struct {
	App *config.AppConfig
}

// NewRepo creates a new repository
func NewRepo(a *config.AppConfig) *Repository {
	return &Repository{
		App: a,
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
	OK      bool   `json:"ok"`
	Message string `json:"message"`
	Date    string `json:"data"`
	Time    string `json:"time"`
}

// PostAvailability handles post
func (m *Repository) PostAvailability(w http.ResponseWriter, r *http.Request) {
	date := r.Form.Get("reservation-dates")
	time := r.Form.Get("selected-time")

	// Dump all the post data
	for key, values := range r.Form {
		for _, value := range values {
			fmt.Printf("%s = %s\n", key, value)
		}
	}

	// Place the DB lookup logic here and sned it back to client
	resp := jsonResponse{
		OK:      true,
		Message: "Available!",
		Date:    date,
		Time:    time,
	}

	out, err := json.MarshalIndent(resp, "", "     ")
	if err != nil {
		log.Println(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(out)

	fmt.Println("Respon JSON:", string(out))

	//w.Write([]byte(fmt.Sprintf("Date is %s and time is %s", data, time)))
}

// Contact renders the contact page
func (m *Repository) Contact(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, r, "contact.page.tmpl", &models.TemplateData{})
}
