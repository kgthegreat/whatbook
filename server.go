package main

import (
	r "github.com/dancannon/gorethink"
	"log"
	"net/http"
	"fmt"
)

var (
	session *r.Session
)

func init() {
	var err error

	session, err = r.Connect(r.ConnectOpts{
		Address:  "localhost:28015",
		Database: "whatbook",
		MaxOpen:  40,
	})
	if err != nil {
		log.Fatalln(err.Error())
	}
}

func NewServer(addr string) *http.Server {
	// Setup router
	initRouting()

	// Create and start server
	return &http.Server{
		Addr:    addr,
	}
}

func StartServer(server *http.Server) {
	log.Println("Starting server")
	err := server.ListenAndServe()
	if err != nil {
		log.Fatalln("Error: %v", err)
	}
}

func initRouting() {
	cssHandler := http.FileServer(http.Dir("./static/css/"))
//	jsHandler := http.FileServer(http.Dir("./static/js/"))
	imagesHandler := http.FileServer(http.Dir("./static/images/"))
	fontsHandler := http.FileServer(http.Dir("./static/fonts/"))

	http.Handle("/css/", http.StripPrefix("/css/", cssHandler))
	http.Handle("/fonts/", http.StripPrefix("/fonts/", fontsHandler))
//	http.Handle("/js/", http.StripPrefix("/js/", jsHandler))
	http.Handle("/images/", http.StripPrefix("/images/", imagesHandler))

	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/answer", answerHandler)
	http.HandleFunc("/bulkupload", bulkUploadHandler)
}

func indexHandler(w http.ResponseWriter, req *http.Request) {
	renderTemplate(w, "index", map[string]interface{}{
		"ImageUrl": "/images/coetzee.jpg",
		"Title": "Youth",
		"Author": "J.M. Coetzee",
		"Blurb": "The second installment of J. M. Coetzee's fictionalized memoir explores a young man's struggle to experience life to its full intensity and transform it into art.",
	})

} 

func answerHandler(w http.ResponseWriter, req *http.Request) {
	var preference string
	if req.FormValue("7") == "on" {
		preference = "7"
	} else if req.FormValue("6") == "on" {
		preference = "6"
	} else if req.FormValue("5") == "on" {
		preference = "5"
	}

	var data = map[string]interface{}{
		"book_id": "1",
		"preference": preference,
	}

	result, err := r.Table("answers").Insert(data).RunWrite(session)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("*** Insert result: ***")
	fmt.Println("*** Insert result: ** %v", result)

	renderTemplate(w, "index", map[string]interface{}{
		"ImageUrl": "/images/windup-bird.jpg",
		"Title": "The Wind Up Bird Chronicles",
		"Author": "Haruki Murakami",
	})

}
