package main

import (
	"html/template"
	"net/http"
	"redditclone/pkg/handlers"
	"redditclone/pkg/items"
	"redditclone/pkg/middleware"
	"redditclone/pkg/session"
	"redditclone/pkg/user"

	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

func main() {
	templates := template.Must(template.ParseFiles("template/index.html"))

	sm := session.NewSessionsMem()
	zapLogger, _ := zap.NewProduction()
	defer zapLogger.Sync()
	logger := zapLogger.Sugar()
	userRepo := user.NewUserRepo()
	itemsRepo := items.NewRepo()

	userHandler := &handlers.UserHandler{
		Tmpl:     templates,
		UserRepo: userRepo,
		Logger:   logger,
		Sessions: sm,
	}

	handlers := &handlers.ItemsHandler{
		Tmpl:      templates,
		Logger:    logger,
		ItemsRepo: itemsRepo,
	}

	staticHandler := http.StripPrefix(
		"/static",
		http.FileServer(http.Dir("./template/static")),
	)
	r := mux.NewRouter()
	r.Handle("/static/js/{file}", staticHandler)
	r.Handle("/static/css/{file}", staticHandler)

	r.HandleFunc("/", userHandler.Index)
	r.HandleFunc("/api/register", userHandler.SignUp).Methods("POST")
	r.HandleFunc("/api/login", userHandler.Login).Methods("POST")

	r.HandleFunc("/api/posts/", handlers.List).Methods("GET")
	r.HandleFunc("/api/posts", handlers.Add).Methods("POST")
	r.HandleFunc("/api/post/{id}", handlers.Read).Methods("GET")
	r.HandleFunc("/a/{catName}/{id}", handlers.Read).Methods("GET")

	r.HandleFunc("/api/post/{id}", handlers.Delete).Methods("DELETE")
	r.HandleFunc("/api/post/{id}", handlers.AddComment).Methods("POST")
	r.HandleFunc("/api/post/{id}/upvote", handlers.Upvote).Methods("GET")
	r.HandleFunc("/api/post/{id}/downvote", handlers.Downvote).Methods("GET")
	r.HandleFunc("/api/post/{id}/unvote", handlers.Unvote).Methods("GET")
	r.HandleFunc("/api/posts/{catName}", handlers.Category).Methods("GET")
	r.HandleFunc("/api/user/{username}", handlers.UserItems).Methods("GET")
	r.HandleFunc("/api/post/{postID}/{comID}", handlers.DeleteComment).Methods("DELETE")

	mux := middleware.Auth(sm, r)
	mux = middleware.AccessLog(logger, mux)
	mux = middleware.Panic(mux)

	addr := ":8091"
	http.ListenAndServe(addr, mux)

}
