package main

import (
	"bytes"
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"testing"
)

type response struct {
	Status  bool   `json:"status"`
	Message string `json:"message"`
}

func TestValidationSuccess(t *testing.T) {
	err := Copy("valid_task.yaml", "validation_success.yaml")
	if err != nil {
		t.Log(err)
	}
	file, err := os.Open("validation_success.yaml")
	if err != nil {
		t.Log(err)
	}
	defer file.Close()
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("validation_success.yaml", filepath.Base(file.Name()))
	if err != nil {
		t.Log(err)
	}
	io.Copy(part, file)
	writer.Close()
	url := "http://localhost:5001/validate/task/validation_success.yaml"
	request, err := http.NewRequest("POST", url, body)
	if err != nil {
		t.Log(err)
	}
	request.Header.Add("Content-Type", writer.FormDataContentType())
	client := &http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		t.Log(err)
	}
	if status := resp.StatusCode; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
	res := response{}
	err = json.NewDecoder(resp.Body).Decode(&res)
	if err != nil {
		t.Log(err)
	}
	t.Log(res.Message)
	if res.Message != "Success" {
		t.Errorf("Validation failed: got %v want %v",
			res.Message, "Success")
	}
}

func Copy(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	if err != nil {
		return err
	}
	return out.Close()
}
