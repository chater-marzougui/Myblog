package main

import (
	"log"
	"myblog/internal/handlers"
	"myblog/internal/models"
	"net/http"
)

func main() {
	models.InitDB()
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("internal/templates"))))

	//https://unsplash.it/800/800/?random

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		handlers.Login(models.DB)(w, r)
	})
	http.HandleFunc("/post/", handlers.ViewPost)
	http.HandleFunc("/posts/", handlers.ListPosts)
	http.HandleFunc("/post/edit/", handlers.EditPost)
	http.HandleFunc("/post/delete/", handlers.DeletePost)
	http.HandleFunc("/new", handlers.NewPost)
	http.HandleFunc("/register", handlers.Register)
	http.HandleFunc("/editProfile", handlers.EditProfile)
	http.HandleFunc("/profile", handlers.ViewProfile)

	log.Println("Starting server on :8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatalf("Could not start server: %s\n", err)
	}
}
