package main

import (
	"encoding/json"
	"fmt"
	_ "goassignment/docs" // Path to generated docs
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	"github.com/rs/cors" // CORS package
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// Appointment represents the model for an appointment
// @Description This model is used to store appointment information
// @Param name query string true "Name of the person making the appointment"
// @Param email query string true "Email address of the person"
// @Param phone query string true "Phone number of the person"
// @Param doctor query string true "Name of the doctor"
// @Param date_time query string true "Date and time of the appointment"

// Appointment represents a scheduled meeting with a doctor.
// It contains details such as the appointment ID, patient's name, email, phone number,
// the doctor's name, and the date and time of the appointment.
type Appointment struct {
	ID       uint   `json:"id" gorm:"primaryKey"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	Doctor   string `json:"doctor"`
	DateTime string `json:"date_time"`
}

// db is a global variable that holds the database connection instance
// using the GORM library. It is used to interact with the database
// throughout the application.
var db *gorm.DB

// initDatabase initializes the database connection and migrates the schema
func initDatabase() {
	var err error
	err = godotenv.Load() // Load environment variables from .env file
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASS")
	dbName := os.Getenv("DB_NAME")
	dbHost := os.Getenv("DB_HOST")

	// Print the values for debugging
	fmt.Println("DB User:", dbUser)
	fmt.Println("DB Password:", dbPassword)
	fmt.Println("DB Host:", dbHost)
	fmt.Println("DB Name:", dbName)

	// Create the Data Source Name (DSN) for MySQL
	dns := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		dbUser, dbPassword, dbHost, dbName)

	// Open a connection to the database
	db, err = gorm.Open(mysql.Open(dns), &gorm.Config{})
	if err != nil {
		log.Fatal("failed to connect to database:", err)
	}
	fmt.Println("Database connection established.")

	// Migrate the schema
	err = db.AutoMigrate(&Appointment{})
	if err != nil {
		log.Fatal("failed to migrate database:", err)
	}
	fmt.Println("Database migration completed.")
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

// createAppointment handles the creation of a new appointment.
// It decodes the JSON payload from the request body into an Appointment struct,
// validates the payload, and inserts the new appointment into the database.
// If successful, it returns the created appointment with a 201 status code.
// If the JSON payload is invalid, it returns a 400 status code.
// If there is an error creating the appointment in the database, it returns a 500 status code.
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

// getAllAppointments handles HTTP requests to retrieve all appointments from the database.
// It responds with a JSON array of appointments or an error message if the retrieval fails.
//
// Parameters:
//   - w: http.ResponseWriter to write the HTTP response.
//   - r: *http.Request representing the HTTP request.
//
// Response:
//   - On success: HTTP status 200 with a JSON array of appointments.
//   - On failure: HTTP status 500 with an error message.
func getAllAppointments(w http.ResponseWriter, r *http.Request) {
	var appointments []Appointment
	if err := db.Find(&appointments).Error; err != nil {
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

// getAppointmentByID handles HTTP requests to retrieve an appointment by its ID.
// It expects an "id" parameter in the query string, which should be a positive integer.
// If the "id" parameter is missing or invalid, it responds with a "400 Bad Request" status.
// If the appointment with the given ID is not found, it responds with a "404 Not Found" status.
// On success, it responds with the appointment details in JSON format and a "200 OK" status.
//
// Parameters:
// - w: http.ResponseWriter to write the HTTP response.
// - r: *http.Request containing the HTTP request details.
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

// updateAppointment handles HTTP requests to update an existing appointment.
// It decodes the JSON payload from the request body into an Appointment struct,
// validates the payload, and updates the corresponding appointment in the database.
// If successful, it returns the updated appointment with a 200 status code.
// If the JSON payload is invalid, it returns a 400 status code.
// If the appointment with the given ID is not found, it returns a 404 status code.
// If there is an error updating the appointment in the database, it returns a 500 status code.
//
// Parameters:
// - w: http.ResponseWriter to write the HTTP response.
// - r: *http.Request containing the HTTP request details.
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

// deleteAppointment handles the deletion of an appointment based on the provided ID.
// It expects an "id" parameter in the URL query string, which should be a valid integer.
// If the ID is invalid or less than 1, it responds with a "400 Bad Request" status.
// If the appointment with the given ID is not found, it responds with a "404 Not Found" status.
// If the appointment is successfully deleted, it responds with a "200 OK" status.
//
// Parameters:
// - w: http.ResponseWriter to write the HTTP response.
// - r: *http.Request containing the HTTP request details.
func deleteAppointment(w http.ResponseWriter, r *http.Request) {
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
	db.Delete(&appointment)
	w.WriteHeader(http.StatusOK)
}

func main() {
	initDatabase()

	// Enable CORS for Swagger UI and API routes
	// c is an instance of the cors.Cors struct, which is configured with specific options.
	// The AllowedOrigins option is set to allow all origins by using the wildcard "*",
	// meaning that any domain can access the resources.
	// The AllowedMethods option specifies the HTTP methods that are permitted,
	// which include GET, POST, PUT, DELETE, and OPTIONS.
	c := cors.New(cors.Options{
		AllowedOrigins: []string{"*"}, // Allow all origins
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
	})

	// Serve Swagger UI
	http.Handle("/swagger/", http.StripPrefix("/swagger", http.FileServer(http.Dir("./docs"))))

	// API routes
	// Define the HTTP endpoints and their corresponding handler functions.
	// These routes handle various CRUD operations for appointments.
	http.HandleFunc("/appointments", createAppointment)        // Create a new appointment
	http.HandleFunc("/appointments/all", getAllAppointments)   // Get all appointments
	http.HandleFunc("/appointments/get", getAppointmentByID)   // Get an appointment by ID
	http.HandleFunc("/appointments/update", updateAppointment) // Update an existing appointment
	http.HandleFunc("/appointments/delete", deleteAppointment) // Delete an appointment by ID

	// Start the server with CORS middleware applied
	// The server listens on port 8080 and uses the CORS handler to manage cross-origin requests.
	fmt.Println("Server is running on http://localhost:8080")
	err := http.ListenAndServe(":8080", c.Handler(http.DefaultServeMux))
	if err != nil {
		log.Fatal("Error starting server: ", err)
	}
}
