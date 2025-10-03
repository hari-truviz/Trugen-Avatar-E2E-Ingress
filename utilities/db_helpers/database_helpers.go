package DBHelpers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	_ "github.com/lib/pq" // registers "postgres"
)

var DB *sql.DB

// InitDB initializes the database connection
func InitDB(connectionString string) error {
	var err error
	DB, err = sql.Open("postgres", connectionString)
	if err != nil {
		log.Printf("Error opening database: %v", err)
		return err
	}

	log.Println("Attempting to ping database...")
	err = DB.Ping()
	if err != nil {
		log.Printf("Error pinging database: %v", err)
		return err
	}

	log.Println("Database connection established successfully.")
	return nil
}

// HandleCreateDemoFeedback creates a new demo feedback entry
func HandleCreateDemoFeedback(w http.ResponseWriter, r *http.Request) {
	feedback := &struct {
		ConversationID string `json:"conversation_id"`
		Email          string `json:"email"`
		Feedback       string `json:"feedback"`
		Rating         int    `json:"rating"`
	}{}

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
		Name    string  `json:"name"`
		Email   string  `json:"email"`
		Company *string `json:"company"`
	}

	// Decode request body
	if err := json.NewDecoder(r.Body).Decode(&waitlist); err != nil {
		log.Printf("Error decoding request body: %v", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validation
	if waitlist.Name == "" {
		http.Error(w, "Name is required", http.StatusBadRequest)
		return
	}

	// Validation
	if waitlist.Email == "" {
		http.Error(w, "Email is required", http.StatusBadRequest)
		return
	}

	query := `
        INSERT INTO waitlists (
            id, name, email, created_at, company
        ) VALUES (
            gen_random_uuid(), $1, $2, CURRENT_TIMESTAMP, $3
        ) RETURNING id`

	var id string
	err := DB.QueryRow(
		query,
		waitlist.Name,
		waitlist.Email,
		waitlist.Company,
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
	if r.Body == nil {
		fmt.Println("Empty request body")
		http.Error(w, "Empty request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()
	demouser := &struct {
		Name  string `json:"name"`
		Email string `json:"email"`
	}{}
	// Decode request body
	if err := json.NewDecoder(r.Body).Decode(demouser); err != nil {
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

// HandleCreatePreviewUser creates a new preview user entry
func HandleCreatePreviewUser(w http.ResponseWriter, r *http.Request) {
	if r.Body == nil {
		fmt.Println("Empty request body")
		http.Error(w, "Empty request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()
	previewuser := &struct {
		Name      string `json:"name"`
		Email     string `json:"email"`
		InviteKey string `json:"invitekey"`
	}{}
	// Decode request body
	if err := json.NewDecoder(r.Body).Decode(previewuser); err != nil {
		log.Printf("Error decoding request body: %v", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validation
	if previewuser.Name == "" {
		http.Error(w, "Name is required", http.StatusBadRequest)
		return
	}

	// Validation
	if previewuser.Email == "" {
		http.Error(w, "Email is required", http.StatusBadRequest)
		return
	}

	// Validation
	if previewuser.InviteKey == "" {
		http.Error(w, "Invite Key is required", http.StatusBadRequest)
		return
	}

	var (
		uuid string
	)

	queryKey := `
		SELECT id 
		FROM preview_invite_keys 
		WHERE id = $1 AND is_active = true
	`
	errFetch := DB.QueryRow(queryKey, previewuser.InviteKey).Scan(&uuid)
	if errFetch != nil {
		http.Error(w, "Not a valid Invite Key", http.StatusNotFound)
		return
	}

	query := `
        INSERT INTO demo_users (
            id, name, email, invite_key, created_at
        ) VALUES (
            gen_random_uuid(), $1, $2, $3, CURRENT_TIMESTAMP
        ) RETURNING id`

	var id string
	err := DB.QueryRow(
		query,
		previewuser.Name,
		previewuser.Email,
		previewuser.InviteKey,
	).Scan(&id)

	if err != nil {
		log.Printf("Error creating preview user: %v", err)
		http.Error(w, "Failed to create preview user", http.StatusInternalServerError)
		return
	}

	// Respond
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"id":      id,
		"message": "Preview user submitted successfully",
	})
}

// HandleCreatePreviewConversation creates a new preview user conversation
func HandleCreatePreviewConversation(w http.ResponseWriter, r *http.Request) {
	if r.Body == nil {
		fmt.Println("Empty request body")
		http.Error(w, "Empty request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()
	previewuser := &struct {
		Email          string `json:"email"`
		InviteKey      string `json:"invitekey"`
		ConversationID string `json:"conversation_id"`
		Status         string `json:"status"`
	}{}
	// Decode request body
	if err := json.NewDecoder(r.Body).Decode(previewuser); err != nil {
		log.Printf("Error decoding request body: %v", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validation
	if previewuser.ConversationID == "" {
		http.Error(w, "Conversation ID is required", http.StatusBadRequest)
		return
	}

	// Validation
	if previewuser.Email == "" {
		http.Error(w, "Email is required", http.StatusBadRequest)
		return
	}

	// Validation
	if previewuser.InviteKey == "" {
		http.Error(w, "Invite Key is required", http.StatusBadRequest)
		return
	}

	var (
		uuid string
	)

	queryKey := `
		SELECT id 
		FROM preview_invite_keys 
		WHERE id = $1 AND is_active = true
	`
	errFetch := DB.QueryRow(queryKey, previewuser.InviteKey).Scan(&uuid)
	if errFetch != nil {
		http.Error(w, "Not a valid Invite Key", http.StatusNotFound)
		return
	}

	query := `
        INSERT INTO preview_conversations (
            id, email, invite_key, conversation_id, created_at, status
        ) VALUES (
            gen_random_uuid(), $1, $2, $3, CURRENT_TIMESTAMP, $4
        ) RETURNING id`

	var id string
	err := DB.QueryRow(
		query,
		previewuser.Email,
		previewuser.InviteKey,
		previewuser.ConversationID,
		previewuser.Status,
	).Scan(&id)

	if err != nil {
		log.Printf("Error creating preview user conversation: %v", err)
		http.Error(w, "Failed to create preview user conversation", http.StatusInternalServerError)
		return
	}

	// Respond
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"id":      id,
		"message": "Preview user Conversation Started",
	})
}

// HandleUpdatePreviewConversation updates an existing preview conversation
func HandleUpdatePreviewConversation(w http.ResponseWriter, r *http.Request) {

	var requestBody struct {
		ConversationId string          `json:"conversation_id"`
		RequestId      string          `json:"request_id"`
		Usage          json.RawMessage `json:"usage"`
		Status         string          `json:"status"`
	}
	convid := &requestBody.ConversationId
	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		log.Printf("Error decoding request body: %v", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	query := `
        UPDATE preview_conversations
        SET status = $1,
            usage = $2,
            updated_at = CURRENT_TIMESTAMP,
			requestid = $3
        WHERE conversation_id = $4
        RETURNING id`

	var id string
	err := DB.QueryRow(
		query,
		requestBody.Status,
		requestBody.Usage,
		requestBody.RequestId,
		convid,
	).Scan(&id)

	if err == sql.ErrNoRows {
		log.Printf("No Conversation found with ID: %s", convid)
		http.Error(w, "Conversation not found", http.StatusNotFound)
		return
	} else if err != nil {
		log.Printf("Error updating Conversation config: %v", err)
		http.Error(w, "Failed to update Conversation configuration", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"id":      id,
		"message": "Conversation updated successfully",
	})
}

// HandleCreatePreviewKey creates a new preview key
func HandleCreatePreviewKey(w http.ResponseWriter, r *http.Request) {
	if r.Body == nil {
		fmt.Println("Empty request body")
		http.Error(w, "Empty request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()
	previewkey := &struct {
		Name string `json:"name"`
	}{}
	// Decode request body
	if err := json.NewDecoder(r.Body).Decode(previewkey); err != nil {
		log.Printf("Error decoding request body: %v", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validation
	if previewkey.Name == "" {
		http.Error(w, "Name is required", http.StatusBadRequest)
		return
	}

	query := `
        INSERT INTO preview_invite_keys (
            id, name, is_active, created_at
        ) VALUES (
            gen_random_uuid(), $1, true, CURRENT_TIMESTAMP
        ) RETURNING id`

	var id string
	err := DB.QueryRow(
		query,
		previewkey.Name,
	).Scan(&id)

	if err != nil {
		log.Printf("Error creating preview invite key: %v", err)
		http.Error(w, "Failed to create preview invite key", http.StatusInternalServerError)
		return
	}

	// Respond
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"id":      id,
		"message": "Preview Invite Key created",
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

// HandleGetPreviewConversation retrieves all Preview Conversation from the database
func HandleGetPreviewConversation(w http.ResponseWriter, r *http.Request) {

	query := `
        SELECT DISTINCT c.id, c.email, name, c.invite_key, conversation_id, c.created_at, usage, COALESCE(status, '') AS status, c.updated_at
        FROM preview_conversations c
		join public.demo_users u on u.email = c.email
        ORDER BY created_at DESC;
    `
	rows, err := DB.Query(query)
	if err != nil {
		log.Printf("Error querying preview conversations: %v", err)
		http.Error(w, "Failed to retrieve preview conversations", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var apiKeys []map[string]interface{}
	for rows.Next() {
		var id, email, name, conversation_id, invite_key, usage, status string
		var createdAt, updatedAt time.Time

		if err := rows.Scan(&id, &email, &name, &invite_key, &conversation_id, &createdAt, &usage, &status, &updatedAt); err != nil {
			log.Printf("Error scanning preview conversations row: %v", err)
			http.Error(w, "Error scanning preview conversations", http.StatusInternalServerError)
			return
		}

		apiKeys = append(apiKeys, map[string]interface{}{
			"id":              id,
			"invite_key":      invite_key,
			"email":           email,
			"name":            name,
			"conversation_id": conversation_id,
			"usage":           usage,
			"status":          status,
			"created":         createdAt.Format("2006-01-02 15:04:05"),
			"updated":         updatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(apiKeys)
}

// HandleGetWaitlist retrieves all waitlists from the database
func HandleGetWaitlist(w http.ResponseWriter, r *http.Request) {

	query := `
        SELECT id, name, email, created_at, COALESCE(company, '') AS company
        FROM waitlists 
        ORDER BY created_at DESC;
    `
	rows, err := DB.Query(query)
	if err != nil {
		log.Printf("Error querying waitlists: %v", err)
		http.Error(w, "Failed to retrieve waitlists", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var apiKeys []map[string]interface{}
	for rows.Next() {
		var id, name, email string
		var company string
		var createdAt time.Time

		if err := rows.Scan(&id, &name, &email, &createdAt, &company); err != nil {
			log.Printf("Error scanning waitlists row: %v", err)
			http.Error(w, "Error scanning waitlists", http.StatusInternalServerError)
			return
		}

		apiKeys = append(apiKeys, map[string]interface{}{
			"id":      id,
			"name":    name,
			"email":   email,
			"company": company,
			"created": createdAt.Format("2006-01-02 15:04:05"),
		})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(apiKeys)
}

// HandleGetDemoUser retrieves all demo user from the database
func HandleGetDemoUser(w http.ResponseWriter, r *http.Request) {

	query := `
        SELECT id, name, email, created_at, COALESCE(invite_key, '') AS invite_key
        FROM demo_users 
        ORDER BY created_at DESC;
    `
	rows, err := DB.Query(query)
	if err != nil {
		log.Printf("Error querying demo users: %v", err)
		http.Error(w, "Failed to retrieve demo users", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var apiKeys []map[string]interface{}
	for rows.Next() {
		var id, name, email string
		var invite_key string
		var createdAt time.Time

		if err := rows.Scan(&id, &name, &email, &createdAt, &invite_key); err != nil {
			log.Printf("Error scanning demo users row: %v", err)
			http.Error(w, "Error scanning demo users", http.StatusInternalServerError)
			return
		}

		apiKeys = append(apiKeys, map[string]interface{}{
			"id":         id,
			"name":       name,
			"invite_key": invite_key,
			"email":      email,
			"created":    createdAt.Format("2006-01-02 15:04:05"),
		})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(apiKeys)
}

// HandleGetDemoFeedback retrieves all demofeedback from the database
func HandleGetDemoFeedback(w http.ResponseWriter, r *http.Request) {

	query := `
        SELECT distinct f.id, f.email, name, feedback, rating, conversation_id, f.created_at
		FROM public.demo_feedbacks f
		join public.demo_users u on u.email = f.email
        ORDER BY created_at DESC;
    `
	rows, err := DB.Query(query)
	if err != nil {
		log.Printf("Error querying demofeedback: %v", err)
		http.Error(w, "Failed to retrieve demofeedback", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var apiKeys []map[string]interface{}
	for rows.Next() {
		var id, email, name, feedback, conversationid string
		var rating int
		var createdAt time.Time

		if err := rows.Scan(&id, &email, &name, &feedback, &rating, &conversationid, &createdAt); err != nil {
			log.Printf("Error scanning demofeedback row: %v", err)
			http.Error(w, "Error scanning demofeedback", http.StatusInternalServerError)
			return
		}

		apiKeys = append(apiKeys, map[string]interface{}{
			"id":             id,
			"email":          email,
			"name":           name,
			"feedback":       feedback,
			"rating":         rating,
			"conversationid": conversationid,
			"created":        createdAt.Format("2006-01-02 15:04:05"),
		})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(apiKeys)
}

// HandleGetContact retrieves all contact us from the database
func HandleGetContact(w http.ResponseWriter, r *http.Request) {

	query := `
        SELECT id, first_name, last_name, email, company, country, intrested_in, get_notify, created_at, description
		FROM public.contacts
        ORDER BY created_at DESC;
    `
	rows, err := DB.Query(query)
	if err != nil {
		log.Printf("Error querying contacts: %v", err)
		http.Error(w, "Failed to retrieve contacts", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var apiKeys []map[string]interface{}
	for rows.Next() {
		var id, first_name, last_name, email, company, country, intrested_in, description string
		var get_notify bool
		var createdAt time.Time

		if err := rows.Scan(&id, &first_name, &last_name, &email, &company, &country, &intrested_in, &get_notify, &createdAt, &description); err != nil {
			log.Printf("Error scanning contacts row: %v", err)
			http.Error(w, "Error scanning contacts", http.StatusInternalServerError)
			return
		}

		apiKeys = append(apiKeys, map[string]interface{}{
			"id":           id,
			"first_name":   first_name,
			"last_name":    last_name,
			"email":        email,
			"company":      company,
			"country":      country,
			"intrested_in": intrested_in,
			"get_notify":   get_notify,
			"description":  description,
			"created":      createdAt.Format("2006-01-02 15:04:05"),
		})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(apiKeys)
}
