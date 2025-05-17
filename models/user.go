package models

type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type UserMovie struct {
	ID      int    `json:"id"`
	UserId  int    `json:"user_id"`
	MovieID int    `json:"movie_id"`
	Rank    int    `json:"rank"`
	Notes   string `json:"notes"`
}
