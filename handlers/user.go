package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	database "mytop100movies/database"
	models "mytop100movies/models"
	"net/http"
)

func CreateUser(w http.ResponseWriter, r *http.Request) {
	var user models.User

	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid Request Body", http.StatusBadRequest)
		return
	}
	query := `INSERT INTO users (username, password) VALUES ($1, $2) RETURNING id;`
	err := database.DB.QueryRow(query, user.Username, user.Password).Scan(&user.ID)

	if err != nil {
		http.Error(w, fmt.Sprintf("Error creating user: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "applicatiom/json")
	json.NewEncoder(w).Encode(user)
}

func GetUsers(w http.ResponseWriter, r *http.Request) {
	db := database.GetDB()

	rows, err := db.Query("SELECT id, username, password FROM users")
	if err != nil {
		http.Error(w, "Failed to fetch users", http.StatusInternalServerError)
		log.Println("Error querying users:", err)
		return
	}

	defer rows.Close()
	//Data Processing

	var users []models.User //creates a slice to hold user
	for rows.Next() {       //Iterates over rows
		var user models.User
		if err := rows.Scan(&user.ID, &user.Username, &user.Password); err != nil {
			http.Error(w, "Error Scanning users", http.StatusInternalServerError)
			log.Println("Error scanning user", err)
			return
		}
		users = append(users, user)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

func DeleteUser(w http.ResponseWriter, r *http.Request) {
	//db := database.GetDB()

	if r.Method != http.MethodDelete {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	//Get ID from query parameter
	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "Missing ID parameter", http.StatusBadRequest)
		return
	}
	//pREPARE THE DELTE QWUERY
	query := `DELETE FROM users WHERE id = $1;`

	//Execute the query
	_, err := database.DB.Exec(query, id)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error Deleting the user: %v", err), http.StatusInternalServerError)
		return
	}
}

func AddMovieToUserTop100(w http.ResponseWriter, r *http.Request) {
	var input struct {
		UserID    int    `json:"user_id"`
		TMDBID    int    `json:"tmdb_id"`
		Title     string `json:"title"`
		Overview  string `json:"overview"`
		PosterURL string `json:"poster_url"`
		Notes     string `json:"notes"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	// Check if user already has this TMDB movie
	var movieID int
	err := database.DB.QueryRow(`
		SELECT id FROM movies WHERE user_id = $1 AND tmdb_id = $2
	`, input.UserID, input.TMDBID).Scan(&movieID)

	if err == sql.ErrNoRows {
		// Insert new movie for this user
		err = database.DB.QueryRow(`
			INSERT INTO movies (tmdb_id, title, overview, poster_url, user_id)
			VALUES ($1, $2, $3, $4, $5)
			RETURNING id
		`, input.TMDBID, input.Title, input.Overview, input.PosterURL, input.UserID).Scan(&movieID)

		if err != nil {
			http.Error(w, fmt.Sprintf("Movie insert failed: %v", err), http.StatusInternalServerError)
			return
		}
	} else if err != nil {
		http.Error(w, "Error checking movie existence", http.StatusInternalServerError)
		return
	}

	// Check if already in user_movies
	var exists bool
	err = database.DB.QueryRow(`
		SELECT EXISTS(SELECT 1 FROM user_movies WHERE user_id = $1 AND movie_id = $2)
	`, input.UserID, movieID).Scan(&exists)

	if err != nil {
		http.Error(w, "DB error", http.StatusInternalServerError)
		return
	}
	if exists {
		http.Error(w, "This movie already exists in user's top 100", http.StatusConflict)
		return
	}

	//Count current entries to assign next available rank

	var count int
	err = database.DB.QueryRow(`SELECT COUNT(*) FROM user_movies WHERE user_id = $1`, input.UserID).Scan(&count)

	if err != nil {
		http.Error(w, "Failed to count current top 100 list", http.StatusInternalServerError)
		return
	}

	if count >= 100 {
		http.Error(w, "Top 100 list is already full", http.StatusBadRequest)
		return
	}
	rank := count + 1

	// Insert into user_movies
	var userMovieID int
	err = database.DB.QueryRow(`
		INSERT INTO user_movies (user_id, movie_id, rank, notes)
		VALUES ($1, $2, $3, $4)
		RETURNING id
	`, input.UserID, movieID, rank, input.Notes).Scan(&userMovieID)

	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to insert user_movie: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"user_movie_id": userMovieID,
		"movie_id":      movieID,
		"rank":          rank,
	})
}

func UpdateUserMovie(w http.ResponseWriter, r *http.Request) {
	var userMovie models.UserMovie

	if err := json.NewDecoder(r.Body).Decode(&userMovie); err != nil {
		http.Error(w, "Invalid Input", http.StatusBadRequest)
		return
	}

	updateQuery := `UPDATE user_movies SET rank = $1, notes = $2 WHERE user_id = $3 AND movie_id=$4`

	result, err := database.DB.Exec(updateQuery, userMovie.Rank, userMovie.Notes, userMovie.UserId, userMovie.MovieID)
	if err != nil {
		http.Error(w, "Update Failed", http.StatusInternalServerError)
		return
	}
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		http.Error(w, "No Record Found To Be Updated", http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Updated Succesfully"))
}

func GetUserTop100(w http.ResponseWriter, r *http.Request) {
	userIDStr := r.URL.Query().Get("user_id")
	if userIDStr == "" {
		http.Error(w, "user_id is required", http.StatusBadRequest)
		return
	}

	query := `
	SELECT 
		um.id AS user_movie_id,
		um.rank,
		um.notes,
		m.id AS movie_id,
		m.tmdb_id,
		m.title,
		m.overview,
		m.poster_url
	FROM user_movies um
	JOIN movies m ON um.movie_id = m.id
	WHERE um.user_id = $1
	ORDER BY um.rank
	`

	rows, err := database.DB.Query(query, userIDStr)
	if err != nil {
		http.Error(w, "Failed to fetch top 100", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var result []map[string]interface{}

	for rows.Next() {
		var (
			userMovieID     int
			rank            int
			notes           string
			movieID         int
			tmdbID          int
			title, overview string
			posterURL       string
		)
		if err := rows.Scan(&userMovieID, &rank, &notes, &movieID, &tmdbID, &title, &overview, &posterURL); err != nil {
			http.Error(w, "Row scan error", http.StatusInternalServerError)
			return
		}

		result = append(result, map[string]interface{}{
			"user_movie_id": userMovieID,
			"rank":          rank,
			"notes":         notes,
			"movie": map[string]interface{}{
				"id":         movieID,
				"tmdb_id":    tmdbID,
				"title":      title,
				"overview":   overview,
				"poster_url": posterURL,
			},
		})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}
