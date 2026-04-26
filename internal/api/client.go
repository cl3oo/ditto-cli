package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/rfcku/ditto/cli/internal/types"
)

var (
	ErrUnauthorized = errors.New("unauthorized")
	ErrNotFound     = errors.New("not found")
)

type Client struct {
	BaseURL    string
	Token      string
	HTTPClient *http.Client
	Logger     *log.Logger
}

func NewClient(baseURL string) *Client {
	var logger *log.Logger
	if logFile, err := os.OpenFile("ditto.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666); err == nil {
		logger = log.New(logFile, "[API] ", log.Ldate|log.Ltime|log.Lshortfile)
	}

	return &Client{
		BaseURL: baseURL,
		HTTPClient: &http.Client{
			Timeout: 15 * time.Second,
		},
		Logger: logger,
	}
}

func (c *Client) SetToken(token string) {
	c.Token = token
}

func (c *Client) Request(method, path string, body interface{}, target interface{}) error {
	var bodyReader io.Reader
	var bodyBytes []byte
	if body != nil {
		var err error
		bodyBytes, err = json.Marshal(body)
		if err != nil {
			return fmt.Errorf("marshal body: %w", err)
		}
		bodyReader = bytes.NewBuffer(bodyBytes)
	}

	url := fmt.Sprintf("%s/v1%s", c.BaseURL, path)
	req, err := http.NewRequest(method, url, bodyReader)
	if err != nil {
		return fmt.Errorf("create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	if c.Token != "" {
		req.Header.Set("Authorization", "Bearer "+c.Token)
	}

	if c.Logger != nil {
		c.Logger.Printf("--> %s %s", method, url)
		if len(bodyBytes) > 0 {
			c.Logger.Printf("Body: %s", string(bodyBytes))
		}
	}

	start := time.Now()
	resp, err := c.HTTPClient.Do(req)
	duration := time.Since(start)

	if err != nil {
		if c.Logger != nil {
			c.Logger.Printf("<-- %s ERROR: %v (%v)", method, err, duration)
		}
		return fmt.Errorf("do request: %w", err)
	}
	defer resp.Body.Close()

	if c.Logger != nil {
		c.Logger.Printf("<-- %d %s (%v)", resp.StatusCode, method, duration)
	}

	if resp.StatusCode >= 400 {
		if resp.StatusCode == http.StatusUnauthorized {
			return ErrUnauthorized
		}
		if resp.StatusCode == http.StatusNotFound {
			return ErrNotFound
		}

		var errResp types.ErrorResponse
		bodyBits, _ := io.ReadAll(resp.Body)
		if c.Logger != nil && len(bodyBits) > 0 {
			c.Logger.Printf("Error Body: %s", string(bodyBits))
		}

		if err := json.Unmarshal(bodyBits, &errResp); err == nil && errResp.Error != "" {
			return fmt.Errorf("API error (%d): %s", resp.StatusCode, errResp.Error)
		}
		return fmt.Errorf("API error: %d", resp.StatusCode)
	}

	if target != nil {
		if err := json.NewDecoder(resp.Body).Decode(target); err != nil {
			return fmt.Errorf("decode response: %w", err)
		}
	}

	return nil
}

func (c *Client) Login(username, password string) (string, error) {
	body := map[string]string{
		"username": username,
		"password": password,
	}
	var res types.TokenResponse
	err := c.Request("POST", "/auth/authorize", body, &res)
	if err != nil {
		return "", err
	}

	if res.Token != "" {
		c.SetToken(res.Token)
	}
	return res.Token, nil
}

func (c *Client) GetTrendingPosts() ([]types.Post, error) {
	var res struct {
		Data []types.Post `json:"data"`
	}
	err := c.Request("GET", "/posts/trending", nil, &res)
	return res.Data, err
}

func (c *Client) GetCommunities() ([]types.Community, error) {
	var res struct {
		Data []types.Community `json:"data"`
	}
	err := c.Request("GET", "/communities", nil, &res)
	return res.Data, err
}

func (c *Client) Vote(targetID string, targetType int, value int) error {
	path := fmt.Sprintf("/votes/%s?type=%d&value=%d", targetID, targetType, value)
	return c.Request("POST", path, nil, nil)
}

func (c *Client) CreateComment(targetID string, targetType int, content string) error {
	body := map[string]string{"content": content}
	path := fmt.Sprintf("/comments/%s?type=%d", targetID, targetType)
	return c.Request("POST", path, body, nil)
}

func (c *Client) JoinCommunity(communityID string) error {
	return c.Request("POST", fmt.Sprintf("/communities/%s/join", communityID), nil, nil)
}

func (c *Client) CreatePost(title, content, communityName string) error {
	body := map[string]string{
		"title":          title,
		"content":        content,
		"community_name": communityName,
	}
	return c.Request("POST", "/posts", body, nil)
}

func (c *Client) GetFeed() ([]types.Post, error) {
	var res struct {
		Data []types.Post `json:"data"`
	}
	err := c.Request("GET", "/feed", nil, &res)
	return res.Data, err
}

func (c *Client) GetPost(id string) (*types.Post, error) {
	var post types.Post
	err := c.Request("GET", fmt.Sprintf("/posts/%s", id), nil, &post)
	return &post, err
}

func (c *Client) GetComments(postID string) ([]types.Comment, error) {
	var res struct {
		Data []types.Comment `json:"data"`
	}
	err := c.Request("GET", fmt.Sprintf("/posts/%s/comments", postID), nil, &res)
	return res.Data, err
}

func (c *Client) GetMe() (*types.User, error) {
	var user types.User
	err := c.Request("GET", "/users/me", nil, &user)
	return &user, err
}
