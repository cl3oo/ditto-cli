package views

import (
	"strings"
	"testing"
	"time"

	"github.com/rfcku/ditto-cli/internal/types"
)

func TestPostDetailAppendCommentsMergesDuplicateParents(t *testing.T) {
	now := time.Now()
	m := NewPostDetailModel()
	m.SetContent(types.Post{ID: "post-1", Title: "Test"}, []types.Comment{{
		ID:        "parent",
		Content:   "parent",
		Author:    types.UserMin{Username: "root"},
		CreatedAt: now,
		Children: []types.Comment{{
			ID:        "child-1",
			Content:   "first reply",
			Author:    types.UserMin{Username: "alpha"},
			CreatedAt: now,
		}},
	}})

	m.AppendComments([]types.Comment{{
		ID:        "parent",
		Content:   "parent refreshed",
		Author:    types.UserMin{Username: "root"},
		CreatedAt: now,
		Children: []types.Comment{{
			ID:        "child-2",
			Content:   "second reply",
			Author:    types.UserMin{Username: "beta"},
			CreatedAt: now,
		}},
	}})

	if got := len(m.Comments); got != 1 {
		t.Fatalf("expected 1 top-level comment after merge, got %d", got)
	}
	if got := len(m.Comments[0].Children); got != 2 {
		t.Fatalf("expected merged children, got %d", got)
	}
	if got := len(m.FlattenedComments); got != 3 {
		t.Fatalf("expected flattened parent + 2 children, got %d", got)
	}
	if got := m.FlattenedComments[2].ID; got != "child-2" {
		t.Fatalf("expected appended child to stay visible, got %q", got)
	}
}

func TestPostDetailRenderCommentItemShowsTreeAndWrapsContent(t *testing.T) {
	now := time.Now()
	m := NewPostDetailModel()
	m.Width = 38
	m.SetContent(types.Post{ID: "post-1", Title: "Test"}, []types.Comment{{
		ID:        "root",
		Content:   "Parent comment",
		Author:    types.UserMin{Username: "root"},
		CreatedAt: now,
		Children: []types.Comment{{
			ID:        "child",
			Content:   "This nested reply is intentionally long so it has to wrap across multiple lines to stay visible.",
			Author:    types.UserMin{Username: "child"},
			CreatedAt: now,
		}},
	}})

	rendered := m.renderCommentItem(m.FlattenedComments[1], true)

	if !strings.Contains(rendered, "└─ u/child") {
		t.Fatalf("expected tree branch in rendered comment, got:\n%s", rendered)
	}
	if !strings.Contains(rendered, "    This nested reply is") {
		t.Fatalf("expected wrapped content prefix to stay visible, got:\n%s", rendered)
	}
	if strings.Count(rendered, "\n") < 2 {
		t.Fatalf("expected wrapped multi-line comment, got:\n%s", rendered)
	}
}

func TestPostDetailViewShowsMediaPreview(t *testing.T) {
	m := NewPostDetailModel()
	m.SetSize(50, 20)
	m.SetContent(types.Post{ID: "post-1", Title: "Preview post"}, nil)
	m.SetMediaPreview("@@\n..")

	view := m.View()
	if !strings.Contains(view, "Preview:") {
		t.Fatalf("expected preview heading, got:\n%s", view)
	}
	if !strings.Contains(view, "@@") {
		t.Fatalf("expected preview content, got:\n%s", view)
	}
}
