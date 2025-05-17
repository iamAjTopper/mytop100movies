package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	database "mytop100movies/database"
	models "mytop100movies/models"
)

func CreateMovie(w http.ResponseWriter, r *http.Request) {
	var movie models.Movie

	//Decode the JSON request body into the `movie` struct

	if err := json.NewDecoder(r.Body).Decode(&movie); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	//Inser into the database

	query := `
		INSERT INTO movies (tmdb_id, title, overview, poster_url, rank, user_id)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id;
	`

	err := database.DB.QueryRow(
		query,
		movie.TMDBID,
		movie.Title,
		movie.Overview,
		movie.PosterURL,
		movie.Rank,
		movie.UserId,
	).Scan(&movie.ID)

	if err != nil {
		http.Error(w, fmt.Sprintf("Datebase Insert Error: %v", err), http.StatusInternalServerError)
		return
	}

	// Respond with the inserted movie (including the generated ID)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(movie)
}

func GetMovies(w http.ResponseWriter, r *http.Request) {
	//	UserId := 1 // hardcoded for now, later you can get it from session/token

	rows, err := database.DB.Query(`SELECT id, tmdb_id, title, overview, poster_url, rank, user_id FROM movies`)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error Querying movies: %v", err), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var movies []models.Movie

	for rows.Next() {
		var movie models.Movie
		err := rows.Scan(&movie.ID, &movie.TMDBID, &movie.Title, &movie.Overview, &movie.PosterURL, &movie.Rank, &movie.UserId)
		if err != nil {
			http.Error(w, "Error Scanning movies", http.StatusInternalServerError)
			return
		}
		movies = append(movies, movie)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(movies)

}

func UpdateMovie(w http.ResponseWriter, r *http.Request) {
	//Extract movie ID from the url
	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "Missing Movie ID", http.StatusBadRequest)
		return
	}

	//Decode the updated movie from the request Body
	var updatedMovie models.Movie
	if err := json.NewDecoder(r.Body).Decode(&updatedMovie); err != nil {
		http.Error(w, "Invaid Request Body", http.StatusBadRequest)
		return
	}

	//Update the movie in the datatbase

	query := `
		UPDATE movies 
		SET tmdb_id=$1, title=$2, overview=$3, poster_url=$4, rank=$5 
		WHERE id=$6 AND user_id=$7
	`

	_, err := database.DB.Exec(
		query,
		updatedMovie.TMDBID,
		updatedMovie.Title,
		updatedMovie.Overview,
		updatedMovie.PosterURL,
		updatedMovie.Rank,
		id,
		1, //hardcoded user_id for now
	)

	if err != nil {
		http.Error(w, fmt.Sprintf("Error Updating Movies: %v", err), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message":"Movie updated succesfully"}`))
}

func DeleteMovie(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}

	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "Missing ID Perameter", http.StatusBadRequest)
		return
	}

	query := `DELETE FROM movies WHERE id = $1;`

	_, err := database.DB.Exec(query, id)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error Deleting the Movie: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("Movie with ID %s is deleted Succesfully", id)))
}
