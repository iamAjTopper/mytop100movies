package models

type Movie struct {
	ID        int    `json:"id"`
	TMDBID    int    `json:"tmdb_id"`
	Title     string `json:"title"`
	Overview  string `json:"overview"`
	PosterURL string `json:"poster_url"`
	Rank      int    `json:"rank"`    // 1 to 100
	UserId    int    `json:"user_id"` // optional for now
}
