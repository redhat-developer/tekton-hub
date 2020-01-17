package models

import "testing"

func TestGetAllResources(t *testing.T) {
	resources := GetAllResources()
	if len(resources) == 0 {
		t.Errorf("Return all tasks Expected: %v , Got: %v", len(resources), 0)
	}
}
