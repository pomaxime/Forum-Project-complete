package main

import (
	"fmt"
	"log"
	"net/http"

	"forum/database"
	"forum/handlers"
	"forum/middleware"
)

func main() {
	// Initialisation de la base de données
	db, err := database.Init()
	if err != nil {
		log.Fatalf("Erreur initialisation base de données : %v", err)
	}
	defer db.Close()

	// Création des tables
	if err := database.CreateTables(db); err != nil {
		log.Fatalf("Erreur création tables : %v", err)
	}

	// Handlers avec injection de la DB
	authHandler := handlers.NewAuthHandler(db)
	postHandler := handlers.NewPostHandler(db)

	mux := http.NewServeMux()

	// Fichiers statiques
	fs := http.FileServer(http.Dir("./static"))
	mux.Handle("/static/", http.StripPrefix("/static/", fs))

	// Routes publiques
	mux.HandleFunc("/", postHandler.Index)
	mux.HandleFunc("/register", authHandler.Register)
	mux.HandleFunc("/login", authHandler.Login)
	mux.HandleFunc("/post/", postHandler.Show)

	// Routes protégées (middleware auth)
	mux.Handle("/logout", middleware.Auth(db, http.HandlerFunc(authHandler.Logout)))
	mux.Handle("/profile", middleware.Auth(db, http.HandlerFunc(authHandler.Profile)))
	mux.Handle("/post/create", middleware.Auth(db, http.HandlerFunc(postHandler.Create)))

	fmt.Println("Serveur démarré sur http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
