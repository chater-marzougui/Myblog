package main

import (
	"log"
	"myblog/internal/handlers"
	"myblog/internal/models"
	"net/http"
)

var userID int

func main() {
	models.InitDB()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		handlers.Login(models.DB)(w, r)
	})
	http.HandleFunc("/post/", handlers.ViewPost)
	http.HandleFunc("/posts/", handlers.ListPosts)
	http.HandleFunc("/post/edit/", handlers.EditPost)
	http.HandleFunc("/post/delete/", handlers.DeletePost)
	http.HandleFunc("/new", handlers.NewPost)
	http.HandleFunc("/register", handlers.Register)

	log.Println("Starting server on :8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatalf("Could not start server: %s\n", err)
	}
}
