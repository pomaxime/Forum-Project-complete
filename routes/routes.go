package routes

import (
	"database/sql"
	"net/http"

	"forum/handlers"
	"forum/middleware"
)

func Register(mux *http.ServeMux, db *sql.DB) {
	authHandler := handlers.NewAuthHandler(db)
	postHandler := handlers.NewPostHandler(db)

	fs := http.FileServer(http.Dir("./static"))
	mux.Handle("/static/", http.StripPrefix("/static/", fs))

	mux.HandleFunc("/", postHandler.Index)
	mux.HandleFunc("/register", authHandler.Register)
	mux.HandleFunc("/login", authHandler.Login)
	mux.HandleFunc("/post/", postHandler.Show)
	mux.HandleFunc("/categories", postHandler.Categories)

	mux.Handle("/comment", middleware.Auth(db, http.HandlerFunc(postHandler.Comment)))
	mux.Handle("/logout", middleware.Auth(db, http.HandlerFunc(authHandler.Logout)))
	mux.Handle("/profile", middleware.Auth(db, http.HandlerFunc(authHandler.Profile)))
	mux.Handle("/post/create", middleware.Auth(db, http.HandlerFunc(postHandler.Create)))
}