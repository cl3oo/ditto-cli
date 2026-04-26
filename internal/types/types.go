package types

import "time"

type User struct {
	ID        string    `json:"id"`
	Username  string    `json:"username"`
	Avatar    string    `json:"avatar"`
	CreatedAt time.Time `json:"created_at"`
}

type Community struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	CreatorID   string    `json:"creator_id"`
	CreatedAt   time.Time `json:"created_at"`
	MemberCount int       `json:"member_count"`
}

type Post struct {
	ID            string    `json:"id"`
	Title         string    `json:"title"`
	Content       string    `json:"content"`
	AuthorID      string    `json:"author_id"`
	AuthorName    string    `json:"author_name"`
	CommunityID   string    `json:"community_id"`
	CommunityName string    `json:"community_name"`
	Score         int       `json:"score"`
	CreatedAt     time.Time `json:"created_at"`
}

type Comment struct {
	ID         string    `json:"id"`
	Content    string    `json:"content"`
	AuthorID   string    `json:"author_id"`
	AuthorName string    `json:"author_name"`
	Score      int       `json:"score"`
	CreatedAt  time.Time `json:"created_at"`
	Children   []Comment `json:"children"`
}

type TokenResponse struct {
	Token string `json:"token"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}
