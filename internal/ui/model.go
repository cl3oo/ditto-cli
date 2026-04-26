package ui

import (
	"errors"
	"os/exec"
	"runtime"

	"github.com/charmbracelet/bubbletea"
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
)

type MainModel struct {
	State         State
	Client        *api.Client
	Error         error
	Width         int
	Height        int
	CommandBuffer string

	// Sub-models
	FeedModel       views.FeedModel
	LoginModel      views.LoginModel
	PostDetailModel views.PostDetailModel
	CommunityModel  views.CommunityModel
	CreatePostModel views.CreatePostModel
}

func NewMainModel(baseURL string) MainModel {
	return MainModel{
		State:           StateLoading,
		Client:          api.NewClient(baseURL),
		FeedModel:       views.NewFeedModel(),
		LoginModel:      views.NewLoginModel(),
		PostDetailModel: views.NewPostDetailModel(),
		CommunityModel:  views.NewCommunityModel(),
		CreatePostModel: views.NewCreatePostModel(),
	}
}

func (m MainModel) Init() tea.Cmd {
	return tea.Batch(
		m.fetchFeed(),
		m.LoginModel.Init(),
	)
}

func (m MainModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		if m.Error != nil {
			m.Error = nil
		}

		key := msg.String()

		// Handle command prefix
		if m.CommandBuffer != "" {
			m.CommandBuffer += key
			executed := true
			switch m.CommandBuffer {
			case ":q":
				return m, tea.Quit
			case ":L":
				m.State = StateLogin
			case ":F":
				m.State = StateFeed
			case ":C":
				m.State = StateLoading
				m.CommandBuffer = ""
				return m, m.fetchCommunities()
			case ":n":
				m.State = StateCreatePost
			case ":r":
				if m.State == StateFeed {
					m.State = StateLoading
					m.CommandBuffer = ""
					return m, m.fetchFeed()
				}
				if m.State == StateCommunities {
					m.State = StateLoading
					m.CommandBuffer = ""
					return m, m.fetchCommunities()
				}
			default:
				executed = false
			}

			if executed || len(m.CommandBuffer) > 2 {
				m.CommandBuffer = ""
			}
			if executed {
				return m, nil
			}
		}

		if key == ":" {
			m.CommandBuffer = ":"
			return m, nil
		}

		switch key {
		case "ctrl+c":
			return m, tea.Quit
		case "a": // Upvote stays as single key for quick interaction, or could be :a
			if m.State == StateFeed {
				if item, ok := m.FeedModel.List.SelectedItem().(views.PostItem); ok {
					return m, m.performVote(item.ID, 0, 1) // 0 = post
				}
			}
			if m.State == StatePostDetail {
				return m, m.performVote(m.PostDetailModel.Post.ID, 0, 1)
			}
		case "z": // Downvote
			if m.State == StateFeed {
				if item, ok := m.FeedModel.List.SelectedItem().(views.PostItem); ok {
					return m, m.performVote(item.ID, 0, -1)
				}
			}
			if m.State == StatePostDetail {
				return m, m.performVote(m.PostDetailModel.Post.ID, 0, -1)
			}
		case "r":
			if m.State == StateFeed {
				m.State = StateLoading
				return m, m.fetchFeed()
			}
			if m.State == StateCommunities {
				m.State = StateLoading
				return m, m.fetchCommunities()
			}
		case "enter":
			if m.State == StateLogin {
				if m.LoginModel.Focused == 2 {
					m.State = StateLoading
					return m, m.performLogin()
				}
				if m.LoginModel.Focused == 3 {
					return m, m.openBrowser("https://ditto.social/register")
				}
			}
			if m.State == StateCreatePost && m.CreatePostModel.Focused == 3 {
				m.State = StateLoading
				return m, m.performCreatePost()
			}
			if m.State == StateFeed {
				if item, ok := m.FeedModel.List.SelectedItem().(views.PostItem); ok {
					m.State = StateLoading
					return m, m.fetchPostDetail(item.ID)
				}
			}
		case "esc", "backspace":
			if m.State == StatePostDetail || m.State == StateCommunities || m.State == StateCreatePost {
				m.State = StateFeed
				return m, nil
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

	case postDetailMsg:
		m.State = StatePostDetail
		m.PostDetailModel.SetContent(msg.post, msg.comments)
		return m, nil

	case loginSuccessMsg:
		m.State = StateFeed
		m.Client.SetToken(string(msg))
		
		// Persist token
		cfg, _ := config.LoadConfig()
		_ = cfg.UpdateToken(string(msg))
		
		return m, m.fetchFeed()

	case errorMsg:
		if errors.Is(msg, api.ErrUnauthorized) {
			m.State = StateLogin
			m.Error = errors.New("Session expired, please login again")
			m.Client.SetToken("")
			cfg, _ := config.LoadConfig()
			_ = cfg.UpdateToken("")
			return m, nil
		}
		m.Error = msg
		// Don't change state, just show error on current screen
		return m, nil
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
	}

	return m, tea.Batch(cmds...)
}

type feedMsg []types.Post
type communitiesMsg []types.Community
type postDetailMsg struct {
	post     types.Post
	comments []types.Comment
}
type loginSuccessMsg string
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
		// For simplicity, just refresh after voting
		if m.State == StateFeed {
			return m.fetchFeed()()
		}
		return m.fetchPostDetail(id)()
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
