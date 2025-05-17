package main

import (
	"fmt"
	"log"
	setup "mytop100movies/database"
	handlers "mytop100movies/handlers"
	handle "mytop100movies/utils"
	"net/http"

	"github.com/joho/godotenv"
)

func main() {
	setup.InitDB() //connect to database

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	http.HandleFunc("/movies", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			handlers.CreateMovie(w, r)
		case http.MethodGet:
			handlers.GetMovies(w, r)
		case http.MethodPut:
			handlers.UpdateMovie(w, r)
		case http.MethodDelete:
			handlers.DeleteMovie(w, r)
		default:
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	})

	//http.HandleFunc("/movies/update", handlers.UpdateMovie)

	http.HandleFunc("/users", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			handlers.CreateUser(w, r)
		case http.MethodGet:
			handlers.GetUsers(w, r)
		case http.MethodDelete:
			handlers.DeleteUser(w, r)
		}
	})

	http.HandleFunc("/top100/add", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			handlers.AddMovieToUserTop100(w, r)
		}
	})

	http.HandleFunc("/top100/update", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPut {
			handlers.UpdateUserMovie(w, r)
		}
	})

	http.HandleFunc("/top100/get", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			handlers.GetUserTop100(w, r)
		}
	})

	http.HandleFunc("/top100/delete", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodDelete {
			handlers.DeleteUserTop100(w, r)
		}
	})

	http.HandleFunc("/search", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			handle.SearchTMDbMovie(w, r)
		} else {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			return
		}
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Server is running and DB is connected!"))
	})

	fmt.Println("🚀 Server started on port 8080")
	handlerWithCORS := handlers.CorsMiddleware(http.DefaultServeMux)
	log.Fatal(http.ListenAndServe(":8080", handlerWithCORS))
}
