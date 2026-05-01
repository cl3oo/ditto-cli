package ui

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/rfcku/ditto-cli/internal/api"
	"github.com/rfcku/ditto-cli/internal/config"
	"github.com/rfcku/ditto-cli/internal/types"
	"github.com/rfcku/ditto-cli/internal/ui/views"
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
	StateHelp
	StateConfirm
	StateCommandPalette
)

type ConfirmAction int

const (
	ConfirmNone ConfirmAction = iota
	ConfirmDeletePost
	ConfirmDeleteCommunity
	ConfirmDeleteComment
	ConfirmBanUser
	ConfirmModDeletePost
	ConfirmDeleteAccount
)

type ConfirmDialog struct {
	Action      ConfirmAction
	Previous    State
	Title       string
	Body        string
	TargetID    string
	TargetLabel string
	Extra       string
	Step        int
	Steps       int
}

type MainModel struct {
	State         State
	PreviousState State
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
	Wallet        *types.Wallet

	// Sub-models
	FeedModel          views.FeedModel
	LoginModel         views.LoginModel
	RegisterModel      views.RegisterModel
	PostDetailModel    views.PostDetailModel
	CommunityModel     views.CommunityModel
	CreatePostModel    views.CreatePostModel
	EditPostModel      views.CreatePostModel
	EditCommunityModel views.CreatePostModel
	SettingsModel      views.SettingsModel
	HelpModel          views.HelpModel
	PaletteModel       views.CommandPaletteModel
	ConfirmDialog      ConfirmDialog
}

func NewMainModel(cfg *config.Config) MainModel {
	m := MainModel{
		State:              StateLoading,
		Client:             api.NewClient(cfg.BaseURL),
		Config:             cfg,
		Theme:              NewTheme(cfg.Appearance),
		Keys:               NewKeyMap(cfg.Keys),
		FeedModel:          views.NewFeedModel(),
		LoginModel:         views.NewLoginModel(),
		RegisterModel:      views.NewRegisterModel(),
		PostDetailModel:    views.NewPostDetailModel(),
		CommunityModel:     views.NewCommunityModel(),
		CreatePostModel:    views.NewCreatePostModel(),
		EditPostModel:      views.NewCreatePostModel(),
		EditCommunityModel: views.NewCreatePostModel(),
		SettingsModel:      views.NewSettingsModel(),
		HelpModel:          views.NewHelpModel(),
		PaletteModel:       views.NewCommandPaletteModel(),
	}
	m.Client.SetToken(cfg.Token)
	m.FeedModel.SetTheme(m.Theme.Selected)
	m.CommunityModel.SetTheme(m.Theme.Selected)
	m.PaletteModel.SetTheme(m.Theme.Selected)
	m.PostDetailModel.SetTheme(m.Theme.Accent, m.Theme.Selected, m.Theme.Markdown)
	return m
}

func (m MainModel) Init() tea.Cmd {
	var cmds []tea.Cmd
	cmds = append(cmds, m.fetchFeed())
	cmds = append(cmds, m.LoginModel.Init())
	if m.Client.Token != "" {
		cmds = append(cmds, m.fetchMe())
		cmds = append(cmds, m.fetchWallet())
	}
	return tea.Batch(cmds...)
}

type tickMsg struct{}
type statusMsg string
type actionStatusMsg struct {
	status    string
	nextState State
	refresh   string
	postID    string
}

func (m MainModel) clearStatus() tea.Cmd {
	return tea.Tick(3*time.Second, func(t time.Time) tea.Msg {
		return tickMsg{}
	})
}

func (m MainModel) errorCmd(err error) tea.Cmd {
	return func() tea.Msg {
		return errorMsg(err)
	}
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
		switch m.State {
		case StateFeed:
			refreshCmd = m.fetchFeed()
		case StatePostDetail:
			refreshCmd = m.fetchPostDetail(m.PostDetailModel.Post.ID)
		}
		return m, tea.Batch(refreshCmd, m.clearStatus())

	case actionStatusMsg:
		m.StatusMessage = msg.status
		m.State = msg.nextState
		var refreshCmd tea.Cmd
		switch msg.refresh {
		case "feed":
			refreshCmd = m.fetchFeed()
		case "communities":
			refreshCmd = m.fetchCommunities()
		case "post":
			refreshCmd = m.fetchPostDetail(msg.postID)
		}
		return m, tea.Batch(refreshCmd, m.clearStatus())

	case tea.KeyMsg:
		if m.Error != nil {
			m.Error = nil
		}

		if m.State == StateConfirm {
			return m.handleConfirmKey(msg)
		}

		if m.State == StateCommandPalette {
			return m.handlePaletteKey(msg)
		}

		if key.Matches(msg, m.Keys.Palette) {
			m.PreviousState = m.State
			m.State = StateCommandPalette
			m.PaletteModel.SetSize(m.Width, m.Height)
			return m, nil
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
				case ":man":
					m.State = StateHelp
					m.HelpModel.SetSize(m.Width, m.Height-4)
					return m, nil
				case ":logout":
					m.Client.SetToken("")
					_ = m.Config.UpdateToken("")
					m.Me = nil
					m.Wallet = nil
					m.State = StateLogin
					m.StatusMessage = "Logged out"
					return m, m.clearStatus()
				case ":l", ":L", ":login":
					m.State = StateLogin
				case ":f", ":F", ":feed":
					m.State = StateLoading
					return m, m.fetchFeed()
				case ":c", ":C", ":communities":
					m.State = StateLoading
					return m, m.fetchCommunities()
				case ":n", ":new":
					m.State = StateCreatePost
				case ":r", ":refresh":
					prevState := m.State
					m.State = StateLoading
					switch prevState {
					case StateFeed:
						return m, m.fetchFeed()
					case StateCommunities:
						return m, m.fetchCommunities()
					case StatePostDetail:
						return m, m.fetchPostDetail(m.PostDetailModel.Post.ID)
					}
					// Default fallback if we don't know what to refresh
					return m, m.fetchFeed()
				case ":delete":
					if m.State == StatePostDetail {
						return m.openConfirm(ConfirmDialog{
							Action:      ConfirmDeletePost,
							Previous:    StatePostDetail,
							Title:       "Delete post?",
							Body:        "This permanently removes the current post.",
							TargetID:    m.PostDetailModel.Post.ID,
							TargetLabel: m.PostDetailModel.Post.Title,
							Step:        1,
							Steps:       1,
						}), nil
					}
					if m.State == StateCommunities {
						if item, ok := m.CommunityModel.List.SelectedItem().(views.CommunityItem); ok {
							return m.openConfirm(ConfirmDialog{
								Action:      ConfirmDeleteCommunity,
								Previous:    StateCommunities,
								Title:       "Delete community?",
								Body:        "This permanently removes the selected community.",
								TargetID:    item.ID,
								TargetLabel: item.Name,
								Step:        1,
								Steps:       1,
							}), nil
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
							m.EditCommunityModel.CommunityID.SetValue(item.Name)
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
					switch prevState {
					case StateFeed:
						return m, m.fetchRandomPosts()
					case StateCommunities:
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
				case ":ban":
					if len(parts) > 1 {
						username := parts[1]
						return m.openConfirm(ConfirmDialog{
							Action:      ConfirmBanUser,
							Previous:    m.State,
							Title:       "Ban user?",
							Body:        "This removes the user from the current community.",
							TargetID:    username,
							TargetLabel: username,
							Step:        1,
							Steps:       1,
						}), nil
					}
				case ":mod":
					if len(parts) > 2 && parts[1] == "add" {
						username := parts[2]
						m.State = StateLoading
						return m, m.performAddModerator(username)
					}
				case ":lock":
					if m.State == StatePostDetail {
						m.State = StateLoading
						return m, m.performLockPost(m.PostDetailModel.Post.ID, true)
					}
				case ":unlock":
					if m.State == StatePostDetail {
						m.State = StateLoading
						return m, m.performLockPost(m.PostDetailModel.Post.ID, false)
					}
				case ":mod-delete":
					if m.State == StatePostDetail && len(parts) > 1 {
						reason := strings.Join(parts[1:], " ")
						return m.openConfirm(ConfirmDialog{
							Action:      ConfirmModDeletePost,
							Previous:    StatePostDetail,
							Title:       "Moderator delete post?",
							Body:        "This removes the post using moderator powers.",
							TargetID:    m.PostDetailModel.Post.ID,
							TargetLabel: m.PostDetailModel.Post.Title,
							Extra:       reason,
							Step:        1,
							Steps:       1,
						}), nil
					}
				case ":award":
					if m.State == StatePostDetail && len(parts) > 1 {
						awardID := parts[1]
						m.State = StateLoading
						targetID := m.PostDetailModel.Post.ID
						targetType := types.TargetTypePost
						if m.PostDetailModel.SelectedIdx >= 0 {
							targetID = m.PostDetailModel.FlattenedComments[m.PostDetailModel.SelectedIdx].ID
							targetType = types.TargetTypeComment
						}
						return m, m.performGiveAward(targetID, targetType, awardID)
					}
				case ":report":
					if len(parts) > 1 {
						reason := strings.Join(parts[1:], " ")
						if m.State == StatePostDetail {
							m.State = StateLoading
							return m, m.performReport(m.PostDetailModel.Post.ID, 2, reason)
						}
						if m.State == StateCommunities {
							if item, ok := m.CommunityModel.List.SelectedItem().(views.CommunityItem); ok {
								m.State = StateLoading
								return m, m.performReport(item.ID, 1, reason)
							}
						}
					}
				case ":delete-comment":
					if m.State == StatePostDetail && len(parts) > 1 {
						commentID := parts[1]
						return m.openConfirm(ConfirmDialog{
							Action:      ConfirmDeleteComment,
							Previous:    StatePostDetail,
							Title:       "Delete comment?",
							Body:        "This permanently removes the selected comment.",
							TargetID:    commentID,
							TargetLabel: commentID,
							Step:        1,
							Steps:       1,
						}), nil
					}
				case ":settings":
					if m.Me != nil {
						m.SettingsModel.Avatar.SetValue(m.Me.Avatar)
					}
					m.State = StateProfileSettings
					return m, nil
				case ":delete-account":
					if m.Me != nil {
						return m.openConfirm(ConfirmDialog{
							Action:      ConfirmDeleteAccount,
							Previous:    StateProfileSettings,
							Title:       "Delete account?",
							Body:        "This is irreversible. You'll be logged out and your account will be removed.",
							TargetID:    m.Me.ID,
							TargetLabel: m.Me.Username,
							Step:        1,
							Steps:       2,
						}), nil
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
				case ":share":
					if m.State == StatePostDetail {
						url := fmt.Sprintf("%s/posts/%s", m.Client.BaseURL, m.PostDetailModel.Post.ID)
						return m, m.performShare(url)
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
		case msg.String() == "q":
			// 'q' is now the default Back key
			if m.State == StatePostDetail {
				if m.PostDetailModel.ShowCommentInput {
					m.PostDetailModel.SetShowCommentInput(false)
					return m, nil
				}
				m.State = StateFeed
				return m, nil
			}
			if m.State == StateCommunities || m.State == StateCreatePost || m.State == StateRegister || m.State == StateEditPost || m.State == StateEditCommunity || m.State == StateProfileSettings || m.State == StateHelp {
				m.State = StateFeed
				return m, nil
			}
			return m, nil
		case msg.String() == "r":
			if m.State == StatePostDetail && !m.PostDetailModel.ShowCommentInput {
				if m.PostDetailModel.SelectedIdx >= 0 && m.PostDetailModel.SelectedIdx < len(m.PostDetailModel.FlattenedComments) {
					author := m.PostDetailModel.FlattenedComments[m.PostDetailModel.SelectedIdx].Author.Username
					m.PostDetailModel.CommentInput.Placeholder = "Replying to u/" + author + "..."
				} else {
					m.PostDetailModel.CommentInput.Placeholder = "Write a comment..."
				}
				m.PostDetailModel.SetShowCommentInput(true)
				return m, nil
			}
		case msg.String() == "p":
			if m.State == StatePostDetail && !m.PostDetailModel.ShowCommentInput {
				var username string
				if m.PostDetailModel.SelectedIdx >= 0 && m.PostDetailModel.SelectedIdx < len(m.PostDetailModel.FlattenedComments) {
					username = m.PostDetailModel.FlattenedComments[m.PostDetailModel.SelectedIdx].Author.Username
				} else {
					username = m.PostDetailModel.Post.Author.Username
				}
				m.StatusMessage = "Viewing profile for u/" + username + " (Coming soon!)"
				return m, m.clearStatus()
			}
		case msg.String() == "s":
			if m.State == StatePostDetail && !m.PostDetailModel.ShowCommentInput {
				url := fmt.Sprintf("%s/posts/%s", m.Client.BaseURL, m.PostDetailModel.Post.ID)
				return m, m.performShare(url)
			}
		case msg.String() == "R":
			if m.State == StatePostDetail && !m.PostDetailModel.ShowCommentInput {
				var targetID string
				var targetType int
				if m.PostDetailModel.SelectedIdx >= 0 {
					targetID = m.PostDetailModel.FlattenedComments[m.PostDetailModel.SelectedIdx].ID
					targetType = types.TargetTypeComment
				} else {
					targetID = m.PostDetailModel.Post.ID
					targetType = types.TargetTypePost
				}
				m.State = StateLoading
				return m, m.performReport(targetID, targetType, "Reported from TUI")
			}
		case msg.String() == "L":
			if m.State == StatePostDetail && !m.PostDetailModel.ShowCommentInput {
				m.State = StateLoading
				return m, m.fetchMoreComments()
			}
		case key.Matches(msg, m.Keys.Upvote):
			if m.State == StateFeed {
				if item, ok := m.FeedModel.List.SelectedItem().(views.PostItem); ok {
					return m, m.performVote(item.ID, types.TargetTypePost, 1)
				}
			}
			if m.State == StatePostDetail && !m.PostDetailModel.ShowCommentInput {
				var targetID string
				var targetType int
				if m.PostDetailModel.SelectedIdx >= 0 {
					targetID = m.PostDetailModel.FlattenedComments[m.PostDetailModel.SelectedIdx].ID
					targetType = types.TargetTypeComment
				} else {
					targetID = m.PostDetailModel.Post.ID
					targetType = types.TargetTypePost
				}
				return m, m.performVote(targetID, targetType, 1)
			}
		case key.Matches(msg, m.Keys.Downvote):
			if m.State == StateFeed {
				if item, ok := m.FeedModel.List.SelectedItem().(views.PostItem); ok {
					return m, m.performVote(item.ID, types.TargetTypePost, -1)
				}
			}
			if m.State == StatePostDetail && !m.PostDetailModel.ShowCommentInput {
				var targetID string
				var targetType int
				if m.PostDetailModel.SelectedIdx >= 0 {
					targetID = m.PostDetailModel.FlattenedComments[m.PostDetailModel.SelectedIdx].ID
					targetType = types.TargetTypeComment
				} else {
					targetID = m.PostDetailModel.Post.ID
					targetType = types.TargetTypePost
				}
				return m, m.performVote(targetID, targetType, -1)
			}
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
					if m.LoginModel.LoggedIn {
						m.State = StateLoading
						return m, tea.Batch(m.fetchFeed(), m.fetchMe(), m.fetchWallet())
					}
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

	case searchMsg:
		m.State = StateCommunities
		m.CommunityModel.SetItems([]list.Item(msg))
		return m, nil

	case postDetailMsg:
		m.State = StatePostDetail
		m.PostDetailModel.SetContent(msg.post, msg.comments)
		m.PostDetailModel.SetShowCommentInput(false)
		return m, nil

	case appendCommentsMsg:
		m.State = StatePostDetail
		m.PostDetailModel.AppendComments(msg)
		return m, nil

	case loginSuccessMsg:
		token := string(msg)
		m.Client.SetToken(token)

		// Persist token
		if err := m.Config.UpdateToken(token); err != nil {
			return m, m.errorCmd(fmt.Errorf("failed to save config: %w", err))
		}

		m.LoginModel.LoggedIn = true
		m.LoginModel.SuccessToken = token
		m.LoginModel.Error = ""
		m.State = StateLogin // STAY on login screen

		displayToken := token
		if len(token) > 8 {
			displayToken = token[:8] + "..."
		}
		m.StatusMessage = fmt.Sprintf("Successfully logged in! Token: %s", displayToken)
		return m, m.clearStatus()

	case meMsg:
		m.Me = (*types.User)(msg)
		return m, nil

	case walletMsg:
		m.Wallet = (*types.Wallet)(msg)
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
			// If we are loading and get unauthorized, maybe only one request failed.
			// Don't wipe the token immediately during initial boot unless we are sure.
			if m.State != StateLoading {
				if m.Client.Token != "" {
					m.StatusMessage = "Session expired, please login again"
					m.Client.SetToken("")
					// DO NOT call m.Config.UpdateToken("") here to avoid wiping the file on transient errors
				}
				if m.State != StateLogin && m.State != StateRegister {
					m.State = StateLogin
				}
			} else {
				// If we are loading, just show error but keep the token for now
				m.StatusMessage = "Some profile data failed to load (Unauthorized)"
			}
			return m, m.clearStatus()
		}
		if m.State == StateLogin {
			m.LoginModel.Error = msg.Error()
		}
		if m.State == StateRegister {
			m.RegisterModel.Error = msg.Error()
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
	case StateHelp:
		m.HelpModel, cmd = m.HelpModel.Update(msg)
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
type searchMsg []list.Item
type postDetailMsg struct {
	post     types.Post
	comments []types.Comment
}
type appendCommentsMsg []types.Comment
type loginSuccessMsg string
type commentSuccessMsg string
type meMsg *types.User
type walletMsg *types.Wallet
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

func (m MainModel) fetchWallet() tea.Cmd {
	return func() tea.Msg {
		wallet, err := m.Client.GetWallet()
		if err != nil {
			return errorMsg(err)
		}
		return walletMsg(wallet)
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
		comments, err := m.Client.GetComments(id, 1, 50)
		if err != nil {
			// Still show post even if comments fail
			return postDetailMsg{post: *post, comments: []types.Comment{}}
		}
		return postDetailMsg{post: *post, comments: comments}
	}
}

func (m MainModel) fetchMoreComments() tea.Cmd {
	return func() tea.Msg {
		id := m.PostDetailModel.Post.ID
		page := m.PostDetailModel.Page + 1
		comments, err := m.Client.GetComments(id, page, 50)
		if err != nil {
			return errorMsg(err)
		}
		return appendCommentsMsg(comments)
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

func (m MainModel) performSearch(query string) tea.Cmd {
	return func() tea.Msg {
		communities, _ := m.Client.SearchCommunities(query)
		posts, _ := m.Client.SearchPosts(query)

		var items []list.Item
		for _, c := range communities {
			items = append(items, views.CommunityItem{Community: c})
		}
		for _, p := range posts {
			items = append(items, views.PostItem{Post: p})
		}

		if len(items) == 0 {
			return errorMsg(fmt.Errorf("no results found for '%s'", query))
		}

		return searchMsg(items)
	}
}

func (m MainModel) fetchCommunityDetail(communityID string) tea.Cmd {
	return func() tea.Msg {
		filter := map[string]interface{}{
			"target_id":   communityID,
			"target_type": types.TargetTypeCommunity,
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
		targetID := m.PostDetailModel.Post.ID
		targetType := types.TargetTypePost

		if m.PostDetailModel.SelectedIdx >= 0 && m.PostDetailModel.SelectedIdx < len(m.PostDetailModel.FlattenedComments) {
			targetID = m.PostDetailModel.FlattenedComments[m.PostDetailModel.SelectedIdx].ID
			targetType = types.TargetTypeComment
		}

		err := m.Client.CreateComment(
			targetID,
			targetType,
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

func (m MainModel) openConfirm(dialog ConfirmDialog) MainModel {
	m.ConfirmDialog = dialog
	m.State = StateConfirm
	return m
}

func (m MainModel) closeConfirm() MainModel {
	m.State = m.ConfirmDialog.Previous
	m.ConfirmDialog = ConfirmDialog{}
	return m
}

func (m MainModel) handleConfirmKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "esc", "q":
		return m.closeConfirm(), nil
	case "enter":
		if m.ConfirmDialog.Action == ConfirmDeleteAccount && m.ConfirmDialog.Step < m.ConfirmDialog.Steps {
			m.ConfirmDialog.Step++
			m.ConfirmDialog.Title = "Delete account, really?"
			m.ConfirmDialog.Body = "Second confirmation required. Press enter again only if you really want to permanently delete this account."
			return m, nil
		}
		return m.executeConfirmedAction()
	}

	return m, nil
}

func (m MainModel) handlePaletteKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	if msg.String() == "esc" {
		m.State = m.PreviousState
		return m, nil
	}

	if msg.String() == "enter" {
		item, ok := m.PaletteModel.List.SelectedItem().(views.CommandItem)
		if !ok {
			return m, nil
		}

		m.State = m.PreviousState // Close palette before executing

		switch item.Action {
		case views.ActionGoFeed:
			m.State = StateLoading
			return m, m.fetchFeed()
		case views.ActionGoCommunities:
			m.State = StateLoading
			return m, m.fetchCommunities()
		case views.ActionGoLogin:
			m.State = StateLogin
		case views.ActionGoRegister:
			m.State = StateRegister
		case views.ActionCreatePost:
			m.State = StateCreatePost
		case views.ActionGoSettings:
			m.State = StateProfileSettings
		case views.ActionGoHelp:
			m.State = StateHelp
		case views.ActionQuit:
			return m, tea.Quit
		}
		return m, nil
	}

	var cmd tea.Cmd
	m.PaletteModel, cmd = m.PaletteModel.Update(msg)
	return m, cmd
}

func (m MainModel) executeConfirmedAction() (tea.Model, tea.Cmd) {
	dialog := m.ConfirmDialog
	m.ConfirmDialog = ConfirmDialog{}
	m.State = dialog.Previous

	switch dialog.Action {
	case ConfirmDeletePost:
		return m, m.performDeletePost(dialog.TargetID)
	case ConfirmDeleteCommunity:
		return m, m.performDeleteCommunity(dialog.TargetID)
	case ConfirmDeleteComment:
		return m, m.performDeleteComment(dialog.TargetID, m.PostDetailModel.Post.ID)
	case ConfirmBanUser:
		return m, m.performBanUser(dialog.TargetID, dialog.Previous)
	case ConfirmModDeletePost:
		return m, m.performModDeletePost(dialog.TargetID, dialog.Extra, dialog.Previous)
	case ConfirmDeleteAccount:
		return m, m.performDeleteUser(dialog.TargetID)
	default:
		return m.closeConfirm(), nil
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
		return actionStatusMsg{status: "Post deleted", nextState: StateFeed, refresh: "feed"}
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
		return actionStatusMsg{status: "Community deleted", nextState: StateCommunities, refresh: "communities"}
	}
}

func (m MainModel) performDeleteComment(id, postID string) tea.Cmd {
	return func() tea.Msg {
		err := m.Client.DeleteComment(id)
		if err != nil {
			return errorMsg(err)
		}
		return actionStatusMsg{status: "Comment deleted", nextState: StatePostDetail, refresh: "post", postID: postID}
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

func (m MainModel) performBanUser(userID string, previous State) tea.Cmd {
	return func() tea.Msg {
		var communityID string
		if item, ok := m.CommunityModel.List.SelectedItem().(views.CommunityItem); ok {
			communityID = item.ID
		} else if previous == StatePostDetail {
			communityID = m.PostDetailModel.Post.CommunityID
		} else {
			return errorMsg(errors.New("no community context for ban"))
		}

		err := m.Client.BanUser(communityID, userID)
		if err != nil {
			return errorMsg(err)
		}
		refresh := "communities"
		postID := ""
		if previous == StatePostDetail {
			refresh = "post"
			postID = m.PostDetailModel.Post.ID
		}
		return actionStatusMsg{status: "User banned from community", nextState: previous, refresh: refresh, postID: postID}
	}
}

func (m MainModel) performAddModerator(userID string) tea.Cmd {
	return func() tea.Msg {
		var communityID string
		if item, ok := m.CommunityModel.List.SelectedItem().(views.CommunityItem); ok {
			communityID = item.ID
		} else if m.State == StatePostDetail {
			communityID = m.PostDetailModel.Post.CommunityID
		} else {
			return errorMsg(errors.New("no community context for mod add"))
		}

		err := m.Client.AddModerator(communityID, userID)
		if err != nil {
			return errorMsg(err)
		}
		return statusMsg("Moderator added")
	}
}

func (m MainModel) performLockPost(postID string, lock bool) tea.Cmd {
	return func() tea.Msg {
		err := m.Client.LockPost(postID, lock)
		if err != nil {
			return errorMsg(err)
		}
		return m.fetchPostDetail(postID)()
	}
}

func (m MainModel) performModDeletePost(postID, reason string, previous State) tea.Cmd {
	return func() tea.Msg {
		err := m.Client.ModDeletePost(postID, reason)
		if err != nil {
			return errorMsg(err)
		}
		nextState := previous
		refresh := ""
		if previous == StatePostDetail {
			nextState = StateFeed
			refresh = "feed"
		}
		return actionStatusMsg{status: "Post deleted by moderator", nextState: nextState, refresh: refresh}
	}
}

func (m MainModel) performGiveAward(targetID string, targetType int, awardID string) tea.Cmd {
	return func() tea.Msg {
		err := m.Client.GiveAward(targetID, targetType, awardID)
		if err != nil {
			return errorMsg(err)
		}
		return tea.Batch(
			m.fetchWallet(),             // Update balance
			m.fetchPostDetail(targetID), // Update award count
			func() tea.Msg { return statusMsg("Award given!") },
		)
	}
}

func (m MainModel) performReport(targetID string, targetType int, reason string) tea.Cmd {
	return func() tea.Msg {
		err := m.Client.ReportResource(targetID, targetType, reason)
		if err != nil {
			return errorMsg(err)
		}
		return statusMsg("Report submitted")
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
		if err := m.Config.UpdateToken(""); err != nil {
			return errorMsg(err)
		}
		return actionStatusMsg{status: "Account deleted", nextState: StateLogin}
	}
}

func (m MainModel) performShare(url string) tea.Cmd {
	return func() tea.Msg {
		if err := copyToClipboard(url); err != nil {
			return statusMsg("Link: " + url)
		}
		return statusMsg("Link copied to clipboard!")
	}
}
