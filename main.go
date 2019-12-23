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

// ValidationResponse represents reponse from Validation service
type ValidationResponse struct {
	Status  bool   `json:"status"`
	Message string `json:"message"`
}

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/validate/{type}/{fileName}", validate).Methods("POST")
	http.ListenAndServe(":5001", router)
}
func checkPipelineSchema(fileName string) error {
	b, err := ioutil.ReadFile(fileName)
	if err != nil {
		log.Println(err)
	}
	var pipeline v1alpha1.Pipeline
	err = yaml.Unmarshal(b, &pipeline)
	if err != nil {
		log.Println("Invalid Pipeline schema")
		return err
	}
	ctx := context.Background()
	if err := pipeline.Validate(ctx); err != nil {
		return err
	}
	return nil
}
func checkTaskSchema(fileName string) error {
	b, err := ioutil.ReadFile(fileName)
	if err != nil {
		log.Println(err)
	}
	var task v1alpha1.Task
	err = yaml.Unmarshal(b, &task)
	if err != nil {
		log.Println("Invalid Task schema")
		return err
	}
	ctx := context.Background()
	task.SetDefaults(ctx)
	if err := task.Validate(ctx); err != nil {
		return err
	}
	return nil
}

func validate(w http.ResponseWriter, r *http.Request) {
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
	if result != "Success\n" {
		response := ValidationResponse{false, result}
		json.NewEncoder(w).Encode(response)
		os.Remove(fileName)
		return
	}
	resourceType := mux.Vars(r)["type"]
	if resourceType == "task" {
		err = checkTaskSchema(fileName)
	} else if resourceType == "pipeline" {
		err = checkPipelineSchema(fileName)
	}
	if err != nil {
		response := ValidationResponse{false, err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}
	resp := ValidationResponse{true, "Success"}
	json.NewEncoder(w).Encode(resp)
}

func checkLint(fileName string) string {
	cmd := exec.Command("/bin/sh", "validation.sh", fileName)
	result, err := cmd.Output()
	if err != nil {
		log.Println(err)
	}
	return string(result)
}
