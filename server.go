package main

import (
	r "gopkg.in/dancannon/gorethink.v2"
	//r "github.com/dancannon/gorethink"
	"log"
	"net/http"
	"math/rand"
	"math"
)

var (
	session *r.Session
	lscale = 7
	iscale = 5
	nextGenre = "culture"
	genreMap = map[string]int{
		"adventure": 0,
		"thriller": 1,
		"mystery": 2,
		"horror": 3,
		"drama": 4,
		"plays": 5,
		"epics": 6,
		"culture": 7,
		"historical": 8,
		"travel": 9,
		"humor": 10,
		"psychology": 11,
		"dystopia": 12,
		"inspiration": 13,
		"short stories": 14,
		"graphic": 15, 
		"novellas": 16,
		"young adult": 17,
		"fantasy": 18,
		"magic realism": 19,
		"science fiction": 20,
	}

	genreArray = []string{"adventure", 
		"thriller", 
		"mystery", 
		"horror", 
		"drama", 
		"plays", 
		"epics", 
		"culture", 
		"historical", 
		"travel", 
		"humor", 
		"psychology", 
		"philosophy", 
		"dystopia", 
		"inspiration", 
		"short stories", 
		"graphic", 
		"novellas", 
		"young adult", 
		"fantasy", 
		"magic realism", 
		"science fiction"}

	iteration = 1
	//displayedBooks []string
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
//	initRouting()

	cssHandler := http.FileServer(http.Dir("./static/css/"))
//	jsHandler := http.FileServer(http.Dir("./static/js/"))
	imagesHandler := http.FileServer(http.Dir("./static/images/"))
	fontsHandler := http.FileServer(http.Dir("./static/fonts/"))

	http.Handle("/css/", http.StripPrefix("/css/", cssHandler))
	http.Handle("/fonts/", http.StripPrefix("/fonts/", fontsHandler))
//	http.Handle("/js/", http.StripPrefix("/js/", jsHandler))
	http.Handle("/images/", http.StripPrefix("/images/", imagesHandler))

	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/quiz", quizHandler)
	http.HandleFunc("/answer", answerHandler)
	http.HandleFunc("/recommendation", recommendationHandler)
	http.HandleFunc("/bulkupload", bulkUploadHandler)

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

func quizHandler(w http.ResponseWriter, req *http.Request) {
	books := getBooks()
	counter := 0
	var empty []string
	for(len(books) == 0 && counter < 10) {
		log.Println("DB does not have such a book. Trying Again..")
		nextGenre = similarGenre(nextGenre)
		books = getBooks()
		counter = counter + 1
	}
	
	log.Printf("Successfully got %v books\n", len(books))
	if len(books) != 0 {
		book := books[rand.Intn(len(books))]
		log.Printf("We are showing %v of %v and iscale %v as question %v\n", book.Title, book.Genre, book.Iscale, iteration )
		renderTemplate(w, "quiz", book)
	} else {
		
		renderTemplate(w, "error", empty)
	}

//	iteration = iteration + 1


} 

func recommendationHandler(w http.ResponseWriter, req *http.Request) {
	iscale = iscale + 1
	books := getBooks()
	for len(books) == 0 {
		log.Println("DB does not have such a book for rec. Trying Again..")
		nextGenre = similarGenre(nextGenre)
		books = getBooks()
	}
	log.Println("Successfully got a book for rec")
	book := books[rand.Intn(len(books))]
	log.Printf("We are recommending %v of %v and iscale %v\n", book.Title, book.Genre, book.Iscale )
	iscale = 5
	nextGenre = "culture"
	iteration = 1
	renderTemplate(w, "recommendation", book)
}

func getBooks() []Book {
	log.Printf("Trying to get a book of genre %s and iscale %v \n", nextGenre, iscale)
	var books []Book
	query := r.Table("books").Filter(r.Row.Field("Iscale").Eq(iscale).And(r.Row.Field("Lscale").Ge(lscale)).And(r.Row.Field("Genre").Eq(nextGenre)))
	result, err := query.Run(session)
	if err != nil {
//		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Println(err.Error())
//		return
	}
	err = result.All(&books)
	if err != nil {
//		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Println(err.Error())
//		return
	}
	return books
}
func indexHandler(w http.ResponseWriter, req *http.Request) {
	var empty []string
	renderTemplate(w, "index", empty)
}

func answerHandler(w http.ResponseWriter, req *http.Request) {
	var preference string
	id := req.FormValue("id")
	genre := req.FormValue("genre")
	title := req.FormValue("title")
	if req.FormValue("7") == "on" {
		nextGenre = similarGenre(genre)
		iscale = iscale + 1
		preference = "7"
	} else if req.FormValue("6") == "on" {
		nextGenre = similarGenre(genre)
		iscale = iscale + 1
		preference = "6"
	} else if req.FormValue("5") == "on" {
		nextGenre = changeGenre(genre)
		iscale = iscale + 1
		preference = "5"
	} else if req.FormValue("4") == "on" {
		nextGenre = changeGenre(genre)
		iscale = iscale + 1
		preference = "4"
	} else if req.FormValue("3") == "on" {
		nextGenre = similarGenre(genre)
	//	iscale = iscale - 1
		preference = "3"
	} else if req.FormValue("2") == "on" {
		nextGenre = changeGenre(genre)
		iscale = iscale - 1
		preference = "2"
	} else if req.FormValue("1") == "on" {
		nextGenre = changeGenre(genre)
//		iscale = iscale - 1
		preference = "1"
	}

	log.Printf("**** User said %s for %s ****\n", preference, title )

	var data = map[string]interface{}{
		"book_id": id,
		"preference": preference,
	}
	
//	displayedBook = append(displayedBook, id)
	_, err := r.Table("answers").Insert(data).RunWrite(session)
	if err != nil {
		log.Println(err)
	}
	iteration = iteration + 1
	if iteration <= 10 {
		http.Redirect(w, req, "/quiz", http.StatusFound)
	} else {
		http.Redirect(w, req, "/recommendation", http.StatusFound)
	}
}

func changeGenre(genre string) string {
	ran := int(math.Floor(float64(len(genreArray)/2))) - genreMap[genre]
	_, lo := determineHiLo(ran, len(genreArray))
	return genreArray[rand.Intn(5) + lo]
}

func similarGenre(genre string) string {
	ran := genreMap[genre]
	_, lo := determineHiLo(ran, len(genreArray))
	
	return genreArray[rand.Intn(5) + lo]

}

func determineHiLo(num int, max int) (int, int) {
	var hi,lo int
	if num <= 2 {
		lo = num
	} else {
		lo = num - 2
	}

	if num >= max - 2 {
		hi = num
	} else {
		hi = num + 2
	}
	return hi, lo
}
	
	
