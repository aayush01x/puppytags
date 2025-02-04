package main

import (
	"log"
	"net/http"
	"os"

	"github.com/pclubiitk/puppylove_tags/db"
	"github.com/pclubiitk/puppylove_tags/router"
	"github.com/pclubiitk/puppylove_tags/services"
)

func main() {
	database := db.InitDB()
	userSimService := services.NewUserSimilarityService(database)
	r := router.NewRouter(userSimService)
	r.RegisterRoutes()
	port := os.Getenv("PORT")
	if port == "" {
		port = "9090"
	}

	log.Printf("Server is running on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
