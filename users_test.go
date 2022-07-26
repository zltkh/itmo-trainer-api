package itmoTrainerApi

import (
	"net/http"
	"testing"
)

func TestEmptyGetUser(t *testing.T) {
	response := GetUser("")
	if response.Body != "Parameter userId is empty" {
		t.Error("Expected \"Parameter userId is empty\", got", response.Body)
	}
}

func TestGetUserNotFound(t *testing.T) {
	response := GetUser("notFoundId")
	if response.StatusCode != http.StatusNotFound {
		t.Error("Expected status code 404, got", response.StatusCode)
	}
}
