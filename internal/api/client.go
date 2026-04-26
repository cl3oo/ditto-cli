package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/rfcku/ditto/cli/internal/types"
)

type Client struct {
	BaseURL    string
	Token      string
	HTTPClient *http.Client
}

func NewClient(baseURL string) *Client {
	return &Client{
		BaseURL: baseURL,
		HTTPClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (c *Client) SetToken(token string) {
	c.Token = token
}

func (c *Client) Request(method, path string, body interface{}, target interface{}) error {
	var bodyReader io.Reader
	if body != nil {
		jsonBody, err := json.Marshal(body)
		if err != nil {
			return err
		}
		bodyReader = bytes.NewBuffer(jsonBody)
	}

	req, err := http.NewRequest(method, fmt.Sprintf("%s/v1%s", c.BaseURL, path), bodyReader)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	if c.Token != "" {
		req.Header.Set("Authorization", "Bearer "+c.Token)
	}

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		var errResp types.ErrorResponse
		if err := json.NewDecoder(resp.Body).Decode(&errResp); err == nil && errResp.Error != "" {
			return fmt.Errorf("API error (%d): %s", resp.StatusCode, errResp.Error)
		}
		return fmt.Errorf("API error: %d", resp.StatusCode)
	}

	if target != nil {
		return json.NewDecoder(resp.Body).Decode(target)
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
