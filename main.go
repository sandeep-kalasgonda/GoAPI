package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	_ "goassignment/docs" // Path to generated docs

	"github.com/rs/cors" // CORS package
	// Swagger UI
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// Appointment represents the model for an appointment
// @Description This model is used to store appointment information
// @Param name query string true "Name of the person making the appointment"
// @Param email query string true "Email address of the person"
// @Param phone query string true "Phone number of the person"
// @Param doctor query string true "Name of the doctor"
// @Param date_time query string true "Date and time of the appointment"
type Appointment struct {
	ID       uint   `json:"id" gorm:"primaryKey"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	Doctor   string `json:"doctor"`
	DateTime string `json:"date_time"`
}

var db *gorm.DB

// initDatabase initializes the database connection and migrates the schema
func initDatabase() {
	var err error
	db, err = gorm.Open(sqlite.Open("appointments.db"), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect database:", err)
	}

	err = db.AutoMigrate(&Appointment{})
	if err != nil {
		log.Fatal("Failed to migrate database:", err)
	}
	fmt.Println("Database connection established and schema migrated.")
}

// @Summary Create a new appointment
// @Description Create a new appointment with customer details
// @Accept json
// @Produce json
// @Param appointment body Appointment true "Appointment Details"
// @Success 201 {object} Appointment
// @Failure 400 {string} string "Invalid JSON payload"
// @Failure 500 {string} string "Failed to create appointment"
// @Router /appointments [post]
func createAppointment(w http.ResponseWriter, r *http.Request) {
	var appointment Appointment
	if err := json.NewDecoder(r.Body).Decode(&appointment); err != nil {
		http.Error(w, "Invalid JSON payload", http.StatusBadRequest)
		return
	}
	result := db.Create(&appointment)
	if result.Error != nil {
		http.Error(w, "Failed to create appointment", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(appointment)
}

// @Summary Get all appointments
// @Description Get a list of all appointments
// @Produce json
// @Success 200 {array} Appointment
// @Failure 500 {string} string "Failed to retrieve appointments"
// @Router /appointments/all [get]
func getAllAppointments(w http.ResponseWriter, r *http.Request) {
	var appointments []Appointment
	result := db.Find(&appointments)
	if result.Error != nil {
		http.Error(w, "Failed to retrieve appointments", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(appointments)
}

// @Summary Get an appointment by ID
// @Description Get an appointment by its ID
// @Produce json
// @Param id query int true "Appointment ID"
// @Success 200 {object} Appointment
// @Failure 400 {string} string "Invalid appointment ID"
// @Failure 404 {string} string "Appointment not found"
// @Router /appointments/get [get]
func getAppointmentByID(w http.ResponseWriter, r *http.Request) {
	idParam := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idParam)
	if err != nil || id < 1 {
		http.Error(w, "Invalid appointment ID", http.StatusBadRequest)
		return
	}
	var appointment Appointment
	result := db.First(&appointment, id)
	if result.Error != nil {
		http.Error(w, "Failed to retrieve appointment", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(appointment)
}

// @Summary Update an appointment
// @Description Update an existing appointment
// @Accept json
// @Produce json
// @Param appointment body Appointment true "Appointment Details"
// @Success 200 {object} Appointment
// @Failure 400 {string} string "Invalid JSON payload"
// @Failure 404 {string} string "Appointment not found"
// @Failure 500 {string} string "Failed to update appointment"
// @Router /appointments/update [put]
func updateAppointment(w http.ResponseWriter, r *http.Request) {
	var updateAppointment Appointment
	if err := json.NewDecoder(r.Body).Decode(&updateAppointment); err != nil {
		http.Error(w, "Invalid JSON payload", http.StatusBadRequest)
		return
	}
	var existingAppointment Appointment
	result := db.First(&existingAppointment, updateAppointment.ID)
	if result.Error != nil {
		http.Error(w, "Failed to retrieve appointment", http.StatusNotFound)
		return
	}

	// Update fields
	existingAppointment.Name = updateAppointment.Name
	existingAppointment.Email = updateAppointment.Email
	existingAppointment.Phone = updateAppointment.Phone
	existingAppointment.DateTime = updateAppointment.DateTime
	existingAppointment.Doctor = updateAppointment.Doctor

	db.Save(&existingAppointment)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(existingAppointment)
}

func main() {
	initDatabase()

	// Enable CORS for Swagger UI and API routes
	c := cors.New(cors.Options{
		AllowedOrigins: []string{"*"}, // Allow all origins
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
	})

	// Serve Swagger UI
	http.Handle("/swagger/", http.StripPrefix("/swagger", http.FileServer(http.Dir("./docs"))))

	// API routes
	http.HandleFunc("/appointments", createAppointment)
	http.HandleFunc("/appointments/all", getAllAppointments)
	http.HandleFunc("/appointments/get", getAppointmentByID)
	http.HandleFunc("/appointments/update", updateAppointment)

	// Start the server with CORS middleware applied
	fmt.Println("Server is running on http://localhost:8080")
	err := http.ListenAndServe(":8080", c.Handler(http.DefaultServeMux))
	if err != nil {
		log.Fatal("Error starting server: ", err)
	}
}
