package api

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestClient_Request(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"status": "ok", "data": "test"}`))
	}))
	defer server.Close()

	client := NewClient(server.URL)
	var resp map[string]interface{}
	err := client.Request("GET", "/test", nil, &resp)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if resp["status"] != "ok" {
		t.Fatalf("Expected status ok, got %v", resp["status"])
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
