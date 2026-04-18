package handler

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

func CreateConsultationHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req struct {
			FullName      string `json:"full_name"`
			BusinessEmail string `json:"business_email"`
			ServiceType   string `json:"service_type"`
			Message       string `json:"message"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		id := uuid.New().String()

		// --- INTEGRASI NLP (START) ---
		// Kita panggil AI Python untuk menganalisa pesan (Kategori, Skor, dan Urgensi)
		nlpCategory, leadScore, isUrgent := getNLPAnalysis(req.Message, req.ServiceType)
		// --- INTEGRASI NLP (END) ---

		query := `
			INSERT INTO consultation_requests (id, full_name, business_email, service_type, message, nlp_category, lead_score, is_urgent)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		`
		_, err := db.Exec(query, id, req.FullName, req.BusinessEmail, req.ServiceType, req.Message, nlpCategory, leadScore, isUrgent)
		if err != nil {
			log.Printf("Error saving consultation: %v\n", err)
			http.Error(w, "Failed to save request: "+err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]string{"message": "Konsultasi berhasil dikirim", "id": id})
	}
}

// getNLPAnalysis adalah fungsi asisten untuk memanggil NLP Service Python
func getNLPAnalysis(message string, service string) (string, float64, bool) {
	// Alamat API Python (NLP Service)
	url := "http://localhost:5000/api/nlp/predict"

	// Menyiapkan data JSON
	payload, _ := json.Marshal(map[string]string{
		"message": message,
		"service": service,
	})

	// Mengatur batas waktu (timeout) 2 detik agar tidak menghambat sistem utama
	client := http.Client{
		Timeout: 2 * time.Second,
	}

	// Memanggil API (POST)
	resp, err := client.Post(url, "application/json", bytes.NewBuffer(payload))
	if err != nil {
		fmt.Printf("⚠️ NLP Service sedang offline atau bermasalah: %v\n", err)
		return "Manual Check Required", 0.0, false
	}
	defer resp.Body.Close()

	// Membaca hasil dari Python
	var result struct {
		Category        string  `json:"category"`
		Score           float64 `json:"score"`
		UrgencyDetected bool    `json:"urgency_detected"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "Error Profiling", 0.0, false
	}

	return result.Category, result.Score, result.UrgencyDetected
}

// GetConsultationsHandler mengambil semua data konsultasi dari database
func GetConsultationsHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rows, err := db.Query("SELECT id, full_name, business_email, service_type, message, nlp_category, lead_score, is_urgent, status, admin_notes, created_at FROM consultation_requests ORDER BY lead_score DESC")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var results []map[string]interface{}
		for rows.Next() {
			var id string
			var name, email, svc, msg string
			var cat, status, notes sql.NullString
			var score sql.NullFloat64
			var urgent sql.NullBool
			var createdAt time.Time

			if err := rows.Scan(&id, &name, &email, &svc, &msg, &cat, &score, &urgent, &status, &notes, &createdAt); err != nil {
				log.Printf("Error scanning row: %v", err)
				continue // Skip problematic rows instead of failing entirely
			}

			// Handle NullStrings
			finalCat := "General"
			if cat.Valid && cat.String != "" {
				finalCat = cat.String
			}
			finalStatus := "New"
			if status.Valid && status.String != "" {
				finalStatus = status.String
			}
			finalNotes := ""
			if notes.Valid {
				finalNotes = notes.String
			}
			finalScore := 0.0
			if score.Valid {
				finalScore = score.Float64
			}
			finalUrgent := false
			if urgent.Valid {
				finalUrgent = urgent.Bool
			}

			results = append(results, map[string]interface{}{
				"id":             id,
				"full_name":      name,
				"business_email": email,
				"service_type":   svc,
				"message":        msg,
				"nlp_category":   finalCat,
				"lead_score":     finalScore,
				"is_urgent":      finalUrgent,
				"status":         finalStatus,
				"admin_notes":    finalNotes,
				"created_at":     createdAt.Format(time.RFC3339),
			})
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(results)
	}
}

// GetAdminStatsHandler mengambil ringkasan statistik untuk dashboard cards
func GetAdminStatsHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var total, corp, sme, urgent, today int

		db.QueryRow("SELECT COUNT(*) FROM consultation_requests").Scan(&total)
		db.QueryRow("SELECT COUNT(*) FROM consultation_requests WHERE nlp_category = 'Corporate'").Scan(&corp)
		db.QueryRow("SELECT COUNT(*) FROM consultation_requests WHERE nlp_category = 'SME'").Scan(&sme)
		db.QueryRow("SELECT COUNT(*) FROM consultation_requests WHERE is_urgent = true").Scan(&urgent)

		// Query untuk menghitung data yang masuk HARI INI saja
		db.QueryRow("SELECT COUNT(*) FROM consultation_requests WHERE created_at >= CURRENT_DATE").Scan(&today)

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"total":     total,
			"corporate": corp,
			"sme":       sme,
			"urgent":    urgent,
			"today":     today,
		})
	}
}

// UpdateConsultationStatusHandler updates the status and notes of a consultation
func UpdateConsultationStatusHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var input struct {
			ID     int    `json:"id"`
			Status string `json:"status"`
			Notes  string `json:"notes"`
		}

		if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		_, err := db.Exec("UPDATE consultation_requests SET status = $1, admin_notes = $2 WHERE id = $3",
			input.Status, input.Notes, input.ID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

	}
}

// DeleteConsultationHandler handles both single and bulk deletion
func DeleteConsultationHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var input struct {
			IDs []string `json:"ids"`
		}

		if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if len(input.IDs) == 0 {
			http.Error(w, "No IDs provided", http.StatusBadRequest)
			return
		}

		// Delete using IN clause for bulk efficiency
		query := "DELETE FROM consultation_requests WHERE id = ANY($1)"
		_, err := db.Exec(query, pq.Array(input.IDs))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"message": "Data deleted successfully"})
	}
}
