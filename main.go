package main

import (
	"context"
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"

	"github.com/ghodss/yaml"
	"github.com/gorilla/mux"
	"github.com/tektoncd/pipeline/pkg/apis/pipeline/v1alpha1"
)

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/validate/task/{fileName}", validateTask).Methods("POST")
	http.ListenAndServe(":5001", router)
}

func checkSchema(fileName string) error {
	b, err := ioutil.ReadFile("check.yaml") // just pass the file name
	if err != nil {
		log.Println(err)
	}
	var task v1alpha1.Task
	err = yaml.Unmarshal(b, &task)
	if err != nil {
		log.Println(err)
		return err
	}
	ctx := context.Background()
	task.SetDefaults(ctx)
	if err := task.Validate(ctx); err != nil {
		log.Println(err)
		return err
	}
	return nil
}

// CheckLintErrors checks lint errors based on .yamllint config file
func validateTask(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	fileName := mux.Vars(r)["fileName"]
	file, header, err := r.FormFile(fileName)
	if err != nil {
		log.Println(err)
	}
	defer file.Close()
	f, err := os.OpenFile(header.Filename, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		log.Println(err)
	}
	io.Copy(f, file)
	result := checkLint(fileName)
	err = checkSchema(fileName)
	os.Remove(header.Filename)
	if err != nil {
		json.NewEncoder(w).Encode(err)
	}
	json.NewEncoder(w).Encode(result)
}

func checkLint(fileName string) string {
	cmd := exec.Command("/bin/sh", "validation.sh", fileName)
	result, err := cmd.Output()
	if err != nil {
		log.Println(err)
	}
	return string(result)
}
