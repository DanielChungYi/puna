package handlers

import (
	"encoding/json"
	"fmt"
	"slices"
	"strconv"
	"time"

	"log"
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

	render.RenderTemplate(w, r, "home.page.tmpl", &models.TemplateData{
		UserName:        m.App.Session.GetString(r.Context(), "user_name"),
		IsAuthenticated: m.App.Session.GetBool(r.Context(), "IsAuthenticated"),
	})
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
	minHour := 8
	maxHour := 22

	var hours []models.HourOption
	for h := minHour; h <= maxHour; h++ {
		hours = append(hours, models.HourOption{
			Value: h,
			Label: fmt.Sprintf("%02d:00", h),
		})
	}
	td := models.TemplateData{
		Hours: hours,
	}
	render.RenderTemplate(w, r, "search-availability.page.tmpl", &td)
}
func (m *Repository) CheckAvailability(w http.ResponseWriter, r *http.Request) {
	// 取得查詢參數
	selectedDate := r.URL.Query().Get("date")
	startHourStr := r.URL.Query().Get("start_hour")

	fmt.Println("📥 [CheckAvailability - GET] Query Params:")
	fmt.Println("  📅 date =", selectedDate)
	fmt.Println("  ⏰ start_hour =", startHourStr)

	if selectedDate == "" || startHourStr == "" {
		fmt.Println("❌ Missing required fields")
		http.Error(w, "date and start_hour are required", http.StatusBadRequest)
		return
	}

	// Parse date and hour
	bookingDate, err := time.Parse("2006-01-02", selectedDate)
	if err != nil {
		fmt.Println("❌ Invalid date format:", err)
		http.Error(w, "Invalid date format", http.StatusBadRequest)
		return
	}
	startHour, err := strconv.Atoi(startHourStr)
	if err != nil {
		fmt.Println("❌ Invalid start_hour:", err)
		http.Error(w, "Invalid start_hour", http.StatusBadRequest)
		return
	}

	// 計算 endHour，不能超過關門時間
	rules := m.App.BookingRules
	endHour := startHour + rules.MaxBookingHours
	if endHour > rules.CloseHour {
		endHour = rules.CloseHour
	}

	fmt.Println("📆 Parsed bookingDate:", bookingDate)
	fmt.Printf("⏰ Querying from %02d:00 to %02d:00\n", startHour, endHour)

	// 查詢資料庫
	availability, err := m.DB.GetCourtAvailabilityMapByTime(bookingDate, startHour, endHour)
	if err != nil {
		fmt.Println("❌ Failed to fetch availability:", err)
		http.Error(w, "Failed to fetch availability", http.StatusInternalServerError)
		return
	}

	fmt.Println("📊 Court availability map:", availability)

	type HourData struct {
		Value int    `json:"value"`
		Label string `json:"label"`
	}

	var result []HourData
	for h := startHour; h < endHour; h++ {
		label := fmt.Sprintf("%02d:00（可用 %d 面）", h, availability[h])
		result = append(result, HourData{Value: h, Label: label})
	}

	fmt.Println("✅ Final availability response:", result)

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(result)
}

// // Check Availablilty
// func (m *Repository) CheckAvailability(w http.ResponseWriter, r *http.Request) {
// 	// Parse form
// 	err := r.ParseForm()
// 	if err != nil {
// 		http.Error(w, "Failed to parse form", http.StatusBadRequest)
// 		return
// 	}

// 	selectedDate := r.PostForm.Get("date")
// 	startHourStr := r.PostForm.Get("start_hour")

// 	if selectedDate == "" || startHourStr == "" {
// 		http.Error(w, "date and start_hour are required", http.StatusBadRequest)
// 		return
// 	}

// 	// Parse date and hour
// 	bookingDate, err := time.Parse("2006-01-02", selectedDate)
// 	if err != nil {
// 		http.Error(w, "Invalid date format", http.StatusBadRequest)
// 		return
// 	}
// 	startHour, err := strconv.Atoi(startHourStr)
// 	if err != nil {
// 		http.Error(w, "Invalid start_hour", http.StatusBadRequest)
// 		return
// 	}

// 	// 計算 endHour，不能超過關門時間
// 	rules := m.App.BookingRules
// 	endHour := startHour + rules.MaxBookingHours
// 	if endHour > rules.CloseHour {
// 		endHour = rules.CloseHour
// 	}

// 	// 查詢資料庫
// 	availability, err := m.DB.GetCourtAvailabilityMapByTime(bookingDate, startHour, endHour)
// 	if err != nil {
// 		http.Error(w, "Failed to fetch availability", http.StatusInternalServerError)
// 		return
// 	}

// 	type HourData struct {
// 		Value int    `json:"value"`
// 		Label string `json:"label"`
// 	}

// 	var result []HourData
// 	for h := startHour; h < endHour; h++ {
// 		label := fmt.Sprintf("%02d:00（可用 %d 面）", h, availability[h])
// 		result = append(result, HourData{Value: h, Label: label})
// 	}

// 	w.Header().Set("Content-Type", "application/json")
// 	_ = json.NewEncoder(w).Encode(result)
// }

// Make Reservation
func (m *Repository) MakeReservation(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Failed to parse form", http.StatusBadRequest)
		return
	}

	selectedDate := r.PostForm.Get("date")
	startHourStr := r.PostForm.Get("start_hour")
	endHourStr := r.PostForm.Get("end_hour")

	fmt.Println("📥 Form Inputs:")
	fmt.Println("  📅 date =", selectedDate)
	fmt.Println("  ⏰ start_hour =", startHourStr)
	fmt.Println("  ⏰ end_hour =", endHourStr)

	bookingDate, err := time.Parse("2006-01-02", selectedDate)
	if err != nil {
		http.Error(w, "Invalid date format", http.StatusBadRequest)
		return
	}

	startHour, err := strconv.Atoi(startHourStr)
	endHour, err2 := strconv.Atoi(endHourStr)
	if err != nil || err2 != nil || endHour <= startHour {
		http.Error(w, "Invalid time range", http.StatusBadRequest)
		return
	}

	userID := m.App.Session.GetInt(r.Context(), "user_id")
	if userID == 0 {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// 查詢已有預約的場地 ID
	reservedCourtIDs, err := m.DB.GetReservedCourtIDs(bookingDate, startHour, endHour)
	if err != nil {
		http.Error(w, "Failed to query reservations", http.StatusInternalServerError)
		return
	}

	// 找第一個還沒被預約的場地
	var assignedCourtID uint
	for i := 1; i <= 8; i++ {
		if !slices.Contains(reservedCourtIDs, uint(i)) {
			assignedCourtID = uint(i)
			break
		}
	}

	if assignedCourtID == 0 {
		http.Error(w, "No available court", http.StatusConflict)
		return
	}

	res := models.Reservation{
		UserID:      uint(userID),
		CourtID:     assignedCourtID,
		BookingDate: bookingDate,
		StartHour:   startHour,
		EndHour:     endHour,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := m.DB.InsertReservation(res); err != nil {
		http.Error(w, "DB insert error", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"ok":         true,
		"date":       selectedDate,
		"start_hour": startHour,
		"end_hour":   endHour,
		"court_id":   assignedCourtID,
	})
}

// Login renders the login page
func (m *Repository) Login(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, r, "login.page.tmpl", &models.TemplateData{})
}

// Login authentication
func (m *Repository) PostLogin(w http.ResponseWriter, r *http.Request) {
	// Parse form
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Unable to parse form", http.StatusBadRequest)
		return
	}

	email := r.Form.Get("email")
	password := r.Form.Get("password")
	// Use the DB method to authenticate
	id, userEmail, Name, err := m.DB.Authenticate(email, password)
	if err != nil {
		log.Println("❌ Authentication failed:", err)
		http.Error(w, "Invalid login credentials", http.StatusUnauthorized)
		return
	}

	log.Printf("🔐 Logging in with Id: %d, Email: %s, Password HEX: %x", id, email, []byte(password))

	m.App.Session.Put(r.Context(), "user_id", id)
	m.App.Session.Put(r.Context(), "user_email", userEmail)
	m.App.Session.Put(r.Context(), "user_name", Name)
	m.App.Session.Put(r.Context(), "IsAuthenticated", true)

	log.Printf("✅ Login successful for user ID %d, mail:%s, name:%s", id, userEmail, Name)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// Logout renders the login page
func (m *Repository) Logout(w http.ResponseWriter, r *http.Request) {
	m.App.Session.Destroy(r.Context())
	m.App.Session.Put(r.Context(), "flash", "您已成功登出")
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// ShowRegister renders the register page
func (m *Repository) Register(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, r, "register.page.tmpl", &models.TemplateData{})
}

// PostShowRegister handles the registration form submission
func (m *Repository) PostRegister(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Unable to parse form", http.StatusBadRequest)
		return
	}

	Name := r.Form.Get("name")
	email := r.Form.Get("email")
	password := r.Form.Get("password")
	confirmPassword := r.Form.Get("confirm_password") // from HTML if you want double check on server too

	// 🪵 Debug logging of user input
	log.Printf("Register New account: First Name=%s, Email=%s, Password Len=%d, Confirm Password Len=%d", Name, email, len(password), len(confirmPassword))

	// Check if passwords match (optional, client already validated)
	if password != confirmPassword {
		log.Println("❌ Passwords do not match")
		http.Error(w, "Passwords do not match", http.StatusBadRequest)
		return
	}

	// Attempt to create the user
	userID, err := m.DB.CreateAccount(Name, email, password)
	if err != nil {
		log.Println("❌ Failed to create user:", err)
		http.Error(w, "Account creation failed: "+err.Error(), http.StatusInternalServerError)
		return
	}

	log.Printf("Register: User created successfully with ID: %d", userID)

	// Redirect to home or dashboard
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

type jsonResponse struct {
	Date string `json:"date"`
	Time string `json:"time"`
}

// Contact renders the contact page
func (m *Repository) Contact(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, r, "contact.page.tmpl", &models.TemplateData{})
}

func (m *Repository) AdminDashboard(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, r, "admin-dashboard.page.tmpl", &models.TemplateData{
		UserName:        m.App.Session.GetString(r.Context(), "user_name"),
		IsAuthenticated: m.App.Session.GetBool(r.Context(), "IsAuthenticated"),
	})
}
