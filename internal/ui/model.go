package ui

import (
	"errors"
	"fmt"
	"os/exec"
	"runtime"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/rfcku/ditto/cli/internal/api"
	"github.com/rfcku/ditto/cli/internal/config"
	"github.com/rfcku/ditto/cli/internal/types"
	"github.com/rfcku/ditto/cli/internal/ui/views"
)

type State int

const (
	StateLoading State = iota
	StateFeed
	StateLogin
	StatePostDetail
	StateCommunities
	StateCreatePost
	StateSelection
	StateRegister
	StateEditPost
	StateEditCommunity
	StateProfileSettings
)

type MainModel struct {
	State         State
	Client        *api.Client
	Config        *config.Config
	Theme         Theme
	Keys          KeyMap
	Error         error
	StatusMessage string
	StatusTimer   int
	Width         int
	Height        int
	CommandBuffer string
	Me            *types.User

	// Sub-models
	FeedModel       views.FeedModel
	LoginModel      views.LoginModel
	RegisterModel   views.RegisterModel
	PostDetailModel views.PostDetailModel
	CommunityModel  views.CommunityModel
	CreatePostModel views.CreatePostModel
	EditPostModel      views.CreatePostModel
	EditCommunityModel views.CreatePostModel
	SettingsModel      views.SettingsModel
}

func NewMainModel(cfg *config.Config) MainModel {
	m := MainModel{
		State:           StateLoading,
		Client:          api.NewClient(cfg.BaseURL),
		Config:          cfg,
		Theme:           NewTheme(cfg.Appearance),
		Keys:            NewKeyMap(cfg.Keys),
		FeedModel:       views.NewFeedModel(),
		LoginModel:      views.NewLoginModel(),
		RegisterModel:   views.NewRegisterModel(),
		PostDetailModel: views.NewPostDetailModel(),
		CommunityModel:  views.NewCommunityModel(),
		CreatePostModel: views.NewCreatePostModel(),
		EditPostModel:      views.NewCreatePostModel(),
		EditCommunityModel: views.NewCreatePostModel(),
		SettingsModel:      views.NewSettingsModel(),
	}
	m.Client.SetToken(cfg.Token)
	m.FeedModel.SetTheme(m.Theme.Selected)
	m.CommunityModel.SetTheme(m.Theme.Selected)
	m.PostDetailModel.SetTheme(m.Theme.Accent, m.Theme.Selected, m.Theme.Markdown)
	return m
}

func (m MainModel) Init() tea.Cmd {
	var cmds []tea.Cmd
	cmds = append(cmds, m.fetchFeed())
	cmds = append(cmds, m.LoginModel.Init())
	if m.Client.Token != "" {
		cmds = append(cmds, m.fetchMe())
	}
	return tea.Batch(cmds...)
}

type tickMsg struct{}
type statusMsg string

func (m MainModel) clearStatus() tea.Cmd {
	return tea.Tick(3*time.Second, func(t time.Time) tea.Msg {
		return tickMsg{}
	})
}

func (m MainModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tickMsg:
		m.StatusMessage = ""
		return m, nil

	case statusMsg:
		m.StatusMessage = string(msg)
		var refreshCmd tea.Cmd
		if m.State == StateFeed {
			refreshCmd = m.fetchFeed()
		} else if m.State == StatePostDetail {
			refreshCmd = m.fetchPostDetail(m.PostDetailModel.Post.ID)
		}
		return m, tea.Batch(refreshCmd, m.clearStatus())

	case tea.KeyMsg:
		if m.Error != nil {
			m.Error = nil
		}

		// Handle command buffer first (prefix :)
		if m.CommandBuffer != "" {
			if msg.String() == "enter" {
				fullCmd := m.CommandBuffer
				m.CommandBuffer = ""
				parts := strings.Fields(fullCmd)
				if len(parts) == 0 {
					return m, nil
				}
				cmd := parts[0]

				switch cmd {
				case ":q", ":quit":
					return m, tea.Quit
				case ":L", ":login":
					m.State = StateLogin
				case ":F", ":feed":
					m.State = StateLoading
					return m, m.fetchFeed()
				case ":C", ":communities":
					m.State = StateLoading
					return m, m.fetchCommunities()
				case ":n", ":new":
					m.State = StateCreatePost
				case ":r", ":refresh":
					m.State = StateLoading
					if m.State == StateFeed {
						return m, m.fetchFeed()
					} else if m.State == StateCommunities {
						return m, m.fetchCommunities()
					} else if m.State == StatePostDetail {
						return m, m.fetchPostDetail(m.PostDetailModel.Post.ID)
					}
					m.State = StateFeed // default to feed
					return m, m.fetchFeed()
				case ":delete":
					if m.State == StatePostDetail {
						return m, m.performDeletePost(m.PostDetailModel.Post.ID)
					}
					if m.State == StateCommunities {
						if item, ok := m.CommunityModel.List.SelectedItem().(views.CommunityItem); ok {
							return m, m.performDeleteCommunity(item.ID)
						}
					}
				case ":edit":
					if m.State == StatePostDetail {
						m.EditPostModel.Title.SetValue(m.PostDetailModel.Post.Title)
						m.EditPostModel.Content.SetValue(m.PostDetailModel.Post.Content)
						m.EditPostModel.CommunityID.SetValue(m.PostDetailModel.Post.Community.Name)
						m.EditPostModel.CommunityID.Blur()
						m.EditPostModel.Title.Focus()
						m.State = StateEditPost
						return m, nil
					}
					if m.State == StateCommunities {
						if item, ok := m.CommunityModel.List.SelectedItem().(views.CommunityItem); ok {
							m.EditCommunityModel.Title.SetValue(item.Community.Title)
							m.EditCommunityModel.CommunityID.SetValue(item.Community.Name)
							m.EditCommunityModel.Content.SetValue(item.Community.Description)
							m.EditCommunityModel.CommunityID.Blur()
							m.EditCommunityModel.Title.Focus()
							m.State = StateEditCommunity
							return m, nil
						}
					}
				case ":random":
					prevState := m.State
					m.State = StateLoading
					if prevState == StateFeed {
						return m, m.fetchRandomPosts()
					} else if prevState == StateCommunities {
						return m, m.fetchRandomCommunities()
					}
					m.State = prevState
				case ":follow":
					if m.State == StatePostDetail {
						return m, m.performToggleFollow(m.PostDetailModel.Post.Author.ID)
					}
				case ":following":
					m.State = StateLoading
					return m, m.fetchFollowedUsers()
				case ":joined":
					m.State = StateLoading
					return m, m.fetchJoinedCommunities()
				case ":delete-comment":
					if m.State == StatePostDetail && len(parts) > 1 {
						return m, m.performDeleteComment(parts[1])
					}
				case ":settings":
					if m.Me != nil {
						m.SettingsModel.Avatar.SetValue(m.Me.Avatar)
					}
					m.State = StateProfileSettings
					return m, nil
				case ":delete-account":
					if m.Me != nil {
						return m, m.performDeleteUser(m.Me.ID)
					}
				case ":cc":
					if m.State == StatePostDetail {
						m.PostDetailModel.SetShowCommentInput(true)
						return m, nil
					}
				case ":s", ":search":
					if len(parts) > 1 {
						query := strings.Join(parts[1:], " ")
						m.State = StateLoading
						return m, m.performSearch(query)
					}
				case ":u", ":user":
					if len(parts) > 1 {
						userID := parts[1]
						m.State = StateLoading
						// return m, m.fetchUserProfile(userID) // Placeholder for Step 3
						m.StatusMessage = "User profile view coming soon for: " + userID
						m.State = StateFeed
						return m, m.clearStatus()
					}
				default:
					m.StatusMessage = "Unknown command: " + cmd
					return m, m.clearStatus()
				}
				return m, nil
			}
			
			if msg.String() == "backspace" {
				if len(m.CommandBuffer) > 1 {
					m.CommandBuffer = m.CommandBuffer[:len(m.CommandBuffer)-1]
				} else {
					m.CommandBuffer = ""
				}
				return m, nil
			}

			if msg.String() == "esc" {
				m.CommandBuffer = ""
				return m, nil
			}

			// Add space support for commands with arguments
			if msg.String() == " " {
				m.CommandBuffer += " "
				return m, nil
			}

			if len(msg.String()) == 1 {
				m.CommandBuffer += msg.String()
			}
			return m, nil
		}

		if msg.String() == ":" {
			m.CommandBuffer = ":"
			return m, nil
		}

		if m.State == StateSelection {
			switch msg.String() {
			case "p":
				m.State = StateCreatePost
				return m, nil
			case "c":
				m.State = StateCommunities
				return m, nil
			case "esc":
				m.State = StateFeed
				return m, nil
			}
		}

		// Basic Navigation Keys only
		switch {
		case key.Matches(msg, m.Keys.Quit):
			return m, tea.Quit
		case key.Matches(msg, m.Keys.Back):
			if msg.String() == "backspace" {
				// Don't go back if we are typing in an input
				if (m.State == StatePostDetail && m.PostDetailModel.ShowCommentInput) ||
					m.State == StateLogin || m.State == StateRegister || m.State == StateCreatePost {
					break
				}
			}

			if m.State == StatePostDetail {
				if m.PostDetailModel.ShowCommentInput {
					m.PostDetailModel.SetShowCommentInput(false)
					return m, nil
				}
				m.State = StateFeed
				return m, nil
			}

			if m.State == StateCommunities || m.State == StateCreatePost || m.State == StateSelection || m.State == StateLogin || m.State == StateRegister {
				m.State = StateFeed
				return m, nil
			}
		case msg.String() == "enter":
			if m.State == StateLogin {
				if m.LoginModel.Focused == 2 {
					m.State = StateLoading
					return m, m.performLogin()
				}
				if m.LoginModel.Focused == 3 {
					m.State = StateRegister
					return m, nil
				}
			}
			if m.State == StateRegister {
				if m.RegisterModel.Focused == 4 {
					if m.RegisterModel.Password.Value() != m.RegisterModel.Confirm.Value() {
						m.RegisterModel.Error = "Passwords do not match"
						return m, nil
					}
					m.State = StateLoading
					return m, m.performRegister()
				}
				if m.RegisterModel.Focused == 5 {
					m.State = StateLogin
					return m, nil
				}
			}
			if m.State == StateCreatePost && m.CreatePostModel.Focused == 3 {
				m.State = StateLoading
				return m, m.performCreatePost()
			}
			if m.State == StateEditPost && m.EditPostModel.Focused == 3 {
				m.State = StateLoading
				return m, m.performUpdatePost()
			}
			if m.State == StateEditCommunity && m.EditCommunityModel.Focused == 3 {
				m.State = StateLoading
				return m, m.performUpdateCommunity()
			}
			if m.State == StateProfileSettings && m.SettingsModel.Focused == 1 {
				m.State = StateLoading
				return m, m.performUpdateUser()
			}
			if m.State == StateFeed {
				if item, ok := m.FeedModel.List.SelectedItem().(views.PostItem); ok {
					m.State = StateLoading
					return m, m.fetchPostDetail(item.ID)
				}
			}
			if m.State == StateCommunities {
				if item, ok := m.CommunityModel.List.SelectedItem().(views.CommunityItem); ok {
					m.State = StateLoading
					return m, m.fetchCommunityDetail(item.ID)
				}
			}
			if m.State == StatePostDetail {
				if m.PostDetailModel.ShowCommentInput && m.PostDetailModel.CommentInput.Value() != "" {
					return m, m.performCreateComment()
				}
			}
		}

	case tea.WindowSizeMsg:
		m.Width = msg.Width
		m.Height = msg.Height
		m.FeedModel.SetSize(msg.Width, msg.Height-4)
		m.PostDetailModel.SetSize(msg.Width, msg.Height-4)
		m.CommunityModel.SetSize(msg.Width, msg.Height-4)

	case feedMsg:
		m.State = StateFeed
		m.FeedModel.SetPosts(msg)
		return m, nil

	case communitiesMsg:
		m.State = StateCommunities
		m.CommunityModel.SetCommunities(msg)
		return m, nil

	case usersMsg:
		m.State = StateCommunities
		m.CommunityModel.SetUsers(msg)
		return m, nil

	case postDetailMsg:
		m.State = StatePostDetail
		m.PostDetailModel.SetContent(msg.post, msg.comments)
		m.PostDetailModel.SetShowCommentInput(false)
		return m, nil

	case loginSuccessMsg:
		m.State = StateFeed
		m.Client.SetToken(string(msg))
		
		// Persist token
		m.Config.UpdateToken(string(msg))
		
		return m, tea.Batch(m.fetchFeed(), m.fetchMe())

	case meMsg:
		m.Me = (*types.User)(msg)
		return m, nil

	case commentSuccessMsg:
		m.StatusMessage = string(msg)
		m.PostDetailModel.CommentInput.SetValue("")
		m.PostDetailModel.SetShowCommentInput(false)
		return m, tea.Batch(
			m.fetchPostDetail(m.PostDetailModel.Post.ID),
			m.clearStatus(),
		)

	case errorMsg:
		if errors.Is(msg, api.ErrUnauthorized) {
			if m.Client.Token != "" {
				m.StatusMessage = "Session expired, please login again"
				m.Client.SetToken("")
				m.Config.UpdateToken("")
			}
			if m.State != StateLogin && m.State != StateRegister {
				m.State = StateLogin
			}
			return m, m.clearStatus()
		}
		m.StatusMessage = fmt.Sprintf("Error: %v", msg)
		if m.State == StateLoading {
			m.State = StateFeed
		}
		return m, m.clearStatus()
	}

	// Delegate to sub-models
	switch m.State {
	case StateFeed:
		m.FeedModel, cmd = m.FeedModel.Update(msg)
		cmds = append(cmds, cmd)
	case StateLogin:
		m.LoginModel, cmd = m.LoginModel.Update(msg)
		cmds = append(cmds, cmd)
	case StatePostDetail:
		m.PostDetailModel, cmd = m.PostDetailModel.Update(msg)
		cmds = append(cmds, cmd)
	case StateCommunities:
		m.CommunityModel, cmd = m.CommunityModel.Update(msg)
		cmds = append(cmds, cmd)
	case StateCreatePost:
		m.CreatePostModel, cmd = m.CreatePostModel.Update(msg)
		cmds = append(cmds, cmd)
	case StateEditPost:
		m.EditPostModel, cmd = m.EditPostModel.Update(msg)
		cmds = append(cmds, cmd)
	case StateEditCommunity:
		m.EditCommunityModel, cmd = m.EditCommunityModel.Update(msg)
		cmds = append(cmds, cmd)
	case StateProfileSettings:
		m.SettingsModel, cmd = m.SettingsModel.Update(msg)
		cmds = append(cmds, cmd)
	case StateRegister:
		m.RegisterModel, cmd = m.RegisterModel.Update(msg)
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

type feedMsg []types.Post
type communitiesMsg []types.Community
type usersMsg []types.User
type postDetailMsg struct {
	post     types.Post
	comments []types.Comment
}
type loginSuccessMsg string
type commentSuccessMsg string
type meMsg *types.User
type errorMsg error

func (m MainModel) fetchFeed() tea.Cmd {
	return func() tea.Msg {
		posts, err := m.Client.GetTrendingPosts()
		if err != nil {
			return errorMsg(err)
		}
		return feedMsg(posts)
	}
}

func (m MainModel) fetchMe() tea.Cmd {
	return func() tea.Msg {
		user, err := m.Client.GetMe()
		if err != nil {
			return errorMsg(err)
		}
		return meMsg(user)
	}
}

func (m MainModel) fetchCommunities() tea.Cmd {
	return func() tea.Msg {
		communities, err := m.Client.GetCommunities()
		if err != nil {
			return errorMsg(err)
		}
		return communitiesMsg(communities)
	}
}

func (m MainModel) fetchPostDetail(id string) tea.Cmd {
	return func() tea.Msg {
		post, err := m.Client.GetPost(id)
		if err != nil {
			return errorMsg(err)
		}
		comments, err := m.Client.GetComments(id)
		if err != nil {
			// Still show post even if comments fail
			return postDetailMsg{post: *post, comments: []types.Comment{}}
		}
		return postDetailMsg{post: *post, comments: comments}
	}
}

func (m MainModel) performVote(id string, targetType int, value int) tea.Cmd {
	return func() tea.Msg {
		err := m.Client.Vote(id, targetType, value)
		if err != nil {
			return errorMsg(err)
		}

		msg := "Upvoted!"
		if value < 0 {
			msg = "Downvoted!"
		}

		return statusMsg(msg)
	}
}

func (m MainModel) performCreatePost() tea.Cmd {
	return func() tea.Msg {
		err := m.Client.CreatePost(
			m.CreatePostModel.Title.Value(),
			m.CreatePostModel.Content.Value(),
			m.CreatePostModel.CommunityID.Value(),
		)
		if err != nil {
			return errorMsg(err)
		}
		return m.fetchFeed()()
	}
}

func (m MainModel) performLogin() tea.Cmd {
	return func() tea.Msg {
		token, err := m.Client.Login(m.LoginModel.Username.Value(), m.LoginModel.Password.Value())
		if err != nil {
			return errorMsg(err)
		}
		return loginSuccessMsg(token)
	}
}

func (m MainModel) openBrowser(url string) tea.Cmd {
	return func() tea.Msg {
		var err error
		switch runtime.GOOS {
		case "linux":
			err = exec.Command("xdg-open", url).Start()
		case "windows":
			err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
		case "darwin":
			err = exec.Command("open", url).Start()
		}
		if err != nil {
			return errorMsg(err)
		}
		return nil
	}
}

func (m MainModel) performSearch(query string) tea.Cmd {
	return func() tea.Msg {
		communities, err := m.Client.SearchCommunities(query)
		if err != nil {
			return errorMsg(err)
		}
		return communitiesMsg(communities)
	}
}

func (m MainModel) fetchCommunityDetail(communityID string) tea.Cmd {
	return func() tea.Msg {
		filter := map[string]interface{}{
			"target_id":   communityID,
			"target_type": 1, // Community
		}
		posts, err := m.Client.GetPostsFiltered(filter)
		if err != nil {
			return errorMsg(err)
		}
		return feedMsg(posts)
	}
}

func (m MainModel) performCreateComment() tea.Cmd {
	return func() tea.Msg {
		err := m.Client.CreateComment(
			m.PostDetailModel.Post.ID,
			2, // Post type
			m.PostDetailModel.CommentInput.Value(),
		)
		if err != nil {
			return errorMsg(err)
		}
		
		return commentSuccessMsg("Comment posted!")
	}
}

func (m MainModel) performRegister() tea.Cmd {
	return func() tea.Msg {
		token, err := m.Client.Register(
			m.RegisterModel.Username.Value(),
			m.RegisterModel.Email.Value(),
			m.RegisterModel.Password.Value(),
		)
		if err != nil {
			return errorMsg(err)
		}
		return loginSuccessMsg(token)
	}
}

func (m MainModel) performUpdatePost() tea.Cmd {
	return func() tea.Msg {
		err := m.Client.UpdatePost(
			m.PostDetailModel.Post.ID,
			m.EditPostModel.Title.Value(),
			m.EditPostModel.Content.Value(),
		)
		if err != nil {
			return errorMsg(err)
		}
		return m.fetchPostDetail(m.PostDetailModel.Post.ID)()
	}
}

func (m MainModel) performDeletePost(id string) tea.Cmd {
	return func() tea.Msg {
		err := m.Client.DeletePost(id)
		if err != nil {
			return errorMsg(err)
		}
		return statusMsg("Post deleted")
	}
}

func (m MainModel) performUpdateCommunity() tea.Cmd {
	return func() tea.Msg {
		// Map from our "CreatePostModel" fields to Community fields
		data := map[string]interface{}{
			"title":       m.EditCommunityModel.Title.Value(),
			"description": m.EditCommunityModel.Content.Value(),
		}
		// Community Name is typically not editable after creation in many systems, 
		// but we used CommunityID field for it.
		
		var communityID string
		if item, ok := m.CommunityModel.List.SelectedItem().(views.CommunityItem); ok {
			communityID = item.ID
		} else {
			return errorMsg(errors.New("no community selected"))
		}

		err := m.Client.UpdateCommunity(communityID, data)
		if err != nil {
			return errorMsg(err)
		}
		return m.fetchCommunities()()
	}
}

func (m MainModel) performDeleteCommunity(id string) tea.Cmd {
	return func() tea.Msg {
		err := m.Client.DeleteCommunity(id)
		if err != nil {
			return errorMsg(err)
		}
		return statusMsg("Community deleted")
	}
}

func (m MainModel) performDeleteComment(id string) tea.Cmd {
	return func() tea.Msg {
		err := m.Client.DeleteComment(id)
		if err != nil {
			return errorMsg(err)
		}
		return statusMsg("Comment deleted")
	}
}

func (m MainModel) fetchRandomPosts() tea.Cmd {
	return func() tea.Msg {
		posts, err := m.Client.GetRandomPosts(10)
		if err != nil {
			return errorMsg(err)
		}
		return feedMsg(posts)
	}
}

func (m MainModel) fetchRandomCommunities() tea.Cmd {
	return func() tea.Msg {
		communities, err := m.Client.GetRandomCommunities(10)
		if err != nil {
			return errorMsg(err)
		}
		return communitiesMsg(communities)
	}
}

func (m MainModel) fetchFollowedUsers() tea.Cmd {
	return func() tea.Msg {
		users, err := m.Client.GetFollowedUsers()
		if err != nil {
			return errorMsg(err)
		}
		return usersMsg(users)
	}
}

func (m MainModel) fetchJoinedCommunities() tea.Cmd {
	return func() tea.Msg {
		communities, err := m.Client.GetJoinedCommunities()
		if err != nil {
			return errorMsg(err)
		}
		return communitiesMsg(communities)
	}
}

func (m MainModel) performToggleFollow(userID string) tea.Cmd {
	return func() tea.Msg {
		err := m.Client.ToggleFollow(userID)
		if err != nil {
			return errorMsg(err)
		}
		return statusMsg("Toggle follow success")
	}
}

func (m MainModel) performUpdateUser() tea.Cmd {
	return func() tea.Msg {
		if m.Me == nil {
			return errorMsg(errors.New("not logged in"))
		}
		data := map[string]interface{}{
			"avatar": m.SettingsModel.Avatar.Value(),
		}
		err := m.Client.UpdateUser(m.Me.ID, data)
		if err != nil {
			return errorMsg(err)
		}
		return m.fetchMe()()
	}
}

func (m MainModel) performDeleteUser(id string) tea.Cmd {
	return func() tea.Msg {
		err := m.Client.DeleteUser(id)
		if err != nil {
			return errorMsg(err)
		}
		// Logout after deletion
		m.Client.SetToken("")
		m.Config.UpdateToken("")
		return statusMsg("Account deleted")
	}
}
