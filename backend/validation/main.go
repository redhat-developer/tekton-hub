package main

import (
	"context"
	"encoding/json"
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
	log.Println("Successfull Connection")
	router.HandleFunc("/validate/{type}/{id}", validate).Methods("POST")
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
	content, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println(err)
	}
	w.Header().Set("Content-Type", "application/json")
	resourceID := mux.Vars(r)["id"]
	filePath := "resources/" + resourceID + ".yaml"
	_, err = os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE, 0777)
	if err != nil {
		log.Println(err)
	}
	err = ioutil.WriteFile(filePath, content, 0777)
	// result := checkLint(filePath)
	// if result != "Success\n" {
	// 	response := ValidationResponse{false, result}
	// 	json.NewEncoder(w).Encode(response)
	// 	os.Remove(filePath)
	// 	return
	// }
	resourceType := mux.Vars(r)["type"]
	if resourceType == "task" {
		err = checkTaskSchema(filePath)
	} else if resourceType == "pipeline" {
		err = checkPipelineSchema(filePath)
	}
	if err != nil {
		response := ValidationResponse{false, err.Error()}
		json.NewEncoder(w).Encode(response)
		os.Remove(filePath)
		return
	}
	os.Remove(filePath)
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
