package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"net/http"
	"time"
)

func bulkUploadHandler(w http.ResponseWriter, res *http.Request) {
	// read data from CSV file

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
	var allRecords []Book
	var genre []string

	for _, each := range csvData {
		oneRecord.Title = each[0]
		oneRecord.Author = each[1]
		oneRecord.Iscale, _ = strconv.ParseFloat(each[2], 32)
		lscale, _ := strconv.ParseFloat(each[3], 32)
		oneRecord.Lscale = 2*lscale
		oneRecord.Genre = append(genre, each[4])
		oneRecord.Created = time.Now()
/*		_, err := r.Table("books").Insert(one).RunWrite(session)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}*/
		
		allRecords = append(allRecords, oneRecord)
	}


	jsondata, err := json.Marshal(allRecords) // convert to JSON

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// sanity check
	// NOTE : You can stream the JSON data to http service as well instead of saving to file
	fmt.Println(string(jsondata))

	// now write to JSON file

	jsonFile, err := os.Create("./data.json")

	if err != nil {
		fmt.Println(err)
	}
	defer jsonFile.Close()

	jsonFile.Write(jsondata)
	jsonFile.Close()
	fmt.Fprintf(w, "All Uploaded")
}
