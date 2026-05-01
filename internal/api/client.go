package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httptrace"
	"net/url"
	"os"
	"path/filepath"
	"strings"
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
	home, _ := os.UserHomeDir()
	logPath := filepath.Join(home, ".config", "ditto-cli", "ditto.log")
	if logFile, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666); err == nil {
		logger = log.New(logFile, "[API] ", log.Ldate|log.Ltime|log.Lshortfile)
	}

	return NewClientWithLogger(baseURL, logger)
}

func NewClientWithLogger(baseURL string, logger *log.Logger) *Client {
	return &Client{
		BaseURL: baseURL,
		HTTPClient: &http.Client{
			Timeout: 15 * time.Second,
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				if logger != nil {
					logger.Printf("Redirecting to: %s", req.URL)
				}
				if len(via) >= 10 {
					return errors.New("stopped after 10 redirects")
				}
				return nil
			},
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

	baseURL := strings.TrimSuffix(c.BaseURL, "/")
	path = "/" + strings.TrimPrefix(path, "/")
	
	url := fmt.Sprintf("%s%s", baseURL, path)
	if !strings.HasPrefix(path, "/v1") && !strings.Contains(baseURL, "/v1") {
	        url = fmt.Sprintf("%s/v1%s", baseURL, path)
	}
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

	trace := &httptrace.ClientTrace{
		DNSStart: func(info httptrace.DNSStartInfo) {
			if c.Logger != nil {
				c.Logger.Printf("DNS Start: %s", info.Host)
			}
		},
		DNSDone: func(info httptrace.DNSDoneInfo) {
			if c.Logger != nil {
				c.Logger.Printf("DNS Done: %v, err: %v", info.Addrs, info.Err)
			}
		},
		ConnectStart: func(network, addr string) {
			if c.Logger != nil {
				c.Logger.Printf("Connect Start: %s %s", network, addr)
			}
		},
		ConnectDone: func(network, addr string, err error) {
			if c.Logger != nil {
				c.Logger.Printf("Connect Done: %s %s, err: %v", network, addr, err)
			}
		},
		GotFirstResponseByte: func() {
			if c.Logger != nil {
				c.Logger.Printf("Got First Response Byte")
			}
		},
	}
	req = req.WithContext(httptrace.WithClientTrace(req.Context(), trace))

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

func (c *Client) Register(username, email, password string) (string, error) {
        body := map[string]string{
                "username": username,
                "email":    email,
                "password": password,
        }
        var res types.TokenResponse
        err := c.Request("POST", "/auth/register", body, &res)
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
        err := c.Request("GET", "/posts/", nil, &res)
        return res.Data, err
}

func (c *Client) GetCommunities() ([]types.Community, error) {
        var res struct {
                Data []types.Community `json:"data"`
        }
        err := c.Request("GET", "/communities/", nil, &res)
        return res.Data, err
}

func (c *Client) Vote(targetID string, targetType int, value int) error {
        path := fmt.Sprintf("/votes/%s/?type=%d&value=%d", targetID, targetType, value)
        return c.Request("POST", path, nil, nil)
}

func (c *Client) CreateComment(targetID string, targetType int, content string) error {
        body := map[string]string{"content": content}
        path := fmt.Sprintf("/comments/%s/?type=%d", targetID, targetType)
        return c.Request("POST", path, body, nil)
}

func (c *Client) JoinCommunity(communityID string) error {
        return c.Request("POST", fmt.Sprintf("/communities/%s/join/", communityID), nil, nil)
}

func (c *Client) CreateCommunity(name, title, description string) error {
        body := map[string]string{
                "name":        name,
                "title":       title,
                "description": description,
        }
        return c.Request("POST", "/communities/", body, nil)
}

func (c *Client) CreatePost(title, content, communityName string) error {
	communities, err := c.GetCommunities()
	if err != nil {
		return fmt.Errorf("failed to fetch communities: %w", err)
	}

	var communityID string
	for _, comm := range communities {
		if comm.Name == communityName {
			communityID = comm.ID
			break
		}
	}

	if communityID == "" {
		return fmt.Errorf("community not found: %s", communityName)
	}

	body := map[string]string{
		"title":   title,
		"content": content,
	}
	path := fmt.Sprintf("/posts/%s/?type=1", communityID)
	return c.Request("POST", path, body, nil)
}


func (c *Client) GetFeed() ([]types.Post, error) {
	var res struct {
		Data []types.Post `json:"data"`
	}
	err := c.Request("GET", "/feed/", nil, &res)
	return res.Data, err
}

func (c *Client) GetPost(id string) (*types.Post, error) {
	var post types.Post
	err := c.Request("GET", fmt.Sprintf("/posts/%s/", id), nil, &post)
	return &post, err
}

func (c *Client) GetComments(postID string) ([]types.Comment, error) {
	var res struct {
		Data []types.Comment `json:"data"`
	}
	err := c.Request("GET", fmt.Sprintf("/comments?id=%s&type=2", postID), nil, &res)
	return res.Data, err
}

func (c *Client) GetMe() (*types.User, error) {
	var user types.User
	err := c.Request("GET", "/users/me/", nil, &user)
	return &user, err
}

func (c *Client) GetPostsFiltered(filter map[string]interface{}) ([]types.Post, error) {
	filterJSON, err := json.Marshal(filter)
	if err != nil {
		return nil, err
	}
	path := fmt.Sprintf("/posts/?filter=%s", url.QueryEscape(string(filterJSON)))
	var res struct {
		Data []types.Post `json:"data"`
	}
	err = c.Request("GET", path, nil, &res)
	return res.Data, err
}

func (c *Client) SearchCommunities(query string) ([]types.Community, error) {
	path := fmt.Sprintf("/search?q=%s", url.QueryEscape(query))
	var res struct {
		Data []types.Community `json:"data"`
	}
	err := c.Request("GET", path, nil, &res)
	return res.Data, err
}

func (c *Client) GetUser(id string) (*types.User, error) {
	var user types.User
	err := c.Request("GET", fmt.Sprintf("/users/%s/", id), nil, &user)
	return &user, err
}

func (c *Client) GetCommunity(id string) (*types.Community, error) {
	var community types.Community
	err := c.Request("GET", fmt.Sprintf("/communities/%s/", id), nil, &community)
	return &community, err
}

func (c *Client) GetTrendingCommunities() ([]types.Community, error) {
	var res struct {
		Data []types.Community `json:"data"`
	}
	err := c.Request("GET", "/communities/trending/", nil, &res)
	return res.Data, err
}

func (c *Client) UpdateCommunity(id string, data map[string]interface{}) error {
	return c.Request("PUT", fmt.Sprintf("/communities/%s/", id), data, nil)
}

func (c *Client) DeleteCommunity(id string) error {
	return c.Request("DELETE", fmt.Sprintf("/communities/%s/", id), nil, nil)
}

func (c *Client) GetRandomCommunities(num int) ([]types.Community, error) {
	var res struct {
		Data []types.Community `json:"data"`
	}
	err := c.Request("GET", fmt.Sprintf("/communities/random?num=%d", num), nil, &res)
	return res.Data, err
}

func (c *Client) UpdatePost(id, title, content string) error {
	body := map[string]string{
		"title":   title,
		"content": content,
	}
	return c.Request("PUT", fmt.Sprintf("/posts/%s/", id), body, nil)
}

func (c *Client) DeletePost(id string) error {
	return c.Request("DELETE", fmt.Sprintf("/posts/%s/", id), nil, nil)
}

func (c *Client) GetRandomPosts(num int) ([]types.Post, error) {
	var res struct {
		Data []types.Post `json:"data"`
	}
	err := c.Request("GET", fmt.Sprintf("/posts/random?num=%d", num), nil, &res)
	return res.Data, err
}

func (c *Client) DeleteComment(id string) error {
	return c.Request("DELETE", fmt.Sprintf("/comments/%s/", id), nil, nil)
}

func (c *Client) GetRandomComments(num int) ([]types.Comment, error) {
	var res struct {
		Data []types.Comment `json:"data"`
	}
	err := c.Request("GET", fmt.Sprintf("/comments/random?num=%d", num), nil, &res)
	return res.Data, err
}

func (c *Client) UpdateUser(id string, data map[string]interface{}) error {
	return c.Request("PUT", fmt.Sprintf("/users/%s/", id), data, nil)
}

func (c *Client) DeleteUser(id string) error {
	return c.Request("DELETE", fmt.Sprintf("/users/%s/", id), nil, nil)
}

func (c *Client) GetRandomUsers(num int) ([]types.User, error) {
	var res struct {
		Data []types.User `json:"data"`
	}
	err := c.Request("GET", fmt.Sprintf("/users/random?num=%d", num), nil, &res)
	return res.Data, err
}
