package api

import (
	"errors"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
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

func TestClient_GetTrendingPosts(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"data": [{"id": "1", "title": "Test Post"}]}`))
	}))
	defer server.Close()

	client := NewClient(server.URL)
	posts, err := client.GetTrendingPosts()
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(posts) != 1 || posts[0].Title != "Test Post" {
		t.Errorf("Unexpected posts: %+v", posts)
	}
}

func TestClient_CreatePost(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Logf("Mock server received: %s %s", r.Method, r.URL.Path)
		path := strings.TrimSuffix(r.URL.Path, "/")
		if r.Method == "GET" && strings.HasSuffix(path, "/communities") {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{"data": [{"id": "c1", "name": "Community"}]}`))
			return
		}
		if r.Method == "POST" && strings.HasSuffix(path, "/posts/c1") && r.URL.Query().Get("type") == "1" {
			w.WriteHeader(http.StatusCreated)
			return
		}
		w.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()

	client := NewClient(server.URL)
	err := client.CreatePost("Title", "Content", "Community")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
}

func TestClient_GetMe(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"id": "u1", "username": "admin"}`))
	}))
	defer server.Close()

	client := NewClient(server.URL)
	user, err := client.GetMe()
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if user.Username != "admin" {
		t.Errorf("Expected username admin, got %s", user.Username)
	}
}

func TestClient_GetPost(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"id": "p1", "title": "Single Post"}`))
	}))
	defer server.Close()

	client := NewClient(server.URL)
	post, err := client.GetPost("p1")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if post.ID != "p1" || post.Title != "Single Post" {
		t.Errorf("Unexpected post: %+v", post)
	}
}

func TestClient_GetComments(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"data": [{"id": "c1", "content": "Comment 1"}]}`))
	}))
	defer server.Close()

	client := NewClient(server.URL)
	comments, err := client.GetComments("p1")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(comments) != 1 || comments[0].Content != "Comment 1" {
		t.Errorf("Unexpected comments: %+v", comments)
	}
}

func TestClient_Logging(t *testing.T) {
	logFile := "ditto.log"
	_ = os.Remove(logFile)
	defer os.Remove(logFile)

	f, err := os.OpenFile(logFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		t.Fatalf("Failed to create log file: %v", err)
	}
	logger := log.New(f, "[TEST] ", log.LstdFlags)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client := NewClientWithLogger(server.URL, logger)
	_ = client.Request("GET", "/log-test", nil, nil)
	f.Close()

	data, err := os.ReadFile(logFile)
	if err != nil {
		t.Fatalf("Failed to read log file: %v", err)
	}

	if len(data) == 0 {
		t.Fatal("Expected log file to be non-empty")
	}
}

func TestClient_UploadMedia(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.Contains(r.Header.Get("Content-Type"), "multipart/form-data") {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	// Create a temporary file to upload
	tmpFile, _ := os.CreateTemp("", "test-upload")
	defer os.Remove(tmpFile.Name())
	_, _ = tmpFile.Write([]byte("test data"))
	tmpFile.Close()

	client := NewClient(server.URL)
	err := client.UploadMedia("p1", 2, tmpFile.Name())
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
}

func TestClient_DownloadMedia(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("test data"))
	}))
	defer server.Close()

	tmpFile := filepath.Join(t.TempDir(), "downloaded")
	client := NewClient(server.URL)
	err := client.DownloadMedia("m1", tmpFile)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	data, _ := os.ReadFile(tmpFile)
	if string(data) != "test data" {
		t.Errorf("Expected test data, got %s", string(data))
	}
}

func TestClient_UpdateCommunity(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "PUT" || !strings.Contains(r.URL.Path, "/communities/c1") {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client := NewClient(server.URL)
	err := client.UpdateCommunity("c1", map[string]interface{}{"title": "New Title"})
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
}

func TestClient_DeletePost(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "DELETE" || !strings.Contains(r.URL.Path, "/posts/p1") {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client := NewClient(server.URL)
	err := client.DeletePost("p1")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
}

func TestClient_GetRandomPosts(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"data": [{"id": "r1", "title": "Random Post"}]}`))
	}))
	defer server.Close()

	client := NewClient(server.URL)
	posts, err := client.GetRandomPosts(1)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(posts) != 1 || posts[0].Title != "Random Post" {
		t.Errorf("Unexpected posts: %+v", posts)
	}
}

func TestClient_CheckFollowStatus(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"status": true}`))
	}))
	defer server.Close()

	client := NewClient(server.URL)
	status, err := client.CheckFollowStatus("u1")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if !status {
		t.Error("Expected status true, got false")
	}
}

func TestClient_GetJoinedCommunities(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"data": [{"id": "c1", "name": "Joined"}]}`))
	}))
	defer server.Close()

	client := NewClient(server.URL)
	comms, err := client.GetJoinedCommunities()
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(comms) != 1 || comms[0].Name != "Joined" {
		t.Errorf("Unexpected communities: %+v", comms)
	}
}

func TestClient_BanUser(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" || !strings.Contains(r.URL.Path, "/communities/c1/ban/u1") {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client := NewClient(server.URL)
	err := client.BanUser("c1", "u1")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
}

func TestClient_LockPost(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" || !strings.Contains(r.URL.Path, "/posts/p1/lock") || r.URL.Query().Get("lock") != "true" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client := NewClient(server.URL)
	err := client.LockPost("p1", true)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
}
