package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"net/http"
	"time"
	r "gopkg.in/dancannon/gorethink.v2"
	//r "github.com/dancannon/gorethink"
	"strings"
)

func bulkUploadHandler(w http.ResponseWriter, res *http.Request) {

	csvFile, err := os.Open("./data.csv")
	if err != nil {
		fmt.Println(err)
	}
	defer csvFile.Close()
	reader := csv.NewReader(csvFile)
	reader.FieldsPerRecord = -1
	csvData, err := reader.ReadAll()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	var oneRecord Book
//	var genre []string

	for _, each := range csvData {
		oneRecord.Title = each[0]
		oneRecord.Author = each[1]
		oneRecord.Iscale, _ = strconv.ParseFloat(each[2], 32)
		lscale, _ := strconv.ParseFloat(each[3], 32)
		oneRecord.Lscale = 2*lscale
		oneRecord.Genre = strings.ToLower(each[4])
		oneRecord.Created = time.Now()
		_, err := r.Table("books").Insert(oneRecord).RunWrite(session)
		fmt.Println("Writing %v", oneRecord.Title)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		
	}
	fmt.Fprintf(w, "All Uploaded")
}
