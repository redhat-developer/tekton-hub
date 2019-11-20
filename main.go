package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"

	"github.com/gorilla/mux"
)

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/lint/{fileName}", checkLintErrors).Methods("POST")
	http.ListenAndServe(":5001", router)
}

// CheckLintErrors checks lint errors based on .yamllint config file
func checkLintErrors(w http.ResponseWriter, r *http.Request) {
	fileName := mux.Vars(r)["fileName"]
	file, header, err := r.FormFile(fileName)
	if err != nil {
		log.Fatalln(err)
	}
	defer file.Close()
	f, err := os.OpenFile(header.Filename, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		log.Fatalln(err)
	}
	io.Copy(f, file)
	result := check(fileName)
	os.Remove(header.Filename)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

func check(fileName string) string {
	cmd := exec.Command("/bin/sh", "validation.sh", fileName)
	result, err := cmd.Output()
	if err != nil {
		log.Fatalln(err)
	}
	return string(result)
}
