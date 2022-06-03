package main

import (
	"Inshorts/controllers"
	"Inshorts/utils"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/articles", utils.SetHeaders(controllers.CreateArticle)).Methods("POST")
	router.HandleFunc("/articles/search", utils.SetHeaders(controllers.SearchArticle)).Methods("GET")
	router.HandleFunc("/articles/{id}", utils.SetHeaders(controllers.GetArticleById)).Methods("GET")
	router.HandleFunc("/articles", utils.SetHeaders(controllers.GetAllArticles)).Methods("GET")

	http.ListenAndServe(":8000", router)
}
