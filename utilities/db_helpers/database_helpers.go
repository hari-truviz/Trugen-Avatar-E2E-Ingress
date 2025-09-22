package DBHelpers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
)

var DB *sql.DB

// HandleCreateDemoFeedback creates a new demo feedback entry
func HandleCreateDemoFeedback(w http.ResponseWriter, r *http.Request) {
	var feedback struct {
		ConversationID string `json:"conversation_id"`
		Email          string `json:"email"`
		Feedback       string `json:"feedback"`
		Rating         int    `json:"rating"`
	}

	// Decode request body
	if err := json.NewDecoder(r.Body).Decode(feedback); err != nil {
		log.Printf("Error decoding request body: %v", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validation
	if feedback.Email == "" {
		http.Error(w, "Email is required", http.StatusBadRequest)
		return
	}
	if feedback.Rating < 0 || feedback.Rating > 5 { // assuming 1–5 scale
		http.Error(w, "rating must be between 0 and 5", http.StatusBadRequest)
		return
	}

	query := `
        INSERT INTO demo_feedbacks (
            id, conversation_id, email, feedback, rating, created_at
        ) VALUES (
            gen_random_uuid(), $1, $2, $3, $4, CURRENT_TIMESTAMP
        ) RETURNING id`

	var id string
	err := DB.QueryRow(
		query,
		feedback.ConversationID,
		feedback.Email,
		feedback.Feedback,
		feedback.Rating,
	).Scan(&id)

	if err != nil {
		log.Printf("Error creating demo feedback: %v", err)
		http.Error(w, "Failed to create demo feedback", http.StatusInternalServerError)
		return
	}

	// Respond
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"id":      id,
		"message": "Demo Feedback submitted successfully",
	})
}

// HandleCreateWaitlist creates a new waitlist entry
func HandleCreateWaitlist(w http.ResponseWriter, r *http.Request) {
	var waitlist struct {
		Name  string `json:"name"`
		Email string `json:"email"`
	}

	// Decode request body
	if err := json.NewDecoder(r.Body).Decode(&waitlist); err != nil {
		log.Printf("Error decoding request body: %v", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validation
	if waitlist.Email == "" {
		http.Error(w, "Email is required", http.StatusBadRequest)
		return
	}
	// Validation
	if waitlist.Name == "" {
		http.Error(w, "Name is required", http.StatusBadRequest)
		return
	}

	query := `
        INSERT INTO waitlists (
            id, name, email, created_at
        ) VALUES (
            gen_random_uuid(), $1, $2, CURRENT_TIMESTAMP
        ) RETURNING id`

	var id string
	err := DB.QueryRow(
		query,
		waitlist.Name,
		waitlist.Email,
	).Scan(&id)

	if err != nil {
		log.Printf("Error creating waitlist: %v", err)
		http.Error(w, "Failed to create waitlist", http.StatusInternalServerError)
		return
	}

	// Respond
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"id":      id,
		"message": "Waitlist submitted successfully",
	})
}

// HandleCreateDemoUser creates a new demo user entry
func HandleCreateDemoUser(w http.ResponseWriter, r *http.Request) {
	var demouser struct {
		Name  string `json:"name"`
		Email string `json:"email"`
	}

	// Decode request body
	if err := json.NewDecoder(r.Body).Decode(&demouser); err != nil {
		log.Printf("Error decoding request body: %v", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validation
	if demouser.Name == "" {
		http.Error(w, "Name is required", http.StatusBadRequest)
		return
	}

	// Validation
	if demouser.Email == "" {
		http.Error(w, "Email is required", http.StatusBadRequest)
		return
	}

	query := `
        INSERT INTO demo_users (
            id, name, email, created_at
        ) VALUES (
            gen_random_uuid(), $1, $2, CURRENT_TIMESTAMP
        ) RETURNING id`

	var id string
	err := DB.QueryRow(
		query,
		demouser.Name,
		demouser.Email,
	).Scan(&id)

	if err != nil {
		log.Printf("Error creating demouser: %v", err)
		http.Error(w, "Failed to create demouser", http.StatusInternalServerError)
		return
	}

	// Respond
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"id":      id,
		"message": "Demo user submitted successfully",
	})
}

// HandleCreateWaitlist creates a new waitlist entry
func HandleCreateContact(w http.ResponseWriter, r *http.Request) {
	var waitlist struct {
		FirstName   string `json:"first_name"`
		LastName    string `json:"last_name"`
		Email       string `json:"email"`
		Description string `json:"description"`
		Company     string `json:"company"`
		Country     string `json:"country"`
		IntrestedIn string `json:"intrested_in"`
		GetNotify   bool   `json:"get_notify"`
	}

	// Decode request body
	if err := json.NewDecoder(r.Body).Decode(&waitlist); err != nil {
		log.Printf("Error decoding request body: %v", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validation
	if waitlist.FirstName == "" {
		http.Error(w, "Name is required", http.StatusBadRequest)
		return
	}
	if waitlist.Email == "" {
		http.Error(w, "Email is required", http.StatusBadRequest)
		return
	}

	query := `
        INSERT INTO contacts (
            id, first_name, last_name, email, company, country, intrested_in, get_notify, description, created_at
        ) VALUES (
            gen_random_uuid(), $1, $2, $3, $4, $5, $6, $7, $8, CURRENT_TIMESTAMP
        ) RETURNING id`

	var id string
	err := DB.QueryRow(
		query,
		waitlist.FirstName,
		waitlist.LastName,
		waitlist.Email,
		waitlist.Company,
		waitlist.Country,
		waitlist.IntrestedIn,
		waitlist.GetNotify,
		waitlist.Description,
	).Scan(&id)

	if err != nil {
		log.Printf("Error creating contact us: %v", err)
		http.Error(w, "Failed to create contact us", http.StatusInternalServerError)
		return
	}

	// Respond
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"id":      id,
		"message": "Contact submitted successfully",
	})
}
