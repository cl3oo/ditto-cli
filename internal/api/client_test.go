package api

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestClient_Request_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"status": "ok"}`))
	}))
	defer server.Close()

	client := NewClient(server.URL)
	var resp map[string]string
	err := client.Request("GET", "/test", nil, &resp)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if resp["status"] != "ok" {
		t.Fatalf("Expected status ok, got %v", resp["status"])
	}
}

func TestClient_Request_Unauthorized(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
	}))
	defer server.Close()

	client := NewClient(server.URL)
	err := client.Request("GET", "/secure", nil, nil)
	if !errors.Is(err, ErrUnauthorized) {
		t.Fatalf("Expected ErrUnauthorized, got %v", err)
	}
}

func TestClient_Request_NotFound(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()

	client := NewClient(server.URL)
	err := client.Request("GET", "/missing", nil, nil)
	if !errors.Is(err, ErrNotFound) {
		t.Fatalf("Expected ErrNotFound, got %v", err)
	}
}

func TestClient_Request_APIError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(`{"error": "invalid parameter"}`))
	}))
	defer server.Close()

	client := NewClient(server.URL)
	err := client.Request("POST", "/bad", nil, nil)
	expected := "API error (400): invalid parameter"
	if err == nil || err.Error() != expected {
		t.Fatalf("Expected '%s', got '%v'", expected, err)
	}
}

func TestClient_Login(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"token": "fake-token"}`))
	}))
	defer server.Close()

	client := NewClient(server.URL)
	token, err := client.Login("user", "pass")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if token != "fake-token" {
		t.Errorf("Expected token fake-token, got %s", token)
	}

	if client.Token != "fake-token" {
		t.Errorf("Expected client token to be set, got %s", client.Token)
	}
}

func TestClient_Logging(t *testing.T) {
	logFile := "ditto.log"
	_ = os.Remove(logFile)
	defer os.Remove(logFile)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client := NewClient(server.URL)
	_ = client.Request("GET", "/log-test", nil, nil)

	data, err := os.ReadFile(logFile)
	if err != nil {
		t.Fatalf("Failed to read log file: %v", err)
	}

	if len(data) == 0 {
		t.Fatal("Expected log file to be non-empty")
	}
}
