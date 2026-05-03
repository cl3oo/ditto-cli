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
	now := time.Date(2026, time.May, 3, 14, 5, 0, 0, time.UTC)
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

	rendered := m.renderCommentItem(m.FlattenedComments[0], true)

	if !strings.Contains(rendered, "• u/root") {
		t.Fatalf("expected root comment header, got:\n%s", rendered)
	}
	expectedTime := formatDetailTimestamp(now)
	if !strings.Contains(rendered, expectedTime) {
		t.Fatalf("expected exact timestamp metadata %q, got:\n%s", expectedTime, rendered)
	}
	if !strings.Contains(rendered, "1 reply") {
		t.Fatalf("expected reply count metadata, got:\n%s", rendered)
	}
	if !strings.Contains(rendered, "  Parent comment") {
		t.Fatalf("expected content prefix to stay visible, got:\n%s", rendered)
	}
}

func TestFormatDetailTimestamp(t *testing.T) {
	ts := time.Date(2026, time.May, 3, 14, 5, 0, 0, time.UTC)
	got := formatDetailTimestamp(ts)
	if got != ts.Local().Format("2006-01-02 15:04") {
		t.Fatalf("unexpected timestamp %q", got)
	}
}
