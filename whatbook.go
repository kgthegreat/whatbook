package main

import (
	"log"
	"net/http"
	"fmt"
	r "github.com/dancannon/gorethink"
)



var session *r.Session

type Answer struct {
    Id    string `gorethink:"id,omitempty"`
    Book_Id  string `gorethink:"book_id"`
    Preference string `gorethink:"preference"`
}

func init() {
    var err error
    session, err = r.Connect(r.ConnectOpts{
        Address:  "localhost:28015",
        Database: "whatbook",
    })
    if err != nil {
        fmt.Println(err)
        return
    }
}


func indexHandler(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "index", map[string]interface{}{
		"ImageUrl": "/images/coetzee.jpg",
		"Title": "Youth",
		"Author": "J.M. Coetzee",
	})

} 

func answerHandler(w http.ResponseWriter, res *http.Request) {
	var preference string
	if res.FormValue("yes") == "on" {
		preference = "yes"
		fmt.Println("In if")
	} else if res.FormValue("neutral") == "on" {
		preference = "neutral"
	} else if res.FormValue("no") == "on" {
		preference = "no"
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
	http.HandleFunc("/answer", answerHandler)
	log.Fatal(http.ListenAndServe(":8081", nil))
}
