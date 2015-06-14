package main

import (
	"log"
	"net/http"
	"html/template"
)

func indexHandler(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "index")
} 



func renderTemplate(w http.ResponseWriter, tmpl string) {
    t, err := template.ParseFiles(tmpl + ".html")
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
	err = t.Execute(w, nil)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
    }
}

func main() {
	cssHandler := http.FileServer(http.Dir("./css/"))
//	jsHandler := http.FileServer(http.Dir("./js/"))
	imagesHandler := http.FileServer(http.Dir("./images/"))
	fontsHandler := http.FileServer(http.Dir("./fonts/"))

	http.Handle("/css/", http.StripPrefix("/css/", cssHandler))
	http.Handle("/fonts/", http.StripPrefix("/fonts/", fontsHandler))
//	http.Handle("/js/", http.StripPrefix("/js/", jsHandler))
	http.Handle("/images/", http.StripPrefix("/images/", imagesHandler))

	http.HandleFunc("/", indexHandler)
	log.Fatal(http.ListenAndServe(":8081", nil))
}
