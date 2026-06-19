package main

import (
	"fmt"
	"log"
	"time"
	"net/http"

	"forum/database"
	"forum/handlers"
	"forum/middleware"
)

func main() {
	db, err := database.Init()
	if err != nil {
		log.Fatalf("Erreur initialisation base de données : %v", err)
	}
	defer db.Close()

	if err := database.CreateTables(db); err != nil {
		log.Fatalf("Erreur création tables : %v", err)
	}

	authHandler := handlers.NewAuthHandler(db)
	postHandler := handlers.NewPostHandler(db)
	likeHandler := handlers.NewLikeHandler(db)

	mux := http.NewServeMux()

	fs := http.FileServer(http.Dir("./static"))
	mux.Handle("/static/", http.StripPrefix("/static/", fs))

	mux.HandleFunc("/", postHandler.Index)
	mux.HandleFunc("/register", authHandler.Register)
	mux.HandleFunc("/login", authHandler.Login)
	mux.HandleFunc("/post/", postHandler.Show)
	mux.HandleFunc("/categories", postHandler.Categories)

	// Routes protégées (middleware auth)
	mux.Handle("/comment", middleware.Auth(db, http.HandlerFunc(postHandler.Comment)))
	mux.Handle("/comment/like", middleware.Auth(db, http.HandlerFunc(likeHandler.LikeComment)))
	mux.Handle("/comment/dislike", middleware.Auth(db, http.HandlerFunc(likeHandler.DislikeComment)))
	mux.Handle("/like", middleware.Auth(db, http.HandlerFunc(likeHandler.Like)))
	mux.Handle("/dislike", middleware.Auth(db, http.HandlerFunc(likeHandler.Dislike)))
	mux.Handle("/logout", middleware.Auth(db, http.HandlerFunc(authHandler.Logout)))
	mux.Handle("/profile", middleware.Auth(db, http.HandlerFunc(authHandler.Profile)))
	mux.Handle("/post/create", middleware.Auth(db, http.HandlerFunc(postHandler.Create)))

	fmt.Println("Serveur démarré sur http://localhost:8080")
    srv := &http.Server{
        Addr:         ":8080",
        Handler:      mux,
        ReadTimeout:  5 * time.Second,
        WriteTimeout: 10 * time.Second,
    }
    if err := srv.ListenAndServe(); err != nil {
        log.Fatalf("Erreur serveur : %v", err)
    }
}
