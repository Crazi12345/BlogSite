package main

import (
	"html/template"
	"log"
	"net/http"
)

type Todo struct {
	ID    int
	Title string
	Done  bool
}

func main() {
	var todos = []Todo{
		{1, "Learn Go", false},
		{2, "Build a Todo App", false},
	}
	http.HandleFunc("/todos", func(w http.ResponseWriter, r *http.Request) {
		tmpl, err := template.ParseFiles("todos.html")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		tmpl.Execute(w, todos)
	})
	log.Println("Starting server on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
