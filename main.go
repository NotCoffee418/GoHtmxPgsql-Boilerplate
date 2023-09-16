package main

import (
	"html/template"
	"log"
	"net/http"
)

type Page struct {
	Title       string
	Description string
	Body        string
}

func main() {
	//http.ListenAndServe(":8080", http.FileServer(http.Dir("public"))
	//h1 := func(w http.ResponseWriter, r *http.Request) {
	//	templ := template.Must((template.ParseFiles("index.html")))
	//
	//	data := map[string]Page{
	//		"Post": {
	//			Title:       "My first post",
	//			Description: "My first post description",
	//			Body:        "My first post body",
	//		},
	//	}
	//
	//	templ.Execute(w, data)
	}

	// Handlers
	//http.HandleFunc("/", h1)
	//
	//// Start server
	//log.Fatal(http.ListenAndServe(":8080", nil))

	mux := http.NewServeMux()

}
