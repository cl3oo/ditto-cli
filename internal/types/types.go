package types

import "time"

type User struct {
	ID        string    `json:"id"`
	Username  string    `json:"username"`
	Avatar    string    `json:"avatar"`
	CreatedAt time.Time `json:"created_at"`
}

type UserMin struct {
	ID       string `json:"id"`
	Username string `json:"username"`
}

type CommunityMin struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type Score struct {
	ID             string    `json:"id"`
	VoteCount      int       `json:"vote_count"`
	CommentCount   int       `json:"comment_count"`
	VoteScore      int       `json:"vote_score"`
	SubCount       int       `json:"sub_count"`
	MediaCount     int       `json:"media_count"`
	PostCount      int       `json:"post_count"`
	CommunityCount int       `json:"community_count"`
	AwardCount     int       `json:"award_count"`
	TargetID       string    `json:"target_id"`
	TargetType     int       `json:"target_type"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

type Community struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	AuthorID    string    `json:"author_id"`
	CreatedAt   time.Time `json:"created_at"`
	Scores      Score     `json:"scores"`
}

type Media struct {
	ID   string `json:"id"`
	Type string `json:"type"`
	Name string `json:"name"`
	Path string `json:"path"`
}

type Post struct {
	ID            string       `json:"id"`
	Title         string       `json:"title"`
	Content       string       `json:"content"`
	AuthorID      string       `json:"author_id"`
	Author        UserMin      `json:"author"`
	CommunityID   string       `json:"community_id"`
	Community     CommunityMin `json:"community"`
	Media         []Media      `json:"media,omitempty"`
	Scores        Score        `json:"scores"`
	Locked        bool         `json:"locked"`
	Voted         bool         `json:"voted"`
	VoteDirection int          `json:"vote_direction"`
	CreatedAt     time.Time    `json:"created_at"`
	UpdatedAt     time.Time    `json:"updated_at"`
}

type Comment struct {
	ID            string    `json:"id"`
	Content       string    `json:"content"`
	AuthorID      string    `json:"author_id"`
	Author        UserMin   `json:"author"`
	Scores        Score     `json:"scores"`
	Voted         bool      `json:"voted"`
	VoteDirection int       `json:"vote_direction"`
	CreatedAt     time.Time `json:"created_at"`
	Children      []Comment `json:"children"`
}

type TokenResponse struct {
	Token string `json:"token"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

type Wallet struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	Tokens    int       `json:"tokens"`
	Coins     int       `json:"coins"`
	UpdatedAt time.Time `json:"updated_at"`
}

type MetaDataResult struct {
	ID         string `json:"id"`
	TargetID   string `json:"target_id"`
	TargetType int    `json:"target_type"`
	VoteScore  int    `json:"vote_score"`
	AwardCount int    `json:"award_count"`
}

type Award struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Icon string `json:"icon"`
	Cost int    `json:"cost"`
}

const (
	TargetTypeUser      = 0
	TargetTypeCommunity = 1
	TargetTypePost      = 2
	TargetTypeComment   = 3
)
