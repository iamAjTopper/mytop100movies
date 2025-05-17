package utils

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

type TMDbMovie struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	PosterPath  string `json:"poster_path"`
	ReleaseDate string `json:"release_date"`
}

type SearchResponse struct {
	Results []TMDbMovie `json:"results"`
}

type MovieResponse struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	PosterURL   string `json:"poster_url"`
	ReleaseDate string `json:"release_date"`
}

func SearchTMDbMovie(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("query")
	if query == "" {
		http.Error(w, "Missing Query Parameter", http.StatusBadRequest)
		return
	}

	apiKey := os.Getenv("TMDB_API_KEY")
	url := fmt.Sprintf("https://api.themoviedb.org/3/search/movie?api_key=%s&query=%s", apiKey, query)

	fmt.Println("Fetching URL:", url)

	resp, err := http.Get(url)
	if err != nil || resp.StatusCode != http.StatusOK {
		http.Error(w, "Failed to fetch from TMDb", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	var SearchRes SearchResponse

	if err := json.NewDecoder(resp.Body).Decode(&SearchRes); err != nil {
		http.Error(w, "Failed to Decode TMDb Response", http.StatusInternalServerError)
		return
	}

	const imageBaseURL = "https://image.tmdb.org/t/p/w500"
	var results []MovieResponse
	for _, movie := range SearchRes.Results {
		posterUrl := ""
		if movie.PosterPath != "" {
			posterUrl = imageBaseURL + movie.PosterPath
		}

		results = append(results, MovieResponse{

			ID:          movie.ID,
			Title:       movie.Title,
			PosterURL:   posterUrl,
			ReleaseDate: movie.ReleaseDate,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(results)

}
