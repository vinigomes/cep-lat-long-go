package main

import (
	"github.com/gorilla/mux"
	"html/template"
	"log"
	"net/http"
)

var tpl = template.Must(template.ParseFiles("index.html"))

func main() {
	handleRequests()
}

func handleRequests() {
	myRouter := mux.NewRouter().StrictSlash(true)
	myRouter.HandleFunc("/", homePage)
	myRouter.HandleFunc("/upload-csv-file", uploadCSVFile).Methods("POST")
	log.Fatal(http.ListenAndServe(":10000", myRouter))
}

func homePage(w http.ResponseWriter, r *http.Request){
	ex := new(Exception)
	tpl.Execute(w, ex)
}

func uploadCSVFile(w http.ResponseWriter, r *http.Request) {
	file, handler, err := r.FormFile("file")
	if err != nil {
		throwException(w, err, "Couldn't find the file key in form-data")
		return
	}
	defer file.Close()
	fileName := handler.Filename
	addresses, err := ReadCepFromCsv(fileName)
	if err != nil {
		throwException(w, err,"Couldn't read the file")
		return
	}
	addresses, err = ConvertCsvWithCepToLatitudeLongitude(addresses)
	if err != nil {
		throwException(w, err,"Couldn't convert the csv file")
		return
	}
	err = WriteCsv(addresses)
	if err != nil {
		throwException(w, err,"Couldn't convert the csv file")
		return
	}
	w.Header().Set("Content-Disposition", "attachment; filename="+fileName)
	http.ServeFile(w, r, "output.csv")
}

func throwException(w http.ResponseWriter, err error, message string) {
	log.Printf("fatal error: %s", err)
	ex := new(Exception)
	ex.Status = true
	ex.Message = message
	tpl.Execute(w, ex)
}
